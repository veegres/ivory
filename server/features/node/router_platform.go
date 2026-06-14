package node

import "github.com/gin-gonic/gin"

func (r *Router) GetPlatformMetrics(context *gin.Context) {
	handlePlatformParamRequest(context, r.service.PlatformMetrics)
}

func (r *Router) PostPlatformCopyId(context *gin.Context) {
	handlePlatformBodyRequest(context, r.service.PlatformCopyId)
}

func (r *Router) StreamPlatformContainerUp(context *gin.Context) {
	handleStreamRequest(context, r.service.PlatformContainerUp)
}

func (r *Router) PostPlatformContainerDown(context *gin.Context) {
	handlePlatformBodyRequest(context, r.service.PlatformContainerDown)
}

func (r *Router) PostPlatformContainerStart(context *gin.Context) {
	handlePlatformBodyRequest(context, r.service.PlatformContainerStart)
}

func (r *Router) PostPlatformContainerStop(context *gin.Context) {
	handlePlatformBodyRequest(context, r.service.PlatformContainerStop)
}

func (r *Router) PostPlatformContainerRestart(context *gin.Context) {
	handlePlatformBodyRequest(context, r.service.PlatformContainerRestart)
}

func (r *Router) GetPlatformContainerList(context *gin.Context) {
	handlePlatformParamRequest(context, r.service.PlatformContainerList)
}

func (r *Router) StreamPlatformContainerLogs(context *gin.Context) {
	handleStreamRequest(context, r.service.PlatformContainerLogs)
}
