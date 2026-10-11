package secret

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type HttpRouter struct {
	secretService *Service
}

func NewHttpRouter(
	secretService *Service,
) *HttpRouter {
	return &HttpRouter{
		secretService: secretService,
	}
}

func (r *HttpRouter) ExistMiddleware() gin.HandlerFunc {
	return func(context *gin.Context) {
		if r.secretService.IsSecretEmpty() {
			context.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "usage restricted when the secret key is not specified"})
			return
		}
		context.Next()
	}
}

func (r *HttpRouter) EmptyMiddleware() gin.HandlerFunc {
	return func(context *gin.Context) {
		if !r.secretService.IsSecretEmpty() {
			context.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "usage restricted when the secret key is specified"})
			return
		}
		context.Next()
	}
}

func (r *HttpRouter) GetSecretStatus(context *gin.Context) {
	context.JSON(http.StatusOK, gin.H{"response": r.secretService.Status()})
}

func (r *HttpRouter) SkipSecret(context *gin.Context) {
	err := r.secretService.SetDefault()
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, gin.H{"response": "the secret was set"})
}

func (r *HttpRouter) SetSecret(context *gin.Context) {
	var body RouterSetRequest
	errBind := context.ShouldBindJSON(&body)
	if errBind != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "body isn't correct"})
		return
	}

	err := r.secretService.Set(body.Key)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, gin.H{"response": "the secret was set"})
}
