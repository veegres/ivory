package router

import (
	"github.com/gin-gonic/gin"
	. "ivory/model"
	"ivory/service"
	"net/http"
)

func (r routes) ProxyGroup(group *gin.RouterGroup) {
	node := group.Group("/instance")
	node.GET("/info", getInstanceInfo)
	node.GET("/overview", getInstanceOverview)
	node.GET("/config", getInstanceConfig)
	node.PATCH("/config", patchInstanceConfig)
	node.POST("/switchover", postInstanceSwitchover)
	node.POST("/reinitialize", postInstanceReinitialize)
}

func getInstanceInfo(context *gin.Context) {
	var instance InstanceRequest
	err := context.ShouldBindJSON(&instance)

	body, err := service.Patroni.Info(instance)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, gin.H{"response": body})
}

func getInstanceOverview(context *gin.Context) {
	var instance InstanceRequest
	err := context.ShouldBindJSON(&instance)

	body, err := service.Patroni.Overview(instance)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, gin.H{"response": body})
}

func getInstanceConfig(context *gin.Context) {
	var instance InstanceRequest
	err := context.ShouldBindJSON(&instance)

	body, err := service.Patroni.Config(instance)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, gin.H{"response": body})
}

func patchInstanceConfig(context *gin.Context) {
	var instance InstanceRequest
	err := context.ShouldBindJSON(&instance)

	body, err := service.Patroni.ConfigUpdate(instance)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, gin.H{"response": body})
}

func postInstanceSwitchover(context *gin.Context) {
	var instance InstanceRequest
	err := context.ShouldBindJSON(&instance)

	body, err := service.Patroni.Switchover(instance)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, gin.H{"response": body})
}

func postInstanceReinitialize(context *gin.Context) {
	var instance InstanceRequest
	err := context.ShouldBindJSON(&instance)

	body, err := service.Patroni.Reinitialize(instance)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, gin.H{"response": body})
}
