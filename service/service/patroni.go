package service

import (
	. "ivory/model"
)

var Patroni = &ptr{}

type ptr struct{}

func (p ptr) Info(instance InstanceRequest) (interface{}, error) {
	return proxy.Get(instance, "/cluster")
}

func (p ptr) Overview(instance InstanceRequest) (interface{}, error) {
	return proxy.Get(instance, "/patroni")
}

func (p ptr) Config(instance InstanceRequest) (interface{}, error) {
	return proxy.Get(instance, "/config")
}

func (p ptr) ConfigUpdate(instance InstanceRequest) (interface{}, error) {
	return proxy.Patch(instance, "/config")
}

func (p ptr) Switchover(instance InstanceRequest) (interface{}, error) {
	return proxy.Post(instance, "/switchover")
}

func (p ptr) Reinitialize(instance InstanceRequest) (interface{}, error) {
	return proxy.Post(instance, "/reinitialize")
}
