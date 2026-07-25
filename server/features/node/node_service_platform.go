package node

import (
	"errors"
	"fmt"
	"ivory/clients/console"
	"ivory/core/service/job"
	"ivory/plugins/keeper"
	"ivory/plugins/platform"
	"maps"
	"strings"

	"github.com/google/uuid"
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

func (s *Service) PlatformVmInfo(r PlatformInfoRequest) (PlatformInfoResponse, error) {
	adapter, conn, err := s.getPlatformAdapter(r)
	if err != nil {
		return nil, err
	}
	info, err := adapter.Info(conn)
	if err != nil {
		return nil, err
	}
	return mapPlatformInfo(info), nil
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

func (s *Service) PlatformContainerUp(r PlatformUpRequest) ([]string, error) {
	if r.Connection.Host == "" {
		return nil, errors.New("host is empty")
	}

	adapter, conn, err := s.getPlatformAdapter(r.Connection)
	if err != nil {
		return nil, err
	}

	values, err := s.getExecutionValues(r.Connection.Host, r.Vaults, r.Values)
	if err != nil {
		return nil, err
	}

	options := s.normalizeDatabaseOptions(keeper.Interpolate(r.Options, values))
	if unresolved := keeper.UnresolvedPlaceholders(options); len(unresolved) > 0 {
		return nil, fmt.Errorf("missing values for placeholders: %s", strings.Join(unresolved, ", "))
	}

	// NOTE: unlike options, entryScript is not run through
	// normalizeDatabaseOptions - it can be a multi-line shell script (see
	// keeper.DeploymentSpec.EntryScript) whose newlines are meaningful
	// statement separators (e.g. before "then"/"fi"), not just formatting.
	entryScript := keeper.Interpolate(r.EntryScript, values)
	if unresolved := keeper.UnresolvedPlaceholders(entryScript); len(unresolved) > 0 {
		return nil, fmt.Errorf("missing values for placeholders: %s", strings.Join(unresolved, ", "))
	}

	return s.executeCommand(adapter.UpContainer(conn, options, r.Image, entryScript))
}

// PlatformContainerExec runs one command inside the named deployment,
// interpolating the command template exactly like PlatformContainerUp does
// for options (host and vault credentials are injected server-side).
func (s *Service) PlatformContainerExec(r PlatformExecRequest) ([]string, error) {
	adapter, conn, err := s.getPlatformAdapter(r.Connection)
	if err != nil {
		return nil, err
	}

	values, err := s.getExecutionValues(r.Connection.Host, r.Vaults, r.Values)
	if err != nil {
		return nil, err
	}

	command := keeper.Interpolate(r.Command, values)
	if unresolved := keeper.UnresolvedPlaceholders(command); len(unresolved) > 0 {
		return nil, fmt.Errorf("missing values for placeholders: %s", strings.Join(unresolved, ", "))
	}

	return s.executeCommand(adapter.ExecContainer(conn, r.Name, command))
}

// getExecutionValues finalizes interpolation values for execution: the host
// comes from the connection and the database credentials from the vault, so
// they cannot be spoofed through request values.
func (s *Service) getExecutionValues(host string, vaults Vaults, requestValues map[string]string) (map[string]string, error) {
	values := make(map[string]string, len(requestValues)+3)
	maps.Copy(values, requestValues)
	values[string(keeper.VarHost)] = host

	// NOTE: the vault is optional for keeper plugins that consume no database
	// credentials, unused values just keep their placeholders unresolved
	if vaults.DatabaseId != uuid.Nil {
		dbCred, err := s.vaultService.GetDecrypted(vaults.DatabaseId)
		if err != nil {
			return nil, fmt.Errorf("failed to get database credentials from vault: %v", err)
		}
		values[string(keeper.VarDbUser)] = dbCred.Username
		values[string(keeper.VarDbPass)] = dbCred.Secret
	}
	return values, nil
}

func (s *Service) normalizeDatabaseOptions(options string) string {
	return strings.Join(strings.Fields(options), " ")
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
