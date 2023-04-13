package service

import (
	"context"
	"errors"
	"github.com/golang/glog"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	. "ivory/src/model"
	"strconv"
)

type PostgresGateway struct {
	clusterService  *ClusterService
	passwordService *PasswordService
}

func NewPostgresGateway(
	clusterService *ClusterService,
	passwordService *PasswordService,
) *PostgresGateway {
	return &PostgresGateway{
		clusterService:  clusterService,
		passwordService: passwordService,
	}
}

func (s *PostgresGateway) GetMany(clusterName string, db Database, query string, args ...any) ([]string, error) {
	rows, _, errReq := s.sendRequest(clusterName, db, query, args...)
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

func (s *PostgresGateway) GetOne(clusterName string, db Database, query string) (any, error) {
	rows, _, errReq := s.sendRequest(clusterName, db, query)
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

func (s *PostgresGateway) GetFields(clusterName string, db Database, query string) (*QueryRunResponse, error) {
	rows, typeMap, errReq := s.sendRequest(clusterName, db, query)
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

func (s *PostgresGateway) Cancel(pid int, clusterName string, db Database) error {
	_, _, err := s.sendRequest(clusterName, db, "SELECT pg_cancel_backend("+strconv.Itoa(pid)+")")
	return err
}

func (s *PostgresGateway) Terminate(pid int, clusterName string, db Database) error {
	_, _, err := s.sendRequest(clusterName, db, "SELECT pg_terminate_backend("+strconv.Itoa(pid)+")")
	return err
}

func (s *PostgresGateway) sendRequest(clusterName string, db Database, query string, args ...any) (pgx.Rows, *pgtype.Map, error) {
	conn, errConn := s.getConnection(clusterName, db)
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

func (s *PostgresGateway) getConnection(clusterName string, db Database) (*pgx.Conn, error) {
	// TODO think about cache connections
	cluster, errCluster := s.clusterService.Get(clusterName)
	if errCluster != nil {
		return nil, errCluster
	}
	if cluster.Credentials.PostgresId == nil {
		return nil, errors.New("there is no password for this cluster")
	}

	cred, errCred := s.passwordService.GetDecrypted(*cluster.Credentials.PostgresId)
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
	return pgx.Connect(context.Background(), connString)
}

func (s *PostgresGateway) closeConnection(conn *pgx.Conn, ctx context.Context) {
	err := conn.Close(ctx)
	if err != nil {
		glog.Warning(err)
	}
}
