package bloat

import (
	"ivory/src/features/bloat/job"

	"github.com/google/uuid"
)

// COMMON (WEB AND SERVER)

type DbConnection struct {
	Host         string    `json:"host"`
	Port         int       `json:"port"`
	CredentialId uuid.UUID `json:"credentialId"`
}

type BloatTarget struct {
	DbName        string `json:"dbName"`
	Schema        string `json:"schema"`
	Table         string `json:"table"`
	ExcludeSchema string `json:"excludeSchema"`
	ExcludeTable  string `json:"excludeTable"`
}

type BloatOptions struct {
	Force           bool `json:"force"`
	NoReindex       bool `json:"noReindex"`
	NoInitialVacuum bool `json:"noInitialVacuum"`
	InitialReindex  bool `json:"initialReindex"`
	RoutineVacuum   bool `json:"routineVacuum"`
	DelayRatio      int  `json:"delayRatio"`
	MinTableSize    int  `json:"minTableSize"`
	MaxTableSize    int  `json:"maxTableSize"`
}

type BloatRequest struct {
	Cluster    string       `json:"cluster"`
	Connection DbConnection `json:"connection"`
	Target     *BloatTarget `json:"target"`
	Options    BloatOptions `json:"options"`
}

type Bloat struct {
	Uuid         uuid.UUID     `json:"uuid"`
	CredentialId uuid.UUID     `json:"credentialId"`
	Cluster      string        `json:"cluster"`
	Status       job.JobStatus `json:"status"`
	Command      string        `json:"command"`
	CommandArgs  []string      `json:"commandArgs"`
	LogsPath     string        `json:"logsPath"`
	CreatedAt    int64         `json:"createdAt"`
}

// SPECIFIC (SERVER)
