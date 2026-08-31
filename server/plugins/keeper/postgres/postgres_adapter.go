package postgres

import (
	"context"
	"errors"
	"ivory/clients/postgres"
	"ivory/plugins/keeper"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

var ErrNodeIsLeader = errors.New("failover target must be a replica; this node is the leader")
var ErrReloadFailed = errors.New("postgres refused to reload the configuration")
var ErrPromoteRefused = errors.New("postgres refused to start the promotion")
var ErrCredentialsRequired = errors.New("native postgres requires keeper credentials; add database password to keeper vault and configure it in your cluster")

// requestTimeout bounds every keeper operation (connect + query), so an
// unreachable node or a non-postgres port cannot hang the cluster overview.
const requestTimeout = 5 * time.Second

const configQuery = `SELECT name, setting FROM pg_settings ORDER BY name`

// listQuery detects the node role and, for replicas, the replication lag as
// the difference between received and replayed wal in bytes. NOTE: the lag
// unit differs from patroni, which reports its own /cluster lag value.
const listQuery = `SELECT pg_is_in_recovery(),
       CASE WHEN pg_is_in_recovery()
            THEN GREATEST(COALESCE(pg_wal_lsn_diff(pg_last_wal_receive_lsn(), pg_last_wal_replay_lsn()), 0), 0)::bigint
            ELSE 0 END`

// syncStandbyQuery reads the primary's live view of connected standbys.
// sync_state is only ever populated on the primary's own connection - a
// standby has no way to determine its own synchronous status by querying
// itself - so this is only run once listQuery has established the connected
// node is not in recovery (i.e. it is the cluster's primary).
//
// application_name, not client_addr, is used to identify which standby a row
// belongs to (see mapSyncStandby's doc for why), so rows with no
// application_name set can't be attributed to any node and are skipped.
const syncStandbyQuery = `SELECT application_name, sync_state FROM pg_stat_replication WHERE application_name IS NOT NULL AND application_name != ''`

// sqlStateCannotConnectNow is postgres' ERRCODE_CANNOT_CONNECT_NOW, returned
// for a new connection attempt while the database is starting up, shutting
// down, or still in crash/archive recovery. All three cases share this one
// SQLSTATE and are only distinguished by the FATAL message text, so this is
// a normal transient condition (e.g. right after a container restart), not a
// connectivity failure, and List reports it as a node state instead of an error.
const sqlStateCannotConnectNow = "57P03"

// NOTE: validate that is matches interface in compile-time
var _ keeper.Adapter = (*Plugin)(nil)

// Plugin provides the native postgres experience alongside patroni by
// talking to postgres directly. The keeper connection host/port is the
// postgres host/port (keeperPort == dbPort convention) and the keeper vault
// holds database credentials. Operations that require orchestration or OS
// access are not supported and excluded from SupportedFeatures - in
// particular, Ivory never configures streaming replication for native
// postgres (DeploymentSpec only deploys a bare postgres image), so
// identifying a connected standby's Sync status (see mapSyncStandby) depends
// entirely on the operator manually setting each standby's primary_conninfo
// application_name to match the Host Ivory has configured for it.
type Plugin struct{}

func NewPlugin() *Plugin {
	return &Plugin{}
}

func (p *Plugin) List(request keeper.Request) ([]keeper.Response, int, error) {
	var inRecovery bool
	var lag int64
	err := p.queryRow(request, listQuery, func(row pgx.Row) error {
		return row.Scan(&inRecovery, &lag)
	})
	if err != nil {
		if state, ok := mapUnavailableState(err); ok {
			// NOTE: err is still returned alongside the mapped response so
			// the caller keeps the reason as a warning instead of it being
			// silently discarded just because we could still report a state.
			return []keeper.Response{mapUnavailableNode(request.Host, request.Port, state)}, http.StatusServiceUnavailable, err
		}
		return nil, http.StatusBadRequest, err
	}
	responses := []keeper.Response{mapNode(request.Host, request.Port, inRecovery, lag)}
	if !inRecovery {
		responses = append(responses, p.listSyncStandbys(request)...)
	}
	return responses, http.StatusOK, nil
}

// listSyncStandbys reports every standby currently connected to this
// primary, so replicas queried independently (see List) can be told apart as
// synchronous vs asynchronous - something they cannot determine about
// themselves. It is best-effort: pg_stat_replication may be restricted for
// the configured credentials, or briefly fail during a topology change, in
// which case the primary's own List response above is still returned as-is.
func (p *Plugin) listSyncStandbys(request keeper.Request) []keeper.Response {
	var standbys []keeper.Response
	err := p.query(request, syncStandbyQuery, func(rows pgx.Rows) error {
		for rows.Next() {
			var applicationName, syncState string
			if errScan := rows.Scan(&applicationName, &syncState); errScan != nil {
				return errScan
			}
			// NOTE: pg_stat_replication only has the standby's ephemeral
			// replication client_port, not its listening port, so the
			// standby's real keeper/db port is left undiscovered here rather
			// than guessed - the discovery layer resolves this response to
			// the right configured node by name instead (see
			// cluster.mergeKeeperSync).
			standbys = append(standbys, mapSyncStandby(applicationName, syncState))
		}
		return rows.Err()
	})
	if err != nil {
		return nil
	}
	return standbys
}

// mapUnavailableState reports the keeper.State represented by a "cannot
// connect now" postgres error, disambiguating "starting up" from "shutting
// down" by message text since both share sqlStateCannotConnectNow. The
// second return value is false for any other error, including genuine
// connectivity failures (host unreachable, wrong port, auth failure).
func mapUnavailableState(err error) (keeper.State, bool) {
	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) || pgErr.Code != sqlStateCannotConnectNow {
		return "", false
	}
	if strings.Contains(pgErr.Message, "shutting down") {
		return keeper.StateStopping, true
	}
	return keeper.StateStarting, true
}

func (p *Plugin) Config(request keeper.Request) (any, int, error) {
	settings := map[string]string{}
	err := p.query(request, configQuery, func(rows pgx.Rows) error {
		for rows.Next() {
			var name, setting string
			if errScan := rows.Scan(&name, &setting); errScan != nil {
				return errScan
			}
			settings[name] = setting
		}
		return rows.Err()
	})
	if err != nil {
		return nil, http.StatusBadRequest, err
	}
	return settings, http.StatusOK, nil
}

func (p *Plugin) Reload(request keeper.Request) (*string, int, error) {
	var reloaded bool
	err := p.queryRow(request, "SELECT pg_reload_conf()", func(row pgx.Row) error {
		return row.Scan(&reloaded)
	})
	if err != nil {
		return nil, http.StatusBadRequest, err
	}
	if !reloaded {
		return nil, http.StatusInternalServerError, ErrReloadFailed
	}
	response := "reloaded"
	return &response, http.StatusOK, nil
}

func (p *Plugin) Failover(request keeper.Request) (*string, int, error) {
	var inRecovery bool
	err := p.queryRow(request, "SELECT pg_is_in_recovery()", func(row pgx.Row) error {
		return row.Scan(&inRecovery)
	})
	if err != nil {
		return nil, http.StatusBadRequest, err
	}
	if !inRecovery {
		return nil, http.StatusBadRequest, ErrNodeIsLeader
	}
	// NOTE: pg_promote(false) only signals the promotion and returns
	// immediately, so it fits into the request timeout
	var promoted bool
	errPromote := p.queryRow(request, "SELECT pg_promote(false)", func(row pgx.Row) error {
		return row.Scan(&promoted)
	})
	if errPromote != nil {
		return nil, http.StatusBadRequest, errPromote
	}
	if !promoted {
		return nil, http.StatusInternalServerError, ErrPromoteRefused
	}
	response := "promotion started"
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

func (p *Plugin) Activate(keeper.Request) (*string, int, error) {
	return nil, http.StatusNotImplemented, keeper.ErrNotSupported
}

func (p *Plugin) Pause(keeper.Request) (*string, int, error) {
	return nil, http.StatusNotImplemented, keeper.ErrNotSupported
}

func (p *Plugin) queryRow(request keeper.Request, query string, scan func(row pgx.Row) error) error {
	return p.withConnection(request, func(ctx context.Context, conn *pgx.Conn) error {
		return scan(conn.QueryRow(ctx, query))
	})
}

func (p *Plugin) query(request keeper.Request, query string, parse func(rows pgx.Rows) error) error {
	return p.withConnection(request, func(ctx context.Context, conn *pgx.Conn) error {
		rows, err := conn.Query(ctx, query)
		if err != nil {
			return err
		}
		defer rows.Close()
		return parse(rows)
	})
}

func (p *Plugin) withConnection(request keeper.Request, action func(ctx context.Context, conn *pgx.Conn) error) error {
	var username, password string
	if request.Credentials != nil {
		username = request.Credentials.Username
		password = request.Credentials.Password
	}
	if username == "" {
		return ErrCredentialsRequired
	}
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()
	conn, _, err := postgres.Connect(ctx, postgres.Config{
		Host:     request.Host,
		Port:     request.Port,
		Username: username,
		Password: password,
		AppName:  "ivory-keeper",
		TLS:      request.TlsConfig,
	})
	if err != nil {
		return err
	}
	defer postgres.Close(conn)
	return action(ctx, conn)
}

// mapUnavailableNode reports a reachable node whose postgres is not
// currently accepting queries (starting up or shutting down), so the
// cluster overview can show its real state instead of flagging it as
// unreachable while its role/lag are not yet knowable. List still returns
// the triggering error alongside this response so callers surface it as a
// warning instead of discarding it.
func mapUnavailableNode(host string, port int, state keeper.State) keeper.Response {
	var status keeper.Status = keeper.Active
	key := host + ":" + strconv.Itoa(port)
	return keeper.Response{
		Key:                  &key,
		Status:               &status,
		State:                state,
		Role:                 keeper.Unknown,
		DiscoveredHost:       &host,
		DiscoveredKeeperPort: &port,
		DiscoveredDbPort:     &port,
	}
}

// mapSyncStandby reports one attribute of another node, not a node. Sync state
// is visible from the primary alone - a standby cannot determine its own
// synchronous status by querying itself - so it is the one thing a node's own
// connection can never report, and the only thing this response claims.
//
// It states no Status, State or Role for the same reason it states no port: it
// knows none of them. Declaring a standby Active and running because the primary
// still lists it would mask one whose postgres is actually down, and the
// standby's own connection reports all three anyway. A response claiming no
// state of its own is how the discovery layer tells an attribute from a node
// (see cluster.mergeKeeperSync).
//
// It names its node rather than addressing it. application_name is self-declared
// by the standby's own primary_conninfo, exactly like a patroni member name, and
// the shipped template sets it to the node's name; a name is unique within a
// cluster, so it identifies one node where a host cannot - every node of a
// single-host cluster shares one. The alternative, client_addr, is worse than
// ambiguous: it is the network-observed source address of the replication
// connection, never a domain name, and an external hop or a shared docker
// network can rewrite it into something that corresponds to no configured node
// at all. A standby whose application_name matches no configured node is
// skipped, same as one that sets none.
func mapSyncStandby(name string, syncState string) keeper.Response {
	return keeper.Response{
		Sync:           syncState == "sync" || syncState == "quorum",
		DiscoveredName: &name,
	}
}

func mapNode(host string, port int, inRecovery bool, lag int64) keeper.Response {
	role := keeper.Leader
	if inRecovery {
		role = keeper.Replica
	} else {
		lag = 0
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
