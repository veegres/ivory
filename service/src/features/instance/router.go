package instance

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

func (r *Router) GetMetrics(context *gin.Context) {
	handleParamRequestOf(context, r.service.Metrics)
}

func handleParamRequest[T any](context *gin.Context, action func(instance Request) (T, int, error)) {
	query := context.Query("request")
	var instance Request
	errBind := json.Unmarshal([]byte(query), &instance)
	if errBind != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": errBind.Error()})
		return
	}
	body, status, err := action(instance)
	handleResponse(context, body, status, err)
}

func handleBodyRequest[T any](context *gin.Context, action func(instance Request) (T, int, error)) {
	handleBodyRequestOf(context, action)
}

func handleParamRequestOf[R any, T any](context *gin.Context, action func(request R) (T, int, error)) {
	query := context.Query("request")
	var request R
	errBind := json.Unmarshal([]byte(query), &request)
	if errBind != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": errBind.Error()})
		return
	}
	body, status, err := action(request)
	handleResponse(context, body, status, err)
}

func handleBodyRequestOf[R any, T any](context *gin.Context, action func(request R) (T, int, error)) {
	var request R
	errBind := context.ShouldBindJSON(&request)
	if errBind != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": errBind.Error()})
		return
	}
	body, status, err := action(request)
	handleResponse(context, body, status, err)
}

func handleResponse(context *gin.Context, body any, status int, err error) {
	if err != nil {
		context.JSON(status, gin.H{"error": err.Error()})
		return
	}
	context.JSON(status, gin.H{"response": body})
}
