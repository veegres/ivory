package password

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Router struct {
	passwordService *Service
}

func NewRouter(passwordService *Service) *Router {
	return &Router{passwordService: passwordService}
}

func (r *Router) GetCredentials(context *gin.Context) {
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

func (r *Router) PostCredential(context *gin.Context) {
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

func (r *Router) PatchCredential(context *gin.Context) {
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

func (r *Router) DeleteCredential(context *gin.Context) {
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
