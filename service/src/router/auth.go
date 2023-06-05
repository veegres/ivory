package router

import (
	"github.com/gin-gonic/gin"
	. "ivory/src/model"
	"ivory/src/service"
	"net/http"
)

type AuthRouter struct {
	authService    *service.AuthService
	generalService *service.GeneralService
}

func NewAuthRouter(
	authService *service.AuthService,
	generalService *service.GeneralService,
) *AuthRouter {
	return &AuthRouter{
		generalService: generalService,
		authService:    authService,
	}
}

func (r *AuthRouter) Middleware() gin.HandlerFunc {
	appConfig, errConfig := r.generalService.GetAppConfig()
	if errConfig != nil {
		return func(context *gin.Context) {
			context.Next()
		}
	}

	switch appConfig.Auth.Type {
	case NONE:
		return func(context *gin.Context) {
			context.Next()
		}
	case BASIC:
		return func(context *gin.Context) {
			authHeader := context.Request.Header.Get("Authorization")
			token, errHeader := r.authService.GetTokenFromHeader(authHeader)
			if errHeader != nil {
				context.Header("WWW-Authenticate", "JWT realm="+r.authService.GetIssuer())
				context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": errHeader.Error()})
				return
			}
			errValidate := r.authService.ValidateToken(token, appConfig.Auth.Body["username"])
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

	appConfig, errConfig := r.generalService.GetAppConfig()
	if errConfig != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": errConfig.Error()})
	}

	// TODO consider moving this to AuthService
	if login.Username != appConfig.Auth.Body["username"] || login.Password != appConfig.Auth.Body["password"] {
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