package patroni

import (
	"encoding/json"
	"errors"
	"ivory/src/clients/database"
	"ivory/src/clients/http"
	"ivory/src/clients/sidecar"
	nethttp "net/http"
	"strconv"
	"strings"
)

var ErrBodyShouldBeEmpty = errors.New("body should be empty")

// NOTE: validate that is matches interface in compile-time
var _ sidecar.Client = (*Client)(nil)

type Client struct {
	httpClient *http.Client
}

func NewClient(httpClient *http.Client) *Client {
	return &Client{httpClient: httpClient}
}

func (p *Client) Overview(request sidecar.Request) ([]sidecar.Instance, int, error) {
	var overview []sidecar.Instance

	response, status, err := http.NewJSONRequest[PatroniCluster](p.httpClient).Get(sidecar.Map(request, "/cluster"))
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
			return nil, nethttp.StatusBadRequest, errCast
		}

		overview = append(overview, sidecar.Instance{
			State:               patroniInstance.State,
			Role:                p.mapRole(patroniInstance.Role),
			Lag:                 p.mapLag(patroniInstance.Lag),
			PendingRestart:      patroniInstance.PendingRestart,
			Database:            database.Database{Host: patroniInstance.Host, Port: patroniInstance.Port},
			Sidecar:             sidecar.Sidecar{Host: host, Port: port, Name: &patroniInstance.Name, Status: &sidecarStatus},
			ScheduledRestart:    p.mapRestart(patroniInstance.ScheduledRestart),
			ScheduledSwitchover: p.mapSwitchover(host, response.ScheduledSwitchover),
			Tags:                patroniInstance.Tags,
		})
	}

	return overview, status, err
}

func (p *Client) mapLag(lag json.RawMessage) int64 {
	if num, errInt := strconv.ParseInt(string(lag), 10, 64); errInt == nil {
		return num
	}
	return -1
}

func (p *Client) mapRestart(restart *PatroniScheduledRestart) *sidecar.InstanceScheduledRestart {
	var scheduledRestart *sidecar.InstanceScheduledRestart
	if restart != nil {
		scheduledRestart = &sidecar.InstanceScheduledRestart{
			PendingRestart: restart.RestartPending,
			At:             restart.Schedule,
		}
	}
	return scheduledRestart
}

func (p *Client) mapSwitchover(host string, switchover *PatroniScheduledSwitchover) *sidecar.InstanceScheduledSwitchover {
	var scheduledSwitchover *sidecar.InstanceScheduledSwitchover
	if switchover != nil && switchover.From == host {
		to := switchover.To
		if to == "" {
			to = "(random selection)"
		}

		scheduledSwitchover = &sidecar.InstanceScheduledSwitchover{
			At: switchover.At,
			To: to,
		}
	}
	return scheduledSwitchover
}

func (p *Client) mapRole(role string) sidecar.Role {
	switch role {
	case "leader", "master":
		return sidecar.Leader
	case "replica":
		return sidecar.Replica
	default:
		return sidecar.Unknown
	}
}

func (p *Client) Config(request sidecar.Request) (any, int, error) {
	return http.NewJSONRequest[any](p.httpClient).Get(sidecar.Map(request, "/config"))
}

func (p *Client) ConfigUpdate(request sidecar.Request) (any, int, error) {
	return http.NewJSONRequest[any](p.httpClient).Patch(sidecar.Map(request, "/config"))
}

func (p *Client) Switchover(request sidecar.Request) (*string, int, error) {
	return http.NewJSONRequest[string](p.httpClient).Post(sidecar.Map(request, "/switchover"))
}

func (p *Client) DeleteSwitchover(request sidecar.Request) (*string, int, error) {
	return http.NewJSONRequest[string](p.httpClient).Delete(sidecar.Map(request, "/switchover"))
}

func (p *Client) Reinitialize(request sidecar.Request) (*string, int, error) {
	return http.NewJSONRequest[string](p.httpClient).Post(sidecar.Map(request, "/reinitialize"))
}

func (p *Client) Restart(request sidecar.Request) (*string, int, error) {
	return http.NewJSONRequest[string](p.httpClient).Post(sidecar.Map(request, "/restart"))
}

func (p *Client) DeleteRestart(request sidecar.Request) (*string, int, error) {
	return http.NewJSONRequest[string](p.httpClient).Delete(sidecar.Map(request, "/restart"))
}

func (p *Client) Reload(request sidecar.Request) (*string, int, error) {
	return http.NewJSONRequest[string](p.httpClient).Post(sidecar.Map(request, "/reload"))
}

func (p *Client) Failover(request sidecar.Request) (*string, int, error) {
	return http.NewJSONRequest[string](p.httpClient).Post(sidecar.Map(request, "/failover"))
}

func (p *Client) Activate(request sidecar.Request) (*string, int, error) {
	if request.Body != nil {
		return nil, nethttp.StatusBadRequest, ErrBodyShouldBeEmpty
	}
	request.Body = ConfigPause{Pause: false}
	return http.NewJSONRequest[string](p.httpClient).Patch(sidecar.Map(request, "/config"))
}

func (p *Client) Pause(request sidecar.Request) (*string, int, error) {
	if request.Body != nil {
		return nil, nethttp.StatusBadRequest, ErrBodyShouldBeEmpty
	}
	request.Body = ConfigPause{Pause: true}
	return http.NewJSONRequest[string](p.httpClient).Patch(sidecar.Map(request, "/config"))
}
