package query

import (
	"testing"

	"github.com/google/uuid"
)

func TestServiceLogLifecycle(t *testing.T) {
	env := createTestQueryService(t, nil)
	id := uuid.New()

	t.Run("a fresh id has no log", func(t *testing.T) {
		if env.service.HasLog(id) {
			t.Fatalf("expected no log to exist yet")
		}
		log, err := env.service.GetLog(id)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(log) != 0 {
			t.Fatalf("expected an empty log, got %v", log)
		}
	})

	t.Run("adding an entry makes it retrievable", func(t *testing.T) {
		if err := env.service.AddLog(id, DbResponse{Url: "select 1"}); err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if !env.service.HasLog(id) {
			t.Fatalf("expected a log to exist")
		}
		log, err := env.service.GetLog(id)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(log) != 1 || log[0].Url != "select 1" {
			t.Fatalf("unexpected log: %v", log)
		}
	})

	t.Run("deleting the log removes it", func(t *testing.T) {
		if err := env.service.DeleteLog(id); err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if env.service.HasLog(id) {
			t.Fatalf("expected the log to be gone")
		}
	})
}

func TestServiceDeleteAllLogs(t *testing.T) {
	env := createTestQueryService(t, nil)
	id1 := uuid.New()
	id2 := uuid.New()

	if err := env.service.AddLog(id1, DbResponse{Url: "q1"}); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if err := env.service.AddLog(id2, DbResponse{Url: "q2"}); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if err := env.service.DeleteAllLogs(); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if env.service.HasLog(id1) || env.service.HasLog(id2) {
		t.Fatalf("expected all logs to be removed")
	}
}
