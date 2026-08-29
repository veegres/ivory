package node

import (
	"errors"
	"fmt"
	"ivory/plugins/keeper"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

var ErrKeeperDeployDatabaseCredentialsRequired = errors.New("database credentials are required")
var ErrKeeperDeployKeeperCredentialsRequired = errors.New("keeper credentials are required")

var ErrKeeperDeployPortsRequired = errors.New("keeper, database and ssh ports are required")

// KeeperDeploySpec reports what the deploy forms need to know about the
// engine: its default endpoints and whether it consumes credentials.
func (s *Service) KeeperDeploySpec(r KeeperDeploySpecRequest) (*KeeperDeploySpecResponse, error) {
	metadata, err := s.keeperRegistry.Get(r.Plugin)
	if err != nil {
		return nil, err
	}
	req := metadata.Requirements()
	return &KeeperDeploySpecResponse{
		DbPort:            req.DbPort,
		KeeperPort:        req.KeeperPort,
		KeeperCredentials: req.KeeperCredentials,
		KeeperUser:        req.KeeperUser,
		DbCredentials:     req.DbCredentials,
		DbUser:            req.DbUser,
	}, nil
}

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
// container. It reports failures as log lines rather than an error, so one
// node's bootstrap step can't fail a whole batch that otherwise deployed.
func (s *Service) KeeperPostDeploy(r KeeperDeployRequest) []string {
	if r.PostScript == "" {
		return nil
	}
	if unknown := keeper.UnknownPlaceholders(r.PostScript); len(unknown) > 0 {
		return []string{fmt.Sprintf("post-deploy initialization failed: unknown variables in post script: %s", strings.Join(unknown, ", "))}
	}
	logs, err := s.PlatformContainerExec(PlatformExecRequest{
		Name:       r.Name,
		Connection: r.Connection,
		Vaults:     r.Vaults,
		Command:    r.PostScript,
		Values:     s.getValues(r),
	})
	if err != nil {
		return []string{fmt.Sprintf("post-deploy initialization failed: %v", err)}
	}
	return logs
}

// KeeperDeploy is the self-contained single-node deploy action: run the
// command, then its own post-script.
func (s *Service) KeeperDeploy(r KeeperDeployRequest) ([]string, error) {
	metadata, err := s.keeperRegistry.Get(r.Plugin)
	if err != nil {
		return nil, err
	}
	req := metadata.Requirements()
	if req.KeeperCredentials && r.Vaults.KeeperId == uuid.Nil {
		return nil, ErrKeeperDeployKeeperCredentialsRequired
	}
	if req.DbCredentials && r.Vaults.DatabaseId == uuid.Nil {
		return nil, ErrKeeperDeployDatabaseCredentialsRequired
	}
	if err := s.ValidateKeeperLockedCredentials(req.KeeperUser, r.Vaults.KeeperId); err != nil {
		return nil, err
	}
	if err := s.ValidateKeeperLockedCredentials(req.DbUser, r.Vaults.DatabaseId); err != nil {
		return nil, err
	}

	logs, err := s.KeeperDeployUp(r)
	if err != nil {
		return nil, err
	}
	return append(logs, s.KeeperPostDeploy(r)...), nil
}

// ValidateKeeperLockedCredentials rejects a vault whose username doesn't match
// the one the engine locks itself to. It serves the keeper and the database
// vault alike - each is checked against its own locked username.
func (s *Service) ValidateKeeperLockedCredentials(requiredUser string, vaultId uuid.UUID) error {
	if requiredUser == "" || vaultId == uuid.Nil {
		return nil
	}
	credentials, err := s.vaultService.Get(vaultId)
	if err != nil {
		return err
	}
	if credentials.Username != requiredUser {
		return fmt.Errorf("vault username %q is not allowed: the keeper plugin locks it to %q", credentials.Username, requiredUser)
	}
	return nil
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
