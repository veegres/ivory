package clickhouse

import (
	"context"
	"errors"
	chclient "ivory/clients/clickhouse"
	"ivory/plugins/database"
	"regexp"
	"strings"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
)

var ErrNotSupported = errors.New("this operation is not supported for clickhouse")

// NOTE: validate that is matches interface in compile-time
var _ database.Adapter = (*Adapter)(nil)

// requestTimeout bounds every non-console-query request (schema/session
// calls), so an unreachable node cannot hang the caller. GetFields itself is
// user-driven and uses the connection's own default timeout instead (see
// clickhouse_adapter_query.go).
const requestTimeout = 10 * time.Second

type Adapter struct{}

func NewAdapter() *Adapter {
	return &Adapter{}
}

func (a *Adapter) connect(ctx database.Context) (driver.Conn, string, error) {
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

	return chclient.Connect(context.Background(), chclient.Config{
		Host:     db.Host,
		Port:     db.Port,
		Database: dbName,
		Username: credentials.Username,
		Password: credentials.Password,
		TLS:      connection.TlsConfig,
	})
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
			regex := regexp.MustCompile(`\s*;\s*$`)
			return regex.ReplaceAllString(query, replace), &limit
		}
		return query + replace, &limit
	}
	return query, nil
}

func (a *Adapter) trimQuery(query string) string {
	commentRegex := regexp.MustCompile("--.*")
	query = commentRegex.ReplaceAllString(query, " ")
	return strings.Join(strings.Fields(query), " ")
}

func (a *Adapter) parseQuery(query string) database.QueryAnalysis {
	lowerQuery := strings.ToLower(query)
	words := strings.Fields(lowerQuery)
	parsed := database.QueryAnalysis{}
	for i, word := range words {
		// NOTE: we need this check to avoid params rename confusion
		if i-1 > 0 && words[i-1] == "as" {
			continue
		}

		switch word {
		case "limit":
			parsed.LIMIT += 1
		case "update", "alter":
			parsed.UPDATE += 1
		case "insert":
			parsed.INSERT += 1
		case "delete", "truncate", "drop":
			parsed.DELETE += 1
		case "select":
			parsed.SELECT += 1
		case "from":
			parsed.FROM += 1
		case "explain":
			parsed.EXPLAIN += 1
		}
	}
	if len(words) == 0 {
		return parsed
	}
	lastWord := words[len(words)-1]
	if lastWord[len(lastWord)-1:] == ";" {
		parsed.Semicolon = true
	}
	return parsed
}
