package router

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	. "ivory/src/model"
	"ivory/src/service"
	"net/http"
	"strconv"
)

type BloatRouter struct {
	bloatService *service.BloatService
}

func NewBloatRouter(bloatService *service.BloatService) *BloatRouter {
	return &BloatRouter{bloatService: bloatService}
}

func (r *BloatRouter) GetCompactTableLogs(context *gin.Context) {
	jobUuid, errUuid := uuid.Parse(context.Param("uuid"))
	if errUuid != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": errUuid.Error()})
		return
	}

	model, errModel := r.bloatService.Get(jobUuid)
	if errModel != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": errModel.Error()})
		return
	}

	context.Writer.Header().Set("Cache-Control", "no-transform")
	context.File(model.LogsPath)
}

func (r *BloatRouter) GetCompactTableList(context *gin.Context) {
	list, err := r.bloatService.List()
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, gin.H{"response": list})
}

func (r *BloatRouter) GetCompactTableListByCluster(context *gin.Context) {
	cluster := context.Param("name")
	list, err := r.bloatService.ListByCluster(cluster)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, gin.H{"response": list})
}

func (r *BloatRouter) GetCompactTable(context *gin.Context) {
	jobUuid, parseErr := uuid.Parse(context.Param("uuid"))
	if parseErr != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": parseErr.Error()})
		return
	}
	compactTable, err := r.bloatService.Get(jobUuid)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, gin.H{"response": compactTable})
}

func (r *BloatRouter) StartJob(context *gin.Context) {
	var cli BloatRequest
	parseErr := context.ShouldBindJSON(&cli)
	if parseErr != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": parseErr.Error()})
		return
	}

	sb := []string{
		"--host", cli.Connection.Host,
		"--port", strconv.Itoa(cli.Connection.Port),
	}
	isDefaultTarget := true
	if cli.Target != nil {
		if cli.Target.DbName != "" {
			sb = append(sb, "--dbname", cli.Target.DbName)
			isDefaultTarget = false
		}
		if cli.Target.Schema != "" {
			sb = append(sb, "--schema", cli.Target.Schema)
			isDefaultTarget = false
		}
		if cli.Target.Table != "" {
			sb = append(sb, "--table", cli.Target.Table)
			isDefaultTarget = false
		}
		if cli.Target.ExcludeSchema != "" {
			sb = append(sb, "--exclude-schema", cli.Target.ExcludeSchema)
			isDefaultTarget = false
		}
		if cli.Target.ExcludeTable != "" && cli.Target.Schema != "" {
			sb = append(sb, "--exclude-table", cli.Target.Schema+"."+cli.Target.ExcludeTable)
			isDefaultTarget = false
		}
	}
	if isDefaultTarget {
		sb = append(sb, "--all")
	}
	if cli.Ratio != 0 {
		sb = append(sb, "--delay-ratio", strconv.Itoa(cli.Ratio))
	}
	sb = append(sb, "--verbose")

	model, errStart := r.bloatService.Start(cli.Connection.CredId, cli.Cluster, sb)
	if errStart != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": errStart.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"response": model})
}

func (r *BloatRouter) StopJob(context *gin.Context) {
	jobUuid, errUuid := uuid.Parse(context.Param("uuid"))
	if errUuid != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": errUuid.Error()})
		return
	}
	errStop := r.bloatService.Stop(jobUuid)
	if errStop != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": errStop.Error()})
		return
	}
}

func (r *BloatRouter) DeleteJob(context *gin.Context) {
	jobUuid, errUuid := uuid.Parse(context.Param("uuid"))
	if errUuid != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": errUuid.Error()})
		return
	}

	err := r.bloatService.Delete(jobUuid)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
}

func (r *BloatRouter) StreamJob(context *gin.Context) {
	// notify proxy that it shouldn't enable any caching
	context.Writer.Header().Set("Cache-Control", "no-transform")
	// force using correct event-stream if there is no proxy
	context.Writer.Header().Set("Content-Type", "text/event-stream")

	// start and end stream
	context.SSEvent(STREAM.String(), START.String())
	defer context.SSEvent(STREAM.String(), END.String())
	context.Writer.Flush()

	// find stream job
	jobUuid, err := uuid.Parse(context.Param("uuid"))
	if err != nil {
		context.SSEvent(SERVER.String(), "Logs streaming failed: Cannot parse UUID")
		context.SSEvent(STATUS.String(), UNKNOWN.String())
		return
	}

	r.bloatService.Stream(jobUuid, func(event Event) {
		context.SSEvent(event.Type.String(), event.Message)
		context.Writer.Flush()
	})

	// finish stream (final flush)
	context.Writer.Flush()
}
