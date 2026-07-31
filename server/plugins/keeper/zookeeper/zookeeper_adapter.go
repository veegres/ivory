package zookeeper

import (
	"context"
	"errors"
	zkclient "ivory/clients/zookeeper"
	"ivory/plugins/keeper"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var ErrCommandNotAvailable = errors.New("zookeeper did not return the expected admin command output; ensure ZOO_4LW_COMMANDS_WHITELIST includes 'mntr' and 'conf'")

// requestTimeout bounds every keeper operation (connect + command), so an
// unreachable node or a non-zookeeper port cannot hang the cluster overview.
const requestTimeout = 5 * time.Second

// NOTE: validate that is matches interface in compile-time
var _ keeper.Adapter = (*Adapter)(nil)

// Adapter talks to zookeeper's own admin ("four-letter word") text protocol
// directly, the same way native etcd does: there is no separate
// orchestrator, and the keeper connection host/port is zookeeper's client
// port (keeperPort == dbPort convention). There are no credentials to
// authenticate with - four-letter commands have no auth concept at all,
// only an optional server-side command whitelist (ZOO_4LW_COMMANDS_WHITELIST,
// see zookeeper_metadata.go). Leader election is fully automatic (the ZAB
// protocol) with no manual trigger anywhere in the admin surface, so
// Switchover/Failover are not supported - unlike etcd, which can force a
// raft leader change via MoveLeader.
type Adapter struct{}

func NewAdapter() *Adapter {
	return &Adapter{}
}

func (a *Adapter) List(request keeper.Request) ([]keeper.Response, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	output, err := zkclient.FourLetterCommand(ctx, zkclient.Config{Host: request.Host, Port: request.Port}, "mntr")
	if err != nil {
		return nil, http.StatusBadRequest, err
	}
	fields := parseLines(output, "\t")
	state, ok := fields["zk_server_state"]
	if !ok {
		return nil, http.StatusBadRequest, ErrCommandNotAvailable
	}

	return []keeper.Response{mapNode(request.Host, request.Port, state)}, http.StatusOK, nil
}

func (a *Adapter) Config(request keeper.Request) (any, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	output, err := zkclient.FourLetterCommand(ctx, zkclient.Config{Host: request.Host, Port: request.Port}, "conf")
	if err != nil {
		return nil, http.StatusBadRequest, err
	}
	settings := parseLines(output, "=")
	if len(settings) == 0 {
		return nil, http.StatusBadRequest, ErrCommandNotAvailable
	}
	return settings, http.StatusOK, nil
}

func (a *Adapter) ConfigUpdate(keeper.Request) (any, int, error) {
	return nil, http.StatusNotImplemented, keeper.ErrNotSupported
}

func (a *Adapter) Switchover(keeper.Request) (*string, int, error) {
	return nil, http.StatusNotImplemented, keeper.ErrNotSupported
}

func (a *Adapter) DeleteSwitchover(keeper.Request) (*string, int, error) {
	return nil, http.StatusNotImplemented, keeper.ErrNotSupported
}

func (a *Adapter) Reinitialize(keeper.Request) (*string, int, error) {
	return nil, http.StatusNotImplemented, keeper.ErrNotSupported
}

func (a *Adapter) Restart(keeper.Request) (*string, int, error) {
	return nil, http.StatusNotImplemented, keeper.ErrNotSupported
}

func (a *Adapter) DeleteRestart(keeper.Request) (*string, int, error) {
	return nil, http.StatusNotImplemented, keeper.ErrNotSupported
}

// Reload has no clean zookeeper equivalent: the closest admin surface,
// "reconfig", changes ensemble membership, not general config, and is a
// stateful multi-step operation - not a simple reload.
func (a *Adapter) Reload(keeper.Request) (*string, int, error) {
	return nil, http.StatusNotImplemented, keeper.ErrNotSupported
}

func (a *Adapter) Failover(keeper.Request) (*string, int, error) {
	return nil, http.StatusNotImplemented, keeper.ErrNotSupported
}

func (a *Adapter) Activate(keeper.Request) (*string, int, error) {
	return nil, http.StatusNotImplemented, keeper.ErrNotSupported
}

func (a *Adapter) Pause(keeper.Request) (*string, int, error) {
	return nil, http.StatusNotImplemented, keeper.ErrNotSupported
}

// parseLines parses "key<sep>value" lines - mntr's tab-separated metrics,
// conf's `=`-separated settings - into a map. Lines without the separator
// (e.g. conf's bare "membership:" section header, or repeated "server.N=..."
// membership entries which conf also emits - these collide on no single key
// anyway) are skipped; conf's repeated server.N lines still parse fine here
// since N makes each key distinct.
func parseLines(output string, sep string) map[string]string {
	fields := map[string]string{}
	for _, line := range strings.Split(output, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		idx := strings.Index(line, sep)
		if idx < 0 {
			continue
		}
		fields[line[:idx]] = line[idx+len(sep):]
	}
	return fields
}

// mapNode maps zookeeper's own zk_server_state (leader/follower/observer/
// standalone) onto keeper.Role. Lag is left at 0: unlike etcd's raft index
// or postgres' WAL position, mntr reports no per-follower offset a member
// can read about itself - only leader-side aggregates like
// zk_synced_followers, which isn't a per-member lag figure.
func mapNode(host string, port int, state string) keeper.Response {
	var role keeper.Role = keeper.Unknown
	switch state {
	case "leader", "standalone":
		role = keeper.Leader
	case "follower", "observer":
		role = keeper.Replica
	}
	var status keeper.Status = keeper.Active
	key := host + ":" + strconv.Itoa(port)
	return keeper.Response{
		Key:                  &key,
		Status:               &status,
		State:                keeper.StateRunning,
		Role:                 role,
		DiscoveredHost:       &host,
		DiscoveredKeeperPort: &port,
		DiscoveredDbPort:     &port,
	}
}
