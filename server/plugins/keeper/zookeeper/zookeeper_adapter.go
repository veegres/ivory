package zookeeper

import (
	"context"
	"errors"
	"ivory/clients/zookeeper"
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
var _ keeper.Adapter = (*Plugin)(nil)

// Plugin talks to zookeeper's own admin ("four-letter word") text protocol
// directly, the same way native etcd does: there is no separate
// orchestrator, and the keeper connection host/port is zookeeper's client
// port (keeperPort == dbPort convention). There are no credentials to
// authenticate with - four-letter commands have no auth concept at all,
// only an optional server-side command whitelist (ZOO_4LW_COMMANDS_WHITELIST,
// see zookeeper_metadata.go). Leader election is fully automatic (the ZAB
// protocol) with no manual trigger anywhere in the admin surface, so
// Switchover/Failover are not supported - unlike etcd, which can force a
// raft leader change via MoveLeader.
type Plugin struct{}

func NewPlugin() *Plugin {
	return &Plugin{}
}

func (p *Plugin) List(request keeper.Request) ([]keeper.Response, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	output, err := zookeeper.FourLetterCommand(ctx, zookeeper.Config{Host: request.Host, Port: request.Port}, "mntr")
	if err != nil {
		return nil, http.StatusBadRequest, err
	}
	fields := parseLines(output, "\t")
	state, ok := fields["zk_server_state"]
	if !ok {
		return nil, http.StatusBadRequest, ErrCommandNotAvailable
	}

	return []keeper.Response{mapNode(request.Host, request.Port, state, fields)}, http.StatusOK, nil
}

func (p *Plugin) Config(request keeper.Request) (any, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	output, err := zookeeper.FourLetterCommand(ctx, zookeeper.Config{Host: request.Host, Port: request.Port}, "conf")
	if err != nil {
		return nil, http.StatusBadRequest, err
	}
	settings := parseLines(output, "=")
	if len(settings) == 0 {
		return nil, http.StatusBadRequest, ErrCommandNotAvailable
	}
	return settings, http.StatusOK, nil
}

func (p *Plugin) ConfigUpdate(keeper.Request) (any, int, error) {
	return nil, http.StatusNotImplemented, keeper.ErrNotSupported
}

func (p *Plugin) Switchover(keeper.Request) (*string, int, error) {
	return nil, http.StatusNotImplemented, keeper.ErrNotSupported
}

func (p *Plugin) DeleteSwitchover(keeper.Request) (*string, int, error) {
	return nil, http.StatusNotImplemented, keeper.ErrNotSupported
}

func (p *Plugin) Reinitialize(keeper.Request) (*string, int, error) {
	return nil, http.StatusNotImplemented, keeper.ErrNotSupported
}

func (p *Plugin) Restart(keeper.Request) (*string, int, error) {
	return nil, http.StatusNotImplemented, keeper.ErrNotSupported
}

func (p *Plugin) DeleteRestart(keeper.Request) (*string, int, error) {
	return nil, http.StatusNotImplemented, keeper.ErrNotSupported
}

// Reload has no clean zookeeper equivalent: the closest admin surface,
// "reconfig", changes ensemble membership, not general config, and is a
// stateful multi-step operation - not a simple reload.
func (p *Plugin) Reload(keeper.Request) (*string, int, error) {
	return nil, http.StatusNotImplemented, keeper.ErrNotSupported
}

func (p *Plugin) Failover(keeper.Request) (*string, int, error) {
	return nil, http.StatusNotImplemented, keeper.ErrNotSupported
}

func (p *Plugin) Activate(keeper.Request) (*string, int, error) {
	return nil, http.StatusNotImplemented, keeper.ErrNotSupported
}

func (p *Plugin) Pause(keeper.Request) (*string, int, error) {
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
func mapNode(host string, port int, state string, fields map[string]string) keeper.Response {
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
		Tags:                 mapTags(state, fields),
		DiscoveredHost:       &host,
		DiscoveredKeeperPort: &port,
		DiscoveredDbPort:     &port,
	}
}

// mapTags reports the mntr metrics that say something Role/State cannot. A
// standalone server is reported as a leader because it is the only node that
// accepts writes, and an observer as a replica because it never accepts them -
// but a standalone node that was meant to join an ensemble, and an observer
// that was meant to vote, are both configuration mistakes that read as a
// healthy cluster without serverState. The leader's follower counts are the one
// place a member that dropped out of the quorum is visible, since every
// zookeeper node otherwise only ever describes itself.
func mapTags(state string, fields map[string]string) *map[string]any {
	tags := map[string]any{}
	if state != "leader" && state != "follower" {
		tags["serverState"] = state
	}
	if version := fields["zk_version"]; version != "" {
		tags["version"] = shortVersion(version)
	}
	if znodes := fields["zk_znode_count"]; znodes != "" {
		tags["znodes"] = znodes
	}
	if connections := fields["zk_num_alive_connections"]; connections != "" {
		tags["connections"] = connections
	}
	if synced := fields["zk_synced_followers"]; synced != "" {
		if followers := fields["zk_followers"]; followers != "" {
			tags["syncedFollowers"] = synced + "/" + followers
		} else {
			tags["syncedFollowers"] = synced
		}
	}
	if len(tags) == 0 {
		return nil
	}
	return &tags
}

// shortVersion keeps the release number alone: zk_version reports the build
// hash and build date after it ("3.9.1-abc123, built on ..."), which says
// nothing a reader comparing one node against another needs.
func shortVersion(version string) string {
	if i := strings.IndexAny(version, "-,"); i >= 0 {
		return version[:i]
	}
	return version
}
