package router

import (
	"github.com/gin-gonic/gin"
	. "ivory/model"
	"ivory/persistence"
	"net/http"
)

func (r routes) ClusterListGroup(group *gin.RouterGroup) {
	cluster := group.Group("/cluster")
	cluster.GET("", getClusterList)
	cluster.GET("/:name", getClusterByName)
	cluster.PUT("", putClusterByName)
	cluster.DELETE("/:name", deleteClusterByName)
}

func getClusterList(context *gin.Context) {
	tags := context.Request.URL.Query()["tags"]

	listMap := make(map[string]bool)
	list := make([]ClusterModel, 0)
	if len(tags) == 0 {
		listAll, err := persistence.BoltDB.Cluster.List()
		if err != nil {
			context.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		list = listAll
	} else {
		for _, tag := range tags {
			tagClusters, err := persistence.BoltDB.Tag.Get(tag)
			if err != nil {
				context.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
				return
			}

			clusters, err := persistence.BoltDB.Cluster.ListByName(tagClusters)
			if err != nil {
				context.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
				return
			}

			for _, cluster := range clusters {
				if !listMap[cluster.Name] {
					listMap[cluster.Name] = true
					list = append(list, cluster)
				}
			}
		}
	}

	context.JSON(http.StatusOK, gin.H{"response": list})
}

func getClusterByName(context *gin.Context) {
	name := context.Param("name")
	cluster, err := persistence.BoltDB.Cluster.Get(name)
	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, gin.H{"response": cluster})
}

func putClusterByName(context *gin.Context) {
	var cluster ClusterModel
	errParse := context.ShouldBindJSON(&cluster)
	if errParse != nil {
		context.JSON(http.StatusNotFound, gin.H{"error": errParse.Error()})
		return
	}

	// NOTE: remove duplicates
	tagMap := make(map[string]bool)
	for _, tag := range cluster.Tags {
		tagMap[tag] = true
	}
	tagList := make([]string, 0)
	for key := range tagMap {
		tagList = append(tagList, key)
	}

	// NOTE: create tags in db with cluster name
	err := persistence.BoltDB.Tag.UpdateCluster(cluster.Name, tagList)
	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	cluster.Tags = tagList

	// NOTE: update cluster
	errCluster := persistence.BoltDB.Cluster.Update(cluster)
	if errCluster != nil {
		context.JSON(http.StatusNotFound, gin.H{"error": errCluster.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"response": cluster})
}

func deleteClusterByName(context *gin.Context) {
	name := context.Param("name")
	err := persistence.BoltDB.Cluster.Delete(name)
	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, gin.H{"response": "deleted"})
}
