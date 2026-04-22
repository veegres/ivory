package node

import (
	"ivory/src/clients/keeper"
	"ivory/src/clients/ssh"
	"ivory/src/features/cert"

	"github.com/google/uuid"
)

// COMMON (WEB AND SERVER)

type Request struct {
	Keeper       keeper.Keeper `json:"keeper" form:"keeper"`
	CredentialId *uuid.UUID    `json:"credentialId" form:"credentialId"`
	Certs        *cert.Certs   `json:"certs" form:"certs"`
	Body         any           `json:"body" form:"body"`
}

type Node = keeper.Node

type SshRequest struct {
	VmId uuid.UUID `json:"vmId" form:"vmId" binding:"required"`
}

type SshMetrics = ssh.Metrics

type DockerRequest struct {
	VmId      uuid.UUID `json:"vmId" form:"vmId" binding:"required"`
	Image     string    `json:"image" form:"image"`
	Container string    `json:"container" form:"container"`
	Options   string    `json:"options" form:"options"`
}

type DockerLogsRequest struct {
	VmId      uuid.UUID `json:"vmId" form:"vmId" binding:"required"`
	Container string    `json:"container" form:"container" binding:"required"`
	Tail      int       `json:"tail" form:"tail"`
}

type DockerResult struct {
	Stdout   string `json:"stdout"`
	Stderr   string `json:"stderr"`
	ExitCode int    `json:"exitCode"`
}

// SPECIFIC (SERVER)

type NodeAutoRequest struct {
	Keepers      []keeper.Keeper `json:"keepers" form:"keepers"`
	CredentialId *uuid.UUID      `json:"credentialId" form:"credentialId"`
	Certs        *cert.Certs     `json:"certs" form:"certs"`
	Body         any             `json:"body" form:"body"`
}
