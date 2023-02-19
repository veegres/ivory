package model

type FileUsageType int8

const (
	UPLOAD CredentialType = iota
	PATH
)

type CertType int8

const (
	CLIENT_CA CredentialType = iota
	CLIENT_CERT
	CLIENT_KEY
)

type CertModel struct {
	FileName      string        `json:"fileName"`
	FileUsageType FileUsageType `json:"fileUsageType"`
	Path          string        `json:"path"`
	Type          CertType      `json:"type"`
}

type CertRequest struct {
	Path string   `json:"path"`
	Type CertType `json:"type"`
}
