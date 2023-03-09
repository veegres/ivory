package router

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	. "ivory/src/model"
	"ivory/src/service"
	"net/http"
	"strconv"
)

type QueryRouter struct {
	queryService *service.QueryService
}

func NewQueryRouter(queryService *service.QueryService) *QueryRouter {
	return &QueryRouter{queryService: queryService}
}

func (r *QueryRouter) PutQuery(context *gin.Context) {
	queryUuid, parseErr := uuid.Parse(context.Param("uuid"))
	if parseErr != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": parseErr.Error()})
		return
	}
	var query QueryRequest
	err := context.ShouldBindJSON(&query)
	_, response, err := r.queryService.Update(queryUuid, query)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, gin.H{"response": response})
}

func (r *QueryRouter) PostQuery(context *gin.Context) {
	var query QueryRequest
	err := context.ShouldBindJSON(&query)
	_, response, err := r.queryService.Create(Manual, query)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, gin.H{"response": response})
}

func (r *QueryRouter) GetQueryList(context *gin.Context) {
	queryTypeStr := context.Request.URL.Query().Get("type")

	if queryTypeStr != "" {
		number, _ := strconv.Atoi(queryTypeStr)
		queryType := QueryType(number)
		queryList, err := r.queryService.GetList(&queryType)
		if err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		context.JSON(http.StatusOK, gin.H{"response": queryList})
	} else {
		queryList, err := r.queryService.GetList(nil)
		if err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		context.JSON(http.StatusOK, gin.H{"response": queryList})
	}
}

func (r *QueryRouter) DeleteQuery(context *gin.Context) {
	queryUuid, parseErr := uuid.Parse(context.Param("uuid"))
	if parseErr != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": parseErr.Error()})
		return
	}
	err := r.queryService.Delete(queryUuid)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, gin.H{"response": "query was deleted"})
}

func (r *QueryRouter) PostRunQuery(context *gin.Context) {
	var req QueryRunRequest
	err := context.ShouldBindJSON(&req)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fields, rows, err := r.queryService.Run(req.QueryUuid, req.ClusterName, req.Db)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"response": QueryRunResponse{
		Fields: fields,
		Rows:   rows,
	}})
}
