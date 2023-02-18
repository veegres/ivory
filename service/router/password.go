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

type PasswordRouter struct {
	repository    *persistence.PasswordRepository
	secretService *service.SecretService
	encryption    *service.Encryption
}

func NewPasswordRouter(
	repository *persistence.PasswordRepository,
	secretService *service.SecretService,
	encryption *service.Encryption,
) *PasswordRouter {
	return &PasswordRouter{repository: repository, secretService: secretService, encryption: encryption}
}

func (r *PasswordRouter) GetCredentials(context *gin.Context) {
	credentialType := context.Request.URL.Query().Get("type")

	var err error
	var credentials map[string]Credential
	if credentialType != "" {
		number, _ := strconv.Atoi(credentialType)
		credentials, err = r.repository.ListByType(CredentialType(number))
	} else {
		credentials, err = r.repository.List()
	}

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"response": credentials})
}

func (r *PasswordRouter) PostCredential(context *gin.Context) {
	var credential Credential
	err := context.ShouldBindJSON(&credential)
	encryptedPassword, err := r.encryption.Encrypt(credential.Password, r.secretService.Get())
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	encryptedCredential := Credential{Username: credential.Username, Password: encryptedPassword, Type: credential.Type}
	key, cred, err := r.repository.Create(encryptedCredential)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"response": gin.H{"key": key.String(), "credential": cred}})
}

func (r *PasswordRouter) PatchCredential(context *gin.Context) {
	credUuid, parseErr := uuid.Parse(context.Param("uuid"))
	if parseErr != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": parseErr.Error()})
		return
	}

	var credential Credential
	err := context.ShouldBindJSON(&credential)
	encryptedPassword, err := r.encryption.Encrypt(credential.Password, r.secretService.Get())
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	encryptedCredential := Credential{Username: credential.Username, Password: encryptedPassword, Type: credential.Type}
	_, cred, err := r.repository.Update(credUuid, encryptedCredential)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, gin.H{"response": cred})
}

func (r *PasswordRouter) DeleteCredential(context *gin.Context) {
	credUuid, parseErr := uuid.Parse(context.Param("uuid"))
	if parseErr != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": parseErr.Error()})
		return
	}
	deleteErr := r.repository.Delete(credUuid)
	if deleteErr != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": parseErr.Error()})
		return
	}
}
