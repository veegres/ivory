package cluster

import (
	"encoding/json"
	"ivory/src/clients/sidecar"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Router struct {
	clusterService *Service
}

func NewRouter(clusterService *Service) *Router {
	return &Router{clusterService: clusterService}
}

func (r *Router) GetClusterList(context *gin.Context) {
	tags := context.Request.URL.Query()["tags[]"]

	var list []Cluster
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

func (r *Router) GetClusterOverview(context *gin.Context) {
	name := context.Param("name")
	query := context.Query("sidecar")
	var side *sidecar.Sidecar
	if query != "" {
		errBind := json.Unmarshal([]byte(query), &side)
		if errBind != nil {
			context.JSON(http.StatusBadRequest, gin.H{"error": errBind.Error()})
			return
		}
	}

	overview, err := r.clusterService.Overview(name, side)
	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, gin.H{"response": overview})
}

func (r *Router) GetClusterByName(context *gin.Context) {
	name := context.Param("name")
	cluster, err := r.clusterService.Get(name)
	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, gin.H{"response": cluster})
}

func (r *Router) PutClusterByName(context *gin.Context) {
	var cluster Cluster
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

func (r *Router) PostClusterAutoCreate(context *gin.Context) {
	var cluster ClusterAuto
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

func (r *Router) PostClusterAutoFix(context *gin.Context) {
	name := context.Param("name")
	response, errRes := r.clusterService.FixAuto(name)
	if errRes != nil {
		context.JSON(http.StatusNotFound, gin.H{"error": errRes.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"response": response})
}

func (r *Router) DeleteClusterByName(context *gin.Context) {
	name := context.Param("name")
	err := r.clusterService.Delete(name)
	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, gin.H{"response": "deleted"})
}
