package service

import (
	. "ivory/model"
	"strconv"
	"strings"
)

type patroniInstanceService struct {
	proxy *Proxy
}

func NewPatroniService(proxy *Proxy) InstanceService {
	return &patroniInstanceService{
		proxy: proxy,
	}
}

func (p *patroniInstanceService) Info(instance InstanceRequest) (InstanceInfo, int, error) {
	var info InstanceInfo

	response, status, err := NewProxyRequest[PatroniInfo](p.proxy).Get(instance, "/patroni")
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
	var err error
	var overview []Instance

	response, status, err := NewProxyRequest[PatroniCluster](p.proxy).Get(instance, "/cluster")
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

func (p *patroniInstanceService) Config(instance InstanceRequest) (any, int, error) {
	return NewProxyRequest[any](p.proxy).Get(instance, "/config")
}

func (p *patroniInstanceService) ConfigUpdate(instance InstanceRequest) (any, int, error) {
	return NewProxyRequest[any](p.proxy).Patch(instance, "/config")
}

func (p *patroniInstanceService) Switchover(instance InstanceRequest) (any, int, error) {
	return NewProxyRequest[string](p.proxy).Post(instance, "/switchover")
}

func (p *patroniInstanceService) Reinitialize(instance InstanceRequest) (any, int, error) {
	return NewProxyRequest[string](p.proxy).Post(instance, "/reinitialize")
}
