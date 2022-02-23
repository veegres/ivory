package service

import (
	. "ivory/model"
	"sync"
)

type Job struct {
	Events      map[chan Event]bool
	Status      JobStatus
	Subscribers *int
	Mutex       *sync.Mutex
}

func (j Job) Subscribe() chan Event {
	channel := make(chan Event)
	j.Mutex.Lock()
	j.Events[channel] = true
	j.incrementSubs()
	j.Mutex.Unlock()
	return channel
}

func (j Job) Unsubscribe(channel chan Event) {
	j.Mutex.Lock()
	delete(j.Events, channel)
	j.decrementSubs()
	j.Mutex.Unlock()
}

func (j Job) Size() int {
	return *j.Subscribers
}

func (j Job) incrementSubs() int {
	*j.Subscribers++
	return *j.Subscribers
}

func (j Job) decrementSubs() int {
	*j.Subscribers--
	return *j.Subscribers
}
