package cluster

import (
	"fmt"
	"ivory/src/features/node"
	"strings"
	"sync"
)

func (s *Service) Deploy(r DeployRequest) (*DeployResponse, error) {
	_, err := s.Update(r.Cluster)
	if err != nil {
		return nil, err
	}

	results := make(map[string]string)
	var mu sync.Mutex
	var wg sync.WaitGroup

	for _, n := range r.Cluster.Nodes {
		wg.Add(1)
		go func(n NodeConfig) {
			defer wg.Done()

			nodeKey := s.getNodeKey(n.Host, n.KeeperPort)

			if n.Host == "" {
				mu.Lock()
				results[nodeKey] = "host is empty"
				mu.Unlock()
				return
			}
			if n.KeeperPort == nil {
				mu.Lock()
				results[nodeKey] = "keeper port is empty"
				mu.Unlock()
				return
			}
			if n.DbPort == nil {
				mu.Lock()
				results[nodeKey] = "db port is empty"
				mu.Unlock()
				return
			}
			if n.SshPort == nil {
				mu.Lock()
				results[nodeKey] = "ssh port is empty"
				mu.Unlock()
				return
			}

			envs := append([]string{}, r.Image.CommonEnvs...)
			if uEnvs, ok := r.Image.UniqueEnvs[nodeKey]; ok {
				envs = append(envs, uEnvs...)
			}

			options := ""
			if n.DbPort != nil {
				options += fmt.Sprintf("-p %d:%d ", *n.DbPort, r.Image.DbPort)
			}
			if n.KeeperPort != nil {
				options += fmt.Sprintf("-p %d:%d ", *n.KeeperPort, r.Image.KeeperPort)
			}
			if r.Image.Volume != "" {
				options += fmt.Sprintf("-v %s ", r.Image.Volume)
			}
			if r.Image.Restart != "" {
				options += fmt.Sprintf("--restart %s ", r.Image.Restart)
			}
			for _, e := range envs {
				options += fmt.Sprintf("-e %s ", e)
			}

			dockerReq := node.DockerRequest{
				Connection: node.SshConnection{
					Host:    n.Host,
					Port:    *n.SshPort,
					VaultId: r.Cluster.Vaults.SshKeyId,
				},
				Image:   r.Image.Uri,
				Options: strings.TrimSpace(options),
			}

			resp, err := s.nodeService.DockerRun(dockerReq)

			mu.Lock()
			if err != nil {
				results[nodeKey] = err.Error()
			} else if resp.ExitCode != 0 {
				results[nodeKey] = resp.Stderr
			} else {
				results[nodeKey] = "OK"
			}
			mu.Unlock()
		}(n)
	}
	wg.Wait()

	return &DeployResponse{Nodes: results}, nil
}
