package job

import (
	"bufio"
	"fmt"
	"ivory/clients/console"
	"ivory/clients/storage"
	"log/slog"
	"maps"
	"os"
	"regexp"
	"slices"
	"strings"
	"sync"
	"time"
)

// ansiEscape matches ANSI CSI escape sequences (cursor movement, color,
// clear-line, ...) that tools like "docker pull" write to a terminal for a
// progress bar; a job's subscribers and its persisted log file are plain
// text, not a terminal, so these would otherwise land verbatim.
var ansiEscape = regexp.MustCompile("\x1b\\[[0-9;]*[a-zA-Z]")

type Job struct {
	cmd         console.Command
	subscribers map[SubscriberID]chan Message
	status      Status
	mu          sync.RWMutex
	storage     *storage.FileStorage

	keepAliveDuration      time.Duration
	keepAliveCheckInterval time.Duration
	keepAliveBegin         time.Time
}

func NewJob(cmd console.Command, storage *storage.FileStorage) *Job {
	return &Job{
		subscribers:            make(map[SubscriberID]chan Message),
		cmd:                    cmd,
		status:                 PENDING,
		storage:                storage,
		keepAliveDuration:      time.Second * 5,
		keepAliveCheckInterval: time.Second,
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
		defer func() {
			if err := logFile.Close(); err != nil {
				slog.Error("failed to close log file", "error", err)
			}
		}()
	}

	reader, errStart := j.cmd.Start()
	if errStart != nil {
		j.setStatus(FAILED)
		j.broadcast(SERVER, errStart.Error())
		return
	}

	// NOTE: run killer which will kill stale jobs
	if !j.cmd.KeepAlive() {
		go j.killer()
	}

	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		msg := j.stripAnsi(scanner.Text())
		j.broadcast(LOG, msg)
		if logFile != nil {
			_, _ = logFile.WriteString(msg + "\n")
		}
	}

	if j.getStatus() == STOPPED {
		return
	}

	if err := scanner.Err(); err != nil {
		// NOTE: Ignore PTY EIO error on process exit
		if !strings.Contains(err.Error(), "input/output error") {
			j.setStatus(FAILED)
			j.broadcast(SERVER, err.Error())
			return
		}
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
	ticker := time.NewTicker(j.keepAliveCheckInterval)
	defer ticker.Stop()

	for {
		<-ticker.C

		shouldStop := false
		j.mu.Lock()
		if j.status != RUNNING {
			j.mu.Unlock()
			return
		}
		if len(j.subscribers) == 0 {
			now := time.Now()
			if j.keepAliveBegin.IsZero() {
				j.keepAliveBegin = now
			}
			shouldStop = now.Sub(j.keepAliveBegin) > j.keepAliveDuration
		} else {
			j.keepAliveBegin = time.Time{}
		}
		j.mu.Unlock()

		if shouldStop {
			if err := j.Stop(); err != nil {
				slog.Error("failed to stop job", "error", err)
			}
			return
		}
	}
}

func (j *Job) stripAnsi(line string) string {
	return ansiEscape.ReplaceAllString(line, "")
}

func (j *Job) broadcast(t EventType, m string) {
	j.mu.RLock()
	defer j.mu.RUnlock()
	for _, ch := range j.subscribers {
		select {
		case ch <- Message{Type: t, Timestamp: time.Now(), Message: m}:
		default: // slow subscriber: skip it
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

func (j *Job) addSubscriber(id SubscriberID) (*Subscription, bool) {
	j.mu.Lock()
	defer j.mu.Unlock()
	if j.status != PENDING && j.status != RUNNING {
		return nil, false
	}

	if old, exists := j.subscribers[id]; exists {
		close(old)
	}

	channel := make(chan Message, 256)
	j.subscribers[id] = channel
	return &Subscription{Messages: channel, job: j, id: id, ch: channel}, true
}

func (j *Job) removeSubscriber(id SubscriberID, ch chan Message) {
	j.mu.Lock()
	if current, exists := j.subscribers[id]; exists && current == ch {
		close(current)
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
