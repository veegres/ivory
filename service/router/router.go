package router

import (
	"github.com/gin-gonic/gin"
	"ivory/service"
	"net/http"
	"os"
)

type routes struct {
	router *gin.Engine
}

func Start() {
	r := routes{router: gin.Default()}
	r.router.UseH2C = true

	api := r.router.Group("/api")
	api.GET("/ping", pong)
	api.GET("/info", info)
	r.ProxyGroup(api)
	r.ClusterGroup(api)
	r.CliGroup(api)
	r.CredentialGroup(api)

	_ = r.router.Run()
}

func pong(context *gin.Context) { context.JSON(http.StatusOK, gin.H{"message": "pong"}) }
func info(context *gin.Context) {
	context.JSON(http.StatusOK, gin.H{"response": gin.H{
		"company": os.Getenv("IVORY_COMPANY_LABEL"),
		"secret":  service.Secret.Status(),
	}})
}
