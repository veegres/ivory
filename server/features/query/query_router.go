package query

import (
	"ivory/features/config"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Router struct {
	service       *Service
	configService *config.Service
}

func NewRouter(
	service *Service,
	configService *config.Service,
) *Router {
	return &Router{
		service:       service,
		configService: configService,
	}
}

func (r *Router) withQueryContext(ctx *gin.Context, con Connection, handle func(Context)) {
	session, errSession := ctx.Cookie("session")
	if errSession != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": errSession.Error()})
		return
	}
	handle(Context{
		Connection: con,
		Session:    session,
	})
}
