package router

import (
	"github.com/gin-gonic/gin"
	. "ivory/model"
	"ivory/persistence"
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
	list, err := persistence.Database.Cluster.List()
	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	} else {
		context.JSON(http.StatusOK, gin.H{"response": list})
	}
}

func getClusterByHost(context *gin.Context) {
	host := context.Param("host")
	cluster, err := persistence.Database.Cluster.Get(host)
	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	} else {
		context.JSON(http.StatusOK, gin.H{"response": cluster})
	}
}

func putClusterByHost(context *gin.Context) {
	var cluster ClusterModel
	_ = context.ShouldBindJSON(&cluster)
	_ = persistence.Database.Cluster.Update(cluster)
	context.JSON(http.StatusOK, gin.H{"response": cluster})
}

func deleteClusterByHost(context *gin.Context) {
	host := context.Param("host")
	_ = persistence.Database.Cluster.Delete(host)
}
