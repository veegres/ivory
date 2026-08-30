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
var _ keeper.Adapter = (*Plugin)(nil)

type Plugin struct {
	httpClient *http.Client
}

func NewPlugin(httpClient *http.Client) *Plugin {
	return &Plugin{httpClient: httpClient}
}

func (p *Plugin) List(request keeper.Request) ([]keeper.Response, int, error) {
	var overview []keeper.Response

	response, status, err := http.NewJSONRequest[cluster](p.httpClient).Get(keeper.Map(request, "/cluster"))
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
			State:                p.mapState(patroniInstance.State),
			Role:                 p.mapRole(patroniInstance.Role),
			Sync:                 p.mapSync(patroniInstance.Role),
			Lag:                  p.mapLag(patroniInstance.Lag),
			PendingRestart:       patroniInstance.PendingRestart,
			ScheduledRestart:     p.mapRestart(patroniInstance.ScheduledRestart),
			ScheduledSwitchover:  p.mapSwitchover(host, response.ScheduledSwitchover),
			Tags:                 patroniInstance.Tags,
			DiscoveredHost:       &patroniInstance.Host,
			DiscoveredName:       &patroniInstance.Name,
			DiscoveredDbPort:     &patroniInstance.Port,
			DiscoveredKeeperPort: &port,
		})
	}

	return overview, status, err
}

func (p *Plugin) mapLag(lag json.RawMessage) int64 {
	if num, errInt := strconv.ParseInt(string(lag), 10, 64); errInt == nil {
		return num
	}
	return -1
}

func (p *Plugin) mapRestart(restart *scheduledRestart) *keeper.ScheduledRestart {
	var scheduledRestart *keeper.ScheduledRestart
	if restart != nil {
		scheduledRestart = &keeper.ScheduledRestart{
			PendingRestart: restart.RestartPending,
			At:             restart.Schedule,
		}
	}
	return scheduledRestart
}

func (p *Plugin) mapSwitchover(host string, switchover *scheduledSwitchover) *keeper.ScheduledSwitchover {
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
func (p *Plugin) mapState(state string) keeper.State {
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
func (p *Plugin) mapRole(role string) keeper.Role {
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
func (p *Plugin) mapSync(role string) bool {
	switch role {
	case "sync_standby", "quorum_standby":
		return true
	default:
		return false
	}
}

func (p *Plugin) Config(request keeper.Request) (any, int, error) {
	return http.NewJSONRequest[any](p.httpClient).Get(keeper.Map(request, "/config"))
}

func (p *Plugin) ConfigUpdate(request keeper.Request) (any, int, error) {
	return http.NewJSONRequest[any](p.httpClient).Patch(keeper.Map(request, "/config"))
}

func (p *Plugin) Switchover(request keeper.Request) (*string, int, error) {
	return http.NewJSONRequest[string](p.httpClient).Post(keeper.Map(request, "/switchover"))
}

func (p *Plugin) DeleteSwitchover(request keeper.Request) (*string, int, error) {
	return http.NewJSONRequest[string](p.httpClient).Delete(keeper.Map(request, "/switchover"))
}

func (p *Plugin) Reinitialize(request keeper.Request) (*string, int, error) {
	return http.NewJSONRequest[string](p.httpClient).Post(keeper.Map(request, "/reinitialize"))
}

func (p *Plugin) Restart(request keeper.Request) (*string, int, error) {
	return http.NewJSONRequest[string](p.httpClient).Post(keeper.Map(request, "/restart"))
}

func (p *Plugin) DeleteRestart(request keeper.Request) (*string, int, error) {
	return http.NewJSONRequest[string](p.httpClient).Delete(keeper.Map(request, "/restart"))
}

func (p *Plugin) Reload(request keeper.Request) (*string, int, error) {
	return http.NewJSONRequest[string](p.httpClient).Post(keeper.Map(request, "/reload"))
}

func (p *Plugin) Failover(request keeper.Request) (*string, int, error) {
	return http.NewJSONRequest[string](p.httpClient).Post(keeper.Map(request, "/failover"))
}

func (p *Plugin) Activate(request keeper.Request) (*string, int, error) {
	if request.Body != nil {
		return nil, nethttp.StatusBadRequest, keeper.ErrBodyShouldBeEmpty
	}
	request.Body = configPause{Pause: false}
	return http.NewJSONRequest[string](p.httpClient).Patch(keeper.Map(request, "/config"))
}

func (p *Plugin) Pause(request keeper.Request) (*string, int, error) {
	if request.Body != nil {
		return nil, nethttp.StatusBadRequest, keeper.ErrBodyShouldBeEmpty
	}
	request.Body = configPause{Pause: true}
	return http.NewJSONRequest[string](p.httpClient).Patch(keeper.Map(request, "/config"))
}
