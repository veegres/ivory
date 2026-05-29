package job

import (
	"io"
	"strconv"
	"time"
)

// COMMON (WEB AND SERVER)

// JobID is the SHA256
type JobID string
type SubscriberID string

type Status int8

const (
	PENDING Status = iota
	RUNNING
	FINISHED
	FAILED
	STOPPED
	UNKNOWN
)

func (j Status) String() string {
	return strconv.Itoa(int(j))
}

type Command interface {
	// Id returns a stable string used to identify the job.
	Id() string
	// Persist returns true if the job output should be saved to a file.
	Persist() bool
	// Start begins the command and returns a reader over its output.
	Start() (io.Reader, error)
	// Wait blocks until the command exits and returns its exit error.
	Wait() error
	// Abort kills the running command immediately.
	Abort() error
	// Execute runs the command synchronously and returns the output as string array.
	Execute() ([]string, error)
}

type EventType int8

const (
	SERVER EventType = iota
	STATUS
	LOG
	STREAM
)

func (s EventType) String() string {
	return []string{"server", "status", "log", "stream"}[s]
}

type EventStreamType int8

const (
	START EventStreamType = iota
	END
)

func (s EventStreamType) String() string {
	return []string{"start", "end"}[s]
}

// SPECIFIC (SERVER)

type Message struct {
	Type      EventType
	Timestamp time.Time
	Message   string
}
