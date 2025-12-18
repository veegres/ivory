package management

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

func (r *Router) Erase(context *gin.Context) {
	err := r.service.Erase()
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"response": "all data was erased"})
}

func (r *Router) Free(context *gin.Context) {
	err := r.service.Free()
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"response": "data was cleaned"})
}

func (r *Router) ChangeSecret(context *gin.Context) {
	var request SecretUpdateRequest
	errBind := context.ShouldBindJSON(&request)
	if errBind != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": errBind.Error()})
		return
	}
	err := r.service.ChangeSecret(request.PreviousKey, request.NewKey)
	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, gin.H{"response": "the secret was changed"})
}

func (r *Router) GetAppInfo(context *gin.Context) {
	info := r.service.GetAppInfo(context)
	context.JSON(http.StatusOK, gin.H{"response": info})
}
