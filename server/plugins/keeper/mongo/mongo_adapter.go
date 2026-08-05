package mongo

import (
	"context"
	"errors"
	mongoclient "ivory/clients/mongo"
	"ivory/plugins/keeper"
	"net"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

var ErrNodeIsNotLeader = errors.New("failover target must be the current leader; mongo has no candidate-targeted step-down, connect to the primary to step it down")
var ErrConfigBodyRequired = errors.New("replica set reconfiguration requires a body with the new config document")

// requestTimeout bounds every keeper operation (connect + command), so an
// unreachable node or a non-mongo port cannot hang the cluster overview.
const requestTimeout = 5 * time.Second

const adminDb = "admin"

// NOTE: validate that is matches interface in compile-time
var _ keeper.Adapter = (*Adapter)(nil)

// Adapter talks to mongod's own replica set admin surface directly, the same
// way native etcd/zookeeper do: there is no separate orchestrator, and the
// keeper connection host/port is mongod's own port (keeperPort == dbPort
// convention). Automatic election is native to the replica set protocol;
// Ivory only exposes replSetStepDown as a manual trigger (Failover) - there
// is no candidate-targeted equivalent of Patroni's Switchover, since
// replSetStepDown only ever steps the current primary down and lets the set
// elect whichever secondary has the best position, it cannot name a specific
// target the way Patroni's DCS-mediated switchover can.
type Adapter struct{}

func NewAdapter() *Adapter {
	return &Adapter{}
}

func (a *Adapter) List(request keeper.Request) ([]keeper.Response, int, error) {
	client, err := a.connect(request)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}
	defer mongoclient.Close(client)

	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	status, errStatus := a.replSetGetStatus(ctx, client)
	if errStatus != nil {
		return nil, http.StatusBadRequest, errStatus
	}
	return mapStatus(status), http.StatusOK, nil
}

func (a *Adapter) Config(request keeper.Request) (any, int, error) {
	client, err := a.connect(request)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}
	defer mongoclient.Close(client)

	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	var result bson.M
	errCmd := client.Database(adminDb).RunCommand(ctx, bson.D{{Key: "replSetGetConfig", Value: 1}}).Decode(&result)
	if errCmd != nil {
		return nil, http.StatusBadRequest, errCmd
	}
	return result["config"], http.StatusOK, nil
}

// ConfigUpdate runs replSetReconfig with request.Body as the new replica set
// config document (the same shape replSetGetConfig's own "config" field
// returns), the real mongo primitive for changing member list/priorities/
// votes - unlike Patroni's DCS-stored dynamic config, this always requires
// the caller to submit the whole document, not a partial patch.
func (a *Adapter) ConfigUpdate(request keeper.Request) (any, int, error) {
	if request.Body == nil {
		return nil, http.StatusBadRequest, ErrConfigBodyRequired
	}
	client, err := a.connect(request)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}
	defer mongoclient.Close(client)

	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	var result bson.M
	errCmd := client.Database(adminDb).RunCommand(ctx, bson.D{{Key: "replSetReconfig", Value: request.Body}}).Decode(&result)
	if errCmd != nil {
		return nil, http.StatusBadRequest, errCmd
	}
	response := "reconfigured"
	return &response, http.StatusOK, nil
}

func (a *Adapter) DeleteSwitchover(keeper.Request) (*string, int, error) {
	return nil, http.StatusNotImplemented, keeper.ErrNotSupported
}

func (a *Adapter) Switchover(keeper.Request) (*string, int, error) {
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

// Reload has no mongo equivalent: setParameter changes individual runtime
// parameters one at a time, it does not reload a config file the way
// postgres' pg_reload_conf() does.
func (a *Adapter) Reload(keeper.Request) (*string, int, error) {
	return nil, http.StatusNotImplemented, keeper.ErrNotSupported
}

// Failover runs replSetStepDown against the node the request targets, which
// must currently be the primary - unlike redis/postgres, where Failover
// connects to a replica and promotes it directly, mongo has no "promote
// yourself" command on a secondary; only the primary can step itself down
// and let the set elect a replacement.
//
// NOTE: a successful step-down makes mongod immediately close its client
// connections, so the RunCommand call below can legitimately return a
// network-level error even though the step-down itself succeeded. That
// error is still returned as-is rather than guessed away, since the same
// error shape also covers genuine failures (e.g. no secondary is caught up
// enough to take over).
func (a *Adapter) Failover(request keeper.Request) (*string, int, error) {
	client, err := a.connect(request)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}
	defer mongoclient.Close(client)

	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	status, errStatus := a.replSetGetStatus(ctx, client)
	if errStatus != nil {
		return nil, http.StatusBadRequest, errStatus
	}
	self, ok := selfMember(status)
	if !ok || self.StateStr != "PRIMARY" {
		return nil, http.StatusBadRequest, ErrNodeIsNotLeader
	}

	errStepDown := client.Database(adminDb).RunCommand(ctx, bson.D{
		{Key: "replSetStepDown", Value: int32(60)},
		{Key: "secondaryCatchUpPeriodSecs", Value: int32(10)},
	}).Err()
	if errStepDown != nil {
		return nil, http.StatusBadRequest, errStepDown
	}
	response := "stepped down"
	return &response, http.StatusOK, nil
}

func (a *Adapter) Activate(keeper.Request) (*string, int, error) {
	return nil, http.StatusNotImplemented, keeper.ErrNotSupported
}

func (a *Adapter) Pause(keeper.Request) (*string, int, error) {
	return nil, http.StatusNotImplemented, keeper.ErrNotSupported
}

func (a *Adapter) connect(request keeper.Request) (*mongoclient.Client, error) {
	var username, password string
	if request.Credentials != nil {
		username = request.Credentials.Username
		password = request.Credentials.Password
	}
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()
	client, _, err := mongoclient.Connect(ctx, mongoclient.Config{
		Host:     request.Host,
		Port:     request.Port,
		Username: username,
		Password: password,
		TLS:      request.TlsConfig,
	})
	return client, err
}

func (a *Adapter) replSetGetStatus(ctx context.Context, client *mongoclient.Client) (*replSetStatus, error) {
	var status replSetStatus
	err := client.Database(adminDb).RunCommand(ctx, bson.D{{Key: "replSetGetStatus", Value: 1}}).Decode(&status)
	if err != nil {
		return nil, err
	}
	return &status, nil
}

func selfMember(status *replSetStatus) (replSetMember, bool) {
	for _, m := range status.Members {
		if m.Self {
			return m, true
		}
	}
	return replSetMember{}, false
}

// primaryOptime finds the current primary's optimeDate, used to compute
// every secondary's lag - mongo reports no ready-made per-member lag value
// the way patroni's /cluster does, only the raw optime each member last
// applied.
func primaryOptime(status *replSetStatus) (time.Time, bool) {
	for _, m := range status.Members {
		if m.StateStr == "PRIMARY" {
			return m.OptimeDate, true
		}
	}
	return time.Time{}, false
}

// mapStatus turns one replSetGetStatus reply into a Response per member,
// the same "whole cluster view from one connection" shape Patroni's
// /cluster and etcd's member list already return.
func mapStatus(status *replSetStatus) []keeper.Response {
	primaryTime, havePrimary := primaryOptime(status)

	responses := make([]keeper.Response, 0, len(status.Members))
	for _, m := range status.Members {
		responses = append(responses, mapMember(m, primaryTime, havePrimary))
	}
	return responses
}

func mapMember(m replSetMember, primaryTime time.Time, havePrimary bool) keeper.Response {
	role, state := mapRoleState(m.StateStr, m.Health)

	var lag int64
	if role == keeper.Replica && havePrimary {
		if diff := primaryTime.Sub(m.OptimeDate).Seconds(); diff > 0 {
			lag = int64(diff)
		}
	}

	var status keeper.Status = keeper.Active
	key := m.Name
	response := keeper.Response{
		Key:    &key,
		Status: &status,
		State:  state,
		Role:   role,
		Lag:    lag,
	}

	host, port, errSplit := net.SplitHostPort(m.Name)
	if errSplit == nil {
		portNum, errPort := net.LookupPort("tcp", port)
		if errPort == nil {
			response.DiscoveredHost = &host
			response.DiscoveredKeeperPort = &portNum
			response.DiscoveredDbPort = &portNum
		}
	}
	return response
}

// mapRoleState maps mongo's own member stateStr (see the replica set member
// states reference) onto keeper.Role/keeper.State. A member reporting
// health 0 is unreachable regardless of its last-known stateStr.
func mapRoleState(stateStr string, health float64) (keeper.Role, keeper.State) {
	if health == 0 {
		return keeper.Unknown, keeper.StateUnreachable
	}
	switch stateStr {
	case "PRIMARY":
		return keeper.Leader, keeper.StateRunning
	case "SECONDARY":
		return keeper.Replica, keeper.StateRunning
	case "RECOVERING", "ROLLBACK":
		return keeper.Replica, keeper.StateRestarting
	case "STARTUP", "STARTUP2":
		return keeper.Unknown, keeper.StateStarting
	case "ARBITER":
		return keeper.Unknown, keeper.StateRunning
	case "DOWN":
		return keeper.Unknown, keeper.StateUnreachable
	case "REMOVED":
		return keeper.Unknown, keeper.StateStopped
	default:
		return keeper.Unknown, keeper.StateUnknown
	}
}
