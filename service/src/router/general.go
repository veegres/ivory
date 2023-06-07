package router

import (
	"github.com/gin-gonic/gin"
	. "ivory/src/model"
	"ivory/src/service"
	"net/http"
)

type GeneralRouter struct {
	service *service.GeneralService
}

func NewGeneralRouter(service *service.GeneralService) *GeneralRouter {
	return &GeneralRouter{service: service}
}

func (r *GeneralRouter) Erase(context *gin.Context) {
	err := r.service.Erase()
	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"response": "all data was erased"})
}

func (r *GeneralRouter) ChangeSecret(context *gin.Context) {
	var secret SecretUpdateRequest
	errBind := context.ShouldBindJSON(&secret)
	if errBind != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": errBind.Error()})
		return
	}
	err := r.service.ChangeSecret(secret.PreviousKey, secret.NewKey)
	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, gin.H{"response": "the secret was changed"})
}

func (r *GeneralRouter) GetAppInfo(context *gin.Context) {
	authHeader := context.Request.Header.Get("Authorization")
	info := r.service.GetAppInfo(authHeader)
	context.JSON(http.StatusOK, gin.H{"response": info})
}

func (r *GeneralRouter) SetAppConfig(context *gin.Context) {
	var config AppConfig
	errBind := context.ShouldBindJSON(&config)
	if errBind != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": errBind.Error()})
		return
	}

	errConfig := r.service.SetAppConfig(config)
	if errConfig != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": errConfig.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"response": "config was successfully set up"})
}
