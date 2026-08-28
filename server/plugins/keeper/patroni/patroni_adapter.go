package patroni

import (
	"encoding/json"
	"ivory/clients/http"
	"ivory/plugins/keeper"
	nethttp "net/http"
	"net/url"
	"strconv"
)

// NOTE: validate that is matches interface in compile-time
var _ keeper.Adapter = (*Adapter)(nil)

type Adapter struct {
	httpClient *http.Client
}

func NewAdapter(httpClient *http.Client) *Adapter {
	return &Adapter{httpClient: httpClient}
}

func (a *Adapter) List(request keeper.Request) ([]keeper.Response, int, error) {
	var overview []keeper.Response

	response, status, err := http.NewJSONRequest[cluster](a.httpClient).Get(keeper.Map(request, "/cluster"))
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
		parsed, errParse := url.Parse(patroniInstance.ApiUrl)
		if errParse != nil {
			return nil, nethttp.StatusBadRequest, errParse
		}
		host := parsed.Hostname()
		port, errCast := strconv.Atoi(parsed.Port())
		if errCast != nil {
			return nil, nethttp.StatusBadRequest, errCast
		}

		overview = append(overview, keeper.Response{
			Key:                  &patroniInstance.Name,
			Status:               &keeperStatus,
			State:                a.mapState(patroniInstance.State),
			Role:                 a.mapRole(patroniInstance.Role),
			Sync:                 a.mapSync(patroniInstance.Role),
			Lag:                  a.mapLag(patroniInstance.Lag),
			PendingRestart:       patroniInstance.PendingRestart,
			ScheduledRestart:     a.mapRestart(patroniInstance.ScheduledRestart),
			ScheduledSwitchover:  a.mapSwitchover(host, response.ScheduledSwitchover),
			Tags:                 patroniInstance.Tags,
			DiscoveredHost:       &patroniInstance.Host,
			DiscoveredName:       &patroniInstance.Name,
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

func (a *Adapter) mapRestart(restart *scheduledRestart) *keeper.ScheduledRestart {
	var scheduledRestart *keeper.ScheduledRestart
	if restart != nil {
		scheduledRestart = &keeper.ScheduledRestart{
			PendingRestart: restart.RestartPending,
			At:             restart.Schedule,
		}
	}
	return scheduledRestart
}

func (a *Adapter) mapSwitchover(host string, switchover *scheduledSwitchover) *keeper.ScheduledSwitchover {
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

// mapState normalizes Patroni's member "state" string onto keeper.State.
// Patroni has renamed/added state names across releases without preserving
// backward compatibility (e.g. older releases report a streaming replica as
// "running", newer ones report "streaming"), so any state Ivory hasn't seen
// before must fall back to StateUnknown rather than being passed through
// verbatim. Values taken from patroni/postgresql/__init__.py's STATE_*
// constants.
func (a *Adapter) mapState(state string) keeper.State {
	switch state {
	case "running", "streaming", "in archive recovery":
		return keeper.StateRunning
	case "starting", "initializing new cluster", "running custom bootstrap script", "creating replica":
		return keeper.StateStarting
	case "restarting":
		return keeper.StateRestarting
	case "stopping":
		return keeper.StateStopping
	case "stopped":
		return keeper.StateStopped
	case "stop failed", "start failed", "restart failed", "crashed", "initdb failed", "custom bootstrap failed":
		return keeper.StateFailed
	default:
		return keeper.StateUnknown
	}
}

// mapRole normalizes Patroni's member "role" string onto keeper.Role.
// standby_leader is the local leader of a standby cluster (cascading from an
// external primary), so it maps to Leader like a regular primary. sync_standby
// and quorum_standby are still replicas from a role standpoint - whether they
// belong to the synchronous replication set is reported separately by mapSync.
func (a *Adapter) mapRole(role string) keeper.Role {
	switch role {
	case "leader", "master", "standby_leader":
		return keeper.Leader
	case "replica", "sync_standby", "quorum_standby":
		return keeper.Replica
	default:
		return keeper.Unknown
	}
}

// mapSync reports whether a replica currently belongs to Patroni's
// synchronous replication set (synchronous_mode / synchronous_mode: quorum),
// as opposed to a plain asynchronous replica.
func (a *Adapter) mapSync(role string) bool {
	switch role {
	case "sync_standby", "quorum_standby":
		return true
	default:
		return false
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
	request.Body = configPause{Pause: false}
	return http.NewJSONRequest[string](a.httpClient).Patch(keeper.Map(request, "/config"))
}

func (a *Adapter) Pause(request keeper.Request) (*string, int, error) {
	if request.Body != nil {
		return nil, nethttp.StatusBadRequest, keeper.ErrBodyShouldBeEmpty
	}
	request.Body = configPause{Pause: true}
	return http.NewJSONRequest[string](a.httpClient).Patch(keeper.Map(request, "/config"))
}
