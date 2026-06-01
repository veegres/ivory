package job

import (
	"bufio"
	"fmt"
	"ivory/clients/console"
	"ivory/core/store"
	"sync"
)

type Service struct {
	mu      *sync.Mutex
	jobs    map[JobID]*Job
	storage *store.FileStorage
}

func NewService(storage *store.FileStorage) *Service {
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
func (s *Service) Stream(id JobID, subID SubscriberID, send func(Message)) {
	send(Message{Type: STREAM, Message: START.String()})
	defer send(Message{Type: STREAM, Message: END.String()})

	// 1. Stream from file if it exists
	if s.storage != nil {
		if file, err := s.storage.OpenByName(string(id)); err == nil {
			defer file.Close()
			send(Message{Type: SERVER, Message: "Streaming from the file started"})
			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				send(Message{Type: LOG, Message: scanner.Text()})
			}
			send(Message{Type: SERVER, Message: "Streaming from the file finished"})
		}
	}

	// 2. Subscribe to live job if it's running
	liveChan, err := s.Subscribe(id, subID)
	if err == nil && liveChan != nil {
		send(Message{Type: SERVER, Message: "Streaming from the console started"})
		for event := range liveChan {
			send(event)
		}
		send(Message{Type: SERVER, Message: "Streaming from the console finished"})
	}
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

// Subscribe attaches a subscriber to a running job.
func (s *Service) Subscribe(id JobID, subscriberID SubscriberID) (<-chan Message, error) {
	s.mu.Lock()
	job, ok := s.jobs[id]
	s.mu.Unlock()

	if !ok {
		return nil, fmt.Errorf("job %s not found", id)
	}
	return job.addSubscriber(subscriberID), nil
}

// Unsubscribe detaches a subscriber from a job without stopping it.
func (s *Service) Unsubscribe(id JobID, subscriberID SubscriberID) {
	s.mu.Lock()
	job, ok := s.jobs[id]
	s.mu.Unlock()

	if !ok {
		return
	}
	job.removeSubscriber(subscriberID)
}

func (s *Service) Status(id JobID) Status {
	s.mu.Lock()
	job, ok := s.jobs[id]
	s.mu.Unlock()

	if !ok {
		return UNKNOWN
	}
	return job.getStatus()
}

func (s *Service) GetLogsPath(id JobID) (string, error) {
	if s.storage == nil {
		return "", fmt.Errorf("storage not initialized")
	}
	return s.storage.GetPathByName(string(id))
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
