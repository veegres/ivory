package node

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
