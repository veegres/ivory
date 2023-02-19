package router

import (
	"github.com/gin-gonic/gin"
	. "ivory/src/model"
	"ivory/src/service"
	"net/http"
)

type SecretRouter struct {
	secretService *service.SecretService
}

func NewSecretRouter(secretService *service.SecretService) *SecretRouter {
	return &SecretRouter{secretService: secretService}
}

func (r *SecretRouter) GetStatus(context *gin.Context) {
	context.JSON(http.StatusOK, gin.H{"response": r.secretService.Status()})
}

func (r *SecretRouter) SetSecret(context *gin.Context) {
	var body SecretSetRequest
	err := context.ShouldBindJSON(&body)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "body isn't correct"})
		return
	}
	if body.Key == "" {
		context.JSON(http.StatusBadRequest, gin.H{"error": "please provide key"})
		return
	}

	err = r.secretService.Set(body.Key, body.Ref)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, gin.H{"response": "the secret was set"})
}

func (r *SecretRouter) UpdateSecret(context *gin.Context) {
	var secret SecretUpdateRequest
	errParse := context.ShouldBindJSON(&secret)
	if errParse != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": errParse.Error()})
		return
	}
	err := r.secretService.Update(secret.PreviousKey, secret.NewKey)
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
