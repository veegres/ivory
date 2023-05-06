package model

// COMMON (WEB AND SERVER)

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

// SPECIFIC (SERVER)
