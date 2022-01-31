package router

import (
	"github.com/gin-gonic/gin"
	"ivory/database"
	. "ivory/model"
	"net/http"
)

func (r routes) ClusterGroup(group *gin.RouterGroup) {
	cluster := group.Group("/cluster")
	cluster.GET("", getClusterList)
	cluster.GET("/:host", getClusterByHost)
	cluster.PUT("", putClusterByHost)
	cluster.DELETE("/:host", deleteClusterByHost)
}

func getClusterList(context *gin.Context) {
	context.JSON(http.StatusOK, gin.H{"response": database.ClusterRepository().List()})
}

func getClusterByHost(context *gin.Context) {
	host := context.Param("host")
	cluster := database.ClusterRepository().Get(host)
	if cluster == nil {
		context.JSON(http.StatusNotFound, gin.H{"error": "Not Found"})
	} else {
		context.JSON(http.StatusOK, gin.H{"response": cluster})
	}
}

func putClusterByHost(context *gin.Context) {
	var cluster Cluster
	_ = context.ShouldBindJSON(&cluster)
	database.ClusterRepository().Update(cluster)
	context.JSON(http.StatusOK, gin.H{"response": cluster})
}

func deleteClusterByHost(context *gin.Context) {
	host := context.Param("host")
	database.ClusterRepository().Delete(host)
}
