package clickhouse

import "ivory/plugins/database"

const GetAllDatabases = `SELECT name FROM system.databases WHERE name ILIKE ? ORDER BY name`
const GetAllTables = `SELECT name FROM system.tables WHERE database = ? AND name ILIKE ? ORDER BY name`

func (a *Adapter) ListDatabases(ctx database.Context, name string) ([]string, error) {
	return a.GetMany(ctx, GetAllDatabases, []any{"%" + name + "%"})
}

// ListSchemas has nothing to report: clickhouse has no schema layer between
// database and table - a database already plays that role, which is why
// ListTables takes a database name in its schema argument.
func (a *Adapter) ListSchemas(ctx database.Context, name string) ([]string, error) {
	return nil, ErrNotSupported
}

func (a *Adapter) ListTables(ctx database.Context, schema string, name string) ([]string, error) {
	return a.GetMany(ctx, GetAllTables, []any{schema, "%" + name + "%"})
}
