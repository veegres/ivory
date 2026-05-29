package node

import "github.com/gin-gonic/gin"

func (r *Router) GetMetrics(context *gin.Context) {
	handleParamRequestOf(context, r.service.PlatformMetrics)
}

func (r *Router) GetPlatformStop(context *gin.Context) {
	handleStreamParamRequestOf(context, r.service.PlatformStop)
}

func (r *Router) GetPlatformDeploy(context *gin.Context) {
	handleStreamParamRequestOf(context, r.service.PlatformDeploy)
}

func (r *Router) GetPlatformDelete(context *gin.Context) {
	handleStreamParamRequestOf(context, r.service.PlatformDelete)
}

func (r *Router) GetPlatformList(context *gin.Context) {
	handleStreamParamRequestOf(context, r.service.PlatformList)
}

func (r *Router) GetPlatformLogs(context *gin.Context) {
	handleStreamParamRequestOf(context, r.service.PlatformLogs)
}
