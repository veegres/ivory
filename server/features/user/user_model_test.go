package user

import (
	"testing"
	"time"
)

func TestUserToResponse(t *testing.T) {
	t.Run("carries the username and the superuser flag, never the password", func(t *testing.T) {
		response := User{Username: "alice", Password: "encrypted", Superuser: true}.toResponse()

		if response.Username != "alice" {
			t.Fatalf("expected username 'alice', got %q", response.Username)
		}
		if !response.Superuser {
			t.Fatalf("expected a superuser")
		}
	})
}

func TestLinkMappers(t *testing.T) {
	created := time.Date(2026, 9, 3, 10, 0, 0, 0, time.UTC)
	expires := created.Add(12 * time.Hour)
	link := Link{Username: "alice", Superuser: true, CreatedAt: created, ExpiresAt: expires}

	t.Run("toResponse carries the id and never a token", func(t *testing.T) {
		response := link.toResponse("id-1")

		if response.Id != "id-1" {
			t.Fatalf("expected id 'id-1', got %q", response.Id)
		}
		if response.Username != "alice" || !response.Superuser {
			t.Fatalf("expected the link facts to survive, got %+v", response)
		}
		if !response.CreatedAt.Equal(created) || !response.ExpiresAt.Equal(expires) {
			t.Fatalf("expected the timestamps to survive, got %+v", response)
		}
	})

	t.Run("toPayload states what the public page shows", func(t *testing.T) {
		payload := link.toPayload()

		if payload.Username != "alice" || !payload.Superuser {
			t.Fatalf("expected the link facts to survive, got %+v", payload)
		}
		if !payload.ExpiresAt.Equal(expires) {
			t.Fatalf("expected the expiration to survive, got %v", payload.ExpiresAt)
		}
	})

	t.Run("toCreatedResponse is the only shape carrying the token", func(t *testing.T) {
		response := link.toCreatedResponse("id-1", "the-token")

		if response.Id != "id-1" || response.Token != "the-token" {
			t.Fatalf("expected the id and the token, got %+v", response)
		}
		if response.Username != "alice" || !response.Superuser {
			t.Fatalf("expected the link facts to survive, got %+v", response)
		}
	})
}
