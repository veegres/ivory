package postgres

import (
	"ivory/plugins/database"
	"testing"
)

func TestIsProtectedRoleName(t *testing.T) {
	if !isProtectedRoleName("postgres") {
		t.Error("postgres should be protected")
	}
	if !isProtectedRoleName("pg_signal_backend") {
		t.Error("pg_* roles should be protected")
	}
	if isProtectedRoleName("app_user") {
		t.Error("app_user should not be protected")
	}
}

func TestIsProtectedDatabaseName(t *testing.T) {
	for _, name := range []string{"postgres", "template0", "template1"} {
		if !isProtectedDatabaseName(name) {
			t.Errorf("%q should be protected", name)
		}
	}
	if isProtectedDatabaseName("mydb") {
		t.Error("mydb should not be protected")
	}
}

func TestQuotePostgresLiteral(t *testing.T) {
	if got := quotePostgresLiteral("abc"); got != "'abc'" {
		t.Errorf("got %q", got)
	}
	if got := quotePostgresLiteral("a'b"); got != "'a''b'" {
		t.Errorf("got %q", got)
	}
}

func TestValidateRoleDeleteSessionRole(t *testing.T) {
	ctx := database.Context{
		Connection: &database.Connection{
			Credentials: &database.Credentials{Username: "ivory"},
		},
	}
	if err := validateRoleDelete(ctx, "ivory"); err != database.ErrCannotDropSessionRole {
		t.Errorf("got %v, want ErrCannotDropSessionRole", err)
	}
	if err := validateRoleDelete(ctx, "other"); err != nil {
		t.Errorf("got %v, want nil", err)
	}
}

func TestValidateRoleDeleteProtected(t *testing.T) {
	ctx := database.Context{Connection: &database.Connection{}}
	if err := validateRoleDelete(ctx, "postgres"); err != database.ErrProtectedRole {
		t.Errorf("got %v, want ErrProtectedRole", err)
	}
}
