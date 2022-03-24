package router

import (
	"bufio"
	"errors"
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
	node.GET("/bloat/:uuid", getCompactTable)
	node.GET("/bloat/:uuid/logs", getCompactTableLogs)

	node.POST("/bloat/job/start", w.StartJob)
	node.POST("/bloat/job/:uuid/stop", w.StopJob)
	node.DELETE("/bloat/job/:uuid/delete", w.DeleteJob)
	node.GET("/bloat/job/:uuid/stream", w.StreamJob)
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

type worker struct {
	start    chan uuid.UUID
	stop     chan uuid.UUID
	elements map[uuid.UUID]*element
	mutex    *sync.Mutex
}

func (w worker) Build() *worker {
	w.start = make(chan uuid.UUID)
	w.stop = make(chan uuid.UUID)
	w.elements = make(map[uuid.UUID]*element)
	w.mutex = &sync.Mutex{}

	go w.initializer()
	go w.runner()
	go w.cleaner()
	go w.stopper()
	return &w
}

func (w worker) StartJob(context *gin.Context) {
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
	w.addElement(compactTableModel)

	context.AbortWithStatusJSON(http.StatusOK, gin.H{"response": compactTableModel})
}

// StreamJob TODO we should stream logs from file if job already running and we missed some logs
func (w worker) StreamJob(context *gin.Context) {
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
	job := w.elements[jobUuid].job
	if err != nil || job == nil {
		context.SSEvent(SERVER.String(), "Logs streaming failed: 404 Not Found")
		return
	}

	// subscribe to stream and stream logs
	context.SSEvent(STATUS.String(), job.GetStatus())
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
	context.SSEvent(STATUS.String(), job.GetStatus())
}

func (w worker) DeleteJob(context *gin.Context) {
	jobUuid, _ := uuid.Parse(context.Param("uuid"))

	element := w.elements[jobUuid]
	if element != nil && element.job.IsJobActive() {
		_ = context.AbortWithError(http.StatusForbidden, errors.New("job is active"))
		return
	}

	errFile := persistence.File.CompactTable.Delete(jobUuid)
	if errFile != nil {
		_ = context.AbortWithError(http.StatusBadRequest, errFile)
	}
	errDb := persistence.Database.CompactTable.Delete(jobUuid)
	if errDb != nil {
		_ = context.AbortWithError(http.StatusBadRequest, errDb)
	}
}

func (w worker) StopJob(context *gin.Context) {
	jobUuid, _ := uuid.Parse(context.Param("uuid"))
	w.stop <- jobUuid
}

func (w worker) runner() {
	for id := range w.start {
		model := w.elements[id].model
		go func() {
			jobUuid := model.Uuid

			// Open file to write command logs
			file, errFileOpen := persistence.File.CompactTable.Open(model.LogsPath)
			if errFileOpen != nil {
				w.jobStatusHandler(model, FAILED, true, errFileOpen)
				return
			}

			// Run command
			cmd := exec.Command("pgcompacttable", model.CommandArgs...)
			pipe, errPipe := cmd.StdoutPipe()
			if errPipe != nil {
				w.jobStatusHandler(model, FAILED, true, errPipe)
				return
			}
			if errStart := cmd.Start(); errStart != nil {
				w.jobStatusHandler(model, FAILED, true, errStart)
				return
			}
			w.elements[jobUuid].job.SetCommand(cmd)

			// Read and write logs from command
			w.jobStatusHandler(model, RUNNING, false, nil)
			reader := bufio.NewReader(pipe)
			line, _, errReadLine := reader.ReadLine()
			for errReadLine == nil {
				lineString := string(line)
				_, errFileWrite := file.WriteString(lineString + "\n")
				if errFileWrite != nil {
					w.jobStatusHandler(model, FAILED, true, errFileWrite)
					return
				}
				w.sendEvents(jobUuid, LOG, lineString)
				line, _, errReadLine = reader.ReadLine()
			}

			// Wait when pipe will be closed
			if errWait := cmd.Wait(); errWait != nil {
				w.jobStatusHandler(model, FAILED, true, errWait)
				return
			}
			w.jobStatusHandler(model, FINISHED, true, nil)
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
		go func() { w.addElement(pendingJob) }()
	}
}

func (w worker) cleaner() {
	ticket := time.NewTicker(10 * time.Second)
	for range ticket.C {
		for id, element := range w.elements {
			if !element.job.IsJobActive() {
				w.removeElement(id)
			}
		}
	}
}

func (w worker) stopper() {
	for id := range w.stop {
		id := id
		go func() {
			job := w.elements[id].job
			if cmd := job.GetCommand(); cmd != nil {
				err := job.GetCommand().Process.Kill()
				if err == nil {
					w.jobStatusHandler(w.elements[id].model, STOPPED, true, nil)
				} else {
					if job.IsJobActive() {
						w.jobStatusHandler(w.elements[id].model, FINISHED, true, err)
					} else {
						w.jobStatusHandler(w.elements[id].model, job.GetStatus(), true, err)
					}
				}
			}
		}()
	}
}

func (w worker) jobStatusHandler(model *CompactTableModel, status JobStatus, close bool, err error) {
	element := w.elements[model.Uuid]
	if element != nil {
		dbErr := persistence.Database.CompactTable.UpdateStatus(*model, status)
		if dbErr != nil {
			w.sendEvents(model.Uuid, SERVER, dbErr.Error())
		}
		if err != nil {
			w.sendEvents(model.Uuid, SERVER, err.Error())
		}

		element.job.SetStatus(status)

		if close == true {
			w.closeEvents(model.Uuid)
		}
	}
}

func (w worker) sendEvents(id uuid.UUID, eventType EventType, message string) {
	element := w.elements[id]
	if element != nil {
		for subscriber := range element.job.Subscribers() {
			subscriber <- Event{Name: eventType, Message: message}
		}
	}
}

func (w worker) closeEvents(id uuid.UUID) {
	element := w.elements[id]
	if element != nil {
		for subscriber := range element.job.Subscribers() {
			// TODO we can close it earlier then it will finish streaming
			close(subscriber)
		}
	}
}

type element struct {
	job   *Job
	model *CompactTableModel
}

func (w worker) addElement(model CompactTableModel) {
	w.mutex.Lock()
	w.elements[model.Uuid] = &element{job: Job{}.Create(), model: &model}
	w.mutex.Unlock()
	w.start <- model.Uuid
}

func (w worker) removeElement(id uuid.UUID) {
	w.mutex.Lock()
	delete(w.elements, id)
	w.mutex.Unlock()
}
