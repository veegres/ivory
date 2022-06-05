package service

import (
	. "ivory/model"
	"os/exec"
	"sync"
)

type Job interface {
	SetCommand(cmd *exec.Cmd)
	GetCommand() *exec.Cmd
	GetStatus() JobStatus
	SetStatus(status JobStatus)
	Subscribers() map[chan Event]bool
	Subscribe() chan Event
	Unsubscribe(channel chan Event)
	Size() int
	IsJobActive() bool
}

type job struct {
	subscribers map[chan Event]bool
	status      JobStatus
	size        int
	command     *exec.Cmd
	mutex       *sync.Mutex
}

func CreateJob() Job {
	return &job{
		subscribers: make(map[chan Event]bool),
		size:        0,
		status:      PENDING,
		mutex:       &sync.Mutex{},
	}
}

func (j *job) SetCommand(cmd *exec.Cmd) {
	j.mutex.Lock()
	j.command = cmd
	j.mutex.Unlock()
}

func (j *job) GetCommand() *exec.Cmd {
	return j.command
}

func (j *job) GetStatus() JobStatus {
	return j.status
}

func (j *job) SetStatus(status JobStatus) {
	j.mutex.Lock()
	j.status = status
	j.mutex.Unlock()
}

func (j *job) Subscribers() map[chan Event]bool {
	return j.subscribers
}

func (j *job) Subscribe() chan Event {
	var channel chan Event
	j.mutex.Lock()
	if j.IsJobActive() {
		channel = make(chan Event)
		j.subscribers[channel] = true
		j.incrementSubs()
	}
	j.mutex.Unlock()
	return channel
}

func (j *job) Unsubscribe(channel chan Event) {
	j.mutex.Lock()
	delete(j.subscribers, channel)
	j.decrementSubs()
	j.mutex.Unlock()
}

func (j *job) Size() int {
	return j.size
}

func (j *job) IsJobActive() bool {
	return j.GetStatus() == PENDING || j.GetStatus() == RUNNING
}

func (j *job) incrementSubs() int {
	j.size++
	return j.size
}

func (j *job) decrementSubs() int {
	j.size--
	return j.size
}
