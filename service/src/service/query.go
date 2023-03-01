package service

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"ivory/src/constant"
	. "ivory/src/model"
	"ivory/src/persistence"
	"strconv"
)

type QueryService struct {
	queryRepository   *persistence.QueryRepository
	clusterRepository *persistence.ClusterRepository
	passwordService   *PasswordService
	secretService     *SecretService
}

func NewQueryService(
	queryRepository *persistence.QueryRepository,
	clusterRepository *persistence.ClusterRepository,
	passwordService *PasswordService,
	secretService *SecretService,
) *QueryService {
	queryService := &QueryService{
		queryRepository:   queryRepository,
		clusterRepository: clusterRepository,
		passwordService:   passwordService,
		secretService:     secretService,
	}
	err := queryService.createDefaultQueries()
	if err != nil {
		panic("Cannot create default queries: " + err.Error())
	}
	return queryService
}

func (s *QueryService) GetMap(queryType *QueryType) (map[string]Query, error) {
	if queryType == nil {
		return s.queryRepository.Map()
	} else {
		return s.queryRepository.MapByType(*queryType)
	}
}

func (s *QueryService) Run(queryUuid uuid.UUID, clusterName string, db Database) ([]QueryField, [][]any, error) {
	query, errQuery := s.queryRepository.Get(queryUuid)
	if errQuery != nil {
		return nil, nil, errQuery
	}
	cluster, errCluster := s.clusterRepository.Get(clusterName)
	if errCluster != nil {
		return nil, nil, errCluster
	}
	if cluster.Credentials.PostgresId == nil {
		return nil, nil, errors.New("there is no password for this cluster")
	}

	cred, errCred := s.passwordService.GetDecrypted(*cluster.Credentials.PostgresId)
	if errCred != nil {
		return nil, nil, errCred
	}

	connString := "postgres://" + cred.Username + ":" + cred.Password + "@" + db.Host + ":" + strconv.Itoa(db.Port)
	conn, errConn := pgx.Connect(context.Background(), connString)
	if errConn != nil {
		return nil, nil, errConn
	}
	defer conn.Close(context.Background())

	res, errRes := conn.Query(context.Background(), query.Custom)
	if errRes != nil {
		return nil, nil, errRes
	}
	defer res.Close()

	typeMap := conn.TypeMap()
	fields := make([]QueryField, 0)
	for _, field := range res.FieldDescriptions() {
		dataType, ok := typeMap.TypeForOID(field.DataTypeOID)
		if !ok {
			return nil, nil, errors.New("cannot parse data type for field: " + field.Name)
		}
		fields = append(fields, QueryField{Name: field.Name, DataType: dataType.Name, DataTypeOID: field.DataTypeOID})
	}

	rows := make([][]any, 0)
	for res.Next() {
		val, err := res.Values()
		if err != nil {
			return nil, nil, err
		}
		rows = append(rows, val)
	}

	return fields, rows, nil
}

func (s *QueryService) Create(creation QueryCreation, query QueryRequest) (*uuid.UUID, *Query, error) {
	if query.Name == nil || query.Type == nil || query.Description == nil {
		return nil, nil, errors.New("all fields have to be filled")
	}

	return s.queryRepository.Create(Query{
		Name:        *query.Name,
		Type:        *query.Type,
		Creation:    creation,
		Description: *query.Description,
		Default:     query.Query,
		Custom:      query.Query,
	})
}

func (s *QueryService) Update(key uuid.UUID, query QueryRequest) (*uuid.UUID, *Query, error) {
	currentQuery, err := s.queryRepository.Get(key)
	if err != nil {
		return nil, nil, err
	}
	if currentQuery.Creation == System {
		if query.Name != nil {
			return nil, nil, errors.New("name change is not allowed for creation type system")
		}
		if query.Type != nil {
			return nil, nil, errors.New("type change is not allowed for creation type system")
		}
		if query.Description != nil {
			return nil, nil, errors.New("description change is not allowed for creation type system")
		}

		return s.queryRepository.Update(key, Query{
			Name:        currentQuery.Name,
			Type:        currentQuery.Type,
			Creation:    currentQuery.Creation,
			Description: currentQuery.Description,
			Default:     currentQuery.Default,
			Custom:      query.Query,
		})
	} else {
		n := currentQuery.Name
		t := currentQuery.Type
		d := currentQuery.Description

		if query.Name != nil {
			n = *query.Name
		}
		if query.Type != nil {
			t = *query.Type
		}
		if query.Description != nil {
			d = *query.Description
		}

		return s.queryRepository.Update(key, Query{
			Name:        n,
			Type:        t,
			Creation:    currentQuery.Creation,
			Description: d,
			Default:     currentQuery.Default,
			Custom:      query.Query,
		})
	}
}

func (s *QueryService) Delete(key uuid.UUID) error {
	currentQuery, err := s.queryRepository.Get(key)
	if err != nil {
		return err
	}
	if currentQuery.Creation == Manual {
		return s.queryRepository.Delete(key)
	} else {
		return errors.New("deletion of system queries is restricted")
	}
}

func (s *QueryService) DeleteAll() error {
	errDel := s.queryRepository.DeleteAll()
	errDefQueries := s.createDefaultQueries()
	return errors.Join(errDel, errDefQueries)
}

func (s *QueryService) createDefaultQueries() error {
	if !s.secretService.IsRefEmpty() {
		return nil
	}

	_, _, errRun := s.Create(System, s.runningQuery())
	if errRun != nil {
		return errRun
	}
	_, _, errCon := s.Create(System, s.connectionQuery())
	if errCon != nil {
		return errCon
	}

	return nil
}

func (s *QueryService) runningQuery() QueryRequest {
	n, t, d := "Running Queries", ACTIVITY, "This query will check currently running queries"
	return QueryRequest{Name: &n, Type: &t, Description: &d, Query: constant.DefaultRunningQuery}
}

func (s *QueryService) connectionQuery() QueryRequest {
	n, t, d := "All Connections", ACTIVITY, "This query will show list of all connections"
	return QueryRequest{Name: &n, Type: &t, Description: &d, Query: constant.DefaultConnectionQuery}
}
