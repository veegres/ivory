package node

func (s *Service) Metrics(request SshRequest) (*SshResponseMetrics, error) {
	adapter, conn, err := s.getOSAdapter(request.Connection)
	if err != nil {
		return nil, err
	}
	return adapter.Metrics(*conn)
}

func (s *Service) DockerDeploy(request DockerRequest) (*DockerResult, error) {
	adapter, conn, err := s.getOSAdapter(request.Connection)
	if err != nil {
		return nil, err
	}
	return adapter.DockerDeploy(*conn, request.Image, request.Options)
}

func (s *Service) DockerStop(request DockerRequest) (*DockerResult, error) {
	adapter, conn, err := s.getOSAdapter(request.Connection)
	if err != nil {
		return nil, err
	}
	return adapter.DockerStop(*conn, request.Container)
}

func (s *Service) DockerRun(request DockerRequest) (*DockerResult, error) {
	adapter, conn, err := s.getOSAdapter(request.Connection)
	if err != nil {
		return nil, err
	}
	return adapter.DockerRun(*conn, request.Options, request.Image)
}

func (s *Service) DockerDelete(request DockerRequest) (*DockerResult, error) {
	adapter, conn, err := s.getOSAdapter(request.Connection)
	if err != nil {
		return nil, err
	}
	return adapter.DockerDelete(*conn, request.Container)
}

func (s *Service) DockerList(request SshRequest) (*DockerResult, error) {
	adapter, conn, err := s.getOSAdapter(request.Connection)
	if err != nil {
		return nil, err
	}
	return adapter.DockerList(*conn)
}

func (s *Service) DockerLogs(request DockerLogsRequest) (*DockerResult, error) {
	adapter, conn, err := s.getOSAdapter(request.Connection)
	if err != nil {
		return nil, err
	}
	return adapter.DockerLogs(*conn, request.Container, request.Tail)
}
