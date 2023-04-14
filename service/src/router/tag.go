package router

import (
	"github.com/gin-gonic/gin"
	"ivory/src/service"
	"net/http"
)

type TagRouter struct {
	tagService *service.TagService
}

func NewTagRouter(tagService *service.TagService) *TagRouter {
	return &TagRouter{tagService: tagService}
}

func (r *TagRouter) GetTagList(context *gin.Context) {
	list, err := r.tagService.List()
	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, gin.H{"response": list})
}
