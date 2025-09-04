package instance

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

type InstanceRouter struct {
	gateway InstanceClient
}

func NewInstanceRouter(gateway InstanceClient) *InstanceRouter {
	return &InstanceRouter{gateway: gateway}
}

func (r *InstanceRouter) GetInstanceOverview(context *gin.Context) {
	handleParamRequest(context, r.gateway.Overview)
}

func (r *InstanceRouter) GetInstanceConfig(context *gin.Context) {
	handleParamRequest(context, r.gateway.Config)
}

func (r *InstanceRouter) PatchInstanceConfig(context *gin.Context) {
	handleBodyRequest(context, r.gateway.ConfigUpdate)
}

func (r *InstanceRouter) PostInstanceSwitchover(context *gin.Context) {
	handleBodyRequest(context, r.gateway.Switchover)
}

func (r *InstanceRouter) DeleteInstanceSwitchover(context *gin.Context) {
	handleParamRequest(context, r.gateway.DeleteSwitchover)
}

func (r *InstanceRouter) PostInstanceReinitialize(context *gin.Context) {
	handleBodyRequest(context, r.gateway.Reinitialize)
}

func (r *InstanceRouter) PostInstanceRestart(context *gin.Context) {
	handleBodyRequest(context, r.gateway.Restart)
}

func (r *InstanceRouter) DeleteInstanceRestart(context *gin.Context) {
	handleParamRequest(context, r.gateway.DeleteRestart)
}

func (r *InstanceRouter) PostInstanceReload(context *gin.Context) {
	handleBodyRequest(context, r.gateway.Reload)
}

func (r *InstanceRouter) PostInstanceFailover(context *gin.Context) {
	handleBodyRequest(context, r.gateway.Failover)
}

func (r *InstanceRouter) PostInstanceActivate(context *gin.Context) {
	handleBodyRequest(context, r.gateway.Activate)
}

func (r *InstanceRouter) PostInstancePause(context *gin.Context) {
	handleBodyRequest(context, r.gateway.Pause)
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
