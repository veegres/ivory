package router

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	. "ivory/model"
	"ivory/persistence"
	"ivory/service"
	"net/http"
)

func (r routes) CredentialGroup(group *gin.RouterGroup) {
	secret := group.Group("/secret")
	secret.GET("", getExistence)
	secret.POST("/set", setSecret)
	secret.POST("/update", updateSecret)
	secret.POST("/clean", cleanSecret)

	credential := group.Group("/credential")
	credential.GET("", getCredentials)
	credential.POST("", postCredential)
	credential.DELETE("/:uuid", deleteCredential)
}

func getExistence(context *gin.Context) {
	context.JSON(http.StatusOK, gin.H{"response": SecretStatus{
		Key: !service.Secret.IsSecretEmpty(),
		Ref: !service.Secret.IsRefEmpty(),
	}})
}

func setSecret(context *gin.Context) {
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

	err = service.Secret.Set(body.Key, body.Ref)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, gin.H{"response": "the secret was set"})
}

func updateSecret(context *gin.Context) {
	var secret SecretUpdateRequest
	_ = context.ShouldBindJSON(&secret)
	err := service.Secret.Update(secret.PreviousKey, secret.NewKey)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, gin.H{"response": "the secret was set"})
}

func cleanSecret(context *gin.Context) {
	err := service.Secret.Clean()
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, gin.H{"response": "the secret was cleaned"})
}

func getCredentials(context *gin.Context) {
	credentials := persistence.Database.Credential.GetCredentialMap()
	context.JSON(http.StatusOK, gin.H{"response": credentials})
}

func postCredential(context *gin.Context) {
	var credential Credential
	_ = context.ShouldBindJSON(&credential)
	password, err := service.Encrypt(credential.Password, service.Secret.Get())
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	encryptedCredential := Credential{Username: credential.Username, Password: password}
	key, err := persistence.Database.Credential.CreateCredential(encryptedCredential)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, gin.H{"response": key})
}

func deleteCredential(context *gin.Context) {
	jobUuid, parseErr := uuid.Parse(context.Param("uuid"))
	if parseErr != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": parseErr.Error()})
		return
	}
	deleteErr := persistence.Database.Credential.DeleteCredential(jobUuid)
	if deleteErr != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": parseErr.Error()})
		return
	}
}
