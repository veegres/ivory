package router

import (
	"github.com/gin-gonic/gin"
	. "ivory/src/model"
	"ivory/src/service"
	"net/http"
)

type ClusterRouter struct {
	clusterService *service.ClusterService
}

func NewClusterRouter(clusterService *service.ClusterService) *ClusterRouter {
	return &ClusterRouter{clusterService: clusterService}
}

func (r *ClusterRouter) GetClusterList(context *gin.Context) {
	tags := context.Request.URL.Query()["tags[]"]

	var list []ClusterModel
	var err error
	if len(tags) == 0 {
		list, err = r.clusterService.List()
	} else {
		list, err = r.clusterService.ListByTag(tags)
	}
	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, gin.H{"response": list})
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

	response, errRes := r.clusterService.Update(cluster)
	if errRes != nil {
		context.JSON(http.StatusNotFound, gin.H{"error": errRes.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"response": response})
}

func (r *ClusterRouter) PostClusterAuto(context *gin.Context) {
	var cluster ClusterAutoModel
	errParse := context.ShouldBindJSON(&cluster)
	if errParse != nil {
		context.JSON(http.StatusNotFound, gin.H{"error": errParse.Error()})
		return
	}

	response, errRes := r.clusterService.CreateAuto(cluster)
	if errRes != nil {
		context.JSON(http.StatusNotFound, gin.H{"error": errRes.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"response": response})
}

func (r *ClusterRouter) DeleteClusterByName(context *gin.Context) {
	name := context.Param("name")
	err := r.clusterService.Delete(name)
	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, gin.H{"response": "deleted"})
}
