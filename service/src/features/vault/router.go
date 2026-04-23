package vault

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Router struct {
	vaultService *Service
}

func NewRouter(vaultService *Service) *Router {
	return &Router{vaultService: vaultService}
}

func (r *Router) GetVaultList(context *gin.Context) {
	stringType := context.Request.URL.Query().Get("type")
	if stringType != "" {
		number, errAtoi := strconv.Atoi(stringType)
		if errAtoi != nil {
			context.JSON(http.StatusBadRequest, gin.H{"error": errAtoi.Error()})
			return
		}
		credType := VaultType(number)
		vaults, err := r.vaultService.Map(&credType)
		if err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		context.JSON(http.StatusOK, gin.H{"response": vaults})
	} else {
		vaults, err := r.vaultService.Map(nil)
		if err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		context.JSON(http.StatusOK, gin.H{"response": vaults})
	}
}

func (r *Router) PostVault(context *gin.Context) {
	var vault Vault
	errBind := context.ShouldBindJSON(&vault)
	if errBind != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": errBind.Error()})
		return
	}

	key, cred, err := r.vaultService.Create(vault)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"response": gin.H{"key": key.String(), "vault": cred}})
}

func (r *Router) PatchVault(context *gin.Context) {
	credUuid, parseErr := uuid.Parse(context.Param("uuid"))
	if parseErr != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": parseErr.Error()})
		return
	}

	var vault Vault
	errBind := context.ShouldBindJSON(&vault)
	if errBind != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": errBind.Error()})
		return
	}

	_, cred, err := r.vaultService.Update(credUuid, vault)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"response": cred})
}

func (r *Router) DeleteVault(context *gin.Context) {
	credUuid, parseErr := uuid.Parse(context.Param("uuid"))
	if parseErr != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": parseErr.Error()})
		return
	}
	deleteErr := r.vaultService.Delete(credUuid)
	if deleteErr != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": deleteErr.Error()})
		return
	}
}
