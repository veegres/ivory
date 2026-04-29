package node

import (
	"ivory/src/features/cert"
	"ivory/src/plugins/keeper"
	"ivory/src/plugins/os"

	"github.com/google/uuid"
)

// COMMON (WEB AND SERVER)

type Connection struct {
	Host       string     `json:"host" form:"host"`
	SshKeyId   *uuid.UUID `json:"sshKeyId,omitempty" form:"sshKeyId"`
	SshPort    int        `json:"sshPort,omitempty" form:"sshPort"`
	KeeperPort int        `json:"keeperPort,omitempty" form:"keeperPort"`
	DbPort     int        `json:"dbPort,omitempty" form:"dbPort"`
}

type KeeperRequest struct {
	Type    keeper.Type `json:"type" form:"type"`
	Host    string      `json:"host" form:"host"`
	Port    int         `json:"port" form:"port"`
	VaultId *uuid.UUID  `json:"vaultId" form:"vaultId"`
	Certs   *cert.Certs `json:"certs" form:"certs"`
	Body    any         `json:"body" form:"body"`
}

type KeeperResponse struct {
	Connection Connection      `json:"connection"`
	Keeper     keeper.Response `json:"keeper"`
}

type SshRequest struct {
	Connection Connection `json:"connection" form:"connection" binding:"required"`
}

type SshResponseMetrics = os.Metrics

type DockerRequest struct {
	Connection Connection `json:"connection" form:"connection" binding:"required"`
	Image      string     `json:"image" form:"image"`
	Container  string     `json:"container" form:"container"`
	Options    string     `json:"options" form:"options"`
}

type DockerLogsRequest struct {
	Connection Connection `json:"connection" form:"connection" binding:"required"`
	Container  string     `json:"container" form:"container" binding:"required"`
	Tail       int        `json:"tail" form:"tail"`
}

type DockerResult struct {
	Stdout   string `json:"stdout"`
	Stderr   string `json:"stderr"`
	ExitCode int    `json:"exitCode"`
}

// SPECIFIC (SERVER)

type KeeperAutoRequest struct {
	Connections []Connection `json:"connections" form:"connections"`
	Type        keeper.Type  `json:"type" form:"type"`
	VaultId     *uuid.UUID   `json:"vaultId" form:"vaultId"`
	Certs       *cert.Certs  `json:"certs" form:"certs"`
	Body        any          `json:"body" form:"body"`
}
