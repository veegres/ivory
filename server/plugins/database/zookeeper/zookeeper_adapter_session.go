package zookeeper

import "ivory/plugins/database"

// Cancel/Terminate/ActiveQueries have nothing to report: znode operations
// have no per-query session or pid concept the way a SQL backend does.

func (a *Adapter) Cancel(ctx database.Context, pid int) error {
	return ErrNotSupported
}

func (a *Adapter) Terminate(ctx database.Context, pid int) error {
	return ErrNotSupported
}

func (a *Adapter) ActiveQueries(ctx database.Context, options *database.QueryOptions) (*database.QueryFields, error) {
	return nil, ErrNotSupported
}
