package router

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/exp/slog"
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
	return func(context *gin.Context) {
		appConfig, errConfig := r.generalService.GetAppConfig()
		if errConfig != nil {
			context.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": errConfig.Error()})
			return
		}

		// NOTE: Browser doesn't support headers in EventSource so there is no option to send Auth header.
		// We consider it as safe because EventSource for now is used only in logs streaming where id is
		// unique. This can work if change header to cookie usage.
		if appConfig.Auth.Type != NONE && context.Request.Header.Get("Accept") == "text/event-stream" {
			slog.Warn("unsafe request", "path", context.Request.URL.Path)
			context.Next()
			return
		}

		authHeader := context.Request.Header.Get("Authorization")
		valid, errValid := r.authService.ValidateHeader(authHeader, appConfig.Auth)
		if !valid {
			context.Header("WWW-Authenticate", "Bearer JWT realm="+r.authService.GetIssuer())
			context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": errValid})
			return
		}

		context.Next()
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
		context.JSON(http.StatusForbidden, gin.H{"error": errConfig.Error()})
		return
	}

	token, exp, err := r.authService.Login(login, appConfig.Auth)
	if err != nil {
		context.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"response": gin.H{
		"token":  token,
		"expire": exp.String(),
	}})
}
