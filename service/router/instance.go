package router

import (
	"github.com/gin-gonic/gin"
    . "ivory/model"
    "ivory/service"
    "net/http"
)

func (r routes) ProxyGroup(group *gin.RouterGroup) {
	node := group.Group("/:cluster/instance")
	node.GET("/cluster", getNodeCluster)
	node.GET("/overview", getNodeOverview)
	node.GET("/config", getNodeConfig)
	node.PATCH("/config", patchNodeConfig)
	node.POST("/switchover", postNodeSwitchover)
	node.POST("/reinitialize", postNodeReinitialize)
}

func getNodeCluster(context *gin.Context) {
    cluster := context.Param("cluster")
    var instanse Instance
    err := context.ShouldBindJSON(&instanse)

    body, err := service.Patroni.Cluster(cluster, instanse)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, gin.H{"response": body})
}

func getNodeOverview(context *gin.Context) {
    cluster := context.Param("cluster")
    var instanse Instance
    err := context.ShouldBindJSON(&instanse)

    body, err := service.Patroni.Overview(cluster, instanse)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, gin.H{"response": body})
}

func getNodeConfig(context *gin.Context) {
    cluster := context.Param("cluster")
    var instanse Instance
    err := context.ShouldBindJSON(&instanse)

    body, err := service.Patroni.Config(cluster, instanse)
    if err != nil {
        context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    context.JSON(http.StatusOK, gin.H{"response": body})
}

func patchNodeConfig(context *gin.Context) {
    cluster := context.Param("cluster")
    var instanse Instance
    err := context.ShouldBindJSON(&instanse)

    body, err := service.Patroni.ConfigUpdate(cluster, instanse)
    if err != nil {
        context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    context.JSON(http.StatusOK, gin.H{"response": body})
}

func postNodeSwitchover(context *gin.Context) {
    cluster := context.Param("cluster")
    var instanse Instance
    err := context.ShouldBindJSON(&instanse)

    body, err := service.Patroni.Switchover(cluster, instanse)
    if err != nil {
        context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    context.JSON(http.StatusOK, gin.H{"response": body})
}

func postNodeReinitialize(context *gin.Context) {
    cluster := context.Param("cluster")
    var instanse Instance
    err := context.ShouldBindJSON(&instanse)

    body, err := service.Patroni.Reinitialize(cluster, instanse)
    if err != nil {
        context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    context.JSON(http.StatusOK, gin.H{"response": body})
}
