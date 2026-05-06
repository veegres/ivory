package bloat

import (
	. "ivory/src/features/tools/bloat/job"
	"ivory/src/plugins/database"
	"os"

	"github.com/google/uuid"
)

// COMMON (WEB AND SERVER)

type Bloat struct {
	Uuid        uuid.UUID  `json:"uuid"`
	Cluster     string     `json:"cluster"`
	VaultId     *uuid.UUID `json:"vaultId"`
	Status      JobStatus  `json:"status"`
	Command     string     `json:"command"`
	CommandArgs []string   `json:"commandArgs"`
	LogsPath    string     `json:"logsPath"`
	CreatedAt   int64      `json:"createdAt"`
}

type BloatRequest struct {
	Cluster string          `json:"cluster"`
	Db      database.Config `json:"db"`
	VaultId *uuid.UUID      `json:"vaultId"`
	Target  *BloatTarget    `json:"target"`
	Options BloatOptions    `json:"options"`
}

type BloatTarget struct {
	Database      string `json:"database"`
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

// SPECIFIC (SERVER)

type element struct {
	job    *Job
	model  *Bloat
	writer *os.File
}
