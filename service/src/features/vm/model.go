package vm

import "github.com/google/uuid"

// COMMON (WEB AND SERVER)

type VM struct {
	Id       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	Host     string    `json:"host"`
	Port     int       `json:"port"`
	Username string    `json:"username"`
	SshKey   string    `json:"sshKey"`
}

// SPECIFIC (SERVER)
