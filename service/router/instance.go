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
	err := context.ShouldBindQuery(&instance)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	body, status, err := service.PatroniInstanceApiImpl.Info(instance)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	context.JSON(status, gin.H{"response": body})
}

func getInstanceOverview(context *gin.Context) {
	var instance InstanceRequest
	err := context.ShouldBindQuery(&instance)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	body, status, err := service.PatroniInstanceApiImpl.Overview(instance)
	handleResponse(context, body, status, err)
}

func getInstanceConfig(context *gin.Context) {
	var instance InstanceRequest
	err := context.ShouldBindQuery(&instance)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	body, status, err := service.PatroniInstanceApiImpl.Config(instance)
	handleResponse(context, body, status, err)
}

func patchInstanceConfig(context *gin.Context) {
	var instance InstanceRequest
	err := context.ShouldBindJSON(&instance)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	body, status, err := service.PatroniInstanceApiImpl.ConfigUpdate(instance)
	handleResponse(context, body, status, err)
}

func postInstanceSwitchover(context *gin.Context) {
	var instance InstanceRequest
	err := context.ShouldBindJSON(&instance)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	body, status, err := service.PatroniInstanceApiImpl.Switchover(instance)
	handleResponse(context, body, status, err)
}

func postInstanceReinitialize(context *gin.Context) {
	var instance InstanceRequest
	err := context.ShouldBindJSON(&instance)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	body, status, err := service.PatroniInstanceApiImpl.Reinitialize(instance)
	handleResponse(context, body, status, err)
}

func handleResponse(context *gin.Context, body any, status int, err error) {
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if status >= 200 && status < 300 {
		context.JSON(status, gin.H{"response": body})
	} else {
		context.JSON(status, gin.H{"error": body})
	}
}
