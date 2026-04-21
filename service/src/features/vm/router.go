package vm

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Router struct {
	service *Service
}

func NewRouter(service *Service) *Router {
	return &Router{service: service}
}

func (r *Router) GetVmList(context *gin.Context) {
	list, err := r.service.List()
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, gin.H{"response": list})
}

func (r *Router) GetVm(context *gin.Context) {
	id, err := uuid.Parse(context.Param("uuid"))
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	model, err := r.service.Get(id)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, gin.H{"response": model})
}

func (r *Router) PostVm(context *gin.Context) {
	var model VM
	if err := context.ShouldBindJSON(&model); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	id, created, err := r.service.Create(model)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, gin.H{"response": gin.H{"key": id.String(), "vm": created}})
}

func (r *Router) PatchVm(context *gin.Context) {
	id, err := uuid.Parse(context.Param("uuid"))
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var model VM
	if err := context.ShouldBindJSON(&model); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, updated, err := r.service.Update(id, model)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, gin.H{"response": updated})
}

func (r *Router) DeleteVm(context *gin.Context) {
	id, err := uuid.Parse(context.Param("uuid"))
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err = r.service.Delete(id)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, gin.H{"response": "deleted"})
}
