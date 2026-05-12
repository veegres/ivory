package cluster

import (
	"errors"
	"fmt"
	"ivory/src/features/node"
	"ivory/src/features/vault"
	"ivory/src/storage/db"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

func (s *Service) Deploy(r DeployRequest) ([]string, error) {
	// 1. Validation & Preparation
	if r.Cluster.Name == "" {
		return nil, errors.New("cluster name not provided")
	}
	if len(r.Cluster.Nodes) == 0 {
		return nil, errors.New("cluster nodes not provided")
	}
	if r.Image.Uri == "" {
		return nil, errors.New("cluster image not provided")
	}
	if r.Cluster.Vaults.SshKeyId == nil && (r.Cred.Ssh.Username == "" || r.Cred.Ssh.Password == "") {
		return nil, errors.New("ssh credentials are required")
	}
	if r.Cluster.Vaults.DatabaseId == nil && (r.Cred.Db.Username == "" || r.Cred.Db.Password == "") {
		return nil, errors.New("database credentials are required")
	}
	if _, e := s.Get(r.Cluster.Name); !errors.Is(e, db.ErrNotFound) {
		return nil, errors.New("cluster name is already taken")
	}

	var mu sync.Mutex
	logs := make([]string, 0)
	appendLog := func(ctx string, msg string) {
		mu.Lock()
		logs = append(logs, fmt.Sprintf("%s | %s | %s", time.Now().Format("2006-01-02 15:04:05"), ctx, msg))
		mu.Unlock()
	}

	// 2. Handle SSH Key
	var sshPubKey string
	var needSshCopy bool
	if r.Cluster.Vaults.SshKeyId == nil {
		appendLog("system", "generating ssh key and saving it to vault")
		sshVault := vault.Vault{Type: vault.SSH_KEY, Username: r.Cred.Ssh.Username}
		id, v, err := s.vaultService.Create(sshVault)
		if err != nil {
			return nil, err
		}
		r.Cluster.Vaults.SshKeyId = id
		if v.Metadata == nil {
			return nil, errors.New("ssh key from vault is missing metadata (public key)")
		}
		sshPubKey = *v.Metadata
		needSshCopy = true
	} else {
		appendLog("system", "using ssh key from vault")
		v, err := s.vaultService.Get(*r.Cluster.Vaults.SshKeyId)
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
	if r.Cluster.Vaults.DatabaseId == nil {
		appendLog("system", "saving database credentials to vault")
		dbVault := vault.Vault{Type: vault.DATABASE_PASSWORD, Username: r.Cred.Db.Username, Secret: r.Cred.Db.Password}
		id, _, err := s.vaultService.Create(dbVault)
		if err != nil {
			return nil, err
		}
		r.Cluster.Vaults.DatabaseId = id
	} else {
		appendLog("system", "using database credentials from vault")
		v, err := s.vaultService.GetDecrypted(*r.Cluster.Vaults.DatabaseId)
		if err != nil {
			return nil, fmt.Errorf("failed to get database credentials from vault: %w", err)
		}
		r.Cred.Db.Username = v.Username
		r.Cred.Db.Password = v.Secret
	}

	// 4. Update Cluster
	appendLog("system", "updating cluster configuration")
	_, err := s.Update(r.Cluster)
	if err != nil {
		return nil, err
	}

	// 5. Deploy Docker containers
	var wg sync.WaitGroup
	for _, n := range r.Cluster.Nodes {
		wg.Add(1)
		go func(n NodeConfig) {
			defer wg.Done()
			nodeKey := s.getNodeKey(n.Host, n.KeeperPort)

			if n.Host == "" {
				appendLog(nodeKey, "host is empty, skipping node")
				return
			}
			if n.SshPort == nil {
				appendLog(nodeKey, "ssh port is empty, skipping node")
				return
			}

			if needSshCopy {
				appendLog(nodeKey, "saving ssh key to vm")
				conn := node.SshCredConnection{
					Host:     n.Host,
					Port:     *n.SshPort,
					Username: r.Cred.Ssh.Username,
					Password: r.Cred.Ssh.Password,
				}
				err := s.nodeService.SshCopyId(conn, sshPubKey)
				if err != nil {
					appendLog(nodeKey, fmt.Sprintf("failed to copy ssh key: %v", err))
					return
				}
			}

			// Interpolate options
			template, ok := r.Image.Options[nodeKey]
			if !ok {
				appendLog(nodeKey, "no options provided for node")
				return
			}

			values := map[string]string{
				"cluster":  r.Cluster.Name,
				"node":     n.Host,
				"username": r.Cred.Db.Username,
				"password": r.Cred.Db.Password,
			}
			if n.KeeperPort != nil {
				values["keeperPort"] = strconv.Itoa(*n.KeeperPort)
			}
			if n.DbPort != nil {
				values["dbPort"] = strconv.Itoa(*n.DbPort)
			}

			options := s.normalizeDockerOptions(s.getInterpolatedString(template, values))

			appendLog(nodeKey, fmt.Sprintf("deploying with options: %s", options))

			dockerReq := node.DockerRequest{
				Connection: node.SshVaultConnection{
					Host:    n.Host,
					Port:    *n.SshPort,
					VaultId: r.Cluster.Vaults.SshKeyId,
				},
				Image:   r.Image.Uri,
				Options: options,
			}

			resp, err := s.nodeService.DockerRun(dockerReq)

			if err != nil {
				appendLog(nodeKey, fmt.Sprintf("%v", err))
			} else if resp.ExitCode != 0 {
				appendLog(nodeKey, resp.Stderr)
			} else {
				appendLog(nodeKey, "OK")
			}
		}(n)
	}
	wg.Wait()

	return logs, nil
}

func (s *Service) normalizeDockerOptions(options string) string {
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
