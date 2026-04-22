package node

import (
	"crypto/tls"
	"errors"
	"fmt"
	"ivory/src/clients/keeper"
)

func (s *Service) OverviewAuto(request NodeAutoRequest) ([]Node, int, *keeper.Keeper, error) {
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
	var overview []keeper.Node
	var detectedBy *keeper.Keeper
	var statusCode int
	var errorChain error
	for i, instance := range request.Keepers {
		r := keeper.Request{Keeper: instance, Body: request.Body, TlsConfig: tlsConfig, Credentials: cred}
		var err error
		overview, statusCode, err = s.keeperClient.Overview(r)
		if err == nil {
			detectedBy = &instance
			break
		}
		errorChain = errors.Join(errorChain, fmt.Errorf("#%d failed %d: %w", i, statusCode, err))
	}
	if overview == nil {
		return nil, statusCode, nil, errorChain
	}
	return overview, statusCode, detectedBy, nil
}

func (s *Service) Overview(node Request) ([]Node, int, error) {
	request, err := s.getKeeperRequest(node)
	if err != nil {
		return nil, 0, err
	}
	return s.keeperClient.Overview(request)
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
