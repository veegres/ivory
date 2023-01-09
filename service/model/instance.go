package model

// InstanceApi TODO add common return types to cast interface to them and create mappers for each impl (patroni)
type InstanceApi interface {
	Info(instance InstanceRequest) (InstanceInfo, error)
	Overview(instance InstanceRequest) ([]Instance, error)
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

type Instance struct {
	State    string   `json:"state"`
	Role     string   `json:"role"`
	Lag      int32    `json:"lag"`
	Database Database `json:"database"`
	Sidecar  Sidecar  `json:"sidecar"`
}

type Database struct {
	Host string `json:"host"`
	Port int8   `json:"port"`
}

type Sidecar struct {
	Host string `json:"host"`
	Port int8   `json:"port"`
}

type InstanceInfo struct {
	State   string  `json:"state"`
	Role    string  `json:"role"`
	Sidecar Sidecar `json:"sidecar"`
}
