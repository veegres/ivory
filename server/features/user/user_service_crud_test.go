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
	repository := NewRepository(storage.NewDbBucket[User](db, "User"))
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

// seedUser registers somebody who can already sign in, which is what most of
// these cases are about - the registration itself is tested where it lives.
func seedUser(t *testing.T, env *testUserEnv, username string, password string, superuser bool) {
	t.Helper()
	if _, err := env.service.CreateOutright(username, password, []AuthType{AuthBasic}, superuser); err != nil {
		t.Fatalf("failed to seed %s: %v", username, err)
	}
}

func TestServiceCreate(t *testing.T) {
	t.Run("registers a user and hands out a registration when they sign in with a password", func(t *testing.T) {
		env := newTestUserEnv(t)

		response, err := env.service.Create(CreateRequest{Username: "alice", AuthTypes: []AuthType{AuthBasic}}, "")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if response.User.Username != "alice" || response.User.Superuser {
			t.Fatalf("expected a regular user 'alice', got %+v", response.User)
		}
		if response.Registration == nil || response.Registration.Token == "" {
			t.Fatalf("expected a registration to hand out, got %+v", response.Registration)
		}
		if response.User.Registration == nil || response.User.Registration.Status != RegistrationPending {
			t.Fatalf("expected the user to be waiting on their registration, got %+v", response.User.Registration)
		}

		stored, errGet := env.repository.Get("alice")
		if errGet != nil {
			t.Fatalf("expected no error, got %v", errGet)
		}
		if stored.Password != "" {
			t.Fatalf("expected no password until its owner sets one, got %q", stored.Password)
		}
	})

	t.Run("a user who signs in elsewhere gets no registration at all", func(t *testing.T) {
		env := newTestUserEnv(t)

		response, err := env.service.Create(CreateRequest{Username: "bob", AuthTypes: []AuthType{AuthLdap}}, "")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if response.Registration != nil {
			t.Fatalf("expected no registration for an ldap-only user, got %+v", response.Registration)
		}
		if response.User.Registration != nil {
			t.Fatalf("expected no registration state for an ldap-only user, got %+v", response.User.Registration)
		}
	})

	t.Run("trims the username", func(t *testing.T) {
		env := newTestUserEnv(t)

		response, err := env.service.Create(CreateRequest{Username: "  alice  ", AuthTypes: []AuthType{AuthLdap}}, "")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if response.User.Username != "alice" {
			t.Fatalf("expected username 'alice', got %q", response.User.Username)
		}
	})

	t.Run("rejects an empty username", func(t *testing.T) {
		env := newTestUserEnv(t)

		_, err := env.service.Create(CreateRequest{Username: "   ", AuthTypes: []AuthType{AuthLdap}}, "")
		if !errors.Is(err, ErrUsernameCannotBeEmpty) {
			t.Fatalf("expected ErrUsernameCannotBeEmpty, got %v", err)
		}
	})

	t.Run("rejects a user with no way of signing in", func(t *testing.T) {
		env := newTestUserEnv(t)

		_, err := env.service.Create(CreateRequest{Username: "alice", AuthTypes: []AuthType{}}, "")
		if !errors.Is(err, ErrAuthTypeRequired) {
			t.Fatalf("expected ErrAuthTypeRequired, got %v", err)
		}
	})

	t.Run("rejects a way of signing in that does not exist", func(t *testing.T) {
		env := newTestUserEnv(t)

		_, err := env.service.Create(CreateRequest{Username: "alice", AuthTypes: []AuthType{"smoke-signals"}}, "")
		if !errors.Is(err, ErrAuthTypeInvalid) {
			t.Fatalf("expected ErrAuthTypeInvalid, got %v", err)
		}
	})

	t.Run("rejects a username that is taken", func(t *testing.T) {
		env := newTestUserEnv(t)
		seedUser(t, env, "alice", "password123", false)

		_, err := env.service.Create(CreateRequest{Username: "alice", AuthTypes: []AuthType{AuthLdap}}, "")
		if !errors.Is(err, ErrUserAlreadyExists) {
			t.Fatalf("expected ErrUserAlreadyExists, got %v", err)
		}
	})

	t.Run("only a superuser registers a superuser", func(t *testing.T) {
		env := newTestUserEnv(t)
		seedUser(t, env, "root", "password123", true)
		seedUser(t, env, "alice", "password123", false)

		_, err := env.service.Create(CreateRequest{Username: "newroot", AuthTypes: []AuthType{AuthLdap}, Superuser: true}, "alice")
		if !errors.Is(err, ErrSuperuserRequired) {
			t.Fatalf("expected ErrSuperuserRequired, got %v", err)
		}
		if _, errSuper := env.service.Create(CreateRequest{Username: "newroot", AuthTypes: []AuthType{AuthLdap}, Superuser: true}, "root"); errSuper != nil {
			t.Fatalf("expected a superuser to be allowed, got %v", errSuper)
		}
	})

	t.Run("a new user starts with nothing permitted", func(t *testing.T) {
		env := newTestUserEnv(t)

		if _, err := env.service.Create(CreateRequest{Username: "alice", AuthTypes: []AuthType{AuthLdap}}, ""); err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		permissions, errPerm := env.permissionService.GetUserPermissions("alice", false, nil)
		if errPerm != nil {
			t.Fatalf("expected the record to exist, got %v", errPerm)
		}
		for feature, status := range permissions {
			if status != permission.NOT_PERMITTED {
				t.Fatalf("expected %s to start not permitted, got %v", feature, status)
			}
		}
	})

	t.Run("a new superuser is granted every permission", func(t *testing.T) {
		env := newTestUserEnv(t)

		if _, err := env.service.Create(CreateRequest{Username: "root", AuthTypes: []AuthType{AuthLdap}, Superuser: true}, ""); err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		permissions, errPerm := env.permissionService.GetUserPermissions("root", false, nil)
		if errPerm != nil {
			t.Fatalf("expected the record to exist, got %v", errPerm)
		}
		for feature, status := range permissions {
			if status != permission.GRANTED {
				t.Fatalf("expected %s to be granted to a superuser, got %v", feature, status)
			}
		}
	})
}

func TestServiceCreateOutright(t *testing.T) {
	t.Run("stores the typed password encrypted and hands out no registration", func(t *testing.T) {
		env := newTestUserEnv(t)

		response, err := env.service.CreateOutright("root", "password123", []AuthType{AuthBasic}, true)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if !response.Superuser {
			t.Fatalf("expected a superuser, got %+v", response)
		}
		if response.Registration == nil || response.Registration.Status != RegistrationActive {
			t.Fatalf("expected an account that can sign in already, got %+v", response.Registration)
		}

		stored, errGet := env.repository.Get("root")
		if errGet != nil {
			t.Fatalf("expected no error, got %v", errGet)
		}
		if stored.Password == "password123" {
			t.Fatalf("expected the password to be encrypted at rest")
		}
		if stored.RegistrationId != "" {
			t.Fatalf("expected no registration outstanding, got %q", stored.RegistrationId)
		}
	})

	// NOTE: this is what a restore looks like - a backup carries no password, so
	// the user comes back waiting for a registration somebody has to issue
	t.Run("no password leaves the user waiting for a registration", func(t *testing.T) {
		env := newTestUserEnv(t)

		response, err := env.service.CreateOutright("alice", "", []AuthType{AuthBasic}, false)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if response.Registration == nil || response.Registration.Status != RegistrationMissing {
			t.Fatalf("expected a registration still to issue, got %+v", response.Registration)
		}
	})
}

func TestServiceUpdate(t *testing.T) {
	t.Run("changes the ways a user may sign in", func(t *testing.T) {
		env := newTestUserEnv(t)
		if _, err := env.service.Create(CreateRequest{Username: "alice", AuthTypes: []AuthType{AuthLdap}}, ""); err != nil {
			t.Fatalf("failed to seed alice: %v", err)
		}

		response, err := env.service.Update("alice", []AuthType{AuthLdap, AuthOidc}, "")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(response.User.AuthTypes) != 2 {
			t.Fatalf("expected both ways of signing in, got %v", response.User.AuthTypes)
		}
		if response.Registration != nil {
			t.Fatalf("expected no registration where no password is involved, got %+v", response.Registration)
		}
	})

	t.Run("asking for a password waits for an explicit reset link", func(t *testing.T) {
		env := newTestUserEnv(t)
		if _, err := env.service.Create(CreateRequest{Username: "alice", AuthTypes: []AuthType{AuthLdap}}, ""); err != nil {
			t.Fatalf("failed to seed alice: %v", err)
		}

		response, err := env.service.Update("alice", []AuthType{AuthLdap, AuthBasic}, "")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if response.Registration != nil {
			t.Fatalf("expected no registration to be issued by update, got %+v", response.Registration)
		}
		if response.User.Registration == nil || response.User.Registration.Status != RegistrationMissing {
			t.Fatalf("expected the user to need a reset link, got %+v", response.User.Registration)
		}

		stored, errGet := env.repository.Get("alice")
		if errGet != nil {
			t.Fatalf("expected no error, got %v", errGet)
		}
		if stored.RegistrationId != "" {
			t.Fatalf("expected no registration id until reset, got %q", stored.RegistrationId)
		}
	})

	t.Run("taking the password away drops it along with any registration", func(t *testing.T) {
		env := newTestUserEnv(t)
		seedUser(t, env, "alice", "password123", false)

		if _, err := env.service.Update("alice", []AuthType{AuthLdap}, ""); err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		stored, errGet := env.repository.Get("alice")
		if errGet != nil {
			t.Fatalf("expected no error, got %v", errGet)
		}
		if stored.Password != "" || stored.RegistrationId != "" {
			t.Fatalf("expected the password and registration to be gone, got %+v", stored)
		}
		if err := env.service.VerifyPassword("alice", "password123"); !errors.Is(err, ErrCredentialsNotCorrect) {
			t.Fatalf("expected the password to stop working, got %v", err)
		}
	})

	t.Run("only a superuser changes a superuser", func(t *testing.T) {
		env := newTestUserEnv(t)
		seedUser(t, env, "root", "password123", true)
		seedUser(t, env, "alice", "password123", false)

		if _, err := env.service.Update("root", []AuthType{AuthLdap}, "alice"); !errors.Is(err, ErrSuperuserRequired) {
			t.Fatalf("expected ErrSuperuserRequired, got %v", err)
		}
	})

	t.Run("rejects an unknown user", func(t *testing.T) {
		env := newTestUserEnv(t)

		if _, err := env.service.Update("nobody", []AuthType{AuthLdap}, ""); !errors.Is(err, ErrUserNotFound) {
			t.Fatalf("expected ErrUserNotFound, got %v", err)
		}
	})

	t.Run("rejects leaving a user with no way of signing in", func(t *testing.T) {
		env := newTestUserEnv(t)
		seedUser(t, env, "alice", "password123", false)

		if _, err := env.service.Update("alice", nil, ""); !errors.Is(err, ErrAuthTypeRequired) {
			t.Fatalf("expected ErrAuthTypeRequired, got %v", err)
		}
	})
}

func TestServiceList(t *testing.T) {
	env := newTestUserEnv(t)
	seedUser(t, env, "bob", "password123", false)
	seedUser(t, env, "alice", "password123", true)

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
	seedUser(t, env, "alice", "password123", false)

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

	t.Run("a user who has not set a password yet cannot sign in with one", func(t *testing.T) {
		if _, err := env.service.Create(CreateRequest{Username: "carol", AuthTypes: []AuthType{AuthBasic}}, ""); err != nil {
			t.Fatalf("failed to seed carol: %v", err)
		}
		if err := env.service.VerifyPassword("carol", ""); !errors.Is(err, ErrCredentialsNotCorrect) {
			t.Fatalf("expected ErrCredentialsNotCorrect, got %v", err)
		}
	})

	t.Run("a user registered for another way of signing in cannot use a password", func(t *testing.T) {
		if _, err := env.service.Create(CreateRequest{Username: "dave", AuthTypes: []AuthType{AuthLdap}}, ""); err != nil {
			t.Fatalf("failed to seed dave: %v", err)
		}
		if err := env.service.VerifyPassword("dave", "password123"); !errors.Is(err, ErrCredentialsNotCorrect) {
			t.Fatalf("expected ErrCredentialsNotCorrect, got %v", err)
		}
	})
}

func TestServiceVerifySignIn(t *testing.T) {
	env := newTestUserEnv(t)
	if _, err := env.service.Create(CreateRequest{Username: "alice", AuthTypes: []AuthType{AuthLdap}}, ""); err != nil {
		t.Fatalf("failed to seed alice: %v", err)
	}

	t.Run("a registered way of signing in is allowed", func(t *testing.T) {
		if err := env.service.VerifySignIn("alice", AuthLdap); err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
	})

	t.Run("another way of signing in is refused", func(t *testing.T) {
		if err := env.service.VerifySignIn("alice", AuthOidc); !errors.Is(err, ErrAuthTypeNotAllowed) {
			t.Fatalf("expected ErrAuthTypeNotAllowed, got %v", err)
		}
	})

	// NOTE: a directory will happily confirm somebody Ivory was never told
	// about, which is the whole reason this gate exists
	t.Run("a name Ivory does not hold is refused", func(t *testing.T) {
		if err := env.service.VerifySignIn("stranger", AuthLdap); !errors.Is(err, ErrUserNotRegistered) {
			t.Fatalf("expected ErrUserNotRegistered, got %v", err)
		}
	})

	t.Run("an empty username is refused", func(t *testing.T) {
		if err := env.service.VerifySignIn("", AuthLdap); !errors.Is(err, ErrUsernameCannotBeEmpty) {
			t.Fatalf("expected ErrUsernameCannotBeEmpty, got %v", err)
		}
	})
}

func TestServiceDelete(t *testing.T) {
	t.Run("deletes a regular user", func(t *testing.T) {
		env := newTestUserEnv(t)
		seedUser(t, env, "alice", "password123", false)

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
		seedUser(t, env, "root", "password123", true)

		if err := env.service.Delete("root", ""); !errors.Is(err, ErrLastSuperuser) {
			t.Fatalf("expected ErrLastSuperuser, got %v", err)
		}
	})

	t.Run("deletes a superuser while another one is left", func(t *testing.T) {
		env := newTestUserEnv(t)
		seedUser(t, env, "root", "password123", true)
		seedUser(t, env, "admin", "password123", true)

		if err := env.service.Delete("root", "admin"); err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
	})

	t.Run("nobody deletes themselves", func(t *testing.T) {
		env := newTestUserEnv(t)
		seedUser(t, env, "root", "password123", true)
		seedUser(t, env, "admin", "password123", true)

		if err := env.service.Delete("root", "root"); !errors.Is(err, ErrCannotDeleteYourself) {
			t.Fatalf("expected ErrCannotDeleteYourself, got %v", err)
		}
	})

	t.Run("a regular user cannot delete a superuser", func(t *testing.T) {
		env := newTestUserEnv(t)
		seedUser(t, env, "root", "password123", true)
		seedUser(t, env, "admin", "password123", true)
		seedUser(t, env, "alice", "password123", false)

		if err := env.service.Delete("root", "alice"); !errors.Is(err, ErrSuperuserRequired) {
			t.Fatalf("expected ErrSuperuserRequired, got %v", err)
		}
	})

	t.Run("a regular user can delete another regular user", func(t *testing.T) {
		env := newTestUserEnv(t)
		seedUser(t, env, "alice", "password123", false)
		seedUser(t, env, "bob", "password123", false)

		if err := env.service.Delete("bob", "alice"); err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
	})

	t.Run("deleting a user takes their permissions with them", func(t *testing.T) {
		env := newTestUserEnv(t)
		seedUser(t, env, "root", "password123", true)
		seedUser(t, env, "alice", "password123", false)

		if err := env.service.Delete("alice", "root"); err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		all, err := env.permissionService.GetAllUserPermissions()
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		for _, up := range all {
			if up.Username == "alice" {
				t.Fatalf("expected the permissions of a deleted user to be gone, got %+v", up)
			}
		}
	})
}

func TestServiceUpdatePassword(t *testing.T) {
	t.Run("replaces the password when the previous one is right", func(t *testing.T) {
		env := newTestUserEnv(t)
		seedUser(t, env, "alice", "password123", false)

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
		seedUser(t, env, "alice", "password123", false)

		if err := env.service.UpdatePassword("alice", "nope", "newpassword"); !errors.Is(err, ErrCredentialsNotCorrect) {
			t.Fatalf("expected ErrCredentialsNotCorrect, got %v", err)
		}
		if err := env.service.VerifyPassword("alice", "password123"); err != nil {
			t.Fatalf("expected the old password to still work, got %v", err)
		}
	})

	t.Run("rejects an empty new password", func(t *testing.T) {
		env := newTestUserEnv(t)
		seedUser(t, env, "alice", "password123", false)

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
	seedUser(t, env, "root", "password123", true)
	seedUser(t, env, "alice", "password123", false)

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

	seedUser(t, env, "alice", "password123", false)
	has, err = env.service.HasSuperuser()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if has {
		t.Fatalf("expected a regular user not to count as a superuser")
	}

	seedUser(t, env, "root", "password123", true)
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
	seedUser(t, env, "alice", "password123", false)

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
	seedUser(t, env, "root", "password123", true)

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
	seedUser(t, env, "alice", "password123", false)
	// NOTE: a user with no password yet has nothing to re-encrypt, and used to
	// take the whole rotation down with an empty string
	if _, err := env.service.Create(CreateRequest{Username: "carol", AuthTypes: []AuthType{AuthBasic}}, ""); err != nil {
		t.Fatalf("failed to seed carol: %v", err)
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

// TestServiceSyncsPermissionsOnStart covers the startup pass: a record that is
// missing is created, and a feature added since the last start is filled in
// with the status that user's kind gets - which is how a superuser comes to
// hold a permission that did not exist when they were registered.
func TestServiceSyncsPermissionsOnStart(t *testing.T) {
	env := newTestUserEnv(t)
	seedUser(t, env, "root", "password123", true)
	seedUser(t, env, "alice", "password123", false)

	// a restart is a fresh service reading the same store
	fresh := permission.NewService(
		permission.NewRepository(storage.NewDbBucket[permission.PermissionMap](env.db, "PermissionFresh")),
	)
	NewService(env.repository, encryption.NewService(), env.secretService, fresh, env.tokenService)

	superuserPermissions, errSuper := fresh.GetUserPermissions("root", false, nil)
	if errSuper != nil {
		t.Fatalf("expected the superuser record to exist, got %v", errSuper)
	}
	for feature, status := range superuserPermissions {
		if status != permission.GRANTED {
			t.Fatalf("expected %s to be granted to a superuser after a restart, got %v", feature, status)
		}
	}

	userPermissions, errUser := fresh.GetUserPermissions("alice", false, nil)
	if errUser != nil {
		t.Fatalf("expected the user record to exist, got %v", errUser)
	}
	for feature, status := range userPermissions {
		if status != permission.NOT_PERMITTED {
			t.Fatalf("expected %s to stay not permitted for a regular user, got %v", feature, status)
		}
	}
}
