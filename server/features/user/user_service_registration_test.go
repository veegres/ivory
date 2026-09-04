package user

import (
	"errors"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func TestServicePasswordResetIssue(t *testing.T) {
	t.Run("hands an existing user a new registration, which is what a reset is", func(t *testing.T) {
		env := newTestUserEnv(t)
		seedUser(t, env, "alice", "password123", false)

		registration, err := env.service.PasswordResetIssue("alice", "")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if registration.Username != "alice" || registration.Token == "" {
			t.Fatalf("expected a registration for alice, got %+v", registration)
		}

		// NOTE: the password they already have keeps working until the new one
		// is set - a reset is an offer, not a lockout
		if errVerify := env.service.VerifyPassword("alice", "password123"); errVerify != nil {
			t.Fatalf("expected the old password to still work, got %v", errVerify)
		}
	})

	t.Run("the second registration makes the first one useless", func(t *testing.T) {
		env := newTestUserEnv(t)
		seedUser(t, env, "alice", "password123", false)

		first, errFirst := env.service.PasswordResetIssue("alice", "")
		if errFirst != nil {
			t.Fatalf("failed to issue the first: %v", errFirst)
		}
		if _, errSecond := env.service.PasswordResetIssue("alice", ""); errSecond != nil {
			t.Fatalf("failed to issue the second: %v", errSecond)
		}

		if _, err := env.service.RegistrationVerify(first.Token); !errors.Is(err, ErrRegistrationObsolete) {
			t.Fatalf("expected ErrRegistrationObsolete, got %v", err)
		}
	})

	t.Run("a user who does not sign in with a password has nothing to register", func(t *testing.T) {
		env := newTestUserEnv(t)
		if _, err := env.service.Create(CreateRequest{Username: "bob", AuthTypes: []AuthType{AuthLdap}}, ""); err != nil {
			t.Fatalf("failed to seed bob: %v", err)
		}

		if _, err := env.service.PasswordResetIssue("bob", ""); !errors.Is(err, ErrRegistrationNotNeeded) {
			t.Fatalf("expected ErrRegistrationNotNeeded, got %v", err)
		}
	})

	t.Run("reports an unknown user", func(t *testing.T) {
		env := newTestUserEnv(t)

		if _, err := env.service.PasswordResetIssue("nobody", ""); !errors.Is(err, ErrUserNotFound) {
			t.Fatalf("expected ErrUserNotFound, got %v", err)
		}
	})

	// NOTE: a registration takes the account over, so handing one out for a
	// superuser is a superuser's own right, exactly as deleting one is
	t.Run("only a superuser registers a superuser again", func(t *testing.T) {
		env := newTestUserEnv(t)
		seedUser(t, env, "root", "password123", true)
		seedUser(t, env, "alice", "password123", false)

		if _, err := env.service.PasswordResetIssue("root", "alice"); !errors.Is(err, ErrSuperuserRequired) {
			t.Fatalf("expected ErrSuperuserRequired, got %v", err)
		}
		if _, err := env.service.PasswordResetIssue("root", "root"); err != nil {
			t.Fatalf("expected a superuser to be allowed, got %v", err)
		}
	})
}

func TestServicePasswordResetRevoke(t *testing.T) {
	t.Run("makes an outstanding registration useless straight away", func(t *testing.T) {
		env := newTestUserEnv(t)
		response, errRegister := env.service.Create(CreateRequest{Username: "alice", AuthTypes: []AuthType{AuthBasic}}, "")
		if errRegister != nil {
			t.Fatalf("failed to register alice: %v", errRegister)
		}

		if err := env.service.PasswordResetRevoke("alice", ""); err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if _, err := env.service.RegistrationVerify(response.Registration.Token); !errors.Is(err, ErrRegistrationObsolete) {
			t.Fatalf("expected ErrRegistrationObsolete, got %v", err)
		}
	})

	t.Run("leaves the password the user already has alone", func(t *testing.T) {
		env := newTestUserEnv(t)
		seedUser(t, env, "alice", "password123", false)
		if _, err := env.service.PasswordResetIssue("alice", ""); err != nil {
			t.Fatalf("failed to issue: %v", err)
		}

		if err := env.service.PasswordResetRevoke("alice", ""); err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if err := env.service.VerifyPassword("alice", "password123"); err != nil {
			t.Fatalf("expected the password to still work, got %v", err)
		}
	})

	t.Run("reports that there is nothing outstanding", func(t *testing.T) {
		env := newTestUserEnv(t)
		seedUser(t, env, "alice", "password123", false)

		if err := env.service.PasswordResetRevoke("alice", ""); !errors.Is(err, ErrRegistrationObsolete) {
			t.Fatalf("expected ErrRegistrationObsolete, got %v", err)
		}
	})

	t.Run("reports an unknown user", func(t *testing.T) {
		env := newTestUserEnv(t)

		if err := env.service.PasswordResetRevoke("nobody", ""); !errors.Is(err, ErrUserNotFound) {
			t.Fatalf("expected ErrUserNotFound, got %v", err)
		}
	})

	t.Run("only a superuser revokes a superuser's registration", func(t *testing.T) {
		env := newTestUserEnv(t)
		seedUser(t, env, "root", "password123", true)
		seedUser(t, env, "alice", "password123", false)
		if _, err := env.service.PasswordResetIssue("root", "root"); err != nil {
			t.Fatalf("failed to issue: %v", err)
		}

		if err := env.service.PasswordResetRevoke("root", "alice"); !errors.Is(err, ErrSuperuserRequired) {
			t.Fatalf("expected ErrSuperuserRequired, got %v", err)
		}
	})
}

func TestServiceRegistrationVerify(t *testing.T) {
	t.Run("returns what the page has to show", func(t *testing.T) {
		env := newTestUserEnv(t)
		response, errRegister := env.service.Create(CreateRequest{Username: "alice", AuthTypes: []AuthType{AuthBasic}}, "")
		if errRegister != nil {
			t.Fatalf("failed to register alice: %v", errRegister)
		}

		payload, err := env.service.RegistrationVerify(response.Registration.Token)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if payload.Username != "alice" {
			t.Fatalf("expected alice, got %+v", payload)
		}
		if payload.ExpiresAt.IsZero() {
			t.Fatalf("expected an expiration, got %+v", payload)
		}
	})

	t.Run("rejects an empty token", func(t *testing.T) {
		env := newTestUserEnv(t)

		if _, err := env.service.RegistrationVerify(""); !errors.Is(err, ErrRegistrationInvalid) {
			t.Fatalf("expected ErrRegistrationInvalid, got %v", err)
		}
	})

	t.Run("rejects garbage", func(t *testing.T) {
		env := newTestUserEnv(t)

		if _, err := env.service.RegistrationVerify("garbage"); !errors.Is(err, ErrRegistrationInvalid) {
			t.Fatalf("expected ErrRegistrationInvalid, got %v", err)
		}
	})

	t.Run("reports an expired token", func(t *testing.T) {
		env := newTestUserEnv(t)
		expired, _, errToken := env.tokenService.Generate("alice", jwt.MapClaims{"jti": "id-1"}, -time.Minute)
		if errToken != nil {
			t.Fatalf("failed to sign: %v", errToken)
		}

		if _, err := env.service.RegistrationVerify(expired); !errors.Is(err, ErrRegistrationExpired) {
			t.Fatalf("expected ErrRegistrationExpired, got %v", err)
		}
	})

	// NOTE: the token is half of a registration and the id on the record is the
	// other - a signature Ivory made itself opens nothing on its own
	t.Run("rejects a token whose id the user no longer carries", func(t *testing.T) {
		env := newTestUserEnv(t)
		seedUser(t, env, "alice", "password123", false)
		forged, _, errToken := env.tokenService.Generate("alice", jwt.MapClaims{"jti": "id-nobody-issued"}, time.Hour)
		if errToken != nil {
			t.Fatalf("failed to sign: %v", errToken)
		}

		if _, err := env.service.RegistrationVerify(forged); !errors.Is(err, ErrRegistrationObsolete) {
			t.Fatalf("expected ErrRegistrationObsolete, got %v", err)
		}
	})

	t.Run("rejects a token naming somebody who does not exist", func(t *testing.T) {
		env := newTestUserEnv(t)
		forged, _, errToken := env.tokenService.Generate("mallory", jwt.MapClaims{"jti": "id-1"}, time.Hour)
		if errToken != nil {
			t.Fatalf("failed to sign: %v", errToken)
		}

		if _, err := env.service.RegistrationVerify(forged); !errors.Is(err, ErrRegistrationObsolete) {
			t.Fatalf("expected ErrRegistrationObsolete, got %v", err)
		}
	})

	t.Run("rejects a token without a registration id", func(t *testing.T) {
		env := newTestUserEnv(t)
		anonymous, _, errToken := env.tokenService.Generate("alice", nil, time.Hour)
		if errToken != nil {
			t.Fatalf("failed to sign: %v", errToken)
		}

		if _, err := env.service.RegistrationVerify(anonymous); !errors.Is(err, ErrRegistrationInvalid) {
			t.Fatalf("expected ErrRegistrationInvalid, got %v", err)
		}
	})

	// NOTE: the record carries the expiration as well as the token, because the
	// token is the half a stranger holds
	t.Run("reports a registration the record says has run out", func(t *testing.T) {
		env := newTestUserEnv(t)
		response, errRegister := env.service.Create(CreateRequest{Username: "alice", AuthTypes: []AuthType{AuthBasic}}, "")
		if errRegister != nil {
			t.Fatalf("failed to register alice: %v", errRegister)
		}
		stored, errGet := env.repository.Get("alice")
		if errGet != nil {
			t.Fatalf("failed to read alice: %v", errGet)
		}
		stored.RegistrationExpiresAt = time.Now().Add(-time.Minute)
		if err := env.repository.Update(stored); err != nil {
			t.Fatalf("failed to age the registration: %v", err)
		}

		if _, err := env.service.RegistrationVerify(response.Registration.Token); !errors.Is(err, ErrRegistrationExpired) {
			t.Fatalf("expected ErrRegistrationExpired, got %v", err)
		}
	})
}

func TestServiceRegistrationApply(t *testing.T) {
	t.Run("sets the password and spends the registration", func(t *testing.T) {
		env := newTestUserEnv(t)
		response, errRegister := env.service.Create(CreateRequest{Username: "alice", AuthTypes: []AuthType{AuthBasic}, Superuser: false}, "")
		if errRegister != nil {
			t.Fatalf("failed to register alice: %v", errRegister)
		}

		user, err := env.service.RegistrationApply(response.Registration.Token, "password123")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if user.Username != "alice" {
			t.Fatalf("expected alice, got %+v", user)
		}
		if user.Registration == nil || user.Registration.Status != RegistrationActive {
			t.Fatalf("expected an active account, got %+v", user.Registration)
		}
		if errVerify := env.service.VerifyPassword("alice", "password123"); errVerify != nil {
			t.Fatalf("expected the new password to verify, got %v", errVerify)
		}
	})

	t.Run("the same token cannot be used twice", func(t *testing.T) {
		env := newTestUserEnv(t)
		response, errRegister := env.service.Create(CreateRequest{Username: "alice", AuthTypes: []AuthType{AuthBasic}}, "")
		if errRegister != nil {
			t.Fatalf("failed to register alice: %v", errRegister)
		}
		if _, err := env.service.RegistrationApply(response.Registration.Token, "password123"); err != nil {
			t.Fatalf("failed to apply: %v", err)
		}

		if _, err := env.service.RegistrationApply(response.Registration.Token, "another"); !errors.Is(err, ErrRegistrationObsolete) {
			t.Fatalf("expected ErrRegistrationObsolete, got %v", err)
		}
	})

	// NOTE: a reset keeps everything the account already is - that is the whole
	// difference between answering a forgotten password and starting again
	t.Run("a reset keeps the permissions and the superuser flag", func(t *testing.T) {
		env := newTestUserEnv(t)
		seedUser(t, env, "root", "password123", true)
		seedUser(t, env, "admin", "password123", true)
		registration, errIssue := env.service.PasswordResetIssue("root", "admin")
		if errIssue != nil {
			t.Fatalf("failed to issue: %v", errIssue)
		}

		user, err := env.service.RegistrationApply(registration.Token, "newpassword")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if !user.Superuser {
			t.Fatalf("expected root to still be a superuser, got %+v", user)
		}
		if errVerify := env.service.VerifyPassword("root", "newpassword"); errVerify != nil {
			t.Fatalf("expected the new password to verify, got %v", errVerify)
		}
	})

	t.Run("rejects an empty password", func(t *testing.T) {
		env := newTestUserEnv(t)
		response, errRegister := env.service.Create(CreateRequest{Username: "alice", AuthTypes: []AuthType{AuthBasic}}, "")
		if errRegister != nil {
			t.Fatalf("failed to register alice: %v", errRegister)
		}

		if _, err := env.service.RegistrationApply(response.Registration.Token, ""); !errors.Is(err, ErrPasswordCannotBeEmpty) {
			t.Fatalf("expected ErrPasswordCannotBeEmpty, got %v", err)
		}
	})

	t.Run("rejects a registration that is not valid", func(t *testing.T) {
		env := newTestUserEnv(t)

		if _, err := env.service.RegistrationApply("garbage", "password123"); !errors.Is(err, ErrRegistrationInvalid) {
			t.Fatalf("expected ErrRegistrationInvalid, got %v", err)
		}
	})

	// NOTE: a name deleted and registered again is a different person, and the
	// new record carries an id of its own - the old token opens nothing
	t.Run("a registration written for a deleted account opens nothing", func(t *testing.T) {
		env := newTestUserEnv(t)
		seedUser(t, env, "root", "password123", true)
		response, errRegister := env.service.Create(CreateRequest{Username: "alice", AuthTypes: []AuthType{AuthBasic}}, "root")
		if errRegister != nil {
			t.Fatalf("failed to register alice: %v", errRegister)
		}
		if err := env.service.Delete("alice", "root"); err != nil {
			t.Fatalf("failed to delete alice: %v", err)
		}
		if _, err := env.service.Create(CreateRequest{Username: "alice", AuthTypes: []AuthType{AuthBasic}}, "root"); err != nil {
			t.Fatalf("failed to register alice again: %v", err)
		}

		if _, err := env.service.RegistrationApply(response.Registration.Token, "password123"); !errors.Is(err, ErrRegistrationObsolete) {
			t.Fatalf("expected ErrRegistrationObsolete, got %v", err)
		}
	})
}
