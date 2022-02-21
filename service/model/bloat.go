package model

import (
	"github.com/google/uuid"
	"strconv"
	"sync/atomic"
)

type DbConnection struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type Target struct {
	DbName        string `json:"dbName"`
	Schema        string `json:"schema"`
	Table         string `json:"table"`
	ExcludeSchema string `json:"excludeSchema"`
	ExcludeTable  string `json:"excludeTable"`
}

type CompactTableRequest struct {
	Connection DbConnection `json:"connection"`
	Target     *Target      `json:"target"`
	Ratio      int          `json:"ratio"`
}

type JobStatus int8

const (
	PENDING JobStatus = iota
	RUNNING
	FINISHED
	FAILED
)

func (j JobStatus) String() string {
	return strconv.Itoa(int(j))
}

type Job struct {
	Event       chan Event
	Subscribers int32
}

func (j Job) IncrementSubs() int32 {
	return atomic.AddInt32(&j.Subscribers, 1)
}

func (j Job) DecrementSubs() int32 {
	return atomic.AddInt32(&j.Subscribers, -1)
}

func (j Job) GetSubs() int32 {
	return atomic.LoadInt32(&j.Subscribers)
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

type CompactTableModel struct {
	Uuid        uuid.UUID `json:"uuid"`
	Status      JobStatus `json:"status"`
	Command     string    `json:"command"`
	CommandArgs []string  `json:"commandArgs"`
	LogsPath    string    `json:"logsPath"`
}
