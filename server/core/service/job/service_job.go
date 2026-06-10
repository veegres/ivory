package job

import (
	"bufio"
	"fmt"
	"ivory/clients/console"
	"ivory/core/store"
	"maps"
	"os"
	"slices"
	"sync"
	"time"
)

type Job struct {
	cmd         console.Command
	subscribers map[SubscriberID]chan Message
	status      Status
	mu          sync.RWMutex
	storage     *store.FileStorage

	keepAliveDuration time.Duration
	keepAliveBegin    time.Time
}

func NewJob(cmd console.Command, storage *store.FileStorage) *Job {
	return &Job{
		subscribers:       make(map[SubscriberID]chan Message),
		cmd:               cmd,
		status:            PENDING,
		storage:           storage,
		keepAliveDuration: time.Second * 5,
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
// Intended to be called in its own goroutine by the job.Service.
func (j *Job) Run() {
	j.setStatus(RUNNING)
	j.broadcast(SERVER, fmt.Sprintf("running the job: %s", j.cmd.Id()))
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

	// NOTE: run killer which will kill stale jobs
	go j.killer()

	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		msg := scanner.Text()
		j.broadcast(LOG, msg)
		if logFile != nil {
			_, _ = logFile.WriteString(msg + "\n")
		}
	}

	if j.getStatus() == STOPPED {
		return
	}

	if err := scanner.Err(); err != nil {
		j.setStatus(FAILED)
		j.broadcast(SERVER, err.Error())
		return
	}

	if err := j.cmd.Wait(); err != nil {
		j.setStatus(FAILED)
		j.broadcast(SERVER, err.Error())
		return
	}

	j.broadcast(SERVER, fmt.Sprintf("finished the job: %s", j.cmd.Id()))
	j.setStatus(FINISHED)
}

// Stop marks the job as STOPPED before calling Abort() so that when
// Read() unblocks with an error in Run(), the status check can
// distinguish an intentional stop from a real failure.
func (j *Job) Stop() error {
	j.setStatus(STOPPED)
	return j.cmd.Abort()
}

func (j *Job) killer() {
	for {
		j.mu.Lock()
		if len(j.subscribers) == 0 {
			now := time.Now()
			if j.keepAliveBegin.IsZero() {
				j.keepAliveBegin = now
			}
			if now.Sub(j.keepAliveBegin) > j.keepAliveDuration {
				j.mu.Unlock()
				_ = j.Stop()
				return
			}
		} else {
			j.keepAliveBegin = time.Time{}
		}
		j.mu.Unlock()
		time.Sleep(time.Second)
	}
}

func (j *Job) broadcast(t EventType, m string) {
	j.mu.RLock()
	defer j.mu.RUnlock()
	for id, ch := range j.subscribers {
		select {
		case ch <- Message{Type: t, Timestamp: time.Now(), Message: m}:
		default: // slow subscriber: remove it
			j.removeSubscriber(id)
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
