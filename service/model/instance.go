package model

// InstanceApi TODO add common return types to cast interface to them and create mappers for each impl (patroni)
type InstanceApi interface {
	Info(instance InstanceRequest) (interface{}, error)
	Overview(instance InstanceRequest) (interface{}, error)
	Config(instance InstanceRequest) (interface{}, error)
	ConfigUpdate(instance InstanceRequest) (interface{}, error)
	Switchover(instance InstanceRequest) (interface{}, error)
	Reinitialize(instance InstanceRequest) (interface{}, error)
}

type InstanceRequest struct {
	Cluster string `json:"cluster"`
	Host    string `json:"host"`
	Port    int8   `json:"port"`
	Body    string `json:"body"`
}
