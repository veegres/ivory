package service

import (
	. "ivory/model"
	"strconv"
	"strings"
)

var PatroniInstanceApiImpl = &patroniInstanceApiImpl{}

type patroniInstanceApiImpl struct{}

func (p patroniInstanceApiImpl) Info(instance InstanceRequest) (InstanceInfo, error) {
	response, err := Get[PatroniInfo](instance, "/patroni")
	info := InstanceInfo{
		Sidecar: Sidecar{Host: instance.Host, Port: instance.Port},
		Role:    response.Role,
		State:   response.State,
	}
	return info, err
}

func (p patroniInstanceApiImpl) Overview(instance InstanceRequest) ([]Instance, error) {
	var err error
	response, err := Get[PatroniCluster](instance, "/cluster")
	var overview []Instance
	for _, patroniInstance := range response.Members {
		domainString := strings.Split(patroniInstance.ApiUrl, "/")[2]
		domain := strings.Split(domainString, ":")
		host := domain[0]
		port, errCast := strconv.Atoi(domain[1])
		if errCast == nil {
			err = errCast
			break
		}

		overview = append(overview, Instance{
			State:    patroniInstance.State,
			Role:     patroniInstance.Role,
			Lag:      patroniInstance.Lag,
			Database: Database{Host: patroniInstance.Host, Port: patroniInstance.Port},
			Sidecar:  Sidecar{Host: host, Port: int8(port)},
		})
	}
	return overview, err
}

func (p patroniInstanceApiImpl) Config(instance InstanceRequest) (any, error) {
	return Get[any](instance, "/config")
}

func (p patroniInstanceApiImpl) ConfigUpdate(instance InstanceRequest) (any, error) {
	return Patch[any](instance, "/config")
}

func (p patroniInstanceApiImpl) Switchover(instance InstanceRequest) (any, error) {
	return Post[any](instance, "/switchover")
}

func (p patroniInstanceApiImpl) Reinitialize(instance InstanceRequest) (any, error) {
	return Post[any](instance, "/reinitialize")
}
