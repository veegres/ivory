package router

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	. "ivory/model"
	"ivory/persistence"
	"ivory/service"
	"net/http"
	"strconv"
)

func (r routes) CredentialGroup(group *gin.RouterGroup) {
	secret := group.Group("/secret")
	secret.GET("", getStatus)
	secret.POST("/set", setSecret)
	secret.POST("/update", updateSecret)
	secret.POST("/clean", cleanSecret)

	credential := group.Group("/credential")
	credential.GET("", getCredentials)
	credential.POST("", postCredential)
	credential.PATCH("/:uuid", patchCredential)
	credential.DELETE("/:uuid", deleteCredential)
}

func getStatus(context *gin.Context) {
	context.JSON(http.StatusOK, gin.H{"response": service.Secret.Status()})
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
	credentialType := context.Request.URL.Query().Get("type")

	var credentials map[string]Credential
	if credentialType != "" {
		number, _ := strconv.Atoi(credentialType)
		credentials = persistence.Database.Credential.ListByType(CredentialType(number))
	} else {
		credentials = persistence.Database.Credential.List()
	}

	context.JSON(http.StatusOK, gin.H{"response": credentials})
}

func postCredential(context *gin.Context) {
	var credential Credential
	err := context.ShouldBindJSON(&credential)
	encryptedPassword, err := service.Encrypt(credential.Password, service.Secret.Get())
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	encryptedCredential := Credential{Username: credential.Username, Password: encryptedPassword, Type: credential.Type}
	key, cred, err := persistence.Database.Credential.Create(encryptedCredential)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"response": gin.H{"key": key.String(), "credential": cred}})
}

func patchCredential(context *gin.Context) {
	credUuid, parseErr := uuid.Parse(context.Param("uuid"))
	if parseErr != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": parseErr.Error()})
		return
	}

	var credential Credential
	err := context.ShouldBindJSON(&credential)
	encryptedPassword, err := service.Encrypt(credential.Password, service.Secret.Get())
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	encryptedCredential := Credential{Username: credential.Username, Password: encryptedPassword, Type: credential.Type}
	_, cred, err := persistence.Database.Credential.Update(credUuid, encryptedCredential)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, gin.H{"response": cred})
}

func deleteCredential(context *gin.Context) {
	credUuid, parseErr := uuid.Parse(context.Param("uuid"))
	if parseErr != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": parseErr.Error()})
		return
	}
	deleteErr := persistence.Database.Credential.Delete(credUuid)
	if deleteErr != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": parseErr.Error()})
		return
	}
}
