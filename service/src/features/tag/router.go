package tag

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Router struct {
	tagService *Service
}

func NewRouter(tagService *Service) *Router {
	return &Router{tagService: tagService}
}

func (r *Router) GetTagList(context *gin.Context) {
	list, err := r.tagService.List()
	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, gin.H{"response": list})
}
