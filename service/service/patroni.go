package service

import (
	. "ivory/model"
)

var Patroni = &ptr{}
type ptr struct{}

func (p ptr) Info(cluster string, instance Instance) (interface{}, error) {
    return proxy.Get(cluster, instance, "/cluster")
}

func (p ptr) Overview(cluster string, instance Instance) (interface{}, error) {
    return proxy.Get(cluster, instance, "/patroni")
}

func (p ptr) Config(cluster string, instance Instance) (interface{}, error) {
    return proxy.Get(cluster, instance, "/config")
}

func (p ptr) ConfigUpdate(cluster string, instance Instance) (interface{}, error) {
    return proxy.Patch(cluster, instance, "/config")
}

func (p ptr) Switchover(cluster string, instance Instance) (interface{}, error) {
    return proxy.Post(cluster, instance, "/switchover")
}

func (p ptr) Reinitialize(cluster string, instance Instance) (interface{}, error) {
    return proxy.Post(cluster, instance, "/reinitialize")
}
