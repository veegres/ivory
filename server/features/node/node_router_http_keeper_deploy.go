package node

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// PostNodeKeeperDeploy answers 207 for a node that came up without its
// post-scripts, the way PostClusterDeploy does for a partial cluster.
func (r *HttpRouter) PostNodeKeeperDeploy(context *gin.Context) {
	var request KeeperDeployRequest
	if err := context.ShouldBindJSON(&request); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	logs, complete, err := r.service.KeeperDeploy(request)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	status := http.StatusOK
	if !complete {
		status = http.StatusMultiStatus
	}
	context.JSON(status, gin.H{"response": logs})
}
