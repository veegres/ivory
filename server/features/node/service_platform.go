package node

import (
	"ivory/plugins/platform"
)

func (s *Service) PlatformCopyId(c PlatformCredConnection, publicKey string) error {
	adapter, err := s.platformRegistry.Get(platform.Onprem)
	if err != nil {
		return err
	}
	con := s.getPlatformCredConnection(c)
	return adapter.CopyId(*con, publicKey)
}

func (s *Service) Metrics(c PlatformVaultConnection) (*MetricsResponse, error) {
	adapter, conn, err := s.getPlatformAdapter(c)
	if err != nil {
		return nil, err
	}
	return adapter.Metrics(*conn)
}

func (s *Service) PlatformStop(request PlatformDeployRequest) (*PlatformResponse, error) {
	adapter, conn, err := s.getPlatformAdapter(request.Connection)
	if err != nil {
		return nil, err
	}
	return adapter.Stop(*conn, request.Name)
}

func (s *Service) PlatformDeploy(request PlatformDeployRequest) (*PlatformResponse, error) {
	adapter, conn, err := s.getPlatformAdapter(request.Connection)
	if err != nil {
		return nil, err
	}
	return adapter.Deploy(*conn, request.Options, request.Image)
}

func (s *Service) PlatformDelete(request PlatformDeployRequest) (*PlatformResponse, error) {
	adapter, conn, err := s.getPlatformAdapter(request.Connection)
	if err != nil {
		return nil, err
	}
	return adapter.Delete(*conn, request.Name)
}

func (s *Service) PlatformList(c PlatformVaultConnection) (*PlatformResponse, error) {
	adapter, conn, err := s.getPlatformAdapter(c)
	if err != nil {
		return nil, err
	}
	return adapter.List(*conn)
}

func (s *Service) PlatformLogs(request PlatformLogsRequest) (*PlatformResponse, error) {
	adapter, conn, err := s.getPlatformAdapter(request.Connection)
	if err != nil {
		return nil, err
	}
	return adapter.Logs(*conn, request.Name, request.Tail)
}
