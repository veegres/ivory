package node

import (
	"crypto/tls"
	"errors"
	"fmt"
	"ivory/src/clients/keeper"
)

func (s *Service) OverviewAuto(request AutoRequest) ([]Node, int, *keeper.Keeper, error) {
	var tlsConfig *tls.Config
	if request.Certs != nil {
		err := s.certService.EnrichTLSConfig(&tlsConfig, request.Certs)
		if err != nil {
			return nil, 0, nil, err
		}
	}
	var cred *keeper.Credentials
	if request.CredentialId != nil {
		pass, err := s.passwordService.GetDecrypted(*request.CredentialId)
		if err != nil {
			return nil, 0, nil, err
		}
		cred = &keeper.Credentials{Username: pass.Username, Password: pass.Password}
	}
	var keeperResponses []keeper.Response
	var detectedBy *keeper.Keeper
	var statusCode int
	var errorChain error
	for i, connection := range request.Connections {
		r := keeper.Request{Host: connection.Host, Port: connection.KeeperPort, Body: request.Body, TlsConfig: tlsConfig, Credentials: cred}
		var err error
		keeperResponses, statusCode, err = s.keeperClient.Overview(r)
		if err == nil {
			detectedBy = &keeper.Keeper{Host: connection.Host, Port: connection.KeeperPort}
			break
		}
		errorChain = errors.Join(errorChain, fmt.Errorf("#%d failed %d: %w", i, statusCode, err))
	}
	if keeperResponses == nil {
		return nil, statusCode, nil, errorChain
	}

	nodes := make([]Node, 0, len(keeperResponses))
	for _, kr := range keeperResponses {
		nodes = append(nodes, Node{
			Connection: Connection{
				Host:       kr.DiscoveredHost,
				KeeperPort: kr.DiscoveredKeeperPort,
				DbPort:     kr.DiscoveredDbPort,
			},
			Keeper: kr,
		})
	}

	return nodes, statusCode, detectedBy, nil
}

func (s *Service) Overview(node Request) ([]Node, int, error) {
	request, err := s.getKeeperRequest(node)
	if err != nil {
		return nil, 0, err
	}
	keeperResponses, status, err := s.keeperClient.Overview(request)
	if err != nil {
		return nil, status, err
	}

	nodes := make([]Node, 0, len(keeperResponses))
	for _, kr := range keeperResponses {
		nodes = append(nodes, Node{
			Connection: Connection{
				Host:       kr.DiscoveredHost,
				KeeperPort: kr.DiscoveredKeeperPort,
				DbPort:     kr.DiscoveredDbPort,
			},
			Keeper: kr,
		})
	}

	return nodes, status, nil
}

func (s *Service) Config(node Request) (any, int, error) {
	request, err := s.getKeeperRequest(node)
	if err != nil {
		return nil, 0, err
	}
	return s.keeperClient.Config(request)
}

func (s *Service) ConfigUpdate(node Request) (any, int, error) {
	request, err := s.getKeeperRequest(node)
	if err != nil {
		return nil, 0, err
	}
	return s.keeperClient.ConfigUpdate(request)
}

func (s *Service) Switchover(node Request) (*string, int, error) {
	request, err := s.getKeeperRequest(node)
	if err != nil {
		return nil, 0, err
	}
	return s.keeperClient.Switchover(request)
}

func (s *Service) DeleteSwitchover(node Request) (*string, int, error) {
	request, err := s.getKeeperRequest(node)
	if err != nil {
		return nil, 0, err
	}
	return s.keeperClient.DeleteSwitchover(request)
}

func (s *Service) Reinitialize(node Request) (*string, int, error) {
	request, err := s.getKeeperRequest(node)
	if err != nil {
		return nil, 0, err
	}
	return s.keeperClient.Reinitialize(request)
}

func (s *Service) Restart(node Request) (*string, int, error) {
	request, err := s.getKeeperRequest(node)
	if err != nil {
		return nil, 0, err
	}
	return s.keeperClient.Restart(request)
}

func (s *Service) DeleteRestart(node Request) (*string, int, error) {
	request, err := s.getKeeperRequest(node)
	if err != nil {
		return nil, 0, err
	}
	return s.keeperClient.DeleteRestart(request)
}

func (s *Service) Reload(node Request) (*string, int, error) {
	request, err := s.getKeeperRequest(node)
	if err != nil {
		return nil, 0, err
	}
	return s.keeperClient.Reload(request)
}

func (s *Service) Failover(node Request) (*string, int, error) {
	request, err := s.getKeeperRequest(node)
	if err != nil {
		return nil, 0, err
	}
	return s.keeperClient.Failover(request)
}

func (s *Service) Activate(node Request) (*string, int, error) {
	request, err := s.getKeeperRequest(node)
	if err != nil {
		return nil, 0, err
	}
	return s.keeperClient.Activate(request)
}

func (s *Service) Pause(node Request) (*string, int, error) {
	request, err := s.getKeeperRequest(node)
	if err != nil {
		return nil, 0, err
	}
	return s.keeperClient.Pause(request)
}
