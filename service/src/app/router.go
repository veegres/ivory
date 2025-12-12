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
	"ivory/src/features/permission"
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
	unsafe.POST("/logout", di.authRouter.Logout)

	initial := unsafe.Group("/", di.secretRouter.EmptyMiddleware())
	initialRouter(initial, di.secretRouter, di.managementRouter)

	general := unsafe.Group("/", di.secretRouter.ExistMiddleware())
	generalRouter(general, di.authRouter, di.configRouter)

	safe := general.Group("/", di.configRouter.InitialiseMiddleware(), di.authRouter.ValidateMiddleware(), di.permissionRouter.ValidateMiddleware())
	managementRouter(safe, di.permissionRouter, di.secretRouter, di.managementRouter)
	clusterRouter(safe, di.permissionRouter, di.clusterRouter)
	instanceRouter(safe, di.permissionRouter, di.instanceRouter)
	tagRouter(safe, di.permissionRouter, di.tagRouter)
	bloatRouter(safe, di.permissionRouter, di.bloatRouter)
	certRouter(safe, di.permissionRouter, di.certRouter)
	passwordRouter(safe, di.permissionRouter, di.passwordRouter)
	permissionRouter(safe, di.permissionRouter, di.permissionRouter)
	queryRouter(safe, di.permissionRouter, di.queryRouter)

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
	g.POST("/config", rg.SetAppConfig)

	g.POST("/basic/login", ra.BasicLogin)
	g.POST("/ldap/login", ra.LdapLogin)
	g.GET("/oidc/login", ra.OidcLogin)
	g.GET("/oidc/callback", ra.OidcCallback)
	g.POST("/oidc/connect", ra.OidcConnect)
	g.POST("/ldap/connect", ra.LdapConnect)
}

func initialRouter(g *gin.RouterGroup, rs *secret.Router, rg *management.Router) {
	group := g.Group("/initial")
	group.POST("/secret", rs.SetSecret)
	group.DELETE("/erase", rg.Erase)
}

func managementRouter(g *gin.RouterGroup, rp *permission.Router, rs *secret.Router, rm *management.Router) {
	group := g.Group("/management")
	group.GET("/secret", rp.ValidateMethodMiddleware("view.management.secret"), rs.GetSecretStatus)
	group.POST("/secret", rp.ValidateMethodMiddleware("manage.management.secret"), rm.ChangeSecret)
	group.DELETE("/erase", rp.ValidateMethodMiddleware("manage.management.erase"), rm.Erase)
}

func clusterRouter(g *gin.RouterGroup, rp *permission.Router, r *cluster.Router) {
	group := g.Group("/cluster")
	group.GET("", rp.ValidateMethodMiddleware("view.cluster.list"), r.GetClusterList)
	group.GET("/:name", rp.ValidateMethodMiddleware("view.cluster.item"), r.GetClusterByName)
	group.PUT("", rp.ValidateMethodMiddleware("manage.cluster.update"), r.PutClusterByName)
	group.POST("/auto", rp.ValidateMethodMiddleware("manage.cluster.create"), r.PostClusterAuto)
	group.DELETE("/:name", rp.ValidateMethodMiddleware("manage.cluster.delete"), r.DeleteClusterByName)
}

func bloatRouter(g *gin.RouterGroup, rp *permission.Router, r *bloat.Router) {
	group := g.Group("/cli")
	group.GET("/bloat", rp.ValidateMethodMiddleware("view.bloat.list"), r.GetBloatList)
	group.GET("/bloat/cluster/:name", rp.ValidateMethodMiddleware("view.bloat.list"), r.GetBloatListByCluster)
	group.GET("/bloat/:uuid", rp.ValidateMethodMiddleware("view.bloat.item"), r.GetBloat)
	group.GET("/bloat/:uuid/logs", rp.ValidateMethodMiddleware("view.bloat.logs"), r.GetBloatLogs)
	group.POST("/bloat/job/start", rp.ValidateMethodMiddleware("manage.bloat.job"), r.PostJobStart)
	group.POST("/bloat/job/:uuid/stop", rp.ValidateMethodMiddleware("manage.bloat.job"), r.PostJobStop)
	group.DELETE("/bloat/job/:uuid/delete", rp.ValidateMethodMiddleware("manage.bloat.job"), r.DeleteJob)
	group.GET("/bloat/job/:uuid/stream", rp.ValidateMethodMiddleware("manage.bloat.job"), r.GetJobStream)
}

func certRouter(g *gin.RouterGroup, rp *permission.Router, r *cert.Router) {
	group := g.Group("/cert")
	group.GET("", rp.ValidateMethodMiddleware("view.cert.list"), r.GetCertList)
	group.POST("/upload", rp.ValidateMethodMiddleware("manage.cert.create"), r.PostUploadCert)
	group.POST("/add", rp.ValidateMethodMiddleware("manage.cert.create"), r.PostAddCert)
	group.DELETE("/:uuid", rp.ValidateMethodMiddleware("manage.cert.delete"), r.DeleteCert)
}

func passwordRouter(g *gin.RouterGroup, rp *permission.Router, r *password.Router) {
	group := g.Group("/password")
	group.GET("", rp.ValidateMethodMiddleware("view.password.list"), r.GetPasswordList)
	group.POST("", rp.ValidateMethodMiddleware("manage.password.create"), r.PostPassword)
	group.PATCH("/:uuid", rp.ValidateMethodMiddleware("manage.password.update"), r.PatchPassword)
	group.DELETE("/:uuid", rp.ValidateMethodMiddleware("manage.password.delete"), r.DeletePassword)
}

func permissionRouter(g *gin.RouterGroup, rp *permission.Router, r *permission.Router) {
	group := g.Group("/permission")
	group.POST("/request", r.RequestUserPermission)
	group.GET("/users", rp.ValidateMethodMiddleware("view.permission.list"), r.GetAllUserPermissions)
	group.POST("/users/:username/approve", rp.ValidateMethodMiddleware("manage.permission.update"), r.ApproveUserPermission)
	group.POST("/users/:username/reject", rp.ValidateMethodMiddleware("manage.permission.update"), r.RejectUserPermission)
	group.DELETE("/users/:username", rp.ValidateMethodMiddleware("manage.permission.delete"), r.DeleteUserPermissions)
}

func tagRouter(g *gin.RouterGroup, rp *permission.Router, r *tag.Router) {
	group := g.Group("/tag")
	group.GET("", rp.ValidateMethodMiddleware("view.tag.list"), r.GetTagList)
}

func instanceRouter(g *gin.RouterGroup, rp *permission.Router, r *instance.Router) {
	group := g.Group("/instance")
	group.GET("/overview", rp.ValidateMethodMiddleware("view.instance.overview"), r.GetInstanceOverview)
	group.GET("/config", rp.ValidateMethodMiddleware("view.instance.config"), r.GetInstanceConfig)
	group.PATCH("/config", rp.ValidateMethodMiddleware("manage.instance.config.update"), r.PatchInstanceConfig)
	group.POST("/switchover", rp.ValidateMethodMiddleware("manage.instance.switchover"), r.PostInstanceSwitchover)
	group.DELETE("/switchover", rp.ValidateMethodMiddleware("manage.instance.switchover"), r.DeleteInstanceSwitchover)
	group.POST("/reinitialize", rp.ValidateMethodMiddleware("manage.instance.reinitialize"), r.PostInstanceReinitialize)
	group.POST("/restart", rp.ValidateMethodMiddleware("manage.instance.restart"), r.PostInstanceRestart)
	group.DELETE("/restart", rp.ValidateMethodMiddleware("manage.instance.restart"), r.DeleteInstanceRestart)
	group.POST("/reload", rp.ValidateMethodMiddleware("manage.instance.reload"), r.PostInstanceReload)
	group.POST("/failover", rp.ValidateMethodMiddleware("manage.instance.failover"), r.PostInstanceFailover)
	group.POST("/activate", rp.ValidateMethodMiddleware("manage.instance.activation"), r.PostInstanceActivate)
	group.POST("/pause", rp.ValidateMethodMiddleware("manage.instance.activation"), r.PostInstancePause)
}

func queryRouter(g *gin.RouterGroup, rp *permission.Router, r *query.Router) {
	group := g.Group("/query")
	group.GET("", rp.ValidateMethodMiddleware("view.query.list"), r.GetQueryList)
	group.POST("", rp.ValidateMethodMiddleware("manage.query.create"), r.PostQuery)
	group.PUT("/:uuid", rp.ValidateMethodMiddleware("manage.query.update"), r.PutQuery)
	group.DELETE("/:uuid", rp.ValidateMethodMiddleware("manage.query.delete"), r.DeleteQuery)

	executeGroup := group.Group("/execute")
	executeGroup.POST("/console", rp.ValidateMethodMiddleware("manage.query.execute.console"), r.PostExecuteConsoleQuery)
	executeGroup.POST("/template", rp.ValidateMethodMiddleware("manage.query.execute.template"), r.PostExecuteTempateQuery)
	executeGroup.POST("/activity", rp.ValidateMethodMiddleware("view.query.execute.info"), r.PostActivityQuery)
	executeGroup.POST("/databases", rp.ValidateMethodMiddleware("view.query.execute.info"), r.PostDatabasesQuery)
	executeGroup.POST("/schemas", rp.ValidateMethodMiddleware("view.query.execute.info"), r.PostSchemasQuery)
	executeGroup.POST("/tables", rp.ValidateMethodMiddleware("view.query.execute.info"), r.PostTablesQuery)
	executeGroup.POST("/chart", rp.ValidateMethodMiddleware("view.query.execute.chart"), r.PostChartQuery)
	executeGroup.POST("/cancel", rp.ValidateMethodMiddleware("manage.query.execute.cancel"), r.PostCancelQuery)
	executeGroup.POST("/terminate", rp.ValidateMethodMiddleware("manage.query.execute.terminate"), r.PostTerminateQuery)

	logGroup := group.Group("/log")
	logGroup.GET("/:uuid", rp.ValidateMethodMiddleware("view.query.log.list"), r.GetQueryLog)
	logGroup.DELETE("/:uuid", rp.ValidateMethodMiddleware("manage.query.log.delete"), r.DeleteQueryLog)
}
