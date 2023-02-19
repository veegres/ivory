package router

import (
	"github.com/gin-gonic/gin"
	"ivory/src/persistence"
	"net/http"
)

type TagRouter struct {
	repository *persistence.TagRepository
}

func NewTagRouter(repository *persistence.TagRepository) *TagRouter {
	return &TagRouter{repository: repository}
}

func (r *TagRouter) GetTagList(context *gin.Context) {
	list, err := r.repository.List()
	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, gin.H{"response": list})
}
