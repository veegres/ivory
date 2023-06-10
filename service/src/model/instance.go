package model

import "github.com/google/uuid"

// COMMON (WEB AND SERVER)

type InstanceRequest struct {
	Host         string     `json:"host" form:"host"`
	Port         int        `json:"port" form:"port"`
	CredentialId *uuid.UUID `json:"credentialId" form:"credentialId"`
	Certs        Certs      `json:"certs" form:"certs"`
	Body         any        `json:"body" form:"body"`
}

type InstanceQueryRequest struct {
	Json InstanceRequest `json:"json" form:"json"`
}

type Instance struct {
	State    string   `json:"state"`
	Role     string   `json:"role"`
	Lag      int64    `json:"lag"`
	Database Database `json:"database"`
	Sidecar  Sidecar  `json:"sidecar"`
}

type InstanceInfo struct {
	State   string  `json:"state"`
	Role    string  `json:"role"`
	Sidecar Sidecar `json:"sidecar"`
}

// SPECIFIC (SERVER)

// InstanceGateway TODO add common return types to cast interface to them and create mappers for each impl (patroni)
type InstanceGateway interface {
	Info(instance InstanceRequest) (InstanceInfo, int, error)
	Overview(instance InstanceRequest) ([]Instance, int, error)
	Config(instance InstanceRequest) (any, int, error)
	ConfigUpdate(instance InstanceRequest) (any, int, error)
	Switchover(instance InstanceRequest) (any, int, error)
	Reinitialize(instance InstanceRequest) (any, int, error)
}
