package node

func (s *Service) Metrics(c SshConnection) (*MetricsResponse, error) {
	adapter, conn, err := s.getOSAdapter(c)
	if err != nil {
		return nil, err
	}
	return adapter.Metrics(*conn)
}

func (s *Service) DockerDeploy(request DockerRequest) (*DockerResponse, error) {
	adapter, conn, err := s.getOSAdapter(request.Connection)
	if err != nil {
		return nil, err
	}
	return adapter.DockerDeploy(*conn, request.Image, request.Options)
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

func (s *Service) DockerList(c SshConnection) (*DockerResponse, error) {
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
