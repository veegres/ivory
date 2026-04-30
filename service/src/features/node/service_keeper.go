package node

import (
	"crypto/tls"
	"errors"
	"fmt"
	"ivory/src/plugins/keeper"
)

func (s *Service) OverviewAuto(request KeeperAutoRequest) ([]KeeperResponse, int, *Connection, error) {
	client, errClient := s.keeperRegistry.Get(request.Plugin)
	if errClient != nil {
		return nil, 0, nil, errClient
	}
	var tlsConfig *tls.Config
	if request.Certs != nil {
		err := s.certService.EnrichTLSConfig(&tlsConfig, request.Certs)
		if err != nil {
			return nil, 0, nil, err
		}
	}
	var cred *keeper.Credentials
	if request.VaultId != nil {
		pass, err := s.vaultService.GetDecrypted(*request.VaultId)
		if err != nil {
			return nil, 0, nil, err
		}
		cred = &keeper.Credentials{Username: pass.Username, Password: pass.Secret}
	}
	var keeperResponses []keeper.Response
	var detectedBy *Connection
	var statusCode int
	var errorChain error
	for i, connection := range request.Connections {
		r := keeper.Request{Host: connection.Host, Port: connection.KeeperPort, Body: request.Body, TlsConfig: tlsConfig, Credentials: cred}
		var err error
		keeperResponses, statusCode, err = client.Overview(r)
		if err == nil {
			detectedBy = &connection
			break
		}
		errorChain = errors.Join(errorChain, fmt.Errorf("#%d failed %d: %w", i, statusCode, err))
	}
	if keeperResponses == nil {
		return nil, statusCode, nil, errorChain
	}

	nodes := make([]KeeperResponse, 0, len(keeperResponses))
	for _, kr := range keeperResponses {
		nodes = append(nodes, KeeperResponse{
			Connection: Connection{
				Host:       *kr.DiscoveredHost,
				KeeperPort: *kr.DiscoveredKeeperPort,
				DbPort:     *kr.DiscoveredDbPort,
			},
			Keeper: kr,
		})
	}

	return nodes, statusCode, detectedBy, nil
}

func (s *Service) Overview(r KeeperConnection) ([]KeeperResponse, int, error) {
	client, request, err := s.getKeeperAdapter(r)
	if err != nil {
		return nil, 0, err
	}
	keeperResponses, status, err := client.Overview(request)
	if err != nil {
		return nil, status, err
	}

	nodes := make([]KeeperResponse, 0, len(keeperResponses))
	for _, kr := range keeperResponses {
		nodes = append(nodes, KeeperResponse{
			Connection: Connection{
				Host:       *kr.DiscoveredHost,
				KeeperPort: *kr.DiscoveredKeeperPort,
				DbPort:     *kr.DiscoveredDbPort,
				SshPort:    22,
			},
			Keeper: kr,
		})
	}

	return nodes, status, nil
}

func (s *Service) Config(r KeeperConnection) (any, int, error) {
	client, request, err := s.getKeeperAdapter(r)
	if err != nil {
		return nil, 0, err
	}
	return client.Config(request)
}

func (s *Service) ConfigUpdate(r KeeperConnection) (any, int, error) {
	client, request, err := s.getKeeperAdapter(r)
	if err != nil {
		return nil, 0, err
	}
	return client.ConfigUpdate(request)
}

func (s *Service) Switchover(r KeeperConnection) (*string, int, error) {
	client, request, err := s.getKeeperAdapter(r)
	if err != nil {
		return nil, 0, err
	}
	return client.Switchover(request)
}

func (s *Service) DeleteSwitchover(r KeeperConnection) (*string, int, error) {
	client, request, err := s.getKeeperAdapter(r)
	if err != nil {
		return nil, 0, err
	}
	return client.DeleteSwitchover(request)
}

func (s *Service) Reinitialize(r KeeperConnection) (*string, int, error) {
	client, request, err := s.getKeeperAdapter(r)
	if err != nil {
		return nil, 0, err
	}
	return client.Reinitialize(request)
}

func (s *Service) Restart(r KeeperConnection) (*string, int, error) {
	client, request, err := s.getKeeperAdapter(r)
	if err != nil {
		return nil, 0, err
	}
	return client.Restart(request)
}

func (s *Service) DeleteRestart(r KeeperConnection) (*string, int, error) {
	client, request, err := s.getKeeperAdapter(r)
	if err != nil {
		return nil, 0, err
	}
	return client.DeleteRestart(request)
}

func (s *Service) Reload(r KeeperConnection) (*string, int, error) {
	client, request, err := s.getKeeperAdapter(r)
	if err != nil {
		return nil, 0, err
	}
	return client.Reload(request)
}

func (s *Service) Failover(r KeeperConnection) (*string, int, error) {
	client, request, err := s.getKeeperAdapter(r)
	if err != nil {
		return nil, 0, err
	}
	return client.Failover(request)
}

func (s *Service) Activate(r KeeperConnection) (*string, int, error) {
	client, request, err := s.getKeeperAdapter(r)
	if err != nil {
		return nil, 0, err
	}
	return client.Activate(request)
}

func (s *Service) Pause(r KeeperConnection) (*string, int, error) {
	client, request, err := s.getKeeperAdapter(r)
	if err != nil {
		return nil, 0, err
	}
	return client.Pause(request)
}
