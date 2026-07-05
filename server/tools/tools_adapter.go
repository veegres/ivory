package tools

import "ivory/core/config"

type Tool string

const (
	PgCompactTable Tool = "pg_compacttable"
)

type Adapter interface {
	SupportedFeatures(plugin env.Plugin) []env.Feature
	DeleteAll() error
}
