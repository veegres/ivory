package vault

import "github.com/google/uuid"

// COMMON (WEB AND SERVER)

type Vault struct {
	Type     VaultType `json:"type" binding:"valid"`
	Username string    `json:"username" binding:"required,notblank"`
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

func (t VaultType) Valid() bool {
	return t >= DATABASE_PASSWORD && t <= SSH_KEY
}

type VaultMap map[string]Vault

// RouterListRequest is the tool's twin of the http route's ?type= query.
type RouterListRequest struct {
	Type *VaultType `json:"type,omitempty" binding:"omitempty,valid" jsonschema:"only vaults of this type: 0 database password, 1 keeper password, 2 ssh password, 3 ssh key"`
}

// SPECIFIC (SERVER)

type RepositoryVault struct {
	Id uuid.UUID `json:"id"`
	Vault
}
