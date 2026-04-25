package patroni

import (
	"encoding/json"
	"errors"
	"ivory/src/clients/http"
	"ivory/src/clients/keeper"
	"ivory/src/features"
	nethttp "net/http"
	"strconv"
	"strings"
)

var ErrBodyShouldBeEmpty = errors.New("body should be empty")

// NOTE: validate that is matches interface in compile-time
var _ keeper.Client = (*Client)(nil)

type Client struct {
	httpClient *http.Client
}

func NewClient(httpClient *http.Client) *Client {
	return &Client{httpClient: httpClient}
}

func (p *Client) SupportedFeatures() []features.Feature {
	return []features.Feature{
		features.ViewNodeDbOverview,
		features.ViewNodeDbConfig,
		features.ManageNodeDbConfigUpdate,
		features.ManageNodeDbSwitchover,
		features.ManageNodeDbReinitialize,
		features.ManageNodeDbRestart,
		features.ManageNodeDbReload,
		features.ManageNodeDbFailover,
		features.ManageNodeDbActivation,
	}
}

func (p *Client) Overview(request keeper.Request) ([]keeper.Response, int, error) {
	var overview []keeper.Response

	response, status, err := http.NewJSONRequest[PatroniCluster](p.httpClient).Get(keeper.Map(request, "/cluster"))
	if err != nil {
		return overview, status, err
	}

	var keeperStatus keeper.Status
	if response.Pause == false {
		keeperStatus = keeper.Active
	} else {
		keeperStatus = keeper.Paused
	}

	for _, patroniInstance := range response.Members {
		domainString := strings.Split(patroniInstance.ApiUrl, "/")[2]
		domain := strings.Split(domainString, ":")
		host := domain[0]
		port, errCast := strconv.Atoi(domain[1])
		if errCast != nil {
			return nil, nethttp.StatusBadRequest, errCast
		}

		overview = append(overview, keeper.Response{
			Name:                 &patroniInstance.Name,
			Status:               &keeperStatus,
			State:                patroniInstance.State,
			Role:                 p.mapRole(patroniInstance.Role),
			Lag:                  p.mapLag(patroniInstance.Lag),
			PendingRestart:       patroniInstance.PendingRestart,
			ScheduledRestart:     p.mapRestart(patroniInstance.ScheduledRestart),
			ScheduledSwitchover:  p.mapSwitchover(host, response.ScheduledSwitchover),
			Tags:                 patroniInstance.Tags,
			DiscoveredHost:       patroniInstance.Host,
			DiscoveredDbPort:     patroniInstance.Port,
			DiscoveredKeeperPort: port,
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

func (p *Client) mapRestart(restart *PatroniScheduledRestart) *keeper.ScheduledRestart {
	var scheduledRestart *keeper.ScheduledRestart
	if restart != nil {
		scheduledRestart = &keeper.ScheduledRestart{
			PendingRestart: restart.RestartPending,
			At:             restart.Schedule,
		}
	}
	return scheduledRestart
}

func (p *Client) mapSwitchover(host string, switchover *PatroniScheduledSwitchover) *keeper.ScheduledSwitchover {
	var scheduledSwitchover *keeper.ScheduledSwitchover
	if switchover != nil && switchover.From == host {
		to := switchover.To
		if to == "" {
			to = "(random selection)"
		}

		scheduledSwitchover = &keeper.ScheduledSwitchover{
			At: switchover.At,
			To: to,
		}
	}
	return scheduledSwitchover
}

func (p *Client) mapRole(role string) keeper.Role {
	switch role {
	case "leader", "master":
		return keeper.Leader
	case "replica":
		return keeper.Replica
	default:
		return keeper.Unknown
	}
}

func (p *Client) Config(request keeper.Request) (any, int, error) {
	return http.NewJSONRequest[any](p.httpClient).Get(keeper.Map(request, "/config"))
}

func (p *Client) ConfigUpdate(request keeper.Request) (any, int, error) {
	return http.NewJSONRequest[any](p.httpClient).Patch(keeper.Map(request, "/config"))
}

func (p *Client) Switchover(request keeper.Request) (*string, int, error) {
	return http.NewJSONRequest[string](p.httpClient).Post(keeper.Map(request, "/switchover"))
}

func (p *Client) DeleteSwitchover(request keeper.Request) (*string, int, error) {
	return http.NewJSONRequest[string](p.httpClient).Delete(keeper.Map(request, "/switchover"))
}

func (p *Client) Reinitialize(request keeper.Request) (*string, int, error) {
	return http.NewJSONRequest[string](p.httpClient).Post(keeper.Map(request, "/reinitialize"))
}

func (p *Client) Restart(request keeper.Request) (*string, int, error) {
	return http.NewJSONRequest[string](p.httpClient).Post(keeper.Map(request, "/restart"))
}

func (p *Client) DeleteRestart(request keeper.Request) (*string, int, error) {
	return http.NewJSONRequest[string](p.httpClient).Delete(keeper.Map(request, "/restart"))
}

func (p *Client) Reload(request keeper.Request) (*string, int, error) {
	return http.NewJSONRequest[string](p.httpClient).Post(keeper.Map(request, "/reload"))
}

func (p *Client) Failover(request keeper.Request) (*string, int, error) {
	return http.NewJSONRequest[string](p.httpClient).Post(keeper.Map(request, "/failover"))
}

func (p *Client) Activate(request keeper.Request) (*string, int, error) {
	if request.Body != nil {
		return nil, nethttp.StatusBadRequest, ErrBodyShouldBeEmpty
	}
	request.Body = ConfigPause{Pause: false}
	return http.NewJSONRequest[string](p.httpClient).Patch(keeper.Map(request, "/config"))
}

func (p *Client) Pause(request keeper.Request) (*string, int, error) {
	if request.Body != nil {
		return nil, nethttp.StatusBadRequest, ErrBodyShouldBeEmpty
	}
	request.Body = ConfigPause{Pause: true}
	return http.NewJSONRequest[string](p.httpClient).Patch(keeper.Map(request, "/config"))
}
