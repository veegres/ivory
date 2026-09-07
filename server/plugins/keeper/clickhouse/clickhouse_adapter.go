package clickhouse

import (
	"context"
	"errors"
	"fmt"
	"ivory/clients/clickhouse"
	"ivory/plugins/keeper"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
)

var ErrCredentialsRequired = errors.New("clickhouse requires keeper credentials; add database password to keeper vault and configure it in your cluster")

// requestTimeout bounds every keeper operation (connect + query), so an
// unreachable node or a non-clickhouse port cannot hang the cluster overview.
const requestTimeout = 5 * time.Second

// NOTE: validate that is matches interface in compile-time
var _ keeper.Adapter = (*Plugin)(nil)

// Plugin talks to clickhouse-server directly, the same way native postgres
// does: there is no separate orchestrator, the keeper connection host/port
// is the clickhouse native-protocol host/port (keeperPort == dbPort
// convention), and the keeper vault holds clickhouse credentials. ClickHouse
// has no single-primary replication model - every replica accepts writes
// and coordinates through ClickHouse Keeper/ZooKeeper - so there is no
// leader to switch or fail over to, and Switchover/Failover are excluded
// from SupportedFeatures. List reports every reachable node as Replica,
// which is what each one is; Unknown would claim Ivory could not tell.
type Plugin struct{}

func NewPlugin() *Plugin {
	return &Plugin{}
}

func (p *Plugin) List(request keeper.Request) ([]keeper.Response, int, error) {
	conn, err := p.connect(request)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}
	defer clickhouse.Close(conn)

	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	// NOTE: system.replicas has one row per replicated table; a node with no
	// replicated tables at all returns zero rows, and ClickHouse's max()/min()
	// over zero rows is 0 for every column here - which is exactly "not
	// read-only, no delay, nothing to compare active replicas against", the
	// right default for a standalone node.
	//
	// active_replicas/total_replicas (via min()/max() across every replicated
	// table on this node) reports a peer that has lost its session with the
	// coordination store - a crashed or partitioned node. It is a distinct
	// signal from the replication queue below: confirmed live (see
	// clickhouse_adapter's connectivity investigation this session) by
	// stopping a healthy replica's container - active_replicas dropped inside
	// 5s - and separately by leaving a replica's own session healthy but its
	// interserver HTTP port unreachable, where active_replicas stayed
	// unchanged because the coordination-store session was never affected.
	var isReadonly uint8
	var absoluteDelay uint64
	var activeReplicas, totalReplicas uint32
	errReplicas := conn.QueryRow(ctx, `
		SELECT max(is_readonly), max(absolute_delay), min(active_replicas), max(total_replicas)
		FROM system.replicas`,
	).Scan(&isReadonly, &absoluteDelay, &activeReplicas, &totalReplicas)
	if errReplicas != nil {
		return nil, http.StatusBadRequest, errReplicas
	}

	// NOTE: last_exception is set only once a queued fetch/merge has actually
	// failed and been retried - never while one is merely in flight - so a
	// table still catching up after a fresh deploy reports nothing here. This
	// is what catches a replica whose coordination-store session is fine but
	// whose interserver data path is broken (the exact shape of the
	// unpublished-9009 bug this template shipped with): active_replicas alone
	// missed it entirely, because the session it measures was never affected.
	var stuckCount uint64
	var sampleError string
	errQueue := conn.QueryRow(ctx, `
		SELECT count(), any(last_exception) FROM system.replication_queue
		WHERE last_exception != ''`,
	).Scan(&stuckCount, &sampleError)
	if errQueue != nil {
		return nil, http.StatusBadRequest, errQueue
	}

	response := mapNode(request.Host, request.Port, isReadonly > 0, absoluteDelay)
	response.Warnings = replicaWarnings(activeReplicas, totalReplicas, stuckCount, sampleError)
	return []keeper.Response{response}, http.StatusOK, nil
}

// replicaWarnings turns the two independent connectivity signals system.replicas
// and system.replication_queue actually carry into the same short-sentence
// warnings the rest of Ivory's overview uses - a peer that dropped its
// coordination-store session, and a peer whose data never crosses despite that
// session being fine. Either can fire alone, and a healthy cluster reports
// neither.
func replicaWarnings(activeReplicas, totalReplicas uint32, stuckCount uint64, sampleError string) []string {
	var warnings []string
	if totalReplicas > 0 && activeReplicas < totalReplicas {
		warnings = append(warnings, fmt.Sprintf(
			"%d of %d cluster replicas have no active session with the coordination store",
			totalReplicas-activeReplicas, totalReplicas))
	}
	if stuckCount > 0 {
		warnings = append(warnings, fmt.Sprintf(
			"replication queue has %d stuck task(s): %s", stuckCount, firstLine(sampleError)))
	}
	return warnings
}

// firstLine trims a clickhouse exception down to the one line worth reading -
// last_exception carries its own multi-thousand-character C++ stack trace
// after ", Stack trace", confirmed live against a genuinely stuck fetch (see
// List's doc), and neither belongs in a one-line cluster warning.
func firstLine(text string) string {
	if i := strings.Index(text, ", Stack trace"); i >= 0 {
		text = text[:i]
	}
	if i := strings.IndexByte(text, '\n'); i >= 0 {
		text = text[:i]
	}
	return text
}

func (p *Plugin) Config(request keeper.Request) (any, int, error) {
	conn, err := p.connect(request)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}
	defer clickhouse.Close(conn)

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
func (p *Plugin) Reload(request keeper.Request) (*string, int, error) {
	conn, err := p.connect(request)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}
	defer clickhouse.Close(conn)

	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	if errExec := conn.Exec(ctx, `SYSTEM RELOAD CONFIG`); errExec != nil {
		return nil, http.StatusBadRequest, errExec
	}
	response := "reloaded"
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

func (p *Plugin) Failover(keeper.Request) (*string, int, error) {
	return nil, http.StatusNotImplemented, keeper.ErrNotSupported
}

func (p *Plugin) Activate(keeper.Request) (*string, int, error) {
	return nil, http.StatusNotImplemented, keeper.ErrNotSupported
}

func (p *Plugin) Pause(keeper.Request) (*string, int, error) {
	return nil, http.StatusNotImplemented, keeper.ErrNotSupported
}

func (p *Plugin) connect(request keeper.Request) (driver.Conn, error) {
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
	conn, _, err := clickhouse.Connect(ctx, clickhouse.Config{
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
