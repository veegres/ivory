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

func (s *Service) PlatformSystemCopyId(r PlatformCopyIdRequest) (string, error) {
	adapter, err := s.platformRegistry.Get(DefaultPlatform)
	if err != nil {
		return "", err
	}
	con := s.getPlatformCredConnection(r.PlatformCredConnection)
	return "ok", adapter.CopyId(con, r.PublicKey)
}

func (s *Service) PlatformSystemMetrics(r PlatformMetricsRequest) (*PlatformMetricsResponse, error) {
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

func (s *Service) PlatformSystemProcesses(r PlatformProcessesRequest) (PlatformProcessesResponse, error) {
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

func (s *Service) PlatformSystemInfo(r PlatformInfoRequest) (PlatformInfoResponse, error) {
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

func (s *Service) PlatformSystemLogs(r PlatformLogsRequest, subscriberID job.SubscriberID, close <-chan struct{}, send func(event job.Message)) {
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
	command, err := resolveCommand(r.Command, values)
	if err != nil {
		return nil, err
	}
	return s.executeCommand(adapter.UpContainer(conn, command))
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

	return s.execContainerCommand(adapter, conn, r.Name, r.Command, values)
}

// execContainerCommand interpolates and runs one command against an
// already-resolved adapter/connection/values set, so a caller that already
// resolved vault credentials (e.g. KeeperPostDeploy) doesn't re-fetch and
// re-decrypt them just to run its command.
func (s *Service) execContainerCommand(adapter platform.Adapter, conn platform.Connection, name string, commandTemplate string, values keeper.Values) ([]string, error) {
	command, err := resolveCommand(commandTemplate, values)
	if err != nil {
		return nil, err
	}
	return s.executeCommand(adapter.ExecContainer(conn, name, command))
}

// resolveCommand turns a template into the exact arguments that will run. The
// split comes first and the values go into arguments that are already
// separated, so a credential can never introduce an argument boundary or close
// a span the template author opened - which is what makes escaping
// unnecessary. The platform is handed the finished command and never sees a
// placeholder.
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

// getExecutionValues finalizes interpolation values for execution: the host
// comes from the connection and the keeper/database credentials from their own
// vaults, so they cannot be spoofed through request values.
//
// Nothing is escaped on the way out. The adapter interpolates into an argument
// that has already been split off, so a value is never seen by a parser and
// there is nothing for an escape to protect it from.
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

// PlatformContainerDown removes a deployment, stopping it first. The stop is
// what makes Down usable at all: removing a running deployment is refused by
// the platform, so Down used to fail on exactly the deployments a user wants
// to remove. It is a graceful stop rather than a forced removal because these
// are databases, and the platform's own shutdown timeout is worth waiting for.
// A stop failure is not fatal - an already-stopped deployment still removes.
func (s *Service) PlatformContainerDown(r PlatformActionRequest) ([]string, error) {
	adapter, conn, err := s.getPlatformAdapter(r.Connection)
	if err != nil {
		return nil, err
	}
	stopLogs, _ := s.executeCommand(adapter.StopContainer(conn, r.Name))
	downLogs, err := s.executeCommand(adapter.DownContainer(conn, r.Name))
	if err != nil {
		return nil, err
	}
	return append(stopLogs, downLogs...), nil
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
