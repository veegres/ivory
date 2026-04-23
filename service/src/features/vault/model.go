package vault

import "github.com/google/uuid"

// COMMON (WEB AND SERVER)

type Vault struct {
	Username string    `json:"username"`
	Secret   string    `json:"secret"`
	Type     VaultType `json:"type"`
}

type VaultType int

const (
	DATABASE_PASSWORD VaultType = iota
	KEEPER_PASSWORD
	SSH_PASSWORD
	SSH_KEY
)

type VaultMap map[string]Vault

// SPECIFIC (SERVER)

type RepositoryVault struct {
	Id uuid.UUID `json:"id"`
	Vault
}
