package node

import "github.com/gin-gonic/gin"

func (r *HttpRouter) GetNodeOverview(context *gin.Context) {
	handleKeeperParamRequest(context, r.service.KeeperNodeList)
}

func (r *HttpRouter) GetNodeConfig(context *gin.Context) {
	handleKeeperParamRequest(context, r.service.KeeperConfigGet)
}

func (r *HttpRouter) PatchNodeConfig(context *gin.Context) {
	handleKeeperBodyRequest(context, r.service.KeeperConfigUpdate)
}

func (r *HttpRouter) PostNodeSwitchover(context *gin.Context) {
	handleKeeperBodyRequest(context, r.service.KeeperSwitchover)
}

func (r *HttpRouter) DeleteNodeSwitchover(context *gin.Context) {
	handleKeeperParamRequest(context, r.service.KeeperSwitchoverDelete)
}

func (r *HttpRouter) PostNodeReinitialize(context *gin.Context) {
	handleKeeperBodyRequest(context, r.service.KeeperReinitialize)
}

func (r *HttpRouter) PostNodeRestart(context *gin.Context) {
	handleKeeperBodyRequest(context, r.service.KeeperRestart)
}

func (r *HttpRouter) DeleteNodeRestart(context *gin.Context) {
	handleKeeperParamRequest(context, r.service.KeeperRestartDelete)
}

func (r *HttpRouter) PostNodeReload(context *gin.Context) {
	handleKeeperBodyRequest(context, r.service.KeeperReload)
}

func (r *HttpRouter) PostNodeFailover(context *gin.Context) {
	handleKeeperBodyRequest(context, r.service.KeeperFailover)
}

func (r *HttpRouter) PostNodeActivate(context *gin.Context) {
	handleKeeperBodyRequest(context, r.service.KeeperActivate)
}

func (r *HttpRouter) PostNodePause(context *gin.Context) {
	handleKeeperBodyRequest(context, r.service.KeeperPause)
}
