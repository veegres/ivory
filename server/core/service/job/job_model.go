package job

import (
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
