package pg_compacttable

import (
	"ivory/core/service/job"
	"ivory/plugins/database"

	"github.com/google/uuid"
)

// COMMON (WEB AND SERVER)

type Response struct {
	Uuid        uuid.UUID  `json:"uuid"`
	JobId       job.JobID  `json:"jobId"`
	Cluster     string     `json:"cluster"`
	VaultId     *uuid.UUID `json:"vaultId"`
	Status      job.Status `json:"status"`
	Command     string     `json:"command"`
	CommandArgs []string   `json:"commandArgs"`
	CreatedAt   int64      `json:"createdAt"`
}

type RunRequest struct {
	Cluster string          `json:"cluster"`
	Db      database.Config `json:"db"`
	VaultId *uuid.UUID      `json:"vaultId"`
	Target  *Target         `json:"target"`
	Options Options         `json:"options"`
}

type Target struct {
	Database      string `json:"database"`
	Schema        string `json:"schema"`
	Table         string `json:"table"`
	ExcludeSchema string `json:"excludeSchema"`
	ExcludeTable  string `json:"excludeTable"`
}

type Options struct {
	Force           bool `json:"force"`
	NoReindex       bool `json:"noReindex"`
	NoInitialVacuum bool `json:"noInitialVacuum"`
	InitialReindex  bool `json:"initialReindex"`
	RoutineVacuum   bool `json:"routineVacuum"`
	DelayRatio      int  `json:"delayRatio"`
	MinTableSize    int  `json:"minTableSize"`
	MaxTableSize    int  `json:"maxTableSize"`
}
