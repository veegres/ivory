package tools

import "ivory/features/tools/pg_compacttable"

type Router struct {
	Bloat *pg_compacttable.Router
}

func NewRouter(service *Service) *Router {
	return &Router{
		Bloat: pg_compacttable.NewRouter(service.pgCompactTable),
	}
}
