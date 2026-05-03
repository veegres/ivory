package node

import (
	"crypto/tls"
	"errors"
	"fmt"
	"ivory/src/plugins/keeper"
)

func (s *Service) ListAuto(request KeeperAutoRequest) ([]KeeperResponse, int, *KeeperConnection, error) {
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
	var detectedBy *KeeperConnection
	var statusCode int
	var errorChain error
	for i, connection := range request.Connections {
		r := keeper.Request{Host: connection.Host, Port: connection.Port, Body: request.Body, TlsConfig: tlsConfig, Credentials: cred}
		var err error
		keeperResponses, statusCode, err = client.List(r)
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
		nodes = append(nodes, kr)
	}

	return nodes, statusCode, detectedBy, errorChain
}

func (s *Service) List(r KeeperRequest) ([]KeeperResponse, int, error) {
	client, request, err := s.getKeeperAdapter(r)
	if err != nil {
		return nil, 0, err
	}
	keeperResponses, status, err := client.List(request)
	if err != nil {
		return nil, status, err
	}

	nodes := make([]KeeperResponse, 0, len(keeperResponses))
	for _, kr := range keeperResponses {
		nodes = append(nodes, kr)
	}

	return nodes, status, nil
}

func (s *Service) Config(r KeeperRequest) (any, int, error) {
	client, request, err := s.getKeeperAdapter(r)
	if err != nil {
		return nil, 0, err
	}
	return client.Config(request)
}

func (s *Service) ConfigUpdate(r KeeperRequest) (any, int, error) {
	client, request, err := s.getKeeperAdapter(r)
	if err != nil {
		return nil, 0, err
	}
	return client.ConfigUpdate(request)
}

func (s *Service) Switchover(r KeeperRequest) (*string, int, error) {
	client, request, err := s.getKeeperAdapter(r)
	if err != nil {
		return nil, 0, err
	}
	return client.Switchover(request)
}

func (s *Service) DeleteSwitchover(r KeeperRequest) (*string, int, error) {
	client, request, err := s.getKeeperAdapter(r)
	if err != nil {
		return nil, 0, err
	}
	return client.DeleteSwitchover(request)
}

func (s *Service) Reinitialize(r KeeperRequest) (*string, int, error) {
	client, request, err := s.getKeeperAdapter(r)
	if err != nil {
		return nil, 0, err
	}
	return client.Reinitialize(request)
}

func (s *Service) Restart(r KeeperRequest) (*string, int, error) {
	client, request, err := s.getKeeperAdapter(r)
	if err != nil {
		return nil, 0, err
	}
	return client.Restart(request)
}

func (s *Service) DeleteRestart(r KeeperRequest) (*string, int, error) {
	client, request, err := s.getKeeperAdapter(r)
	if err != nil {
		return nil, 0, err
	}
	return client.DeleteRestart(request)
}

func (s *Service) Reload(r KeeperRequest) (*string, int, error) {
	client, request, err := s.getKeeperAdapter(r)
	if err != nil {
		return nil, 0, err
	}
	return client.Reload(request)
}

func (s *Service) Failover(r KeeperRequest) (*string, int, error) {
	client, request, err := s.getKeeperAdapter(r)
	if err != nil {
		return nil, 0, err
	}
	return client.Failover(request)
}

func (s *Service) Activate(r KeeperRequest) (*string, int, error) {
	client, request, err := s.getKeeperAdapter(r)
	if err != nil {
		return nil, 0, err
	}
	return client.Activate(request)
}

func (s *Service) Pause(r KeeperRequest) (*string, int, error) {
	client, request, err := s.getKeeperAdapter(r)
	if err != nil {
		return nil, 0, err
	}
	return client.Pause(request)
}
