package deployment

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Router struct {
	service *Service
}

func NewRouter(service *Service) *Router {
	return &Router{service: service}
}

func (r *Router) GetDeploymentTemplateList(context *gin.Context) {
	var criteria ListRequest
	if err := context.ShouldBindQuery(&criteria); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	list, err := r.service.List(criteria)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, gin.H{"response": list})
}

func (r *Router) PostDeploymentTemplate(context *gin.Context) {
	var request TemplateRequest
	if err := context.ShouldBindJSON(&request); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	response, err := r.service.Create(request)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, gin.H{"response": response})
}

func (r *Router) PutDeploymentTemplate(context *gin.Context) {
	id, err := parseTemplateId(context)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var request TemplateRequest
	if err := context.ShouldBindJSON(&request); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	response, err := r.service.Update(id, request)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, gin.H{"response": response})
}

func (r *Router) DeleteDeploymentTemplate(context *gin.Context) {
	id, err := parseTemplateId(context)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := r.service.Delete(id); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, gin.H{"response": "deleted"})
}

// parseTemplateId turns a shipped template's synthetic id into the read-only
// error rather than an unhelpful uuid parse failure, so a write aimed at one
// says why it is refused.
func parseTemplateId(context *gin.Context) (uuid.UUID, error) {
	param := context.Param("uuid")
	if isDefaultId(param) {
		return uuid.Nil, ErrTemplateReadOnly
	}
	id, err := uuid.Parse(param)
	if err != nil {
		return uuid.Nil, errors.New("invalid template id")
	}
	return id, nil
}
