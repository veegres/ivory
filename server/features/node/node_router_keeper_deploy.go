package node

import "github.com/gin-gonic/gin"

func (r *Router) GetNodeKeeperDeploySpec(context *gin.Context) {
	handlePlatformParamRequest(context, r.service.KeeperDeploySpec)
}

func (r *Router) PostNodeKeeperDeploy(context *gin.Context) {
	handlePlatformBodyRequest(context, r.service.KeeperDeploy)
}
