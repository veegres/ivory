package tag

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type HttpRouter struct {
	tagService *Service
}

func NewHttpRouter(tagService *Service) *HttpRouter {
	return &HttpRouter{tagService: tagService}
}

func (r *HttpRouter) GetTagList(context *gin.Context) {
	list, err := r.tagService.List()
	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, gin.H{"response": list})
}
