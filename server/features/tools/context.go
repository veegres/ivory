package tools

import (
	"ivory/clients/shell"
	"ivory/features/job"
	"ivory/features/tools/pg_compacttable"
	"ivory/features/vault"
	"ivory/storage/db"
)

type Context struct {
	Service *Service
	Router  *Router
}

func NewContext(
	shellClient *shell.Client,
	vaultService *vault.Service,
	jobManager *job.Manager,
) *Context {
	st := db.NewStorage("tools.db")
	pgCompactTableBucket := db.NewBucket[pg_compacttable.Response](st, "PgCompactTable")
	pgCompactTableRepo := pg_compacttable.NewRepository(pgCompactTableBucket)
	pgCompactTableService := pg_compacttable.NewService(pgCompactTableRepo, shellClient, vaultService, jobManager)

	service := NewService(pgCompactTableService)
	router := NewRouter(service)

	return &Context{
		Service: service,
		Router:  router,
	}
}
