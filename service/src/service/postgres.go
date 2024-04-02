package service

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"golang.org/x/exp/slog"
	. "ivory/src/model"
	"strconv"
	"time"
)

type PostgresClient struct {
	passwordService *PasswordService
}

func NewPostgresClient(
	passwordService *PasswordService,
) *PostgresClient {
	return &PostgresClient{
		passwordService: passwordService,
	}
}

func (s *PostgresClient) GetMany(credentialId uuid.UUID, db Database, query string, args ...any) ([]string, error) {
	rows, _, _, errReq := s.sendRequest(credentialId, db, query, args...)
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

func (s *PostgresClient) GetOne(credentialId uuid.UUID, db Database, query string) (any, error) {
	rows, _, _, errReq := s.sendRequest(credentialId, db, query)
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

func (s *PostgresClient) GetFields(credentialId uuid.UUID, db Database, query string, args ...any) (*QueryFields, error) {
	startTime := time.Now().UnixMilli()
	rows, typeMap, url, errReq := s.sendRequest(credentialId, db, query, args...)
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

func (s *PostgresClient) Cancel(pid int, credentialId uuid.UUID, db Database) error {
	_, _, _, err := s.sendRequest(credentialId, db, "SELECT pg_cancel_backend("+strconv.Itoa(pid)+")")
	return err
}

func (s *PostgresClient) Terminate(pid int, credentialId uuid.UUID, db Database) error {
	_, _, _, err := s.sendRequest(credentialId, db, "SELECT pg_terminate_backend("+strconv.Itoa(pid)+")")
	return err
}

func (s *PostgresClient) sendRequest(credentialId uuid.UUID, db Database, query string, args ...any) (pgx.Rows, *pgtype.Map, string, error) {
	conn, connUrl, errConn := s.getConnection(credentialId, db)
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

func (s *PostgresClient) getConnection(credentialId uuid.UUID, db Database) (*pgx.Conn, string, error) {
	cred, errCred := s.passwordService.GetDecrypted(credentialId)
	if errCred != nil {
		return nil, "unknown", errors.New("password problems, check if it is exists")
	}

	if db.Port == 0 || db.Host == "" || db.Host == "-" {
		return nil, "unknown", errors.New("database host or port are not specified")
	}

	dbName := "postgres"
	if db.Name != nil && *db.Name != "" {
		dbName = *db.Name
	}

	connProtocol := "postgres://"
	connCredentials := cred.Username + ":" + cred.Password + "@"
	connHost := db.Host + ":" + strconv.Itoa(db.Port) + "/" + dbName
	connUrl := connProtocol + connHost
	connString := connProtocol + connCredentials + connHost

	// TODO think about to cache connections, check if we have cache provided by lib
	conn, err := pgx.Connect(context.Background(), connString)
	return conn, connUrl, err
}

func (s *PostgresClient) closeConnection(conn *pgx.Conn, ctx context.Context) {
	err := conn.Close(ctx)
	if err != nil {
		slog.Warn("postgres close connection", err)
	}
}
