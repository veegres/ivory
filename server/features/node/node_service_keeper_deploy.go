package node

import (
	"errors"
	"fmt"
	"ivory/plugins/keeper"
	"strconv"
	"strings"
)

var ErrKeeperDeployPortsRequired = errors.New("keeper, database and ssh ports are required")

// KeeperDeployUp deploys one node by running its own command. Every value the
// command can reference belongs to this request alone, so one node's values
// can never reach another node's command.
func (s *Service) KeeperDeployUp(r KeeperDeployRequest) ([]string, error) {
	if r.Connection.Host == "" {
		return nil, errors.New("host not provided for node")
	}
	if r.KeeperPort <= 0 || r.DbPort <= 0 || r.Connection.Port <= 0 {
		return nil, ErrKeeperDeployPortsRequired
	}
	if unknown := keeper.UnknownPlaceholders(r.Command); len(unknown) > 0 {
		return nil, fmt.Errorf("unknown variables in command: %s", strings.Join(unknown, ", "))
	}
	return s.PlatformContainerUp(PlatformUpRequest{
		Connection: r.Connection,
		Vaults:     r.Vaults,
		Command:    r.Command,
		Values:     s.getValues(r),
	})
}

// KeeperPostDeploy runs a deployment's post-script inside its already-running
// container. It keeps reporting the failure as a log line so one node's
// bootstrap step can't fail a whole batch that otherwise deployed, and returns
// it as well so the caller can say so in its result: a post script that
// silently does nothing leaves a cluster that looks deployed and is not
// initialized.
func (s *Service) KeeperPostDeploy(r KeeperDeployRequest) ([]string, error) {
	if len(r.PostScripts) == 0 {
		return nil, nil
	}
	// NOTE: the steps run as separate executions rather than one chained
	// script, so a deployment whose image has no shell can still be
	// initialized - etcd's holds only etcd, etcdctl and etcdutl
	logs := make([]string, 0, len(r.PostScripts))
	for _, script := range r.PostScripts {
		if unknown := keeper.UnknownPlaceholders(script); len(unknown) > 0 {
			err := fmt.Errorf("unknown variables in post script: %s", strings.Join(unknown, ", "))
			return append(logs, fmt.Sprintf("post-deploy initialization failed: %v", err)), err
		}
		stepLogs, err := s.PlatformContainerExec(PlatformExecRequest{
			Name:       r.Name,
			Connection: r.Connection,
			Vaults:     r.Vaults,
			Command:    script,
			Values:     s.getValues(r),
		})
		logs = append(logs, stepLogs...)
		if err != nil {
			return append(logs, fmt.Sprintf("post-deploy initialization failed: %v", err)), err
		}
	}
	return logs, nil
}

// KeeperDeploy is the self-contained single-node deploy action: run the
// command, then its own post-script.
func (s *Service) KeeperDeploy(r KeeperDeployRequest) ([]string, error) {
	if err := s.ValidateKeeperPlugin(r.Plugin); err != nil {
		return nil, err
	}
	logs, err := s.KeeperDeployUp(r)
	if err != nil {
		return nil, err
	}
	postLogs, postErr := s.KeeperPostDeploy(r)
	logs = append(logs, postLogs...)
	if postErr != nil {
		logs = append(logs, "the node is running, but its post-deploy initialization did not complete")
	}
	return logs, nil
}

// getValues builds one command's complete interpolation scope from its own
// request. Host and ssh port come from the connection rather than the body so
// they always describe the machine actually being deployed to; the credentials
// are left out here and resolved from the vault at execution time.
func (s *Service) getValues(r KeeperDeployRequest) keeper.Values {
	return keeper.Values{
		Cluster:    r.Cluster,
		Name:       r.Name,
		Host:       r.Connection.Host,
		SshPort:    strconv.Itoa(r.Connection.Port),
		KeeperPort: strconv.Itoa(r.KeeperPort),
		DbPort:     strconv.Itoa(r.DbPort),
	}
}
