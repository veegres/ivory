package router

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	. "ivory/model"
	"ivory/persistence"
	"ivory/service"
	"net/http"
	"strconv"
)

func (r routes) CliGroup(group *gin.RouterGroup) {
	w := service.CreateJobWorker()

	node := group.Group("/cli")
	node.GET("/bloat", getCompactTableList)
	node.GET("/bloat/:uuid", getCompactTable)
	node.GET("/bloat/:uuid/logs", getCompactTableLogs)
	node.GET("/bloat/cluster/:name", getCompactTableListByCluster)

	node.POST("/bloat/job/start", startJob(w))
	node.POST("/bloat/job/:uuid/stop", stopJob(w))
	node.DELETE("/bloat/job/:uuid/delete", deleteJob(w))
	node.GET("/bloat/job/:uuid/stream", streamJob(w))
}

func getCompactTableLogs(context *gin.Context) {
	jobUuid, errUuid := uuid.Parse(context.Param("uuid"))
	if errUuid != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": errUuid.Error()})
		return
	}

	model, errModel := persistence.Database.CompactTable.Get(jobUuid)
	if errModel != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": errModel.Error()})
		return
	}

	context.Writer.Header().Set("Cache-Control", "no-transform")
	context.File(model.LogsPath)
}

func getCompactTableList(context *gin.Context) {
	list, err := persistence.Database.CompactTable.List()
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, gin.H{"response": list})
}

func getCompactTableListByCluster(context *gin.Context) {
	cluster := context.Param("name")
	list, err := persistence.Database.CompactTable.ListByCluster(cluster)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, gin.H{"response": list})
}

func getCompactTable(context *gin.Context) {
	jobUuid, parseErr := uuid.Parse(context.Param("uuid"))
	if parseErr != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": parseErr.Error()})
		return
	}
	compactTable, err := persistence.Database.CompactTable.Get(jobUuid)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, gin.H{"response": compactTable})
}

func startJob(jobWorker service.JobWorker) func(context *gin.Context) {
	return func(context *gin.Context) {
		var cli CompactTableRequest
		_ = context.ShouldBindJSON(&cli)

		// TODO move decryption to another method (encapsulate)
		user := cli.Connection.Username
		pass, errCred := service.Encrypt(cli.Connection.Password, service.Secret.Get())
		// TODO think to move it and do from frontend
		credId, errCred := persistence.Database.Credential.CreateCredential(Credential{Username: user, Password: pass})
		if errCred != nil {
			context.JSON(http.StatusBadRequest, gin.H{"error": errCred.Error()})
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
				sb = append(sb, "--excludeSchema", cli.Target.ExcludeSchema)
				isDefaultTarget = false
			}
			if cli.Target.ExcludeTable != "" {
				sb = append(sb, "--excludeTable", cli.Target.ExcludeTable)
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

		model, errStart := jobWorker.Start(credId, cli.Cluster, sb)
		if errStart != nil {
			context.JSON(http.StatusBadRequest, gin.H{"error": errStart.Error()})
			return
		}

		context.JSON(http.StatusOK, gin.H{"response": model})

	}
}

func stopJob(jobWorker service.JobWorker) func(context *gin.Context) {
	return func(context *gin.Context) {
		jobUuid, errUuid := uuid.Parse(context.Param("uuid"))
		if errUuid != nil {
			context.JSON(http.StatusBadRequest, gin.H{"error": errUuid.Error()})
			return
		}
		errStop := jobWorker.Stop(jobUuid)
		if errStop != nil {
			context.JSON(http.StatusBadRequest, gin.H{"error": errStop.Error()})
			return
		}
	}
}

func deleteJob(jobWorker service.JobWorker) func(context *gin.Context) {
	return func(context *gin.Context) {
		jobUuid, errUuid := uuid.Parse(context.Param("uuid"))
		if errUuid != nil {
			context.JSON(http.StatusBadRequest, gin.H{"error": errUuid.Error()})
			return
		}

		err := jobWorker.Delete(jobUuid)
		if err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}
}

func streamJob(jobWorker service.JobWorker) func(context *gin.Context) {
	return func(context *gin.Context) {
		// notify proxy that it shouldn't enable any caching
		context.Writer.Header().Set("Cache-Control", "no-transform")
		// force using correct event-stream if there is no proxy
		context.Writer.Header().Set("Content-Type", "text/event-stream")

		// start and end stream
		context.SSEvent(STREAM.String(), "start")
		defer context.SSEvent(STREAM.String(), "end")
		context.Writer.Flush()

		// find stream job
		jobUuid, err := uuid.Parse(context.Param("uuid"))
		if err != nil {
			context.SSEvent(SERVER.String(), "Logs streaming failed: Cannot parse UUID")
			context.SSEvent(STATUS.String(), UNKNOWN.String())
			return
		}

		jobWorker.Stream(jobUuid, func(event Event) {
			context.SSEvent(event.Name.String(), event.Message)
			context.Writer.Flush()
		})

		// finish stream (final flush)
		context.Writer.Flush()
	}
}
