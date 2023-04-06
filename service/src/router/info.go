package router

import (
	"github.com/gin-gonic/gin"
	"ivory/src/config"
	"ivory/src/service"
	"net/http"
)

type InfoRouter struct {
	env           *config.Env
	secretService *service.SecretService
}

func NewInfoRouter(env *config.Env, secretService *service.SecretService) *InfoRouter {
	return &InfoRouter{env: env, secretService: secretService}
}

func (r *InfoRouter) Info(context *gin.Context) {
	context.JSON(http.StatusOK, gin.H{"response": gin.H{
		"company": r.env.Company,
		"auth":    r.env.Auth,
		"version": gin.H{
			"tag":    r.env.Tag,
			"commit": r.env.Commit,
		},
		"secret": r.secretService.Status(),
	}})
}
