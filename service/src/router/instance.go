package router

import (
	"encoding/json"
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

func (r *InstanceRouter) GetInstanceOverview(context *gin.Context) {
	handleParamRequest(context, r.instanceService.Overview)
}

func (r *InstanceRouter) GetInstanceConfig(context *gin.Context) {
	handleParamRequest(context, r.instanceService.Config)
}

func (r *InstanceRouter) PatchInstanceConfig(context *gin.Context) {
	handleBodyRequest(context, r.instanceService.ConfigUpdate)
}

func (r *InstanceRouter) PostInstanceSwitchover(context *gin.Context) {
	handleBodyRequest(context, r.instanceService.Switchover)
}

func (r *InstanceRouter) DeleteInstanceSwitchover(context *gin.Context) {
	handleParamRequest(context, r.instanceService.DeleteSwitchover)
}

func (r *InstanceRouter) PostInstanceReinitialize(context *gin.Context) {
	handleBodyRequest(context, r.instanceService.Reinitialize)
}

func (r *InstanceRouter) PostInstanceRestart(context *gin.Context) {
	handleBodyRequest(context, r.instanceService.Restart)
}

func (r *InstanceRouter) DeleteInstanceRestart(context *gin.Context) {
	handleParamRequest(context, r.instanceService.DeleteRestart)
}

func (r *InstanceRouter) PostInstanceReload(context *gin.Context) {
	handleBodyRequest(context, r.instanceService.Reload)
}

func (r *InstanceRouter) PostInstanceFailover(context *gin.Context) {
	handleBodyRequest(context, r.instanceService.Failover)
}

func handleParamRequest[T any](context *gin.Context, action func(instance InstanceRequest) (T, int, error)) {
	query := context.Query("request")
	var instance InstanceRequest
	errBind := json.Unmarshal([]byte(query), &instance)
	if errBind != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": errBind.Error()})
		return
	}
	body, status, err := action(instance)
	handleResponse(context, body, status, err)
}

func handleBodyRequest[T any](context *gin.Context, action func(instance InstanceRequest) (T, int, error)) {
	var instance InstanceRequest
	errBind := context.ShouldBindJSON(&instance)
	if errBind != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": errBind.Error()})
		return
	}
	body, status, err := action(instance)
	handleResponse(context, body, status, err)
}

func handleResponse(context *gin.Context, body any, status int, err error) {
	if err != nil {
		context.JSON(status, gin.H{"error": err.Error()})
		return
	}
	context.JSON(status, gin.H{"response": body})
}
