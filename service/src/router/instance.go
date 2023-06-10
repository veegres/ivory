package router

import (
	"github.com/gin-gonic/gin"
	. "ivory/src/model"
	"net/http"
)

type InstanceRouter struct {
	instanceService InstanceGateway
}

func NewInstanceRouter(instanceService InstanceGateway) *InstanceRouter {
	return &InstanceRouter{instanceService: instanceService}
}

func (r *InstanceRouter) GetInstanceInfo(context *gin.Context) {
	var instance InstanceQueryRequest
	errBind := context.ShouldBindQuery(&instance)
	if errBind != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": errBind.Error()})
		return
	}

	body, status, err := r.instanceService.Info(instance.Json)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	context.JSON(status, gin.H{"response": body})
}

func (r *InstanceRouter) GetInstanceOverview(context *gin.Context) {
	var instance InstanceQueryRequest
	errBind := context.ShouldBindQuery(&instance)
	if errBind != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": errBind.Error()})
		return
	}

	body, status, err := r.instanceService.Overview(instance.Json)
	r.handleResponse(context, body, status, err)
}

func (r *InstanceRouter) GetInstanceConfig(context *gin.Context) {
	var instance InstanceQueryRequest
	errBind := context.ShouldBindQuery(&instance)
	if errBind != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": errBind.Error()})
		return
	}

	body, status, err := r.instanceService.Config(instance.Json)
	r.handleResponse(context, body, status, err)
}

func (r *InstanceRouter) PatchInstanceConfig(context *gin.Context) {
	var instance InstanceRequest
	errBind := context.ShouldBindJSON(&instance)
	if errBind != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": errBind.Error()})
		return
	}

	body, status, err := r.instanceService.ConfigUpdate(instance)
	r.handleResponse(context, body, status, err)
}

func (r *InstanceRouter) PostInstanceSwitchover(context *gin.Context) {
	var instance InstanceRequest
	errBind := context.ShouldBindJSON(&instance)
	if errBind != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": errBind.Error()})
		return
	}

	body, status, err := r.instanceService.Switchover(instance)
	r.handleResponse(context, body, status, err)
}

func (r *InstanceRouter) PostInstanceReinitialize(context *gin.Context) {
	var instance InstanceRequest
	errBind := context.ShouldBindJSON(&instance)
	if errBind != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": errBind.Error()})
		return
	}

	body, status, err := r.instanceService.Reinitialize(instance)
	r.handleResponse(context, body, status, err)
}

func (r *InstanceRouter) handleResponse(context *gin.Context, body any, status int, err error) {
	if err != nil {
		context.JSON(status, gin.H{"error": err.Error()})
		return
	}
	context.JSON(status, gin.H{"response": body})
}
