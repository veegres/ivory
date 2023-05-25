package router

import (
	"github.com/gin-gonic/gin"
	"ivory/src/config"
	. "ivory/src/model"
	"ivory/src/service"
	"net/http"
)

type AuthRouter struct {
	env         *config.Env
	authService *service.AuthService
}

func NewAuthRouter(env *config.Env, authService *service.AuthService) *AuthRouter {
	return &AuthRouter{
		env:         env,
		authService: authService,
	}
}

func (r *AuthRouter) Middleware() gin.HandlerFunc {
	switch r.env.Auth {
	case "none":
		return func(context *gin.Context) {
			context.Next()
		}
	case "basic":
		return func(context *gin.Context) {
			authHeader := context.Request.Header.Get("Authorization")
			token, errHeader := r.authService.GetTokenFromHeader(authHeader)
			if errHeader != nil {
				context.Header("WWW-Authenticate", "JWT realm="+r.authService.GetIssuer())
				context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": errHeader.Error()})
				return
			}
			errValidate := r.authService.ValidateToken(token, *r.env.Username)
			if errValidate != nil {
				context.Header("WWW-Authenticate", "JWT realm="+r.authService.GetIssuer())
				context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": errValidate.Error()})
				return
			}
			context.Next()
		}
	default:
		return func(context *gin.Context) {
			context.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "authentication type is not specified or incorrect"})
			return
		}
	}

}

func (r *AuthRouter) Login(context *gin.Context) {
	var login Login
	parseErr := context.ShouldBindJSON(&login)
	if parseErr != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": parseErr.Error()})
		return
	}

	if login.Username != *r.env.Username || login.Password != *r.env.Password {
		context.JSON(http.StatusUnauthorized, gin.H{"error": "credentials are not correct"})
		return
	}

	token, exp, err := r.authService.GenerateToken(login.Username)
	if err != nil {
		context.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"response": gin.H{
		"token":  token,
		"expire": exp.String(),
	}})
}
