package keeper

import (
	"crypto/tls"
)

// COMMON (WEB AND SERVER)

type KeeperStatus string

const (
	Paused KeeperStatus = "PAUSED"
	Active              = "ACTIVE"
)

type Role string

const (
	Leader  Role = "leader"
	Replica      = "replica"
	Unknown      = "unknown"
)

type Response struct {
	Name                *string                  `json:"name"`
	Status              *KeeperStatus            `json:"status,omitempty"`
	State               string                   `json:"state"`
	Role                Role                     `json:"role"`
	Lag                 int64                    `json:"lag"`
	PendingRestart      bool                     `json:"pendingRestart"`
	ScheduledSwitchover *NodeScheduledSwitchover `json:"scheduledSwitchover,omitempty"`
	ScheduledRestart    *NodeScheduledRestart    `json:"scheduledRestart,omitempty"`
	Tags                *map[string]any          `json:"tags,omitempty"`

	// Discovered Topology (Crucial for Auto-Creation)
	DiscoveredHost       string `json:"discoveredHost"`
	DiscoveredKeeperPort int    `json:"discoveredKeeperPort"`
	DiscoveredDbPort     int    `json:"discoveredDbPort"`
}

type Keeper struct {
	Host   string        `json:"host"`
	Port   int           `json:"port"`
	Name   *string       `json:"name"`
	Status *KeeperStatus `json:"status,omitempty"`
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
	Host        string       `json:"host" form:"host"`
	Port        int          `json:"port" form:"port"`
	Credentials *Credentials `json:"credentials" form:"credentials"`
	TlsConfig   *tls.Config  `json:"tlsConfig" form:"tlsConfig"`
	Body        any          `json:"body" form:"body"`
}

type Credentials struct {
	Username string `json:"username" form:"username"`
	Password string `json:"password" form:"password"`
}
