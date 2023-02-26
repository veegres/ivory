package service

import (
	"errors"
	"github.com/google/uuid"
	. "ivory/src/model"
	"ivory/src/persistence"
)

type QueryService struct {
	queryRepository *persistence.QueryRepository
	secretService   *SecretService
}

func NewQueryService(queryRepository *persistence.QueryRepository, secretService *SecretService) *QueryService {
	queryService := &QueryService{queryRepository: queryRepository, secretService: secretService}
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

func (s *QueryService) Create(query Query) (uuid.UUID, Query, error) {
	return s.queryRepository.Create(query)
}

func (s *QueryService) Update(key uuid.UUID, query Query) (uuid.UUID, Query, error) {
	return s.queryRepository.Update(key, query)
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

	_, _, errRun := s.queryRepository.Create(s.runningQuery())
	if errRun != nil {
		return errRun
	}
	_, _, errCon := s.queryRepository.Create(s.connectionQuery())
	if errCon != nil {
		return errCon
	}

	return nil
}

const defaultRunningQuery = `SELECT
	pid AS process_id,
	pg_blocking_pids(pid) AS blocked_by_process_id,
	now() - pg_stat_activity.query_start AS duration,
	query,
	state,
	usename,
	application_name,
	client_addr
FROM pg_stat_activity
WHERE now() - pg_stat_activity.query_start IS NOT NULL
AND state <> 'idle'
AND query NOT LIKE 'START_REPLICATION SLOT%'
ORDER BY now() - pg_stat_activity.query_start DESC;`

func (s *QueryService) runningQuery() Query {
	return Query{
		Name:        "Running Queries",
		Type:        ACTIVITY,
		Creation:    System,
		Edited:      nil,
		Description: "This query will check currently running queries",
		Default:     defaultRunningQuery,
	}
}

const defaultConnectionQuery = `SELECT datname, state, count(*)
FROM pg_stat_activity
GROUP BY 1, 2 ORDER BY 1, 2;`

func (s *QueryService) connectionQuery() Query {
	return Query{
		Name:        "All Connections",
		Type:        ACTIVITY,
		Creation:    System,
		Edited:      nil,
		Description: "This query will show list of all connections",
		Default:     defaultConnectionQuery,
	}
}
