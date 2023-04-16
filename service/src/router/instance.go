package router

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
	var instance InstanceRequestTmp
	err := context.ShouldBindQuery(&instance)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var cred *uuid.UUID
	if instance.CredentialId != nil {
		credTmp, _ := uuid.Parse(*instance.CredentialId)
		cred = &credTmp
	}
	request := InstanceRequest{
		Host:         instance.Host,
		Port:         instance.Port,
		CredentialId: cred,
		Certs:        instance.Certs,
		Body:         instance.Body,
	}
	body, status, err := r.instanceService.Info(request)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	context.JSON(status, gin.H{"response": body})
}

func (r *InstanceRouter) GetInstanceOverview(context *gin.Context) {
	var instance InstanceRequestTmp
	err := context.ShouldBindQuery(&instance)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var cred *uuid.UUID
	if instance.CredentialId != nil {
		credTmp, _ := uuid.Parse(*instance.CredentialId)
		cred = &credTmp
	}
	request := InstanceRequest{
		Host:         instance.Host,
		Port:         instance.Port,
		CredentialId: cred,
		Certs:        instance.Certs,
		Body:         instance.Body,
	}
	body, status, err := r.instanceService.Overview(request)
	r.handleResponse(context, body, status, err)
}

func (r *InstanceRouter) GetInstanceConfig(context *gin.Context) {
	var instance InstanceRequestTmp
	err := context.ShouldBindQuery(&instance)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var cred *uuid.UUID
	if instance.CredentialId != nil {
		credTmp, _ := uuid.Parse(*instance.CredentialId)
		cred = &credTmp
	}
	request := InstanceRequest{
		Host:         instance.Host,
		Port:         instance.Port,
		CredentialId: cred,
		Certs:        instance.Certs,
		Body:         instance.Body,
	}
	body, status, err := r.instanceService.Config(request)
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
		context.JSON(status, gin.H{"error": err.Error()})
		return
	}
	context.JSON(status, gin.H{"response": body})
}
