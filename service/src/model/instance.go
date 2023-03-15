package model

// InstanceService TODO add common return types to cast interface to them and create mappers for each impl (patroni)
type InstanceService interface {
	Info(instance InstanceRequest) (InstanceInfo, int, error)
	Overview(instance InstanceRequest) ([]Instance, int, error)
	Config(instance InstanceRequest) (any, int, error)
	ConfigUpdate(instance InstanceRequest) (any, int, error)
	Switchover(instance InstanceRequest) (any, int, error)
	Reinitialize(instance InstanceRequest) (any, int, error)
}

type InstanceRequest struct {
	Cluster string `json:"cluster" form:"cluster"`
	Host    string `json:"host" form:"host"`
	Port    int    `json:"port" form:"port"`
	Body    any    `json:"body" form:"body"`
}

type Instance struct {
	State    string   `json:"state"`
	Role     string   `json:"role"`
	Lag      int      `json:"lag"`
	Database Database `json:"database"`
	Sidecar  Sidecar  `json:"sidecar"`
}

type Database struct {
	Host     string  `json:"host"`
	Port     int     `json:"port"`
	Database *string `json:"database"`
}

type Sidecar struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

type InstanceInfo struct {
	State   string  `json:"state"`
	Role    string  `json:"role"`
	Sidecar Sidecar `json:"sidecar"`
}
