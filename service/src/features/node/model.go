package node

import (
	"ivory/src/features/cert"
	"ivory/src/plugins/cloud"
	"ivory/src/plugins/keeper"

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

type CloudVaultConnection struct {
	Host    string     `json:"host" form:"host"`
	Port    int        `json:"port" form:"port"`
	VaultId *uuid.UUID `json:"vaultId" form:"vaultId"`
}

type CloudCredConnection struct {
	Host     string `json:"host" form:"host"`
	Port     int    `json:"port" form:"port"`
	Username string `json:"username" form:"username"`
	Password string `json:"password" form:"password"`
}

type MetricsResponse = cloud.Metrics

type ContainerRequest struct {
	Connection CloudVaultConnection `json:"connection" form:"connection" binding:"required"`
	Image      string               `json:"image" form:"image"`
	Container  string               `json:"container" form:"container"`
	Options    string               `json:"options" form:"options"`
}

type ContainerLogsRequest struct {
	Connection CloudVaultConnection `json:"connection" form:"connection" binding:"required"`
	Container  string               `json:"container" form:"container" binding:"required"`
	Tail       int                  `json:"tail" form:"tail"`
}

type ContainerResponse = cloud.Container

// SPECIFIC (SERVER)
