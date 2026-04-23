package app

import (
	"ivory/src/features/auth"
	"ivory/src/features/bloat"
	"ivory/src/features/cert"
	"ivory/src/features/cluster"
	"ivory/src/features/config"
	"ivory/src/features/management"
	"ivory/src/features/node"
	"ivory/src/features/permission"
	"ivory/src/features/query"
	"ivory/src/features/secret"
	"ivory/src/features/tag"
	"ivory/src/features/vault"
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
	nodeRouter(safe, di.permissionRouter, di.nodeRouter)
	tagRouter(safe, di.permissionRouter, di.tagRouter)
	bloatRouter(safe, di.permissionRouter, di.bloatRouter)
	certRouter(safe, di.permissionRouter, di.certRouter)
	vaultRouter(safe, di.permissionRouter, di.vaultRouter)
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
	group.POST("/skip", rs.SkipSecret)
	group.POST("/secret", rs.SetSecret)
	group.DELETE("/erase", rg.Erase)
}

func managementRouter(g *gin.RouterGroup, rp *permission.Router, rs *secret.Router, rm *management.Router) {
	group := g.Group("/management")
	group.GET("/secret", rp.ValidateMethodMiddleware(permission.ViewManagementSecret), rs.GetSecretStatus)
	group.POST("/secret", rp.ValidateMethodMiddleware(permission.ManageManagementSecret), rm.ChangeSecret)
	group.DELETE("/erase", rp.ValidateMethodMiddleware(permission.ManageManagementErase), rm.Erase)
	group.DELETE("/free", rp.ValidateMethodMiddleware(permission.ManageManagementFree), rm.Free)
	group.POST("/export", rp.ValidateMethodMiddleware(permission.ManageManagementBackup), rm.Export)
	group.POST("/import", rp.ValidateMethodMiddleware(permission.ManageManagementBackup), rm.Import)
}

func clusterRouter(g *gin.RouterGroup, rp *permission.Router, r *cluster.Router) {
	group := g.Group("/cluster")

	group.GET("", rp.ValidateMethodMiddleware(permission.ViewClusterList), r.GetClusterList)
	group.GET("/:name", rp.ValidateMethodMiddleware(permission.ViewClusterItem), r.GetClusterByName)
	group.PUT("", rp.ValidateMethodMiddleware(permission.ManageClusterUpdate), r.PutClusterByName)
	group.DELETE("/:name", rp.ValidateMethodMiddleware(permission.ManageClusterDelete), r.DeleteClusterByName)
	group.GET("/overview/:name", rp.ValidateMethodMiddleware(permission.ViewClusterOverview), r.GetClusterOverview)
	group.POST("/auto", rp.ValidateMethodMiddleware(permission.ManageClusterCreate), r.PostClusterAutoCreate)
	group.POST("/auto/:name", rp.ValidateMethodMiddleware(permission.ManageClusterUpdate), r.PostClusterAutoFix)
}

func bloatRouter(g *gin.RouterGroup, rp *permission.Router, r *bloat.Router) {
	group := g.Group("/cli")
	group.GET("/bloat", rp.ValidateMethodMiddleware(permission.ViewBloatList), r.GetBloatList)
	group.GET("/bloat/cluster/:name", rp.ValidateMethodMiddleware(permission.ViewBloatList), r.GetBloatListByCluster)
	group.GET("/bloat/:uuid", rp.ValidateMethodMiddleware(permission.ViewBloatItem), r.GetBloat)
	group.GET("/bloat/:uuid/logs", rp.ValidateMethodMiddleware(permission.ViewBloatLogs), r.GetBloatLogs)
	group.GET("/bloat/job/:uuid/stream", rp.ValidateMethodMiddleware(permission.ViewBloatLogs), r.GetJobStream)
	group.POST("/bloat/job/start", rp.ValidateMethodMiddleware(permission.ManageBloatJob), r.PostJobStart)
	group.POST("/bloat/job/:uuid/stop", rp.ValidateMethodMiddleware(permission.ManageBloatJob), r.PostJobStop)
	group.DELETE("/bloat/job/:uuid/delete", rp.ValidateMethodMiddleware(permission.ManageBloatJob), r.DeleteJob)
}

func certRouter(g *gin.RouterGroup, rp *permission.Router, r *cert.Router) {
	group := g.Group("/cert")
	group.GET("", rp.ValidateMethodMiddleware(permission.ViewCertList), r.GetCertList)
	group.POST("/upload", rp.ValidateMethodMiddleware(permission.ManageCertCreate), r.PostUploadCert)
	group.POST("/add", rp.ValidateMethodMiddleware(permission.ManageCertCreate), r.PostAddCert)
	group.DELETE("/:uuid", rp.ValidateMethodMiddleware(permission.ManageCertDelete), r.DeleteCert)
}

func vaultRouter(g *gin.RouterGroup, rp *permission.Router, r *vault.Router) {
	group := g.Group("/vault")
	group.GET("", rp.ValidateMethodMiddleware(permission.ViewVaultList), r.GetVaultList)
	group.POST("", rp.ValidateMethodMiddleware(permission.ManageVaultCreate), r.PostVault)
	group.PATCH("/:uuid", rp.ValidateMethodMiddleware(permission.ManageVaultUpdate), r.PatchVault)
	group.DELETE("/:uuid", rp.ValidateMethodMiddleware(permission.ManageVaultDelete), r.DeleteVault)
}

func permissionRouter(g *gin.RouterGroup, rp *permission.Router, r *permission.Router) {
	group := g.Group("/permission")
	group.POST("/request", r.RequestUserPermission)
	group.GET("/users", rp.ValidateMethodMiddleware(permission.ViewPermissionList), r.GetAllUserPermissions)
	group.POST("/users/:permUsername/approve", rp.ValidateMethodMiddleware(permission.ManagePermissionUpdate), r.ApproveUserPermission)
	group.POST("/users/:permUsername/reject", rp.ValidateMethodMiddleware(permission.ManagePermissionUpdate), r.RejectUserPermission)
	group.DELETE("/users/:permUsername", rp.ValidateMethodMiddleware(permission.ManagePermissionDelete), r.DeleteUserPermissions)
}

func tagRouter(g *gin.RouterGroup, rp *permission.Router, r *tag.Router) {
	group := g.Group("/tag")
	group.GET("", rp.ValidateMethodMiddleware(permission.ViewTagList), r.GetTagList)
}

func nodeRouter(g *gin.RouterGroup, rp *permission.Router, r *node.Router) {
	group := g.Group("/node")

	dbGroup := group.Group("/db")
	dbGroup.GET("/overview", rp.ValidateMethodMiddleware(permission.ViewNodeDbOverview), r.GetNodeOverview)
	dbGroup.GET("/config", rp.ValidateMethodMiddleware(permission.ViewNodeDbConfig), r.GetNodeConfig)
	dbGroup.PATCH("/config", rp.ValidateMethodMiddleware(permission.ManageNodeDbConfigUpdate), r.PatchNodeConfig)
	dbGroup.POST("/switchover", rp.ValidateMethodMiddleware(permission.ManageNodeDbSwitchover), r.PostNodeSwitchover)
	dbGroup.DELETE("/switchover", rp.ValidateMethodMiddleware(permission.ManageNodeDbSwitchover), r.DeleteNodeSwitchover)
	dbGroup.POST("/reinitialize", rp.ValidateMethodMiddleware(permission.ManageNodeDbReinitialize), r.PostNodeReinitialize)
	dbGroup.POST("/restart", rp.ValidateMethodMiddleware(permission.ManageNodeDbRestart), r.PostNodeRestart)
	dbGroup.DELETE("/restart", rp.ValidateMethodMiddleware(permission.ManageNodeDbRestart), r.DeleteNodeRestart)
	dbGroup.POST("/reload", rp.ValidateMethodMiddleware(permission.ManageNodeDbReload), r.PostNodeReload)
	dbGroup.POST("/failover", rp.ValidateMethodMiddleware(permission.ManageNodeDbFailover), r.PostNodeFailover)
	dbGroup.POST("/activate", rp.ValidateMethodMiddleware(permission.ManageNodeDbActivation), r.PostNodeActivate)
	dbGroup.POST("/pause", rp.ValidateMethodMiddleware(permission.ManageNodeDbActivation), r.PostNodePause)

	sshGroup := group.Group("/ssh")
	sshGroup.GET("/metrics", rp.ValidateMethodMiddleware(permission.ViewNodeSshMetrics), r.GetMetrics)

	dockerGroup := sshGroup.Group("/docker")
	dockerGroup.GET("", rp.ValidateMethodMiddleware(permission.ViewNodeSshDocker), r.GetDockerList)
	dockerGroup.GET("/logs", rp.ValidateMethodMiddleware(permission.ViewNodeSshDocker), r.GetDockerLogs)
	dockerGroup.POST("/deploy", rp.ValidateMethodMiddleware(permission.ManageNodeSshDocker), r.PostDockerDeploy)
	dockerGroup.POST("/stop", rp.ValidateMethodMiddleware(permission.ManageNodeSshDocker), r.PostDockerStop)
	dockerGroup.POST("/run", rp.ValidateMethodMiddleware(permission.ManageNodeSshDocker), r.PostDockerRun)
	dockerGroup.POST("/delete", rp.ValidateMethodMiddleware(permission.ManageNodeSshDocker), r.PostDockerDelete)
}

func queryRouter(g *gin.RouterGroup, rp *permission.Router, r *query.Router) {
	group := g.Group("/query")
	group.GET("", rp.ValidateMethodMiddleware(permission.ViewQueryList), r.GetQueryList)
	group.POST("", rp.ValidateMethodMiddleware(permission.ManageQueryCreate), r.PostQuery)
	group.PUT("/:uuid", rp.ValidateMethodMiddleware(permission.ManageQueryUpdate), r.PutQuery)
	group.DELETE("/:uuid", rp.ValidateMethodMiddleware(permission.ManageQueryDelete), r.DeleteQuery)

	executeGroup := group.Group("/execute")
	executeGroup.POST("/console", rp.ValidateMethodMiddleware(permission.ManageQueryExecuteConsole), r.PostExecuteConsoleQuery)
	executeGroup.POST("/template", rp.ValidateMethodMiddleware(permission.ManageQueryExecuteTemplate), r.PostExecuteTemplateQuery)
	executeGroup.POST("/activity", rp.ValidateMethodMiddleware(permission.ViewQueryExecuteInfo), r.PostActivityQuery)
	executeGroup.POST("/databases", rp.ValidateMethodMiddleware(permission.ViewQueryExecuteInfo), r.PostDatabasesQuery)
	executeGroup.POST("/schemas", rp.ValidateMethodMiddleware(permission.ViewQueryExecuteInfo), r.PostSchemasQuery)
	executeGroup.POST("/tables", rp.ValidateMethodMiddleware(permission.ViewQueryExecuteInfo), r.PostTablesQuery)
	executeGroup.POST("/chart", rp.ValidateMethodMiddleware(permission.ViewQueryExecuteChart), r.PostChartQuery)
	executeGroup.POST("/cancel", rp.ValidateMethodMiddleware(permission.ManageQueryExecuteCancel), r.PostCancelQuery)
	executeGroup.POST("/terminate", rp.ValidateMethodMiddleware(permission.ManageQueryExecuteTerminate), r.PostTerminateQuery)

	logGroup := group.Group("/log")
	logGroup.GET("/:uuid", rp.ValidateMethodMiddleware(permission.ViewQueryLogList), r.GetQueryLog)
	logGroup.DELETE("/:uuid", rp.ValidateMethodMiddleware(permission.ManageQueryLogDelete), r.DeleteQueryLog)
}
