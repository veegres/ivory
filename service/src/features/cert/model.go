package cert

import "github.com/google/uuid"

// COMMON (WEB AND SERVER)

type Certs struct {
	ClientCAId   *uuid.UUID `json:"clientCAId"`
	ClientKeyId  *uuid.UUID `json:"clientKeyId"`
	ClientCertId *uuid.UUID `json:"clientCertId"`
}

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

type CertAddRequest struct {
	Path string   `json:"path"`
	Type CertType `json:"type"`
}

// SPECIFIC (SERVER)
