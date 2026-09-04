package user

import (
	"errors"
	"ivory/clients/storage"
	"testing"
	"time"
)

func newTestLinkRepository(t *testing.T) *Repository {
	t.Helper()
	return newTestRepository(t)
}

func testLink(username string) Link {
	now := time.Now()
	return Link{Username: username, CreatedAt: now, ExpiresAt: now.Add(time.Hour)}
}

func TestLinkRepositoryCreateAndGet(t *testing.T) {
	r := newTestLinkRepository(t)

	if _, err := r.LinkCreate("id-1", testLink("alice")); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	link, errGet := r.LinkGet("id-1")
	if errGet != nil {
		t.Fatalf("expected no error, got %v", errGet)
	}
	if link.Username != "alice" {
		t.Fatalf("expected username 'alice', got %q", link.Username)
	}
}

func TestLinkRepositoryGetMissing(t *testing.T) {
	r := newTestLinkRepository(t)

	if _, err := r.LinkGet("nothing"); !errors.Is(err, storage.ErrNotFound) {
		t.Fatalf("expected storage.ErrNotFound, got %v", err)
	}
}

func TestLinkRepositoryMap(t *testing.T) {
	r := newTestLinkRepository(t)

	if _, err := r.LinkCreate("id-1", testLink("alice")); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if _, err := r.LinkCreate("id-2", testLink("bob")); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	linkMap, err := r.LinkMap()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(linkMap) != 2 || linkMap["id-1"].Username != "alice" || linkMap["id-2"].Username != "bob" {
		t.Fatalf("expected both links back, got %+v", linkMap)
	}
}

func TestLinkRepositoryDelete(t *testing.T) {
	r := newTestLinkRepository(t)

	if _, err := r.LinkCreate("id-1", testLink("alice")); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if err := r.LinkDelete("id-1"); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if _, err := r.LinkGet("id-1"); !errors.Is(err, storage.ErrNotFound) {
		t.Fatalf("expected storage.ErrNotFound, got %v", err)
	}
}

func TestLinkRepositoryDeleteAll(t *testing.T) {
	r := newTestLinkRepository(t)

	if _, err := r.LinkCreate("id-1", testLink("alice")); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if err := r.LinkDeleteAll(); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	linkMap, err := r.LinkMap()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(linkMap) != 0 {
		t.Fatalf("expected no links, got %d", len(linkMap))
	}
}
