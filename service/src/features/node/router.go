package node

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Router struct {
	service *Service
}

func NewRouter(service *Service) *Router {
	return &Router{service: service}
}

func handleParamRequest[T any](context *gin.Context, action func(node KeeperConnection) (T, int, error)) {
	query := context.Query("request")
	var node KeeperConnection
	errBind := json.Unmarshal([]byte(query), &node)
	if errBind != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": errBind.Error()})
		return
	}
	body, status, err := action(node)
	if err != nil {
		context.JSON(status, gin.H{"error": err.Error()})
		return
	}
	context.JSON(status, gin.H{"response": body})
}

func handleBodyRequest[T any](context *gin.Context, action func(node KeeperConnection) (T, int, error)) {
	var request KeeperConnection
	errBind := context.ShouldBindJSON(&request)
	if errBind != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": errBind.Error()})
		return
	}
	body, status, err := action(request)
	if err != nil {
		context.JSON(status, gin.H{"error": err.Error()})
		return
	}
	context.JSON(status, gin.H{"response": body})
}

func handleParamRequestOf[R any, T any](context *gin.Context, action func(request R) (T, error)) {
	query := context.Query("request")
	var request R
	errBind := json.Unmarshal([]byte(query), &request)
	if errBind != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": errBind.Error()})
		return
	}
	body, err := action(request)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, gin.H{"response": body})
}

func handleBodyRequestOf[R any, T any](context *gin.Context, action func(request R) (T, error)) {
	var request R
	errBind := context.ShouldBindJSON(&request)
	if errBind != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": errBind.Error()})
		return
	}
	body, err := action(request)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, gin.H{"response": body})
}
