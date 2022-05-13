package model

type Credential struct {
	Username string `json:"username"`
	Password string `json:"password"`
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
