package service

import (
	. "ivory/model"
	"os/exec"
	"sync"
)

type Job struct {
	subscribers map[chan Event]bool
	status      JobStatus
	size        int
	command     *exec.Cmd
	mutex       *sync.Mutex
}

func NewJob() *Job {
	return &Job{
		subscribers: make(map[chan Event]bool),
		size:        0,
		status:      PENDING,
		mutex:       &sync.Mutex{},
	}
}

func (j *Job) SetCommand(cmd *exec.Cmd) {
	j.mutex.Lock()
	j.command = cmd
	j.mutex.Unlock()
}

func (j *Job) GetCommand() *exec.Cmd {
	return j.command
}

func (j *Job) GetStatus() JobStatus {
	return j.status
}

func (j *Job) SetStatus(status JobStatus) {
	j.mutex.Lock()
	j.status = status
	j.mutex.Unlock()
}

func (j *Job) Subscribers() map[chan Event]bool {
	return j.subscribers
}

func (j *Job) Subscribe() chan Event {
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

func (j *Job) Unsubscribe(channel chan Event) {
	j.mutex.Lock()
	delete(j.subscribers, channel)
	j.decrementSubs()
	j.mutex.Unlock()
}

func (j *Job) Size() int {
	return j.size
}

func (j *Job) IsJobActive() bool {
	return j.GetStatus() == PENDING || j.GetStatus() == RUNNING
}

func (j *Job) incrementSubs() int {
	j.size++
	return j.size
}

func (j *Job) decrementSubs() int {
	j.size--
	return j.size
}
