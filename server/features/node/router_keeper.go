package node

import "github.com/gin-gonic/gin"

func (r *Router) GetNodeOverview(context *gin.Context) {
	handleParamRequest(context, r.service.List)
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
