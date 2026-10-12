package query

import (
	"ivory/plugins/database"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (r *HttpRouter) PutQuery(context *gin.Context) {
	queryUuid, parseErr := uuid.Parse(context.Param("uuid"))
	if parseErr != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": parseErr.Error()})
		return
	}
	var query Request
	errBind := context.ShouldBindJSON(&query)
	if errBind != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": errBind.Error()})
		return
	}
	_, response, err := r.service.Update(queryUuid, query)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, gin.H{"response": response})
}

func (r *HttpRouter) PostQuery(context *gin.Context) {
	var query Request
	errBind := context.ShouldBindJSON(&query)
	if errBind != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": errBind.Error()})
		return
	}
	_, response, err := r.service.Create(Manual, query)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, gin.H{"response": response})
}

func (r *HttpRouter) GetQueryList(context *gin.Context) {
	var request RouterListRequest
	if err := context.ShouldBindQuery(&request); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// NOTE: when the plugin is not provided, filter to postgres so clients
	// that predate multi-plugin queries keep seeing exactly what they saw
	if request.Plugin == nil {
		plugin := database.POSTGRES
		request.Plugin = &plugin
	}

	queryList, err := r.service.GetList(request.Type, request.Plugin)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, gin.H{"response": queryList})
}

func (r *HttpRouter) DeleteQuery(context *gin.Context) {
	queryUuid, parseErr := uuid.Parse(context.Param("uuid"))
	if parseErr != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": parseErr.Error()})
		return
	}
	err := r.service.Delete(queryUuid)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, gin.H{"response": "the query was deleted"})
}
