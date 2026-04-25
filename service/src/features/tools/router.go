package tools

import "ivory/src/features/tools/bloat"

type Router struct {
	Bloat *bloat.Router
}

func NewRouter(service *Service) *Router {
	return &Router{
		Bloat: bloat.NewRouter(service.bloat),
	}
}
