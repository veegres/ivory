package clickhouse

import (
	"ivory/clients/clickhouse"
	"ivory/plugins/database"
)

// Cancel/Terminate/ActiveQueries are not supported: unlike postgres' integer backend pid,
// clickhouse identifies a running query by its string query_id
// (system.processes.query_id, matched by KILL QUERY WHERE query_id = ...) -
// there is no integer identifier to satisfy this interface's pid parameter,
// so Cancel/Terminate/ActiveQueries would have nothing correct to do with it.

func (a *Adapter) Cancel(ctx database.Context, pid int) error {
	return ErrNotSupported
}

func (a *Adapter) Terminate(ctx database.Context, pid int) error {
	return ErrNotSupported
}

func (a *Adapter) ActiveQueries(ctx database.Context, options *database.QueryOptions) (*database.QueryFields, error) {
	return nil, ErrNotSupported
}

func (a *Adapter) Ping(ctx database.Context) error {
	conn, _, err := a.connect(ctx)
	if err != nil {
		return err
	}
	clickhouse.Close(conn)
	return nil
}
