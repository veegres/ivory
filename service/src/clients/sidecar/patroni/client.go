package patroni

import (
	"errors"
	. "ivory/src/clients/sidecar"
	. "ivory/src/features/instance"
	"ivory/src/features/query"
	"net/http"
	"strconv"
	"strings"
)

type patroniInstanceService struct {
	client *SidecarGateway
}

func NewPatroniClient(client *SidecarGateway) InstanceClient {
	return &patroniInstanceService{client: client}
}

func (p *patroniInstanceService) Overview(instance InstanceRequest) ([]Instance, int, error) {
	var overview []Instance

	response, status, err := NewSidecarRequest[PatroniCluster](p.client).Get(instance, "/cluster")
	if err != nil {
		return overview, status, err
	}

	var sidecarStatus SidecarStatus
	if response.Pause == false {
		sidecarStatus = Active
	} else {
		sidecarStatus = Paused
	}

	for _, patroniInstance := range response.Members {
		domainString := strings.Split(patroniInstance.ApiUrl, "/")[2]
		domain := strings.Split(domainString, ":")
		host := domain[0]
		port, errCast := strconv.Atoi(domain[1])
		if errCast != nil {
			return nil, http.StatusBadRequest, errCast
		}

		var lag int64
		if num, errInt := strconv.ParseInt(string(patroniInstance.Lag), 10, 64); errInt == nil {
			lag = num
		} else {
			lag = -1
		}

		var scheduledRestart *InstanceScheduledRestart
		if patroniInstance.ScheduledRestart != nil {
			scheduledRestart = &InstanceScheduledRestart{
				PendingRestart: patroniInstance.ScheduledRestart.RestartPending,
				At:             patroniInstance.ScheduledRestart.Schedule,
			}
		}

		var scheduledSwitchover *InstanceScheduledSwitchover
		if response.ScheduledSwitchover != nil && response.ScheduledSwitchover.From == host {
			to := response.ScheduledSwitchover.To
			if to == "" {
				to = "(random selection)"
			}

			scheduledSwitchover = &InstanceScheduledSwitchover{
				At: response.ScheduledSwitchover.At,
				To: to,
			}
		}

		overview = append(overview, Instance{
			State:               patroniInstance.State,
			Role:                patroniInstance.Role,
			Lag:                 lag,
			PendingRestart:      patroniInstance.PendingRestart,
			Database:            query.Database{Host: patroniInstance.Host, Port: patroniInstance.Port},
			Sidecar:             Sidecar{Host: host, Port: port, Name: &patroniInstance.Name, Status: &sidecarStatus},
			ScheduledRestart:    scheduledRestart,
			ScheduledSwitchover: scheduledSwitchover,
			Tags:                patroniInstance.Tags,
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

func (p *patroniInstanceService) Switchover(instance InstanceRequest) (*string, int, error) {
	return NewSidecarRequest[string](p.client).Post(instance, "/switchover")
}

func (p *patroniInstanceService) DeleteSwitchover(instance InstanceRequest) (*string, int, error) {
	return NewSidecarRequest[string](p.client).Delete(instance, "/switchover")
}

func (p *patroniInstanceService) Reinitialize(instance InstanceRequest) (*string, int, error) {
	return NewSidecarRequest[string](p.client).Post(instance, "/reinitialize")
}

func (p *patroniInstanceService) Restart(instance InstanceRequest) (*string, int, error) {
	return NewSidecarRequest[string](p.client).Post(instance, "/restart")
}

func (p *patroniInstanceService) DeleteRestart(instance InstanceRequest) (*string, int, error) {
	return NewSidecarRequest[string](p.client).Delete(instance, "/restart")
}

func (p *patroniInstanceService) Reload(instance InstanceRequest) (*string, int, error) {
	return NewSidecarRequest[string](p.client).Post(instance, "/reload")
}

func (p *patroniInstanceService) Failover(instance InstanceRequest) (*string, int, error) {
	return NewSidecarRequest[string](p.client).Post(instance, "/failover")
}

func (p *patroniInstanceService) Activate(instance InstanceRequest) (*string, int, error) {
	if instance.Body != nil {
		return nil, http.StatusBadRequest, errors.New("body should be empty")
	}
	instance.Body = ConfigPause{Pause: false}
	return NewSidecarRequest[string](p.client).Patch(instance, "/config")
}

func (p *patroniInstanceService) Pause(instance InstanceRequest) (*string, int, error) {
	if instance.Body != nil {
		return nil, http.StatusBadRequest, errors.New("body should be empty")
	}
	instance.Body = ConfigPause{Pause: true}
	return NewSidecarRequest[string](p.client).Patch(instance, "/config")
}
