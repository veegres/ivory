package router

import (
	"github.com/gin-gonic/gin"
	"ivory/src/service"
	"net/http"
)

type EraseRouter struct {
	service *service.EraseService
}

func NewEraseRouter(service *service.EraseService) *EraseRouter {
	return &EraseRouter{service: service}
}

func (r *EraseRouter) Erase(context *gin.Context) {
	err := r.service.Erase()
	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"response": "all data was erased"})
}
