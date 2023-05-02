package service

import (
	. "ivory/src/model"
	"strconv"
	"strings"
)

type patroniInstanceService struct {
	client *SidecarClient
}

func NewPatroniGateway(client *SidecarClient) InstanceGateway {
	return &patroniInstanceService{client: client}
}

func (p *patroniInstanceService) Info(instance InstanceRequest) (InstanceInfo, int, error) {
	var info InstanceInfo

	response, status, err := NewSidecarRequest[PatroniInfo](p.client).Get(instance, "/patroni")
	if err == nil {
		info = InstanceInfo{
			Sidecar: Sidecar{Host: instance.Host, Port: instance.Port},
			Role:    response.Role,
			State:   response.State,
		}
	}

	return info, status, err
}

func (p *patroniInstanceService) Overview(instance InstanceRequest) ([]Instance, int, error) {
	var overview []Instance

	response, status, err := NewSidecarRequest[PatroniCluster](p.client).Get(instance, "/cluster")
	if err != nil {
		return overview, status, err
	}

	for _, patroniInstance := range response.Members {
		domainString := strings.Split(patroniInstance.ApiUrl, "/")[2]
		domain := strings.Split(domainString, ":")
		host := domain[0]
		port, errCast := strconv.Atoi(domain[1])
		if errCast != nil {
			return nil, 0, errCast
		}

		var lag int
		if ok := patroniInstance.Lag; ok == nil {
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

	return overview, status, err
}

func (p *patroniInstanceService) Config(instance InstanceRequest) (any, int, error) {
	return NewSidecarRequest[any](p.client).Get(instance, "/config")
}

func (p *patroniInstanceService) ConfigUpdate(instance InstanceRequest) (any, int, error) {
	return NewSidecarRequest[any](p.client).Patch(instance, "/config")
}

func (p *patroniInstanceService) Switchover(instance InstanceRequest) (any, int, error) {
	return NewSidecarRequest[string](p.client).Post(instance, "/switchover")
}

func (p *patroniInstanceService) Reinitialize(instance InstanceRequest) (any, int, error) {
	return NewSidecarRequest[string](p.client).Post(instance, "/reinitialize")
}
