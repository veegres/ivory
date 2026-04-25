package node

import "github.com/gin-gonic/gin"

func (r *Router) GetMetrics(context *gin.Context) {
	handleParamRequestOf(context, r.service.Metrics)
}

func (r *Router) PostDockerDeploy(context *gin.Context) {
	handleBodyRequestOf(context, r.service.DockerDeploy)
}

func (r *Router) PostDockerStop(context *gin.Context) {
	handleBodyRequestOf(context, r.service.DockerStop)
}

func (r *Router) PostDockerRun(context *gin.Context) {
	handleBodyRequestOf(context, r.service.DockerRun)
}

func (r *Router) PostDockerDelete(context *gin.Context) {
	handleBodyRequestOf(context, r.service.DockerDelete)
}

func (r *Router) GetDockerList(context *gin.Context) {
	handleParamRequestOf(context, r.service.DockerList)
}

func (r *Router) GetDockerLogs(context *gin.Context) {
	handleParamRequestOf(context, r.service.DockerLogs)
}
