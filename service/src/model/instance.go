package model

import "github.com/google/uuid"

// COMMON (WEB AND SERVER)

type InstanceRequest struct {
	Sidecar      Sidecar    `json:"sidecar" form:"sidecar"`
	CredentialId *uuid.UUID `json:"credentialId" form:"credentialId"`
	Certs        Certs      `json:"certs" form:"certs"`
	Body         any        `json:"body" form:"body"`
}

type Instance struct {
	State               string                       `json:"state"`
	Role                string                       `json:"role"`
	Lag                 int64                        `json:"lag"`
	PendingRestart      bool                         `json:"pendingRestart"`
	Database            Database                     `json:"database"`
	Sidecar             Sidecar                      `json:"sidecar"`
	ScheduledSwitchover *InstanceScheduledSwitchover `json:"scheduledSwitchover,omitempty"`
	ScheduledRestart    *InstanceScheduledRestart    `json:"scheduledRestart,omitempty"`
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

type InstanceGateway interface {
	Overview(instance InstanceRequest) ([]Instance, int, error)
	Config(instance InstanceRequest) (any, int, error)
	ConfigUpdate(instance InstanceRequest) (any, int, error)
	Switchover(instance InstanceRequest) (*string, int, error)
	DeleteSwitchover(instance InstanceRequest) (*string, int, error)
	Reinitialize(instance InstanceRequest) (*string, int, error)
	Restart(instance InstanceRequest) (*string, int, error)
	DeleteRestart(instance InstanceRequest) (*string, int, error)
	Reload(instance InstanceRequest) (*string, int, error)
	Failover(instance InstanceRequest) (*string, int, error)
	Activate(instance InstanceRequest) (*string, int, error)
	Pause(instance InstanceRequest) (*string, int, error)
}
