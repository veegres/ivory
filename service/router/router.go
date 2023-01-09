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
		user := os.Getenv("IVORY_BASIC_USER")
		password := os.Getenv("IVORY_BASIC_PASSWORD")
		if user == "" || password == "" {
			panic("Provide IVORY_BASIC_USER and IVORY_BASIC_PASSWORD when you use IVORY_AUTHENTICATION=basic")
		}
		authHandler = gin.BasicAuth(gin.Accounts{user: password})
	} else {
		auth = "none"
		authHandler = func(context *gin.Context) {}
	}

	r := routes{router: gin.Default()}
	r.router.UseH2C = true

	api := r.router.Group("/api")
	api.GET("/ping", pong)
	api.GET("/info", info(companyLabel, auth))

	// TODO we shouldn't allow any request if there is no secret set think of creating mw
	authApi := api.Group("/", authHandler)
	r.ProxyGroup(authApi)
	r.ClusterGroup(authApi)
	r.CliGroup(authApi)
	r.CredentialGroup(authApi)
	r.CertGroup(authApi)

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
