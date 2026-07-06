package postgres

import "ivory/plugins/database"

func (a *Adapter) ListDatabases(ctx database.Context, name string) ([]string, error) {
	return a.GetMany(ctx, GetAllDatabases, []any{"%" + name + "%"})
}

func (a *Adapter) ListSchemas(ctx database.Context, name string) ([]string, error) {
	return a.GetMany(ctx, GetAllSchemas, []any{"%" + name + "%"})
}

func (a *Adapter) ListTables(ctx database.Context, schema string, name string) ([]string, error) {
	return a.GetMany(ctx, GetAllTables, []any{schema, "%" + name + "%"})
}
