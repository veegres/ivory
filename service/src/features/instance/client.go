package instance

type InstanceClient interface {
	Overview(instance InstanceRequest) ([]Instance, int, error)
	Config(instance InstanceRequest) (any, int, error)
	ConfigUpdate(instance InstanceRequest) (any, int, error)
	Switchover(instance InstanceRequest) (*string, int, error)
	DeleteSwitchover(instance InstanceRequest) (*string, int, error)
	Reinitialize(instance InstanceRequest) (*string, int, error)
	Restart(instance InstanceRequest) (*string, int, error)
	DeleteRestart(instance InstanceRequest) (*string, int, error)
	Reload(instance InstanceRequest) (*string, int, error)
	Failover(instance InstanceRequest) (*string, int, error)
	Activate(instance InstanceRequest) (*string, int, error)
	Pause(instance InstanceRequest) (*string, int, error)
}
