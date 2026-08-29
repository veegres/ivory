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

	// NOTE: the command is not whitespace-normalized here - it may embed a
	// multi-line startup script whose newlines are real statement separators.
	// The adapter collapses only the whitespace outside quoted spans.
	command := keeper.Interpolate(r.Command, values)
	if unresolved := keeper.UnresolvedPlaceholders(command); len(unresolved) > 0 {
		return nil, fmt.Errorf("missing values for placeholders: %s", strings.Join(unresolved, ", "))
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
	command := keeper.Interpolate(commandTemplate, values)
	if unresolved := keeper.UnresolvedPlaceholders(command); len(unresolved) > 0 {
		return nil, fmt.Errorf("missing values for placeholders: %s", strings.Join(unresolved, ", "))
	}

	return s.executeCommand(adapter.ExecContainer(conn, name, command))
}

// getExecutionValues finalizes interpolation values for execution: the host
// comes from the connection and the keeper/database credentials from their own
// vaults, so they cannot be spoofed through request values.
func (s *Service) getExecutionValues(host string, vaults Vaults, values keeper.Values) (keeper.Values, error) {
	values.Host = host

	// NOTE: either vault is optional for keeper plugins that consume no such
	// credentials, unused values just keep their placeholders unresolved
	if vaults.KeeperId != uuid.Nil {
		keeperCred, err := s.vaultService.GetDecrypted(vaults.KeeperId)
		if err != nil {
			return values, fmt.Errorf("failed to get keeper credentials from vault: %v", err)
		}
		values.KeeperUser = escapeInterpolatedValue(keeperCred.Username)
		values.KeeperPass = escapeInterpolatedValue(keeperCred.Secret)
	}
	if vaults.DatabaseId != uuid.Nil {
		dbCred, err := s.vaultService.GetDecrypted(vaults.DatabaseId)
		if err != nil {
			return values, fmt.Errorf("failed to get database credentials from vault: %v", err)
		}
		values.DbUser = escapeInterpolatedValue(dbCred.Username)
		values.DbPass = escapeInterpolatedValue(dbCred.Secret)
	}
	return values, nil
}

// escapeInterpolatedValue neutralizes characters that could break out of a
// keeper plugin's own hand-written quoting once substituted into a command or
// post-script template (e.g. an etcd post-script wraps {{dbUser}}:{{dbPass}}
// in literal quotes - an unescaped quote in the value corrupts the command
// it's exec'd with). Each one is marked with platform.ValueEscape rather than
// a backslash, so splitShellFields can tell a credential's quote from a
// backslash the template author wrote: escaping both the same way is what
// used to strip a post script's own `\"` on the way to the inner shell.
//
// A value never carries the marker itself, so any occurrence is dropped
// before escaping rather than trusted.
func escapeInterpolatedValue(v string) string {
	v = strings.ReplaceAll(v, string(platform.ValueEscape), "")
	var escaped strings.Builder
	for _, r := range v {
		if r == '\'' || r == '"' || r == '\\' {
			escaped.WriteRune(platform.ValueEscape)
		}
		escaped.WriteRune(r)
	}
	return escaped.String()
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
