package node

import (
	"fmt"
	"strconv"
)

func (s *Service) Metrics(request SshRequest) (*SshMetrics, int, error) {
	connection, err := s.getSshConnection(request.VmId)
	if err != nil {
		return nil, 0, err
	}
	result, err := s.sshClient.Metrics(*connection)
	if err != nil {
		return nil, 0, err
	}
	return result, 200, nil
}

func (s *Service) DockerDeploy(request DockerRequest) (*DockerResult, int, error) {
	connection, err := s.getSshConnection(request.VmId)
	if err != nil {
		return nil, 0, err
	}
	// Pull then run
	command := fmt.Sprintf("pull %s && run -d %s %s", request.Image, request.Options, request.Image)
	res, err := s.sshClient.ExecuteDocker(*connection, command)
	if err != nil {
		return nil, 0, err
	}
	return &DockerResult{Stdout: res.Stdout, Stderr: res.Stderr, ExitCode: res.ExitCode}, 200, nil
}

func (s *Service) DockerStop(request DockerRequest) (*DockerResult, int, error) {
	connection, err := s.getSshConnection(request.VmId)
	if err != nil {
		return nil, 0, err
	}
	res, err := s.sshClient.ExecuteDocker(*connection, "stop "+request.Container)
	if err != nil {
		return nil, 0, err
	}
	return &DockerResult{Stdout: res.Stdout, Stderr: res.Stderr, ExitCode: res.ExitCode}, 200, nil
}

func (s *Service) DockerRun(request DockerRequest) (*DockerResult, int, error) {
	connection, err := s.getSshConnection(request.VmId)
	if err != nil {
		return nil, 0, err
	}
	res, err := s.sshClient.ExecuteDocker(*connection, fmt.Sprintf("run -d %s %s", request.Options, request.Image))
	if err != nil {
		return nil, 0, err
	}
	return &DockerResult{Stdout: res.Stdout, Stderr: res.Stderr, ExitCode: res.ExitCode}, 200, nil
}

func (s *Service) DockerDelete(request DockerRequest) (*DockerResult, int, error) {
	connection, err := s.getSshConnection(request.VmId)
	if err != nil {
		return nil, 0, err
	}
	res, err := s.sshClient.ExecuteDocker(*connection, "rm "+request.Container)
	if err != nil {
		return nil, 0, err
	}
	return &DockerResult{Stdout: res.Stdout, Stderr: res.Stderr, ExitCode: res.ExitCode}, 200, nil
}

func (s *Service) DockerList(request SshRequest) (*DockerResult, int, error) {
	connection, err := s.getSshConnection(request.VmId)
	if err != nil {
		return nil, 0, err
	}
	res, err := s.sshClient.ExecuteDocker(*connection, "ps -a")
	if err != nil {
		return nil, 0, err
	}
	return &DockerResult{Stdout: res.Stdout, Stderr: res.Stderr, ExitCode: res.ExitCode}, 200, nil
}

func (s *Service) DockerLogs(request DockerLogsRequest) (*DockerResult, int, error) {
	connection, err := s.getSshConnection(request.VmId)
	if err != nil {
		return nil, 0, err
	}
	command := "logs "
	if request.Tail > 0 {
		command += "--tail " + strconv.Itoa(request.Tail) + " "
	}
	command += request.Container
	res, err := s.sshClient.ExecuteDocker(*connection, command)
	if err != nil {
		return nil, 0, err
	}
	return &DockerResult{Stdout: res.Stdout, Stderr: res.Stderr, ExitCode: res.ExitCode}, 200, nil
}
