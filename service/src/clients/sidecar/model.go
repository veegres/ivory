package sidecar

import (
	"crypto/tls"
	"ivory/src/clients/database"
)

// COMMON (WEB AND SERVER)

type SidecarStatus string

const (
	Paused SidecarStatus = "PAUSED"
	Active               = "ACTIVE"
)

type Sidecar struct {
	Host   string         `json:"host" form:"host"`
	Port   int            `json:"port" form:"port"`
	Name   *string        `json:"name" form:"name"`
	Status *SidecarStatus `json:"status,omitempty" form:"status"`
}

type Role string

const (
	Leader  Role = "leader"
	Replica      = "replica"
	Unknown      = "unknown"
)

type Instance struct {
	State               string                       `json:"state"`
	Role                Role                         `json:"role"`
	Lag                 int64                        `json:"lag"`
	PendingRestart      bool                         `json:"pendingRestart"`
	Database            database.Database            `json:"database"`
	Sidecar             Sidecar                      `json:"sidecar"`
	ScheduledSwitchover *InstanceScheduledSwitchover `json:"scheduledSwitchover,omitempty"`
	ScheduledRestart    *InstanceScheduledRestart    `json:"scheduledRestart,omitempty"`
	Tags                *map[string]any              `json:"tags,omitempty"`
}

type InstanceScheduledSwitchover struct {
	At string `json:"at"`
	To string `json:"to"`
}

type InstanceScheduledRestart struct {
	At             string `json:"at"`
	PendingRestart bool   `json:"pendingRestart"`
}

// SPECIFIC (SERVER)

type Request struct {
	Sidecar     Sidecar      `json:"sidecar" form:"sidecar"`
	Credentials *Credentials `json:"credentials" form:"credentials"`
	TlsConfig   *tls.Config  `json:"tlsConfig" form:"tlsConfig"`
	Body        any          `json:"body" form:"body"`
}

type Credentials struct {
	Username string `json:"username" form:"username"`
	Password string `json:"password" form:"password"`
}
