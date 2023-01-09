package service

import (
	. "ivory/model"
	"strconv"
	"strings"
)

var PatroniInstanceApiImpl = &patroniInstanceApiImpl{}

type patroniInstanceApiImpl struct{}

func (p patroniInstanceApiImpl) Info(instance InstanceRequest) (InstanceInfo, int, error) {
	var info InstanceInfo

	response, status, err := Get[PatroniInfo](instance, "/patroni")
	if err == nil {
		info = InstanceInfo{
			Sidecar: Sidecar{Host: instance.Host, Port: instance.Port},
			Role:    response.Role,
			State:   response.State,
		}
	}

	return info, status, err
}

func (p patroniInstanceApiImpl) Overview(instance InstanceRequest) ([]Instance, int, error) {
	var err error
	var overview []Instance

	response, status, err := Get[PatroniCluster](instance, "/cluster")
	if err == nil {
		for _, patroniInstance := range response.Members {
			domainString := strings.Split(patroniInstance.ApiUrl, "/")[2]
			domain := strings.Split(domainString, ":")
			host := domain[0]
			port, errCast := strconv.Atoi(domain[1])
			if errCast != nil {
				err = errCast
				break
			}

			var lag int
			if _, ok := patroniInstance.Lag.(string); ok {
				lag = -1
			}

			overview = append(overview, Instance{
				State:    patroniInstance.State,
				Role:     patroniInstance.Role,
				Lag:      lag,
				Database: Database{Host: patroniInstance.Host, Port: patroniInstance.Port},
				Sidecar:  Sidecar{Host: host, Port: port},
			})
		}
	}

	return overview, status, err
}

func (p patroniInstanceApiImpl) Config(instance InstanceRequest) (any, int, error) {
	return Get[any](instance, "/config")
}

func (p patroniInstanceApiImpl) ConfigUpdate(instance InstanceRequest) (any, int, error) {
	return Patch[any](instance, "/config")
}

func (p patroniInstanceApiImpl) Switchover(instance InstanceRequest) (any, int, error) {
	return Post[string](instance, "/switchover")
}

func (p patroniInstanceApiImpl) Reinitialize(instance InstanceRequest) (any, int, error) {
	return Post[string](instance, "/reinitialize")
}
