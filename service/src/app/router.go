package app

import (
	"github.com/gin-gonic/gin"
	"ivory/src/router"
	"net/http"
)

func NewRouter(di *Context) {
	engine := gin.Default()
	engine.UseH2C = true

	unsafe := engine.Group("/api")
	unsafe.GET("/ping", pong)
	unsafe.GET("/info", di.infoRouter.Info)

	initial := unsafe.Group("/", di.secretRouter.EmptyMiddleware())
	initialRouter(initial, di.secretRouter, di.generalRouter)

	login := unsafe.Group("/", di.secretRouter.ExistMiddleware())
	login.POST("/login", di.authRouter.Login)

	safe := login.Group("/", di.authRouter.Middleware())
	safeRouter(safe, di.secretRouter, di.generalRouter)
	clusterRouter(safe, di.clusterRouter)
	bloatRouter(safe, di.bloatRouter)
	certRouter(safe, di.certRouter)
	passwordRouter(safe, di.passwordRouter)
	tagRouter(safe, di.tagRouter)
	instanceRouter(safe, di.instanceRouter)
	queryRouter(safe, di.queryRouter)

	err := engine.Run()
	if err != nil {
		panic(err)
	}
}

func pong(context *gin.Context) {
	context.JSON(http.StatusOK, gin.H{"message": "pong"})
}

func initialRouter(g *gin.RouterGroup, rs *router.SecretRouter, rg *router.GeneralRouter) {
	initial := g.Group("/initial")
	initial.POST("/secret", rs.SetSecret)
	initial.DELETE("/erase", rg.Erase)
}

func safeRouter(g *gin.RouterGroup, rs *router.SecretRouter, rg *router.GeneralRouter) {
	safe := g.Group("/safe")
	safe.GET("/secret", rs.GetStatus)
	safe.POST("/secret", rg.ChangeSecret)
	safe.DELETE("/erase", rg.Erase)
}

func clusterRouter(g *gin.RouterGroup, r *router.ClusterRouter) {
	cluster := g.Group("/cluster")
	cluster.GET("", r.GetClusterList)
	cluster.GET("/:name", r.GetClusterByName)
	cluster.PUT("", r.PutClusterByName)
	cluster.POST("/auto", r.PostClusterAuto)
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

func queryRouter(g *gin.RouterGroup, r *router.QueryRouter) {
	query := g.Group("/query")
	query.GET("", r.GetQueryList)
	query.POST("", r.PostQuery)
	query.PUT("/:uuid", r.PutQuery)
	query.DELETE("/:uuid", r.DeleteQuery)
	query.POST("/run", r.PostRunQuery)
	query.POST("/databases", r.PostDatabasesQuery)
	query.POST("/schemas", r.PostSchemasQuery)
	query.POST("/tables", r.PostTablesQuery)
	query.POST("/chart/common", r.PostCommonChartQuery)
	query.POST("/chart/database", r.PostDatabaseChartQuery)
	query.POST("/cancel", r.PostCancelQuery)
	query.POST("/terminate", r.PostTerminateQuery)
}
