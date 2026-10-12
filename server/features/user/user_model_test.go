package user

import (
	"ivory/core/utils"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin/binding"
)

func TestUserToResponse(t *testing.T) {
	t.Run("carries the username, the ways in and the superuser flag, never the password", func(t *testing.T) {
		stored := User{
			Username:  "alice",
			Password:  "encrypted",
			AuthTypes: []AuthType{AuthBasic, AuthLdap},
			Superuser: true,
		}

		response := stored.toResponse()

		if response.Username != "alice" {
			t.Fatalf("expected username 'alice', got %q", response.Username)
		}
		if !response.Superuser {
			t.Fatalf("expected a superuser")
		}
		if len(response.AuthTypes) != 2 {
			t.Fatalf("expected both ways of signing in, got %v", response.AuthTypes)
		}
	})
}

func TestUserRegistrationState(t *testing.T) {
	future := time.Now().Add(time.Hour)
	past := time.Now().Add(-time.Hour)

	tests := []struct {
		name     string
		user     User
		expected *RegistrationStatus
	}{
		{
			name: "a user who does not sign in with a password is not asked about one",
			user: User{AuthTypes: []AuthType{AuthLdap, AuthOidc}, Password: ""},
		},
		{
			name:     "a password set and nothing outstanding is an active account",
			user:     User{AuthTypes: []AuthType{AuthBasic}, Password: "encrypted"},
			expected: statusPtr(RegistrationActive),
		},
		{
			name:     "an outstanding registration is pending",
			user:     User{AuthTypes: []AuthType{AuthBasic}, RegistrationId: "id-1", RegistrationExpiresAt: future},
			expected: statusPtr(RegistrationPending),
		},
		{
			name:     "an outstanding registration past its time has expired",
			user:     User{AuthTypes: []AuthType{AuthBasic}, RegistrationId: "id-1", RegistrationExpiresAt: past},
			expected: statusPtr(RegistrationExpired),
		},
		{
			// NOTE: this is what a restored user looks like - a backup carries no
			// password, so somebody has to hand them a registration
			name:     "no password and nothing outstanding is a registration still to issue",
			user:     User{AuthTypes: []AuthType{AuthBasic}},
			expected: statusPtr(RegistrationMissing),
		},
		{
			// NOTE: a registration outstanding beside a password set is a reset
			// somebody asked for - the old password keeps working until it is spent
			name:     "a registration beside a password is still pending",
			user:     User{AuthTypes: []AuthType{AuthBasic}, Password: "encrypted", RegistrationId: "id-1", RegistrationExpiresAt: future},
			expected: statusPtr(RegistrationPending),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			state := tt.user.registrationState()

			if tt.expected == nil {
				if state != nil {
					t.Fatalf("expected no registration state, got %+v", state)
				}
				return
			}
			if state == nil {
				t.Fatalf("expected %v, got no state at all", *tt.expected)
			}
			if state.Status != *tt.expected {
				t.Fatalf("expected %v, got %v", *tt.expected, state.Status)
			}
		})
	}
}

func TestUserAllows(t *testing.T) {
	user := User{AuthTypes: []AuthType{AuthLdap}}

	if !user.allows(AuthLdap) {
		t.Fatalf("expected ldap to be allowed")
	}
	if user.allows(AuthBasic) {
		t.Fatalf("expected basic not to be allowed")
	}
}

func TestAuthTypeValid(t *testing.T) {
	for _, authType := range AuthTypes {
		if !authType.Valid() {
			t.Fatalf("expected %q to be valid", authType)
		}
	}
	if AuthType("smoke-signals").Valid() {
		t.Fatalf("expected an unknown way of signing in to be invalid")
	}
}

func statusPtr(status RegistrationStatus) *RegistrationStatus {
	return &status
}

func TestRequestBinding(t *testing.T) {
	utils.RegisterValidators()
	tests := []struct {
		name    string
		request any
		want    string
	}{
		{name: "create passes", request: &CreateRequest{Username: "alice", AuthTypes: []AuthType{AuthBasic, AuthOidc}}},
		{name: "create needs a username", request: &CreateRequest{Username: " ", AuthTypes: []AuthType{AuthBasic}}, want: "'CreateRequest.Username' Error:Field validation for 'Username' failed on the 'notblank' tag"},
		{name: "create needs an auth type", request: &CreateRequest{Username: "alice", AuthTypes: []AuthType{}}, want: "'CreateRequest.AuthTypes' Error:Field validation for 'AuthTypes' failed on the 'min' tag"},
		{name: "create refuses an unknown auth type", request: &CreateRequest{Username: "alice", AuthTypes: []AuthType{"kerberos"}}, want: "'CreateRequest.AuthTypes[0]' Error:Field validation for 'AuthTypes[0]' failed on the 'valid' tag"},
		{name: "nobody is registered for a minted session", request: &CreateRequest{Username: "alice", AuthTypes: []AuthType{AuthSession}}, want: "'CreateRequest.AuthTypes[0]' Error:Field validation for 'AuthTypes[0]' failed on the 'ne' tag"},
		{name: "update refuses a minted session too", request: &RouterUpdateRequest{AuthTypes: []AuthType{AuthLdap, AuthSession}}, want: "'RouterUpdateRequest.AuthTypes[1]' Error:Field validation for 'AuthTypes[1]' failed on the 'ne' tag"},
		{name: "password update needs both passwords", request: &RouterPasswordUpdateRequest{NewPassword: "x"}, want: "'RouterPasswordUpdateRequest.PreviousPassword' Error:Field validation for 'PreviousPassword' failed on the 'required' tag"},
		{name: "registration needs a token", request: &RouterRegistrationPasswordRequest{Password: "x"}, want: "'RouterRegistrationPasswordRequest.Token' Error:Field validation for 'Token' failed on the 'required' tag"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := binding.Validator.ValidateStruct(test.request)
			if test.want == "" {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				return
			}
			if err == nil || !strings.Contains(err.Error(), test.want) {
				t.Fatalf("expected an error containing %q, got %v", test.want, err)
			}
		})
	}
}
