package app

import (
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"ivory/src/router"
	"log/slog"
	"net/http"
)

func NewRouter(di *Context) {
	engine := gin.Default()
	engine.UseH2C = true

	// NOTE: Serving ivory static files to web
	if di.env.Config.StaticFilesPath != "" {
		slog.Info("Ivory serves static files from '" + di.env.Config.StaticFilesPath + "' under sub path '" + di.env.Config.UrlPath + "'")
		engine.Use(static.Serve(di.env.Config.UrlPath, static.LocalFile(di.env.Config.StaticFilesPath, true)))
		engine.NoRoute(func(context *gin.Context) {
			// NOTE: if files wasn't found and NoRoute come here we need throw 404 and prevent endless redirect
			if context.Request.URL.Path != di.env.Config.UrlPath {
				context.Redirect(http.StatusMovedPermanently, di.env.Config.UrlPath)
			}
		})
	}

	// NOTE: Setup default sub path for reverse proxies, default "/"
	slog.Info("Ivory serves backend api '/api' under sub path '" + di.env.Config.UrlPath + "'")
	path := engine.Group(di.env.Config.UrlPath)

	unsafe := path.Group("/api", gin.Recovery())
	unsafe.GET("/ping", pong)
	unsafe.GET("/info", di.generalRouter.GetAppInfo)

	initial := unsafe.Group("/", di.secretRouter.EmptyMiddleware())
	initialRouter(initial, di.secretRouter, di.generalRouter)

	general := unsafe.Group("/", di.secretRouter.ExistMiddleware())
	generalRouter(general, di.authRouter, di.generalRouter)

	safe := general.Group("/", di.authRouter.Middleware())
	safeRouter(safe, di.secretRouter, di.generalRouter)
	clusterRouter(safe, di.clusterRouter)
	bloatRouter(safe, di.bloatRouter)
	certRouter(safe, di.certRouter)
	passwordRouter(safe, di.passwordRouter)
	tagRouter(safe, di.tagRouter)
	instanceRouter(safe, di.instanceRouter)
	queryRouter(safe, di.queryRouter)

	if di.env.Config.CertKeyFilePath != "" && di.env.Config.CertFilePath != "" {
		slog.Info("Ivory serves https connection under address " + di.env.Config.UrlAddress)
		err := engine.RunTLS(di.env.Config.UrlAddress, di.env.Config.CertFilePath, di.env.Config.CertKeyFilePath)
		if err != nil {
			panic(err)
		}
	} else {
		slog.Info("Ivory serves http connection under address " + di.env.Config.UrlAddress)
		err := engine.Run(di.env.Config.UrlAddress)
		if err != nil {
			panic(err)
		}
	}
}

func pong(context *gin.Context) {
	context.JSON(http.StatusOK, gin.H{"message": "pong"})
}

func generalRouter(g *gin.RouterGroup, ra *router.AuthRouter, rg *router.GeneralRouter) {
	g.POST("/login", ra.Login)

	initial := g.Group("/initial")
	initial.POST("/config", rg.SetAppConfig)
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
	instance.GET("/overview", r.GetInstanceOverview)
	instance.GET("/config", r.GetInstanceConfig)
	instance.PATCH("/config", r.PatchInstanceConfig)
	instance.POST("/switchover", r.PostInstanceSwitchover)
	instance.DELETE("/switchover", r.DeleteInstanceSwitchover)
	instance.POST("/reinitialize", r.PostInstanceReinitialize)
	instance.POST("/restart", r.PostInstanceRestart)
	instance.DELETE("/restart", r.DeleteInstanceRestart)
	instance.POST("/reload", r.PostInstanceReload)
	instance.POST("/failover", r.PostInstanceFailover)
	instance.POST("/activate", r.PostInstanceActivate)
	instance.POST("/pause", r.PostInstancePause)
}

func queryRouter(g *gin.RouterGroup, r *router.QueryRouter) {
	query := g.Group("/query")
	query.GET("", r.GetQueryList)
	query.DELETE("/:uuid", r.DeleteQuery)
	query.POST("/run", r.PostRunQuery)
	query.POST("/activity", r.PostAllRunningQueriesByApplicationName)
	query.POST("/databases", r.PostDatabasesQuery)
	query.POST("/schemas", r.PostSchemasQuery)
	query.POST("/tables", r.PostTablesQuery)
	query.POST("/chart", r.PostChartQuery)
	query.POST("/cancel", r.PostCancelQuery)
	query.POST("/terminate", r.PostTerminateQuery)

	queryManual := query.Group("", r.ManualMiddleware())
	queryManual.POST("", r.PostQuery)
	queryManual.PUT("/:uuid", r.PutQuery)

	log := query.Group("/log")
	log.GET("/:uuid", r.GetQueryLog)
	log.DELETE("/:uuid", r.DeleteQueryLog)
}
