package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"golang.org/x/exp/slog"
	"ivory/src/config"
	. "ivory/src/model"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type PostgresClient struct {
	appName         string
	passwordService *PasswordService
	certService     *CertService
}

func NewPostgresClient(
	env *config.Env,
	passwordService *PasswordService,
	certService *CertService,
) *PostgresClient {
	return &PostgresClient{
		appName:         env.Version.Label,
		passwordService: passwordService,
		certService:     certService,
	}
}

func (s *PostgresClient) GetMany(ctx QueryContext, query string, queryParams []any) ([]string, error) {
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

func (s *PostgresClient) GetOne(ctx QueryContext, query string) (any, error) {
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

func (s *PostgresClient) GetFields(ctx QueryContext, query string, options *QueryOptions) (*QueryFields, error) {
	startTime := time.Now().UnixMilli()

	fields := make([]QueryField, 0)
	rowList := make([][]any, 0)
	url := "-"

	// NOTE: we need this object ot avoid fatal errors when options is `nil`
	tmpOptions := QueryOptions{}
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
				fields = append(fields, QueryField{Name: field.Name, DataType: "unknown", DataTypeOID: field.DataTypeOID})
			} else {
				fields = append(fields, QueryField{Name: field.Name, DataType: dataType.Name, DataTypeOID: field.DataTypeOID})
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

	options.Limit = normLimit
	endTime := time.Now().UnixMilli()
	res := &QueryFields{
		Fields:    fields,
		Rows:      rowList,
		Url:       url,
		StartTime: startTime,
		EndTime:   endTime,
		Options:   options,
	}
	return res, nil
}

func (s *PostgresClient) Cancel(ctx QueryContext, pid int) error {
	return s.sendRequest(ctx, "SELECT pg_cancel_backend("+strconv.Itoa(pid)+")", nil, nil)
}

func (s *PostgresClient) Terminate(ctx QueryContext, pid int) error {
	return s.sendRequest(ctx, "SELECT pg_terminate_backend("+strconv.Itoa(pid)+")", nil, nil)
}

type fn func(pgx.Rows, *pgtype.Map, string) error

func (s *PostgresClient) sendRequest(ctx QueryContext, query string, queryParams []any, parse fn) error {
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

func (s *PostgresClient) normalizeQuery(query string, trim *bool, limit *string) (string, *string, error) {
	if trim == nil || *trim == false {
		if limit != nil {
			return "", limit, errors.New("cannot limit query without trimming it")
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

func (s *PostgresClient) addLimitToQuery(query string, queryParsed QueryParsed, limit string) (string, *string) {
	if queryParsed.LIMIT == 0 && queryParsed.SELECT > 0 && queryParsed.FROM > 0 && queryParsed.EXPLAIN == 0 &&
		queryParsed.DELETE == 0 && queryParsed.UPDATE == 0 && queryParsed.INSERT == 0 {
		replace := " LIMIT " + limit + ";"
		if queryParsed.Semicolon {
			// NOTE: removing all spaces and semicolon at the end
			regex := regexp.MustCompile("\\s*;\\s*$")
			return regex.ReplaceAllString(query, replace), &limit
		} else {
			return query + replace, &limit
		}
	}
	return query, nil
}

func (s *PostgresClient) trimQuery(query string) string {
	// NOTE: removing all comments from the query with spaces and enters or end of string
	regex := regexp.MustCompile("\\s*--.*\\n*")
	return regex.ReplaceAllString(query, "")
}

func (s *PostgresClient) parseQuery(query string) QueryParsed {
	lowerQuery := strings.ToLower(query)
	words := strings.Fields(lowerQuery)
	parsed := QueryParsed{LIMIT: 0, UPDATE: 0, SELECT: 0, INSERT: 0, DELETE: 0, Semicolon: false}
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

func (s *PostgresClient) getConnection(ctx QueryContext) (*pgx.Conn, string, error) {
	db := ctx.Connection.Db
	if db.Port == 0 || db.Host == "" || db.Host == "-" {
		return nil, "unknown", errors.New("database host or port are not specified")
	}

	dbName := "postgres"
	if db.Name != nil && *db.Name != "" {
		dbName = *db.Name
	}

	if ctx.Connection.CredentialId == nil {
		return nil, "unknown", errors.New("password is not set")
	}
	cred, errCred := s.passwordService.GetDecrypted(*ctx.Connection.CredentialId)
	if errCred != nil {
		return nil, "unknown", errors.New("password problems, check if it is exists")
	}

	connProtocol := "postgres://"
	connHost := db.Host + ":" + strconv.Itoa(db.Port) + "/" + dbName
	connUrl := connProtocol + connHost

	if ctx.Connection.Certs != nil {
		// NOTE: verify-ca was chosen, because it potentially can protect from machine-in-the-middle attack if
		// it has right CA policy. More info can be found here https://www.postgresql.org/docs/16/libpq-ssl.html#LIBPQ-SSL-PROTECTION
		connUrl += "?sslmode=verify-ca"
	}

	conConfig, errConfig := pgx.ParseConfig(connUrl)
	conConfig.User = cred.Username
	conConfig.Password = cred.Password
	conConfig.RuntimeParams = map[string]string{
		"application_name": s.GetApplicationName(ctx.Token),
	}
	if errConfig != nil {
		return nil, connUrl, errConfig
	}

	errTls := s.certService.EnrichTLSConfig(&conConfig.TLSConfig, ctx.Connection.Certs)
	if errTls != nil {
		return nil, connUrl, errTls
	}

	conn, err := pgx.ConnectConfig(context.Background(), conConfig)
	return conn, connUrl, err
}

func (s *PostgresClient) closeConnection(conn *pgx.Conn, ctx context.Context) {
	err := conn.Close(ctx)
	if err != nil {
		slog.Warn("postgres close connection", err)
	}
}

func (s *PostgresClient) GetApplicationName(token string) string {
	return s.appName + " [" + s.getTokenSignature(token) + "]"
}

func (s *PostgresClient) getTokenSignature(token string) string {
	signature := strings.SplitN(token, ".", 3)[2]
	return fmt.Sprintf("%.9s", signature)
}
