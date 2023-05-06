package model

// COMMON (WEB AND SERVER)

type Cert struct {
	FileName      string        `json:"fileName"`
	FileUsageType FileUsageType `json:"fileUsageType"`
	Path          string        `json:"path"`
	Type          CertType      `json:"type"`
}

type CertMap map[string]Cert

type CertType int8

const (
	CLIENT_CA PasswordType = iota
	CLIENT_CERT
	CLIENT_KEY
)

type CertAddRequest struct {
	Path string   `json:"path"`
	Type CertType `json:"type"`
}

// SPECIFIC (SERVER)
