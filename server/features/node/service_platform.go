package node

import (
	"encoding/json"
	"errors"
	"fmt"
	"ivory/clients/console"
	"ivory/core/service/job"
	"ivory/plugins/platform"
	"regexp"
	"strconv"
	"strings"
)

func (s *Service) PlatformVmCopyId(r PlatformCopyIdRequest) (string, error) {
	adapter, err := s.platformRegistry.Get(platform.Linux)
	if err != nil {
		return "", err
	}
	con := s.getPlatformCredConnection(r.PlatformCredConnection)
	return "ok", adapter.CopyId(con, r.PublicKey)
}

func (s *Service) PlatformVmMetrics(r PlatformMetricsRequest) (*PlatformMetricsResponse, error) {
	adapter, conn, err := s.getPlatformAdapter(r)
	if err != nil {
		return nil, err
	}
	metrics, err := adapter.Metrics(conn)
	if err != nil {
		return nil, err
	}
	return mapPlatformMetrics(metrics), nil
}

func (s *Service) PlatformVmProcesses(r PlatformProcessesRequest) (PlatformProcessesResponse, error) {
	adapter, conn, err := s.getPlatformAdapter(r)
	if err != nil {
		return nil, err
	}
	processes, err := adapter.Processes(conn)
	if err != nil {
		return nil, err
	}
	return mapPlatformProcesses(processes), nil
}

func (s *Service) PlatformLogs(r PlatformLogsRequest, subscriberID job.SubscriberID, close <-chan struct{}, send func(event job.Message)) {
	adapter, conn, err := s.getPlatformAdapter(r.Connection)
	if err != nil {
		send(job.Message{Type: job.SERVER, Message: err.Error()})
		return
	}
	if r.Path == "" {
		send(job.Message{Type: job.SERVER, Message: "path cannot be empty"})
		return
	}
	s.streamCommand(adapter.Logs(conn, r.Path, r.Tail, r.Follow), subscriberID, close, send)
}

func (s *Service) PlatformContainerStop(r PlatformActionRequest) ([]string, error) {
	adapter, conn, err := s.getPlatformAdapter(r.Connection)
	if err != nil {
		return nil, err
	}
	return s.executeCommand(adapter.StopContainer(conn, r.Name))
}

func (s *Service) PlatformContainerRestart(r PlatformActionRequest) ([]string, error) {
	adapter, conn, err := s.getPlatformAdapter(r.Connection)
	if err != nil {
		return nil, err
	}
	return s.executeCommand(adapter.RestartContainer(conn, r.Name))
}

func (s *Service) PlatformContainerStart(r PlatformActionRequest) ([]string, error) {
	adapter, conn, err := s.getPlatformAdapter(r.Connection)
	if err != nil {
		return nil, err
	}
	return s.executeCommand(adapter.StartContainer(conn, r.Name))
}

func (s *Service) PlatformContainerDeployOptions(r PlatformDeployOptionsRequest) (*PlatformDeployOptionsResponse, error) {
	metadata, err := s.keeperMetadataRegistry.Get(r.Plugin)
	if err != nil {
		return nil, err
	}
	adapter, err := s.platformRegistry.Get(platform.Linux)
	if err != nil {
		return nil, err
	}
	spec := metadata.DeploymentSpec()
	return &PlatformDeployOptionsResponse{
		Uri:               spec.DefaultImage,
		DefaultValues:     spec.DefaultValues,
		Options:           adapter.RenderOptions(mapKeeperDeploymentToPlatformSpec(spec, false)),
		OptionsSingleHost: adapter.RenderOptions(mapKeeperDeploymentToPlatformSpec(spec, true)),
	}, nil
}

func (s *Service) PlatformContainerUp(r PlatformUpRequest) ([]string, error) {
	if r.Connection.Host == "" {
		return nil, errors.New("host is empty")
	}

	adapter, conn, err := s.getPlatformAdapter(r.Connection)
	if err != nil {
		return nil, err
	}

	values := ImageOptions{}
	values.Host = r.Connection.Host
	values.Cluster = r.ImageOptions.Cluster
	values.Dcs = r.ImageOptions.Dcs
	values.KeeperPort = strconv.Itoa(r.ImageOptions.KeeperPort)
	values.DbPort = strconv.Itoa(r.ImageOptions.DbPort)

	// 1. Handle DB Password
	dbCred, errDb := s.vaultService.GetDecrypted(r.Vaults.DatabaseId)
	if errDb != nil {
		return nil, fmt.Errorf("failed to get database credentials from vault: %v", errDb)
	}
	values.DbUser = dbCred.Username
	values.DbPass = dbCred.Secret

	interpolatedString, errString := s.getInterpolatedString(r.RawImageOptions, values)
	if errString != nil {
		return nil, errString
	}
	options := s.normalizeDatabaseOptions(interpolatedString)

	return s.executeCommand(adapter.UpContainer(conn, options, r.Image))
}

func (s *Service) normalizeDatabaseOptions(options string) string {
	return strings.Join(strings.Fields(options), " ")
}

func (s *Service) getInterpolatedString(template string, values ImageOptions) (string, error) {
	// 1. Convert struct to JSON bytes
	bytes, _ := json.Marshal(values)
	// 2. Unmarshal bytes into a map
	var resultMap map[string]string
	err := json.Unmarshal(bytes, &resultMap)
	if err != nil {
		return "", err
	}
	return regexp.MustCompile(`{{(\w+)}}`).ReplaceAllStringFunc(template, func(match string) string {
		key := match[2 : len(match)-2]
		if val, ok := resultMap[key]; ok {
			return val
		}
		return match
	}), nil
}

func (s *Service) PlatformContainerDown(r PlatformActionRequest) ([]string, error) {
	adapter, conn, err := s.getPlatformAdapter(r.Connection)
	if err != nil {
		return nil, err
	}
	return s.executeCommand(adapter.DownContainer(conn, r.Name))
}

func (s *Service) PlatformContainerList(c PlatformVaultConnection) ([]string, error) {
	adapter, conn, err := s.getPlatformAdapter(c)
	if err != nil {
		return nil, err
	}
	return s.executeCommand(adapter.ListContainer(conn))
}

func (s *Service) PlatformContainerMetrics(r PlatformActionRequest) (*PlatformMetricsResponse, error) {
	adapter, conn, err := s.getPlatformAdapter(r.Connection)
	if err != nil {
		return nil, err
	}
	metrics, err := adapter.MetricsContainer(conn, r.Name)
	if err != nil {
		return nil, err
	}
	return mapPlatformMetrics(metrics), nil
}

func (s *Service) PlatformContainerLogs(r PlatformLogsRequest, subscriberID job.SubscriberID, close <-chan struct{}, send func(event job.Message)) {
	adapter, conn, err := s.getPlatformAdapter(r.Connection)
	if err != nil {
		send(job.Message{Type: job.SERVER, Message: err.Error()})
		return
	}
	if r.Path == "" {
		send(job.Message{Type: job.SERVER, Message: "path cannot be empty"})
		return
	}
	s.streamCommand(adapter.LogsContainer(conn, r.Path, r.Tail, r.Follow), subscriberID, close, send)
}

func (s *Service) streamCommand(cmd console.Command, subscriberID job.SubscriberID, close <-chan struct{}, send func(event job.Message)) {
	jobID, err := s.jobManager.Start(cmd)
	if err != nil {
		send(job.Message{Type: job.SERVER, Message: err.Error()})
		return
	}
	s.jobManager.Stream(jobID, subscriberID, close, send)
}

func (s *Service) executeCommand(cmd console.Command) ([]string, error) {
	return cmd.Execute()
}
