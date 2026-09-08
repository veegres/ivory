package redis

import (
	"context"
	"errors"
	"ivory/clients/redis"
	"ivory/plugins/keeper"
	"net/http"
	"strconv"
	"time"
)

var ErrNodeIsLeader = errors.New("failover target must be a replica; this node is the leader")
var ErrCredentialsRequired = errors.New("redis requires keeper credentials; add database password to keeper vault and configure it in your cluster")

// requestTimeout bounds every keeper operation (connect + command), so an
// unreachable node or a non-redis port cannot hang the cluster overview.
const requestTimeout = 5 * time.Second

// NOTE: validate that is matches interface in compile-time
var _ keeper.Adapter = (*Plugin)(nil)

// Plugin talks to redis directly, the same way native postgres does: there
// is no separate orchestrator, the keeper connection host/port is the redis
// host/port (keeperPort == dbPort convention), and the keeper vault holds
// redis AUTH credentials. Operations that require orchestration across nodes
// (a Sentinel/Cluster-style coordinator Ivory does not run) are not
// supported and excluded from SupportedFeatures.
type Plugin struct{}

func NewPlugin() *Plugin {
	return &Plugin{}
}

func (p *Plugin) List(request keeper.Request) ([]keeper.Response, int, error) {
	client, err := p.connect(request)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}
	defer redis.Close(client)

	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	// NOTE: INFO with no section returns redis' default section set, which
	// carries the server and memory fields the node's tags report alongside the
	// replication ones its role is read from - one round trip on every redis
	// version, unlike naming several sections, which only redis 7 accepts.
	info, errInfo := client.Info(ctx).Result()
	if errInfo != nil {
		return nil, http.StatusBadRequest, errInfo
	}

	fields := parseInfo(info)
	// NOTE: the cluster first - who else is in it - then this node's own state.
	// A node borrowed from another redis is contradicted here, because the
	// master it names is one nobody configured.
	members := mapMembers(fields)
	response := mapNode(request.Host, request.Port, fields)
	return append([]keeper.Response{response}, members...), http.StatusOK, nil
}

// mapMembers reports the other end of this node's replication: the master a
// replica follows, and deliberately nothing in the other direction. A master
// does list its replicas, but the ip on those slaveN: lines is the source
// address it observed the connection arrive from - ::1 under the single-host
// template's --network host, the docker gateway under bridge networking - never
// an address a replica listens on, so every one of them came back as a node
// nobody configured. Postgres answers the same problem the same way, keying
// standbys by application_name rather than client_addr; redis has no such
// field, so this direction is the connected_slaves count in mapTags instead.
func mapMembers(fields map[string]string) []keeper.Response {
	if master, ok := mapMaster(fields); ok {
		return []keeper.Response{master}
	}
	return nil
}

// mapMaster reports the server this replica follows, taken from master_host/
// master_port - the address out of its own replicaof, so it is what Ivory would
// connect to rather than something observed off a socket.
//
// It claims no role. A replica usually follows the master, but redis allows a
// replica of a replica, so Leader would be a guess and a wrong one puts a second
// leader on the overview. The node's own poll reports its real role anyway; this
// response exists to say the node is there at all.
func mapMaster(fields map[string]string) (keeper.Response, bool) {
	host := fields["master_host"]
	port, err := strconv.Atoi(fields["master_port"])
	if host == "" || err != nil || port <= 0 {
		return keeper.Response{}, false
	}
	key := host + ":" + strconv.Itoa(port)
	return keeper.Response{
		Key:                  &key,
		State:                keeper.StateUnknown,
		Role:                 keeper.Unknown,
		Lag:                  -1,
		DiscoveredHost:       &host,
		DiscoveredKeeperPort: &port,
		DiscoveredDbPort:     &port,
	}, true
}

func (p *Plugin) Config(request keeper.Request) (any, int, error) {
	client, err := p.connect(request)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}
	defer redis.Close(client)

	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	settings, errConfig := client.ConfigGet(ctx, "*").Result()
	if errConfig != nil {
		return nil, http.StatusBadRequest, errConfig
	}
	return settings, http.StatusOK, nil
}

func (p *Plugin) Failover(request keeper.Request) (*string, int, error) {
	client, err := p.connect(request)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}
	defer redis.Close(client)

	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	info, errInfo := client.Info(ctx, "replication").Result()
	if errInfo != nil {
		return nil, http.StatusBadRequest, errInfo
	}
	fields := parseInfo(info)
	if fields["role"] != "slave" {
		return nil, http.StatusBadRequest, ErrNodeIsLeader
	}

	if errPromote := client.ReplicaOf(ctx, "NO", "ONE").Err(); errPromote != nil {
		return nil, http.StatusBadRequest, errPromote
	}
	response := "promoted to master"
	return &response, http.StatusOK, nil
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

// Reload has no redis equivalent: CONFIG REWRITE only persists the current
// in-memory config back to redis.conf, it does not reload settings from
// disk, so there is nothing to map this onto.
func (p *Plugin) Reload(keeper.Request) (*string, int, error) {
	return nil, http.StatusNotImplemented, keeper.ErrNotSupported
}

func (p *Plugin) Activate(keeper.Request) (*string, int, error) {
	return nil, http.StatusNotImplemented, keeper.ErrNotSupported
}

func (p *Plugin) Pause(keeper.Request) (*string, int, error) {
	return nil, http.StatusNotImplemented, keeper.ErrNotSupported
}

func (p *Plugin) connect(request keeper.Request) (*redis.Client, error) {
	var username, password string
	if request.Credentials != nil {
		username = request.Credentials.Username
		password = request.Credentials.Password
	}
	if password == "" {
		return nil, ErrCredentialsRequired
	}
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()
	client, _, err := redis.Connect(ctx, redis.Config{
		Host:     request.Host,
		Port:     request.Port,
		Username: username,
		Password: password,
		TLS:      request.TlsConfig,
	})
	return client, err
}

// parseInfo parses a redis INFO section's "key:value\r\n" body into a map;
// comment lines (starting with "#") and blank lines are skipped.
func parseInfo(info string) map[string]string {
	fields := map[string]string{}
	start := 0
	for i := 0; i <= len(info); i++ {
		if i < len(info) && info[i] != '\n' {
			continue
		}
		line := info[start:i]
		start = i + 1
		line = trimCR(line)
		if line == "" || line[0] == '#' {
			continue
		}
		for j := 0; j < len(line); j++ {
			if line[j] == ':' {
				fields[line[:j]] = line[j+1:]
				break
			}
		}
	}
	return fields
}

func trimCR(s string) string {
	if len(s) > 0 && s[len(s)-1] == '\r' {
		return s[:len(s)-1]
	}
	return s
}

// mapNode builds the Response for the node whose INFO replication section
// was just read. lastIOSeconds (master_last_io_seconds_ago) is reported as
// Lag for a replica - unlike patroni's own lag value or postgres' WAL byte
// diff, it measures time since the replica last heard from its master, since
// a replica has no cheap way to learn how many bytes it is behind without
// also querying the master.
func mapNode(host string, port int, fields map[string]string) keeper.Response {
	role := keeper.Leader
	var lag int64
	if fields["role"] == "slave" {
		role = keeper.Replica
		if seconds, err := strconv.ParseInt(fields["master_last_io_seconds_ago"], 10, 64); err == nil && seconds > 0 {
			lag = seconds
		}
	}
	var status keeper.Status = keeper.Active
	key := host + ":" + strconv.Itoa(port)
	return keeper.Response{
		Key:                  &key,
		Status:               &status,
		State:                keeper.StateRunning,
		Role:                 role,
		Lag:                  lag,
		Tags:                 mapTags(fields, role),
		DiscoveredHost:       &host,
		DiscoveredKeeperPort: &port,
		DiscoveredDbPort:     &port,
	}
}

// mapTags reports what INFO carries that Role/State/Lag cannot say. A replica
// names the master it actually follows and whether that link is up: Lag only
// measures how long ago the last byte arrived, so a replica whose link went
// down seconds ago and one that was never able to connect are indistinguishable
// without it. A master names how many replicas are attached, which is the only
// place a replica missing from the cluster shows up at all. Memory is here
// because reaching maxmemory is what turns a healthy-looking redis into one
// refusing writes.
func mapTags(fields map[string]string, role keeper.Role) *map[string]any {
	tags := map[string]any{}
	if version := fields["redis_version"]; version != "" {
		tags["version"] = version
	}
	if memory := fields["used_memory_human"]; memory != "" {
		tags["memory"] = memory
	}
	if role == keeper.Replica {
		if link := fields["master_link_status"]; link != "" {
			tags["link"] = link
		}
		if masterHost := fields["master_host"]; masterHost != "" {
			tags["master"] = masterHost + ":" + fields["master_port"]
		}
	} else if replicas := fields["connected_slaves"]; replicas != "" {
		tags["replicas"] = replicas
	}
	if len(tags) == 0 {
		return nil
	}
	return &tags
}
