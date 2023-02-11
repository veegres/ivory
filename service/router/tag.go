package router

import (
	"github.com/gin-gonic/gin"
	"ivory/persistence"
	"net/http"
)

func (r routes) TagGroup(group *gin.RouterGroup) {
	tag := group.Group("/tag")
	tag.GET("", getTagList)
}

func getTagList(context *gin.Context) {
	list, err := persistence.BoltDB.Tag.List()
	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, gin.H{"response": list})
}
