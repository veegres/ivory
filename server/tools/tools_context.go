package tools

import (
	"ivory/clients/console/shell"
	"ivory/clients/storage"
	"ivory/core"
	"ivory/core/utils"
	"ivory/tools/pg_compacttable"
)

type HttpRouter struct {
	PgCompactTable *pg_compacttable.HttpRouter
}

type Context struct {
	Registry   *utils.Registry[Tool, Adapter]
	HttpRouter *HttpRouter
}

func NewContext(
	shellClient *shell.Client,
	coreService *core.Service,
) *Context {
	// DB
	st := storage.NewDbStorage("tools.db")

	// PG COMPACT TABLE
	pgCompactTableBucket := storage.NewDbBucket[pg_compacttable.Response](st, "PgCompactTable")
	pgCompactTableRepo := pg_compacttable.NewRepository(pgCompactTableBucket)
	pgCompactTableService := pg_compacttable.NewService(pgCompactTableRepo, shellClient, coreService.Vault, coreService.Job)

	toolRegistry := utils.NewRegistry[Tool, Adapter]()
	toolRegistry.Register(PgCompactTable, pgCompactTableService)

	return &Context{
		Registry: toolRegistry,
		HttpRouter: &HttpRouter{
			PgCompactTable: pg_compacttable.NewHttpRouter(pgCompactTableService),
		},
	}
}
