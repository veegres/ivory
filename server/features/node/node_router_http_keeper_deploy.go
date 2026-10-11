package node

import "github.com/gin-gonic/gin"

func (r *HttpRouter) PostNodeKeeperDeploy(context *gin.Context) {
	handlePlatformBodyRequest(context, r.service.KeeperDeploy)
}
