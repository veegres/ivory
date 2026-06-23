package tools

import (
	"ivory/clients/console/shell"
	"ivory/clients/storage"
	"ivory/core"
	"ivory/core/utils"
	"ivory/tools/pg_compacttable"
)

type Router struct {
	PgCompactTable *pg_compacttable.Router
}

type Context struct {
	Registry *utils.Registry[Tool, Adapter]
	Router   *Router
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
		Router: &Router{
			PgCompactTable: pg_compacttable.NewRouter(pgCompactTableService),
		},
	}
}
