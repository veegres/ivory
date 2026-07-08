package node

import "github.com/gin-gonic/gin"

func (r *Router) GetPlatformMetrics(context *gin.Context) {
	handlePlatformParamRequest(context, r.service.PlatformVmMetrics)
}

func (r *Router) PostPlatformCopyId(context *gin.Context) {
	handlePlatformBodyRequest(context, r.service.PlatformVmCopyId)
}

func (r *Router) StreamPlatformLogs(context *gin.Context) {
	handleStreamRequest(context, r.service.PlatformLogs)
}

func (r *Router) GetPlatformProcesses(context *gin.Context) {
	handlePlatformParamRequest(context, r.service.PlatformVmProcesses)
}

func (r *Router) GetPlatformInfo(context *gin.Context) {
	handlePlatformParamRequest(context, r.service.PlatformVmInfo)
}

func (r *Router) GetPlatformContainerDeployOptions(context *gin.Context) {
	handlePlatformParamRequest(context, r.service.PlatformContainerDeployOptions)
}

func (r *Router) PostPlatformContainerUp(context *gin.Context) {
	handlePlatformBodyRequest(context, r.service.PlatformContainerUp)
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

func (r *Router) GetPlatformContainerMetrics(context *gin.Context) {
	handlePlatformParamRequest(context, r.service.PlatformContainerMetrics)
}
