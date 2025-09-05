package instance

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Router struct {
	service *Service
}

func NewInstanceRouter(service *Service) *Router {
	return &Router{service: service}
}

func (r *Router) GetInstanceOverview(context *gin.Context) {
	handleParamRequest(context, r.service.Overview)
}

func (r *Router) GetInstanceConfig(context *gin.Context) {
	handleParamRequest(context, r.service.Config)
}

func (r *Router) PatchInstanceConfig(context *gin.Context) {
	handleBodyRequest(context, r.service.ConfigUpdate)
}

func (r *Router) PostInstanceSwitchover(context *gin.Context) {
	handleBodyRequest(context, r.service.Switchover)
}

func (r *Router) DeleteInstanceSwitchover(context *gin.Context) {
	handleParamRequest(context, r.service.DeleteSwitchover)
}

func (r *Router) PostInstanceReinitialize(context *gin.Context) {
	handleBodyRequest(context, r.service.Reinitialize)
}

func (r *Router) PostInstanceRestart(context *gin.Context) {
	handleBodyRequest(context, r.service.Restart)
}

func (r *Router) DeleteInstanceRestart(context *gin.Context) {
	handleParamRequest(context, r.service.DeleteRestart)
}

func (r *Router) PostInstanceReload(context *gin.Context) {
	handleBodyRequest(context, r.service.Reload)
}

func (r *Router) PostInstanceFailover(context *gin.Context) {
	handleBodyRequest(context, r.service.Failover)
}

func (r *Router) PostInstanceActivate(context *gin.Context) {
	handleBodyRequest(context, r.service.Activate)
}

func (r *Router) PostInstancePause(context *gin.Context) {
	handleBodyRequest(context, r.service.Pause)
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
