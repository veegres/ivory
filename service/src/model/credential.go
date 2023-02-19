package model

type CredentialType int8

const (
	POSTGRES CredentialType = iota
	PATRONI
)

type Credential struct {
	Username string         `json:"username"`
	Password string         `json:"password"`
	Type     CredentialType `json:"type"`
}

type SecretStatus struct {
	Key bool `json:"key"`
	Ref bool `json:"ref"`
}

type SecretSetRequest struct {
	Key string `json:"key"`
	Ref string `json:"ref"`
}

type SecretUpdateRequest struct {
	PreviousKey string `json:"previousKey"`
	NewKey      string `json:"newKey"`
}
