package cluster

import (
	"errors"
	"fmt"
	"ivory/clients/storage"
	"ivory/core/service/vault"
	"ivory/features/node"
	"slices"
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

	spec, err := s.nodeService.KeeperDeploySpec(node.KeeperDeploySpecRequest{Plugin: r.ClusterOptions.Plugins.Keeper})
	if err != nil {
		return nil, err
	}

	cluster := Request{
		Name:    r.CommonConfig.Cluster,
		Nodes:   mapDeployNodeConfigs(r.Nodes),
		Options: r.ClusterOptions,
	}

	if err := s.validateNodeNames(cluster.Nodes); err != nil {
		return nil, err
	}
	if err := s.validateNodePorts(cluster.Nodes); err != nil {
		return nil, err
	}
	// NOTE: a vault and a pair of inline credentials are two answers to one
	// question - rather than pick one silently, say that only one was asked for
	if cluster.Vaults.SshKeyId != nil && (r.CommonConfig.SshUser != "" || r.CommonConfig.SshPass != "") {
		return nil, ErrSshCredentialsAmbiguous
	}
	if cluster.Vaults.KeeperId != nil && (r.CommonConfig.KeeperUser != "" || r.CommonConfig.KeeperPass != "") {
		return nil, ErrKeeperCredentialsAmbiguous
	}
	if cluster.Vaults.DatabaseId != nil && (r.CommonConfig.DbUser != "" || r.CommonConfig.DbPass != "") {
		return nil, ErrDatabaseCredentialsAmbiguous
	}
	if cluster.Vaults.SshKeyId == nil && (r.CommonConfig.SshUser == "" || r.CommonConfig.SshPass == "") {
		return nil, ErrSshCredentialsRequired
	}
	keeperUser, err := s.resolveLockedUsername(r.CommonConfig.KeeperUser, spec.KeeperUser)
	if err != nil {
		return nil, err
	}
	dbUser, err := s.resolveLockedUsername(r.CommonConfig.DbUser, spec.DbUser)
	if err != nil {
		return nil, err
	}
	// NOTE: an engine that is its own keeper is asked twice - pointing both at
	// one vault entry is the user's answer to give, never the deploy's
	if spec.KeeperCredentials && cluster.Vaults.KeeperId == nil && (keeperUser == "" || r.CommonConfig.KeeperPass == "") {
		return nil, ErrKeeperCredentialsRequired
	}
	if spec.DbCredentials && cluster.Vaults.DatabaseId == nil && (dbUser == "" || r.CommonConfig.DbPass == "") {
		return nil, ErrDatabaseCredentialsRequired
	}
	if err := s.nodeService.ValidateKeeperLockedCredentials(spec.KeeperUser, getVaultId(cluster.Vaults.KeeperId)); err != nil {
		return nil, err
	}
	if err := s.nodeService.ValidateKeeperLockedCredentials(spec.DbUser, getVaultId(cluster.Vaults.DatabaseId)); err != nil {
		return nil, err
	}
	if _, e := s.Get(cluster.Name); e == nil {
		return nil, ErrClusterNameTaken
	} else if !errors.Is(e, storage.ErrNotFound) {
		return nil, e
	}

	logs := &deployLogs{}
	logs.send("system", fmt.Sprintf("deploying %d node(s)", len(r.Nodes)))

	// NOTE: only the vaults this call created are rolled back when nothing
	// deploys - one the user picked is theirs and outlives the attempt
	created := make([]uuid.UUID, 0, 3)

	// 2. Handle SSH Key
	if cluster.Vaults.SshKeyId == nil {
		logs.send("system", "generating ssh key and saving it to vault")
		sshVault := vault.Vault{Type: vault.SSH_KEY, Username: r.CommonConfig.SshUser}
		id, v, err := s.vaultService.Create(sshVault)
		if err != nil {
			return nil, err
		}
		cluster.Vaults.SshKeyId = id
		created = append(created, *id)
		if v.Metadata == nil {
			return nil, ErrSshKeyVaultMissingMetadata
		}
	}

	// 3. Handle Keeper Password
	if spec.KeeperCredentials && cluster.Vaults.KeeperId == nil {
		logs.send("system", "saving keeper credentials to vault")
		keeperVault := vault.Vault{Type: vault.KEEPER_PASSWORD, Username: keeperUser, Secret: r.CommonConfig.KeeperPass}
		id, _, err := s.vaultService.Create(keeperVault)
		if err != nil {
			return nil, err
		}
		cluster.Vaults.KeeperId = id
		created = append(created, *id)
	}

	// 4. Handle DB Password
	if spec.DbCredentials && cluster.Vaults.DatabaseId == nil {
		logs.send("system", "saving database credentials to vault")
		dbVault := vault.Vault{Type: vault.DATABASE_PASSWORD, Username: dbUser, Secret: r.CommonConfig.DbPass}
		id, _, err := s.vaultService.Create(dbVault)
		if err != nil {
			return nil, err
		}
		cluster.Vaults.DatabaseId = id
		created = append(created, *id)
	}

	// 5. Deploy every node by running its own command; the cluster owns only
	// the batch orchestration (parallel/sequential, aggregated timestamped
	// logs) while node executes the individual deployment.
	var up atomic.Int32
	if r.Parallel {
		var wg sync.WaitGroup
		for _, deployNode := range r.Nodes {
			wg.Add(1)
			go func(dn DeployNode) {
				defer wg.Done()
				if s.deployNode(cluster, dn, logs.send) {
					up.Add(1)
				}
			}(deployNode)
		}
		wg.Wait()
	} else {
		for _, deployNode := range r.Nodes {
			if s.deployNode(cluster, deployNode, logs.send) {
				up.Add(1)
			}
		}
	}
	deployed := int(up.Load()) == len(r.Nodes)

	// 6. Register the cluster only once something is actually running under
	// its name. Registering first left a cluster of nodes that do not exist
	// behind every failed attempt, along with the vaults made to reach them.
	if up.Load() == 0 {
		logs.send("system", "no node deployed: the cluster was not registered")
		s.rollbackVaults(created, logs.send)
		return logs.list(), nil
	}
	logs.send("system", "updating cluster configuration")
	if _, err = s.Update(cluster); err != nil {
		return nil, err
	}

	// 7. Post-scripts run after every node is up, in node order - a script
	// that needs the whole cluster running (etcd's auth enable, mongo's
	// rs.initiate) therefore belongs on the last node.
	if deployed {
		s.postDeploy(cluster, r.Nodes, logs.send)
	} else if slices.ContainsFunc(r.Nodes, func(n DeployNode) bool { return n.PostScript != "" }) {
		logs.send("system", "skipping post-deploy initialization: not every node deployed successfully")
	}

	return logs.list(), nil
}

// validateNodePorts requires every endpoint a deploy needs to be stated: a
// keeper port is never taken from the database port, however alike the two
// look on engines whose keeper endpoint is the database itself.
func (s *Service) validateNodePorts(nodes []NodeConfig) error {
	for _, n := range nodes {
		if !isPortProvided(n.KeeperPort) || !isPortProvided(n.DbPort) || !isPortProvided(n.SshPort) {
			return fmt.Errorf("%w: %s", ErrClusterNodePortsNotProvided, n.Name)
		}
	}
	return nil
}

// rollbackVaults removes the vaults this deploy created for itself. A failure
// to remove one is logged rather than returned: the deploy has already
// finished, and a stray vault is not worth losing its logs over.
func (s *Service) rollbackVaults(ids []uuid.UUID, logsSend func(ctx string, msg string)) {
	for _, id := range ids {
		if err := s.vaultService.Delete(id); err != nil {
			logsSend("system", fmt.Sprintf("failed to remove the vault created for this deploy: %v", err))
		}
	}
}

// postDeploy runs every node's post-script and reports, at the end, how many
// of them failed. A failure does not abort the batch - the nodes are already
// up - but it has to be stated: a post script is what turns running processes
// into an initialized cluster (etcd's auth enable, mongo's rs.initiate), so a
// deploy whose scripts all failed otherwise reads as a success.
func (s *Service) postDeploy(cluster Request, nodes []DeployNode, logsSend func(ctx string, msg string)) {
	failed := 0
	total := 0
	for _, n := range nodes {
		if n.PostScript == "" {
			continue
		}
		total++
		nodeKey := s.getNodeKey(n.Host, n.KeeperPort)
		logsSend(nodeKey, "running post-deploy initialization")
		logs, err := s.nodeService.KeeperPostDeploy(s.mapDeployRequest(cluster, n))
		for _, log := range logs {
			logsSend(nodeKey, log)
		}
		if err != nil {
			failed++
			logsSend(nodeKey, "post-deploy initialization did not complete")
			continue
		}
		logsSend(nodeKey, "post-deploy initialization finished")
	}
	if failed > 0 {
		logsSend("system", fmt.Sprintf(
			"post-deploy initialization failed on %d of %d node(s): the cluster is registered and running, but not initialized",
			failed, total))
	}
}

// deployNode deploys one node by delegating to node's KeeperDeployUp; it only
// adds the cluster-level concerns of resolving the node's vault-backed
// connection and aggregating timestamped logs for the batch.
func (s *Service) deployNode(cluster Request, n DeployNode, logsSend func(ctx string, msg string)) bool {
	nodeKey := s.getNodeKey(n.Host, n.KeeperPort)

	if n.Host == "" {
		logsSend(nodeKey, "deploy failed: host not provided for node")
		return false
	}
	if cluster.Vaults.SshKeyId == nil {
		logsSend(nodeKey, "deploy failed: ssh key id not provided for node")
		return false
	}

	logsSend(nodeKey, "deploy started")
	// NOTE: even if connection was closed we do not want to stop deployment
	res, err := s.nodeService.KeeperDeployUp(s.mapDeployRequest(cluster, n))
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

// mapDeployRequest builds node's flat deploy request for one node; both the
// deploy and post-deploy steps run against the identical scope, so they share
// one mapper rather than assembling it twice.
func (s *Service) mapDeployRequest(cluster Request, n DeployNode) node.KeeperDeployRequest {
	return node.KeeperDeployRequest{
		Plugin:     cluster.Options.Plugins.Keeper,
		Cluster:    cluster.Name,
		Name:       n.Name,
		KeeperPort: getPortValue(n.KeeperPort),
		DbPort:     getPortValue(n.DbPort),
		Command:    n.Command,
		PostScript: n.PostScript,
		Connection: s.getNodeConnection(cluster, n),
		Vaults:     s.getNodeVaults(cluster),
	}
}

func getPortValue(port *int) int {
	if port == nil {
		return 0
	}
	return *port
}

func isPortProvided(port *int) bool {
	return port != nil && *port > 0
}

// getVaultId converts to node's Vaults convention, a value type instead of a
// pointer.
func getVaultId(id *uuid.UUID) uuid.UUID {
	if id == nil {
		return uuid.Nil
	}
	return *id
}

// resolveLockedUsername applies an engine-required username (spilo's postgres,
// etcd's root): changing it is rejected rather than silently overridden.
func (s *Service) resolveLockedUsername(typed string, locked string) (string, error) {
	if locked == "" {
		return typed, nil
	}
	if typed != "" && typed != locked {
		return "", fmt.Errorf("username %q is not allowed: the keeper plugin locks it to %q", typed, locked)
	}
	return locked, nil
}

// getNodeConnection resolves the SSH connection for a node from the cluster's
// own SSH vault; it runs only after the SSH vault is guaranteed non-nil
// (either provided or freshly generated in Deploy's "Handle SSH Key" step)
// and after validateNodePorts, so the ssh port is the node's own rather than
// an assumed 22.
func (s *Service) getNodeConnection(cluster Request, n DeployNode) node.PlatformVaultConnection {
	return node.PlatformVaultConnection{Host: n.Host, Port: getPortValue(n.SshPort), VaultId: *cluster.Vaults.SshKeyId}
}

// getNodeVaults resolves the vault ids to pass to a node keeper-deploy action.
// A plugin that needs credentials it was not given fails with an unresolved
// {{keeperUser}}/{{dbUser}} placeholder error instead.
func (s *Service) getNodeVaults(cluster Request) node.Vaults {
	return node.Vaults{
		KeeperId:   getVaultId(cluster.Vaults.KeeperId),
		DatabaseId: getVaultId(cluster.Vaults.DatabaseId),
		SshKeyId:   *cluster.Vaults.SshKeyId,
	}
}
