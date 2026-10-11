package management

import (
	"fmt"
	"ivory/core/config"
	"net/http"

	"github.com/gin-gonic/gin"
)

type HttpRouter struct {
	service *Service
}

func NewHttpRouter(service *Service) *HttpRouter {
	return &HttpRouter{service: service}
}

func (r *HttpRouter) Erase(context *gin.Context) {
	err := r.service.Erase()
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"response": "all data was erased"})
}

func (r *HttpRouter) Free(context *gin.Context) {
	err := r.service.Free()
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"response": "data is cleaned"})
}

func (r *HttpRouter) Import(context *gin.Context) {
	file, errFile := context.FormFile("file")
	if errFile != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": errFile.Error()})
		return
	}
	errImport := r.service.BackupImport(file)
	if errImport != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": errImport.Error()})
		return
	}
	context.JSON(http.StatusOK, gin.H{"response": "data is imported"})
}

func (r *HttpRouter) Export(context *gin.Context) {
	fileName := r.service.BackupFileName()
	context.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileName))
	context.Header("Content-Type", "application/json")
	data, err := r.service.BackupExport()
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	_, _ = context.Writer.Write(data)
}

func (r *HttpRouter) ChangeSecret(context *gin.Context) {
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

func (r *HttpRouter) GetAppInfo(context *gin.Context) {
	username := context.GetString(config.AuthContextKey.Username)
	authEnabled := context.GetBool(config.AuthContextKey.Enabled)
	authorised := context.GetBool(config.AuthContextKey.Authorised)
	authError := ""
	if err := context.GetError(config.AuthContextKey.Error); err != nil {
		authError = err.Error()
	}
	info := r.service.GetAppInfo(authorised, authEnabled, username, authError)
	context.JSON(http.StatusOK, gin.H{"response": info})
}
