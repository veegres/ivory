package postgres

import (
	"context"
	"fmt"
	"ivory/clients/postgres"
	"ivory/plugins/database"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

var _ database.DbRoleAdmin = (*Adapter)(nil)
var _ database.DbDatabaseAdmin = (*Adapter)(nil)

const listDbRolesQuery = `SELECT rolname, rolsuper, rolcreatedb, rolcanlogin
FROM pg_catalog.pg_roles
ORDER BY rolname`

const listDbDatabasesQuery = `SELECT d.datname,
       pg_catalog.pg_get_userbyid(d.datdba) AS owner,
       pg_catalog.pg_size_pretty(pg_catalog.pg_database_size(d.datname)) AS size,
       pg_catalog.pg_encoding_to_char(d.encoding) AS encoding
FROM pg_catalog.pg_database d
WHERE NOT d.datistemplate
ORDER BY d.datname`

func (a *Adapter) ListDbRoles(ctx database.Context) ([]database.DbRole, error) {
	roles := make([]database.DbRole, 0)
	errReq := a.sendRequest(ctx, listDbRolesQuery, nil, func(rows pgx.Rows, _ *pgtype.Map, _ string) error {
		for rows.Next() {
			var role database.DbRole
			if err := rows.Scan(&role.Name, &role.Superuser, &role.CreateDb, &role.Login); err != nil {
				return err
			}
			roles = append(roles, role)
		}
		return nil
	})
	if errReq != nil {
		return nil, errReq
	}
	return roles, nil
}

func (a *Adapter) CreateDbRole(ctx database.Context, name string, password string) error {
	if err := validateRoleName(name); err != nil {
		return err
	}
	sql := fmt.Sprintf("CREATE ROLE %s WITH LOGIN PASSWORD %s",
		pgx.Identifier{name}.Sanitize(), quotePostgresLiteral(password))
	return a.execStatement(ctx, sql)
}

func (a *Adapter) UpdateDbRolePassword(ctx database.Context, name string, password string) error {
	if err := validateRoleName(name); err != nil {
		return err
	}
	sql := fmt.Sprintf("ALTER ROLE %s WITH PASSWORD %s",
		pgx.Identifier{name}.Sanitize(), quotePostgresLiteral(password))
	return a.execStatement(ctx, sql)
}

func (a *Adapter) DeleteDbRole(ctx database.Context, name string) error {
	if err := validateRoleDelete(ctx, name); err != nil {
		return err
	}
	sql := fmt.Sprintf("DROP ROLE %s", pgx.Identifier{name}.Sanitize())
	return a.execStatement(ctx, sql)
}

func (a *Adapter) ListDbDatabases(ctx database.Context) ([]database.DbDatabaseCatalog, error) {
	dbs := make([]database.DbDatabaseCatalog, 0)
	errReq := a.sendRequest(ctx, listDbDatabasesQuery, nil, func(rows pgx.Rows, _ *pgtype.Map, _ string) error {
		for rows.Next() {
			var db database.DbDatabaseCatalog
			if err := rows.Scan(&db.Name, &db.Owner, &db.Size, &db.Encoding); err != nil {
				return err
			}
			dbs = append(dbs, db)
		}
		return nil
	})
	if errReq != nil {
		return nil, errReq
	}
	return dbs, nil
}

func (a *Adapter) CreateDbDatabase(ctx database.Context, name string) error {
	if err := validateDatabaseName(name); err != nil {
		return err
	}
	sql := fmt.Sprintf("CREATE DATABASE %s", pgx.Identifier{name}.Sanitize())
	return a.execStatementAutocommit(ctx, sql)
}

func (a *Adapter) RenameDbDatabase(ctx database.Context, name string, newName string) error {
	if err := validateDatabaseDelete(name); err != nil {
		return err
	}
	if err := validateDatabaseName(newName); err != nil {
		return err
	}
	sql := fmt.Sprintf("ALTER DATABASE %s RENAME TO %s",
		pgx.Identifier{name}.Sanitize(), pgx.Identifier{newName}.Sanitize())
	return a.execStatement(ctx, sql)
}

func (a *Adapter) DeleteDbDatabase(ctx database.Context, name string) error {
	if err := validateDatabaseDelete(name); err != nil {
		return err
	}
	sql := fmt.Sprintf("DROP DATABASE %s", pgx.Identifier{name}.Sanitize())
	return a.execStatementAutocommit(ctx, sql)
}

func quotePostgresLiteral(value string) string {
	return "'" + strings.ReplaceAll(value, "'", "''") + "'"
}

func (a *Adapter) execStatementAutocommit(ctx database.Context, query string, args ...any) error {
	conn, _, errConn := a.getConnection(ctx)
	if errConn != nil {
		return errConn
	}
	defer postgres.Close(conn)

	execCtx := context.Background()
	_, err := conn.Exec(execCtx, query, args...)
	return err
}

func (a *Adapter) execStatement(ctx database.Context, query string, args ...any) error {
	conn, _, errConn := a.getConnection(ctx)
	if errConn != nil {
		return errConn
	}
	defer postgres.Close(conn)

	txCtx := context.Background()
	tx, errTx := conn.Begin(txCtx)
	if errTx != nil {
		return errTx
	}
	var err error
	defer func() { a.closeTransaction(tx, txCtx, err) }()

	if ctx.Connection.Config.Schema != nil {
		safeSchema := pgx.Identifier{*ctx.Connection.Config.Schema}.Sanitize()
		_, err = tx.Exec(txCtx, fmt.Sprintf("SET LOCAL search_path TO %s", safeSchema))
		if err != nil {
			return fmt.Errorf("failed to set schema: %w", err)
		}
	}

	_, err = tx.Exec(txCtx, query, args...)
	return err
}

func validateRoleName(name string) error {
	if isProtectedRoleName(name) {
		return database.ErrProtectedRole
	}
	return nil
}

func validateRoleDelete(ctx database.Context, name string) error {
	if err := validateRoleName(name); err != nil {
		return err
	}
	sessionUser := sessionRoleName(ctx)
	if sessionUser != "" && strings.EqualFold(name, sessionUser) {
		return database.ErrCannotDropSessionRole
	}
	return nil
}

func validateDatabaseName(name string) error {
	if isProtectedDatabaseName(name) {
		return database.ErrProtectedDatabase
	}
	return nil
}

func validateDatabaseDelete(name string) error {
	return validateDatabaseName(name)
}

func isProtectedRoleName(name string) bool {
	if name == "postgres" {
		return true
	}
	return strings.HasPrefix(name, "pg_")
}

func isProtectedDatabaseName(name string) bool {
	switch name {
	case "postgres", "template0", "template1":
		return true
	default:
		return false
	}
}

func sessionRoleName(ctx database.Context) string {
	if ctx.Connection == nil || ctx.Connection.Credentials == nil {
		return ""
	}
	return ctx.Connection.Credentials.Username
}
