package node

import (
	"ivory/core/service/cert"
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

type KeeperOneRequest struct {
	KeeperConnection
	Body any `json:"body" form:"body"`
	KeeperOptions
}

type KeeperOneResponse = keeper.Response

type KeeperMultiRequest struct {
	Connections []KeeperConnection `json:"connections" form:"connections"`
	Body        any                `json:"body" form:"body"`
	KeeperOptions
}

type KeeperMultiResponse struct {
	Connection KeeperConnection  `json:"connection"`
	Response   []keeper.Response `json:"response"`
	Error      error             `json:"error"`
}

type PlatformVaultConnection struct {
	Host    string     `json:"host" form:"host" binding:"required"`
	Port    int        `json:"port" form:"port" binding:"required"`
	VaultId *uuid.UUID `json:"vaultId" form:"vaultId" binding:"required"`
}

type PlatformCredConnection struct {
	Host     string `json:"host" form:"host"`
	Port     int    `json:"port" form:"port"`
	Username string `json:"username" form:"username"`
	Password string `json:"password" form:"password"`
}

type PlatformMetricsRequest = PlatformVaultConnection

type PlatformMetricsResponse = platform.Metrics

type PlatformCopyIdRequest struct {
	PlatformCredConnection
	PublicKey string
}

type PlatformUpRequest struct {
	Connection PlatformVaultConnection `json:"connection" form:"connection" binding:"required"`
	Image      string                  `json:"image" form:"image" binding:"required"`
	Name       string                  `json:"name" form:"name" binding:"required"`
	Options    string                  `json:"options" form:"options"`
}

type PlatformLogsRequest struct {
	Connection PlatformVaultConnection `json:"connection" form:"connection" binding:"required"`
	Name       string                  `json:"name" form:"name" binding:"required"`
	Tail       int                     `json:"tail" form:"tail"`
	Follow     bool                    `json:"follow" form:"follow"`
}

type PlatformActionRequest struct {
	Connection PlatformVaultConnection `json:"connection" form:"connection" binding:"required"`
	Name       string                  `json:"name" form:"name" binding:"required"`
}

type PlatformResponse = []string

// SPECIFIC (SERVER)
