package bloat

import (
	"bufio"
	"errors"
	. "ivory/src/features/bloat/job"
	"ivory/src/features/password"
	"os"
	"os/exec"
	"sync"
	"time"

	"github.com/google/uuid"
)

type BloatService struct {
	start           chan uuid.UUID
	stop            chan uuid.UUID
	elements        map[uuid.UUID]*element
	mutex           *sync.Mutex
	bloatRepository *BloatRepository
	passwordService *password.PasswordService
}

type element struct {
	job    *Job
	model  *Bloat
	writer *os.File
}

func NewBloatService(
	bloatRepository *BloatRepository,
	passwordService *password.PasswordService,
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

func (w *BloatService) List() ([]Bloat, error) {
	return w.bloatRepository.List()
}

func (w *BloatService) ListByStatus(status JobStatus) ([]Bloat, error) {
	return w.bloatRepository.ListByStatus(status)
}

func (w *BloatService) ListByCluster(cluster string) ([]Bloat, error) {
	return w.bloatRepository.ListByCluster(cluster)
}

func (w *BloatService) Get(uuid uuid.UUID) (Bloat, error) {
	return w.bloatRepository.Get(uuid)
}

func (w *BloatService) Start(credentialId uuid.UUID, cluster string, args []string) (*Bloat, error) {
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
		stream(Event{Type: SERVER, Message: "Streaming failed: Stream Not Found"})
		stream(Event{Type: STATUS, Message: UNKNOWN.String()})
		return
	}
	job := element.job
	model := element.model

	// subscribing to stream and showing status,
	// we have to subscribe as soon as possible to prevent Job deletion
	channel := job.Subscribe()
	stream(Event{Type: STATUS, Message: job.GetStatus().String()})

	// stream logs from file
	// NOTE: we need to open the file one more time, because a writer has a cursor at the end
	reader, errFile := w.bloatRepository.GetOpenFile(model.Uuid)
	if errFile != nil {
		stream(Event{Type: SERVER, Message: errFile.Error()})
	} else {
		stream(Event{Type: SERVER, Message: "Streaming from the file started"})
		scanner := bufio.NewScanner(reader)
		for scanner.Scan() {
			stream(Event{Type: LOG, Message: scanner.Text()})
		}
		if errScanner := scanner.Err(); errScanner != nil {
			stream(Event{Type: SERVER, Message: errScanner.Error()})
		}
		errFileClose := reader.Close()
		if errFileClose != nil {
			stream(Event{Type: SERVER, Message: errFileClose.Error()})
		}
		stream(Event{Type: SERVER, Message: "Streaming from the file finished"})
	}

	// subscribe to stream and stream logs
	if channel != nil {
		stream(Event{Type: SERVER, Message: "Streaming from the console started"})
		for event := range channel {
			stream(event)
		}
		stream(Event{Type: SERVER, Message: "Streaming from the console finished"})

		// NOTE: we should get status before unsubscription
		stream(Event{Type: STATUS, Message: job.GetStatus().String()})
		job.Unsubscribe(channel)
	} else {
		stream(Event{Type: STATUS, Message: job.GetStatus().String()})
	}
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
		for s := range e.job.Subscribers() {
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
			defer func() { _ = element.writer.Close() }()
			w.jobStatusHandler(element, RUNNING, nil)

			// Get password
			credential, errCred := w.passwordService.GetDecrypted(model.CredentialId)
			if errCred != nil {
				w.jobStatusHandler(element, FAILED, errors.New("Password error: "+errCred.Error()))
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

			// read and write logs from command
			reader := bufio.NewReader(pipe)
			line, _, errReadLine := reader.ReadLine()
			for errReadLine == nil {
				lineString := string(line)
				w.addLogElement(element, LOG, lineString)
				line, _, errReadLine = reader.ReadLine()
			}

			// Wait when pipe will be closed
			if errWait := cmd.Wait(); errWait != nil && element.job.IsJobActive() {
				w.jobStatusHandler(element, FAILED, errors.New("pgcompacttable execution error: "+errWait.Error()))
				return
			}

			if element.job.IsJobActive() {
				w.jobStatusHandler(element, FINISHED, nil)
			}
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
	_, errFileWrite := element.writer.WriteString(message + "\n")
	if errFileWrite != nil {
		w.sendEvents(element.job, SERVER, errFileWrite.Error())
		return
	}
	w.sendEvents(element.job, eventType, message)
}

func (w *BloatService) sendEvents(job *Job, eventType EventType, message string) {
	for subscriber := range job.Subscribers() {
		subscriber <- Event{Type: eventType, Message: message}
	}
}

func (w *BloatService) closeEvents(job *Job) {
	for subscriber := range job.Subscribers() {
		close(subscriber)
	}
}

func (w *BloatService) addElement(model *Bloat) {
	w.mutex.Lock()
	// NOTE: we can have potential problem here because file can be nil
	file, _ := w.bloatRepository.GetOpenFile(model.Uuid)
	w.elements[model.Uuid] = &element{job: NewJob(), model: model, writer: file}
	w.mutex.Unlock()
	w.start <- model.Uuid
}

func (w *BloatService) removeElement(id uuid.UUID) {
	w.mutex.Lock()
	delete(w.elements, id)
	w.mutex.Unlock()
}
