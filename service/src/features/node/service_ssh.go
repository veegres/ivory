package node

import (
	"ivory/src/plugins/os"
)

func (s *Service) SshCopyId(c SshCredConnection, publicKey string) error {
	adapter, err := s.osRegistry.Get(os.Linux)
	if err != nil {
		return err
	}
	con := s.getSshCredConnection(c)
	return adapter.SshCopyId(*con, publicKey)
}

func (s *Service) Metrics(c SshVaultConnection) (*MetricsResponse, error) {
	adapter, conn, err := s.getOSAdapter(c)
	if err != nil {
		return nil, err
	}
	return adapter.Metrics(*conn)
}

func (s *Service) DockerStop(request DockerRequest) (*DockerResponse, error) {
	adapter, conn, err := s.getOSAdapter(request.Connection)
	if err != nil {
		return nil, err
	}
	return adapter.DockerStop(*conn, request.Container)
}

func (s *Service) DockerRun(request DockerRequest) (*DockerResponse, error) {
	adapter, conn, err := s.getOSAdapter(request.Connection)
	if err != nil {
		return nil, err
	}
	return adapter.DockerRun(*conn, request.Options, request.Image)
}

func (s *Service) DockerDelete(request DockerRequest) (*DockerResponse, error) {
	adapter, conn, err := s.getOSAdapter(request.Connection)
	if err != nil {
		return nil, err
	}
	return adapter.DockerDelete(*conn, request.Container)
}

func (s *Service) DockerList(c SshVaultConnection) (*DockerResponse, error) {
	adapter, conn, err := s.getOSAdapter(c)
	if err != nil {
		return nil, err
	}
	return adapter.DockerList(*conn)
}

func (s *Service) DockerLogs(request DockerLogsRequest) (*DockerResponse, error) {
	adapter, conn, err := s.getOSAdapter(request.Connection)
	if err != nil {
		return nil, err
	}
	return adapter.DockerLogs(*conn, request.Container, request.Tail)
}
