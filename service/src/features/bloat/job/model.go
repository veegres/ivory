package job

import "strconv"

// COMMON (WEB AND SERVER)

type JobStatus int8

const (
	PENDING JobStatus = iota
	RUNNING
	FINISHED
	FAILED
	STOPPED
	UNKNOWN
)

func (j JobStatus) String() string {
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

type Event struct {
	Type    EventType
	Message string
}
