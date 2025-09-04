package secret

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type SecretRouter struct {
	secretService *SecretService
}

func NewSecretRouter(
	secretService *SecretService,
) *SecretRouter {
	return &SecretRouter{
		secretService: secretService,
	}
}

func (r *SecretRouter) ExistMiddleware() gin.HandlerFunc {
	return func(context *gin.Context) {
		if r.secretService.IsSecretEmpty() {
			context.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "usage restricted when secret key is not specified"})
			return
		}
		context.Next()
	}
}

func (r *SecretRouter) EmptyMiddleware() gin.HandlerFunc {
	return func(context *gin.Context) {
		if !r.secretService.IsSecretEmpty() {
			context.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "usage restricted when secret key is specified"})
			return
		}
		context.Next()
	}
}

func (r *SecretRouter) GetStatus(context *gin.Context) {
	context.JSON(http.StatusOK, gin.H{"response": r.secretService.Status()})
}

func (r *SecretRouter) SetSecret(context *gin.Context) {
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
