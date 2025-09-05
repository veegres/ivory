package query

import (
	"ivory/src/clients/database"
	"ivory/src/features/config"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type QueryRouter struct {
	queryService  *QueryService
	configService *config.Service
}

func NewQueryRouter(
	queryService *QueryService,
	configService *config.Service,
) *QueryRouter {
	return &QueryRouter{
		queryService:  queryService,
		configService: configService,
	}
}

func (r *QueryRouter) ManualMiddleware() gin.HandlerFunc {
	return func(context *gin.Context) {
		appConfig, errConfig := r.configService.GetAppConfig()
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
	var query database.QueryRequest
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
	var query database.QueryRequest
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
		number, errParse := strconv.ParseInt(queryTypeStr, 10, 8)
		if errParse != nil {
			context.JSON(http.StatusBadRequest, gin.H{"error": errParse.Error()})
			return
		}
		queryType := database.QueryType(number)
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
	context.JSON(http.StatusOK, gin.H{"response": "the query was deleted"})
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

	var res *database.QueryFields
	var err error
	queryContext := r.getQueryContext(context, req.Connection)
	if req.QueryUuid != nil {
		res, err = r.queryService.RunTemplateQuery(queryContext, *req.QueryUuid, req.QueryOptions)
	}
	if req.Query != nil {
		res, err = r.queryService.RunQuery(queryContext, *req.Query, req.QueryOptions)
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

func (r *QueryRouter) PostAllRunningQueriesByApplicationName(context *gin.Context) {
	var req QueryConnection
	errBind := context.ShouldBindJSON(&req)
	if errBind != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": errBind.Error()})
		return
	}

	queryContext := r.getQueryContext(context, req)
	res, err := r.queryService.GetAllRunningQueriesByApplicationName(queryContext)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"response": res})
}

func (r *QueryRouter) PostDatabasesQuery(context *gin.Context) {
	var req QueryDatabasesRequest
	errBind := context.ShouldBindJSON(&req)
	if errBind != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": errBind.Error()})
		return
	}

	queryContext := r.getQueryContext(context, req.Connection)
	res, err := r.queryService.DatabasesQuery(queryContext, req.Name)
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

	queryContext := r.getQueryContext(context, req.Connection)
	res, err := r.queryService.SchemasQuery(queryContext, req.Name)
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

	queryContext := r.getQueryContext(context, req.Connection)
	res, err := r.queryService.TablesQuery(queryContext, req.Schema, req.Name)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"response": res})
}

func (r *QueryRouter) PostChartQuery(context *gin.Context) {
	var req QueryChartRequest
	errBind := context.ShouldBindJSON(&req)
	if errBind != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": errBind.Error()})
		return
	}

	queryContext := r.getQueryContext(context, req.Connection)
	res, err := r.queryService.ChartQuery(queryContext, *req.Type)
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

	queryContext := r.getQueryContext(context, req.Connection)
	errQuery := r.queryService.Cancel(queryContext, req.Pid)
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

	queryContext := r.getQueryContext(context, req.Connection)
	errQuery := r.queryService.Terminate(queryContext, req.Pid)
	if errQuery != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": errQuery.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"response": "query is terminated"})
}

func (r *QueryRouter) getQueryContext(ctx *gin.Context, con QueryConnection) QueryContext {
	session, errSession := ctx.Cookie("session")
	if errSession != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": errSession.Error()})
	}
	return QueryContext{
		Connection: con,
		Session:    session,
	}
}
