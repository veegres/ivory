package cluster

import (
	"errors"
	"fmt"
	"ivory/clients/storage"
	"ivory/core/service/vault"
	"ivory/features/node"
	"sync"
	"time"
)

func (s *Service) Deploy(r DeployRequest) ([]string, error) {
	// 1. Validation & Preparation
	cluster := Request{
		Name:    r.CommonConfig.Cluster,
		Nodes:   r.NodeConfig,
		Options: r.ClusterOptions,
	}

	if cluster.Name == "" {
		return nil, ErrClusterNameNotProvided
	}
	if len(cluster.Nodes) == 0 {
		return nil, ErrClusterNodesNotProvided
	}
	if r.Uri == "" {
		return nil, ErrClusterImageNotProvided
	}
	if cluster.Vaults.SshKeyId == nil && (r.CommonConfig.SshUser == "" || r.CommonConfig.SshPass == "") {
		return nil, ErrSshCredentialsRequired
	}
	if cluster.Vaults.DatabaseId == nil && (r.CommonConfig.DbUser == "" || r.CommonConfig.DbPass == "") {
		return nil, ErrDatabaseCredentialsRequired
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
	if cluster.Vaults.DatabaseId == nil {
		logsSend("system", "saving database credentials to vault")
		dbVault := vault.Vault{Type: vault.DATABASE_PASSWORD, Username: r.CommonConfig.DbUser, Secret: r.CommonConfig.DbPass}
		id, _, err := s.vaultService.Create(dbVault)
		if err != nil {
			return nil, err
		}
		cluster.Vaults.DatabaseId = id
	}

	// 4. Update Cluster
	logsSend("system", "updating cluster configuration")
	_, err := s.Update(cluster)
	if err != nil {
		return nil, err
	}

	if r.Parallel {
		var wg sync.WaitGroup
		for _, n := range cluster.Nodes {
			wg.Add(1)
			go func(n NodeConfig) {
				defer wg.Done()
				s.deployNode(r, cluster, n, logsSend)
			}(n)
		}
		wg.Wait()
	} else {
		for _, n := range cluster.Nodes {
			s.deployNode(r, cluster, n, logsSend)
		}
	}

	return logs, nil
}

func (s *Service) deployNode(r DeployRequest, cluster Request, n NodeConfig, logsSend func(ctx string, msg string)) {
	nodeKey := s.getNodeKey(n.Host, n.KeeperPort)
	options := r.NodeRawImageOptions[nodeKey]

	if n.SshPort == nil {
		logsSend(nodeKey, fmt.Sprintf("deployment failed: ssh port not provided for node"))
		return
	}
	if cluster.Vaults.SshKeyId == nil {
		logsSend(nodeKey, fmt.Sprintf("deployment failed: ssh key id not provided for node"))
		return
	}
	if n.KeeperPort == nil {
		logsSend(nodeKey, fmt.Sprintf("deployment failed: keeper port not provided for node"))
		return
	}
	if n.KeeperPort == nil {
		logsSend(nodeKey, fmt.Sprintf("deployment failed: keeper port not provided for node"))
		return
	}
	if n.DbPort == nil {
		logsSend(nodeKey, fmt.Sprintf("deployment failed: database port not provided for node"))
		return
	}
	if cluster.Vaults.DatabaseId == nil {
		logsSend(nodeKey, fmt.Sprintf("deployment failed: database credentials not provided for node"))
		return
	}

	platformReq := node.PlatformUpRequest{
		Name:  cluster.Name,
		Image: r.Uri,
		Connection: node.PlatformVaultConnection{
			Host:    n.Host,
			Port:    *n.SshPort,
			VaultId: *cluster.Vaults.SshKeyId,
		},
		ImageOptions: node.ImageOptionsRequest{
			Cluster:    cluster.Name,
			Host:       n.Host,
			Dcs:        r.CommonConfig.Dcs,
			KeeperPort: *n.KeeperPort,
			DbPort:     *n.DbPort,
		},
		Vaults: node.Vaults{
			DatabaseId: *cluster.Vaults.DatabaseId,
			SshKeyId:   *cluster.Vaults.SshKeyId,
		},
		RawImageOptions: options,
	}

	// NOTE: even if connection was closed we do not want to stop deployment
	res, err := s.nodeService.PlatformContainerUp(platformReq)
	if err != nil {
		logsSend(nodeKey, fmt.Sprintf("deployment failed: %v", err))
		return
	}
	for _, log := range res {
		logsSend(nodeKey, log)
	}
	logsSend(nodeKey, "deploy finished")
}
