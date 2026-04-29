package node

import (
	"ivory/src/plugins/os"
)

func (s *Service) Metrics(request SshRequest) (*SshResponseMetrics, error) {
	connection, err := s.getSshConnection(request.Connection)
	if err != nil {
		return nil, err
	}
	// NOTE: today we only support Linux nodes for metrics collection.
	// In the future, the node connection should carry its OS type.
	linuxAdapter, err := s.osRegistry.Get(os.Linux)
	if err != nil {
		return nil, err
	}
	result, err := linuxAdapter.Metrics(*connection)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *Service) DockerDeploy(request DockerRequest) (*DockerResult, error) {
	connection, err := s.getSshConnection(request.Connection)
	if err != nil {
		return nil, err
	}
	linuxAdapter, err := s.osRegistry.Get(os.Linux)
	if err != nil {
		return nil, err
	}
	return linuxAdapter.DockerDeploy(*connection, request.Image, request.Options)
}

func (s *Service) DockerStop(request DockerRequest) (*DockerResult, error) {
	connection, err := s.getSshConnection(request.Connection)
	if err != nil {
		return nil, err
	}
	linuxAdapter, err := s.osRegistry.Get(os.Linux)
	if err != nil {
		return nil, err
	}
	return linuxAdapter.DockerStop(*connection, request.Container)
}

func (s *Service) DockerRun(request DockerRequest) (*DockerResult, error) {
	connection, err := s.getSshConnection(request.Connection)
	if err != nil {
		return nil, err
	}
	linuxAdapter, err := s.osRegistry.Get(os.Linux)
	if err != nil {
		return nil, err
	}
	return linuxAdapter.DockerRun(*connection, request.Options, request.Image)
}

func (s *Service) DockerDelete(request DockerRequest) (*DockerResult, error) {
	connection, err := s.getSshConnection(request.Connection)
	if err != nil {
		return nil, err
	}
	linuxAdapter, err := s.osRegistry.Get(os.Linux)
	if err != nil {
		return nil, err
	}
	return linuxAdapter.DockerDelete(*connection, request.Container)
}

func (s *Service) DockerList(request SshRequest) (*DockerResult, error) {
	connection, err := s.getSshConnection(request.Connection)
	if err != nil {
		return nil, err
	}
	linuxAdapter, err := s.osRegistry.Get(os.Linux)
	if err != nil {
		return nil, err
	}
	return linuxAdapter.DockerList(*connection)
}

func (s *Service) DockerLogs(request DockerLogsRequest) (*DockerResult, error) {
	connection, err := s.getSshConnection(request.Connection)
	if err != nil {
		return nil, err
	}
	linuxAdapter, err := s.osRegistry.Get(os.Linux)
	if err != nil {
		return nil, err
	}
	return linuxAdapter.DockerLogs(*connection, request.Container, request.Tail)
}
