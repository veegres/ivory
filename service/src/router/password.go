package router

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	. "ivory/src/model"
	"ivory/src/service"
	"net/http"
	"strconv"
)

type PasswordRouter struct {
	passwordService *service.PasswordService
}

func NewPasswordRouter(passwordService *service.PasswordService) *PasswordRouter {
	return &PasswordRouter{passwordService: passwordService}
}

func (r *PasswordRouter) GetCredentials(context *gin.Context) {
	stringType := context.Request.URL.Query().Get("type")
	if stringType != "" {
		number, errAtoi := strconv.Atoi(stringType)
		if errAtoi != nil {
			context.JSON(http.StatusBadRequest, gin.H{"error": errAtoi.Error()})
			return
		}
		credType := PasswordType(number)
		credentials, err := r.passwordService.Map(&credType)
		if err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		context.JSON(http.StatusOK, gin.H{"response": credentials})
	} else {
		credentials, err := r.passwordService.Map(nil)
		if err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		context.JSON(http.StatusOK, gin.H{"response": credentials})
	}
}

func (r *PasswordRouter) PostCredential(context *gin.Context) {
	var credential Password
	errBind := context.ShouldBindJSON(&credential)
	if errBind != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": errBind.Error()})
		return
	}

	key, cred, err := r.passwordService.Create(credential)
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

	var credential Password
	errBind := context.ShouldBindJSON(&credential)
	if errBind != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": errBind.Error()})
		return
	}

	_, cred, err := r.passwordService.Update(credUuid, credential)
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
	deleteErr := r.passwordService.Delete(credUuid)
	if deleteErr != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": deleteErr.Error()})
		return
	}
}
