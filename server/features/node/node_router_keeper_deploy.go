package node

import "github.com/gin-gonic/gin"

func (r *Router) PostNodeKeeperDeploy(context *gin.Context) {
	handlePlatformBodyRequest(context, r.service.KeeperDeploy)
}
