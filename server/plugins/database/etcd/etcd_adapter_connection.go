package etcd

import (
	"context"
	"ivory/plugins/database"
	"strconv"
)

func (a *Adapter) Cancel(ctx database.Context, pid int) error {
	return ErrNotSupported
}

func (a *Adapter) Terminate(ctx database.Context, pid int) error {
	return ErrNotSupported
}

func (a *Adapter) ActiveQueries(ctx database.Context, options *database.QueryOptions) (*database.QueryFields, error) {
	return nil, ErrNotSupported
}

// Ping asks the endpoint for its status, since creating an etcd client dials
// lazily and proves nothing on its own.
func (a *Adapter) Ping(ctx database.Context) error {
	client, _, err := a.connect(ctx)
	if err != nil {
		return err
	}
	defer client.Close()
	requestCtx, cancel := context.WithTimeout(context.Background(), client.Timeout)
	defer cancel()
	db := ctx.Connection.Config
	_, err = client.Status(requestCtx, db.Host+":"+strconv.Itoa(db.Port))
	return err
}
