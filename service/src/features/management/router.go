package management

import (
	. "ivory/src/features/secret"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Router struct {
	service *Service
}

func NewRouter(service *Service) *Router {
	return &Router{service: service}
}

func (r *Router) Erase(context *gin.Context) {
	err := r.service.Erase()
	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"response": "all data was erased"})
}

func (r *Router) ChangeSecret(context *gin.Context) {
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

func (r *Router) GetAppInfo(context *gin.Context) {
	authHeader := context.Request.Header.Get("Authorization")
	info := r.service.GetAppInfo(authHeader)
	context.JSON(http.StatusOK, gin.H{"response": info})
}
