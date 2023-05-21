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
	authService   *service.AuthService
}

func NewInfoRouter(
	env *config.Env,
	secretService *service.SecretService,
	authService *service.AuthService,
) *InfoRouter {
	return &InfoRouter{
		env:           env,
		secretService: secretService,
		authService:   authService,
	}
}

func (r *InfoRouter) Info(context *gin.Context) {
	authHeader := context.Request.Header.Get("Authorization")
	authorized, authError := r.authInfo(authHeader)

	context.JSON(http.StatusOK, gin.H{"response": gin.H{
		"company": r.env.Company,
		"auth": gin.H{
			"type":       r.env.Auth,
			"authorized": authorized,
			"error":      authError,
		},
		"version": gin.H{
			"tag":    r.env.Tag,
			"commit": r.env.Commit,
		},
		"secret": r.secretService.Status(),
	}})
}

func (r *InfoRouter) authInfo(header string) (bool, string) {
	token, errHeader := r.authService.GetTokenFromHeader(header)
	if errHeader != nil {
		return false, errHeader.Error()
	} else {
		errValidate := r.authService.ValidateToken(token, *r.env.Username)
		if errValidate != nil {
			return false, errValidate.Error()
		}
	}
	return true, ""
}
