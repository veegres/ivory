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

type KeeperOptions struct {
	Plugin  keeper.Plugin `json:"plugin" form:"plugin"`
	VaultId *uuid.UUID    `json:"vaultId" form:"vaultId"`
	Certs   *cert.Certs   `json:"certs" form:"certs"`
}

type KeeperRequest struct {
	KeeperConnection
	Body any `json:"body" form:"body"`
	KeeperOptions
}

type KeeperResponse = keeper.Response

type KeeperParallelRequest struct {
	Connections []KeeperConnection `json:"connections" form:"connections"`
	Body        any                `json:"body" form:"body"`
	KeeperOptions
}

type KeeperParallelResponse struct {
	Connection KeeperConnection  `json:"connection"`
	Response   []keeper.Response `json:"response"`
	Error      error             `json:"error"`
}

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
