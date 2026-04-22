package vm

import "github.com/google/uuid"

// COMMON (WEB AND SERVER)

type VM struct {
	Id       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	Host     string    `json:"host"`
	Username string    `json:"username"`
	SshPort  int       `json:"sshPort"`
	SshKey   string    `json:"sshKey"`
}

// SPECIFIC (SERVER)
