package router

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type routes struct {
	router *gin.Engine
}

func Start() {
	r := routes{router: gin.Default()}

	api := r.router.Group("/api")
	api.GET("/ping", pong)
	r.ProxyGroup(api)
	r.ClusterGroup(api)

	_ = r.router.Run()
}

func pong(context *gin.Context) { context.JSON(http.StatusOK, gin.H{"message": "pong"}) }
