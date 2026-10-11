package query

import (
	"ivory/features/config"
	"net/http"

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
	session, errSession := ctx.Cookie("session")
	if errSession != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": errSession.Error()})
		return
	}
	handle(Context{
		Connection: con,
		Session:    session,
	})
}
