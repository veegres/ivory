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
		username := context.GetString("username")
		permissions, err := r.permissionService.GetUserPermissions(username)
		if err != nil {
			context.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		context.Set("permissions", permissions)
		context.Next()
	}
}

func (r *Router) ValidateMethodMiddleware(permission string) gin.HandlerFunc {
	return func(context *gin.Context) {
		if val, ok := context.Get("permissions"); ok {
			permissions := val.(map[string]PermissionStatus)
			if permissions[permission] != GRANTED {
				context.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": permission + " is not permitted"})
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
	username := context.Param("username")
	if username == "" {
		context.JSON(http.StatusBadRequest, gin.H{"error": "username is required"})
		return
	}

	var request struct {
		Permission string `json:"permission"`
	}
	if err := context.ShouldBindJSON(&request); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := r.permissionService.RequestUserPermission(username, request.Permission)
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

	var request struct {
		Permission string `json:"permission"`
	}
	if err := context.ShouldBindJSON(&request); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := r.permissionService.ApproveUserPermission(username, request.Permission)
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

	var request struct {
		Permission string `json:"permission"`
	}
	if err := context.ShouldBindJSON(&request); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := r.permissionService.RejectUserPermission(username, request.Permission)
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
