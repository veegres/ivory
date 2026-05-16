package node

import "github.com/gin-gonic/gin"

func (r *Router) GetMetrics(context *gin.Context) {
	handleParamRequestOf(context, r.service.Metrics)
}

func (r *Router) PostPlatformStop(context *gin.Context) {
	handleBodyRequestOf(context, r.service.PlatformStop)
}

func (r *Router) PostPlatformDeploy(context *gin.Context) {
	handleBodyRequestOf(context, r.service.PlatformDeploy)
}

func (r *Router) PostPlatformDelete(context *gin.Context) {
	handleBodyRequestOf(context, r.service.PlatformDelete)
}

func (r *Router) GetPlatformList(context *gin.Context) {
	handleParamRequestOf(context, r.service.PlatformList)
}

func (r *Router) GetPlatformLogs(context *gin.Context) {
	handleParamRequestOf(context, r.service.PlatformLogs)
}
