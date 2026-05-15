package node

import "github.com/gin-gonic/gin"

func (r *Router) GetMetrics(context *gin.Context) {
	handleParamRequestOf(context, r.service.Metrics)
}

func (r *Router) PostContainerStop(context *gin.Context) {
	handleBodyRequestOf(context, r.service.ContainerStop)
}

func (r *Router) PostContainerRun(context *gin.Context) {
	handleBodyRequestOf(context, r.service.ContainerRun)
}

func (r *Router) PostContainerDelete(context *gin.Context) {
	handleBodyRequestOf(context, r.service.ContainerDelete)
}

func (r *Router) GetContainerList(context *gin.Context) {
	handleParamRequestOf(context, r.service.ContainerList)
}

func (r *Router) GetContainerLogs(context *gin.Context) {
	handleParamRequestOf(context, r.service.ContainerLogs)
}
