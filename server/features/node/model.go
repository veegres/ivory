package node

import (
	"ivory/features/cert"
	"ivory/plugins/keeper"
	"ivory/plugins/platform"

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

type PlatformVaultConnection struct {
	Host    string     `json:"host" form:"host"`
	Port    int        `json:"port" form:"port"`
	VaultId *uuid.UUID `json:"vaultId" form:"vaultId"`
}

type PlatformCredConnection struct {
	Host     string `json:"host" form:"host"`
	Port     int    `json:"port" form:"port"`
	Username string `json:"username" form:"username"`
	Password string `json:"password" form:"password"`
}

type MetricsResponse = platform.Metrics

type PlatformDeployRequest struct {
	Connection PlatformVaultConnection `json:"connection" form:"connection" binding:"required"`
	Image      string                  `json:"image" form:"image"`
	Name       string                  `json:"name" form:"name"`
	Options    string                  `json:"options" form:"options"`
}

type PlatformLogsRequest struct {
	Connection PlatformVaultConnection `json:"connection" form:"connection" binding:"required"`
	Name       string                  `json:"name" form:"name" binding:"required"`
	Tail       int                     `json:"tail" form:"tail"`
}

type PlatformResponse = platform.OperationResult

// SPECIFIC (SERVER)
