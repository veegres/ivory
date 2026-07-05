package query

import (
	"errors"
	"ivory/plugins/database"

	"github.com/google/uuid"
)

// GetList filters by type and plugin when provided; nil means no filtering
// for that dimension. Records stored before the plugin field existed are
// returned as postgres queries.
func (s *Service) GetList(queryType *Type, plugin *DbPlugin) ([]Response, error) {
	list, err := s.repository.ListByFilter(queryType, plugin)
	if err != nil {
		return nil, err
	}
	for i := range list {
		if list[i].Plugin == "" {
			list[i].Plugin = database.POSTGRES
		}
	}
	return list, nil
}

func (s *Service) Create(creation CreationType, query Request) (*uuid.UUID, *Response, error) {
	if query.Name == "" || query.Type == nil || query.Query == "" {
		return nil, nil, ErrAllFieldsRequired
	}

	plugin := query.Plugin
	if plugin == "" {
		plugin = database.POSTGRES
	}

	return s.repository.Create(Response{
		Name:        query.Name,
		Type:        *query.Type,
		Plugin:      plugin,
		Creation:    creation,
		Description: query.Description,
		Default:     query.Query,
		Custom:      query.Query,
		Params:      query.Params,
		Varieties:   query.Varieties,
	})
}

func (s *Service) Update(key uuid.UUID, query Request) (*uuid.UUID, *Response, error) {
	currentQuery, err := s.repository.Get(key)
	if err != nil {
		return nil, nil, err
	}
	if currentQuery.Creation == System {
		if query.Name != currentQuery.Name {
			return nil, nil, ErrNameChangeNotAllowed
		}
		if *query.Type != currentQuery.Type {
			return nil, nil, ErrTypeChangeNotAllowed
		}
		if *query.Description != *currentQuery.Description {
			return nil, nil, ErrDescriptionChangeNotAllowed
		}

		return s.repository.Update(key, Response{
			Id:          key,
			Name:        currentQuery.Name,
			Type:        currentQuery.Type,
			Plugin:      currentQuery.Plugin,
			Creation:    currentQuery.Creation,
			Description: currentQuery.Description,
			Default:     currentQuery.Default,
			Custom:      query.Query,
			Varieties:   query.Varieties,
			Params:      query.Params,
			CreatedAt:   currentQuery.CreatedAt,
		})
	}

	n := currentQuery.Name
	t := currentQuery.Type
	d := currentQuery.Description

	if query.Name != "" {
		n = query.Name
	}
	if query.Type != nil {
		t = *query.Type
	}
	if query.Description != nil {
		d = query.Description
	}

	// NOTE: plugin is immutable, it always stays as it was on creation
	return s.repository.Update(key, Response{
		Id:          key,
		Name:        n,
		Type:        t,
		Plugin:      currentQuery.Plugin,
		Creation:    currentQuery.Creation,
		Description: d,
		Default:     currentQuery.Default,
		Custom:      query.Query,
		Varieties:   query.Varieties,
		Params:      query.Params,
		CreatedAt:   currentQuery.CreatedAt,
	})
}

func (s *Service) Delete(key uuid.UUID) error {
	currentQuery, err := s.repository.Get(key)
	if err != nil {
		return err
	}
	if currentQuery.Creation == System {
		return ErrDeletionOfSystemQueriesRestricted
	}

	var errLog error
	if s.HasLog(key) {
		errLog = s.DeleteLog(key)
	}
	errBucket := s.repository.Delete(key)
	return errors.Join(errLog, errBucket)
}

func (s *Service) DeleteAll() error {
	errLog := s.DeleteAllLogs()
	errDel := s.repository.DeleteAll()
	errDefQueries := s.initializeSystemQueries()
	return errors.Join(errLog, errDel, errDefQueries)
}
