package user

import (
	"errors"
	"ivory/clients/storage"
	"ivory/core/config"
	"ivory/features/permission"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func newTestLinkEnv(t *testing.T) *testUserEnv {
	t.Helper()
	return newTestUserEnv(t)
}

func TestLinkServiceCreate(t *testing.T) {
	t.Run("issues a token and stores the record behind it", func(t *testing.T) {
		env := newTestLinkEnv(t)

		link, err := env.service.LinkCreateInvite(LinkRequest{Username: "alice", Superuser: true}, "")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if link.Token == "" || link.Id == "" {
			t.Fatalf("expected a token and an id, got %+v", link)
		}
		if link.Username != "alice" || !link.Superuser {
			t.Fatalf("expected the request to be reflected, got %+v", link)
		}
		if !link.ExpiresAt.After(time.Now()) {
			t.Fatalf("expected a future expiration, got %v", link.ExpiresAt)
		}

		stored, errGet := env.repository.LinkGet(link.Id)
		if errGet != nil {
			t.Fatalf("expected the link to be stored, got %v", errGet)
		}
		if stored.Username != "alice" {
			t.Fatalf("expected the stored link to name alice, got %+v", stored)
		}
	})

	t.Run("the token carries the username and the record id", func(t *testing.T) {
		env := newTestLinkEnv(t)

		link, err := env.service.LinkCreateInvite(LinkRequest{Username: "alice"}, "")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		claims, errParse := env.tokenService.Parse(link.Token)
		if errParse != nil {
			t.Fatalf("expected the token to parse, got %v", errParse)
		}
		if claims["sub"] != "alice" {
			t.Fatalf("expected subject 'alice', got %v", claims["sub"])
		}
		if claims["jti"] != link.Id {
			t.Fatalf("expected jti %q, got %v", link.Id, claims["jti"])
		}
	})

	t.Run("trims the username", func(t *testing.T) {
		env := newTestLinkEnv(t)

		link, err := env.service.LinkCreateInvite(LinkRequest{Username: "  alice  "}, "")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if link.Username != "alice" {
			t.Fatalf("expected username 'alice', got %q", link.Username)
		}
	})

	t.Run("rejects an empty username", func(t *testing.T) {
		env := newTestLinkEnv(t)

		if _, err := env.service.LinkCreateInvite(LinkRequest{Username: "  "}, ""); !errors.Is(err, ErrUsernameCannotBeEmpty) {
			t.Fatalf("expected ErrUsernameCannotBeEmpty, got %v", err)
		}
	})

	t.Run("rejects a username that already has a user", func(t *testing.T) {
		env := newTestLinkEnv(t)
		if _, err := env.service.Create("alice", "password123", false); err != nil {
			t.Fatalf("failed to seed alice: %v", err)
		}

		if _, err := env.service.LinkCreateInvite(LinkRequest{Username: "alice"}, ""); !errors.Is(err, ErrUserAlreadyExists) {
			t.Fatalf("expected ErrUserAlreadyExists, got %v", err)
		}
	})
}

func TestLinkServiceCreateSuperuserIsSuperusersOwnRight(t *testing.T) {
	t.Run("a regular user cannot invite a superuser", func(t *testing.T) {
		env := newTestLinkEnv(t)
		if _, err := env.service.Create("root", "password123", true); err != nil {
			t.Fatalf("failed to seed root: %v", err)
		}
		if _, err := env.service.Create("alice", "password123", false); err != nil {
			t.Fatalf("failed to seed alice: %v", err)
		}

		_, err := env.service.LinkCreateInvite(LinkRequest{Username: "newbie", Superuser: true}, "alice")
		if !errors.Is(err, ErrSuperuserRequired) {
			t.Fatalf("expected ErrSuperuserRequired, got %v", err)
		}
	})

	t.Run("a regular user can invite a regular user", func(t *testing.T) {
		env := newTestLinkEnv(t)
		if _, err := env.service.Create("alice", "password123", false); err != nil {
			t.Fatalf("failed to seed alice: %v", err)
		}

		if _, err := env.service.LinkCreateInvite(LinkRequest{Username: "newbie"}, "alice"); err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
	})

	t.Run("a superuser can invite a superuser", func(t *testing.T) {
		env := newTestLinkEnv(t)
		if _, err := env.service.Create("root", "password123", true); err != nil {
			t.Fatalf("failed to seed root: %v", err)
		}

		link, err := env.service.LinkCreateInvite(LinkRequest{Username: "newbie", Superuser: true}, "root")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if !link.Superuser {
			t.Fatalf("expected a superuser link, got %+v", link)
		}
	})

	t.Run("a regular user cannot revoke a superuser invitation", func(t *testing.T) {
		env := newTestLinkEnv(t)
		if _, err := env.service.Create("root", "password123", true); err != nil {
			t.Fatalf("failed to seed root: %v", err)
		}
		if _, err := env.service.Create("alice", "password123", false); err != nil {
			t.Fatalf("failed to seed alice: %v", err)
		}
		link, errCreate := env.service.LinkCreateInvite(LinkRequest{Username: "newbie", Superuser: true}, "root")
		if errCreate != nil {
			t.Fatalf("failed to create a link: %v", errCreate)
		}

		if err := env.service.LinkRevoke(link.Id, "alice"); !errors.Is(err, ErrSuperuserRequired) {
			t.Fatalf("expected ErrSuperuserRequired, got %v", err)
		}
		if err := env.service.LinkRevoke(link.Id, "root"); err != nil {
			t.Fatalf("expected a superuser to revoke it, got %v", err)
		}
	})
}

func TestLinkServiceReset(t *testing.T) {
	t.Run("sets a new password without touching anything else about the user", func(t *testing.T) {
		env := newTestLinkEnv(t)
		if _, err := env.service.Create("root", "password123", true); err != nil {
			t.Fatalf("failed to seed root: %v", err)
		}
		if _, err := env.service.Create("alice", "password123", false); err != nil {
			t.Fatalf("failed to seed alice: %v", err)
		}
		if _, err := env.permissionService.CreateUserPermissions("basic", "alice"); err != nil {
			t.Fatalf("failed to seed permissions: %v", err)
		}
		if err := env.permissionService.ApproveUserPermissions("basic:alice", []config.Feature{config.ViewClusterList}); err != nil {
			t.Fatalf("failed to grant a permission: %v", err)
		}

		link, errCreate := env.service.LinkCreateReset("alice", "root")
		if errCreate != nil {
			t.Fatalf("expected no error, got %v", errCreate)
		}
		if link.Kind != LinkReset {
			t.Fatalf("expected a reset link, got %+v", link)
		}

		user, errApply := env.service.LinkApply(link.Token, "newpassword")
		if errApply != nil {
			t.Fatalf("expected no error, got %v", errApply)
		}
		if user.Username != "alice" || user.Superuser {
			t.Fatalf("expected alice to stay a regular user, got %+v", user)
		}
		if err := env.service.VerifyPassword("alice", "newpassword"); err != nil {
			t.Fatalf("expected the new password to verify, got %v", err)
		}
		perms, errPerms := env.permissionService.GetUserPermissions("basic", "alice", false)
		if errPerms != nil {
			t.Fatalf("expected no error, got %v", errPerms)
		}
		if perms[config.ViewClusterList] != permission.GRANTED {
			t.Fatalf("expected the granted permission to survive a reset, got %v", perms[config.ViewClusterList])
		}
	})

	t.Run("is refused for a user Ivory does not have", func(t *testing.T) {
		env := newTestLinkEnv(t)

		_, err := env.service.LinkCreateReset("nobody", "")
		if !errors.Is(err, ErrUserNotFound) {
			t.Fatalf("expected ErrUserNotFound, got %v", err)
		}
	})

	t.Run("an invitation is refused for a user that already exists", func(t *testing.T) {
		env := newTestLinkEnv(t)
		if _, err := env.service.Create("alice", "password123", false); err != nil {
			t.Fatalf("failed to seed alice: %v", err)
		}

		_, err := env.service.LinkCreateInvite(LinkRequest{Username: "alice"}, "")
		if !errors.Is(err, ErrUserAlreadyExists) {
			t.Fatalf("expected ErrUserAlreadyExists, got %v", err)
		}
	})

	t.Run("resetting a superuser is a superuser's own right", func(t *testing.T) {
		env := newTestLinkEnv(t)
		if _, err := env.service.Create("root", "password123", true); err != nil {
			t.Fatalf("failed to seed root: %v", err)
		}
		if _, err := env.service.Create("alice", "password123", false); err != nil {
			t.Fatalf("failed to seed alice: %v", err)
		}

		if _, err := env.service.LinkCreateReset("root", "alice"); !errors.Is(err, ErrSuperuserRequired) {
			t.Fatalf("expected ErrSuperuserRequired, got %v", err)
		}
		link, errCreate := env.service.LinkCreateReset("root", "root")
		if errCreate != nil {
			t.Fatalf("expected a superuser to be allowed, got %v", errCreate)
		}
		if !link.Superuser {
			t.Fatalf("expected the link to report the account it names as a superuser, got %+v", link)
		}
	})

	t.Run("a reset link for a user deleted in the meantime cannot be spent", func(t *testing.T) {
		env := newTestLinkEnv(t)
		if _, err := env.service.Create("root", "password123", true); err != nil {
			t.Fatalf("failed to seed root: %v", err)
		}
		if _, err := env.service.Create("alice", "password123", false); err != nil {
			t.Fatalf("failed to seed alice: %v", err)
		}
		link, errCreate := env.service.LinkCreateReset("alice", "root")
		if errCreate != nil {
			t.Fatalf("failed to create a link: %v", errCreate)
		}
		if err := env.service.Delete("alice", "root"); err != nil {
			t.Fatalf("failed to delete alice: %v", err)
		}

		if _, err := env.service.LinkApply(link.Token, "newpassword"); !errors.Is(err, ErrUserNotFound) {
			t.Fatalf("expected ErrUserNotFound, got %v", err)
		}
	})

	t.Run("a reset link cannot open an account created after it was issued", func(t *testing.T) {
		env := newTestLinkEnv(t)
		if _, err := env.service.Create("root", "password123", true); err != nil {
			t.Fatalf("failed to seed root: %v", err)
		}
		if _, err := env.service.Create("alice", "password123", false); err != nil {
			t.Fatalf("failed to seed alice: %v", err)
		}
		link, errCreate := env.service.LinkCreateReset("alice", "root")
		if errCreate != nil {
			t.Fatalf("failed to create a link: %v", errCreate)
		}
		if err := env.service.Delete("alice", "root"); err != nil {
			t.Fatalf("failed to delete alice: %v", err)
		}
		// the same name, a different person
		time.Sleep(2 * time.Millisecond)
		if _, err := env.service.Create("alice", "somebodyelse", false); err != nil {
			t.Fatalf("failed to re-create alice: %v", err)
		}

		if _, err := env.service.LinkApply(link.Token, "takeover"); !errors.Is(err, ErrLinkObsolete) {
			t.Fatalf("expected ErrLinkObsolete, got %v", err)
		}
		if err := env.service.VerifyPassword("alice", "somebodyelse"); err != nil {
			t.Fatalf("expected the new account's password to be untouched, got %v", err)
		}
	})

	t.Run("an invitation left over from before the user existed cannot reset them", func(t *testing.T) {
		env := newTestLinkEnv(t)
		link, errCreate := env.service.LinkCreateInvite(LinkRequest{Username: "alice"}, "")
		if errCreate != nil {
			t.Fatalf("failed to create a link: %v", errCreate)
		}
		if _, err := env.service.Create("alice", "password123", false); err != nil {
			t.Fatalf("failed to seed alice: %v", err)
		}

		if _, err := env.service.LinkApply(link.Token, "takeover"); !errors.Is(err, ErrUserAlreadyExists) {
			t.Fatalf("expected ErrUserAlreadyExists, got %v", err)
		}
		if err := env.service.VerifyPassword("alice", "password123"); err != nil {
			t.Fatalf("expected the password to be untouched, got %v", err)
		}
	})

	t.Run("the kind is reported to the page the link opens", func(t *testing.T) {
		env := newTestLinkEnv(t)
		if _, err := env.service.Create("alice", "password123", false); err != nil {
			t.Fatalf("failed to seed alice: %v", err)
		}
		link, errCreate := env.service.LinkCreateReset("alice", "")
		if errCreate != nil {
			t.Fatalf("failed to create a link: %v", errCreate)
		}

		payload, err := env.service.LinkVerify(link.Token)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if payload.Kind != LinkReset {
			t.Fatalf("expected a reset payload, got %+v", payload)
		}
	})

	t.Run("an invitation says so in what it stores", func(t *testing.T) {
		env := newTestLinkEnv(t)

		link, err := env.service.LinkCreateInvite(LinkRequest{Username: "alice"}, "")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if link.Kind != LinkInvite {
			t.Fatalf("expected an invitation, got %+v", link)
		}
	})
}

func TestLinkServiceVerify(t *testing.T) {
	t.Run("returns what the page has to show", func(t *testing.T) {
		env := newTestLinkEnv(t)
		link, errCreate := env.service.LinkCreateInvite(LinkRequest{Username: "alice", Superuser: true}, "")
		if errCreate != nil {
			t.Fatalf("failed to create a link: %v", errCreate)
		}

		payload, err := env.service.LinkVerify(link.Token)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if payload.Username != "alice" || !payload.Superuser {
			t.Fatalf("expected the link facts, got %+v", payload)
		}
	})

	t.Run("rejects an empty token", func(t *testing.T) {
		env := newTestLinkEnv(t)

		if _, err := env.service.LinkVerify(""); !errors.Is(err, ErrLinkInvalid) {
			t.Fatalf("expected ErrLinkInvalid, got %v", err)
		}
	})

	t.Run("rejects garbage", func(t *testing.T) {
		env := newTestLinkEnv(t)

		if _, err := env.service.LinkVerify("garbage"); !errors.Is(err, ErrLinkInvalid) {
			t.Fatalf("expected ErrLinkInvalid, got %v", err)
		}
	})

	t.Run("reports an expired token", func(t *testing.T) {
		env := newTestLinkEnv(t)
		expired, _, errToken := env.tokenService.Generate("alice", jwt.MapClaims{"jti": "id-1"}, -time.Minute)
		if errToken != nil {
			t.Fatalf("failed to sign: %v", errToken)
		}

		if _, err := env.service.LinkVerify(expired); !errors.Is(err, ErrLinkExpired) {
			t.Fatalf("expected ErrLinkExpired, got %v", err)
		}
	})

	t.Run("reports a revoked link even though the token is still signed", func(t *testing.T) {
		env := newTestLinkEnv(t)
		link, errCreate := env.service.LinkCreateInvite(LinkRequest{Username: "alice"}, "")
		if errCreate != nil {
			t.Fatalf("failed to create a link: %v", errCreate)
		}
		if err := env.service.LinkRevoke(link.Id, ""); err != nil {
			t.Fatalf("failed to revoke: %v", err)
		}

		if _, err := env.service.LinkVerify(link.Token); !errors.Is(err, ErrLinkObsolete) {
			t.Fatalf("expected ErrLinkObsolete, got %v", err)
		}
	})

	t.Run("rejects a token whose subject does not match the record", func(t *testing.T) {
		env := newTestLinkEnv(t)
		link, errCreate := env.service.LinkCreateInvite(LinkRequest{Username: "alice"}, "")
		if errCreate != nil {
			t.Fatalf("failed to create a link: %v", errCreate)
		}
		forged, _, errToken := env.tokenService.Generate("mallory", jwt.MapClaims{"jti": link.Id}, time.Hour)
		if errToken != nil {
			t.Fatalf("failed to sign: %v", errToken)
		}

		if _, err := env.service.LinkVerify(forged); !errors.Is(err, ErrLinkInvalid) {
			t.Fatalf("expected ErrLinkInvalid, got %v", err)
		}
	})

	t.Run("rejects a token without a record id", func(t *testing.T) {
		env := newTestLinkEnv(t)
		anonymous, _, errToken := env.tokenService.Generate("alice", nil, time.Hour)
		if errToken != nil {
			t.Fatalf("failed to sign: %v", errToken)
		}

		if _, err := env.service.LinkVerify(anonymous); !errors.Is(err, ErrLinkInvalid) {
			t.Fatalf("expected ErrLinkInvalid, got %v", err)
		}
	})
}

func TestLinkServiceApply(t *testing.T) {
	t.Run("creates the user and spends the link", func(t *testing.T) {
		env := newTestLinkEnv(t)
		link, errCreate := env.service.LinkCreateInvite(LinkRequest{Username: "alice", Superuser: true}, "")
		if errCreate != nil {
			t.Fatalf("failed to create a link: %v", errCreate)
		}

		user, err := env.service.LinkApply(link.Token, "password123")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if user.Username != "alice" || !user.Superuser {
			t.Fatalf("expected alice to become a superuser, got %+v", user)
		}
		if errVerify := env.service.VerifyPassword("alice", "password123"); errVerify != nil {
			t.Fatalf("expected the new password to verify, got %v", errVerify)
		}
		if _, errGet := env.repository.LinkGet(link.Id); !errors.Is(errGet, storage.ErrNotFound) {
			t.Fatalf("expected the link record to be gone, got %v", errGet)
		}
	})

	t.Run("the same link cannot be used twice", func(t *testing.T) {
		env := newTestLinkEnv(t)
		link, errCreate := env.service.LinkCreateInvite(LinkRequest{Username: "alice"}, "")
		if errCreate != nil {
			t.Fatalf("failed to create a link: %v", errCreate)
		}
		if _, err := env.service.LinkApply(link.Token, "password123"); err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if _, err := env.service.LinkApply(link.Token, "another"); !errors.Is(err, ErrLinkObsolete) {
			t.Fatalf("expected ErrLinkObsolete, got %v", err)
		}
	})

	t.Run("an empty password is rejected and the link survives", func(t *testing.T) {
		env := newTestLinkEnv(t)
		link, errCreate := env.service.LinkCreateInvite(LinkRequest{Username: "alice"}, "")
		if errCreate != nil {
			t.Fatalf("failed to create a link: %v", errCreate)
		}

		if _, err := env.service.LinkApply(link.Token, ""); !errors.Is(err, ErrPasswordCannotBeEmpty) {
			t.Fatalf("expected ErrPasswordCannotBeEmpty, got %v", err)
		}
		if _, errGet := env.repository.LinkGet(link.Id); errGet != nil {
			t.Fatalf("expected the link to survive a failed attempt, got %v", errGet)
		}
	})

	t.Run("an expired link cannot be spent", func(t *testing.T) {
		env := newTestLinkEnv(t)
		expired, _, errToken := env.tokenService.Generate("alice", jwt.MapClaims{"jti": "id-1"}, -time.Minute)
		if errToken != nil {
			t.Fatalf("failed to sign: %v", errToken)
		}

		if _, err := env.service.LinkApply(expired, "password123"); !errors.Is(err, ErrLinkExpired) {
			t.Fatalf("expected ErrLinkExpired, got %v", err)
		}
	})
}

func TestLinkServiceList(t *testing.T) {
	t.Run("returns the outstanding links, newest first, without any token", func(t *testing.T) {
		env := newTestLinkEnv(t)
		if _, err := env.service.LinkCreateInvite(LinkRequest{Username: "alice"}, ""); err != nil {
			t.Fatalf("failed to create a link: %v", err)
		}
		time.Sleep(2 * time.Millisecond)
		if _, err := env.service.LinkCreateInvite(LinkRequest{Username: "bob"}, ""); err != nil {
			t.Fatalf("failed to create a link: %v", err)
		}

		links, err := env.service.LinkList()
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(links) != 2 {
			t.Fatalf("expected 2 links, got %d", len(links))
		}
		if links[0].Username != "bob" || links[1].Username != "alice" {
			t.Fatalf("expected the newest link first, got %+v", links)
		}
	})

	t.Run("drops the links that have expired", func(t *testing.T) {
		env := newTestLinkEnv(t)
		expired := Link{Username: "old", CreatedAt: time.Now().Add(-2 * time.Hour), ExpiresAt: time.Now().Add(-time.Hour)}
		if _, err := env.repository.LinkCreate("id-old", expired); err != nil {
			t.Fatalf("failed to seed: %v", err)
		}
		if _, err := env.service.LinkCreateInvite(LinkRequest{Username: "alice"}, ""); err != nil {
			t.Fatalf("failed to create a link: %v", err)
		}

		links, err := env.service.LinkList()
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(links) != 1 || links[0].Username != "alice" {
			t.Fatalf("expected only the live link, got %+v", links)
		}
		if _, errGet := env.repository.LinkGet("id-old"); !errors.Is(errGet, storage.ErrNotFound) {
			t.Fatalf("expected the expired link to be pruned, got %v", errGet)
		}
	})
}

func TestLinkServiceRevoke(t *testing.T) {
	t.Run("rejects an empty id", func(t *testing.T) {
		env := newTestLinkEnv(t)

		if err := env.service.LinkRevoke("", ""); !errors.Is(err, ErrLinkIdCannotBeEmpty) {
			t.Fatalf("expected ErrLinkIdCannotBeEmpty, got %v", err)
		}
	})

	t.Run("reports a link that is already gone", func(t *testing.T) {
		env := newTestLinkEnv(t)

		if err := env.service.LinkRevoke("id-nothing", ""); !errors.Is(err, ErrLinkObsolete) {
			t.Fatalf("expected ErrLinkObsolete, got %v", err)
		}
	})
}

func TestLinkServiceDeleteAll(t *testing.T) {
	env := newTestLinkEnv(t)
	if _, err := env.service.LinkCreateInvite(LinkRequest{Username: "alice"}, ""); err != nil {
		t.Fatalf("failed to create a link: %v", err)
	}

	if err := env.service.LinkDeleteAll(); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	links, err := env.service.LinkList()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(links) != 0 {
		t.Fatalf("expected no links, got %d", len(links))
	}
}
