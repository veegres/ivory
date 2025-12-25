package bloat

import (
	. "ivory/src/features/bloat/job"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Router struct {
	bloatService *Service
}

func NewRouter(bloatService *Service) *Router {
	return &Router{bloatService: bloatService}
}

func (r *Router) GetBloatLogs(context *gin.Context) {
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

func (r *Router) GetBloatList(context *gin.Context) {
	list, err := r.bloatService.List()
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, gin.H{"response": list})
}

func (r *Router) GetBloatListByCluster(context *gin.Context) {
	cluster := context.Param("name")
	list, err := r.bloatService.ListByCluster(cluster)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, gin.H{"response": list})
}

func (r *Router) GetBloat(context *gin.Context) {
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

func (r *Router) PostJobStart(context *gin.Context) {
	var cli BloatRequest
	parseErr := context.ShouldBindJSON(&cli)
	if parseErr != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": parseErr.Error()})
		return
	}

	if cli.Connection.CredentialId == nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "credentials are required"})
		return
	}

	sb := []string{
		"--host", cli.Connection.Db.Host,
		"--port", strconv.Itoa(cli.Connection.Db.Port),
	}
	isDefaultTarget := true
	if cli.Target != nil {
		if cli.Target.Database != "" {
			sb = append(sb, "--dbname", cli.Target.Database)
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
	if cli.Options.Force {
		sb = append(sb, "--force")
	}
	if cli.Options.NoReindex {
		sb = append(sb, "--no-reindex")
	}
	if cli.Options.NoInitialVacuum {
		sb = append(sb, "--no-initial-vacuum")
	}
	if cli.Options.InitialReindex {
		sb = append(sb, "--initial-reindex")
	}
	if cli.Options.RoutineVacuum {
		sb = append(sb, "--routine-vacuum")
	}
	if cli.Options.DelayRatio != 0 {
		sb = append(sb, "--delay-ratio", strconv.Itoa(cli.Options.DelayRatio))
	}
	if cli.Options.MinTableSize != 0 {
		sb = append(sb, "--min-table-size", strconv.Itoa(cli.Options.MinTableSize))
	}
	if cli.Options.MaxTableSize != 0 {
		sb = append(sb, "--max-table-size", strconv.Itoa(cli.Options.MaxTableSize))
	}
	sb = append(sb, "--verbose")

	model, errStart := r.bloatService.Start(*cli.Connection.CredentialId, cli.Cluster, sb)
	if errStart != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": errStart.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"response": model})
}

func (r *Router) PostJobStop(context *gin.Context) {
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

func (r *Router) DeleteJob(context *gin.Context) {
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

func (r *Router) GetJobStream(context *gin.Context) {
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
		context.SSEvent(SERVER.String(), "Streaming failed: Cannot parse UUID")
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
