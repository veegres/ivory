package query

import (
	"ivory/src/clients/database"
	"ivory/src/features/config"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Router struct {
	service        *Service
	executeService *ExecuteService
	logService     *LogService
	configService  *config.Service
}

func NewRouter(
	service *Service,
	executeService *ExecuteService,
	logService *LogService,
	configService *config.Service,
) *Router {
	return &Router{
		service:        service,
		executeService: executeService,
		logService:     logService,
		configService:  configService,
	}
}

func (r *Router) PutQuery(context *gin.Context) {
	queryUuid, parseErr := uuid.Parse(context.Param("uuid"))
	if parseErr != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": parseErr.Error()})
		return
	}
	var query database.Query
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

func (r *Router) PostQuery(context *gin.Context) {
	var query database.Query
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

func (r *Router) GetQueryList(context *gin.Context) {
	queryTypeStr := context.Request.URL.Query().Get("type")

	if queryTypeStr != "" {
		number, errParse := strconv.ParseInt(queryTypeStr, 10, 8)
		if errParse != nil {
			context.JSON(http.StatusBadRequest, gin.H{"error": errParse.Error()})
			return
		}
		queryType := database.QueryType(number)
		queryList, err := r.service.GetList(&queryType)
		if err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		context.JSON(http.StatusOK, gin.H{"response": queryList})
	} else {
		queryList, err := r.service.GetList(nil)
		if err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		context.JSON(http.StatusOK, gin.H{"response": queryList})
	}
}

func (r *Router) DeleteQuery(context *gin.Context) {
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

func (r *Router) GetQueryLog(context *gin.Context) {
	queryUuid, parseErr := uuid.Parse(context.Param("uuid"))
	if parseErr != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": parseErr.Error()})
		return
	}
	response, err := r.logService.Get(queryUuid)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, gin.H{"response": response})
}

func (r *Router) DeleteQueryLog(context *gin.Context) {
	queryUuid, parseErr := uuid.Parse(context.Param("uuid"))
	if parseErr != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": parseErr.Error()})
		return
	}
	err := r.logService.Delete(queryUuid)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, gin.H{"response": "query log was deleted"})
}

func (r *Router) PostExecuteTempateQuery(context *gin.Context) {
	var req QueryTemplateRequest
	errBind := context.ShouldBindJSON(&req)
	if errBind != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": errBind.Error()})
		return
	}

	if req.QueryUuid == nil || req.QueryUuid.String() == "" {
		context.JSON(http.StatusBadRequest, gin.H{"error": "queryUuid is required"})
		return
	}

	queryContext := r.getQueryContext(context, req.Connection)
	res, err := r.executeService.TemplateQuery(queryContext, *req.QueryUuid, req.QueryOptions)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"response": res})
}

func (r *Router) PostExecuteConsoleQuery(context *gin.Context) {
	var req QueryConsoleRequest
	errBind := context.ShouldBindJSON(&req)
	if errBind != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": errBind.Error()})
		return
	}

	if req.Query == "" {
		context.JSON(http.StatusBadRequest, gin.H{"error": "query is required"})
		return
	}

	queryContext := r.getQueryContext(context, req.Connection)
	res, err := r.executeService.ConsoleQuery(queryContext, req.Query, req.QueryOptions)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"response": res})
}

func (r *Router) PostActivityQuery(context *gin.Context) {
	var req QueryConnection
	errBind := context.ShouldBindJSON(&req)
	if errBind != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": errBind.Error()})
		return
	}

	queryContext := r.getQueryContext(context, req)
	res, err := r.executeService.RunningQueriesByApplicationName(queryContext)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"response": res})
}

func (r *Router) PostDatabasesQuery(context *gin.Context) {
	var req QueryDatabasesRequest
	errBind := context.ShouldBindJSON(&req)
	if errBind != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": errBind.Error()})
		return
	}

	queryContext := r.getQueryContext(context, req.Connection)
	res, err := r.executeService.DatabasesQuery(queryContext, req.Name)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"response": res})
}

func (r *Router) PostSchemasQuery(context *gin.Context) {
	var req QuerySchemasRequest
	errBind := context.ShouldBindJSON(&req)
	if errBind != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": errBind.Error()})
		return
	}

	queryContext := r.getQueryContext(context, req.Connection)
	res, err := r.executeService.SchemasQuery(queryContext, req.Name)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"response": res})
}

func (r *Router) PostTablesQuery(context *gin.Context) {
	var req QueryTablesRequest
	errBind := context.ShouldBindJSON(&req)
	if errBind != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": errBind.Error()})
		return
	}

	queryContext := r.getQueryContext(context, req.Connection)
	res, err := r.executeService.TablesQuery(queryContext, req.Schema, req.Name)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"response": res})
}

func (r *Router) PostChartQuery(context *gin.Context) {
	var req QueryChartRequest
	errBind := context.ShouldBindJSON(&req)
	if errBind != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": errBind.Error()})
		return
	}

	queryContext := r.getQueryContext(context, req.Connection)
	res, err := r.executeService.ChartQuery(queryContext, *req.Type)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"response": res})
}

func (r *Router) PostCancelQuery(context *gin.Context) {
	var req QueryKillRequest
	errBind := context.ShouldBindJSON(&req)
	if errBind != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": errBind.Error()})
		return
	}

	queryContext := r.getQueryContext(context, req.Connection)
	errQuery := r.executeService.CancelQuery(queryContext, req.Pid)
	if errQuery != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": errQuery.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"response": "query is canceled"})
}

func (r *Router) PostTerminateQuery(context *gin.Context) {
	var req QueryKillRequest
	errBind := context.ShouldBindJSON(&req)
	if errBind != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": errBind.Error()})
		return
	}

	queryContext := r.getQueryContext(context, req.Connection)
	errQuery := r.executeService.TerminateQuery(queryContext, req.Pid)
	if errQuery != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": errQuery.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"response": "query is terminated"})
}

func (r *Router) getQueryContext(ctx *gin.Context, con QueryConnection) QueryContext {
	session, errSession := ctx.Cookie("session")
	if errSession != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": errSession.Error()})
	}
	return QueryContext{
		Connection: con,
		Session:    session,
	}
}
