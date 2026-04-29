package patroni

import (
	"encoding/json"
	"ivory/src/clients/http"
	"ivory/src/features"
	"ivory/src/plugins/keeper"
	nethttp "net/http"
	"strconv"
	"strings"
)

// NOTE: validate that is matches interface in compile-time
var _ keeper.Adapter = (*Adapter)(nil)

type Adapter struct {
	httpClient *http.Client
}

func NewAdapter(httpClient *http.Client) *Adapter {
	return &Adapter{httpClient: httpClient}
}

func (a *Adapter) SupportedFeatures() []features.Feature {
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

func (a *Adapter) Overview(request keeper.Request) ([]keeper.Response, int, error) {
	var overview []keeper.Response

	response, status, err := http.NewJSONRequest[PatroniCluster](a.httpClient).Get(keeper.Map(request, "/cluster"))
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
			Key:                  &patroniInstance.Name,
			Status:               &keeperStatus,
			State:                patroniInstance.State,
			Role:                 a.mapRole(patroniInstance.Role),
			Lag:                  a.mapLag(patroniInstance.Lag),
			PendingRestart:       patroniInstance.PendingRestart,
			ScheduledRestart:     a.mapRestart(patroniInstance.ScheduledRestart),
			ScheduledSwitchover:  a.mapSwitchover(host, response.ScheduledSwitchover),
			Tags:                 patroniInstance.Tags,
			DiscoveredHost:       &patroniInstance.Host,
			DiscoveredDbPort:     &patroniInstance.Port,
			DiscoveredKeeperPort: &port,
		})
	}

	return overview, status, err
}

func (a *Adapter) mapLag(lag json.RawMessage) int64 {
	if num, errInt := strconv.ParseInt(string(lag), 10, 64); errInt == nil {
		return num
	}
	return -1
}

func (a *Adapter) mapRestart(restart *PatroniScheduledRestart) *keeper.ScheduledRestart {
	var scheduledRestart *keeper.ScheduledRestart
	if restart != nil {
		scheduledRestart = &keeper.ScheduledRestart{
			PendingRestart: restart.RestartPending,
			At:             restart.Schedule,
		}
	}
	return scheduledRestart
}

func (a *Adapter) mapSwitchover(host string, switchover *PatroniScheduledSwitchover) *keeper.ScheduledSwitchover {
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

func (a *Adapter) mapRole(role string) keeper.Role {
	switch role {
	case "leader", "master":
		return keeper.Leader
	case "replica":
		return keeper.Replica
	default:
		return keeper.Unknown
	}
}

func (a *Adapter) Config(request keeper.Request) (any, int, error) {
	return http.NewJSONRequest[any](a.httpClient).Get(keeper.Map(request, "/config"))
}

func (a *Adapter) ConfigUpdate(request keeper.Request) (any, int, error) {
	return http.NewJSONRequest[any](a.httpClient).Patch(keeper.Map(request, "/config"))
}

func (a *Adapter) Switchover(request keeper.Request) (*string, int, error) {
	return http.NewJSONRequest[string](a.httpClient).Post(keeper.Map(request, "/switchover"))
}

func (a *Adapter) DeleteSwitchover(request keeper.Request) (*string, int, error) {
	return http.NewJSONRequest[string](a.httpClient).Delete(keeper.Map(request, "/switchover"))
}

func (a *Adapter) Reinitialize(request keeper.Request) (*string, int, error) {
	return http.NewJSONRequest[string](a.httpClient).Post(keeper.Map(request, "/reinitialize"))
}

func (a *Adapter) Restart(request keeper.Request) (*string, int, error) {
	return http.NewJSONRequest[string](a.httpClient).Post(keeper.Map(request, "/restart"))
}

func (a *Adapter) DeleteRestart(request keeper.Request) (*string, int, error) {
	return http.NewJSONRequest[string](a.httpClient).Delete(keeper.Map(request, "/restart"))
}

func (a *Adapter) Reload(request keeper.Request) (*string, int, error) {
	return http.NewJSONRequest[string](a.httpClient).Post(keeper.Map(request, "/reload"))
}

func (a *Adapter) Failover(request keeper.Request) (*string, int, error) {
	return http.NewJSONRequest[string](a.httpClient).Post(keeper.Map(request, "/failover"))
}

func (a *Adapter) Activate(request keeper.Request) (*string, int, error) {
	if request.Body != nil {
		return nil, nethttp.StatusBadRequest, keeper.ErrBodyShouldBeEmpty
	}
	request.Body = ConfigPause{Pause: false}
	return http.NewJSONRequest[string](a.httpClient).Patch(keeper.Map(request, "/config"))
}

func (a *Adapter) Pause(request keeper.Request) (*string, int, error) {
	if request.Body != nil {
		return nil, nethttp.StatusBadRequest, keeper.ErrBodyShouldBeEmpty
	}
	request.Body = ConfigPause{Pause: true}
	return http.NewJSONRequest[string](a.httpClient).Patch(keeper.Map(request, "/config"))
}
