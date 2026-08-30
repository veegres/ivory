package clickhouse

import (
	"context"
	"errors"
	chclient "ivory/clients/clickhouse"
	"ivory/plugins/keeper"
	"net/http"
	"strconv"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
)

var ErrCredentialsRequired = errors.New("clickhouse requires keeper credentials; add database password to keeper vault and configure it in your cluster")

// requestTimeout bounds every keeper operation (connect + query), so an
// unreachable node or a non-clickhouse port cannot hang the cluster overview.
const requestTimeout = 5 * time.Second

// NOTE: validate that is matches interface in compile-time
var _ keeper.Adapter = (*Adapter)(nil)

// Adapter talks to clickhouse-server directly, the same way native postgres
// does: there is no separate orchestrator, the keeper connection host/port
// is the clickhouse native-protocol host/port (keeperPort == dbPort
// convention), and the keeper vault holds clickhouse credentials. ClickHouse
// has no single-primary replication model - every replica accepts writes
// and coordinates through ClickHouse Keeper/ZooKeeper - so there is no
// leader to switch or fail over to, and Switchover/Failover are excluded
// from SupportedFeatures. List reports every reachable node as Replica,
// which is what each one is; Unknown would claim Ivory could not tell.
type Adapter struct{}

func NewAdapter() *Adapter {
	return &Adapter{}
}

func (a *Adapter) List(request keeper.Request) ([]keeper.Response, int, error) {
	conn, err := a.connect(request)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}
	defer chclient.Close(conn)

	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	// NOTE: system.replicas has one row per replicated table; a node with no
	// replicated tables at all returns zero rows, and ClickHouse's max()
	// over zero rows is 0 for both columns - which is exactly "not
	// read-only, no delay", the right default for a standalone node.
	var isReadonly uint8
	var absoluteDelay uint64
	errQuery := conn.QueryRow(ctx, `SELECT max(is_readonly), max(absolute_delay) FROM system.replicas`).Scan(&isReadonly, &absoluteDelay)
	if errQuery != nil {
		return nil, http.StatusBadRequest, errQuery
	}

	return []keeper.Response{mapNode(request.Host, request.Port, isReadonly > 0, absoluteDelay)}, http.StatusOK, nil
}

func (a *Adapter) Config(request keeper.Request) (any, int, error) {
	conn, err := a.connect(request)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}
	defer chclient.Close(conn)

	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	rows, errQuery := conn.Query(ctx, `SELECT name, value FROM system.settings ORDER BY name`)
	if errQuery != nil {
		return nil, http.StatusBadRequest, errQuery
	}
	defer rows.Close()

	settings := map[string]string{}
	for rows.Next() {
		var name, value string
		if errScan := rows.Scan(&name, &value); errScan != nil {
			return nil, http.StatusBadRequest, errScan
		}
		settings[name] = value
	}
	if rows.Err() != nil {
		return nil, http.StatusBadRequest, rows.Err()
	}
	return settings, http.StatusOK, nil
}

// Reload runs SYSTEM RELOAD CONFIG, clickhouse's genuine equivalent of
// postgres' pg_reload_conf(): it re-reads config.xml without restarting.
func (a *Adapter) Reload(request keeper.Request) (*string, int, error) {
	conn, err := a.connect(request)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}
	defer chclient.Close(conn)

	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	if errExec := conn.Exec(ctx, `SYSTEM RELOAD CONFIG`); errExec != nil {
		return nil, http.StatusBadRequest, errExec
	}
	response := "reloaded"
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

func (a *Adapter) Failover(keeper.Request) (*string, int, error) {
	return nil, http.StatusNotImplemented, keeper.ErrNotSupported
}

func (a *Adapter) Activate(keeper.Request) (*string, int, error) {
	return nil, http.StatusNotImplemented, keeper.ErrNotSupported
}

func (a *Adapter) Pause(keeper.Request) (*string, int, error) {
	return nil, http.StatusNotImplemented, keeper.ErrNotSupported
}

func (a *Adapter) connect(request keeper.Request) (driver.Conn, error) {
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
	conn, _, err := chclient.Connect(ctx, chclient.Config{
		Host:     request.Host,
		Port:     request.Port,
		Username: username,
		Password: password,
		TLS:      request.TlsConfig,
	})
	return conn, err
}

func mapNode(host string, port int, readonly bool, absoluteDelay uint64) keeper.Response {
	state := keeper.StateRunning
	if readonly {
		state = keeper.StateStopping
	}
	var status keeper.Status = keeper.Active
	key := host + ":" + strconv.Itoa(port)
	return keeper.Response{
		Key:                  &key,
		Status:               &status,
		State:                state,
		Role:                 keeper.Replica,
		Lag:                  int64(absoluteDelay),
		DiscoveredHost:       &host,
		DiscoveredKeeperPort: &port,
		DiscoveredDbPort:     &port,
	}
}
