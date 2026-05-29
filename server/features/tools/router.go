package tools

import "ivory/features/tools/pg_compacttable"

type Router struct {
	Bloat *bloat.Router
}

func NewRouter(service *Service) *Router {
	return &Router{
		Bloat: bloat.NewRouter(service.bloat),
	}
}
