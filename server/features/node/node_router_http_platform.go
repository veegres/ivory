package node

import "github.com/gin-gonic/gin"

func (r *HttpRouter) GetPlatformSystemMetrics(context *gin.Context) {
	handlePlatformParamRequest(context, r.service.PlatformSystemMetrics)
}

func (r *HttpRouter) PostPlatformSystemCopyId(context *gin.Context) {
	handlePlatformBodyRequest(context, r.service.PlatformSystemCopyId)
}

func (r *HttpRouter) StreamPlatformSystemLogs(context *gin.Context) {
	handleStreamRequest(context, r.service.PlatformSystemLogs)
}

func (r *HttpRouter) GetPlatformSystemProcesses(context *gin.Context) {
	handlePlatformParamRequest(context, r.service.PlatformSystemProcesses)
}

func (r *HttpRouter) GetPlatformSystemInfo(context *gin.Context) {
	handlePlatformParamRequest(context, r.service.PlatformSystemInfo)
}

func (r *HttpRouter) PostPlatformContainerUp(context *gin.Context) {
	handlePlatformBodyRequest(context, r.service.PlatformContainerUp)
}

func (r *HttpRouter) PostPlatformContainerDown(context *gin.Context) {
	handlePlatformBodyRequest(context, r.service.PlatformContainerDown)
}

func (r *HttpRouter) PostPlatformContainerStart(context *gin.Context) {
	handlePlatformBodyRequest(context, r.service.PlatformContainerStart)
}

func (r *HttpRouter) PostPlatformContainerStop(context *gin.Context) {
	handlePlatformBodyRequest(context, r.service.PlatformContainerStop)
}

func (r *HttpRouter) PostPlatformContainerRestart(context *gin.Context) {
	handlePlatformBodyRequest(context, r.service.PlatformContainerRestart)
}

func (r *HttpRouter) GetPlatformContainerList(context *gin.Context) {
	handlePlatformParamRequest(context, r.service.PlatformContainerList)
}

func (r *HttpRouter) StreamPlatformContainerLogs(context *gin.Context) {
	handleStreamRequest(context, r.service.PlatformContainerLogs)
}

func (r *HttpRouter) GetPlatformContainerMetrics(context *gin.Context) {
	handlePlatformParamRequest(context, r.service.PlatformContainerMetrics)
}

func (r *HttpRouter) GetPlatformContainerInspect(context *gin.Context) {
	handlePlatformParamRequest(context, r.service.PlatformContainerInspect)
}
