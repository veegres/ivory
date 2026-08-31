package cluster

import (
	"net/http"
	"strconv"

	"ivory/features/node"
	"ivory/features/query"

	"github.com/gin-gonic/gin"
)

type Router struct {
	clusterService *Service
}

func NewRouter(clusterService *Service) *Router {
	return &Router{clusterService: clusterService}
}

func (r *Router) GetClusterList(context *gin.Context) {
	request := SearchRequest{Tags: context.Request.URL.Query()["tags[]"]}
	if keeper := context.Query("keeper"); keeper != "" {
		plugin := node.KeeperPlugin(keeper)
		request.Keeper = &plugin
	}
	if database := context.Query("database"); database != "" {
		plugin := query.DbPlugin(database)
		request.Database = &plugin
	}

	list, err := r.clusterService.Search(request)
	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, gin.H{"response": list})
}

func (r *Router) GetClusterOverview(context *gin.Context) {
	name := context.Param("name")
	host := context.Query("host")
	port := context.Query("port")

	// NOTE: we need here default value because it is the same as empty string
	if port == "" {
		port = "0"
	}
	parsedPort, err := strconv.Atoi(port)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	overview, err := r.clusterService.Overview(name, host, parsedPort)
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
	var cluster Request
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

func (r *Router) PostClusterDetect(context *gin.Context) {
	var cluster CreateAutoRequest
	errParse := context.ShouldBindJSON(&cluster)
	if errParse != nil {
		context.JSON(http.StatusNotFound, gin.H{"error": errParse.Error()})
		return
	}

	response, errRes := r.clusterService.Detect(cluster)
	if errRes != nil {
		context.JSON(http.StatusNotFound, gin.H{"error": errRes.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"response": response})
}

func (r *Router) PostClusterFix(context *gin.Context) {
	name := context.Param("name")
	response, errRes := r.clusterService.Fix(name)
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

func (r *Router) PostClusterDeploy(context *gin.Context) {
	var request DeployRequest
	errParse := context.ShouldBindJSON(&request)
	if errParse != nil {
		context.JSON(http.StatusNotFound, gin.H{"error": errParse.Error()})
		return
	}

	response, complete, errRes := r.clusterService.Deploy(request)
	if errRes != nil {
		context.JSON(http.StatusNotFound, gin.H{"error": errRes.Error()})
		return
	}

	// NOTE: a deploy that started every node and initialized them is the only
	// one that answers 200. Where a node or its post-script failed the logs are
	// still the whole report, so the body is unchanged and the status is what
	// says not to read it as a success - an error is reserved for what Ivory
	// itself rejected, and a container the operator's own command killed is not
	// that.
	status := http.StatusOK
	if !complete {
		status = http.StatusMultiStatus
	}
	context.JSON(status, gin.H{"response": response})
}
