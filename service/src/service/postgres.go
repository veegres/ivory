package service

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"golang.org/x/exp/slog"
	. "ivory/src/model"
	"strconv"
	"time"
)

type PostgresClient struct {
	passwordService *PasswordService
	certService     *CertService
}

func NewPostgresClient(
	passwordService *PasswordService,
	certService *CertService,
) *PostgresClient {
	return &PostgresClient{
		passwordService: passwordService,
		certService:     certService,
	}
}

func (s *PostgresClient) GetMany(connection QueryConnection, query string, args ...any) ([]string, error) {
	rows, _, _, errReq := s.sendRequest(connection, query, args...)
	if errReq != nil {
		return nil, errReq
	}
	defer rows.Close()

	values := make([]string, 0)
	for rows.Next() {
		var value string
		err := rows.Scan(&value)
		if err != nil {
			return nil, err
		}
		values = append(values, value)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return values, nil
}

func (s *PostgresClient) GetOne(connection QueryConnection, query string) (any, error) {
	rows, _, _, errReq := s.sendRequest(connection, query)
	if errReq != nil {
		return nil, errReq
	}
	defer rows.Close()

	var value any
	if rows.Next() {
		values, err := rows.Values()
		if err != nil {
			return nil, err
		}
		value = values[0]
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return value, nil
}

func (s *PostgresClient) GetFields(connection QueryConnection, query string, args ...any) (*QueryFields, error) {
	startTime := time.Now().UnixMilli()
	rows, typeMap, url, errReq := s.sendRequest(connection, query, args...)
	if errReq != nil {
		return nil, errReq
	}
	defer rows.Close()

	fields := make([]QueryField, 0)
	for _, field := range rows.FieldDescriptions() {
		dataType, ok := typeMap.TypeForOID(field.DataTypeOID)
		if !ok {
			fields = append(fields, QueryField{Name: field.Name, DataType: "unknown", DataTypeOID: field.DataTypeOID})
		} else {
			fields = append(fields, QueryField{Name: field.Name, DataType: dataType.Name, DataTypeOID: field.DataTypeOID})
		}
	}

	rowList := make([][]any, 0)
	for rows.Next() {
		val, err := rows.Values()
		if err != nil {
			return nil, err
		}
		rowList = append(rowList, val)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
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

func (s *PostgresClient) Cancel(connection QueryConnection, pid int) error {
	_, _, _, err := s.sendRequest(connection, "SELECT pg_cancel_backend("+strconv.Itoa(pid)+")")
	return err
}

func (s *PostgresClient) Terminate(connection QueryConnection, pid int) error {
	_, _, _, err := s.sendRequest(connection, "SELECT pg_terminate_backend("+strconv.Itoa(pid)+")")
	return err
}

func (s *PostgresClient) sendRequest(connection QueryConnection, query string, args ...any) (pgx.Rows, *pgtype.Map, string, error) {
	conn, connUrl, errConn := s.getConnection(connection)
	if errConn != nil {
		return nil, nil, connUrl, errConn
	}
	defer s.closeConnection(conn, context.Background())

	var res pgx.Rows
	var errRes error
	if args == nil {
		res, errRes = conn.Query(context.Background(), query)
	} else {
		res, errRes = conn.Query(context.Background(), query, args...)
	}
	if errRes != nil {
		return nil, nil, connUrl, errRes
	}
	return res, conn.TypeMap(), connUrl, nil
}

func (s *PostgresClient) getConnection(connection QueryConnection) (*pgx.Conn, string, error) {
	db := connection.Db
	if db.Port == 0 || db.Host == "" || db.Host == "-" {
		return nil, "unknown", errors.New("database host or port are not specified")
	}

	dbName := "postgres"
	if db.Name != nil && *db.Name != "" {
		dbName = *db.Name
	}

	cred, errCred := s.passwordService.GetDecrypted(*connection.CredentialId)
	if errCred != nil {
		return nil, "unknown", errors.New("password problems, check if it is exists")
	}

	connProtocol := "postgres://"
	connCredentials := cred.Username + ":" + cred.Password + "@"
	connHost := db.Host + ":" + strconv.Itoa(db.Port) + "/" + dbName
	connUrl := connProtocol + connHost
	connString := connProtocol + connCredentials + connHost

	config, errConfig := pgx.ParseConfig(connString)
	if errConfig != nil {
		return nil, connUrl, errConfig
	}

	if connection.Certs != nil {
		tlsConfig, errTls := s.certService.GetTLSConfig(*connection.Certs)
		if errTls != nil {
			return nil, connUrl, errTls
		}
		config.TLSConfig = tlsConfig
	}

	conn, err := pgx.ConnectConfig(context.Background(), config)
	return conn, connUrl, err
}

func (s *PostgresClient) closeConnection(conn *pgx.Conn, ctx context.Context) {
	err := conn.Close(ctx)
	if err != nil {
		slog.Warn("postgres close connection", err)
	}
}
