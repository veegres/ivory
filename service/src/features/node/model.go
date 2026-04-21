package node

import (
	"ivory/src/clients/sidecar"
	"ivory/src/clients/ssh"
	"ivory/src/features/cert"

	"github.com/google/uuid"
)

// COMMON (WEB AND SERVER)

type Request struct {
	Sidecar      sidecar.Sidecar `json:"sidecar" form:"sidecar"`
	CredentialId *uuid.UUID      `json:"credentialId" form:"credentialId"`
	Certs        *cert.Certs     `json:"certs" form:"certs"`
	Body         any             `json:"body" form:"body"`
}

type Node = sidecar.Instance

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
	Sidecars     []sidecar.Sidecar `json:"sidecars" form:"sidecars"`
	CredentialId *uuid.UUID        `json:"credentialId" form:"credentialId"`
	Certs        *cert.Certs       `json:"certs" form:"certs"`
	Body         any               `json:"body" form:"body"`
}
