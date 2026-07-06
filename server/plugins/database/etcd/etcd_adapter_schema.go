package etcd

import "ivory/plugins/database"

func (a *Adapter) ListDatabases(ctx database.Context, name string) ([]string, error) {
	return nil, ErrNotSupported
}

func (a *Adapter) ListSchemas(ctx database.Context, name string) ([]string, error) {
	return nil, ErrNotSupported
}

func (a *Adapter) ListTables(ctx database.Context, schema string, name string) ([]string, error) {
	return nil, ErrNotSupported
}
