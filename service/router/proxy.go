package router

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
)

// ProxyGroup TODO make code independent of patroni
func (r routes) ProxyGroup(group *gin.RouterGroup) {
	node := group.Group("/node")
	node.GET("/:host/cluster", getNodeCluster)
	node.GET("/:host/overview", getNodeOverview)
	node.GET("/:host/config", getNodeConfig)
	node.PATCH("/:host/config", patchNodeConfig)
	node.POST("/:host/switchover", postNodeSwitchover)
	node.POST("/:host/reinitialize", postNodeReinitialize)
}

func getNodeCluster(context *gin.Context) {
	host := context.Param("host")
	response, err := http.Get("http://" + host + "/cluster")
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		var body interface{}
		_ = json.NewDecoder(response.Body).Decode(&body)
		context.JSON(http.StatusOK, gin.H{"response": body})
	}
}

func getNodeOverview(context *gin.Context) {
	host := context.Param("host")
	response, err := http.Get("http://" + host + "/patroni")
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		var body interface{}
		_ = json.NewDecoder(response.Body).Decode(&body)
		context.JSON(http.StatusOK, gin.H{"response": body})
	}
}

func getNodeConfig(context *gin.Context) {
	host := context.Param("host")
	response, err := http.Get("http://" + host + "/config")
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		var body interface{}
		_ = json.NewDecoder(response.Body).Decode(&body)
		context.JSON(http.StatusOK, gin.H{"response": body})
	}
}

func patchNodeConfig(context *gin.Context) {
	host := context.Param("host")
	body, _ := ioutil.ReadAll(context.Request.Body)
	req, _ := http.NewRequest(http.MethodPatch, "http://"+host+"/config", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Content-Length", strconv.FormatInt(req.ContentLength, 10))
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		var body interface{}
		_ = json.NewDecoder(response.Body).Decode(&body)
		context.JSON(http.StatusOK, gin.H{"response": body})
	}
}

func postNodeSwitchover(context *gin.Context) {
	host := context.Param("host")
	body, _ := ioutil.ReadAll(context.Request.Body)
	req, _ := http.NewRequest(http.MethodPost, "http://"+host+"/switchover", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Content-Length", strconv.FormatInt(req.ContentLength, 10))
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		bodyBytes, _ := io.ReadAll(response.Body)
		bodyMessage := string(bodyBytes)
		context.JSON(http.StatusOK, gin.H{"response": bodyMessage})
	}
}

func postNodeReinitialize(context *gin.Context) {
	host := context.Param("host")
	body, _ := ioutil.ReadAll(context.Request.Body)
	req, _ := http.NewRequest(http.MethodPost, "http://"+host+"/reinitialize", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Content-Length", strconv.FormatInt(req.ContentLength, 10))
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		bodyBytes, _ := io.ReadAll(response.Body)
		bodyMessage := string(bodyBytes)
		context.JSON(http.StatusOK, gin.H{"response": bodyMessage})
	}
}
