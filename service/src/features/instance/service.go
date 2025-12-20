package instance

import (
	"crypto/tls"
	"errors"
	"fmt"
	"ivory/src/clients/sidecar"
	"ivory/src/features/cert"
	"ivory/src/features/password"
)

type Service struct {
	sidecarClient   sidecar.Client
	passwordService *password.Service
	certService     *cert.Service
}

func NewService(
	sidecarClient sidecar.Client,
	passwordService *password.Service,
	certService *cert.Service,
) *Service {
	return &Service{
		sidecarClient:   sidecarClient,
		passwordService: passwordService,
		certService:     certService,
	}
}

func (s *Service) AutoOverview(request InstanceAutoRequest) ([]sidecar.Instance, int, error) {
	var tlsConfig *tls.Config
	if request.Certs != nil {
		err := s.certService.EnrichTLSConfig(&tlsConfig, request.Certs)
		if err != nil {
			return nil, 0, err
		}
	}
	var cred *sidecar.Credentials
	if request.CredentialId != nil {
		pass, err := s.passwordService.GetDecrypted(*request.CredentialId)
		if err != nil {
			return nil, 0, err
		}
		cred = &sidecar.Credentials{Username: pass.Username, Password: pass.Password}
	}
	var overview []sidecar.Instance
	var statusCode int
	var errorChain error
	for i, instance := range request.Sidecars {
		r := sidecar.Request{Sidecar: instance, Body: request.Body, TlsConfig: tlsConfig, Credentials: cred}
		var err error
		overview, statusCode, err = s.sidecarClient.Overview(r)
		if err == nil {
			break
		}
		errorChain = errors.Join(errorChain, fmt.Errorf("#%d failed %d: %w", i, statusCode, err))
	}
	if overview == nil {
		return nil, statusCode, errorChain
	}
	return overview, statusCode, nil
}

func (s *Service) Overview(instance InstanceRequest) ([]sidecar.Instance, int, error) {
	request, err := s.mapRequest(instance)
	if err != nil {
		return nil, 0, err
	}
	return s.sidecarClient.Overview(request)
}

func (s *Service) Config(instance InstanceRequest) (any, int, error) {
	request, err := s.mapRequest(instance)
	if err != nil {
		return nil, 0, err
	}
	return s.sidecarClient.Config(request)
}

func (s *Service) ConfigUpdate(instance InstanceRequest) (any, int, error) {
	request, err := s.mapRequest(instance)
	if err != nil {
		return nil, 0, err
	}
	return s.sidecarClient.ConfigUpdate(request)
}

func (s *Service) Switchover(instance InstanceRequest) (*string, int, error) {
	request, err := s.mapRequest(instance)
	if err != nil {
		return nil, 0, err
	}
	return s.sidecarClient.Switchover(request)
}

func (s *Service) DeleteSwitchover(instance InstanceRequest) (*string, int, error) {
	request, err := s.mapRequest(instance)
	if err != nil {
		return nil, 0, err
	}
	return s.sidecarClient.DeleteSwitchover(request)
}

func (s *Service) Reinitialize(instance InstanceRequest) (*string, int, error) {
	request, err := s.mapRequest(instance)
	if err != nil {
		return nil, 0, err
	}
	return s.sidecarClient.Reinitialize(request)
}

func (s *Service) Restart(instance InstanceRequest) (*string, int, error) {
	request, err := s.mapRequest(instance)
	if err != nil {
		return nil, 0, err
	}
	return s.sidecarClient.Restart(request)
}

func (s *Service) DeleteRestart(instance InstanceRequest) (*string, int, error) {
	request, err := s.mapRequest(instance)
	if err != nil {
		return nil, 0, err
	}
	return s.sidecarClient.DeleteRestart(request)
}

func (s *Service) Reload(instance InstanceRequest) (*string, int, error) {
	request, err := s.mapRequest(instance)
	if err != nil {
		return nil, 0, err
	}
	return s.sidecarClient.Reload(request)
}

func (s *Service) Failover(instance InstanceRequest) (*string, int, error) {
	request, err := s.mapRequest(instance)
	if err != nil {
		return nil, 0, err
	}
	return s.sidecarClient.Failover(request)
}

func (s *Service) Activate(instance InstanceRequest) (*string, int, error) {
	request, err := s.mapRequest(instance)
	if err != nil {
		return nil, 0, err
	}
	return s.sidecarClient.Activate(request)
}

func (s *Service) Pause(instance InstanceRequest) (*string, int, error) {
	request, err := s.mapRequest(instance)
	if err != nil {
		return nil, 0, err
	}
	return s.sidecarClient.Pause(request)
}

func (s *Service) mapRequest(instance InstanceRequest) (sidecar.Request, error) {
	request := sidecar.Request{Sidecar: instance.Sidecar, Body: instance.Body}
	if instance.Certs != nil {
		err := s.certService.EnrichTLSConfig(&request.TlsConfig, instance.Certs)
		if err != nil {
			return request, err
		}
	}
	if instance.CredentialId != nil {
		pass, err := s.passwordService.GetDecrypted(*instance.CredentialId)
		if err != nil {
			return request, err
		}
		request.Credentials = &sidecar.Credentials{Username: pass.Username, Password: pass.Password}
	}
	return request, nil
}
