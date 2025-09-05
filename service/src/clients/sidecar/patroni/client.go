package patroni

import (
	"errors"
	"ivory/src/clients/database"
	. "ivory/src/clients/sidecar"
	"net/http"
	"strconv"
	"strings"
)

type patroniClient struct {
	gateway *SidecarGateway
}

func NewPatroniClient(gateway *SidecarGateway) SidecarClient {
	return &patroniClient{gateway: gateway}
}

func (p *patroniClient) Overview(request Request) ([]Instance, int, error) {
	var overview []Instance

	response, status, err := NewSidecarRequest[PatroniCluster](p.gateway).Get(request, "/cluster")
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
			Database:            database.Database{Host: patroniInstance.Host, Port: patroniInstance.Port},
			Sidecar:             Sidecar{Host: host, Port: port, Name: &patroniInstance.Name, Status: &sidecarStatus},
			ScheduledRestart:    scheduledRestart,
			ScheduledSwitchover: scheduledSwitchover,
			Tags:                patroniInstance.Tags,
		})
	}

	return overview, status, err
}

func (p *patroniClient) Config(request Request) (any, int, error) {
	return NewSidecarRequest[any](p.gateway).Get(request, "/config")
}

func (p *patroniClient) ConfigUpdate(request Request) (any, int, error) {
	return NewSidecarRequest[any](p.gateway).Patch(request, "/config")
}

func (p *patroniClient) Switchover(request Request) (*string, int, error) {
	return NewSidecarRequest[string](p.gateway).Post(request, "/switchover")
}

func (p *patroniClient) DeleteSwitchover(request Request) (*string, int, error) {
	return NewSidecarRequest[string](p.gateway).Delete(request, "/switchover")
}

func (p *patroniClient) Reinitialize(request Request) (*string, int, error) {
	return NewSidecarRequest[string](p.gateway).Post(request, "/reinitialize")
}

func (p *patroniClient) Restart(request Request) (*string, int, error) {
	return NewSidecarRequest[string](p.gateway).Post(request, "/restart")
}

func (p *patroniClient) DeleteRestart(request Request) (*string, int, error) {
	return NewSidecarRequest[string](p.gateway).Delete(request, "/restart")
}

func (p *patroniClient) Reload(request Request) (*string, int, error) {
	return NewSidecarRequest[string](p.gateway).Post(request, "/reload")
}

func (p *patroniClient) Failover(request Request) (*string, int, error) {
	return NewSidecarRequest[string](p.gateway).Post(request, "/failover")
}

func (p *patroniClient) Activate(request Request) (*string, int, error) {
	if request.Body != nil {
		return nil, http.StatusBadRequest, errors.New("body should be empty")
	}
	request.Body = ConfigPause{Pause: false}
	return NewSidecarRequest[string](p.gateway).Patch(request, "/config")
}

func (p *patroniClient) Pause(request Request) (*string, int, error) {
	if request.Body != nil {
		return nil, http.StatusBadRequest, errors.New("body should be empty")
	}
	request.Body = ConfigPause{Pause: true}
	return NewSidecarRequest[string](p.gateway).Patch(request, "/config")
}
