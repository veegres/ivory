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
	queryService    *service.QueryService
	postgresGateway *service.PostgresClient
}

func NewQueryRouter(
	queryService *service.QueryService,
	postgresGateway *service.PostgresClient,
) *QueryRouter {
	return &QueryRouter{
		queryService:    queryService,
		postgresGateway: postgresGateway,
	}
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

	res, err := r.queryService.RunQuery(req.QueryUuid, req.CredentialId, req.Db)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"response": res})
}

func (r *QueryRouter) PostDatabasesQuery(context *gin.Context) {
	var req QueryDatabasesRequest
	err := context.ShouldBindJSON(&req)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := r.queryService.DatabasesQuery(req.CredentialId, req.Db, req.Name)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"response": res})
}

func (r *QueryRouter) PostSchemasQuery(context *gin.Context) {
	var req QuerySchemasRequest
	err := context.ShouldBindJSON(&req)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := r.queryService.SchemasQuery(req.CredentialId, req.Db, req.Name)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"response": res})
}

func (r *QueryRouter) PostTablesQuery(context *gin.Context) {
	var req QueryTablesRequest
	err := context.ShouldBindJSON(&req)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := r.queryService.TablesQuery(req.CredentialId, req.Db, req.Schema, req.Name)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"response": res})
}

func (r *QueryRouter) PostCommonChartQuery(context *gin.Context) {
	var req QueryChartRequest
	err := context.ShouldBindJSON(&req)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := r.queryService.CommonChartQuery(req.CredentialId, req.Db)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"response": res})
}

func (r *QueryRouter) PostDatabaseChartQuery(context *gin.Context) {
	var req QueryChartRequest
	err := context.ShouldBindJSON(&req)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := r.queryService.DatabaseChartQuery(req.CredentialId, req.Db)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"response": res})
}

func (r *QueryRouter) PostCancelQuery(context *gin.Context) {
	var req QueryKillRequest
	err := context.ShouldBindJSON(&req)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	errQuery := r.postgresGateway.Cancel(req.Pid, req.CredentialId, req.Db)
	if errQuery != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": errQuery.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"response": "query is canceled"})
}

func (r *QueryRouter) PostTerminateQuery(context *gin.Context) {
	var req QueryKillRequest
	err := context.ShouldBindJSON(&req)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	errQuery := r.postgresGateway.Terminate(req.Pid, req.CredentialId, req.Db)
	if errQuery != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": errQuery.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"response": "query is terminated"})
}
