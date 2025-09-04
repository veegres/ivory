package tag

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type TagRouter struct {
	tagService *TagService
}

func NewTagRouter(tagService *TagService) *TagRouter {
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
