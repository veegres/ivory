package router

import (
	"github.com/gin-gonic/gin"
	. "ivory/src/model"
	"ivory/src/service"
	"net/http"
)

type SecretRouter struct {
	secretService   *service.SecretService
	passwordService *service.PasswordService
}

func NewSecretRouter(
	secretService *service.SecretService,
	passwordService *service.PasswordService,
) *SecretRouter {
	return &SecretRouter{
		secretService:   secretService,
		passwordService: passwordService,
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
	if body.Key == "" {
		context.JSON(http.StatusBadRequest, gin.H{"error": "please provide key"})
		return
	}

	err := r.secretService.Set(body.Key, body.Ref)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, gin.H{"response": "the secret was set"})
}

func (r *SecretRouter) UpdateSecret(context *gin.Context) {
	var secret SecretUpdateRequest
	errBind := context.ShouldBindJSON(&secret)
	if errBind != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": errBind.Error()})
		return
	}
	prevSha, newSha, err := r.secretService.Update(secret.PreviousKey, secret.NewKey)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err = r.passwordService.ReEncryptPasswords(prevSha, newSha)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, gin.H{"response": "the secret was set"})
}

func (r *SecretRouter) CleanSecret(context *gin.Context) {
	err := r.secretService.Clean()
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, gin.H{"response": "the secret was cleaned"})
}
