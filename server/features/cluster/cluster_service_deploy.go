package cluster

import (
	"errors"
	"fmt"
	"ivory/clients/storage"
	"ivory/core/service/vault"
	"ivory/features/node"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/google/uuid"
)

// deployLogs collects timestamped log lines from a cluster deploy into a
// single ordered slice; Deploy's parallel mode runs one goroutine per node,
// so appends are mutex-guarded.
type deployLogs struct {
	mu   sync.Mutex
	logs []string
}

func (l *deployLogs) send(ctx string, msg string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.logs = append(l.logs, fmt.Sprintf("%s | %s | %s", time.Now().Format("2006-01-02 15:04:05"), ctx, msg))
}

func (l *deployLogs) list() []string {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.logs
}

func (s *Service) Deploy(r DeployRequest) ([]string, error) {
	// 1. Validation & Preparation
	if r.CommonConfig.Cluster == "" {
		return nil, ErrClusterNameNotProvided
	}
	if len(r.Nodes) == 0 {
		return nil, ErrClusterNodesNotProvided
	}

	plan, err := s.planDeploy(r.CommonConfig.Cluster, r.ClusterOptions, r.SingleHost, r.Image, r.Values, r.Nodes)
	if err != nil {
		return nil, err
	}
	if plan.Image == "" {
		return nil, ErrClusterImageNotProvided
	}
	if len(plan.Warnings) > 0 {
		return nil, fmt.Errorf("%w: %s", ErrClusterDeployPlanHasWarnings, strings.Join(plan.Warnings, "; "))
	}

	clusterOptions := r.ClusterOptions
	clusterOptions.SingleHost = r.SingleHost
	cluster := Request{
		Name:    r.CommonConfig.Cluster,
		Nodes:   mapPlanNodeConfigs(plan.Nodes),
		Options: clusterOptions,
	}

	if cluster.Vaults.SshKeyId == nil && (r.CommonConfig.SshUser == "" || r.CommonConfig.SshPass == "") {
		return nil, ErrSshCredentialsRequired
	}
	// NOTE: an engine-required username (spilo's postgres, etcd's root) is
	// locked: changing it is rejected instead of silently overridden, an
	// empty one is prefilled
	dbUser := r.CommonConfig.DbUser
	requiredUser, dbCredentials := plan.Fields.Defaults[string(node.VarDbUser)]
	if requiredUser != "" {
		if dbUser != "" && dbUser != requiredUser {
			return nil, fmt.Errorf("database username %q is not allowed: the keeper plugin locks it to %q", dbUser, requiredUser)
		}
		dbUser = requiredUser
	}
	if dbCredentials && cluster.Vaults.DatabaseId == nil && (dbUser == "" || r.CommonConfig.DbPass == "") {
		return nil, ErrDatabaseCredentialsRequired
	}
	if err := s.nodeService.ValidateKeeperLockedCredentials(requiredUser, s.getDatabaseId(cluster)); err != nil {
		return nil, err
	}
	if _, e := s.Get(cluster.Name); e == nil {
		return nil, ErrClusterNameTaken
	} else if !errors.Is(e, storage.ErrNotFound) {
		return nil, e
	}

	// NOTE: plan.Warnings is already guaranteed empty here - Deploy refuses to
	// proceed with a warned plan above, so there is nothing left to log.
	logs := &deployLogs{}
	logs.send("system", fmt.Sprintf("deployment plan resolved: image %s, %d node(s)", plan.Image, len(plan.Nodes)))

	// 2. Handle SSH Key
	if cluster.Vaults.SshKeyId == nil {
		logs.send("system", "generating ssh key and saving it to vault")
		sshVault := vault.Vault{Type: vault.SSH_KEY, Username: r.CommonConfig.SshUser}
		id, v, err := s.vaultService.Create(sshVault)
		if err != nil {
			return nil, err
		}
		cluster.Vaults.SshKeyId = id
		if v.Metadata == nil {
			return nil, ErrSshKeyVaultMissingMetadata
		}
	}

	// 3. Handle DB Password
	if dbCredentials && cluster.Vaults.DatabaseId == nil {
		logs.send("system", "saving database credentials to vault")
		dbVault := vault.Vault{Type: vault.DATABASE_PASSWORD, Username: dbUser, Secret: r.CommonConfig.DbPass}
		id, _, err := s.vaultService.Create(dbVault)
		if err != nil {
			return nil, err
		}
		cluster.Vaults.DatabaseId = id
	}
	// NOTE: when the plugin has no separate keeper port its keeper endpoint
	// is the database itself, so the same credentials authenticate both
	if _, hasKeeperEndpoint := plan.Fields.Defaults[string(node.VarKeeperPort)]; !hasKeeperEndpoint && cluster.Vaults.KeeperId == nil {
		cluster.Vaults.KeeperId = cluster.Vaults.DatabaseId
	}

	// 4. Update Cluster
	logs.send("system", "updating cluster configuration")
	_, err = s.Update(cluster)
	if err != nil {
		return nil, err
	}

	// 5. Deploy every node, calling into node's keeper-deploy business logic
	// for each; the cluster owns the batch orchestration (parallel/sequential,
	// aggregated timestamped logs) while node resolves and executes the
	// individual container deploy.
	var deployed atomic.Bool
	deployed.Store(true)
	if r.Parallel {
		var wg sync.WaitGroup
		for _, planNode := range plan.Nodes {
			wg.Add(1)
			go func(pn node.KeeperDeployPlanNode) {
				defer wg.Done()
				if !s.deployNode(cluster, r.Values, plan, pn, logs.send) {
					deployed.Store(false)
				}
			}(planNode)
		}
		wg.Wait()
	} else {
		for _, planNode := range plan.Nodes {
			if !s.deployNode(cluster, r.Values, plan, planNode, logs.send) {
				deployed.Store(false)
			}
		}
	}

	// 6. Post-deploy initialization runs once for the whole cluster (e.g.
	// enabling authentication), not once per node.
	if deployed.Load() && plan.PostScript != "" {
		pn := plan.Nodes[0]
		nodeKey := s.getNodeKey(pn.Host, &pn.KeeperPort)
		logs.send(nodeKey, "running post-deploy initialization")
		for _, log := range s.nodeService.KeeperPostDeploy(node.KeeperPostDeployRequest{
			KeeperDeployExecRequest: node.KeeperDeployExecRequest{
				Cluster:       cluster.Name,
				RequestValues: r.Values,
				PlanValues:    plan.Values,
				Node:          pn,
				Connection:    s.getNodeConnection(cluster, pn),
				Vaults:        s.getNodeVaults(cluster),
			},
			PostScript: plan.PostScript,
		}) {
			logs.send(nodeKey, log)
		}
		logs.send(nodeKey, "post-deploy initialization finished")
	} else if plan.PostScript != "" {
		logs.send("system", "skipping post-deploy initialization: not every node deployed successfully")
	}

	return logs.list(), nil
}

// deployNode deploys one node of a plan already resolved for the whole
// cluster by delegating to node's KeeperDeployUp; it only adds the
// cluster-level concerns of resolving the node's vault-backed connection and
// aggregating timestamped logs for the batch.
func (s *Service) deployNode(cluster Request, requestValues map[string]string, plan *node.KeeperDeployPlanResponse, pn node.KeeperDeployPlanNode, logsSend func(ctx string, msg string)) bool {
	nodeKey := s.getNodeKey(pn.Host, &pn.KeeperPort)

	if pn.Host == "" {
		logsSend(nodeKey, "deploy failed: host not provided for node")
		return false
	}
	if cluster.Vaults.SshKeyId == nil {
		logsSend(nodeKey, "deploy failed: ssh key id not provided for node")
		return false
	}

	logsSend(nodeKey, "deploy started")
	// NOTE: even if connection was closed we do not want to stop deployment
	res, err := s.nodeService.KeeperDeployUp(node.KeeperDeployUpRequest{
		KeeperDeployExecRequest: node.KeeperDeployExecRequest{
			Cluster:       cluster.Name,
			PlanValues:    plan.Values,
			RequestValues: requestValues,
			Node:          pn,
			Connection:    s.getNodeConnection(cluster, pn),
			Vaults:        s.getNodeVaults(cluster),
		},
		Image: plan.Image,
	})
	if err != nil {
		logsSend(nodeKey, fmt.Sprintf("deploy failed: %v", err))
		return false
	}
	for _, log := range res {
		logsSend(nodeKey, log)
	}
	logsSend(nodeKey, "deploy finished")
	return true
}

// getDatabaseId returns the cluster's database vault id, or uuid.Nil when
// none is linked, matching node's Vaults convention (a value type instead of
// a pointer).
func (s *Service) getDatabaseId(cluster Request) uuid.UUID {
	if cluster.Vaults.DatabaseId == nil {
		return uuid.Nil
	}
	return *cluster.Vaults.DatabaseId
}

// getNodeConnection resolves the SSH connection for a planned node from the
// cluster's own SSH vault; shared by deployNode and Deploy's post-deploy
// step, both of which run only after the SSH vault is guaranteed non-nil
// (either provided or freshly generated in Deploy's "Handle SSH Key" step).
func (s *Service) getNodeConnection(cluster Request, pn node.KeeperDeployPlanNode) node.PlatformVaultConnection {
	return node.PlatformVaultConnection{Host: pn.Host, Port: pn.SshPort, VaultId: *cluster.Vaults.SshKeyId}
}

// getNodeVaults resolves the database/SSH vault ids to pass to a node
// keeper-deploy action; shared by deployNode and Deploy's post-deploy step.
// The database vault is optional for keeper plugins that consume no
// database credentials; a plugin that does need them fails with an
// unresolved {{dbUser}}/{{dbPass}} placeholder error instead.
func (s *Service) getNodeVaults(cluster Request) node.Vaults {
	return node.Vaults{DatabaseId: s.getDatabaseId(cluster), SshKeyId: *cluster.Vaults.SshKeyId}
}

func (s *Service) planDeploy(name string, options Options, singleHost bool, image string, values map[string]string, nodes []DeployNode) (*node.KeeperDeployPlanResponse, error) {
	planNodes := make([]node.KeeperDeployPlanNodeRequest, 0, len(nodes))
	for _, n := range nodes {
		planNodes = append(planNodes, node.KeeperDeployPlanNodeRequest{
			Host:       n.Host,
			SshPort:    n.SshPort,
			KeeperPort: n.KeeperPort,
			DbPort:     n.DbPort,
			Options:    n.Options,
		})
	}
	return s.nodeService.KeeperDeployPlan(node.KeeperDeployPlanRequest{
		Plugin:     options.Plugins.Keeper,
		Cluster:    name,
		SingleHost: singleHost,
		Image:      image,
		Values:     values,
		Nodes:      planNodes,
	})
}
