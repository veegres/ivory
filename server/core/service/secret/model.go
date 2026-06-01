package secret

// COMMON (WEB AND SERVER)

type SecretStatus struct {
	Key bool `json:"key"`
	Ref bool `json:"ref"`
}

type SecretSetRequest struct {
	Key string `json:"key"`
}

// SPECIFIC (SERVER)
