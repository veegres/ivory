package query

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (r *Router) PostExecuteTemplateQuery(context *gin.Context) {
	var req TemplateRequest
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
	res, err := r.service.TemplateQuery(queryContext, *req.QueryUuid, req.Options)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"response": res})
}

func (r *Router) PostExecuteConsoleQuery(context *gin.Context) {
	var req ConsoleRequest
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
	res, err := r.service.ConsoleQuery(queryContext, req.Query, req.Options)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"response": res})
}

func (r *Router) PostActivityQuery(context *gin.Context) {
	var req Connection
	errBind := context.ShouldBindJSON(&req)
	if errBind != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": errBind.Error()})
		return
	}

	queryContext := r.getQueryContext(context, req)
	res, err := r.service.RunningQueriesByApplicationName(queryContext)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"response": res})
}

func (r *Router) PostDatabasesQuery(context *gin.Context) {
	var req DatabasesRequest
	errBind := context.ShouldBindJSON(&req)
	if errBind != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": errBind.Error()})
		return
	}

	queryContext := r.getQueryContext(context, req.Connection)
	res, err := r.service.DatabasesQuery(queryContext, req.Name)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"response": res})
}

func (r *Router) PostSchemasQuery(context *gin.Context) {
	var req SchemasRequest
	errBind := context.ShouldBindJSON(&req)
	if errBind != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": errBind.Error()})
		return
	}

	queryContext := r.getQueryContext(context, req.Connection)
	res, err := r.service.SchemasQuery(queryContext, req.Name)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"response": res})
}

func (r *Router) PostTablesQuery(context *gin.Context) {
	var req TablesRequest
	errBind := context.ShouldBindJSON(&req)
	if errBind != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": errBind.Error()})
		return
	}

	queryContext := r.getQueryContext(context, req.Connection)
	res, err := r.service.TablesQuery(queryContext, req.Schema, req.Name)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"response": res})
}

func (r *Router) PostChartQuery(context *gin.Context) {
	var req ChartRequest
	errBind := context.ShouldBindJSON(&req)
	if errBind != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": errBind.Error()})
		return
	}

	queryContext := r.getQueryContext(context, req.Connection)
	res, err := r.service.ChartQuery(queryContext, req.Type)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"response": res})
}

func (r *Router) PostCancelQuery(context *gin.Context) {
	var req KillRequest
	errBind := context.ShouldBindJSON(&req)
	if errBind != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": errBind.Error()})
		return
	}

	queryContext := r.getQueryContext(context, req.Connection)
	errQuery := r.service.CancelQuery(queryContext, req.Pid)
	if errQuery != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": errQuery.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"response": "query is canceled"})
}

func (r *Router) PostTerminateQuery(context *gin.Context) {
	var req KillRequest
	errBind := context.ShouldBindJSON(&req)
	if errBind != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": errBind.Error()})
		return
	}

	queryContext := r.getQueryContext(context, req.Connection)
	errQuery := r.service.TerminateQuery(queryContext, req.Pid)
	if errQuery != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": errQuery.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"response": "query is terminated"})
}
