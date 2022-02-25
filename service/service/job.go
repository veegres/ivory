package service

import (
	. "ivory/model"
	"sync"
)

type Job struct {
	subscribers map[chan Event]bool
	status      JobStatus
	size        int
	mutex       *sync.Mutex
}

func (j Job) Create() *Job {
	j.subscribers = make(map[chan Event]bool)
	j.size = 0
	j.status = PENDING
	j.mutex = &sync.Mutex{}
	return &j
}

func (j *Job) Status() JobStatus {
	return j.status
}

func (j *Job) Next() JobStatus {
	j.mutex.Lock()
	switch j.status {
	case PENDING:
		j.status = RUNNING
	case RUNNING:
		j.status = FINISHED
	}
	tmp := j.status
	j.mutex.Unlock()
	return tmp
}

func (j *Job) Failed() {
	j.mutex.Lock()
	j.status = FAILED
	j.mutex.Unlock()
}

func (j *Job) Subscribers() map[chan Event]bool {
	return j.subscribers
}

func (j *Job) Subscribe() chan Event {
	var channel chan Event
	j.mutex.Lock()
	if j.Status() != FINISHED {
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

func (j *Job) IsFinished() bool {
	return j.Status() == FINISHED && j.Size() == 0
}

func (j *Job) incrementSubs() int {
	j.size++
	return j.size
}

func (j *Job) decrementSubs() int {
	j.size--
	return j.size
}
