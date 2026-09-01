package node

import (
	"errors"
	"fmt"
	"ivory/clients/console"
	"ivory/core/service/job"
	"ivory/plugins/keeper"
	"ivory/plugins/platform"
	"strings"

	"github.com/google/uuid"
)

var ErrPlatformRequired = errors.New("platform is required")

func (s *Service) PlatformSystemCopyId(r PlatformCopyIdRequest) (string, error) {
	if r.Platform == "" {
		return "", ErrPlatformRequired
	}
	plugin, err := s.platformRegistry.Get(r.Platform)
	if err != nil {
		return "", err
	}
	con := s.getPlatformCredConnection(r.PlatformCredConnection)
	return "ok", plugin.CopyId(con, r.PublicKey)
}

func (s *Service) PlatformSystemMetrics(r PlatformMetricsRequest) (*PlatformMetricsResponse, error) {
	plugin, conn, err := s.getPlatformPlugin(r)
	if err != nil {
		return nil, err
	}
	metrics, err := plugin.Metrics(conn)
	if err != nil {
		return nil, err
	}
	return mapPlatformMetrics(metrics), nil
}

func (s *Service) PlatformSystemProcesses(r PlatformProcessesRequest) (PlatformProcessesResponse, error) {
	plugin, conn, err := s.getPlatformPlugin(r)
	if err != nil {
		return nil, err
	}
	processes, err := plugin.Processes(conn)
	if err != nil {
		return nil, err
	}
	return mapPlatformProcesses(processes), nil
}

func (s *Service) PlatformSystemInfo(r PlatformInfoRequest) (PlatformInfoResponse, error) {
	plugin, conn, err := s.getPlatformPlugin(r)
	if err != nil {
		return nil, err
	}
	info, err := plugin.Info(conn)
	if err != nil {
		return nil, err
	}
	return mapPlatformInfo(info), nil
}

func (s *Service) PlatformSystemLogs(r PlatformLogsRequest, subscriberID job.SubscriberID, close <-chan struct{}, send func(event job.Message)) {
	plugin, conn, err := s.getPlatformPlugin(r.Connection)
	if err != nil {
		send(job.Message{Type: job.SERVER, Message: err.Error()})
		return
	}
	if r.Path == "" {
		send(job.Message{Type: job.SERVER, Message: "path cannot be empty"})
		return
	}
	s.streamCommand(plugin.Logs(conn, r.Path, r.Tail, r.Follow), subscriberID, close, send)
}

func (s *Service) PlatformContainerStop(r PlatformActionRequest) ([]string, error) {
	plugin, conn, err := s.getPlatformPlugin(r.Connection)
	if err != nil {
		return nil, err
	}
	return s.executeCommand(plugin.StopContainer(conn, r.Name))
}

func (s *Service) PlatformContainerRestart(r PlatformActionRequest) ([]string, error) {
	plugin, conn, err := s.getPlatformPlugin(r.Connection)
	if err != nil {
		return nil, err
	}
	return s.executeCommand(plugin.RestartContainer(conn, r.Name))
}

func (s *Service) PlatformContainerStart(r PlatformActionRequest) ([]string, error) {
	plugin, conn, err := s.getPlatformPlugin(r.Connection)
	if err != nil {
		return nil, err
	}
	return s.executeCommand(plugin.StartContainer(conn, r.Name))
}

func (s *Service) PlatformContainerUp(r PlatformUpRequest) ([]string, error) {
	if r.Connection.Host == "" {
		return nil, errors.New("host is empty")
	}

	plugin, conn, err := s.getPlatformPlugin(r.Connection)
	if err != nil {
		return nil, err
	}

	values, err := s.getExecutionValues(r.Connection.Host, r.Vaults, r.Values)
	if err != nil {
		return nil, err
	}
	command, err := resolveCommand(r.Command, values)
	if err != nil {
		return nil, err
	}
	return s.executeCommand(plugin.UpContainer(conn, command))
}

func (s *Service) PlatformContainerExec(r PlatformExecRequest) ([]string, error) {
	plugin, conn, err := s.getPlatformPlugin(r.Connection)
	if err != nil {
		return nil, err
	}

	values, err := s.getExecutionValues(r.Connection.Host, r.Vaults, r.Values)
	if err != nil {
		return nil, err
	}

	return s.execContainerCommand(plugin, conn, r.Name, r.Command, values)
}

func (s *Service) execContainerCommand(plugin platform.Plugin, conn platform.Connection, name string, commandTemplate string, values keeper.Values) ([]string, error) {
	command, err := resolveCommand(commandTemplate, values)
	if err != nil {
		return nil, err
	}
	return s.executeCommand(plugin.ExecContainer(conn, name, command))
}

func resolveCommand(template string, values keeper.Values) ([]string, error) {
	command := platform.SplitCommand(template)
	for i, argument := range command {
		command[i] = keeper.Interpolate(argument, values)
	}
	if unresolved := keeper.UnresolvedPlaceholders(strings.Join(command, " ")); len(unresolved) > 0 {
		return nil, fmt.Errorf("missing values for placeholders: %s", strings.Join(unresolved, ", "))
	}
	return command, nil
}

func (s *Service) getExecutionValues(host string, vaults Vaults, values keeper.Values) (keeper.Values, error) {
	values.Host = host

	// NOTE: either vault is optional for keeper plugins that consume no such
	// credentials, unused values just keep their placeholders unresolved
	if vaults.KeeperId != uuid.Nil {
		keeperCred, err := s.vaultService.GetDecrypted(vaults.KeeperId)
		if err != nil {
			return values, fmt.Errorf("failed to get keeper credentials from vault: %v", err)
		}
		values.KeeperUser = keeperCred.Username
		values.KeeperPass = keeperCred.Secret
	}
	if vaults.DatabaseId != uuid.Nil {
		dbCred, err := s.vaultService.GetDecrypted(vaults.DatabaseId)
		if err != nil {
			return values, fmt.Errorf("failed to get database credentials from vault: %v", err)
		}
		values.DbUser = dbCred.Username
		values.DbPass = dbCred.Secret
	}
	return values, nil
}

func (s *Service) PlatformContainerDown(r PlatformActionRequest) ([]string, error) {
	plugin, conn, err := s.getPlatformPlugin(r.Connection)
	if err != nil {
		return nil, err
	}
	stopLogs, _ := s.executeCommand(plugin.StopContainer(conn, r.Name))
	downLogs, err := s.executeCommand(plugin.DownContainer(conn, r.Name))
	if err != nil {
		return nil, err
	}
	return append(stopLogs, downLogs...), nil
}

func (s *Service) PlatformContainerList(c PlatformVaultConnection) ([]string, error) {
	plugin, conn, err := s.getPlatformPlugin(c)
	if err != nil {
		return nil, err
	}
	return s.executeCommand(plugin.ListContainer(conn))
}

func (s *Service) PlatformContainerMetrics(r PlatformActionRequest) (*PlatformMetricsResponse, error) {
	plugin, conn, err := s.getPlatformPlugin(r.Connection)
	if err != nil {
		return nil, err
	}
	metrics, err := plugin.MetricsContainer(conn, r.Name)
	if err != nil {
		return nil, err
	}
	return mapPlatformMetrics(metrics), nil
}

func (s *Service) PlatformContainerLogs(r PlatformLogsRequest, subscriberID job.SubscriberID, close <-chan struct{}, send func(event job.Message)) {
	plugin, conn, err := s.getPlatformPlugin(r.Connection)
	if err != nil {
		send(job.Message{Type: job.SERVER, Message: err.Error()})
		return
	}
	if r.Path == "" {
		send(job.Message{Type: job.SERVER, Message: "path cannot be empty"})
		return
	}
	s.streamCommand(plugin.LogsContainer(conn, r.Path, r.Tail, r.Follow), subscriberID, close, send)
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
