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
	w := worker{}.Build()

	node := group.Group("/cli")
	node.GET("/bloat", getCompactTableList)
	node.POST("/bloat", w.createCompactTable)
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

// TODO fix deletion for PENDING AND RUNNING jobs
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
	models chan CompactTableModel
	jobs   map[uuid.UUID]*Job
	mutex  *sync.Mutex
}

func (w worker) Build() *worker {
	w.models = make(chan CompactTableModel)
	w.jobs = make(map[uuid.UUID]*Job)
	w.mutex = &sync.Mutex{}

	go w.initializer()
	go w.runner()
	go w.cleaner()
	return &w
}

func (w worker) CreateJob(model CompactTableModel) {
	w.models <- model
	w.mutex.Lock()
	w.jobs[model.Uuid] = Job{}.Create()
	w.mutex.Unlock()
}

func (w worker) DeleteJob(id uuid.UUID) {
	w.mutex.Lock()
	delete(w.jobs, id)
	w.mutex.Unlock()
}

func (w worker) createCompactTable(context *gin.Context) {
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
	w.CreateJob(compactTableModel)

	context.AbortWithStatusJSON(http.StatusOK, gin.H{"response": compactTableModel})
}

// TODO we should stream logs from file if job already running and we missed some logs
func (w worker) stream(context *gin.Context) {
	// notify proxy that it shouldn't enable any caching
	context.Writer.Header().Set("Cache-Control", "no-transform")
	// force using correct event-stream if there is no proxy
	context.Writer.Header().Set("Content-Type", "text/event-stream")

	context.SSEvent(STREAM.String(), "start")
	defer context.SSEvent(STREAM.String(), "end")
	context.Writer.Flush()

	jobUuid, err := uuid.Parse(context.Param("uuid"))
	job := w.jobs[jobUuid]
	if err != nil || job == nil {
		context.SSEvent(SERVER.String(), "Logs streaming failed: 404 Not Found")
		return
	}

	context.SSEvent(STATUS.String(), job.Status())
	channel := job.Subscribe()
	if channel != nil {
		context.Stream(func(w io.Writer) bool {
			if event, ok := <-channel; ok {
				context.SSEvent(event.Name.String(), event.Message)
				return true
			} else {
				return false
			}
		})
		job.Unsubscribe(channel)
	}
	context.SSEvent(STATUS.String(), job.Status())
}

func (w worker) runner() {
	for model := range w.models {
		model := model
		go func() {
			jobUuid := model.Uuid

			cmd := exec.Command("pgcompacttable", model.CommandArgs...)
			pipe, _ := cmd.StdoutPipe()
			if errStart := cmd.Start(); errStart != nil {
				w.updateStatus(model, errStart)
				return
			}

			w.updateStatus(model, nil)
			reader := bufio.NewReader(pipe)
			line, _, err := reader.ReadLine()
			file, _ := persistence.File.CompactTable.Open(model.LogsPath)
			for err == nil {
				lineString := string(line)
				_, errFile := file.WriteString(lineString + "\n")
				if errFile != nil {
					w.updateStatus(model, errFile)
					return
				}
				w.sendEvents(jobUuid, LOG, lineString)
				line, _, err = reader.ReadLine()
			}

			w.updateStatus(model, nil)
			w.closeEvents(jobUuid)
		}()
	}
}

func (w worker) initializer() {
	runningJobs, _ := persistence.Database.CompactTable.ListByStatus(RUNNING)
	for _, runningJob := range runningJobs {
		_ = persistence.Database.CompactTable.UpdateStatus(runningJob, FAILED)
	}

	pendingJobs, _ := persistence.Database.CompactTable.ListByStatus(PENDING)
	for _, pendingJob := range pendingJobs {
		pendingJob := pendingJob
		go func() { w.models <- pendingJob }()
	}
}

func (w worker) cleaner() {
	ticket := time.NewTicker(10 * time.Second)
	for range ticket.C {
		for id, job := range w.jobs {
			if job.IsFinished() {
				w.DeleteJob(id)
			}
		}
	}
}

func (w worker) closeEvents(uuid uuid.UUID) {
	for subscriber := range w.jobs[uuid].Subscribers() {
		close(subscriber)
	}
}

func (w worker) sendEvents(uuid uuid.UUID, eventType EventType, message string) {
	for subscriber := range w.jobs[uuid].Subscribers() {
		subscriber <- Event{Name: eventType, Message: message}
	}
}

func (w worker) updateStatus(model CompactTableModel, err error) {
	job := w.jobs[model.Uuid]
	dbErr := persistence.Database.CompactTable.UpdateStatus(model, job.Next())
	if dbErr != nil {
		job.Failed()
		w.sendEvents(model.Uuid, SERVER, dbErr.Error())
	} else {
		if err != nil {
			job.Failed()
			w.sendEvents(model.Uuid, SERVER, err.Error())
		}
	}
}
