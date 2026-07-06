package postgres

import (
	"context"
	"fmt"
	"ivory/clients/postgres"
	"ivory/plugins/database"
	"regexp"
	"strings"

	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

// NOTE: validate that is matches interface in compile-time
var _ database.Adapter = (*Adapter)(nil)

type Adapter struct{}

func NewAdapter() *Adapter {
	return &Adapter{}
}

type fn func(pgx.Rows, *pgtype.Map, string) error

func (a *Adapter) sendRequest(ctx database.Context, query string, queryParams []any, parse fn) error {
	conn, connUrl, errConn := a.getConnection(ctx)
	if errConn != nil {
		return errConn
	}
	defer postgres.Close(conn)

	txCtx := context.Background()
	tx, errTx := conn.Begin(txCtx)
	if errTx != nil {
		return errTx
	}
	defer a.closeTransaction(tx, txCtx)

	if ctx.Connection.Config.Schema != nil {
		safeSchema := pgx.Identifier{*ctx.Connection.Config.Schema}.Sanitize()
		_, errSchema := tx.Exec(txCtx, fmt.Sprintf("SET LOCAL search_path TO %s", safeSchema))
		if errSchema != nil {
			return fmt.Errorf("failed to set schema: %w", errSchema)
		}
	}

	var rows pgx.Rows
	var err error
	if queryParams == nil {
		rows, err = tx.Query(txCtx, query)
	} else {
		rows, err = tx.Query(txCtx, query, queryParams...)
	}
	if err != nil {
		return err
	}

	defer func() {
		rows.Close()
	}()
	if parse != nil {
		errParse := parse(rows, conn.TypeMap(), connUrl)
		if errParse != nil {
			return errParse
		}
	}
	if rows.Err() != nil {
		return rows.Err()
	}

	return nil
}

func (a *Adapter) normalizeQuery(query string, trim *bool, limit *string) (string, *string, error) {
	if trim == nil || *trim == false {
		if limit != nil {
			return "", limit, database.ErrCannotLimitWithoutTrim
		}
		return query, limit, nil
	}
	trimmedQuery := a.trimQuery(query)
	if limit == nil {
		return trimmedQuery, limit, nil
	}
	parsed := a.parseQuery(trimmedQuery)
	newQuery, newLimit := a.addLimitToQuery(trimmedQuery, parsed, *limit)
	return newQuery, newLimit, nil
}

func (a *Adapter) addLimitToQuery(query string, queryAnalysis database.QueryAnalysis, limit string) (string, *string) {
	if queryAnalysis.LIMIT == 0 && queryAnalysis.SELECT > 0 && queryAnalysis.FROM > 0 && queryAnalysis.EXPLAIN == 0 &&
		queryAnalysis.DELETE == 0 && queryAnalysis.UPDATE == 0 && queryAnalysis.INSERT == 0 {
		replace := " LIMIT " + limit + ";"
		if queryAnalysis.Semicolon {
			// NOTE: removing all spaces and semicolon at the end
			regex := regexp.MustCompile("\\s*;\\s*$")
			return regex.ReplaceAllString(query, replace), &limit
		}
		return query + replace, &limit
	}
	return query, nil
}

func (a *Adapter) trimQuery(query string) string {
	// Remove comments
	commentRegex := regexp.MustCompile("--.*")
	query = commentRegex.ReplaceAllString(query, " ")

	// Normalize whitespace (including tabs and newlines)
	return strings.Join(strings.Fields(query), " ")
}

func (a *Adapter) parseQuery(query string) database.QueryAnalysis {
	lowerQuery := strings.ToLower(query)
	words := strings.Fields(lowerQuery)
	parsed := database.QueryAnalysis{LIMIT: 0, UPDATE: 0, SELECT: 0, INSERT: 0, DELETE: 0, Semicolon: false}
	for i, word := range words {
		// NOTE: we need this check to avoid params rename confusion
		if i-1 > 0 && words[i-1] == "as" {
			continue
		}

		switch word {
		case "limit":
			parsed.LIMIT += 1
		case "update":
			parsed.UPDATE += 1
		case "insert":
			parsed.INSERT += 1
		case "delete":
			parsed.DELETE += 1
		case "select":
			parsed.SELECT += 1
		case "from":
			parsed.FROM += 1
		case "explain":
			parsed.EXPLAIN += 1
		}
	}
	lastWord := words[len(words)-1]
	if lastWord[len(lastWord)-1:] == ";" {
		parsed.Semicolon = true
	}
	return parsed
}

func (a *Adapter) getConnection(ctx database.Context) (*pgx.Conn, string, error) {
	connection := ctx.Connection
	db := connection.Config
	if db.Port == 0 || db.Host == "" || db.Host == "-" {
		return nil, "unknown", database.ErrDatabaseHostOrPortNotSpecified
	}

	dbName := ""
	if db.Name != nil {
		dbName = *db.Name
	}

	credentials := connection.Credentials
	if credentials == nil {
		return nil, "unknown", database.ErrPasswordNotSet
	}

	return postgres.Connect(context.Background(), postgres.Config{
		Host:     db.Host,
		Port:     db.Port,
		Database: dbName,
		Username: credentials.Username,
		Password: credentials.Password,
		AppName:  ctx.Application,
		TLS:      connection.TlsConfig,
	})
}

func (a *Adapter) closeTransaction(tx pgx.Tx, txCtx context.Context) {
	err := tx.Rollback(txCtx)
	if err != nil {
		slog.Warn("postgres rollback", "error", err)
	}
}
