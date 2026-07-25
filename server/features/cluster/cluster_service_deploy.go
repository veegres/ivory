package cluster

import (
	"errors"
	"fmt"
	"ivory/clients/storage"
	"ivory/core/service/vault"
	"ivory/features/node"
	"sync"
	"time"

	"github.com/google/uuid"
)

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

	cluster := Request{
		Name:    r.CommonConfig.Cluster,
		Nodes:   mapPlanNodeConfigs(plan.Nodes),
		Options: r.ClusterOptions,
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
	if _, e := s.Get(cluster.Name); !errors.Is(e, storage.ErrNotFound) {
		return nil, ErrClusterNameTaken
	}

	var mu sync.Mutex
	logs := make([]string, 0)
	logsSend := func(ctx string, msg string) {
		mu.Lock()
		logs = append(logs, fmt.Sprintf("%s | %s | %s", time.Now().Format("2006-01-02 15:04:05"), ctx, msg))
		mu.Unlock()
	}
	for _, warning := range plan.Warnings {
		logsSend("system", fmt.Sprintf("warning: %s", warning))
	}

	// 2. Handle SSH Key
	if cluster.Vaults.SshKeyId == nil {
		logsSend("system", "generating ssh key and saving it to vault")
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
		logsSend("system", "saving database credentials to vault")
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
	logsSend("system", "updating cluster configuration")
	_, err = s.Update(cluster)
	if err != nil {
		return nil, err
	}

	// 5. Deploy every node, calling into node's keeper-deploy business logic
	// for each; the cluster owns the batch orchestration (parallel/sequential,
	// aggregated timestamped logs) while node resolves and executes the
	// individual container deploy.
	deployed := true
	if r.Parallel {
		var wg sync.WaitGroup
		for _, planNode := range plan.Nodes {
			wg.Add(1)
			go func(pn node.KeeperDeployPlanNode) {
				defer wg.Done()
				if !s.deployNode(cluster, r.Values, plan, pn, logsSend) {
					mu.Lock()
					deployed = false
					mu.Unlock()
				}
			}(planNode)
		}
		wg.Wait()
	} else {
		for _, planNode := range plan.Nodes {
			if !s.deployNode(cluster, r.Values, plan, planNode, logsSend) {
				deployed = false
			}
		}
	}

	// 6. Post-deploy initialization runs once for the whole cluster (e.g.
	// enabling authentication), not once per node.
	if deployed && len(plan.PostDeploy) > 0 {
		pn := plan.Nodes[0]
		nodeKey := s.getNodeKey(pn.Host, &pn.KeeperPort)
		logsSend(nodeKey, "running post-deploy initialization")
		for _, log := range s.nodeService.KeeperPostDeploy(node.KeeperPostDeployRequest{
			Cluster:       cluster.Name,
			RequestValues: r.Values,
			PlanValues:    plan.Values,
			PostDeploy:    plan.PostDeploy,
			Node:          pn,
			Connection:    node.PlatformVaultConnection{Host: pn.Host, Port: pn.SshPort, VaultId: *cluster.Vaults.SshKeyId},
			Vaults:        node.Vaults{DatabaseId: s.getDatabaseId(cluster), SshKeyId: *cluster.Vaults.SshKeyId},
		}) {
			logsSend(nodeKey, log)
		}
		logsSend(nodeKey, "post-deploy initialization finished")
	} else if len(plan.PostDeploy) > 0 {
		logsSend("system", "skipping post-deploy initialization: not every node deployed successfully")
	}

	return logs, nil
}

// deployNode deploys one node of a plan already resolved for the whole
// cluster by delegating to node's KeeperDeployUp; it only adds the
// cluster-level concerns of resolving the node's vault-backed connection and
// aggregating timestamped logs for the batch.
func (s *Service) deployNode(cluster Request, requestValues map[string]string, plan *node.KeeperDeployPlanResponse, pn node.KeeperDeployPlanNode, logsSend func(ctx string, msg string)) bool {
	nodeKey := s.getNodeKey(pn.Host, &pn.KeeperPort)

	if pn.Host == "" {
		logsSend(nodeKey, "deployment failed: host not provided for node")
		return false
	}
	if cluster.Vaults.SshKeyId == nil {
		logsSend(nodeKey, "deployment failed: ssh key id not provided for node")
		return false
	}

	// NOTE: even if connection was closed we do not want to stop deployment
	res, err := s.nodeService.KeeperDeployUp(node.KeeperDeployUpRequest{
		Cluster:       cluster.Name,
		Image:         plan.Image,
		PlanValues:    plan.Values,
		RequestValues: requestValues,
		Node:          pn,
		Connection: node.PlatformVaultConnection{
			Host:    pn.Host,
			Port:    pn.SshPort,
			VaultId: *cluster.Vaults.SshKeyId,
		},
		Vaults: node.Vaults{
			// NOTE: the database vault is optional for keeper plugins that
			// consume no database credentials; a plugin that does need them
			// fails with an unresolved {{dbUser}}/{{dbPass}} placeholder error
			DatabaseId: s.getDatabaseId(cluster),
			SshKeyId:   *cluster.Vaults.SshKeyId,
		},
	})
	if err != nil {
		logsSend(nodeKey, fmt.Sprintf("deployment failed: %v", err))
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

func mapPlanNodeConfigs(planNodes []node.KeeperDeployPlanNode) []NodeConfig {
	nodes := make([]NodeConfig, 0, len(planNodes))
	for _, pn := range planNodes {
		sshPort, keeperPort, dbPort := pn.SshPort, pn.KeeperPort, pn.DbPort
		nodes = append(nodes, NodeConfig{
			Host:       pn.Host,
			SshPort:    &sshPort,
			KeeperPort: &keeperPort,
			DbPort:     &dbPort,
		})
	}
	return nodes
}
