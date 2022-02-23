package router

import (
	"bufio"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"io"
	. "ivory/model"
	"ivory/persistence"
	. "ivory/service"
	"net/http"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"
)

func (r routes) CliGroup(group *gin.RouterGroup) {
	pendingJobs := make(chan CompactTableModel)
	activeJobs := make(map[uuid.UUID]*Job)
	w := worker{pendingJobs, activeJobs}.Build()

	node := group.Group("/cli")
	node.GET("/bloat", getCompactTableList)
	node.POST("/bloat", w.postCompactTable)
	node.DELETE("/bloat/:uuid", deleteCompactTable)
	node.GET("/bloat/:uuid", getCompactTable)
	node.GET("/bloat/:uuid/logs", getCompactTableLogs)
	node.GET("/bloat/:uuid/stream", w.stream)
}

func getCompactTableLogs(context *gin.Context) {
	jobUuid, _ := uuid.Parse(context.Param("uuid"))
	job, _ := persistence.Database.CompactTable.Get(jobUuid)
	context.Writer.Header().Set("Cache-Control", "no-transform")
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

type worker struct {
	pendingJobs chan CompactTableModel
	activeJobs  map[uuid.UUID]*Job
}

func (w worker) Build() worker {
	go w.initiate()
	go w.runner()
	return w
}

func (w worker) postCompactTable(context *gin.Context) {
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
		CreatedAt:   time.Now().UnixNano(),
	}
	err := persistence.Database.CompactTable.Update(compactTableModel)
	if err != nil {
		_ = context.AbortWithError(http.StatusBadRequest, err)
	}

	w.activeJobs[jobUuid] = &Job{Events: make(map[chan Event]bool), Subscribers: new(int), Status: PENDING, Mutex: &sync.Mutex{}}
	w.pendingJobs <- compactTableModel

	context.AbortWithStatusJSON(http.StatusOK, gin.H{"response": compactTableModel})
}

func (w worker) stream(context *gin.Context) {
	// notify proxy that it shouldn't enable any caching
	context.Writer.Header().Set("Cache-Control", "no-transform")
	// force using correct event-stream if there is no proxy
	context.Writer.Header().Set("Content-Type", "text/event-stream")

	context.SSEvent(STREAM.String(), "start")
	defer context.SSEvent(STREAM.String(), "end")
	context.Writer.Flush()

	jobUuid, err := uuid.Parse(context.Param("uuid"))
	job := w.activeJobs[jobUuid]
	if err != nil || job == nil || job.Events == nil {
		context.SSEvent(SERVER.String(), "Logs streaming failed: 404 Not Found")
		return
	}

	context.SSEvent(STATUS.String(), job.Status)
	channel := job.Subscribe()
	context.Stream(func(w io.Writer) bool {
		if event, ok := <-channel; ok {
			context.SSEvent(event.Name.String(), event.Message)
			return true
		} else {
			return false
		}
	})
	job.Unsubscribe(channel)
	context.SSEvent(STATUS.String(), job.Status)
}

func (w worker) runner() {
	for model := range w.pendingJobs {
		model := model
		go func() {
			jobUuid := model.Uuid

			cmd := exec.Command("pgcompacttable", model.CommandArgs...)
			pipe, _ := cmd.StdoutPipe()
			if errStart := cmd.Start(); errStart != nil {
				w.updateStatus(model, FAILED, errStart.Error())
				return
			}

			w.updateStatus(model, RUNNING, "")
			reader := bufio.NewReader(pipe)
			line, _, err := reader.ReadLine()
			file, _ := persistence.File.CompactTable.Open(model.LogsPath)
			for err == nil {
				lineString := string(line)
				_, errFile := file.WriteString(lineString + "\n")
				if errFile != nil {
					w.updateStatus(model, FAILED, errFile.Error())
					return
				}
				w.sendEvents(jobUuid, LOG, lineString)
				line, _, err = reader.ReadLine()
			}

			w.updateStatus(model, FINISHED, "")
			w.closeEvents(jobUuid)
			w.clean(jobUuid)
		}()
	}
}

func (w worker) initiate() {
	runningJobs, _ := persistence.Database.CompactTable.ListByStatus(RUNNING)
	for _, runningJob := range runningJobs {
		_ = persistence.Database.CompactTable.UpdateStatus(runningJob, FAILED)
	}

	pendingJobs, _ := persistence.Database.CompactTable.ListByStatus(PENDING)
	for _, pendingJob := range pendingJobs {
		pendingJob := pendingJob
		go func() { w.pendingJobs <- pendingJob }()
	}
}

func (w worker) clean(uuid uuid.UUID) {
	go func() {
		for {
			if w.activeJobs[uuid].Size() == 0 {
				// TODO FIX fatal error: concurrent map read and map write
				delete(w.activeJobs, uuid)
				return
			}
		}
	}()
}

func (w worker) closeEvents(uuid uuid.UUID) {
	for subscriber := range w.activeJobs[uuid].Events {
		close(subscriber)
	}
}

func (w worker) sendEvents(uuid uuid.UUID, eventType EventType, message string) {
	for subscriber := range w.activeJobs[uuid].Events {
		subscriber <- Event{Name: eventType, Message: message}
	}
}

func (w worker) updateStatus(model CompactTableModel, status JobStatus, message string) {
	dbErr := persistence.Database.CompactTable.UpdateStatus(model, status)
	if dbErr != nil {
		w.sendEvents(model.Uuid, SERVER, dbErr.Error())
	} else {
		w.activeJobs[model.Uuid].Status = status
		if message != "" {
			w.sendEvents(model.Uuid, SERVER, message)
		}
	}
}
