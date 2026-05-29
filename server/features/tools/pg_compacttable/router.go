package bloat

import (
	"ivory/features/job"
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

	path, errPath := r.bloatService.GetLogsPath(jobUuid)
	if errPath != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": errPath.Error()})
		return
	}

	context.Writer.Header().Set("Cache-Control", "no-transform")
	context.File(path)
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
	var request RunRequest
	parseErr := context.ShouldBindJSON(&request)
	if parseErr != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": parseErr.Error()})
		return
	}

	sb := []string{
		"--host", request.Db.Host,
		"--port", strconv.Itoa(request.Db.Port),
	}
	isDefaultTarget := true
	if request.Target != nil {
		if request.Target.Database != "" {
			sb = append(sb, "--dbname", request.Target.Database)
			isDefaultTarget = false
		}
		if request.Target.Schema != "" {
			sb = append(sb, "--schema", request.Target.Schema)
			isDefaultTarget = false
		}
		if request.Target.Table != "" {
			sb = append(sb, "--table", request.Target.Table)
			isDefaultTarget = false
		}
		if request.Target.ExcludeSchema != "" {
			sb = append(sb, "--exclude-schema", request.Target.ExcludeSchema)
			isDefaultTarget = false
		}
		if request.Target.ExcludeTable != "" && request.Target.Schema != "" {
			sb = append(sb, "--exclude-table", request.Target.Schema+"."+request.Target.ExcludeTable)
			isDefaultTarget = false
		}
	}
	if isDefaultTarget {
		sb = append(sb, "--all")
	}
	if request.Options.Force {
		sb = append(sb, "--force")
	}
	if request.Options.NoReindex {
		sb = append(sb, "--no-reindex")
	}
	if request.Options.NoInitialVacuum {
		sb = append(sb, "--no-initial-vacuum")
	}
	if request.Options.InitialReindex {
		sb = append(sb, "--initial-reindex")
	}
	if request.Options.RoutineVacuum {
		sb = append(sb, "--routine-vacuum")
	}
	if request.Options.DelayRatio != 0 {
		sb = append(sb, "--delay-ratio", strconv.Itoa(request.Options.DelayRatio))
	}
	if request.Options.MinTableSize != 0 {
		sb = append(sb, "--min-table-size", strconv.Itoa(request.Options.MinTableSize))
	}
	if request.Options.MaxTableSize != 0 {
		sb = append(sb, "--max-table-size", strconv.Itoa(request.Options.MaxTableSize))
	}
	sb = append(sb, "--verbose")

	model, errStart := r.bloatService.Start(request.Cluster, request.VaultId, sb)
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

	context.Writer.Flush()

	// find stream
	jobId := context.Param("uuid")
	if jobId == "" {
		context.SSEvent(job.SERVER.String(), "Streaming failed: Cannot parse UUID")
		context.SSEvent(job.STATUS.String(), job.UNKNOWN.String())
		return
	}
	jobUuid, errUuid := uuid.Parse(jobId)
	if errUuid != nil {
		context.SSEvent(job.SERVER.String(), "Streaming failed: Cannot parse UUID")
		context.SSEvent(job.STATUS.String(), job.UNKNOWN.String())
		return
	}

	session := context.GetString("session")
	r.bloatService.Stream(jobUuid, job.SubscriberID(session), func(event job.Message) {
		context.SSEvent(event.Type.String(), event.Message)
		context.Writer.Flush()
	})

	// finish stream (final flush)
	context.Writer.Flush()
}
