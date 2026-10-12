package secret

// COMMON (WEB AND SERVER)

type SecretStatus struct {
	Key bool `json:"key"`
	Ref bool `json:"ref"`
}

type RouterSetRequest struct {
	Key string `json:"key" binding:"required"`
}

// SPECIFIC (SERVER)
