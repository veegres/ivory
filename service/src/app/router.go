package app

import (
	"ivory/src/features/auth"
	"ivory/src/features/bloat"
	"ivory/src/features/cert"
	"ivory/src/features/cluster"
	"ivory/src/features/config"
	"ivory/src/features/instance"
	"ivory/src/features/management"
	"ivory/src/features/password"
	"ivory/src/features/query"
	"ivory/src/features/secret"
	"ivory/src/features/tag"
	"log/slog"
	"net/http"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
)

func NewRouter(di *Context) {
	engine := gin.Default()
	engine.UseH2C = true

	// NOTE: Serving ivory static files to web
	if di.env.Config.StaticFilesPath != "" {
		engine.Use(static.Serve(di.env.Config.UrlPath, static.LocalFile(di.env.Config.StaticFilesPath, true)))
		engine.NoRoute(func(context *gin.Context) {
			// NOTE: if files weren't found and NoRoute come here, we need to throw 404 and prevent endless redirect
			if context.Request.URL.Path != di.env.Config.UrlPath {
				context.Redirect(http.StatusMovedPermanently, di.env.Config.UrlPath)
			}
		})
	}

	// NOTE: Setup default sub path for reverse proxies, default "/"
	path := engine.Group(di.env.Config.UrlPath)
	unsafe := path.Group("/api", gin.Recovery(), di.authRouter.SessionMiddleware())
	unsafe.GET("/ping", pong)
	unsafe.GET("/info", di.managementRouter.GetAppInfo)

	initial := unsafe.Group("/", di.secretRouter.EmptyMiddleware())
	initialRouter(initial, di.secretRouter, di.managementRouter)

	general := unsafe.Group("/", di.secretRouter.ExistMiddleware())
	generalRouter(general, di.authRouter, di.configRouter)

	safe := general.Group("/", di.authRouter.AuthMiddleware())
	safe.POST("/logout", di.authRouter.Logout)
	managementRouter(safe, di.secretRouter, di.managementRouter)
	clusterRouter(safe, di.clusterRouter)
	bloatRouter(safe, di.bloatRouter)
	certRouter(safe, di.certRouter)
	passwordRouter(safe, di.passwordRouter)
	tagRouter(safe, di.tagRouter)
	instanceRouter(safe, di.instanceRouter)
	queryRouter(safe, di.queryRouter)

	slog.Info("Ivory address: " + di.env.Config.UrlAddress)
	slog.Info("Ivory WEB path: " + di.env.Config.UrlPath)
	slog.Info("Ivory API path: " + unsafe.BasePath())

	if di.env.Config.TlsEnabled {
		slog.Info("Ivory connection type: HTTPS")
		err := engine.RunTLS(di.env.Config.UrlAddress, di.env.Config.CertFilePath, di.env.Config.CertKeyFilePath)
		if err != nil {
			panic(err)
		}
	} else {
		slog.Info("Ivory connection type: HTTP")
		err := engine.Run(di.env.Config.UrlAddress)
		if err != nil {
			panic(err)
		}
	}
}

func pong(context *gin.Context) {
	context.JSON(http.StatusOK, gin.H{"message": "pong"})
}

func generalRouter(g *gin.RouterGroup, ra *auth.Router, rg *config.Router) {
	g.POST("/basic/login", ra.BasicLogin)
	g.POST("/ldap/login", ra.LdapLogin)
	g.GET("/oidc/login", ra.OidcLogin)
	g.GET("/oidc/callback", ra.OidcCallback)

	initial := g.Group("/initial")
	initial.POST("/config", rg.SetAppConfig)
}

func initialRouter(g *gin.RouterGroup, rs *secret.Router, rg *management.Router) {
	group := g.Group("/initial")
	// TODO check if we can access this endpoint after secret was set
	group.POST("/secret", rs.SetSecret)
	group.DELETE("/erase", rg.Erase)
}

func managementRouter(g *gin.RouterGroup, rs *secret.Router, rm *management.Router) {
	group := g.Group("/management")
	group.GET("/secret", rs.GetSecretStatus)
	group.POST("/secret", rm.ChangeSecret)
	group.DELETE("/erase", rm.Erase)
}

func clusterRouter(g *gin.RouterGroup, r *cluster.Router) {
	group := g.Group("/cluster")
	group.GET("", r.GetClusterList)
	group.GET("/:name", r.GetClusterByName)
	group.PUT("", r.PutClusterByName)
	group.POST("/auto", r.PostClusterAuto)
	group.DELETE("/:name", r.DeleteClusterByName)
}

func bloatRouter(g *gin.RouterGroup, r *bloat.Router) {
	group := g.Group("/cli")
	group.GET("/bloat", r.GetBloatList)
	group.GET("/bloat/:uuid", r.GetBloat)
	group.GET("/bloat/:uuid/logs", r.GetBloatLogs)
	group.GET("/bloat/cluster/:name", r.GetBloatListByCluster)
	group.POST("/bloat/job/start", r.PostJobStart)
	group.POST("/bloat/job/:uuid/stop", r.PostJobStop)
	group.DELETE("/bloat/job/:uuid/delete", r.DeleteJob)
	group.GET("/bloat/job/:uuid/stream", r.GetJobStream)
}

func certRouter(g *gin.RouterGroup, r *cert.Router) {
	group := g.Group("/cert")
	group.GET("", r.GetCertList)
	group.DELETE("/:uuid", r.DeleteCert)
	group.POST("/upload", r.PostUploadCert)
	group.POST("/add", r.PostAddCert)
}

func passwordRouter(g *gin.RouterGroup, r *password.Router) {
	group := g.Group("/password")
	group.GET("", r.GetCredentials)
	group.POST("", r.PostCredential)
	group.PATCH("/:uuid", r.PatchCredential)
	group.DELETE("/:uuid", r.DeleteCredential)
}

func tagRouter(g *gin.RouterGroup, r *tag.Router) {
	group := g.Group("/tag")
	group.GET("", r.GetTagList)
}

func instanceRouter(g *gin.RouterGroup, r *instance.Router) {
	group := g.Group("/instance")
	group.GET("/overview", r.GetInstanceOverview)
	group.GET("/config", r.GetInstanceConfig)
	group.PATCH("/config", r.PatchInstanceConfig)
	group.POST("/switchover", r.PostInstanceSwitchover)
	group.DELETE("/switchover", r.DeleteInstanceSwitchover)
	group.POST("/reinitialize", r.PostInstanceReinitialize)
	group.POST("/restart", r.PostInstanceRestart)
	group.DELETE("/restart", r.DeleteInstanceRestart)
	group.POST("/reload", r.PostInstanceReload)
	group.POST("/failover", r.PostInstanceFailover)
	group.POST("/activate", r.PostInstanceActivate)
	group.POST("/pause", r.PostInstancePause)
}

func queryRouter(g *gin.RouterGroup, r *query.Router) {
	group := g.Group("/query")
	group.GET("", r.GetQueryList)
	group.POST("/run", r.PostRunQuery)
	group.POST("/activity", r.PostActivityQuery)
	group.POST("/databases", r.PostDatabasesQuery)
	group.POST("/schemas", r.PostSchemasQuery)
	group.POST("/tables", r.PostTablesQuery)
	group.POST("/chart", r.PostChartQuery)
	group.POST("/cancel", r.PostCancelQuery)
	group.POST("/terminate", r.PostTerminateQuery)

	manualGroup := group.Group("", r.ManualMiddleware())
	manualGroup.POST("", r.PostQuery)
	manualGroup.PUT("/:uuid", r.PutQuery)
	manualGroup.DELETE("/:uuid", r.DeleteQuery)

	logGroup := group.Group("/log")
	logGroup.GET("/:uuid", r.GetQueryLog)
	logGroup.DELETE("/:uuid", r.DeleteQueryLog)
}
