package instance

import (
	"ivory/src/clients/sidecar"
	sshclient "ivory/src/clients/ssh"
	"ivory/src/features/cert"

	"github.com/google/uuid"
)

// COMMON (WEB AND SERVER)

type InstanceRequest struct {
	Sidecar      sidecar.Sidecar `json:"sidecar" form:"sidecar"`
	CredentialId *uuid.UUID      `json:"credentialId" form:"credentialId"`
	Certs        *cert.Certs     `json:"certs" form:"certs"`
	Body         any             `json:"body" form:"body"`
}

type VmRequest struct {
	VmId uuid.UUID `json:"vmId" form:"vmId" binding:"required"`
}

type VmCommandRequest struct {
	VmId    uuid.UUID `json:"vmId" binding:"required"`
	Command string    `json:"command" binding:"required"`
}

type VmCommandResult = sshclient.CommandResult
type VmMetrics = sshclient.Metrics

// SPECIFIC (SERVER)

type InstanceAutoRequest struct {
	Sidecars     []sidecar.Sidecar `json:"sidecars" form:"sidecars"`
	CredentialId *uuid.UUID        `json:"credentialId" form:"credentialId"`
	Certs        *cert.Certs       `json:"certs" form:"certs"`
	Body         any               `json:"body" form:"body"`
}
