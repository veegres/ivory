package permission

import (
	"ivory/core/config"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Router struct {
	permissionService *Service
}

func NewRouter(permissionService *Service) *Router {
	return &Router{permissionService: permissionService}
}

func (r *Router) ValidateMiddleware() gin.HandlerFunc {
	return func(context *gin.Context) {
		authEnabled := context.GetBool(config.AuthContextKey.Enabled)
		username := context.GetString(config.AuthContextKey.Username)
		// NOTE: a superuser holds every permission and cannot be talked out of
		// it, so their request is authorised wholesale rather than out of a
		// record somebody else could edit - this feature never learns what a
		// superuser is, the auth middleware that does says so here
		superuser := context.GetBool(config.AuthContextKey.Superuser)
		permissions, err := r.permissionService.GetUserPermissions(username, !authEnabled || superuser, config.FeaturesWithheldWithoutAuth(authEnabled))
		if err != nil {
			context.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		context.Set("permissions", permissions)
		context.Next()
	}
}

func (r *Router) ValidateMethodMiddleware(feature config.Feature) gin.HandlerFunc {
	return func(context *gin.Context) {
		if val, ok := context.Get("permissions"); ok {
			permissions := val.(PermissionMap)
			if permissions[feature] != GRANTED {
				context.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": string(feature) + " is not permitted"})
				return
			}
		} else {
			context.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "permission validation failed"})
			return
		}
		context.Next()
	}
}

func (r *Router) GetAllUserPermissions(context *gin.Context) {
	userPermissions, err := r.permissionService.GetAllUserPermissions()
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, gin.H{"response": userPermissions})
}

func (r *Router) RequestUserPermission(context *gin.Context) {
	username := context.GetString(config.AuthContextKey.Username)
	var request PermissionRequest
	if err := context.ShouldBindJSON(&request); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := r.permissionService.RequestUserPermissions(username, request.Permissions)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"response": "permission request submitted successfully"})
}

func (r *Router) ApproveUserPermission(context *gin.Context) {
	username := context.Param("username")
	if username == "" {
		context.JSON(http.StatusBadRequest, gin.H{"error": "username is required"})
		return
	}

	var request PermissionRequest
	if err := context.ShouldBindJSON(&request); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := r.permissionService.ApproveUserPermissions(username, request.Permissions)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"response": "permission request submitted successfully"})
}

func (r *Router) RejectUserPermission(context *gin.Context) {
	username := context.Param("username")
	if username == "" {
		context.JSON(http.StatusBadRequest, gin.H{"error": "username is required"})
		return
	}

	var request PermissionRequest
	if err := context.ShouldBindJSON(&request); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := r.permissionService.RejectUserPermissions(username, request.Permissions)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"response": "permission request submitted successfully"})
}

func (r *Router) DeleteUserPermissions(context *gin.Context) {
	username := context.Param("username")
	if username == "" {
		context.JSON(http.StatusBadRequest, gin.H{"error": "username is required"})
		return
	}

	err := r.permissionService.DeleteUserPermissions(username)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"response": "permissions deleted successfully"})
}
