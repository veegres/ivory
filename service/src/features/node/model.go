package node

import (
	"ivory/src/features/cert"
	"ivory/src/plugins/keeper"
	"ivory/src/plugins/os"

	"github.com/google/uuid"
)

// COMMON (WEB AND SERVER)

type KeeperConnection struct {
	Host string `json:"host" form:"host"`
	Port int    `json:"port" form:"port"`
}

type KeeperRequest struct {
	KeeperConnection
	Plugin  keeper.Plugin `json:"plugin" form:"plugin"`
	VaultId *uuid.UUID    `json:"vaultId" form:"vaultId"`
	Certs   *cert.Certs   `json:"certs" form:"certs"`
	Body    any           `json:"body" form:"body"`
}

type KeeperResponse = keeper.Response

type SshConnection struct {
	Host    string     `json:"host" form:"host"`
	Port    int        `json:"port" form:"port"`
	VaultId *uuid.UUID `json:"vaultId" form:"vaultId"`
}

type MetricsResponse = os.Metrics

type DockerRequest struct {
	Connection SshConnection `json:"connection" form:"connection" binding:"required"`
	Image      string        `json:"image" form:"image"`
	Container  string        `json:"container" form:"container"`
	Options    string        `json:"options" form:"options"`
}

type DockerLogsRequest struct {
	Connection SshConnection `json:"connection" form:"connection" binding:"required"`
	Container  string        `json:"container" form:"container" binding:"required"`
	Tail       int           `json:"tail" form:"tail"`
}

type DockerResponse = os.Docker

// SPECIFIC (SERVER)

type KeeperAutoRequest struct {
	Connections []KeeperConnection `json:"connections" form:"connections"`
	Plugin      keeper.Plugin      `json:"plugin" form:"plugin"`
	VaultId     *uuid.UUID         `json:"vaultId" form:"vaultId"`
	Certs       *cert.Certs        `json:"certs" form:"certs"`
	Body        any                `json:"body" form:"body"`
}
