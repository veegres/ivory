package node

import "github.com/gin-gonic/gin"

func (r *Router) GetNodeOverview(context *gin.Context) {
	handleKeeperParamRequest(context, r.service.KeeperNodeList)
}

func (r *Router) GetNodeConfig(context *gin.Context) {
	handleKeeperParamRequest(context, r.service.KeeperConfigGet)
}

func (r *Router) PatchNodeConfig(context *gin.Context) {
	handleKeeperBodyRequest(context, r.service.KeeperConfigUpdate)
}

func (r *Router) PostNodeSwitchover(context *gin.Context) {
	handleKeeperBodyRequest(context, r.service.KeeperSwitchover)
}

func (r *Router) DeleteNodeSwitchover(context *gin.Context) {
	handleKeeperParamRequest(context, r.service.KeeperSwitchoverDelete)
}

func (r *Router) PostNodeReinitialize(context *gin.Context) {
	handleKeeperBodyRequest(context, r.service.KeeperReinitialize)
}

func (r *Router) PostNodeRestart(context *gin.Context) {
	handleKeeperBodyRequest(context, r.service.KeeperRestart)
}

func (r *Router) DeleteNodeRestart(context *gin.Context) {
	handleKeeperParamRequest(context, r.service.KeeperRestartDelete)
}

func (r *Router) PostNodeReload(context *gin.Context) {
	handleKeeperBodyRequest(context, r.service.KeeperReload)
}

func (r *Router) PostNodeFailover(context *gin.Context) {
	handleKeeperBodyRequest(context, r.service.KeeperFailover)
}

func (r *Router) PostNodeActivate(context *gin.Context) {
	handleKeeperBodyRequest(context, r.service.KeeperActivate)
}

func (r *Router) PostNodePause(context *gin.Context) {
	handleKeeperBodyRequest(context, r.service.KeeperPause)
}
