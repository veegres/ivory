package config

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Router struct {
	service *Service
}

func NewRouter(service *Service) *Router {
	return &Router{service: service}
}

func (r *Router) InitialiseMiddleware() gin.HandlerFunc {
	return func(context *gin.Context) {
		_, err := r.service.GetAppConfig()
		if err != nil {
			context.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		context.Next()
	}
}

func (r *Router) SetAppConfig(context *gin.Context) {
	var appConfig NewAppConfig
	errBind := context.ShouldBindJSON(&appConfig)
	if errBind != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": errBind.Error()})
		return
	}

	errConfig := r.service.SetAppConfig(appConfig)
	if errConfig != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": errConfig.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"response": "config was successfully set up"})
}
