package instance

import (
	"ivory/src/clients/sidecar"
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

// SPECIFIC (SERVER)
