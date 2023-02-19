package app

import (
	"github.com/gin-gonic/gin"
	"ivory/src/router"
	"net/http"
)

func NewRouter(di *Context) {
	authHandler := func(context *gin.Context) {}
	if di.env.Username != nil && di.env.Password != nil {
		authHandler = gin.BasicAuth(gin.Accounts{*di.env.Username: *di.env.Password})
	}

	engine := gin.Default()
	engine.UseH2C = true

	unsafe := engine.Group("/api")
	unsafe.GET("/ping", pong)
	infoRouter(unsafe, di.infoRouter)

	// TODO we shouldn't allow any request if there is no secret set think of creating mw
	safe := unsafe.Group("/", authHandler)
	clusterRouter(safe, di.clusterRouter)
	bloatRouter(safe, di.bloatRouter)
	certRouter(safe, di.certRouter)
	secretRouter(safe, di.secretRouter)
	passwordRouter(safe, di.passwordRouter)
	tagRouter(safe, di.tagRouter)
	instanceRouter(safe, di.instanceRouter)

	err := engine.Run()
	if err != nil {
		panic(err)
	}
}

func pong(context *gin.Context) {
	context.JSON(http.StatusOK, gin.H{"message": "pong"})
}

func infoRouter(g *gin.RouterGroup, r *router.InfoRouter) {
	info := g.Group("/info")
	info.GET("", r.Info)
}

func clusterRouter(g *gin.RouterGroup, r *router.ClusterRouter) {
	cluster := g.Group("/cluster")
	cluster.GET("", r.GetClusterList)
	cluster.GET("/:name", r.GetClusterByName)
	cluster.PUT("", r.PutClusterByName)
	cluster.DELETE("/:name", r.DeleteClusterByName)
}

func bloatRouter(g *gin.RouterGroup, r *router.BloatRouter) {
	cli := g.Group("/cli")
	cli.GET("/bloat", r.GetCompactTableList)
	cli.GET("/bloat/:uuid", r.GetCompactTable)
	cli.GET("/bloat/:uuid/logs", r.GetCompactTableLogs)
	cli.GET("/bloat/cluster/:name", r.GetCompactTableListByCluster)
	cli.POST("/bloat/job/start", r.StartJob)
	cli.POST("/bloat/job/:uuid/stop", r.StopJob)
	cli.DELETE("/bloat/job/:uuid/delete", r.DeleteJob)
	cli.GET("/bloat/job/:uuid/stream", r.StreamJob)
}

func certRouter(g *gin.RouterGroup, r *router.CertRouter) {
	cert := g.Group("/cert")
	cert.GET("", r.GetCertList)
	cert.DELETE("/:uuid", r.DeleteCert)
	cert.POST("/upload", r.PostUploadCert)
	cert.POST("/add", r.PostAddCert)
}

func secretRouter(g *gin.RouterGroup, r *router.SecretRouter) {
	secret := g.Group("/secret")
	secret.GET("", r.GetStatus)
	secret.POST("/set", r.SetSecret)
	secret.POST("/update", r.UpdateSecret)
	secret.POST("/clean", r.CleanSecret)
}
func passwordRouter(g *gin.RouterGroup, r *router.PasswordRouter) {
	password := g.Group("/password")
	password.GET("", r.GetCredentials)
	password.POST("", r.PostCredential)
	password.PATCH("/:uuid", r.PatchCredential)
	password.DELETE("/:uuid", r.DeleteCredential)
}

func tagRouter(g *gin.RouterGroup, r *router.TagRouter) {
	tag := g.Group("/tag")
	tag.GET("", r.GetTagList)
}

func instanceRouter(g *gin.RouterGroup, r *router.InstanceRouter) {
	instance := g.Group("/instance")
	instance.GET("/info", r.GetInstanceInfo)
	instance.GET("/overview", r.GetInstanceOverview)
	instance.GET("/config", r.GetInstanceConfig)
	instance.PATCH("/config", r.PatchInstanceConfig)
	instance.POST("/switchover", r.PostInstanceSwitchover)
	instance.POST("/reinitialize", r.PostInstanceReinitialize)

}
