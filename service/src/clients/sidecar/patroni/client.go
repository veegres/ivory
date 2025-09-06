package patroni

import (
	"errors"
	"ivory/src/clients/database"
	"ivory/src/clients/sidecar"
	"net/http"
	"strconv"
	"strings"
)

type Client struct {
	gateway *sidecar.Gateway
}

func NewClient(gateway *sidecar.Gateway) *Client {
	return &Client{gateway: gateway}
}

func (p *Client) Overview(request sidecar.Request) ([]sidecar.Instance, int, error) {
	var overview []sidecar.Instance

	response, status, err := sidecar.NewSidecarRequest[PatroniCluster](p.gateway).Get(request, "/cluster")
	if err != nil {
		return overview, status, err
	}

	var sidecarStatus sidecar.SidecarStatus
	if response.Pause == false {
		sidecarStatus = sidecar.Active
	} else {
		sidecarStatus = sidecar.Paused
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

		var scheduledRestart *sidecar.InstanceScheduledRestart
		if patroniInstance.ScheduledRestart != nil {
			scheduledRestart = &sidecar.InstanceScheduledRestart{
				PendingRestart: patroniInstance.ScheduledRestart.RestartPending,
				At:             patroniInstance.ScheduledRestart.Schedule,
			}
		}

		var scheduledSwitchover *sidecar.InstanceScheduledSwitchover
		if response.ScheduledSwitchover != nil && response.ScheduledSwitchover.From == host {
			to := response.ScheduledSwitchover.To
			if to == "" {
				to = "(random selection)"
			}

			scheduledSwitchover = &sidecar.InstanceScheduledSwitchover{
				At: response.ScheduledSwitchover.At,
				To: to,
			}
		}

		overview = append(overview, sidecar.Instance{
			State:               patroniInstance.State,
			Role:                patroniInstance.Role,
			Lag:                 lag,
			PendingRestart:      patroniInstance.PendingRestart,
			Database:            database.Database{Host: patroniInstance.Host, Port: patroniInstance.Port},
			Sidecar:             sidecar.Sidecar{Host: host, Port: port, Name: &patroniInstance.Name, Status: &sidecarStatus},
			ScheduledRestart:    scheduledRestart,
			ScheduledSwitchover: scheduledSwitchover,
			Tags:                patroniInstance.Tags,
		})
	}

	return overview, status, err
}

func (p *Client) Config(request sidecar.Request) (any, int, error) {
	return sidecar.NewSidecarRequest[any](p.gateway).Get(request, "/config")
}

func (p *Client) ConfigUpdate(request sidecar.Request) (any, int, error) {
	return sidecar.NewSidecarRequest[any](p.gateway).Patch(request, "/config")
}

func (p *Client) Switchover(request sidecar.Request) (*string, int, error) {
	return sidecar.NewSidecarRequest[string](p.gateway).Post(request, "/switchover")
}

func (p *Client) DeleteSwitchover(request sidecar.Request) (*string, int, error) {
	return sidecar.NewSidecarRequest[string](p.gateway).Delete(request, "/switchover")
}

func (p *Client) Reinitialize(request sidecar.Request) (*string, int, error) {
	return sidecar.NewSidecarRequest[string](p.gateway).Post(request, "/reinitialize")
}

func (p *Client) Restart(request sidecar.Request) (*string, int, error) {
	return sidecar.NewSidecarRequest[string](p.gateway).Post(request, "/restart")
}

func (p *Client) DeleteRestart(request sidecar.Request) (*string, int, error) {
	return sidecar.NewSidecarRequest[string](p.gateway).Delete(request, "/restart")
}

func (p *Client) Reload(request sidecar.Request) (*string, int, error) {
	return sidecar.NewSidecarRequest[string](p.gateway).Post(request, "/reload")
}

func (p *Client) Failover(request sidecar.Request) (*string, int, error) {
	return sidecar.NewSidecarRequest[string](p.gateway).Post(request, "/failover")
}

func (p *Client) Activate(request sidecar.Request) (*string, int, error) {
	if request.Body != nil {
		return nil, http.StatusBadRequest, errors.New("body should be empty")
	}
	request.Body = ConfigPause{Pause: false}
	return sidecar.NewSidecarRequest[string](p.gateway).Patch(request, "/config")
}

func (p *Client) Pause(request sidecar.Request) (*string, int, error) {
	if request.Body != nil {
		return nil, http.StatusBadRequest, errors.New("body should be empty")
	}
	request.Body = ConfigPause{Pause: true}
	return sidecar.NewSidecarRequest[string](p.gateway).Patch(request, "/config")
}
