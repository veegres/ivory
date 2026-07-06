package postgres

import (
	"ivory/plugins/database"
	"strconv"
)

func (a *Adapter) Cancel(ctx database.Context, pid int) error {
	return a.sendRequest(ctx, "SELECT pg_cancel_backend("+strconv.Itoa(pid)+")", nil, nil)
}

func (a *Adapter) Terminate(ctx database.Context, pid int) error {
	return a.sendRequest(ctx, "SELECT pg_terminate_backend("+strconv.Itoa(pid)+")", nil, nil)
}

func (a *Adapter) ActiveQueries(ctx database.Context, options *database.QueryOptions) (*database.QueryFields, error) {
	return a.GetFields(ctx, GetAllActiveQueriesByApplicationName, options)
}
