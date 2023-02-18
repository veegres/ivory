package router

import (
	"github.com/gin-gonic/gin"
	. "ivory/model"
	"net/http"
)

type InstanceRouter struct {
	instanceService InstanceService
}

func NewInstanceRouter(instanceService InstanceService) *InstanceRouter {
	return &InstanceRouter{instanceService: instanceService}
}

func (r *InstanceRouter) GetInstanceInfo(context *gin.Context) {
	var instance InstanceRequest
	err := context.ShouldBindQuery(&instance)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	body, status, err := r.instanceService.Info(instance)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	context.JSON(status, gin.H{"response": body})
}

func (r *InstanceRouter) GetInstanceOverview(context *gin.Context) {
	var instance InstanceRequest
	err := context.ShouldBindQuery(&instance)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	body, status, err := r.instanceService.Overview(instance)
	r.handleResponse(context, body, status, err)
}

func (r *InstanceRouter) GetInstanceConfig(context *gin.Context) {
	var instance InstanceRequest
	err := context.ShouldBindQuery(&instance)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	body, status, err := r.instanceService.Config(instance)
	r.handleResponse(context, body, status, err)
}

func (r *InstanceRouter) PatchInstanceConfig(context *gin.Context) {
	var instance InstanceRequest
	err := context.ShouldBindJSON(&instance)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	body, status, err := r.instanceService.ConfigUpdate(instance)
	r.handleResponse(context, body, status, err)
}

func (r *InstanceRouter) PostInstanceSwitchover(context *gin.Context) {
	var instance InstanceRequest
	err := context.ShouldBindJSON(&instance)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	body, status, err := r.instanceService.Switchover(instance)
	r.handleResponse(context, body, status, err)
}

func (r *InstanceRouter) PostInstanceReinitialize(context *gin.Context) {
	var instance InstanceRequest
	err := context.ShouldBindJSON(&instance)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	body, status, err := r.instanceService.Reinitialize(instance)
	r.handleResponse(context, body, status, err)
}

func (r *InstanceRouter) handleResponse(context *gin.Context, body any, status int, err error) {
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if status >= 200 && status < 300 {
		context.JSON(status, gin.H{"response": body})
	} else {
		context.JSON(status, gin.H{"error": body})
	}
}
