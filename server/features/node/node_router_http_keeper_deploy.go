package node

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// PostNodeKeeperDeploy answers 200 only when every step is done and 207
// otherwise, the way PostClusterDeploy does.
func (r *HttpRouter) PostNodeKeeperDeploy(context *gin.Context) {
	var request KeeperDeployRequest
	if err := context.ShouldBindJSON(&request); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	response, err := r.service.KeeperDeploy(request)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	status := http.StatusOK
	if !response.Steps.Done() {
		status = http.StatusMultiStatus
	}
	context.JSON(status, gin.H{"response": response})
}
