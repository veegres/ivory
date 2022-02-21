package router

import (
	"bufio"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"io"
	. "ivory/model"
	"ivory/persistence"
	"net/http"
	"os/exec"
	"strconv"
	"strings"
)

func (r routes) CliGroup(group *gin.RouterGroup) {
	pendingJobs := make(chan CompactTableModel)
	activeJobs := make(map[uuid.UUID]Job)
	go job(pendingJobs, activeJobs)

	node := group.Group("/cli")
	node.GET("/bloat", getCompactTableList)
	node.POST("/bloat", postCompactTable(pendingJobs, activeJobs))
	node.DELETE("/bloat/:uuid", deleteCompactTable)
	node.GET("/bloat/:uuid", getCompactTable)
	node.GET("/bloat/:uuid/logs", getCompactTableLogs)
	node.GET("/bloat/:uuid/stream", stream(activeJobs))
}

func getCompactTableLogs(context *gin.Context) {
	jobUuid, _ := uuid.Parse(context.Param("uuid"))
	job, _ := persistence.Database.CompactTable.Get(jobUuid)
	context.File(job.LogsPath)
}

func getCompactTableList(context *gin.Context) {
	list, err := persistence.Database.CompactTable.List()
	if err != nil {
		_ = context.AbortWithError(http.StatusBadRequest, err)
	}
	context.AbortWithStatusJSON(http.StatusOK, gin.H{"response": list})
}

func getCompactTable(context *gin.Context) {
	jobUuid, parseErr := uuid.Parse(context.Param("uuid"))
	if parseErr != nil {
		_ = context.AbortWithError(http.StatusBadRequest, parseErr)
	}
	compactTable, err := persistence.Database.CompactTable.Get(jobUuid)
	if err != nil {
		_ = context.AbortWithError(http.StatusBadRequest, parseErr)
	}

	context.AbortWithStatusJSON(http.StatusOK, gin.H{"response": compactTable})
}

func deleteCompactTable(context *gin.Context) {
	jobUuid, _ := uuid.Parse(context.Param("uuid"))
	errFile := persistence.File.CompactTable.Delete(jobUuid)
	if errFile != nil {
		_ = context.AbortWithError(http.StatusBadRequest, errFile)
	}
	errDb := persistence.Database.CompactTable.Delete(jobUuid)
	if errDb != nil {
		_ = context.AbortWithError(http.StatusBadRequest, errDb)
	}
}

func postCompactTable(pending chan CompactTableModel, active map[uuid.UUID]Job) func(context *gin.Context) {
	return func(context *gin.Context) {
		var cli CompactTableRequest
		_ = context.ShouldBindJSON(&cli)

		sb := []string{
			"--host", cli.Connection.Host,
			"--port", strconv.Itoa(cli.Connection.Port),
			"--user", cli.Connection.Username,
			"--password", cli.Connection.Password,
		}
		if cli.Target != nil {
			if cli.Target.DbName != "" {
				sb = append(sb, "--dbname", cli.Target.DbName)
			}
			if cli.Target.Schema != "" {
				sb = append(sb, "--schema", cli.Target.Schema)
			}
			if cli.Target.Table != "" {
				sb = append(sb, "--table", cli.Target.Table)
			}
			if cli.Target.ExcludeSchema != "" {
				sb = append(sb, "--excludeSchema", cli.Target.ExcludeSchema)
			}
			if cli.Target.ExcludeTable != "" {
				sb = append(sb, "--excludeTable", cli.Target.ExcludeTable)
			}
		}
		if len(sb) == 8 {
			sb = append(sb, "--all")
		}
		if cli.Ratio != 0 {
			sb = append(sb, "--delay-ratio", strconv.Itoa(cli.Ratio))
		}
		sb = append(sb, "--verbose")

		jobUuid := uuid.New()
		compactTableModel := CompactTableModel{
			Uuid:        jobUuid,
			Status:      PENDING,
			Command:     "pgcompacttable " + strings.Join(sb, " "),
			CommandArgs: sb,
			LogsPath:    persistence.File.CompactTable.Create(jobUuid),
		}
		err := persistence.Database.CompactTable.Update(compactTableModel)
		if err != nil {
			_ = context.AbortWithError(http.StatusBadRequest, err)
		}
		active[jobUuid] = Job{Event: make(chan Event)}
		go func() {
			pending <- compactTableModel
		}()

		context.AbortWithStatusJSON(http.StatusOK, gin.H{"response": compactTableModel})
	}
}

func stream(jobs map[uuid.UUID]Job) func(context *gin.Context) {
	return func(context *gin.Context) {
		// notify proxy that it shouldn't enable any caching
		context.Writer.Header().Set("Cache-Control", "no-transform")
		// force using correct event-stream if there is no proxy
		context.Writer.Header().Set("Content-Type", "text/event-stream")
		context.SSEvent(STREAM.String(), "start")
		context.Writer.Flush()

		jobUuid, err := uuid.Parse(context.Param("uuid"))
		job := jobs[jobUuid]
		job.IncrementSubs()
		if err != nil || job.Event == nil {
			context.SSEvent(SERVER.String(), "Logs streaming failed: 404 Not Found")
			return
		}

		context.Stream(func(w io.Writer) bool {
			if event, ok := <-job.Event; ok {
				context.SSEvent(event.Name.String(), event.Message)
				return true
			}

			return false
		})

		job.DecrementSubs()
		context.SSEvent(STREAM.String(), "end")
	}
}

func job(pending chan CompactTableModel, active map[uuid.UUID]Job) {
	initialJobs(pending)

	// TODO change to run multiple jobs
	for model := range pending {
		jobUuid := model.Uuid
		event := active[jobUuid].Event

		cmd := exec.Command("pgcompacttable", model.CommandArgs...)
		pipe, _ := cmd.StdoutPipe()
		if err := cmd.Start(); err != nil {
			event <- Event{Name: SERVER, Message: err.Error()}
			updateStatus(model, event, FAILED)
			return
		}

		updateStatus(model, event, RUNNING)
		reader := bufio.NewReader(pipe)
		line, _, err := reader.ReadLine()
		file, _ := persistence.File.CompactTable.Open(model.LogsPath)
		for err == nil {
			lineString := string(line)
			_, errFile := file.WriteString(lineString + "\n")
			if errFile != nil {
				updateStatus(model, event, FAILED)
				return
			}
			event <- Event{Name: LOG, Message: lineString}
			line, _, err = reader.ReadLine()
		}

		updateStatus(model, event, FINISHED)
		close(event)

		go cleanJob(jobUuid, active)
	}
}

func cleanJob(jobUuid uuid.UUID, jobs map[uuid.UUID]Job) {
	// TODO this thing DOESN't work
	for jobs[jobUuid].GetSubs() == 0 {
		delete(jobs, jobUuid)
		return
	}
}

func initialJobs(pending chan CompactTableModel) {
	runningJobs, _ := persistence.Database.CompactTable.ListByStatus(RUNNING)
	for _, runningJob := range runningJobs {
		_ = persistence.Database.CompactTable.UpdateStatus(runningJob, FAILED)
	}

	pendingJobs, _ := persistence.Database.CompactTable.ListByStatus(PENDING)
	for _, pendingJob := range pendingJobs {
		pending <- pendingJob
	}
}

func updateStatus(model CompactTableModel, event chan Event, status JobStatus) {
	dbErr := persistence.Database.CompactTable.UpdateStatus(model, status)
	if dbErr != nil {
		event <- Event{Name: SERVER, Message: dbErr.Error()}
	} else {
		event <- Event{Name: STATUS, Message: status.String()}
	}
}
