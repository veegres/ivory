package instance

import (
	"ivory/src/features/cert"
	"ivory/src/features/query"

	"github.com/google/uuid"
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

type InstanceRequest struct {
	Sidecar      Sidecar     `json:"sidecar" form:"sidecar"`
	CredentialId *uuid.UUID  `json:"credentialId" form:"credentialId"`
	Certs        *cert.Certs `json:"certs" form:"certs"`
	Body         any         `json:"body" form:"body"`
}

type Instance struct {
	State               string                       `json:"state"`
	Role                string                       `json:"role"`
	Lag                 int64                        `json:"lag"`
	PendingRestart      bool                         `json:"pendingRestart"`
	Database            query.Database               `json:"database"`
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
