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

// Deploy reports whether the batch completed as well as its logs: complete
// means every node came up and every post-script ran, which is what separates a
// deploy the caller can answer 200 to from one it has to flag. An error is
// what Ivory itself rejected before anything started; a container the
// operator's own command killed is not one, and the logs are the report.
func (s *Service) Deploy(r DeployRequest) ([]string, bool, error) {
	// 1. Validation & Preparation
	if r.CommonConfig.Cluster == "" {
		return nil, false, ErrClusterNameNotProvided
	}
	if len(r.Nodes) == 0 {
		return nil, false, ErrClusterNodesNotProvided
	}

	if err := s.nodeService.ValidateKeeperPlugin(r.ClusterOptions.Plugins.Keeper); err != nil {
		return nil, false, err
	}

	cluster := Request{
		Name:    r.CommonConfig.Cluster,
		Nodes:   mapDeployNodeConfigs(r.Nodes),
		Options: r.ClusterOptions,
	}

	if err := s.validateNodeNames(cluster.Nodes); err != nil {
		return nil, false, err
	}
	if err := s.validateNodePorts(cluster.Nodes); err != nil {
		return nil, false, err
	}
	// NOTE: a vault and a pair of inline credentials are two answers to one
	// question - rather than pick one silently, say that only one was asked for
	if cluster.Vaults.SshKeyId != nil && (r.CommonConfig.SshUser != "" || r.CommonConfig.SshPass != "") {
		return nil, false, ErrSshCredentialsAmbiguous
	}
	if cluster.Vaults.KeeperId != nil && (r.CommonConfig.KeeperUser != "" || r.CommonConfig.KeeperPass != "") {
		return nil, false, ErrKeeperCredentialsAmbiguous
	}
	if cluster.Vaults.DatabaseId != nil && (r.CommonConfig.DbUser != "" || r.CommonConfig.DbPass != "") {
		return nil, false, ErrDatabaseCredentialsAmbiguous
	}
	if cluster.Vaults.SshKeyId == nil && (r.CommonConfig.SshUser == "" || r.CommonConfig.SshPass == "") {
		return nil, false, ErrSshCredentialsRequired
	}
	// NOTE: whether a deployment has keeper or database credentials at all is
	// the user's answer, not the engine's - the template names a default
	// username and the deploy screen offers to switch the pair off. What is
	// left to check is that the answer is whole: half a pair would be stored
	// as a vault entry that authenticates nothing.
	if isHalfCredential(r.CommonConfig.KeeperUser, r.CommonConfig.KeeperPass) {
		return nil, false, ErrKeeperCredentialsIncomplete
	}
	if isHalfCredential(r.CommonConfig.DbUser, r.CommonConfig.DbPass) {
		return nil, false, ErrDatabaseCredentialsIncomplete
	}
	if _, e := s.Get(cluster.Name); e == nil {
		return nil, false, ErrClusterNameTaken
	} else if !errors.Is(e, storage.ErrNotFound) {
		return nil, false, e
	}

	logs := &deployLogs{}
	logs.send("system", fmt.Sprintf("deploying %d node(s)", len(r.Nodes)))

	// NOTE: only the vaults this call created can be rolled back at all - one
	// the user picked is theirs and outlives the attempt either way
	created := make([]uuid.UUID, 0, 3)

	// 2. Handle SSH Key
	if cluster.Vaults.SshKeyId == nil {
		logs.send("system", "generating ssh key and saving it to vault")
		sshVault := vault.Vault{Type: vault.SSH_KEY, Username: r.CommonConfig.SshUser}
		id, v, err := s.vaultService.Create(sshVault)
		if err != nil {
			s.rollbackVaults(created, logs.send)
			return nil, false, err
		}
		cluster.Vaults.SshKeyId = id
		created = append(created, *id)
		if v.Metadata == nil {
			s.rollbackVaults(created, logs.send)
			return nil, false, ErrSshKeyVaultMissingMetadata
		}
		if err := s.authorizeSshKey(r.Nodes, r.CommonConfig, *v.Metadata, logs.send); err != nil {
			s.rollbackVaults(created, logs.send)
			return nil, false, err
		}
	}

	// 3. Handle Keeper Password
	//
	// NOTE: a vault is written because the user typed a pair, never because the
	// engine was asked whether it consumes one - an engine that is its own
	// keeper is asked twice and pointing both at one entry is the user's answer
	if cluster.Vaults.KeeperId == nil && r.CommonConfig.KeeperPass != "" {
		logs.send("system", "saving keeper credentials to vault")
		keeperVault := vault.Vault{Type: vault.KEEPER_PASSWORD, Username: r.CommonConfig.KeeperUser, Secret: r.CommonConfig.KeeperPass}
		id, _, err := s.vaultService.Create(keeperVault)
		if err != nil {
			s.rollbackVaults(created, logs.send)
			return nil, false, err
		}
		cluster.Vaults.KeeperId = id
		created = append(created, *id)
	}

	// 4. Handle DB Password
	if cluster.Vaults.DatabaseId == nil && r.CommonConfig.DbPass != "" {
		logs.send("system", "saving database credentials to vault")
		dbVault := vault.Vault{Type: vault.DATABASE_PASSWORD, Username: r.CommonConfig.DbUser, Secret: r.CommonConfig.DbPass}
		id, _, err := s.vaultService.Create(dbVault)
		if err != nil {
			s.rollbackVaults(created, logs.send)
			return nil, false, err
		}
		cluster.Vaults.DatabaseId = id
		created = append(created, *id)
	}

	// 5. Register the cluster with its full configuration before anything is
	// started. A deployment that failed to come up is still there to be opened
	// and read, and that is exactly what an operator needs - which takes the
	// cluster's own config to reach it. Registering afterwards, or registering
	// only the nodes that came up, would leave a failed deploy invisible and
	// untroubleshootable.
	logs.send("system", "updating cluster configuration")
	if _, err := s.Update(cluster); err != nil {
		s.rollbackVaults(created, logs.send)
		return nil, false, err
	}

	// 6. Deploy every node by running its own command; the cluster owns only
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
	if !deployed {
		logs.send("system", fmt.Sprintf(
			"%d of %d node(s) deployed: the cluster is registered, open the ones that failed to see why",
			up.Load(), len(r.Nodes)))
	}

	// 7. Post-scripts run after every node is up, in node order - a script
	// that needs the whole cluster running (etcd's auth enable, mongo's
	// rs.initiate) therefore belongs on the last node.
	initialized := deployed
	if deployed {
		initialized = s.postDeploy(cluster, r.Nodes, logs.send) == 0
	} else if slices.ContainsFunc(r.Nodes, func(n DeployNode) bool { return len(n.PostScripts) > 0 }) {
		logs.send("system", "skipping post-deploy initialization: not every node deployed successfully")
	}

	return logs.list(), deployed && initialized, nil
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

// authorizeSshKey installs the freshly generated public key on every host this
// deploy will reach. A generated key authenticates nothing on its own - the
// host has never seen it - so without this the deploy offers a key that is
// rejected everywhere. The typed ssh password is what installs it, and that is
// the only thing that password is ever used for: every later connection is by
// key, out of the vault.
//
// It runs only where the key was generated. A vault the user picked is one
// whose key they already authorized, and there is no password to install it
// with anyway - a vault and inline credentials are two answers to one question
// and are rejected together.
//
// A failure aborts the deploy rather than being logged and stepped over, like
// every other step of this preparation phase: nothing has been started on any
// host yet, so there is nothing left behind to troubleshoot, and every node
// would otherwise fail with the same authentication error one at a time.
func (s *Service) authorizeSshKey(nodes []DeployNode, c CommonConfig, publicKey string, logsSend func(ctx string, msg string)) error {
	for _, target := range s.sshTargets(nodes, c.SshUser, c.SshPass) {
		nodeKey := s.getNodeKey(target.Host, &target.Port)
		logsSend(nodeKey, "authorizing the generated ssh key on the host")
		if _, err := s.nodeService.PlatformSystemCopyId(node.PlatformCopyIdRequest{
			PlatformCredConnection: target,
			PublicKey:              publicKey,
		}); err != nil {
			return fmt.Errorf("failed to authorize the ssh key on %s: %w", nodeKey, err)
		}
	}
	return nil
}

// sshTargets is every distinct host and ssh port this deploy has to reach.
// Several nodes of one cluster may share a VM, and a key only has to be
// installed on it once; a node with no host is left to deployNode, which
// reports that itself.
func (s *Service) sshTargets(nodes []DeployNode, username string, password string) []node.PlatformCredConnection {
	seen := make(map[string]bool, len(nodes))
	targets := make([]node.PlatformCredConnection, 0, len(nodes))
	for _, n := range nodes {
		port := getPortValue(n.SshPort)
		key := s.getNodeKey(n.Host, &port)
		if n.Host == "" || seen[key] {
			continue
		}
		seen[key] = true
		targets = append(targets, node.PlatformCredConnection{Host: n.Host, Port: port, Username: username, Password: password})
	}
	return targets
}

// rollbackVaults removes the vaults this deploy created for itself. It runs
// only where the call fails before the cluster is registered - nothing then
// points at them and no retry can find them again. Once the cluster exists
// they are its configuration and are never removed, however the deploy goes:
// a deployment that failed to start is still there to be opened and read, and
// reaching it takes these same credentials. A failure to remove one is logged
// rather than returned, since a stray vault is not worth losing the logs over.
func (s *Service) rollbackVaults(ids []uuid.UUID, logsSend func(ctx string, msg string)) {
	for _, id := range ids {
		if err := s.vaultService.Delete(id); err != nil {
			logsSend("system", fmt.Sprintf("failed to remove the vault created for this deploy: %v", err))
		}
	}
}

// postDeploy runs every node's post-script and returns how many of them failed.
// A failure does not abort the batch - the nodes are already up - but it has to
// be stated: a post script is what turns running processes into an initialized
// cluster (etcd's auth enable, mongo's rs.initiate), so a deploy whose scripts
// all failed otherwise reads as a success.
func (s *Service) postDeploy(cluster Request, nodes []DeployNode, logsSend func(ctx string, msg string)) int {
	failed := 0
	total := 0
	for _, n := range nodes {
		if len(n.PostScripts) == 0 {
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
	return failed
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
		Plugin:      cluster.Options.Plugins.Keeper,
		Cluster:     cluster.Name,
		Name:        n.Name,
		KeeperPort:  getPortValue(n.KeeperPort),
		DbPort:      getPortValue(n.DbPort),
		Command:     n.Command,
		PostScripts: n.PostScripts,
		Connection:  s.getNodeConnection(cluster, n),
		Vaults:      s.getNodeVaults(cluster),
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

// isHalfCredential reports a pair answered with only one of its two halves.
func isHalfCredential(username string, password string) bool {
	return (username == "") != (password == "")
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
