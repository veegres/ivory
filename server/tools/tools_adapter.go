package tools

import "ivory/core/config"

type Tool string

const (
	PgCompactTable Tool = "pg_compacttable"
)

type Adapter interface {
	// SupportedFeatures reports, for every tool-related env.Feature this tool
	// knows about, whether it supports it for the given database plugin. A
	// feature absent from the map is not a tool capability at all and is
	// left unrestricted.
	SupportedFeatures(plugin env.Plugin) map[env.Feature]bool
	DeleteAll() error
}
