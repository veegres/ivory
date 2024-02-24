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
	generalService  *service.GeneralService
	postgresGateway *service.PostgresClient
}

func NewQueryRouter(
	queryService *service.QueryService,
	generalService *service.GeneralService,
	postgresGateway *service.PostgresClient,
) *QueryRouter {
	return &QueryRouter{
		queryService:    queryService,
		generalService:  generalService,
		postgresGateway: postgresGateway,
	}
}

func (r *QueryRouter) ManualMiddleware() gin.HandlerFunc {
	return func(context *gin.Context) {
		appConfig, errConfig := r.generalService.GetAppConfig()
		if errConfig != nil {
			context.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": errConfig})
			return
		}
		if !appConfig.Availability.ManualQuery {
			context.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "manual queries are restricted"})
			return
		}
		context.Next()
	}
}

func (r *QueryRouter) PutQuery(context *gin.Context) {
	queryUuid, parseErr := uuid.Parse(context.Param("uuid"))
	if parseErr != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": parseErr.Error()})
		return
	}
	var query QueryRequest
	errBind := context.ShouldBindJSON(&query)
	if errBind != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": errBind.Error()})
		return
	}
	_, response, err := r.queryService.Update(queryUuid, query)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, gin.H{"response": response})
}

func (r *QueryRouter) PostQuery(context *gin.Context) {
	var query QueryRequest
	errBind := context.ShouldBindJSON(&query)
	if errBind != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": errBind.Error()})
		return
	}
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
		number, errAtoi := strconv.Atoi(queryTypeStr)
		if errAtoi != nil {
			context.JSON(http.StatusBadRequest, gin.H{"error": errAtoi.Error()})
			return
		}
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
	errBind := context.ShouldBindJSON(&req)
	if errBind != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": errBind.Error()})
		return
	}

	if req.QueryUuid != nil && req.Query != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Only one parameter allowed either `queryUuid` or `query`"})
		return
	}
	if (req.QueryUuid == nil || req.QueryUuid.String() == "") && (req.Query == nil || *req.Query == "") {
		context.JSON(http.StatusBadRequest, gin.H{"error": "At least `queryUuid` or `query` are required"})
		return
	}

	var res *QueryFields
	var err error
	if req.QueryUuid != nil {
		res, err = r.queryService.RunTemplateQuery(*req.QueryUuid, req.QueryParams, req.CredentialId, req.Db)
	}
	if req.Query != nil {
		res, err = r.queryService.RunQuery(*req.Query, req.QueryParams, req.CredentialId, req.Db)
	}

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"response": res})
}

func (r *QueryRouter) GetQueryLog(context *gin.Context) {
	queryUuid, parseErr := uuid.Parse(context.Param("uuid"))
	if parseErr != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": parseErr.Error()})
		return
	}
	response, err := r.queryService.GetLog(queryUuid)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, gin.H{"response": response})
}

func (r *QueryRouter) DeleteQueryLog(context *gin.Context) {
	queryUuid, parseErr := uuid.Parse(context.Param("uuid"))
	if parseErr != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": parseErr.Error()})
		return
	}
	err := r.queryService.DeleteLog(queryUuid)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, gin.H{"response": "query log was deleted"})
}

func (r *QueryRouter) PostDatabasesQuery(context *gin.Context) {
	var req QueryDatabasesRequest
	errBind := context.ShouldBindJSON(&req)
	if errBind != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": errBind.Error()})
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
	errBind := context.ShouldBindJSON(&req)
	if errBind != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": errBind.Error()})
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
	errBind := context.ShouldBindJSON(&req)
	if errBind != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": errBind.Error()})
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
	errBind := context.ShouldBindJSON(&req)
	if errBind != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": errBind.Error()})
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
	errBind := context.ShouldBindJSON(&req)
	if errBind != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": errBind.Error()})
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
	errBind := context.ShouldBindJSON(&req)
	if errBind != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": errBind.Error()})
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
	errBind := context.ShouldBindJSON(&req)
	if errBind != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": errBind.Error()})
		return
	}

	errQuery := r.postgresGateway.Terminate(req.Pid, req.CredentialId, req.Db)
	if errQuery != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": errQuery.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"response": "query is terminated"})
}
