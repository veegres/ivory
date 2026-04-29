package node

import (
	"fmt"
	"ivory/src/plugins/os"
	"strconv"
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
	// Pull then run
	command := fmt.Sprintf("pull %s && run -d %s %s", request.Image, request.Options, request.Image)
	res, err := s.sshClient.ExecuteDocker(*connection, command)
	if err != nil {
		return nil, err
	}
	return &DockerResult{Stdout: res.Stdout, Stderr: res.Stderr, ExitCode: res.ExitCode}, nil
}

func (s *Service) DockerStop(request DockerRequest) (*DockerResult, error) {
	connection, err := s.getSshConnection(request.Connection)
	if err != nil {
		return nil, err
	}
	res, err := s.sshClient.ExecuteDocker(*connection, "stop "+request.Container)
	if err != nil {
		return nil, err
	}
	return &DockerResult{Stdout: res.Stdout, Stderr: res.Stderr, ExitCode: res.ExitCode}, nil
}

func (s *Service) DockerRun(request DockerRequest) (*DockerResult, error) {
	connection, err := s.getSshConnection(request.Connection)
	if err != nil {
		return nil, err
	}
	res, err := s.sshClient.ExecuteDocker(*connection, fmt.Sprintf("run -d %s %s", request.Options, request.Image))
	if err != nil {
		return nil, err
	}
	return &DockerResult{Stdout: res.Stdout, Stderr: res.Stderr, ExitCode: res.ExitCode}, nil
}

func (s *Service) DockerDelete(request DockerRequest) (*DockerResult, error) {
	connection, err := s.getSshConnection(request.Connection)
	if err != nil {
		return nil, err
	}
	res, err := s.sshClient.ExecuteDocker(*connection, "rm "+request.Container)
	if err != nil {
		return nil, err
	}
	return &DockerResult{Stdout: res.Stdout, Stderr: res.Stderr, ExitCode: res.ExitCode}, nil
}

func (s *Service) DockerList(request SshRequest) (*DockerResult, error) {
	connection, err := s.getSshConnection(request.Connection)
	if err != nil {
		return nil, err
	}
	res, err := s.sshClient.ExecuteDocker(*connection, "ps -a")
	if err != nil {
		return nil, err
	}
	return &DockerResult{Stdout: res.Stdout, Stderr: res.Stderr, ExitCode: res.ExitCode}, nil
}

func (s *Service) DockerLogs(request DockerLogsRequest) (*DockerResult, error) {
	connection, err := s.getSshConnection(request.Connection)
	if err != nil {
		return nil, err
	}
	command := "logs "
	if request.Tail > 0 {
		command += "--tail " + strconv.Itoa(request.Tail) + " "
	}
	command += request.Container
	res, err := s.sshClient.ExecuteDocker(*connection, command)
	if err != nil {
		return nil, err
	}
	return &DockerResult{Stdout: res.Stdout, Stderr: res.Stderr, ExitCode: res.ExitCode}, nil
}
