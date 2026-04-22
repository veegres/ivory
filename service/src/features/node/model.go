package node

import (
	"ivory/src/clients/keeper"
	"ivory/src/clients/ssh"
	"ivory/src/features/cert"

	"github.com/google/uuid"
)

// COMMON (WEB AND SERVER)

type NodeConnection struct {
	VmId       *uuid.UUID `json:"vmId,omitempty" form:"vmId"`
	Host       string     `json:"host" form:"host"`
	SshPort    int        `json:"sshPort,omitempty" form:"sshPort"`
	KeeperPort int        `json:"keeperPort" form:"keeperPort"`
	DbPort     int        `json:"dbPort,omitempty" form:"dbPort"`
}

type Node struct {
	Connection NodeConnection  `json:"connection"`
	Response   keeper.Response `json:"response"`
}

type Request struct {
	Connection   NodeConnection `json:"connection" form:"connection"`
	CredentialId *uuid.UUID     `json:"credentialId" form:"credentialId"`
	Certs        *cert.Certs    `json:"certs" form:"certs"`
	Body         any            `json:"body" form:"body"`
}

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
	Connections  []NodeConnection `json:"connections" form:"connections"`
	CredentialId *uuid.UUID       `json:"credentialId" form:"credentialId"`
	Certs        *cert.Certs      `json:"certs" form:"certs"`
	Body         any              `json:"body" form:"body"`
}
