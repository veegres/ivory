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

var ErrNodeIsPrimary = errors.New("failover target must be a replica; this node is the primary")
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
var _ keeper.Adapter = (*Adapter)(nil)

// Adapter provides the native postgres experience alongside patroni by
// talking to postgres directly. The keeper connection host/port is the
// postgres host/port (keeperPort == dbPort convention) and the keeper vault
// holds database credentials. Operations that require orchestration or OS
// access are not supported and excluded from SupportedFeatures - in
// particular, Ivory never configures streaming replication for native
// postgres (DeploymentSpec only deploys a bare postgres image), so
// identifying a connected standby's Sync status (see mapSyncStandby) depends
// entirely on the operator manually setting each standby's primary_conninfo
// application_name to match the Host Ivory has configured for it.
type Adapter struct{}

func NewAdapter() *Adapter {
	return &Adapter{}
}

func (a *Adapter) List(request keeper.Request) ([]keeper.Response, int, error) {
	var inRecovery bool
	var lag int64
	err := a.queryRow(request, listQuery, func(row pgx.Row) error {
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
		responses = append(responses, a.listSyncStandbys(request)...)
	}
	return responses, http.StatusOK, nil
}

// listSyncStandbys reports every standby currently connected to this
// primary, so replicas queried independently (see List) can be told apart as
// synchronous vs asynchronous - something they cannot determine about
// themselves. It is best-effort: pg_stat_replication may be restricted for
// the configured credentials, or briefly fail during a topology change, in
// which case the primary's own List response above is still returned as-is.
func (a *Adapter) listSyncStandbys(request keeper.Request) []keeper.Response {
	var standbys []keeper.Response
	err := a.query(request, syncStandbyQuery, func(rows pgx.Rows) error {
		for rows.Next() {
			var applicationName, syncState string
			if errScan := rows.Scan(&applicationName, &syncState); errScan != nil {
				return errScan
			}
			// NOTE: pg_stat_replication only has the standby's ephemeral
			// replication client_port, not its listening port, so this
			// reuses request.Port - the port Ivory already dialed *this*
			// primary on - as the standby's port too, since every node in a
			// native postgres cluster shares one port by convention (see
			// Adapter's doc comment). This is data the adapter itself
			// legitimately possesses, not a value borrowed from Ivory's own
			// cluster configuration.
			standbys = append(standbys, mapSyncStandby(applicationName, request.Port, syncState))
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

func (a *Adapter) Config(request keeper.Request) (any, int, error) {
	settings := map[string]string{}
	err := a.query(request, configQuery, func(rows pgx.Rows) error {
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

func (a *Adapter) Reload(request keeper.Request) (*string, int, error) {
	var reloaded bool
	err := a.queryRow(request, "SELECT pg_reload_conf()", func(row pgx.Row) error {
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

func (a *Adapter) Failover(request keeper.Request) (*string, int, error) {
	var inRecovery bool
	err := a.queryRow(request, "SELECT pg_is_in_recovery()", func(row pgx.Row) error {
		return row.Scan(&inRecovery)
	})
	if err != nil {
		return nil, http.StatusBadRequest, err
	}
	if !inRecovery {
		return nil, http.StatusBadRequest, ErrNodeIsPrimary
	}
	// NOTE: pg_promote(false) only signals the promotion and returns
	// immediately, so it fits into the request timeout
	var promoted bool
	errPromote := a.queryRow(request, "SELECT pg_promote(false)", func(row pgx.Row) error {
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

func (a *Adapter) Activate(keeper.Request) (*string, int, error) {
	return nil, http.StatusNotImplemented, keeper.ErrNotSupported
}

func (a *Adapter) Pause(keeper.Request) (*string, int, error) {
	return nil, http.StatusNotImplemented, keeper.ErrNotSupported
}

func (a *Adapter) queryRow(request keeper.Request, query string, scan func(row pgx.Row) error) error {
	return a.withConnection(request, func(ctx context.Context, conn *pgx.Conn) error {
		return scan(conn.QueryRow(ctx, query))
	})
}

func (a *Adapter) query(request keeper.Request, query string, parse func(rows pgx.Rows) error) error {
	return a.withConnection(request, func(ctx context.Context, conn *pgx.Conn) error {
		rows, err := conn.Query(ctx, query)
		if err != nil {
			return err
		}
		defer rows.Close()
		return parse(rows)
	})
}

func (a *Adapter) withConnection(request keeper.Request, action func(ctx context.Context, conn *pgx.Conn) error) error {
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

// mapSyncStandby builds a Response for a standby discovered via the
// primary's pg_stat_replication.
//
// host comes from application_name, not client_addr. client_addr is the
// network-observed source address of the replication TCP connection - it's
// never a domain name, and in Ivory's own Docker-over-SSH deployments it can
// be a container-internal address that doesn't correspond to the node's
// configured Host at all (an external network hop can rewrite it, or a
// single-host cluster's containers can share one docker network). That
// makes it fundamentally unlike Patroni's or etcd's discovered host, which
// is self-declared by the same operator who configured Ivory. application_name
// is also self-declared (by the standby's own primary_conninfo), so it's the
// right kind of value in principle - but Ivory has no orchestration for
// native postgres (see Adapter's doc comment) and never sets it itself, even
// for clusters it deployed, so this only resolves to the right node if the
// operator manually sets each standby's primary_conninfo application_name to
// exactly the Host configured for that node in Ivory. Standbys that don't
// follow this convention are simply skipped by syncStandbyQuery, same as if
// they weren't connected at all.
//
// port is the primary's own connection port, reused for the standby too
// since native postgres has no separate keeper port - every node's keeper
// port and db port are the same port it's dialed on (see Adapter's doc
// comment), and by convention every node in the cluster shares that same
// port number.
func mapSyncStandby(host string, port int, syncState string) keeper.Response {
	var status keeper.Status = keeper.Active
	return keeper.Response{
		Status:               &status,
		State:                keeper.StateRunning,
		Role:                 keeper.Replica,
		Sync:                 syncState == "sync" || syncState == "quorum",
		DiscoveredHost:       &host,
		DiscoveredKeeperPort: &port,
		DiscoveredDbPort:     &port,
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
