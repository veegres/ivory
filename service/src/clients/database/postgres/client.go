package postgres

import (
	"context"
	"errors"
	"fmt"
	"ivory/src/clients/database"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"golang.org/x/exp/slog"
)

var ErrCannotLimitWithoutTrim = errors.New("cannot limit query without trimming it")
var ErrDatabaseHostOrPortNotSpecified = errors.New("database host or port are not specified")
var ErrPasswordNotSet = errors.New("password is not set")

// NOTE: validate that is matches interface in compile-time
var _ database.Client = (*Client)(nil)

type Client struct {
	appName string
}

func NewClient(appName string) *Client {
	return &Client{appName: appName}
}

func (s *Client) GetApplicationName(session string) string {
	return s.appName + " [" + fmt.Sprintf("%.7s", session) + "]"
}

func (s *Client) GetMany(ctx database.Context, query string, queryParams []any) ([]string, error) {
	values := make([]string, 0)
	errReq := s.sendRequest(ctx, query, queryParams, func(rows pgx.Rows, _ *pgtype.Map, _ string) error {
		for rows.Next() {
			var value string
			err := rows.Scan(&value)
			if err != nil {
				return err
			}
			values = append(values, value)
		}
		return nil
	})
	if errReq != nil {
		return nil, errReq
	}
	return values, nil
}

func (s *Client) GetOne(ctx database.Context, query string) (any, error) {
	var value any
	errReq := s.sendRequest(ctx, query, nil, func(rows pgx.Rows, _ *pgtype.Map, _ string) error {
		if rows.Next() {
			values, err := rows.Values()
			if err != nil {
				return err
			}
			value = values[0]
		}
		return nil
	})
	if errReq != nil {
		return nil, errReq
	}
	return value, nil
}

func (s *Client) GetFields(ctx database.Context, query string, options *database.QueryOptions) (*database.QueryFields, error) {
	startTime := time.Now().UnixMilli()

	fields := make([]database.QueryField, 0)
	rowList := make([][]any, 0)
	url := "-"

	// NOTE: we need this object ot avoid fatal errors when the option variable is `nil`
	tmpOptions := database.QueryOptions{}
	if options != nil {
		tmpOptions.Limit = options.Limit
		tmpOptions.Trim = options.Trim
		tmpOptions.Params = options.Params
	}

	normQuery, normLimit, errNormQuery := s.normalizeQuery(query, tmpOptions.Trim, tmpOptions.Limit)
	if errNormQuery != nil {
		return nil, errNormQuery
	}

	errReq := s.sendRequest(ctx, normQuery, tmpOptions.Params, func(rows pgx.Rows, typeMap *pgtype.Map, connUrl string) error {
		url = connUrl
		for _, field := range rows.FieldDescriptions() {
			dataType, ok := typeMap.TypeForOID(field.DataTypeOID)
			if !ok {
				fields = append(fields, database.QueryField{Name: field.Name, DataType: "unknown", DataTypeOID: field.DataTypeOID})
			} else {
				fields = append(fields, database.QueryField{Name: field.Name, DataType: dataType.Name, DataTypeOID: field.DataTypeOID})
			}
		}
		for rows.Next() {
			val, err := rows.Values()
			if err != nil {
				return err
			}
			rowList = append(rowList, val)
		}
		return nil
	})
	if errReq != nil {
		return nil, errReq
	}

	if options != nil {
		options.Limit = normLimit
	}
	endTime := time.Now().UnixMilli()
	res := &database.QueryFields{
		Fields:    fields,
		Rows:      rowList,
		Url:       url,
		StartTime: startTime,
		EndTime:   endTime,
		Options:   options,
	}
	return res, nil
}

func (s *Client) Cancel(ctx database.Context, pid int) error {
	return s.sendRequest(ctx, "SELECT pg_cancel_backend("+strconv.Itoa(pid)+")", nil, nil)
}

func (s *Client) Terminate(ctx database.Context, pid int) error {
	return s.sendRequest(ctx, "SELECT pg_terminate_backend("+strconv.Itoa(pid)+")", nil, nil)
}

type fn func(pgx.Rows, *pgtype.Map, string) error

func (s *Client) sendRequest(ctx database.Context, query string, queryParams []any, parse fn) error {
	conn, connUrl, errConn := s.getConnection(ctx)
	if errConn != nil {
		return errConn
	}
	defer s.closeConnection(conn, context.Background())

	var rows pgx.Rows
	var err error
	if queryParams == nil {
		rows, err = conn.Query(context.Background(), query)
	} else {
		rows, err = conn.Query(context.Background(), query, queryParams...)
	}
	if err != nil {
		return err
	}

	defer rows.Close()
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

func (s *Client) normalizeQuery(query string, trim *bool, limit *string) (string, *string, error) {
	if trim == nil || *trim == false {
		if limit != nil {
			return "", limit, ErrCannotLimitWithoutTrim
		} else {
			return query, limit, nil
		}
	}
	trimmedQuery := s.trimQuery(query)
	if limit == nil {
		return trimmedQuery, limit, nil
	}
	parsed := s.parseQuery(trimmedQuery)
	newQuery, newLimit := s.addLimitToQuery(trimmedQuery, parsed, *limit)
	return newQuery, newLimit, nil
}

func (s *Client) addLimitToQuery(query string, queryAnalysis database.QueryAnalysis, limit string) (string, *string) {
	if queryAnalysis.LIMIT == 0 && queryAnalysis.SELECT > 0 && queryAnalysis.FROM > 0 && queryAnalysis.EXPLAIN == 0 &&
		queryAnalysis.DELETE == 0 && queryAnalysis.UPDATE == 0 && queryAnalysis.INSERT == 0 {
		replace := " LIMIT " + limit + ";"
		if queryAnalysis.Semicolon {
			// NOTE: removing all spaces and semicolon at the end
			regex := regexp.MustCompile("\\s*;\\s*$")
			return regex.ReplaceAllString(query, replace), &limit
		} else {
			return query + replace, &limit
		}
	}
	return query, nil
}

func (s *Client) trimQuery(query string) string {
	// NOTE: removing all comments from the query with spaces and enters or end of string
	regex := regexp.MustCompile("\\s*--.*\\n*")
	return regex.ReplaceAllString(query, "")
}

func (s *Client) parseQuery(query string) database.QueryAnalysis {
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

func (s *Client) getConnection(ctx database.Context) (*pgx.Conn, string, error) {
	connection := ctx.Connection
	db := connection.Database
	if db.Port == 0 || db.Host == "" || db.Host == "-" {
		return nil, "unknown", ErrDatabaseHostOrPortNotSpecified
	}

	dbName := "postgres"
	if db.Name != nil && *db.Name != "" {
		dbName = *db.Name
	}

	credentials := connection.Credentials
	if credentials == nil {
		return nil, "unknown", ErrPasswordNotSet
	}

	connProtocol := "postgres://"
	connHost := db.Host + ":" + strconv.Itoa(db.Port) + "/" + dbName
	connUrl := connProtocol + connHost

	tlsConfig := connection.TlsConfig
	if tlsConfig != nil {
		// NOTE: verify-ca was chosen because it potentially can protect from machine-in-the-middle attack if
		// it has the right CA policy. More info can be found here https://www.postgresql.org/docs/16/libpq-ssl.html#LIBPQ-SSL-PROTECTION
		connUrl += "?sslmode=verify-ca"
	}

	conConfig, errConfig := pgx.ParseConfig(connUrl)
	if errConfig != nil {
		return nil, connUrl, errConfig
	}
	conConfig.User = credentials.Username
	conConfig.Password = credentials.Password
	conConfig.RuntimeParams = map[string]string{
		"application_name": s.GetApplicationName(ctx.Session),
	}
	if tlsConfig != nil {
		// NOTE: we rewrite only RootCAs and Certificates, because pgx.ParseConfig creates proper
		//  tlsConfig for different `sslmode`. For example `verify-ca` should mark `InsecureSkipVerify=true`
		//  and it always sets `ServerName` it required for `verify-full` mode.
		conConfig.TLSConfig.RootCAs = tlsConfig.RootCAs
		conConfig.TLSConfig.Certificates = tlsConfig.Certificates
	}

	conn, err := pgx.ConnectConfig(context.Background(), conConfig)
	return conn, connUrl, err
}

func (s *Client) closeConnection(conn *pgx.Conn, ctx context.Context) {
	err := conn.Close(ctx)
	if err != nil {
		slog.Warn("postgres close connection", err)
	}
}
