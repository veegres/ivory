package node

import (
	"encoding/json"
	"ivory/core/config"
	"ivory/core/service/job"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Router struct {
	service *Service
}

func NewRouter(service *Service) *Router {
	return &Router{service: service}
}

func handleKeeperParamRequest[T any](context *gin.Context, action func(node KeeperOneRequest) (T, int, error)) {
	query := context.Query("request")
	var node KeeperOneRequest
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

func handleKeeperBodyRequest[T any](context *gin.Context, action func(node KeeperOneRequest) (T, int, error)) {
	var request KeeperOneRequest
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

func handlePlatformParamRequest[R any, T any](context *gin.Context, action func(request R) (T, error)) {
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

func handlePlatformBodyRequest[R any, T any](context *gin.Context, action func(request R) (T, error)) {
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

func handleStreamRequest[R any](context *gin.Context, action func(request R, subscriberID job.SubscriberID, close <-chan struct{}, send func(event job.Message))) {
	context.Writer.Header().Set("Cache-Control", "no-transform")
	context.Writer.Header().Set("Content-Type", "text/event-stream")
	context.Writer.Flush()

	query := context.Query("request")
	var request R
	if err := json.Unmarshal([]byte(query), &request); err != nil {
		context.SSEvent(job.SERVER.String(), "Streaming failed: "+err.Error())
		return
	}

	ctx := context.Request.Context()
	session := context.GetString(env.AuthContextKey.Session)
	action(request, job.SubscriberID(session), ctx.Done(), func(event job.Message) {
		context.SSEvent(event.Type.String(), event.Message)
		context.Writer.Flush()
	})

	context.Writer.Flush()
}
