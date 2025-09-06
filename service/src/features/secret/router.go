package secret

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Router struct {
	secretService *Service
}

func NewRouter(
	secretService *Service,
) *Router {
	return &Router{
		secretService: secretService,
	}
}

func (r *Router) ExistMiddleware() gin.HandlerFunc {
	return func(context *gin.Context) {
		if r.secretService.IsSecretEmpty() {
			context.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "usage restricted when secret key is not specified"})
			return
		}
		context.Next()
	}
}

func (r *Router) EmptyMiddleware() gin.HandlerFunc {
	return func(context *gin.Context) {
		if !r.secretService.IsSecretEmpty() {
			context.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "usage restricted when secret key is specified"})
			return
		}
		context.Next()
	}
}

func (r *Router) GetStatus(context *gin.Context) {
	context.JSON(http.StatusOK, gin.H{"response": r.secretService.Status()})
}

func (r *Router) SetSecret(context *gin.Context) {
	var body SecretSetRequest
	errBind := context.ShouldBindJSON(&body)
	if errBind != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "body isn't correct"})
		return
	}

	err := r.secretService.Set(body.Key, body.Ref)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, gin.H{"response": "the secret was set"})
}
