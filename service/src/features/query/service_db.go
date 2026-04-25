package query

import (
	"fmt"
	"ivory/src/clients/database"
	"ivory/src/clients/database/postgres"

	"github.com/google/uuid"
)

func (s *Service) ConsoleQuery(queryCtx QueryContext, query string, options *database.QueryOptions) (*database.QueryFields, error) {
	ctx, err := s.mapContext(queryCtx)
	if err != nil {
		return nil, err
	}
	client, err := s.databaseRegistry.Get(ctx.Connection.Database.Type)
	if err != nil {
		return nil, err
	}
	return client.GetFields(ctx, query, options)
}

func (s *Service) TemplateQuery(ctx QueryContext, uuid uuid.UUID, options *database.QueryOptions) (*database.QueryFields, error) {
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

func (s *Service) CancelQuery(queryCtx QueryContext, pid int) error {
	ctx, err := s.mapContext(queryCtx)
	if err != nil {
		return err
	}
	client, err := s.databaseRegistry.Get(ctx.Connection.Database.Type)
	if err != nil {
		return err
	}
	return client.Cancel(ctx, pid)
}

func (s *Service) TerminateQuery(queryCtx QueryContext, pid int) error {
	ctx, err := s.mapContext(queryCtx)
	if err != nil {
		return err
	}
	client, err := s.databaseRegistry.Get(ctx.Connection.Database.Type)
	if err != nil {
		return err
	}
	return client.Terminate(ctx, pid)
}

func (s *Service) RunningQueriesByApplicationName(queryCtx QueryContext) (*database.QueryFields, error) {
	ctx, err := s.mapContext(queryCtx)
	if err != nil {
		return nil, err
	}
	client, err := s.databaseRegistry.Get(ctx.Connection.Database.Type)
	if err != nil {
		return nil, err
	}
	options := &database.QueryOptions{Params: []any{client.GetApplicationName(ctx.Session)}}
	return client.GetFields(ctx, postgres.GetAllRunningQueriesByApplicationName, options)
}

func (s *Service) DatabasesQuery(queryCtx QueryContext, name string) ([]string, error) {
	ctx, err := s.mapContext(queryCtx)
	if err != nil {
		return nil, err
	}
	client, err := s.databaseRegistry.Get(ctx.Connection.Database.Type)
	if err != nil {
		return nil, err
	}
	return client.GetMany(ctx, postgres.GetAllDatabases, []any{"%" + name + "%"})
}

func (s *Service) SchemasQuery(queryCtx QueryContext, name string) ([]string, error) {
	db := queryCtx.Connection.Db
	if db.Name == nil || *db.Name == "" {
		return []string{}, nil
	}
	ctx, err := s.mapContext(queryCtx)
	if err != nil {
		return nil, err
	}
	client, err := s.databaseRegistry.Get(ctx.Connection.Database.Type)
	if err != nil {
		return nil, err
	}
	return client.GetMany(ctx, postgres.GetAllSchemas, []any{"%" + name + "%"})
}

func (s *Service) TablesQuery(queryCtx QueryContext, schema string, name string) ([]string, error) {
	db := queryCtx.Connection.Db
	if db.Name == nil || *db.Name == "" || schema == "" {
		return []string{}, nil
	}
	ctx, err := s.mapContext(queryCtx)
	if err != nil {
		return nil, err
	}
	client, err := s.databaseRegistry.Get(ctx.Connection.Database.Type)
	if err != nil {
		return nil, err
	}
	return client.GetMany(ctx, postgres.GetAllTables, []any{schema, "%" + name + "%"})
}

func (s *Service) ChartQuery(queryCtx QueryContext, chartType database.QueryChartType) (*QueryChart, error) {
	request, ok := s.chartMap[chartType]
	if !ok {
		return nil, fmt.Errorf("chart %s is not supported", chartType)
	}
	ctx, err := s.mapContext(queryCtx)
	if err != nil {
		return nil, err
	}
	client, err := s.databaseRegistry.Get(ctx.Connection.Database.Type)
	if err != nil {
		return nil, err
	}
	response, err := client.GetOne(ctx, request.Query)
	if err != nil {
		return nil, fmt.Errorf("cannot get %s: %w", request.Name, err)
	}
	return &QueryChart{Name: request.Name, Value: response}, nil
}

func (s *Service) mapContext(queryCtx QueryContext) (database.Context, error) {
	con := database.Connection{Database: queryCtx.Connection.Db}
	ctx := database.Context{Connection: &con, Session: queryCtx.Session}
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
