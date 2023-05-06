package model

import (
	"github.com/google/uuid"
)

// COMMON (WEB AND SERVER)

type BloatTarget struct {
	DbName        string `json:"dbName"`
	Schema        string `json:"schema"`
	Table         string `json:"table"`
	ExcludeSchema string `json:"excludeSchema"`
	ExcludeTable  string `json:"excludeTable"`
}

type BloatRequest struct {
	Cluster    string       `json:"cluster"`
	Connection DbConnection `json:"connection"`
	Target     *BloatTarget `json:"target"`
	Ratio      int          `json:"ratio"`
}

type Bloat struct {
	Uuid         uuid.UUID `json:"uuid"`
	CredentialId uuid.UUID `json:"credentialId"`
	Cluster      string    `json:"cluster"`
	Status       JobStatus `json:"status"`
	Command      string    `json:"command"`
	CommandArgs  []string  `json:"commandArgs"`
	LogsPath     string    `json:"logsPath"`
	CreatedAt    int64     `json:"createdAt"`
}

// SPECIFIC (SERVER)
