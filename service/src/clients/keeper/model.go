package keeper

import (
	"crypto/tls"
	"ivory/src/clients/database"
)

// COMMON (WEB AND SERVER)

type KeeperStatus string

const (
	Paused KeeperStatus = "PAUSED"
	Active              = "ACTIVE"
)

type Keeper struct {
	Host   string        `json:"host" form:"host"`
	Port   int           `json:"port" form:"port"`
	Name   *string       `json:"name" form:"name"`
	Status *KeeperStatus `json:"status,omitempty" form:"status"`
}

type Role string

const (
	Leader  Role = "leader"
	Replica      = "replica"
	Unknown      = "unknown"
)

type Node struct {
	State               string                   `json:"state"`
	Role                Role                     `json:"role"`
	Lag                 int64                    `json:"lag"`
	PendingRestart      bool                     `json:"pendingRestart"`
	Database            database.Database        `json:"database"`
	Keeper              Keeper                   `json:"keeper"`
	ScheduledSwitchover *NodeScheduledSwitchover `json:"scheduledSwitchover,omitempty"`
	ScheduledRestart    *NodeScheduledRestart    `json:"scheduledRestart,omitempty"`
	Tags                *map[string]any          `json:"tags,omitempty"`
}

type NodeScheduledSwitchover struct {
	At string `json:"at"`
	To string `json:"to"`
}

type NodeScheduledRestart struct {
	At             string `json:"at"`
	PendingRestart bool   `json:"pendingRestart"`
}

// SPECIFIC (SERVER)

type Request struct {
	Keeper      Keeper       `json:"keeper" form:"keeper"`
	Credentials *Credentials `json:"credentials" form:"credentials"`
	TlsConfig   *tls.Config  `json:"tlsConfig" form:"tlsConfig"`
	Body        any          `json:"body" form:"body"`
}

type Credentials struct {
	Username string `json:"username" form:"username"`
	Password string `json:"password" form:"password"`
}
