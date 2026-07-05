package vault

import "github.com/google/uuid"

// COMMON (WEB AND SERVER)

type Vault struct {
	Type     VaultType `json:"type"`
	Username string    `json:"username"`
	Secret   string    `json:"secret"`
	Metadata *string   `json:"metadata,omitempty"`
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
