package permission

import (
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
		authEnabled := context.GetBool("auth")
		authType := context.GetString("authType")
		username := context.GetString("username")
		permissions, err := r.permissionService.GetUserPermissions(authType, username, !authEnabled)
		if err != nil {
			context.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		context.Set("permissions", permissions)
		context.Next()
	}
}

func (r *Router) ValidateMethodMiddleware(permission Permission) gin.HandlerFunc {
	return func(context *gin.Context) {
		if val, ok := context.Get("permissions"); ok {
			permissions := val.(PermissionMap)
			if permissions[permission] != GRANTED {
				context.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": string(permission) + " is not permitted"})
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
	username := context.GetString("username")
	prefix := context.GetString("authType")
	var request PermissionRequest
	if err := context.ShouldBindJSON(&request); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := r.permissionService.RequestUserPermissions(prefix, username, request.Permissions)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"response": "permission request submitted successfully"})
}

func (r *Router) ApproveUserPermission(context *gin.Context) {
	permUsername := context.Param("permUsername")
	if permUsername == "" {
		context.JSON(http.StatusBadRequest, gin.H{"error": "permUsername is required"})
		return
	}

	var request PermissionRequest
	if err := context.ShouldBindJSON(&request); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := r.permissionService.ApproveUserPermissions(permUsername, request.Permissions)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"response": "permission request submitted successfully"})
}

func (r *Router) RejectUserPermission(context *gin.Context) {
	permUsername := context.Param("permUsername")
	if permUsername == "" {
		context.JSON(http.StatusBadRequest, gin.H{"error": "permUsername is required"})
		return
	}

	var request PermissionRequest
	if err := context.ShouldBindJSON(&request); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := r.permissionService.RejectUserPermissions(permUsername, request.Permissions)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"response": "permission request submitted successfully"})
}

func (r *Router) DeleteUserPermissions(context *gin.Context) {
	permUsername := context.Param("permUsername")
	if permUsername == "" {
		context.JSON(http.StatusBadRequest, gin.H{"error": "permUsername is required"})
		return
	}

	err := r.permissionService.DeleteUserPermissions(permUsername)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"response": "permissions deleted successfully"})
}
