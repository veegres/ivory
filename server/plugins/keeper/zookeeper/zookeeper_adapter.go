package zookeeper

import (
	"context"
	"errors"
	"ivory/clients/zookeeper"
	"ivory/plugins/keeper"
	"net/http"
	"sort"
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
	config := zookeeper.Config{Host: request.Host, Port: request.Port}

	// NOTE: the cluster first, and on its own round trip: mntr describes nothing
	// but the node answering, so the ensemble has to come out of conf. It is
	// best-effort - a whitelist that omits "conf" must cost the membership check
	// rather than the node.
	var ensemble map[string]string
	if conf, errConf := zookeeper.FourLetterCommand(ctx, config, "conf"); errConf == nil {
		ensemble = parseLines(conf, "=")
	}

	output, err := zookeeper.FourLetterCommand(ctx, config, "mntr")
	if err != nil {
		return nil, http.StatusBadRequest, err
	}
	fields := parseLines(output, "\t")
	state, ok := fields["zk_server_state"]
	if !ok {
		return nil, http.StatusBadRequest, ErrCommandNotAvailable
	}

	response := mapNode(request.Host, request.Port, state, fields)
	response.DiscoveredCluster = ensembleIdentity(ensemble)
	return append([]keeper.Response{response}, mapEnsembleMembers(ensemble)...), http.StatusOK, nil
}

// mapEnsembleMembers reports the other servers of this node's ensemble, one per
// server.N line, so a single node can answer for the cluster the way patroni's
// /cluster does.
//
// Only a member whose line carries a client port can be reported, because that
// is the port Ivory reaches a zookeeper on and the rest of the line is the
// quorum and election pair, which it does not. The multi-host template writes
// ";2181" and so enumerates; the single-host one leaves it off and gives each
// node its own clientPort instead, so there it reports no members at all and
// ensembleIdentity is what catches a node from elsewhere. Guessing the missing
// port would invent nodes matching no configuration.
//
// A member states no role and no state: nothing contacted it, and which one
// holds the election is not in this node's conf. Every configured node reports
// both from its own mntr when Ivory polls it.
func mapEnsembleMembers(settings map[string]string) []keeper.Response {
	self := "server." + settings["serverId"]
	members := make([]keeper.Response, 0)
	for _, key := range sortedServerKeys(settings) {
		if key == self {
			continue
		}
		host, port, ok := memberEndpoint(settings[key])
		if !ok {
			continue
		}
		key := host + ":" + strconv.Itoa(port)
		members = append(members, keeper.Response{
			Key:                  &key,
			State:                keeper.StateUnknown,
			Role:                 keeper.Unknown,
			Lag:                  -1,
			DiscoveredHost:       &host,
			DiscoveredKeeperPort: &port,
			DiscoveredDbPort:     &port,
		})
	}
	if len(members) == 0 {
		return nil
	}
	return members
}

// memberEndpoint reads the client endpoint out of a
// "host:2888:3888:participant;0.0.0.0:2181" server entry. Everything before the
// ";" is the quorum and election pair, which Ivory never connects to; the part
// after it is the client address, written either as a bare port or as
// address:port. An entry without it is not addressable and reports nothing.
func memberEndpoint(server string) (string, int, bool) {
	host, client, found := strings.Cut(server, ";")
	if !found {
		return "", 0, false
	}
	if i := strings.IndexByte(host, ':'); i >= 0 {
		host = host[:i]
	}
	if i := strings.LastIndexByte(client, ':'); i >= 0 {
		client = client[i+1:]
	}
	port, err := strconv.Atoi(client)
	if host == "" || err != nil || port <= 0 {
		return "", 0, false
	}
	return host, port, true
}

func sortedServerKeys(settings map[string]string) []string {
	keys := make([]string, 0)
	for key := range settings {
		if strings.HasPrefix(key, "server.") {
			keys = append(keys, key)
		}
	}
	sort.Strings(keys)
	return keys
}

// ensembleIdentity reports the ensemble a node belongs to as the sorted list of
// its members' quorum endpoints. Every node of one ensemble is configured with
// the same server.N lines and a node from another ensemble never is, so where
// mapEnsembleMembers cannot address the members - a single-host ensemble, whose
// lines carry no client port - this is still enough to tell a stranger from a
// member. It costs nothing extra: the conf it reads was already fetched.
func ensembleIdentity(settings map[string]string) *string {
	endpoints := make([]string, 0)
	for _, key := range sortedServerKeys(settings) {
		endpoints = append(endpoints, quorumEndpoint(settings[key]))
	}
	if len(endpoints) == 0 {
		return nil
	}
	sort.Strings(endpoints)
	identity := strings.Join(endpoints, ",")
	return &identity
}

// quorumEndpoint trims a "host:2888:3888:participant;0.0.0.0:2181" entry down
// to the host:quorumPort every member writes identically. The role and the
// client-address suffix are dropped because a node states them for itself -
// zookeeper writes its own line with 0.0.0.0 - so keeping them would make one
// ensemble look like as many ensembles as it has members.
func quorumEndpoint(server string) string {
	if i := strings.IndexByte(server, ';'); i >= 0 {
		server = server[:i]
	}
	parts := strings.Split(server, ":")
	if len(parts) < 2 {
		return server
	}
	return parts[0] + ":" + parts[1]
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
