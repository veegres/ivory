package query

import (
	"errors"

	"github.com/google/uuid"
)

func (s *Service) GetList(queryType *Type) ([]Response, error) {
	if queryType == nil {
		return s.repository.List()
	}

	return s.repository.ListByType(*queryType)
}

func (s *Service) Create(creation CreationType, query Request) (*uuid.UUID, *Response, error) {
	if query.Name == "" || query.Type == nil || query.Query == "" {
		return nil, nil, ErrAllFieldsRequired
	}

	return s.repository.Create(Response{
		Name:        query.Name,
		Type:        *query.Type,
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

	return s.repository.Update(key, Response{
		Id:          key,
		Name:        n,
		Type:        t,
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
