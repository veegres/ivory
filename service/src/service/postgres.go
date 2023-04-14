package service

import (
	"context"
	"errors"
	"github.com/golang/glog"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	. "ivory/src/model"
	"strconv"
)

type PostgresGateway struct {
	passwordService *PasswordService
}

func NewPostgresGateway(
	passwordService *PasswordService,
) *PostgresGateway {
	return &PostgresGateway{
		passwordService: passwordService,
	}
}

func (s *PostgresGateway) GetMany(credentialId uuid.UUID, db Database, query string, args ...any) ([]string, error) {
	rows, _, errReq := s.sendRequest(credentialId, db, query, args...)
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

	return values, nil
}

func (s *PostgresGateway) GetOne(credentialId uuid.UUID, db Database, query string) (any, error) {
	rows, _, errReq := s.sendRequest(credentialId, db, query)
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

	return value, nil
}

func (s *PostgresGateway) GetFields(credentialId uuid.UUID, db Database, query string) (*QueryRunResponse, error) {
	rows, typeMap, errReq := s.sendRequest(credentialId, db, query)
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

	res := &QueryRunResponse{
		Fields: fields,
		Rows:   rowList,
	}
	return res, nil
}

func (s *PostgresGateway) Cancel(pid int, credentialId uuid.UUID, db Database) error {
	_, _, err := s.sendRequest(credentialId, db, "SELECT pg_cancel_backend("+strconv.Itoa(pid)+")")
	return err
}

func (s *PostgresGateway) Terminate(pid int, credentialId uuid.UUID, db Database) error {
	_, _, err := s.sendRequest(credentialId, db, "SELECT pg_terminate_backend("+strconv.Itoa(pid)+")")
	return err
}

func (s *PostgresGateway) sendRequest(credentialId uuid.UUID, db Database, query string, args ...any) (pgx.Rows, *pgtype.Map, error) {
	conn, errConn := s.getConnection(credentialId, db)
	if errConn != nil {
		return nil, nil, errConn
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
		return nil, nil, errRes
	}
	return res, conn.TypeMap(), nil
}

func (s *PostgresGateway) getConnection(credentialId uuid.UUID, db Database) (*pgx.Conn, error) {
	cred, errCred := s.passwordService.GetDecrypted(credentialId)
	if errCred != nil {
		return nil, errCred
	}

	dbName := "postgres"
	if db.Port == 0 || db.Host == "" || db.Host == "-" {
		return nil, errors.New("database host or port are not specified")
	}
	if db.Database != nil && *db.Database != "" {
		dbName = *db.Database
	}
	connString := "postgres://" + cred.Username + ":" + cred.Password + "@" + db.Host + ":" + strconv.Itoa(db.Port) + "/" + dbName

	// TODO think about cache connections
	return pgx.Connect(context.Background(), connString)
}

func (s *PostgresGateway) closeConnection(conn *pgx.Conn, ctx context.Context) {
	err := conn.Close(ctx)
	if err != nil {
		glog.Warning(err)
	}
}
