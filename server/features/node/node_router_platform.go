package node

import "github.com/gin-gonic/gin"

func (r *Router) GetPlatformSystemMetrics(context *gin.Context) {
	handlePlatformParamRequest(context, r.service.PlatformSystemMetrics)
}

func (r *Router) PostPlatformSystemCopyId(context *gin.Context) {
	handlePlatformBodyRequest(context, r.service.PlatformSystemCopyId)
}

func (r *Router) StreamPlatformSystemLogs(context *gin.Context) {
	handleStreamRequest(context, r.service.PlatformSystemLogs)
}

func (r *Router) GetPlatformSystemProcesses(context *gin.Context) {
	handlePlatformParamRequest(context, r.service.PlatformSystemProcesses)
}

func (r *Router) GetPlatformSystemInfo(context *gin.Context) {
	handlePlatformParamRequest(context, r.service.PlatformSystemInfo)
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
