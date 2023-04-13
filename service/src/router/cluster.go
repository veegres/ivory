package router

import (
	"github.com/gin-gonic/gin"
	. "ivory/src/model"
	"ivory/src/persistence"
	"ivory/src/service"
	"net/http"
)

type ClusterRouter struct {
	clusterService *service.ClusterService
	tagRepository  *persistence.TagRepository
}

func NewClusterRouter(
	clusterService *service.ClusterService,
	tagRepository *persistence.TagRepository,
) *ClusterRouter {
	return &ClusterRouter{
		clusterService: clusterService,
		tagRepository:  tagRepository,
	}
}

func (r *ClusterRouter) GetClusterList(context *gin.Context) {
	tags := context.Request.URL.Query()["tags[]"]

	if len(tags) == 0 {
		list, err := r.clusterService.List()
		if err != nil {
			context.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		context.JSON(http.StatusOK, gin.H{"response": list})
	} else {
		listMap := make(map[string]bool)
		for _, tag := range tags {
			clusters, err := r.tagRepository.Get(tag)
			if err != nil {
				context.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
				return
			}

			for _, c := range clusters {
				if !listMap[c] {
					listMap[c] = true
				}
			}
		}

		listName := make([]string, 0)
		for k := range listMap {
			listName = append(listName, k)
		}

		list, err := r.clusterService.ListByName(listName)
		if err != nil {
			context.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		context.JSON(http.StatusOK, gin.H{"response": list})
	}

}

func (r *ClusterRouter) GetClusterByName(context *gin.Context) {
	name := context.Param("name")
	cluster, err := r.clusterService.Get(name)
	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, gin.H{"response": cluster})
}

func (r *ClusterRouter) PutClusterByName(context *gin.Context) {
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
	err := r.tagRepository.UpdateCluster(cluster.Name, tagList)
	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	cluster.Tags = tagList

	// NOTE: update cluster
	errCluster := r.clusterService.Update(cluster)
	if errCluster != nil {
		context.JSON(http.StatusNotFound, gin.H{"error": errCluster.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"response": cluster})
}

func (r *ClusterRouter) DeleteClusterByName(context *gin.Context) {
	name := context.Param("name")
	errTag := r.tagRepository.UpdateCluster(name, nil)
	if errTag != nil {
		context.JSON(http.StatusNotFound, gin.H{"error": errTag.Error()})
		return
	}
	errCluster := r.clusterService.Delete(name)
	if errCluster != nil {
		context.JSON(http.StatusNotFound, gin.H{"error": errCluster.Error()})
		return
	}
	context.JSON(http.StatusOK, gin.H{"response": "deleted"})
}
