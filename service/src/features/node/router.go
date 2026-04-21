package node

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

func (r *Router) GetNodeOverview(context *gin.Context) {
	handleParamRequest(context, r.service.Overview)
}

func (r *Router) GetNodeConfig(context *gin.Context) {
	handleParamRequest(context, r.service.Config)
}

func (r *Router) PatchNodeConfig(context *gin.Context) {
	handleBodyRequest(context, r.service.ConfigUpdate)
}

func (r *Router) PostNodeSwitchover(context *gin.Context) {
	handleBodyRequest(context, r.service.Switchover)
}

func (r *Router) DeleteNodeSwitchover(context *gin.Context) {
	handleParamRequest(context, r.service.DeleteSwitchover)
}

func (r *Router) PostNodeReinitialize(context *gin.Context) {
	handleBodyRequest(context, r.service.Reinitialize)
}

func (r *Router) PostNodeRestart(context *gin.Context) {
	handleBodyRequest(context, r.service.Restart)
}

func (r *Router) DeleteNodeRestart(context *gin.Context) {
	handleParamRequest(context, r.service.DeleteRestart)
}

func (r *Router) PostNodeReload(context *gin.Context) {
	handleBodyRequest(context, r.service.Reload)
}

func (r *Router) PostNodeFailover(context *gin.Context) {
	handleBodyRequest(context, r.service.Failover)
}

func (r *Router) PostNodeActivate(context *gin.Context) {
	handleBodyRequest(context, r.service.Activate)
}

func (r *Router) PostNodePause(context *gin.Context) {
	handleBodyRequest(context, r.service.Pause)
}

func (r *Router) GetMetrics(context *gin.Context) {
	handleParamRequestOf(context, r.service.Metrics)
}

func handleParamRequest[T any](context *gin.Context, action func(node Request) (T, int, error)) {
	query := context.Query("request")
	var node Request
	errBind := json.Unmarshal([]byte(query), &node)
	if errBind != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": errBind.Error()})
		return
	}
	body, status, err := action(node)
	handleResponse(context, body, status, err)
}

func handleBodyRequest[T any](context *gin.Context, action func(node Request) (T, int, error)) {
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
