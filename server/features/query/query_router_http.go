package query

import (
	coreConfig "ivory/core/config"
	"ivory/features/config"

	"github.com/gin-gonic/gin"
)

type HttpRouter struct {
	service       *Service
	configService *config.Service
}

func NewHttpRouter(
	service *Service,
	configService *config.Service,
) *HttpRouter {
	return &HttpRouter{
		service:       service,
		configService: configService,
	}
}

func (r *HttpRouter) withQueryContext(ctx *gin.Context, con Connection, handle func(Context)) {
	handle(Context{
		Connection: con,
		Session:    ctx.GetString(coreConfig.AuthContextKey.Session),
	})
}
