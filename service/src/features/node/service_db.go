package node

import (
	"crypto/tls"
	"errors"
	"fmt"
	"ivory/src/clients/sidecar"
)

func (s *Service) OverviewAuto(request NodeAutoRequest) ([]Node, int, *sidecar.Sidecar, error) {
	var tlsConfig *tls.Config
	if request.Certs != nil {
		err := s.certService.EnrichTLSConfig(&tlsConfig, request.Certs)
		if err != nil {
			return nil, 0, nil, err
		}
	}
	var cred *sidecar.Credentials
	if request.CredentialId != nil {
		pass, err := s.passwordService.GetDecrypted(*request.CredentialId)
		if err != nil {
			return nil, 0, nil, err
		}
		cred = &sidecar.Credentials{Username: pass.Username, Password: pass.Password}
	}
	var overview []sidecar.Instance
	var detectedBy *sidecar.Sidecar
	var statusCode int
	var errorChain error
	for i, instance := range request.Sidecars {
		r := sidecar.Request{Sidecar: instance, Body: request.Body, TlsConfig: tlsConfig, Credentials: cred}
		var err error
		overview, statusCode, err = s.sidecarClient.Overview(r)
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
	request, err := s.mapRequest(node)
	if err != nil {
		return nil, 0, err
	}
	return s.sidecarClient.Overview(request)
}

func (s *Service) Config(node Request) (any, int, error) {
	request, err := s.mapRequest(node)
	if err != nil {
		return nil, 0, err
	}
	return s.sidecarClient.Config(request)
}

func (s *Service) ConfigUpdate(node Request) (any, int, error) {
	request, err := s.mapRequest(node)
	if err != nil {
		return nil, 0, err
	}
	return s.sidecarClient.ConfigUpdate(request)
}

func (s *Service) Switchover(node Request) (*string, int, error) {
	request, err := s.mapRequest(node)
	if err != nil {
		return nil, 0, err
	}
	return s.sidecarClient.Switchover(request)
}

func (s *Service) DeleteSwitchover(node Request) (*string, int, error) {
	request, err := s.mapRequest(node)
	if err != nil {
		return nil, 0, err
	}
	return s.sidecarClient.DeleteSwitchover(request)
}

func (s *Service) Reinitialize(node Request) (*string, int, error) {
	request, err := s.mapRequest(node)
	if err != nil {
		return nil, 0, err
	}
	return s.sidecarClient.Reinitialize(request)
}

func (s *Service) Restart(node Request) (*string, int, error) {
	request, err := s.mapRequest(node)
	if err != nil {
		return nil, 0, err
	}
	return s.sidecarClient.Restart(request)
}

func (s *Service) DeleteRestart(node Request) (*string, int, error) {
	request, err := s.mapRequest(node)
	if err != nil {
		return nil, 0, err
	}
	return s.sidecarClient.DeleteRestart(request)
}

func (s *Service) Reload(node Request) (*string, int, error) {
	request, err := s.mapRequest(node)
	if err != nil {
		return nil, 0, err
	}
	return s.sidecarClient.Reload(request)
}

func (s *Service) Failover(node Request) (*string, int, error) {
	request, err := s.mapRequest(node)
	if err != nil {
		return nil, 0, err
	}
	return s.sidecarClient.Failover(request)
}

func (s *Service) Activate(node Request) (*string, int, error) {
	request, err := s.mapRequest(node)
	if err != nil {
		return nil, 0, err
	}
	return s.sidecarClient.Activate(request)
}

func (s *Service) Pause(node Request) (*string, int, error) {
	request, err := s.mapRequest(node)
	if err != nil {
		return nil, 0, err
	}
	return s.sidecarClient.Pause(request)
}
