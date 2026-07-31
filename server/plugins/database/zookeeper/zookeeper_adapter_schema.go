package zookeeper

import "ivory/plugins/database"

// ListDatabases/ListSchemas/ListTables have nothing to report: zookeeper has
// no database/schema/table hierarchy, only the znode tree already reachable
// through the console's own `ls` command (see zookeeper_adapter_query.go).

func (a *Adapter) ListDatabases(ctx database.Context, name string) ([]string, error) {
	return nil, ErrNotSupported
}

func (a *Adapter) ListSchemas(ctx database.Context, name string) ([]string, error) {
	return nil, ErrNotSupported
}

func (a *Adapter) ListTables(ctx database.Context, schema string, name string) ([]string, error) {
	return nil, ErrNotSupported
}
