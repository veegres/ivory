package job

import (
	"fmt"
	"io"
	"ivory/storage/files"
	"maps"
	"os"
	"slices"
	"strings"
	"sync"
	"time"
)

type Job struct {
	cmd         Command
	subscribers map[SubscriberID]chan Message
	status      Status
	mu          sync.RWMutex
	storage     *files.Storage
}

func NewJob(cmd Command, storage *files.Storage) *Job {
	return &Job{
		subscribers: make(map[SubscriberID]chan Message),
		cmd:         cmd,
		status:      PENDING,
		storage:     storage,
	}
}

func (j *Job) Subscribers() []SubscriberID {
	j.mu.RLock()
	defer j.mu.RUnlock()
	return slices.Collect(maps.Keys(j.subscribers))
}

func (j *Job) Size() int {
	j.mu.RLock()
	defer j.mu.RUnlock()
	return len(j.subscribers)
}

// Run starts the command and streams output to all subscribers.
// Intended to be called in its own goroutine by the job.Manager.
func (j *Job) Run() {
	j.setStatus(RUNNING)
	defer j.close()

	var logFile *os.File
	if j.cmd.Persist() && j.storage != nil {
		var err error
		logFile, err = j.storage.OpenOrCreateByName(j.cmd.Id())
		if err != nil {
			j.setStatus(FAILED)
			j.broadcast(SERVER, fmt.Sprintf("failed to create log file: %s", err))
			return
		}
		defer logFile.Close()
	}

	reader, errStart := j.cmd.Start()
	if errStart != nil {
		j.setStatus(FAILED)
		j.broadcast(SERVER, errStart.Error())
		return
	}

	buf := make([]byte, 4096)
	for {
		n, errRead := reader.Read(buf)
		if n > 0 {
			msg := strings.TrimSuffix(string(buf[:n]), "\n")
			j.broadcast(LOG, msg)
			if logFile != nil {
				_, _ = logFile.WriteString(msg + "\n")
			}
		}
		if errRead != nil {
			errWait := j.cmd.Wait()
			if j.getStatus() == STOPPED {
				return
			}
			if errRead != io.EOF {
				j.setStatus(FAILED)
				j.broadcast(SERVER, errRead.Error())
				return
			}
			if errWait != nil {
				j.setStatus(FAILED)
				j.broadcast(SERVER, errWait.Error())
				return
			}
			j.setStatus(FINISHED)
			return
		}
	}
}

// Stop marks the job as STOPPED before calling Abort() so that when
// Read() unblocks with an error in Run(), the status check can
// distinguish an intentional stop from a real failure.
func (j *Job) Stop() error {
	j.setStatus(STOPPED)
	return j.cmd.Abort()
}

func (j *Job) broadcast(t EventType, m string) {
	j.mu.RLock()
	defer j.mu.RUnlock()
	for _, ch := range j.subscribers {
		select {
		case ch <- Message{Type: t, Timestamp: time.Now(), Message: m}:
		default: // slow subscriber: drop rather than block broadcast
		}
	}
}

func (j *Job) close() {
	j.mu.Lock()
	defer j.mu.Unlock()
	for id, ch := range j.subscribers {
		close(ch)
		delete(j.subscribers, id)
	}
}

func (j *Job) addSubscriber(id SubscriberID) chan Message {
	var channel chan Message
	j.mu.Lock()
	if j.status == PENDING || j.status == RUNNING {
		channel = make(chan Message, 256)
		j.subscribers[id] = channel
	}
	j.mu.Unlock()
	return channel
}

func (j *Job) removeSubscriber(id SubscriberID) {
	j.mu.Lock()
	if ch, ok := j.subscribers[id]; ok {
		close(ch)
		delete(j.subscribers, id)
	}
	j.mu.Unlock()
}

func (j *Job) getStatus() Status {
	j.mu.RLock()
	defer j.mu.RUnlock()
	return j.status
}

func (j *Job) setStatus(status Status) {
	j.mu.Lock()
	j.status = status
	j.mu.Unlock()
	j.broadcast(STATUS, status.String())
}
