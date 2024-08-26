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

func (s *PostgresClient) GetFields(ctx QueryContext, query string, queryParams []any) (*QueryFields, error) {
	startTime := time.Now().UnixMilli()

	fields := make([]QueryField, 0)
	rowList := make([][]any, 0)
	url := "-"

	errReq := s.sendRequest(ctx, query, queryParams, func(rows pgx.Rows, typeMap *pgtype.Map, connUrl string) error {
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

	endTime := time.Now().UnixMilli()
	res := &QueryFields{
		Fields:    fields,
		Rows:      rowList,
		Url:       url,
		StartTime: startTime,
		EndTime:   endTime,
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

func (s *PostgresClient) getConnection(ctx QueryContext) (*pgx.Conn, string, error) {
	db := ctx.Connection.Db
	if db.Port == 0 || db.Host == "" || db.Host == "-" {
		return nil, "unknown", errors.New("database host or port are not specified")
	}

	dbName := "postgres"
	if db.Name != nil && *db.Name != "" {
		dbName = *db.Name
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
		"application_name": s.getApplicationName(ctx.Token),
	}
	if errConfig != nil {
		return nil, connUrl, errConfig
	}

	if ctx.Connection.Certs != nil {
		errTls := s.certService.SetTLSConfig(conConfig.TLSConfig, *ctx.Connection.Certs)
		if errTls != nil {
			return nil, connUrl, errTls
		}
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

func (s *PostgresClient) getApplicationName(token string) string {
	return s.appName + " [" + s.getTokenSignature(token) + "]"
}

func (s *PostgresClient) getTokenSignature(token string) string {
	signature := strings.SplitN(token, ".", 3)[2]
	return fmt.Sprintf("%.9s", signature)
}
