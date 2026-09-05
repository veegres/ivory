package user

import "slices"

import "time"

// AuthType is a way of signing in that a user was registered for. It is the
// source of truth for the three methods: features/auth spells its own type out
// with it rather than the two drifting apart.
type AuthType string

const (
	AuthBasic AuthType = "basic"
	AuthLdap  AuthType = "ldap"
	AuthOidc  AuthType = "oidc"
)

var AuthTypes = []AuthType{AuthBasic, AuthLdap, AuthOidc}

func (t AuthType) valid() bool {
	return slices.Contains(AuthTypes, t)
}

// RegistrationStatus is where a user stands with the one thing Ivory holds for
// them, their password. It is stated by the server so every screen reads the
// same answer out of the same two fields.
type RegistrationStatus string

const (
	RegistrationActive  RegistrationStatus = "active"
	RegistrationPending RegistrationStatus = "pending"
	RegistrationExpired RegistrationStatus = "expired"
	RegistrationMissing RegistrationStatus = "missing"
)

type Response struct {
	Username     string             `json:"username"`
	AuthTypes    []AuthType         `json:"authTypes"`
	Superuser    bool               `json:"superuser"`
	Registration *RegistrationState `json:"registration,omitempty"`
}

type RegistrationState struct {
	Status    RegistrationStatus `json:"status"`
	ExpiresAt *time.Time         `json:"expiresAt,omitempty"`
}

// RegistrationResponse is the only shape carrying the token: a registration is
// shown once, when it is issued, and is never stored to be shown again.
type RegistrationResponse struct {
	Token     string    `json:"token"`
	Username  string    `json:"username"`
	ExpiresAt time.Time `json:"expiresAt"`
}

// CreateResponse is the user and, when one was issued along with them, the
// registration link to hand out.
type CreateResponse struct {
	User         Response              `json:"user"`
	Registration *RegistrationResponse `json:"registration,omitempty"`
}

// RegistrationPayload is what the public page shows to the person who opened
// the link.
type RegistrationPayload struct {
	Username  string    `json:"username"`
	ExpiresAt time.Time `json:"expiresAt"`
}

// CreateRequest carries no password: a user is created by name and by the ways
// they may sign in, and the password is set by its owner on the page the
// registration opens. Setup is the one exception, and it goes through the
// config feature instead.
type CreateRequest struct {
	Username  string     `json:"username" binding:"required"`
	AuthTypes []AuthType `json:"authTypes" binding:"required"`
	Superuser bool       `json:"superuser"`
}

type UpdateRequest struct {
	AuthTypes []AuthType `json:"authTypes" binding:"required"`
}

type PasswordUpdateRequest struct {
	PreviousPassword string `json:"previousPassword" binding:"required"`
	NewPassword      string `json:"newPassword" binding:"required"`
}

type RegistrationVerifyRequest struct {
	Token string `json:"token" binding:"required"`
}

type RegistrationPasswordRequest struct {
	Token    string `json:"token" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// SPECIFIC (SERVER)

// User is the stored shape. Password is encrypted with the secret key, never
// leaves the server, and is empty until its owner sets it. RegistrationId is
// what makes an outstanding registration one-time and revocable: the token a
// stranger holds is only good while it still matches. CreatedAt tells one
// account from another under the same name.
type User struct {
	Username              string     `json:"username"`
	Password              string     `json:"password"`
	AuthTypes             []AuthType `json:"authTypes"`
	Superuser             bool       `json:"superuser"`
	RegistrationId        string     `json:"registrationId"`
	RegistrationExpiresAt time.Time  `json:"registrationExpiresAt"`
	CreatedAt             time.Time  `json:"createdAt"`
}

func (u User) allows(authType AuthType) bool {
	return slices.Contains(u.AuthTypes, authType)
}

func (u User) toResponse() Response {
	return Response{
		Username:     u.Username,
		AuthTypes:    u.AuthTypes,
		Superuser:    u.Superuser,
		Registration: u.registrationState(),
	}
}

// registrationState answers the question the user list asks: where is this
// person with their password. A user who cannot sign in with one is not asked
// it at all.
func (u User) registrationState() *RegistrationState {
	if !u.allows(AuthBasic) {
		return nil
	}
	if u.RegistrationId != "" {
		expiresAt := u.RegistrationExpiresAt
		if time.Now().After(expiresAt) {
			return &RegistrationState{Status: RegistrationExpired, ExpiresAt: &expiresAt}
		}
		return &RegistrationState{Status: RegistrationPending, ExpiresAt: &expiresAt}
	}
	if u.Password == "" {
		return &RegistrationState{Status: RegistrationMissing}
	}
	return &RegistrationState{Status: RegistrationActive}
}
