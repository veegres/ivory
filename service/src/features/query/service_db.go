package query

import (
	"fmt"
	"ivory/src/plugins/database"

	"github.com/google/uuid"
)

func (s *Service) ConsoleQuery(queryCtx Context, query string, options *DbOptions) (*DbResponse, error) {
	ctx, err := s.mapContext(queryCtx)
	if err != nil {
		return nil, err
	}
	client, err := s.databaseRegistry.Get(ctx.Connection.Config.Type)
	if err != nil {
		return nil, err
	}
	return client.GetFields(ctx, query, options)
}

func (s *Service) TemplateQuery(ctx Context, uuid uuid.UUID, options *DbOptions) (*DbResponse, error) {
	query, errQuery := s.repository.Get(uuid)
	if errQuery != nil {
		return nil, errQuery
	}
	if query.Custom == "" {
		return nil, ErrQueryEmpty
	}
	response, errRun := s.ConsoleQuery(ctx, query.Custom, options)
	if errRun == nil && len(response.Rows) > 0 {
		// NOTE: we don't want fail request if there is some problem with writing to the file
		_ = s.AddLog(uuid, response)
	}
	return response, errRun
}

func (s *Service) CancelQuery(queryCtx Context, pid int) error {
	ctx, err := s.mapContext(queryCtx)
	if err != nil {
		return err
	}
	client, err := s.databaseRegistry.Get(ctx.Connection.Config.Type)
	if err != nil {
		return err
	}
	return client.Cancel(ctx, pid)
}

func (s *Service) TerminateQuery(queryCtx Context, pid int) error {
	ctx, err := s.mapContext(queryCtx)
	if err != nil {
		return err
	}
	client, err := s.databaseRegistry.Get(ctx.Connection.Config.Type)
	if err != nil {
		return err
	}
	return client.Terminate(ctx, pid)
}

func (s *Service) RunningQueriesByApplicationName(queryCtx Context) (*DbResponse, error) {
	ctx, err := s.mapContext(queryCtx)
	if err != nil {
		return nil, err
	}
	client, err := s.databaseRegistry.Get(ctx.Connection.Config.Type)
	if err != nil {
		return nil, err
	}
	options := &database.QueryOptions{Params: []any{ctx.Application}}
	return client.ActiveQueries(ctx, options)
}

func (s *Service) DatabasesQuery(queryCtx Context, name string) ([]string, error) {
	ctx, err := s.mapContext(queryCtx)
	if err != nil {
		return nil, err
	}
	client, err := s.databaseRegistry.Get(ctx.Connection.Config.Type)
	if err != nil {
		return nil, err
	}
	return client.ListDatabases(ctx, name)
}

func (s *Service) SchemasQuery(queryCtx Context, name string) ([]string, error) {
	db := queryCtx.Connection.Db
	if db.Name == nil || *db.Name == "" {
		return []string{}, nil
	}
	ctx, err := s.mapContext(queryCtx)
	if err != nil {
		return nil, err
	}
	client, err := s.databaseRegistry.Get(ctx.Connection.Config.Type)
	if err != nil {
		return nil, err
	}
	return client.ListSchemas(ctx, name)
}

func (s *Service) TablesQuery(queryCtx Context, schema string, name string) ([]string, error) {
	db := queryCtx.Connection.Db
	if db.Name == nil || *db.Name == "" || schema == "" {
		return []string{}, nil
	}
	ctx, err := s.mapContext(queryCtx)
	if err != nil {
		return nil, err
	}
	client, err := s.databaseRegistry.Get(ctx.Connection.Config.Type)
	if err != nil {
		return nil, err
	}
	return client.ListTables(ctx, schema, name)
}

func (s *Service) ChartQuery(queryCtx Context, chartType ChartType) (*Chart, error) {
	dbType := queryCtx.Connection.Db.Type
	typeMap, ok := s.chartMap[dbType]
	if !ok {
		return nil, fmt.Errorf("charts for database type %v are not supported", dbType)
	}
	request, ok := typeMap[chartType]
	if !ok {
		return nil, fmt.Errorf("chart %s is not supported for database type %v", chartType, dbType)
	}
	ctx, err := s.mapContext(queryCtx)
	if err != nil {
		return nil, err
	}
	client, err := s.databaseRegistry.Get(dbType)
	if err != nil {
		return nil, err
	}
	response, err := client.GetOne(ctx, request.Query)
	if err != nil {
		return nil, fmt.Errorf("cannot get %s: %w", request.Name, err)
	}
	return &Chart{Name: request.Name, Value: response}, nil
}

func (s *Service) mapContext(queryCtx Context) (database.Context, error) {
	con := database.Connection{Config: queryCtx.Connection.Db}
	ctx := database.Context{Connection: &con, Application: s.GetApplicationName(queryCtx.Session)}
	if queryCtx.Connection.VaultId != nil {
		cred, errCred := s.vaultService.GetDecrypted(*queryCtx.Connection.VaultId)
		if errCred != nil {
			return ctx, ErrVaultProblems
		}
		con.Credentials = &database.Credentials{Username: cred.Username, Password: cred.Secret}
	}
	if queryCtx.Connection.Certs != nil {
		errTls := s.certService.EnrichTLSConfig(&con.TlsConfig, queryCtx.Connection.Certs)
		if errTls != nil {
			return ctx, errTls
		}
	}
	return ctx, nil
}
