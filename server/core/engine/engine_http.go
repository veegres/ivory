package engine

import (
	"ivory/core"
	coreConfig "ivory/core/config"
	"ivory/core/service/cert"
	"ivory/core/service/secret"
	"ivory/core/service/vault"
	"ivory/features"
	"ivory/features/auth"
	"ivory/features/cluster"
	"ivory/features/config"
	"ivory/features/deployment"
	"ivory/features/management"
	"ivory/features/node"
	"ivory/features/permission"
	"ivory/features/query"
	"ivory/features/tag"
	"ivory/features/user"
	"ivory/tools"
	"log/slog"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
)

func NewHttpServer(env *coreConfig.Environment, cc *core.Router, fc *features.Router, tc *tools.Router) {
	engine := gin.Default()
	engine.UseH2C = true
	urlPath := env.Config.UrlPath
	apiPath := "/api"
	staticFilesPath := env.Config.StaticFilesPath

	// NOTE: Serving ivory static file to web
	if staticFilesPath != "" {
		spaServe(engine, urlPath, apiPath, staticFilesPath)
	}

	// NOTE: Setup default sub path for reverse proxies, default "/"
	root := engine.Group(env.Config.UrlPath)
	unsafe := root.Group(apiPath, gin.Recovery(), fc.Auth.SessionMiddleware())
	unsafe.GET("/ping", pong)
	unsafe.GET("/info", fc.Auth.AuthContextMiddleware(), fc.Management.GetAppInfo)
	unsafe.POST("/logout", fc.Auth.Logout)

	initial := unsafe.Group("/", cc.Secret.EmptyMiddleware())
	initialRouter(initial, cc.Secret)

	general := unsafe.Group("/", cc.Secret.ExistMiddleware())
	generalRouter(general, fc.Auth, fc.Config)
	userPublicRouter(general, fc.User)

	safe := general.Group("/", fc.Config.InitialiseMiddleware(), fc.Auth.ValidateMiddleware(), fc.Permission.ValidateMiddleware())
	toolsRouter(safe, fc.Permission, tc)

	managementRouter(safe, fc.Permission, cc.Secret, fc.Management)
	vaultRouter(safe, fc.Permission, cc.Vault)
	certRouter(safe, fc.Permission, cc.Cert)

	clusterRouter(safe, fc.Permission, fc.Cluster)
	nodeRouter(safe, fc.Permission, fc.Node)
	tagRouter(safe, fc.Permission, fc.Tag)
	permissionRouter(safe, fc.Permission, fc.Permission)
	userRouter(safe, fc.Permission, fc.User)
	queryRouter(safe, fc.Permission, fc.Query)
	deploymentRouter(safe, fc.Permission, fc.Deployment)

	slog.Info("Ivory address: " + env.Config.UrlAddress)
	slog.Info("Ivory WEB path: " + urlPath)
	slog.Info("Ivory API path: " + unsafe.BasePath())

	if env.Config.TlsEnabled {
		slog.Info("Ivory WEB protocol: HTTPS")
		err := engine.RunTLS(env.Config.UrlAddress, env.Config.CertFilePath, env.Config.CertKeyFilePath)
		if err != nil {
			panic(err)
		}
	} else {
		slog.Info("Ivory WEB protocol: HTTP")
		err := engine.Run(env.Config.UrlAddress)
		if err != nil {
			panic(err)
		}
	}
}

func spaServe(engine *gin.Engine, urlPath, apiPath, staticFilesPath string) {
	if urlPath[0] != '/' || apiPath[0] != '/' {
		panic("path should contain `/` at the beginning")
	}

	// Serve static files under urlPath.
	// For urlPath = "/ivory":
	//   /ivory/assets/app.js -> staticFilesPath/assets/app.js
	//   /ivory/favicon.ico   -> staticFilesPath/favicon.ico
	engine.Use(static.Serve(urlPath, static.LocalFile(staticFilesPath, true)))

	// API is served under the SPA base path.
	// urlPath = "/":
	//   /api/...
	// urlPath = "/ivory":
	//   /ivory/api/...
	apiPathFull := urlPath + "/api"
	if urlPath == "/" {
		apiPathFull = "/api"
	}

	// Handle requests that weren't matched by an API route or static file.
	engine.NoRoute(func(c *gin.Context) {
		path := c.Request.URL.Path

		// API is part of the application, but API errors must not
		// be handled by the React SPA.
		if path == apiPathFull || strings.HasPrefix(path, apiPathFull+"/") {
			c.Status(http.StatusNotFound)
			return
		}

		// Only URLs belonging to the SPA should fall back to index.html.
		// urlPath = "/":
		//   /user/123       -> index.html
		//   /clusters       -> index.html
		// urlPath = "/ivory":
		//   /user/123       -> 404
		//   /clusters       -> 404
		//   /ivory/user/123 -> index.html
		//   /ivory/clusters -> index.html
		if urlPath != "/" && path != urlPath && !strings.HasPrefix(path, urlPath+"/") {
			c.Status(http.StatusNotFound)
			return
		}

		// Return index.html without redirecting.
		// React Router receives the original URL.
		c.File(filepath.Join(staticFilesPath, "index.html"))
	})
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

// initialRouter is what a restarted Ivory offers before its secret is back.
// There is deliberately no erase here: wiping everything is a thing only
// somebody signed in may ask for, and whoever has lost the secret word
// reinstalls instead.
func initialRouter(g *gin.RouterGroup, rs *secret.Router) {
	group := g.Group("/initial")
	group.POST("/skip", rs.SkipSecret)
	group.POST("/secret", rs.SetSecret)
}

func toolsRouter(g *gin.RouterGroup, rp *permission.Router, r *tools.Router) {
	group := g.Group("/tool")
	group.GET("/bloat", rp.ValidateMethodMiddleware(coreConfig.ViewToolPgCompactTableList), r.PgCompactTable.GetBloatList)
	group.GET("/bloat/cluster/:name", rp.ValidateMethodMiddleware(coreConfig.ViewToolPgCompactTableList), r.PgCompactTable.GetBloatListByCluster)
	group.GET("/bloat/:uuid", rp.ValidateMethodMiddleware(coreConfig.ViewToolPgCompactTableItem), r.PgCompactTable.GetBloat)
	group.GET("/bloat/:uuid/logs", rp.ValidateMethodMiddleware(coreConfig.ViewToolPgCompactTableLogs), r.PgCompactTable.GetBloatLogs)
	group.GET("/bloat/job/:uuid/stream", rp.ValidateMethodMiddleware(coreConfig.ViewToolPgCompactTableLogs), r.PgCompactTable.GetJobStream)
	group.POST("/bloat/job/start", rp.ValidateMethodMiddleware(coreConfig.ManageToolPgCompactTableJob), r.PgCompactTable.PostJobStart)
	group.POST("/bloat/job/:uuid/stop", rp.ValidateMethodMiddleware(coreConfig.ManageToolPgCompactTableJob), r.PgCompactTable.PostJobStop)
	group.DELETE("/bloat/job/:uuid/delete", rp.ValidateMethodMiddleware(coreConfig.ManageToolPgCompactTableJob), r.PgCompactTable.DeleteJob)
}

func managementRouter(g *gin.RouterGroup, rp *permission.Router, rs *secret.Router, rm *management.Router) {
	group := g.Group("/management")
	group.GET("/secret", rp.ValidateMethodMiddleware(coreConfig.ViewManagementSecret), rs.GetSecretStatus)
	group.POST("/secret", rp.ValidateMethodMiddleware(coreConfig.ManageManagementSecret), rm.ChangeSecret)
	group.DELETE("/erase", rp.ValidateMethodMiddleware(coreConfig.ManageManagementErase), rm.Erase)
	group.DELETE("/free", rp.ValidateMethodMiddleware(coreConfig.ManageManagementFree), rm.Free)
	group.POST("/export", rp.ValidateMethodMiddleware(coreConfig.ManageManagementBackup), rm.Export)
	group.POST("/import", rp.ValidateMethodMiddleware(coreConfig.ManageManagementBackup), rm.Import)
}

func clusterRouter(g *gin.RouterGroup, rp *permission.Router, r *cluster.Router) {
	group := g.Group("/cluster")

	group.GET("", rp.ValidateMethodMiddleware(coreConfig.ViewClusterList), r.GetClusterList)
	group.GET("/:name", rp.ValidateMethodMiddleware(coreConfig.ViewClusterItem), r.GetClusterByName)
	group.POST("", rp.ValidateMethodMiddleware(coreConfig.ManageClusterCreate), r.PostCluster)
	group.PUT("", rp.ValidateMethodMiddleware(coreConfig.ManageClusterUpdate), r.PutClusterByName)
	group.DELETE("/:name", rp.ValidateMethodMiddleware(coreConfig.ManageClusterDelete), r.DeleteClusterByName)
	group.GET("/overview/:name", rp.ValidateMethodMiddleware(coreConfig.ViewClusterOverview), r.GetClusterOverview)
	group.POST("/fix/:name", rp.ValidateMethodMiddleware(coreConfig.ManageClusterUpdate), r.PostClusterFix)
	group.POST("/detect", rp.ValidateMethodMiddleware(coreConfig.ManageClusterCreate), r.PostClusterDetect)
	group.POST("/deploy", rp.ValidateMethodMiddleware(coreConfig.ManageClusterCreate), r.PostClusterDeploy)
}

func certRouter(g *gin.RouterGroup, rp *permission.Router, r *cert.Router) {
	group := g.Group("/cert")
	group.GET("", rp.ValidateMethodMiddleware(coreConfig.ViewCertList), r.GetCertList)
	group.POST("/upload", rp.ValidateMethodMiddleware(coreConfig.ManageCertCreate), r.PostUploadCert)
	group.POST("/add", rp.ValidateMethodMiddleware(coreConfig.ManageCertCreate), r.PostAddCert)
	group.DELETE("/:uuid", rp.ValidateMethodMiddleware(coreConfig.ManageCertDelete), r.DeleteCert)
}

func vaultRouter(g *gin.RouterGroup, rp *permission.Router, r *vault.Router) {
	group := g.Group("/vault")
	group.GET("", rp.ValidateMethodMiddleware(coreConfig.ViewVaultList), r.GetVaultList)
	group.POST("", rp.ValidateMethodMiddleware(coreConfig.ManageVaultCreate), r.PostVault)
	group.PATCH("/:uuid", rp.ValidateMethodMiddleware(coreConfig.ManageVaultUpdate), r.PatchVault)
	group.DELETE("/:uuid", rp.ValidateMethodMiddleware(coreConfig.ManageVaultDelete), r.DeleteVault)
}

func permissionRouter(g *gin.RouterGroup, rp *permission.Router, r *permission.Router) {
	group := g.Group("/permission")
	group.POST("/request", r.RequestUserPermission)
	group.GET("/users", rp.ValidateMethodMiddleware(coreConfig.ViewPermissionList), r.GetAllUserPermissions)
	group.POST("/users/:username/approve", rp.ValidateMethodMiddleware(coreConfig.ManagePermissionUpdate), r.ApproveUserPermission)
	group.POST("/users/:username/reject", rp.ValidateMethodMiddleware(coreConfig.ManagePermissionUpdate), r.RejectUserPermission)
	group.DELETE("/users/:username", rp.ValidateMethodMiddleware(coreConfig.ManagePermissionDelete), r.DeleteUserPermissions)
}

// userPublicRouter carries the two endpoints the registration page needs. They
// are public because the registration is the authorization - the person opening
// it has no session and no password yet - and there is no third one: a
// registration is only ever handed out by somebody signed in, or by the initial
// setup, which registers its superuser with a typed password.
func userPublicRouter(g *gin.RouterGroup, r *user.Router) {
	group := g.Group("/user/registration")
	group.POST("/verify", r.PostUserRegistrationVerify)
	group.POST("/password", r.PostUserRegistrationPassword)
}

func userRouter(g *gin.RouterGroup, rp *permission.Router, r *user.Router) {
	group := g.Group("/user")
	group.GET("", rp.ValidateMethodMiddleware(coreConfig.ViewUserList), r.GetUserList)
	group.POST("", rp.ValidateMethodMiddleware(coreConfig.ManageUserCreate), r.PostUser)
	group.PUT("/:username", rp.ValidateMethodMiddleware(coreConfig.ManageUserUpdate), r.PutUser)
	group.DELETE("/:username", rp.ValidateMethodMiddleware(coreConfig.ManageUserDelete), r.DeleteUser)
	// NOTE: changing your own password needs no permission, exactly as requesting
	// permissions does not - it is the account you are already signed in as
	group.POST("/password", r.PostUserPassword)

	// NOTE: resetting somebody else's password takes their account over, which is
	// a different power from creating a new name, so it holds a permission of its
	// own - what it hands out is a registration like any other
	group.POST("/:username/password/reset", rp.ValidateMethodMiddleware(coreConfig.ManageUserPasswordReset), r.PostUserPasswordReset)
	group.DELETE("/:username/password/reset", rp.ValidateMethodMiddleware(coreConfig.ManageUserPasswordReset), r.DeleteUserPasswordReset)
}

func tagRouter(g *gin.RouterGroup, rp *permission.Router, r *tag.Router) {
	group := g.Group("/tag")
	group.GET("", rp.ValidateMethodMiddleware(coreConfig.ViewTagList), r.GetTagList)
}

func nodeRouter(g *gin.RouterGroup, rp *permission.Router, r *node.Router) {
	group := g.Group("/node")

	keeperGroup := group.Group("/keeper")
	keeperGroup.GET("/overview", rp.ValidateMethodMiddleware(coreConfig.ViewNodeKeeperOverview), r.GetNodeOverview)
	keeperGroup.GET("/config", rp.ValidateMethodMiddleware(coreConfig.ViewNodeKeeperConfig), r.GetNodeConfig)
	keeperGroup.PATCH("/config", rp.ValidateMethodMiddleware(coreConfig.ManageNodeKeeperConfigUpdate), r.PatchNodeConfig)
	keeperGroup.POST("/switchover", rp.ValidateMethodMiddleware(coreConfig.ManageNodeKeeperSwitchover), r.PostNodeSwitchover)
	keeperGroup.DELETE("/switchover", rp.ValidateMethodMiddleware(coreConfig.ManageNodeKeeperSwitchover), r.DeleteNodeSwitchover)
	keeperGroup.POST("/reinitialize", rp.ValidateMethodMiddleware(coreConfig.ManageNodeKeeperReinitialize), r.PostNodeReinitialize)
	keeperGroup.POST("/restart", rp.ValidateMethodMiddleware(coreConfig.ManageNodeKeeperRestart), r.PostNodeRestart)
	keeperGroup.DELETE("/restart", rp.ValidateMethodMiddleware(coreConfig.ManageNodeKeeperRestart), r.DeleteNodeRestart)
	keeperGroup.POST("/reload", rp.ValidateMethodMiddleware(coreConfig.ManageNodeKeeperReload), r.PostNodeReload)
	keeperGroup.POST("/failover", rp.ValidateMethodMiddleware(coreConfig.ManageNodeKeeperFailover), r.PostNodeFailover)
	keeperGroup.POST("/activate", rp.ValidateMethodMiddleware(coreConfig.ManageNodeKeeperActivation), r.PostNodeActivate)
	keeperGroup.POST("/pause", rp.ValidateMethodMiddleware(coreConfig.ManageNodeKeeperActivation), r.PostNodePause)

	platformGroup := group.Group("/platform")

	systemGroup := platformGroup.Group("/system")
	systemGroup.GET("/metrics", rp.ValidateMethodMiddleware(coreConfig.ViewNodeSystem), r.GetPlatformSystemMetrics)
	systemGroup.GET("/logs", rp.ValidateMethodMiddleware(coreConfig.ViewNodeSystem), r.StreamPlatformSystemLogs)
	systemGroup.GET("/processes", rp.ValidateMethodMiddleware(coreConfig.ViewNodeSystem), r.GetPlatformSystemProcesses)
	systemGroup.GET("/info", rp.ValidateMethodMiddleware(coreConfig.ViewNodeSystem), r.GetPlatformSystemInfo)
	systemGroup.POST("/copy-id", rp.ValidateMethodMiddleware(coreConfig.ManageNodeSystem), r.PostPlatformSystemCopyId)

	containerGroup := platformGroup.Group("/container")
	containerGroup.GET("", rp.ValidateMethodMiddleware(coreConfig.ViewNodePlatformContainer), r.GetPlatformContainerList)
	containerGroup.GET("/logs", rp.ValidateMethodMiddleware(coreConfig.ViewNodePlatformContainer), r.StreamPlatformContainerLogs)
	containerGroup.GET("/metrics", rp.ValidateMethodMiddleware(coreConfig.ViewNodePlatformContainer), r.GetPlatformContainerMetrics)
	containerGroup.POST("/start", rp.ValidateMethodMiddleware(coreConfig.ManageNodePlatformContainer), r.PostPlatformContainerStart)
	containerGroup.POST("/stop", rp.ValidateMethodMiddleware(coreConfig.ManageNodePlatformContainer), r.PostPlatformContainerStop)
	containerGroup.POST("/restart", rp.ValidateMethodMiddleware(coreConfig.ManageNodePlatformContainer), r.PostPlatformContainerRestart)
	containerGroup.POST("/up", rp.ValidateMethodMiddleware(coreConfig.ManageNodePlatformContainer), r.PostPlatformContainerUp)
	containerGroup.POST("/down", rp.ValidateMethodMiddleware(coreConfig.ManageNodePlatformContainer), r.PostPlatformContainerDown)

	containerKeeperGroup := containerGroup.Group("/keeper")
	containerKeeperGroup.POST("/deploy", rp.ValidateMethodMiddleware(coreConfig.ManageNodePlatformContainer), r.PostNodeKeeperDeploy)
}

func deploymentRouter(g *gin.RouterGroup, rp *permission.Router, r *deployment.Router) {
	group := g.Group("/deployment/template")
	group.GET("", rp.ValidateMethodMiddleware(coreConfig.ViewDeploymentTemplateList), r.GetDeploymentTemplateList)
	group.POST("", rp.ValidateMethodMiddleware(coreConfig.ManageDeploymentTemplateCreate), r.PostDeploymentTemplate)
	group.PUT("/:uuid", rp.ValidateMethodMiddleware(coreConfig.ManageDeploymentTemplateUpdate), r.PutDeploymentTemplate)
	group.DELETE("/:uuid", rp.ValidateMethodMiddleware(coreConfig.ManageDeploymentTemplateDelete), r.DeleteDeploymentTemplate)
}

func queryRouter(g *gin.RouterGroup, rp *permission.Router, r *query.Router) {
	group := g.Group("/query")
	group.GET("", rp.ValidateMethodMiddleware(coreConfig.ViewQueryCrudList), r.GetQueryList)
	group.POST("", rp.ValidateMethodMiddleware(coreConfig.ManageQueryCrudCreate), r.PostQuery)
	group.PUT("/:uuid", rp.ValidateMethodMiddleware(coreConfig.ManageQueryCrudUpdate), r.PutQuery)
	group.DELETE("/:uuid", rp.ValidateMethodMiddleware(coreConfig.ManageQueryCrudDelete), r.DeleteQuery)

	executeGroup := group.Group("/execute")
	executeGroup.POST("/console", rp.ValidateMethodMiddleware(coreConfig.ManageQueryDbConsole), r.PostExecuteConsoleQuery)
	executeGroup.POST("/template", rp.ValidateMethodMiddleware(coreConfig.ManageQueryDbTemplate), r.PostExecuteTemplateQuery)
	executeGroup.POST("/activity", rp.ValidateMethodMiddleware(coreConfig.ViewQueryDbInfo), r.PostActivityQuery)
	executeGroup.POST("/databases", rp.ValidateMethodMiddleware(coreConfig.ViewQueryDbInfo), r.PostDatabasesQuery)
	executeGroup.POST("/schemas", rp.ValidateMethodMiddleware(coreConfig.ViewQueryDbInfo), r.PostSchemasQuery)
	executeGroup.POST("/tables", rp.ValidateMethodMiddleware(coreConfig.ViewQueryDbInfo), r.PostTablesQuery)
	executeGroup.POST("/chart", rp.ValidateMethodMiddleware(coreConfig.ViewQueryDbChart), r.PostChartQuery)
	executeGroup.POST("/cancel", rp.ValidateMethodMiddleware(coreConfig.ManageQueryDbCancel), r.PostCancelQuery)
	executeGroup.POST("/terminate", rp.ValidateMethodMiddleware(coreConfig.ManageQueryDbTerminate), r.PostTerminateQuery)

	logGroup := group.Group("/log")
	logGroup.GET("/:uuid", rp.ValidateMethodMiddleware(coreConfig.ViewQueryLogList), r.GetQueryLog)
	logGroup.DELETE("/:uuid", rp.ValidateMethodMiddleware(coreConfig.ManageQueryLogDelete), r.DeleteQueryLog)
}
