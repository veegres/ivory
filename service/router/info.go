package router

import (
	"github.com/gin-gonic/gin"
	"ivory/config"
	"ivory/service"
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
		"secret":  r.secretService.Status(),
	}})
}
