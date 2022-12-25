package model

type Proxy interface {
    Info(cluster string, instance Instance) (interface{}, error)
    Overview(cluster string, instance Instance) (interface{}, error)
    Config(cluster string, instance Instance) (interface{}, error)
    ConfigUpdate(cluster string, instance Instance) (interface{}, error)
    Switchover(cluster string, instance Instance) (interface{}, error)
    Reinitialize(cluster string, instance Instance) (interface{}, error)
}

type Instance struct {
	Host   string `json:"host"`
	Port   int8   `json:"port"`
	Body   string `json:"body"`
}
