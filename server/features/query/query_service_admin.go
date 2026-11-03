package query

import (
	"ivory/plugins/database"
)

func (s *Service) DbRoleList(queryCtx Context) ([]DbRoleResponse, error) {
	client, ctx, err := s.getDatabaseAdapter(catalogContext(queryCtx))
	if err != nil {
		return nil, err
	}
	admin, ok := client.(database.DbRoleAdmin)
	if !ok {
		return nil, database.ErrAdminNotSupported
	}
	roles, err := admin.ListDbRoles(ctx)
	if err != nil {
		return nil, err
	}
	return mapDbRoles(roles), nil
}

func (s *Service) DbRoleCreate(queryCtx Context, name string, password string) error {
	client, ctx, err := s.getDatabaseAdapter(catalogContext(queryCtx))
	if err != nil {
		return err
	}
	admin, ok := client.(database.DbRoleAdmin)
	if !ok {
		return database.ErrAdminNotSupported
	}
	return admin.CreateDbRole(ctx, name, password)
}

func (s *Service) DbRoleUpdatePassword(queryCtx Context, name string, password string) error {
	client, ctx, err := s.getDatabaseAdapter(catalogContext(queryCtx))
	if err != nil {
		return err
	}
	admin, ok := client.(database.DbRoleAdmin)
	if !ok {
		return database.ErrAdminNotSupported
	}
	return admin.UpdateDbRolePassword(ctx, name, password)
}

func (s *Service) DbRoleDelete(queryCtx Context, name string) error {
	client, ctx, err := s.getDatabaseAdapter(catalogContext(queryCtx))
	if err != nil {
		return err
	}
	admin, ok := client.(database.DbRoleAdmin)
	if !ok {
		return database.ErrAdminNotSupported
	}
	return admin.DeleteDbRole(ctx, name)
}

func (s *Service) DbDatabaseList(queryCtx Context) ([]DbDatabaseResponse, error) {
	client, ctx, err := s.getDatabaseAdapter(catalogContext(queryCtx))
	if err != nil {
		return nil, err
	}
	admin, ok := client.(database.DbDatabaseAdmin)
	if !ok {
		return nil, database.ErrAdminNotSupported
	}
	dbs, err := admin.ListDbDatabases(ctx)
	if err != nil {
		return nil, err
	}
	return mapDbDatabases(dbs), nil
}

func (s *Service) DbDatabaseCreate(queryCtx Context, name string) error {
	client, ctx, err := s.getDatabaseAdapter(catalogContext(queryCtx))
	if err != nil {
		return err
	}
	admin, ok := client.(database.DbDatabaseAdmin)
	if !ok {
		return database.ErrAdminNotSupported
	}
	return admin.CreateDbDatabase(ctx, name)
}

func (s *Service) DbDatabaseRename(queryCtx Context, name string, newName string) error {
	client, ctx, err := s.getDatabaseAdapter(catalogContext(queryCtx))
	if err != nil {
		return err
	}
	admin, ok := client.(database.DbDatabaseAdmin)
	if !ok {
		return database.ErrAdminNotSupported
	}
	return admin.RenameDbDatabase(ctx, name, newName)
}

func (s *Service) DbDatabaseDelete(queryCtx Context, name string) error {
	client, ctx, err := s.getDatabaseAdapter(catalogContext(queryCtx))
	if err != nil {
		return err
	}
	admin, ok := client.(database.DbDatabaseAdmin)
	if !ok {
		return database.ErrAdminNotSupported
	}
	return admin.DeleteDbDatabase(ctx, name)
}
