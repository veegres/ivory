package query

import (
	"errors"
	"ivory/plugins/database"
	"testing"
)

type fakeDbRoleAdmin struct {
	listFn    func(ctx database.Context) ([]database.DbRole, error)
	createFn  func(ctx database.Context, name string, password string) error
	updateFn  func(ctx database.Context, name string, password string) error
	deleteFn  func(ctx database.Context, name string) error
}

func (f *fakeDbRoleAdmin) ListDbRoles(ctx database.Context) ([]database.DbRole, error) {
	if f.listFn != nil {
		return f.listFn(ctx)
	}
	return nil, nil
}

func (f *fakeDbRoleAdmin) CreateDbRole(ctx database.Context, name string, password string) error {
	if f.createFn != nil {
		return f.createFn(ctx, name, password)
	}
	return nil
}

func (f *fakeDbRoleAdmin) UpdateDbRolePassword(ctx database.Context, name string, password string) error {
	if f.updateFn != nil {
		return f.updateFn(ctx, name, password)
	}
	return nil
}

func (f *fakeDbRoleAdmin) DeleteDbRole(ctx database.Context, name string) error {
	if f.deleteFn != nil {
		return f.deleteFn(ctx, name)
	}
	return nil
}

type fakeDbDatabaseAdmin struct {
	listFn   func(ctx database.Context) ([]database.DbDatabaseCatalog, error)
	createFn func(ctx database.Context, name string) error
	renameFn func(ctx database.Context, name string, newName string) error
	deleteFn func(ctx database.Context, name string) error
}

func (f *fakeDbDatabaseAdmin) ListDbDatabases(ctx database.Context) ([]database.DbDatabaseCatalog, error) {
	if f.listFn != nil {
		return f.listFn(ctx)
	}
	return nil, nil
}

func (f *fakeDbDatabaseAdmin) CreateDbDatabase(ctx database.Context, name string) error {
	if f.createFn != nil {
		return f.createFn(ctx, name)
	}
	return nil
}

func (f *fakeDbDatabaseAdmin) RenameDbDatabase(ctx database.Context, name string, newName string) error {
	if f.renameFn != nil {
		return f.renameFn(ctx, name, newName)
	}
	return nil
}

func (f *fakeDbDatabaseAdmin) DeleteDbDatabase(ctx database.Context, name string) error {
	if f.deleteFn != nil {
		return f.deleteFn(ctx, name)
	}
	return nil
}

type fakeAdminAdapter struct {
	fakeDatabaseAdapter
	*fakeDbRoleAdmin
	*fakeDbDatabaseAdmin
}

func TestServiceDbRoleListUnsupported(t *testing.T) {
	env := createTestQueryService(t, &fakeDatabaseAdapter{})
	_, err := env.service.DbRoleList(Context{Connection: Connection{Db: DbConfig{Plugin: database.POSTGRES}}})
	if !errors.Is(err, database.ErrAdminNotSupported) {
		t.Errorf("got %v, want ErrAdminNotSupported", err)
	}
}

func TestServiceDbRoleListCatalogContext(t *testing.T) {
	dbName := "app"
	schema := "public"
	var gotCtx database.Context
	admin := &fakeAdminAdapter{
		fakeDbRoleAdmin: &fakeDbRoleAdmin{
			listFn: func(ctx database.Context) ([]database.DbRole, error) {
				gotCtx = ctx
				return []database.DbRole{{Name: "u1", Login: true}}, nil
			},
		},
	}
	env := createTestQueryService(t, admin)
	roles, err := env.service.DbRoleList(Context{Connection: Connection{Db: DbConfig{
		Plugin: database.POSTGRES, Host: "h", Port: 5432, Name: &dbName, Schema: &schema,
	}}})
	if err != nil {
		t.Fatal(err)
	}
	if len(roles) != 1 || roles[0].Name != "u1" {
		t.Fatalf("got %+v", roles)
	}
	if gotCtx.Connection.Config.Name != nil || gotCtx.Connection.Config.Schema != nil {
		t.Errorf("expected catalog connection without db/schema, got name=%v schema=%v",
			gotCtx.Connection.Config.Name, gotCtx.Connection.Config.Schema)
	}
}

func TestServiceDbDatabaseCreate(t *testing.T) {
	var created string
	admin := &fakeAdminAdapter{
		fakeDbDatabaseAdmin: &fakeDbDatabaseAdmin{
			createFn: func(ctx database.Context, name string) error {
				created = name
				return nil
			},
		},
	}
	env := createTestQueryService(t, admin)
	err := env.service.DbDatabaseCreate(Context{Connection: Connection{Db: DbConfig{Plugin: database.POSTGRES, Host: "h", Port: 5432}}}, "newdb")
	if err != nil {
		t.Fatal(err)
	}
	if created != "newdb" {
		t.Errorf("got %q", created)
	}
}
