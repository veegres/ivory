package user

import (
	"errors"
	"ivory/clients/storage"
	"ivory/core/service/encryption"
	"ivory/core/service/secret"
	"ivory/core/service/token"
	"ivory/features/permission"
	"path/filepath"
	"testing"

	"github.com/boltdb/bolt"
)

type testUserEnv struct {
	service           *Service
	repository        *Repository
	secretService     *secret.Service
	permissionService *permission.Service
	tokenService      *token.Service
	db                *bolt.DB
}

func newTestUserEnv(t *testing.T) *testUserEnv {
	t.Helper()

	db, err := bolt.Open(filepath.Join(t.TempDir(), "test.db"), 0600, nil)
	if err != nil {
		t.Fatalf("failed to open test database: %v", err)
	}
	t.Cleanup(func() { db.Close() })

	secretService := secret.NewService(
		secret.NewRepository(storage.NewDbBucket[string](db, "Secret")),
		encryption.NewService(),
	)
	if errSecret := secretService.SetDefault(); errSecret != nil {
		t.Fatalf("failed to set default secret: %v", errSecret)
	}

	permissionService := permission.NewService(
		permission.NewRepository(storage.NewDbBucket[permission.PermissionMap](db, "Permission")),
	)
	repository := NewRepository(
		storage.NewDbBucket[User](db, "User"),
		storage.NewDbBucket[Link](db, "UserLink"),
	)
	tokenService := token.NewService(secretService)

	return &testUserEnv{
		service:           NewService(repository, encryption.NewService(), secretService, permissionService, tokenService),
		repository:        repository,
		secretService:     secretService,
		permissionService: permissionService,
		tokenService:      tokenService,
		db:                db,
	}
}

func TestServiceCreate(t *testing.T) {
	t.Run("creates a user and stores the password encrypted", func(t *testing.T) {
		env := newTestUserEnv(t)

		response, err := env.service.Create("alice", "password123", false)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if response.Username != "alice" || response.Superuser {
			t.Fatalf("expected a regular user 'alice', got %+v", response)
		}

		stored, errGet := env.repository.Get("alice")
		if errGet != nil {
			t.Fatalf("expected no error, got %v", errGet)
		}
		if stored.Password == "password123" {
			t.Fatalf("expected the password to be encrypted at rest")
		}
	})

	t.Run("trims the username", func(t *testing.T) {
		env := newTestUserEnv(t)

		response, err := env.service.Create("  alice  ", "password123", false)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if response.Username != "alice" {
			t.Fatalf("expected username 'alice', got %q", response.Username)
		}
	})

	t.Run("rejects an empty username", func(t *testing.T) {
		env := newTestUserEnv(t)

		if _, err := env.service.Create("   ", "password123", false); !errors.Is(err, ErrUsernameCannotBeEmpty) {
			t.Fatalf("expected ErrUsernameCannotBeEmpty, got %v", err)
		}
	})

	t.Run("rejects an empty password", func(t *testing.T) {
		env := newTestUserEnv(t)

		if _, err := env.service.Create("alice", "", false); !errors.Is(err, ErrPasswordCannotBeEmpty) {
			t.Fatalf("expected ErrPasswordCannotBeEmpty, got %v", err)
		}
	})

	t.Run("rejects a username that is taken", func(t *testing.T) {
		env := newTestUserEnv(t)

		if _, err := env.service.Create("alice", "password123", false); err != nil {
			t.Fatalf("failed to seed alice: %v", err)
		}
		if _, err := env.service.Create("alice", "another", false); !errors.Is(err, ErrUserAlreadyExists) {
			t.Fatalf("expected ErrUserAlreadyExists, got %v", err)
		}
	})

	t.Run("a new superuser is granted every permission", func(t *testing.T) {
		env := newTestUserEnv(t)

		if _, err := env.service.Create("root", "password123", true); err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		permissions, errPerm := env.permissionService.CreateUserPermissions("basic", "root")
		if errPerm != nil {
			t.Fatalf("expected no error, got %v", errPerm)
		}
		for feature, status := range permissions {
			if status != permission.GRANTED {
				t.Fatalf("expected %s to be granted to a superuser, got %v", feature, status)
			}
		}
	})
}

func TestServiceList(t *testing.T) {
	env := newTestUserEnv(t)

	if _, err := env.service.Create("bob", "password123", false); err != nil {
		t.Fatalf("failed to seed bob: %v", err)
	}
	if _, err := env.service.Create("alice", "password123", true); err != nil {
		t.Fatalf("failed to seed alice: %v", err)
	}

	users, err := env.service.List()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(users) != 2 {
		t.Fatalf("expected 2 users, got %d", len(users))
	}
	if users[0].Username != "alice" || !users[0].Superuser {
		t.Fatalf("expected alice first and a superuser, got %+v", users[0])
	}
	if users[1].Username != "bob" || users[1].Superuser {
		t.Fatalf("expected bob second and a regular user, got %+v", users[1])
	}
}

func TestServiceVerifyPassword(t *testing.T) {
	env := newTestUserEnv(t)
	if _, err := env.service.Create("alice", "password123", false); err != nil {
		t.Fatalf("failed to seed alice: %v", err)
	}

	t.Run("accepts the right password", func(t *testing.T) {
		if err := env.service.VerifyPassword("alice", "password123"); err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
	})

	t.Run("rejects the wrong password", func(t *testing.T) {
		if err := env.service.VerifyPassword("alice", "nope"); !errors.Is(err, ErrCredentialsNotCorrect) {
			t.Fatalf("expected ErrCredentialsNotCorrect, got %v", err)
		}
	})

	t.Run("rejects an unknown user with the same error", func(t *testing.T) {
		if err := env.service.VerifyPassword("mallory", "password123"); !errors.Is(err, ErrCredentialsNotCorrect) {
			t.Fatalf("expected ErrCredentialsNotCorrect, got %v", err)
		}
	})
}

func TestServiceDelete(t *testing.T) {
	t.Run("deletes a regular user", func(t *testing.T) {
		env := newTestUserEnv(t)
		if _, err := env.service.Create("alice", "password123", false); err != nil {
			t.Fatalf("failed to seed alice: %v", err)
		}

		if err := env.service.Delete("alice", "root"); err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if _, err := env.repository.Get("alice"); !errors.Is(err, storage.ErrNotFound) {
			t.Fatalf("expected the user to be gone, got %v", err)
		}
	})

	t.Run("rejects an empty username", func(t *testing.T) {
		env := newTestUserEnv(t)

		if err := env.service.Delete("", "root"); !errors.Is(err, ErrUsernameCannotBeEmpty) {
			t.Fatalf("expected ErrUsernameCannotBeEmpty, got %v", err)
		}
	})

	t.Run("reports an unknown user", func(t *testing.T) {
		env := newTestUserEnv(t)

		if err := env.service.Delete("nobody", "root"); !errors.Is(err, ErrUserNotFound) {
			t.Fatalf("expected ErrUserNotFound, got %v", err)
		}
	})

	// NOTE: with authentication on, the last superuser is unreachable anyway -
	// only a superuser may delete one and nobody deletes themselves. The floor
	// is what holds when Ivory runs without authentication, where the requester
	// is nobody in particular.
	t.Run("refuses to delete the last superuser", func(t *testing.T) {
		env := newTestUserEnv(t)
		if _, err := env.service.Create("root", "password123", true); err != nil {
			t.Fatalf("failed to seed root: %v", err)
		}

		if err := env.service.Delete("root", ""); !errors.Is(err, ErrLastSuperuser) {
			t.Fatalf("expected ErrLastSuperuser, got %v", err)
		}
	})

	t.Run("deletes a superuser while another one is left", func(t *testing.T) {
		env := newTestUserEnv(t)
		if _, err := env.service.Create("root", "password123", true); err != nil {
			t.Fatalf("failed to seed root: %v", err)
		}
		if _, err := env.service.Create("admin", "password123", true); err != nil {
			t.Fatalf("failed to seed admin: %v", err)
		}

		if err := env.service.Delete("root", "admin"); err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
	})

	t.Run("nobody deletes themselves", func(t *testing.T) {
		env := newTestUserEnv(t)
		if _, err := env.service.Create("root", "password123", true); err != nil {
			t.Fatalf("failed to seed root: %v", err)
		}
		if _, err := env.service.Create("admin", "password123", true); err != nil {
			t.Fatalf("failed to seed admin: %v", err)
		}

		if err := env.service.Delete("root", "root"); !errors.Is(err, ErrCannotDeleteYourself) {
			t.Fatalf("expected ErrCannotDeleteYourself, got %v", err)
		}
	})

	t.Run("a regular user cannot delete a superuser", func(t *testing.T) {
		env := newTestUserEnv(t)
		if _, err := env.service.Create("root", "password123", true); err != nil {
			t.Fatalf("failed to seed root: %v", err)
		}
		if _, err := env.service.Create("admin", "password123", true); err != nil {
			t.Fatalf("failed to seed admin: %v", err)
		}
		if _, err := env.service.Create("alice", "password123", false); err != nil {
			t.Fatalf("failed to seed alice: %v", err)
		}

		if err := env.service.Delete("root", "alice"); !errors.Is(err, ErrSuperuserRequired) {
			t.Fatalf("expected ErrSuperuserRequired, got %v", err)
		}
	})

	t.Run("a regular user can delete another regular user", func(t *testing.T) {
		env := newTestUserEnv(t)
		if _, err := env.service.Create("alice", "password123", false); err != nil {
			t.Fatalf("failed to seed alice: %v", err)
		}
		if _, err := env.service.Create("bob", "password123", false); err != nil {
			t.Fatalf("failed to seed bob: %v", err)
		}

		if err := env.service.Delete("bob", "alice"); err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
	})

	t.Run("deleting a user takes their permissions with them", func(t *testing.T) {
		env := newTestUserEnv(t)
		if _, err := env.service.Create("root", "password123", true); err != nil {
			t.Fatalf("failed to seed root: %v", err)
		}
		if _, err := env.service.Create("alice", "password123", false); err != nil {
			t.Fatalf("failed to seed alice: %v", err)
		}
		if _, err := env.permissionService.CreateUserPermissions("basic", "alice"); err != nil {
			t.Fatalf("failed to seed permissions: %v", err)
		}

		if err := env.service.Delete("alice", "root"); err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		all, err := env.permissionService.GetAllUserPermissions()
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		for _, up := range all {
			if up.Username == "basic:alice" {
				t.Fatalf("expected the permissions of a deleted user to be gone, got %+v", up)
			}
		}
	})
}

func TestServiceUpdatePassword(t *testing.T) {
	t.Run("replaces the password when the previous one is right", func(t *testing.T) {
		env := newTestUserEnv(t)
		if _, err := env.service.Create("alice", "password123", false); err != nil {
			t.Fatalf("failed to seed alice: %v", err)
		}

		if err := env.service.UpdatePassword("alice", "password123", "newpassword"); err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if err := env.service.VerifyPassword("alice", "newpassword"); err != nil {
			t.Fatalf("expected the new password to verify, got %v", err)
		}
		if err := env.service.VerifyPassword("alice", "password123"); !errors.Is(err, ErrCredentialsNotCorrect) {
			t.Fatalf("expected the old password to stop working, got %v", err)
		}
	})

	t.Run("rejects a wrong previous password", func(t *testing.T) {
		env := newTestUserEnv(t)
		if _, err := env.service.Create("alice", "password123", false); err != nil {
			t.Fatalf("failed to seed alice: %v", err)
		}

		if err := env.service.UpdatePassword("alice", "nope", "newpassword"); !errors.Is(err, ErrCredentialsNotCorrect) {
			t.Fatalf("expected ErrCredentialsNotCorrect, got %v", err)
		}
		if err := env.service.VerifyPassword("alice", "password123"); err != nil {
			t.Fatalf("expected the old password to still work, got %v", err)
		}
	})

	t.Run("rejects an empty new password", func(t *testing.T) {
		env := newTestUserEnv(t)
		if _, err := env.service.Create("alice", "password123", false); err != nil {
			t.Fatalf("failed to seed alice: %v", err)
		}

		if err := env.service.UpdatePassword("alice", "password123", ""); !errors.Is(err, ErrPasswordCannotBeEmpty) {
			t.Fatalf("expected ErrPasswordCannotBeEmpty, got %v", err)
		}
	})

	t.Run("rejects an unknown user", func(t *testing.T) {
		env := newTestUserEnv(t)

		if err := env.service.UpdatePassword("mallory", "password123", "newpassword"); !errors.Is(err, ErrCredentialsNotCorrect) {
			t.Fatalf("expected ErrCredentialsNotCorrect, got %v", err)
		}
	})
}

func TestServiceIsSuperuser(t *testing.T) {
	env := newTestUserEnv(t)
	if _, err := env.service.Create("root", "password123", true); err != nil {
		t.Fatalf("failed to seed root: %v", err)
	}
	if _, err := env.service.Create("alice", "password123", false); err != nil {
		t.Fatalf("failed to seed alice: %v", err)
	}

	cases := []struct {
		username string
		expected bool
	}{
		{"root", true},
		{"alice", false},
		{"nobody", false},
	}
	for _, c := range cases {
		is, err := env.service.IsSuperuser(c.username)
		if err != nil {
			t.Fatalf("expected no error for %s, got %v", c.username, err)
		}
		if is != c.expected {
			t.Fatalf("expected %v for %s, got %v", c.expected, c.username, is)
		}
	}
}

func TestServiceHasSuperuser(t *testing.T) {
	env := newTestUserEnv(t)

	has, err := env.service.HasSuperuser()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if has {
		t.Fatalf("expected no superuser in a fresh Ivory")
	}

	if _, errCreate := env.service.Create("alice", "password123", false); errCreate != nil {
		t.Fatalf("failed to seed alice: %v", errCreate)
	}
	has, err = env.service.HasSuperuser()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if has {
		t.Fatalf("expected a regular user not to count as a superuser")
	}

	if _, errCreate := env.service.Create("root", "password123", true); errCreate != nil {
		t.Fatalf("failed to seed root: %v", errCreate)
	}
	has, err = env.service.HasSuperuser()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !has {
		t.Fatalf("expected a superuser")
	}
}

func TestServiceExists(t *testing.T) {
	env := newTestUserEnv(t)
	if _, err := env.service.Create("alice", "password123", false); err != nil {
		t.Fatalf("failed to seed alice: %v", err)
	}

	exist, err := env.service.Exists("alice")
	if err != nil || !exist {
		t.Fatalf("expected alice to exist, got %v (%v)", exist, err)
	}
	exist, err = env.service.Exists("mallory")
	if err != nil || exist {
		t.Fatalf("expected mallory not to exist, got %v (%v)", exist, err)
	}
}

func TestServiceDeleteAll(t *testing.T) {
	env := newTestUserEnv(t)
	if _, err := env.service.Create("root", "password123", true); err != nil {
		t.Fatalf("failed to seed root: %v", err)
	}

	if err := env.service.DeleteAll(); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	users, err := env.service.List()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(users) != 0 {
		t.Fatalf("expected no users, got %d", len(users))
	}
}

func TestServiceReencrypt(t *testing.T) {
	env := newTestUserEnv(t)
	if _, err := env.service.Create("alice", "password123", false); err != nil {
		t.Fatalf("failed to seed alice: %v", err)
	}

	previous, next, errUpdate := env.secretService.Update("ivory", "newsecret")
	if errUpdate != nil {
		t.Fatalf("failed to rotate the secret: %v", errUpdate)
	}
	if err := env.service.Reencrypt(previous, next); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if err := env.service.VerifyPassword("alice", "password123"); err != nil {
		t.Fatalf("expected the password to still verify after the rotation, got %v", err)
	}
}

func TestServiceLoadsSuperusersOnStart(t *testing.T) {
	env := newTestUserEnv(t)
	if _, err := env.service.Create("root", "password123", true); err != nil {
		t.Fatalf("failed to seed root: %v", err)
	}

	// a restart is a fresh service reading the same store
	fresh := permission.NewService(
		permission.NewRepository(storage.NewDbBucket[permission.PermissionMap](env.db, "PermissionFresh")),
	)
	NewService(env.repository, encryption.NewService(), env.secretService, fresh, env.tokenService)

	permissions, err := fresh.CreateUserPermissions("basic", "root")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	for feature, status := range permissions {
		if status != permission.GRANTED {
			t.Fatalf("expected %s to be granted to a superuser after a restart, got %v", feature, status)
		}
	}
}
