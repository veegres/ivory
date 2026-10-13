package cert

import "github.com/google/uuid"

// COMMON (WEB AND SERVER)

type FileUsageType int8

const (
	UPLOAD FileUsageType = iota
	PATH
)

type Cert struct {
	FileName      string        `json:"fileName"`
	FileUsageType FileUsageType `json:"fileUsageType"`
	Path          string        `json:"path"`
	Type          CertType      `json:"type"`
}

type CertMap map[string]Cert

type CertType int8

const (
	CLIENT_CA CertType = iota
	CLIENT_CERT
	CLIENT_KEY
)

func (t CertType) Valid() bool {
	return t >= CLIENT_CA && t <= CLIENT_KEY
}

type RouterAddRequest struct {
	Path string   `json:"path" binding:"required,notblank"`
	Type CertType `json:"type" binding:"valid"`
}

// RouterListRequest is the tool's twin of the http route's ?type= query.
type RouterListRequest struct {
	Type *CertType `json:"type,omitempty" binding:"omitempty,valid" jsonschema:"only certificates of this type: 0 client CA, 1 client certificate, 2 client key"`
}

type Certs struct {
	ClientCAId   *uuid.UUID `json:"clientCAId"`
	ClientKeyId  *uuid.UUID `json:"clientKeyId"`
	ClientCertId *uuid.UUID `json:"clientCertId"`
}

// SPECIFIC (SERVER)
