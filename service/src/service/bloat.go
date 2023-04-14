package service

import (
	"bufio"
	"errors"
	"github.com/google/uuid"
	. "ivory/src/model"
	"ivory/src/persistence"
	"os"
	"os/exec"
	"sync"
	"time"
)

type BloatService struct {
	start           chan uuid.UUID
	stop            chan uuid.UUID
	elements        map[uuid.UUID]*element
	mutex           *sync.Mutex
	bloatRepository *persistence.BloatRepository
	passwordService *PasswordService
}

type element struct {
	job   Job
	model *BloatModel
	file  *os.File
}

func NewBloatService(
	bloatRepository *persistence.BloatRepository,
	passwordService *PasswordService,
) *BloatService {
	worker := &BloatService{
		start:           make(chan uuid.UUID),
		stop:            make(chan uuid.UUID),
		elements:        make(map[uuid.UUID]*element),
		mutex:           &sync.Mutex{},
		bloatRepository: bloatRepository,
		passwordService: passwordService,
	}

	// run channel subscribers
	go worker.initializer()
	go worker.runner()
	go worker.cleaner()
	go worker.stopper()
	return worker
}

func (w *BloatService) Start(credentialId uuid.UUID, cluster string, args []string) (*BloatModel, error) {
	compactTable, err := w.bloatRepository.Create(credentialId, cluster, args)
	if err != nil {
		return nil, err
	}
	w.addElement(compactTable)
	return compactTable, nil
}

func (w *BloatService) Stream(jobUuid uuid.UUID, stream func(event Event)) {
	element := w.elements[jobUuid]
	if element == nil {
		stream(Event{Type: SERVER, Message: "Logs streaming failed: Stream Not Found"})
		stream(Event{Type: STATUS, Message: UNKNOWN.String()})
		return
	}
	job := element.job
	model := element.model

	// subscribe to stream and show status
	// we have to subscribe as soon as possible to prevent Job deletion
	channel := job.Subscribe()
	stream(Event{Type: STATUS, Message: job.GetStatus().String()})

	// stream logs from file
	file, _ := w.bloatRepository.GetFile(model.LogsPath)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		stream(Event{Type: LOG, Message: scanner.Text()})
	}
	if errScanner := scanner.Err(); errScanner != nil {
		stream(Event{Type: LOG, Message: errScanner.Error()})
	}
	errFileClose := file.Close()
	if errFileClose != nil {
		stream(Event{Type: LOG, Message: errFileClose.Error()})
	}

	// subscribe to stream and stream logs
	if channel != nil {
		for event := range channel {
			stream(event)
		}
		job.Unsubscribe(channel)
	}

	stream(Event{Type: STATUS, Message: job.GetStatus().String()})
}

func (w *BloatService) Delete(jobUuid uuid.UUID) error {
	element := w.elements[jobUuid]
	if element != nil && element.job.IsJobActive() {
		return errors.New("job is active")
	}
	return w.bloatRepository.Delete(jobUuid)
}

func (w *BloatService) DeleteAll() error {
	for _, e := range w.elements {
		e.job.SetStatus(FAILED)
		for s, _ := range e.job.subscribers {
			e.job.Unsubscribe(s)
		}
	}
	return w.bloatRepository.DeleteAll()
}

func (w *BloatService) Stop(jobUuid uuid.UUID) error {
	if w.elements[jobUuid] != nil {
		w.stop <- jobUuid
		return nil
	} else {
		return errors.New("there is no such active job")
	}
}

func (w *BloatService) runner() {
	for id := range w.start {
		element := w.elements[id]
		model := w.elements[id].model
		jobUuid := model.Uuid
		go func() {
			w.jobStatusHandler(element, RUNNING, nil)

			// Get password
			credential, errCred := w.passwordService.GetDecrypted(model.CredentialId)
			if errCred != nil {
				w.jobStatusHandler(element, FAILED, errors.New("Credential error: "+errCred.Error()))
				return
			}
			credentialArgs := []string{
				"--user", credential.Username,
				"--password", credential.Password,
			}

			// Run command
			args := append(model.CommandArgs, credentialArgs...)
			cmd := exec.Command("pgcompacttable", args...)
			pipe, errPipe := cmd.StdoutPipe()
			if errPipe != nil {
				w.jobStatusHandler(element, FAILED, errors.New("pgcompacttable execution error: "+errPipe.Error()))
				return
			}
			if errStart := cmd.Start(); errStart != nil {
				w.jobStatusHandler(element, FAILED, errors.New("pgcompacttable execution error: "+errStart.Error()))
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
				w.jobStatusHandler(element, FAILED, errors.New("pgcompacttable execution error: "+errWait.Error()))
				return
			}

			w.jobStatusHandler(element, FINISHED, nil)
		}()
	}
}

func (w *BloatService) initializer() {
	runningJobs, _ := w.bloatRepository.ListByStatus(RUNNING)
	for _, runningJob := range runningJobs {
		_ = w.bloatRepository.UpdateStatus(runningJob, FAILED)
	}

	pendingJobs, _ := w.bloatRepository.ListByStatus(PENDING)
	for _, pendingJob := range pendingJobs {
		pendingJob := pendingJob
		go func() { w.addElement(&pendingJob) }()
	}
}

func (w *BloatService) cleaner() {
	ticket := time.NewTicker(10 * time.Second)
	for range ticket.C {
		for id, element := range w.elements {
			if !element.job.IsJobActive() && element.job.Size() == 0 {
				w.removeElement(id)
			}
		}
	}
}

func (w *BloatService) stopper() {
	for id := range w.stop {
		id := id
		go func() {
			element := w.elements[id]
			if element == nil {
				return
			}
			job := element.job
			if cmd := job.GetCommand(); cmd != nil {
				killErr := job.GetCommand().Process.Kill()
				if killErr == nil {
					w.jobStatusHandler(element, STOPPED, nil)
				} else {
					if job.IsJobActive() {
						w.jobStatusHandler(element, FINISHED, errors.New("Kill pgcompacttable error: "+killErr.Error()))
					} else {
						w.jobStatusHandler(element, job.GetStatus(), errors.New("Kill pgcompacttable error: "+killErr.Error()))
					}
				}
			}
		}()
	}
}

func (w *BloatService) jobStatusHandler(element *element, status JobStatus, err error) {
	model := element.model
	element.job.SetStatus(status)
	dbErr := w.bloatRepository.UpdateStatus(*model, status)
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

func (w *BloatService) addLogElement(element *element, eventType EventType, message string) {
	_, errFileWrite := element.file.WriteString(message + "\n")
	if errFileWrite != nil {
		w.sendEvents(element.job, SERVER, errFileWrite.Error())
		return
	}
	w.sendEvents(element.job, eventType, message)
}

func (w *BloatService) sendEvents(job Job, eventType EventType, message string) {
	for subscriber := range job.Subscribers() {
		subscriber <- Event{Type: eventType, Message: message}
	}
}

func (w *BloatService) closeEvents(job Job) {
	for subscriber := range job.Subscribers() {
		close(subscriber)
	}
}

func (w *BloatService) addElement(model *BloatModel) {
	w.mutex.Lock()
	file, _ := w.bloatRepository.GetFile(model.LogsPath)
	w.elements[model.Uuid] = &element{job: *NewJob(), model: model, file: file}
	w.mutex.Unlock()
	w.start <- model.Uuid
}

func (w *BloatService) removeElement(id uuid.UUID) {
	w.mutex.Lock()
	delete(w.elements, id)
	w.mutex.Unlock()
}
