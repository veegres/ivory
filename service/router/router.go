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
	// TODO think how to move this to UI set up page it can be next after secret page
	companyLabel := os.Getenv("IVORY_COMPANY_LABEL")
	auth := os.Getenv("IVORY_AUTHENTICATION")
	var authHandler gin.HandlerFunc
	if auth == "basic" {
		username := os.Getenv("IVORY_BASIC_USERNAME")
		password := os.Getenv("IVORY_BASIC_PASSWORD")
		if username == "" || password == "" {
			panic("Provide IVORY_BASIC_USERNAME and IVORY_BASIC_PASSWORD when you use IVORY_AUTHENTICATION=basic")
		}
		authHandler = gin.BasicAuth(gin.Accounts{username: password})
	} else {
		auth = "none"
		authHandler = func(context *gin.Context) {}
	}

	r := routes{router: gin.Default()}
	r.router.UseH2C = true

	api := r.router.Group("/api")
	api.GET("/ping", pong)
	api.GET("/info", info(companyLabel, auth))

	authApi := api.Group("/", authHandler)
	r.ProxyGroup(authApi)
	r.ClusterGroup(authApi)
	r.CliGroup(authApi)
	r.CredentialGroup(authApi)

	_ = r.router.Run()
}

func pong(context *gin.Context) {
	context.JSON(http.StatusOK, gin.H{"message": "pong"})
}

func info(company string, auth string) func(context *gin.Context) {
	return func(context *gin.Context) {
		context.JSON(http.StatusOK, gin.H{"response": gin.H{
			"company": company,
			"auth":    auth,
			"secret":  service.Secret.Status(),
		}})
	}
}
