package service

import (
	"bufio"
	"errors"
	"github.com/google/uuid"
	. "ivory/model"
	"ivory/persistence"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"
)

type JobWorker interface {
	Start(credentialId uuid.UUID, cluster string, args []string) (*CompactTableModel, error)
	Stream(jobUuid uuid.UUID, stream func(event Event))
	Delete(jobUuid uuid.UUID) error
	Stop(jobUuid uuid.UUID)
}

type worker struct {
	start    chan uuid.UUID
	stop     chan uuid.UUID
	elements map[uuid.UUID]*element
	mutex    *sync.Mutex
}

type element struct {
	job   Job
	model *CompactTableModel
	file  *os.File
}

func CreateJobWorker() JobWorker {
	worker := &worker{
		start:    make(chan uuid.UUID),
		stop:     make(chan uuid.UUID),
		elements: make(map[uuid.UUID]*element),
		mutex:    &sync.Mutex{},
	}

	// run channel subscribers
	go worker.initializer()
	go worker.runner()
	go worker.cleaner()
	go worker.stopper()
	return worker
}

func (w *worker) Start(credentialId uuid.UUID, cluster string, args []string) (*CompactTableModel, error) {
	jobUuid := uuid.New()
	compactTableModel := CompactTableModel{
		Uuid:         jobUuid,
		CredentialId: credentialId,
		Cluster:      cluster,
		Status:       PENDING,
		Command:      "pgcompacttable " + strings.Join(args, " "),
		CommandArgs:  args,
		LogsPath:     persistence.File.CompactTable.Create(jobUuid),
		CreatedAt:    time.Now().UnixNano(),
	}
	errCompactTable := persistence.Database.CompactTable.Update(compactTableModel)
	if errCompactTable != nil {
		return nil, errCompactTable
	}
	w.addElement(compactTableModel)
	return &compactTableModel, nil
}

func (w *worker) Stream(jobUuid uuid.UUID, stream func(event Event)) {
	element := w.elements[jobUuid]
	if element == nil {
		stream(Event{Name: SERVER, Message: "Logs streaming failed: Stream Not Found"})
		model, modelErr := persistence.Database.CompactTable.Get(jobUuid)
		if modelErr != nil {
			stream(Event{Name: STATUS, Message: UNKNOWN.String()})
		} else {
			stream(Event{Name: STATUS, Message: model.Status.String()})
		}
		return
	}
	job := element.job
	model := element.model

	// subscribe to stream and show status
	// we have to subscribe as soon as possible to prevent Job deletion
	channel := job.Subscribe()
	stream(Event{Name: STATUS, Message: job.GetStatus().String()})

	// stream logs from file
	file, _ := persistence.File.CompactTable.Open(model.LogsPath)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		stream(Event{Name: LOG, Message: scanner.Text()})
	}
	if errScanner := scanner.Err(); errScanner != nil {
		stream(Event{Name: LOG, Message: errScanner.Error()})
	}
	errFileClose := file.Close()
	if errFileClose != nil {
		stream(Event{Name: LOG, Message: errFileClose.Error()})
	}

	// subscribe to stream and stream logs
	if channel != nil {
		for {
			if event, ok := <-channel; ok {
				stream(event)
			} else {
				break
			}
		}
		job.Unsubscribe(channel)
	}

	stream(Event{Name: STATUS, Message: job.GetStatus().String()})
}

func (w *worker) Delete(jobUuid uuid.UUID) error {
	element := w.elements[jobUuid]
	if element != nil && element.job.IsJobActive() {
		return errors.New("job is active")
	}
	errFile := persistence.File.CompactTable.Delete(jobUuid)
	if errFile != nil {
		return errFile
	}
	errDb := persistence.Database.CompactTable.Delete(jobUuid)
	if errDb != nil {
		return errDb
	}
	return nil
}

func (w *worker) Stop(jobUuid uuid.UUID) {
	w.stop <- jobUuid
}

func (w *worker) runner() {
	for id := range w.start {
		element := w.elements[id]
		model := w.elements[id].model
		jobUuid := model.Uuid
		go func() {
			w.jobStatusHandler(element, RUNNING, nil)

			// Get password
			// TODO move decryption to another method (encapsulate)
			credential, _ := persistence.Database.Credential.GetCredential(model.CredentialId)
			password, errDecrypt := Decrypt(credential.Password, Secret.Get())
			if errDecrypt != nil {
				w.jobStatusHandler(element, FAILED, errDecrypt)
				return
			}
			credentialArgs := []string{
				"--user", credential.Username,
				"--password", password,
			}

			// Run command
			args := append(model.CommandArgs, credentialArgs...)
			cmd := exec.Command("pgcompacttable", args...)
			pipe, errPipe := cmd.StdoutPipe()
			if errPipe != nil {
				w.jobStatusHandler(element, FAILED, errPipe)
				return
			}
			if errStart := cmd.Start(); errStart != nil {
				w.jobStatusHandler(element, FAILED, errStart)
				return
			}
			w.elements[jobUuid].job.SetCommand(cmd)

			// Read and write logs from command
			reader := bufio.NewReader(pipe)
			line, _, errReadLine := reader.ReadLine()
			for errReadLine == nil {
				lineString := string(line)
				w.addLogElement(element, LOG, lineString)
				line, _, errReadLine = reader.ReadLine()
			}

			// Wait when pipe will be closed
			if errWait := cmd.Wait(); errWait != nil {
				w.jobStatusHandler(element, FAILED, errWait)
				return
			}

			w.jobStatusHandler(element, FINISHED, nil)
		}()
	}
}

func (w *worker) initializer() {
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

func (w *worker) cleaner() {
	ticket := time.NewTicker(10 * time.Second)
	for range ticket.C {
		for id, element := range w.elements {
			if !element.job.IsJobActive() && element.job.Size() == 0 {
				w.removeElement(id)
			}
		}
	}
}

func (w *worker) stopper() {
	for id := range w.stop {
		id := id
		go func() {
			element := w.elements[id]
			job := element.job
			if cmd := job.GetCommand(); cmd != nil {
				err := job.GetCommand().Process.Kill()
				if err == nil {
					w.jobStatusHandler(element, STOPPED, nil)
				} else {
					if job.IsJobActive() {
						w.jobStatusHandler(element, FINISHED, err)
					} else {
						w.jobStatusHandler(element, job.GetStatus(), err)
					}
				}
			}
		}()
	}
}

func (w *worker) jobStatusHandler(element *element, status JobStatus, err error) {
	model := element.model
	element.job.SetStatus(status)
	dbErr := persistence.Database.CompactTable.UpdateStatus(*model, status)
	if dbErr != nil {
		w.addLogElement(element, SERVER, dbErr.Error())
	}
	if err != nil {
		w.addLogElement(element, SERVER, err.Error())
	}
	if !element.job.IsJobActive() {
		w.closeEvents(element.job)
	}
}

func (w *worker) addLogElement(element *element, eventType EventType, message string) {
	_, errFileWrite := element.file.WriteString(message + "\n")
	if errFileWrite != nil {
		w.sendEvents(element.job, SERVER, errFileWrite.Error())
		return
	}
	w.sendEvents(element.job, eventType, message)
}

func (w *worker) sendEvents(job Job, eventType EventType, message string) {
	for subscriber := range job.Subscribers() {
		subscriber <- Event{Name: eventType, Message: message}
	}
}

func (w *worker) closeEvents(job Job) {
	for subscriber := range job.Subscribers() {
		close(subscriber)
	}
}

func (w *worker) addElement(model CompactTableModel) {
	w.mutex.Lock()
	file, _ := persistence.File.CompactTable.Open(model.LogsPath)
	w.elements[model.Uuid] = &element{job: CreateJob(), model: &model, file: file}
	w.mutex.Unlock()
	w.start <- model.Uuid
}

func (w *worker) removeElement(id uuid.UUID) {
	w.mutex.Lock()
	delete(w.elements, id)
	w.mutex.Unlock()
}
