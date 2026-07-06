package postgres

import (
	"errors"
	"ivory/plugins/database"
	"testing"
)

func noConnectionContext() database.Context {
	return database.Context{Connection: &database.Connection{Config: database.Config{}}}
}

func TestGetManyPropagatesConnectionErrors(t *testing.T) {
	adapter := NewAdapter()

	_, err := adapter.GetMany(noConnectionContext(), "SELECT 1", nil)
	if !errors.Is(err, database.ErrDatabaseHostOrPortNotSpecified) {
		t.Errorf("expected ErrDatabaseHostOrPortNotSpecified, got %v", err)
	}
}

func TestGetOnePropagatesConnectionErrors(t *testing.T) {
	adapter := NewAdapter()

	_, err := adapter.GetOne(noConnectionContext(), "SELECT 1")
	if !errors.Is(err, database.ErrDatabaseHostOrPortNotSpecified) {
		t.Errorf("expected ErrDatabaseHostOrPortNotSpecified, got %v", err)
	}
}

func TestGetFieldsPropagatesConnectionErrors(t *testing.T) {
	adapter := NewAdapter()

	_, err := adapter.GetFields(noConnectionContext(), "SELECT 1", nil)
	if !errors.Is(err, database.ErrDatabaseHostOrPortNotSpecified) {
		t.Errorf("expected ErrDatabaseHostOrPortNotSpecified, got %v", err)
	}
}

func TestGetFieldsRejectsLimitWithoutTrim(t *testing.T) {
	adapter := NewAdapter()
	limit := "100"

	_, err := adapter.GetFields(noConnectionContext(), "SELECT 1", &database.QueryOptions{Limit: &limit})
	if !errors.Is(err, database.ErrCannotLimitWithoutTrim) {
		t.Errorf("expected ErrCannotLimitWithoutTrim, got %v", err)
	}
}
