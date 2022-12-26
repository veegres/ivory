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
	cluster.GET("/:name", getClusterByName)
	cluster.PUT("", putClusterByName)
	cluster.DELETE("/:name", deleteClusterByName)
}

func getClusterList(context *gin.Context) {
	list, err := persistence.Database.Cluster.List()
	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, gin.H{"response": list})
}

func getClusterByName(context *gin.Context) {
	name := context.Param("name")
	cluster, err := persistence.Database.Cluster.Get(name)
	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, gin.H{"response": cluster})
}

func putClusterByName(context *gin.Context) {
	var cluster ClusterModel
	_ = context.ShouldBindJSON(&cluster)
	_ = persistence.Database.Cluster.Update(cluster)
	context.JSON(http.StatusOK, gin.H{"response": cluster})
}

func deleteClusterByName(context *gin.Context) {
	name := context.Param("name")
	err := persistence.Database.Cluster.Delete(name)
	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, gin.H{"response": "deleted"})
}
