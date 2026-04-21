package instance

import (
	"crypto/tls"
	"errors"
	"fmt"
	"ivory/src/clients/sidecar"
	sshclient "ivory/src/clients/ssh"
	"ivory/src/features/cert"
	"ivory/src/features/password"
	"ivory/src/features/vm"
)

type Service struct {
	sidecarClient   sidecar.Client
	sshClient       sshclient.Client
	passwordService *password.Service
	certService     *cert.Service
	vmService       *vm.Service
}

func NewService(
	sidecarClient sidecar.Client,
	sshClient sshclient.Client,
	passwordService *password.Service,
	certService *cert.Service,
	vmService *vm.Service,
) *Service {
	return &Service{
		sidecarClient:   sidecarClient,
		sshClient:       sshClient,
		passwordService: passwordService,
		certService:     certService,
		vmService:       vmService,
	}
}

func (s *Service) OverviewAuto(request InstanceAutoRequest) ([]Instance, int, *sidecar.Sidecar, error) {
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

func (s *Service) Overview(instance Request) ([]Instance, int, error) {
	request, err := s.mapRequest(instance)
	if err != nil {
		return nil, 0, err
	}
	return s.sidecarClient.Overview(request)
}

func (s *Service) Config(instance Request) (any, int, error) {
	request, err := s.mapRequest(instance)
	if err != nil {
		return nil, 0, err
	}
	return s.sidecarClient.Config(request)
}

func (s *Service) ConfigUpdate(instance Request) (any, int, error) {
	request, err := s.mapRequest(instance)
	if err != nil {
		return nil, 0, err
	}
	return s.sidecarClient.ConfigUpdate(request)
}

func (s *Service) Switchover(instance Request) (*string, int, error) {
	request, err := s.mapRequest(instance)
	if err != nil {
		return nil, 0, err
	}
	return s.sidecarClient.Switchover(request)
}

func (s *Service) DeleteSwitchover(instance Request) (*string, int, error) {
	request, err := s.mapRequest(instance)
	if err != nil {
		return nil, 0, err
	}
	return s.sidecarClient.DeleteSwitchover(request)
}

func (s *Service) Reinitialize(instance Request) (*string, int, error) {
	request, err := s.mapRequest(instance)
	if err != nil {
		return nil, 0, err
	}
	return s.sidecarClient.Reinitialize(request)
}

func (s *Service) Restart(instance Request) (*string, int, error) {
	request, err := s.mapRequest(instance)
	if err != nil {
		return nil, 0, err
	}
	return s.sidecarClient.Restart(request)
}

func (s *Service) DeleteRestart(instance Request) (*string, int, error) {
	request, err := s.mapRequest(instance)
	if err != nil {
		return nil, 0, err
	}
	return s.sidecarClient.DeleteRestart(request)
}

func (s *Service) Reload(instance Request) (*string, int, error) {
	request, err := s.mapRequest(instance)
	if err != nil {
		return nil, 0, err
	}
	return s.sidecarClient.Reload(request)
}

func (s *Service) Failover(instance Request) (*string, int, error) {
	request, err := s.mapRequest(instance)
	if err != nil {
		return nil, 0, err
	}
	return s.sidecarClient.Failover(request)
}

func (s *Service) Activate(instance Request) (*string, int, error) {
	request, err := s.mapRequest(instance)
	if err != nil {
		return nil, 0, err
	}
	return s.sidecarClient.Activate(request)
}

func (s *Service) Pause(instance Request) (*string, int, error) {
	request, err := s.mapRequest(instance)
	if err != nil {
		return nil, 0, err
	}
	return s.sidecarClient.Pause(request)
}

func (s *Service) Metrics(request SshRequest) (*SshMetrics, int, error) {
	vmModel, err := s.vmService.GetDecrypted(request.VmId)
	if err != nil {
		return nil, 0, err
	}
	result, err := s.sshClient.Metrics(*vmModel)
	if err != nil {
		return nil, 0, err
	}
	return result, 200, nil
}

func (s *Service) mapRequest(instance Request) (sidecar.Request, error) {
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
