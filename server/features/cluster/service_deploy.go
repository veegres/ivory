package cluster

import (
	"errors"
	"fmt"
	"ivory/core/service/job"
	"ivory/core/service/vault"
	"ivory/core/store"
	"ivory/features/node"
	"regexp"
	"strconv"
	"strings"
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
		return nil, errors.New("cluster name not provided")
	}
	if len(cluster.Nodes) == 0 {
		return nil, errors.New("cluster nodes not provided")
	}
	if r.Uri == "" {
		return nil, errors.New("cluster image not provided")
	}
	if cluster.Vaults.SshKeyId == nil && (r.CommonConfig.SshUser == "" || r.CommonConfig.SshPass == "") {
		return nil, errors.New("ssh credentials are required")
	}
	if cluster.Vaults.DatabaseId == nil && (r.CommonConfig.DbUser == "" || r.CommonConfig.DbPass == "") {
		return nil, errors.New("database credentials are required")
	}
	if _, e := s.Get(cluster.Name); !errors.Is(e, store.ErrNotFound) {
		return nil, errors.New("cluster name is already taken")
	}

	var mu sync.Mutex
	logs := make([]string, 0)
	msgSend := func(ctx string, msg string) {
		mu.Lock()
		logs = append(logs, fmt.Sprintf("%s | %s | %s", time.Now().Format("2006-01-02 15:04:05"), ctx, msg))
		mu.Unlock()
	}

	// 2. Handle SSH Key
	var sshPubKey string
	var needSshCopy bool
	if cluster.Vaults.SshKeyId == nil {
		msgSend("system", "generating ssh key and saving it to vault")
		sshVault := vault.Vault{Type: vault.SSH_KEY, Username: r.CommonConfig.SshUser}
		id, v, err := s.vaultService.Create(sshVault)
		if err != nil {
			return nil, err
		}
		cluster.Vaults.SshKeyId = id
		if v.Metadata == nil {
			return nil, errors.New("ssh key from vault is missing metadata (public key)")
		}
		sshPubKey = *v.Metadata
		needSshCopy = true
	} else {
		msgSend("system", "using ssh key from vault")
		v, err := s.vaultService.Get(*cluster.Vaults.SshKeyId)
		if err != nil {
			return nil, fmt.Errorf("failed to get ssh key from vault: %w", err)
		}
		if v.Metadata == nil {
			return nil, errors.New("ssh key from vault is missing metadata (public key)")
		}
		sshPubKey = *v.Metadata
		needSshCopy = false
	}

	// 3. Handle DB Password
	if cluster.Vaults.DatabaseId == nil {
		msgSend("system", "saving database credentials to vault")
		dbVault := vault.Vault{Type: vault.DATABASE_PASSWORD, Username: r.CommonConfig.DbUser, Secret: r.CommonConfig.DbPass}
		id, _, err := s.vaultService.Create(dbVault)
		if err != nil {
			return nil, err
		}
		cluster.Vaults.DatabaseId = id
	} else {
		msgSend("system", "using database credentials from vault")
		v, err := s.vaultService.GetDecrypted(*cluster.Vaults.DatabaseId)
		if err != nil {
			return nil, fmt.Errorf("failed to get database credentials from vault: %w", err)
		}
		r.CommonConfig.DbUser = v.Username
		r.CommonConfig.DbPass = v.Secret
	}

	// 4. Update Cluster
	msgSend("system", "updating cluster configuration")
	_, err := s.Update(cluster)
	if err != nil {
		return nil, err
	}

	// 5. Deploy database
	if r.Parallel {
		var wg sync.WaitGroup
		for _, n := range cluster.Nodes {
			wg.Add(1)
			go func(n NodeConfig) {
				defer wg.Done()
				s.deployNode(r, cluster, n, needSshCopy, sshPubKey, msgSend)
			}(n)
		}
		wg.Wait()
	} else {
		for _, n := range cluster.Nodes {
			s.deployNode(r, cluster, n, needSshCopy, sshPubKey, msgSend)
		}
	}

	return logs, nil
}

func (s *Service) deployNode(r DeployRequest, cluster Request, n NodeConfig, needSshCopy bool, sshPubKey string, send func(ctx string, msg string)) {
	nodeKey := s.getNodeKey(n.Host, n.KeeperPort)

	if n.Host == "" {
		send(nodeKey, "host is empty, skipping node")
		return
	}
	if n.SshPort == nil {
		send(nodeKey, "ssh port is empty, skipping node")
		return
	}

	if needSshCopy {
		send(nodeKey, "saving ssh key to vm")
		conn := node.PlatformCredConnection{
			Host:     n.Host,
			Port:     *n.SshPort,
			Username: r.CommonConfig.SshUser,
			Password: r.CommonConfig.SshPass,
		}
		req := node.PlatformCopyIdRequest{PlatformCredConnection: conn, PublicKey: sshPubKey}
		_, err := s.nodeService.PlatformCopyId(req)
		if err != nil {
			send(nodeKey, fmt.Sprintf("failed to copy ssh key: %v", err))
			return
		}
	}

	// Interpolate options
	template, ok := r.NodeRawOptions[nodeKey]
	if !ok {
		send(nodeKey, "no options provided for node")
		return
	}

	values := map[string]string{
		"cluster": cluster.Name,
		"dcs":     r.CommonConfig.Dcs,
		"host":    n.Host,
		"node":    n.Host,
		"dbUser":  r.CommonConfig.DbUser,
		"dbPass":  r.CommonConfig.DbPass,
		"sshUser": r.CommonConfig.SshUser,
		"sshPass": r.CommonConfig.SshPass,
	}
	if n.SshPort != nil {
		values["sshPort"] = strconv.Itoa(*n.SshPort)
	}
	if n.KeeperPort != nil {
		values["keeperPort"] = strconv.Itoa(*n.KeeperPort)
	}
	if n.DbPort != nil {
		values["dbPort"] = strconv.Itoa(*n.DbPort)
	}

	options := s.normalizeDatabaseOptions(s.getInterpolatedString(template, values))

	send(nodeKey, fmt.Sprintf("deploying with options: %s", options))

	platformReq := node.PlatformUpRequest{
		Connection: node.PlatformVaultConnection{
			Host:    n.Host,
			Port:    *n.SshPort,
			VaultId: cluster.Vaults.SshKeyId,
		},
		Image:   r.Uri,
		Options: options,
	}

	// NOTE: even if connection was closed we do not want to stop deployment
	mockClose := make(<-chan struct{})
	s.nodeService.PlatformContainerUp(platformReq, "cluster-deploy", mockClose, func(event job.Message) {
		send(nodeKey, event.Message)
	})
	send(nodeKey, "deploy stream finished")
}

func (s *Service) normalizeDatabaseOptions(options string) string {
	return strings.Join(strings.Fields(options), " ")
}

func (s *Service) getInterpolatedString(template string, values map[string]string) string {
	return regexp.MustCompile(`{{(\w+)}}`).ReplaceAllStringFunc(template, func(match string) string {
		key := match[2 : len(match)-2]
		if val, ok := values[key]; ok {
			return val
		}
		return match
	})
}
