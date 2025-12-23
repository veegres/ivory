package management

import (
	"encoding/json"
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

	context.JSON(http.StatusOK, gin.H{"response": "data is cleaned"})
}

func (r *Router) Import(context *gin.Context) {
	file, errFile := context.FormFile("file")
	if errFile != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": errFile.Error()})
		return
	}
	errImport := r.service.Import(file)
	if errImport != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": errImport.Error()})
		return
	}
	context.JSON(http.StatusOK, gin.H{"response": "data is imported"})
}

func (r *Router) Export(context *gin.Context) {
	context.Header("Content-Disposition", "attachment; filename=ivory.bak")
	context.Header("Content-Type", "application/json")
	data, err := r.service.Export()
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	_ = json.NewEncoder(context.Writer).Encode(data)
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
