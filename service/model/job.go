package model

import "strconv"

type JobStatus int8

const (
	PENDING JobStatus = iota
	RUNNING
	FINISHED
	FAILED
	STOPPED
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

type Event struct {
	Name    EventType
	Message string
}
