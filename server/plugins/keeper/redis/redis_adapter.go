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
var _ keeper.Adapter = (*Adapter)(nil)

// Adapter talks to redis directly, the same way native postgres does: there
// is no separate orchestrator, the keeper connection host/port is the redis
// host/port (keeperPort == dbPort convention), and the keeper vault holds
// redis AUTH credentials. Operations that require orchestration across nodes
// (a Sentinel/Cluster-style coordinator Ivory does not run) are not
// supported and excluded from SupportedFeatures.
type Adapter struct{}

func NewAdapter() *Adapter {
	return &Adapter{}
}

func (a *Adapter) List(request keeper.Request) ([]keeper.Response, int, error) {
	client, err := a.connect(request)
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

	response := mapNode(request.Host, request.Port, parseInfo(info))
	return []keeper.Response{response}, http.StatusOK, nil
}

func (a *Adapter) Config(request keeper.Request) (any, int, error) {
	client, err := a.connect(request)
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

func (a *Adapter) Failover(request keeper.Request) (*string, int, error) {
	client, err := a.connect(request)
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

// Reload has no redis equivalent: CONFIG REWRITE only persists the current
// in-memory config back to redis.conf, it does not reload settings from
// disk, so there is nothing to map this onto.
func (a *Adapter) Reload(keeper.Request) (*string, int, error) {
	return nil, http.StatusNotImplemented, keeper.ErrNotSupported
}

func (a *Adapter) Activate(keeper.Request) (*string, int, error) {
	return nil, http.StatusNotImplemented, keeper.ErrNotSupported
}

func (a *Adapter) Pause(keeper.Request) (*string, int, error) {
	return nil, http.StatusNotImplemented, keeper.ErrNotSupported
}

func (a *Adapter) connect(request keeper.Request) (*redis.Client, error) {
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
		DiscoveredHost:       &host,
		DiscoveredKeeperPort: &port,
		DiscoveredDbPort:     &port,
	}
}
