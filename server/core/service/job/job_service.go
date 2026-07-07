package job

import (
	"bufio"
	"errors"
	"fmt"
	"ivory/clients/console"
	"ivory/clients/storage"
	"log/slog"
	"sync"
)

var ErrStorageNotInitialized = errors.New("storage not initialized")

type Service struct {
	mu      *sync.Mutex
	jobs    map[JobID]*Job
	storage *storage.FileStorage
}

func NewService(storage *storage.FileStorage) *Service {
	return &Service{
		mu:      &sync.Mutex{},
		jobs:    make(map[JobID]*Job),
		storage: storage,
	}
}

// Start creates a job for the given command and begins execution.
// If a job for this command is already active, the existing JobID
// is returned without starting a second instance.
func (s *Service) Start(cmd console.Command) (JobID, error) {
	jobID := JobID(cmd.Id())
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.jobs[jobID]; ok {
		return jobID, nil
	}

	job := NewJob(cmd, s.storage)
	s.unsafeAddJob(jobID, job)

	go func() {
		job.Run()
		s.removeJob(jobID)
	}()

	return jobID, nil
}

// Stream orchestrates the streaming of job events.
// It handles START/END events, historical logs from file, and live updates.
func (s *Service) Stream(id JobID, subID SubscriberID, close <-chan struct{}, send func(Message)) {
	send(Message{Type: STREAM, Message: START.String()})
	defer send(Message{Type: STREAM, Message: END.String()})

	// 1. Stream from file if it exists
	if s.storage != nil {
		send(Message{Type: SERVER, Message: "streaming from the file, tries to established connection"})
		file, err := s.storage.OpenByName(string(id))
		if err != nil {
			send(Message{Type: SERVER, Message: "streaming from the file skipped, error: " + err.Error()})
		} else {
			defer func() {
				if err := file.Close(); err != nil {
					slog.Error("failed to close file", "error", err)
				}
			}()
			send(Message{Type: SERVER, Message: "streaming from the file started"})
			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				send(Message{Type: LOG, Message: scanner.Text()})
			}
			send(Message{Type: SERVER, Message: "streaming from the file finished"})
		}
	}

	// NOTE: Race Condition: Streaming Gap --- it is expected behavior
	//  There is a temporal gap finishing the read of the log file and subscribing to the live job. Logs
	//  produced by the command in this tiny window may be neither in the file yet nor captured by the
	//  new subscription, leading to missing log lines for the user.

	// 2. Subscribe to live job if it's running
	sub, err := s.Subscribe(id, subID)
	if err != nil {
		send(Message{Type: SERVER, Message: "streaming from the console error: " + err.Error()})
		return
	}
	defer sub.Close()
	send(Message{Type: SERVER, Message: "streaming from the console started"})
loop:
	for {
		select {
		case <-close:
			break loop
		case event, ok := <-sub.Messages:
			if !ok {
				break loop
			}
			send(event)
		}
	}
	send(Message{Type: SERVER, Message: "streaming from the console finished"})
}

// Stop cancels the job with the given ID.
func (s *Service) Stop(id JobID) error {
	s.mu.Lock()
	job, ok := s.jobs[id]
	s.mu.Unlock()

	if !ok {
		return fmt.Errorf("job %s not found", id)
	}
	return job.Stop()
}

// Subscribe attaches a subscriber to a running job. Call sub.Close() (typically
// via defer) once the caller is done reading from sub.Messages.
func (s *Service) Subscribe(id JobID, subscriberID SubscriberID) (*Subscription, error) {
	job, ok := s.getJob(id)
	if !ok {
		return nil, fmt.Errorf("job %s not found", id)
	}
	sub, ok := job.addSubscriber(subscriberID)
	if !ok {
		return nil, fmt.Errorf("job %s is not running", id)
	}
	return sub, nil
}

func (s *Service) Status(id JobID) Status {
	job, ok := s.getJob(id)
	if !ok {
		return UNKNOWN
	}
	return job.getStatus()
}

func (s *Service) GetLogsPath(id JobID) (string, error) {
	if s.storage == nil {
		return "", ErrStorageNotInitialized
	}
	return s.storage.GetPathByName(string(id))
}

func (s *Service) getJob(id JobID) (*Job, bool) {
	s.mu.Lock()
	job, ok := s.jobs[id]
	s.mu.Unlock()
	return job, ok
}

func (s *Service) unsafeAddJob(id JobID, job *Job) {
	s.jobs[id] = job
}

func (s *Service) unsafeRemoveJob(id JobID) {
	delete(s.jobs, id)
}

func (s *Service) removeJob(id JobID) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.unsafeRemoveJob(id)
}
