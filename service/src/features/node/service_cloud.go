package node

import (
	"ivory/src/plugins/cloud"
)

func (s *Service) CloudCopyId(c CloudCredConnection, publicKey string) error {
	adapter, err := s.cloudRegistry.Get(cloud.Onprem)
	if err != nil {
		return err
	}
	con := s.getCloudCredConnection(c)
	return adapter.CopyId(*con, publicKey)
}

func (s *Service) Metrics(c CloudVaultConnection) (*MetricsResponse, error) {
	adapter, conn, err := s.getCloudAdapter(c)
	if err != nil {
		return nil, err
	}
	return adapter.Metrics(*conn)
}

func (s *Service) ContainerStop(request ContainerRequest) (*ContainerResponse, error) {
	adapter, conn, err := s.getCloudAdapter(request.Connection)
	if err != nil {
		return nil, err
	}
	return adapter.ContainerStop(*conn, request.Container)
}

func (s *Service) ContainerRun(request ContainerRequest) (*ContainerResponse, error) {
	adapter, conn, err := s.getCloudAdapter(request.Connection)
	if err != nil {
		return nil, err
	}
	return adapter.ContainerRun(*conn, request.Options, request.Image)
}

func (s *Service) ContainerDelete(request ContainerRequest) (*ContainerResponse, error) {
	adapter, conn, err := s.getCloudAdapter(request.Connection)
	if err != nil {
		return nil, err
	}
	return adapter.ContainerDelete(*conn, request.Container)
}

func (s *Service) ContainerList(c CloudVaultConnection) (*ContainerResponse, error) {
	adapter, conn, err := s.getCloudAdapter(c)
	if err != nil {
		return nil, err
	}
	return adapter.ContainerList(*conn)
}

func (s *Service) ContainerLogs(request ContainerLogsRequest) (*ContainerResponse, error) {
	adapter, conn, err := s.getCloudAdapter(request.Connection)
	if err != nil {
		return nil, err
	}
	return adapter.ContainerLogs(*conn, request.Container, request.Tail)
}
