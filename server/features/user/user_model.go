package user

import "time"

// COMMON (WEB AND SERVER)

// WebPath is the web sub-path an invitation link points at, the token being the
// segment after it. The frontend builds the link with it, the server has to let
// it through to the app.
const WebPath = "user"

// LinkKind is what the link does when it is spent: invite a user Ivory does not
// have yet, or set the password of one it already has. It is stated when the
// link is created rather than worked out when it is used, so a name that comes
// into existence in between cannot turn an invitation into a password reset.
type LinkKind string

const (
	LinkInvite LinkKind = "invite"
	LinkReset  LinkKind = "reset"
)

type Response struct {
	Username  string `json:"username"`
	Superuser bool   `json:"superuser"`
}

// LinkResponse describes an outstanding link. It never carries the token: the
// token is shown once, when the link is created, and is not stored anywhere.
type LinkResponse struct {
	Id        string    `json:"id"`
	Kind      LinkKind  `json:"kind"`
	Username  string    `json:"username"`
	Superuser bool      `json:"superuser"`
	CreatedAt time.Time `json:"createdAt"`
	ExpiresAt time.Time `json:"expiresAt"`
}

type LinkCreatedResponse struct {
	Id        string    `json:"id"`
	Kind      LinkKind  `json:"kind"`
	Token     string    `json:"token"`
	Username  string    `json:"username"`
	Superuser bool      `json:"superuser"`
	CreatedAt time.Time `json:"createdAt"`
	ExpiresAt time.Time `json:"expiresAt"`
}

// LinkPayload is what the public page shows to the person who opened the link.
type LinkPayload struct {
	Kind      LinkKind  `json:"kind"`
	Username  string    `json:"username"`
	Superuser bool      `json:"superuser"`
	ExpiresAt time.Time `json:"expiresAt"`
}

// LinkRequest asks for an invitation. A reset asks for nothing but a username,
// and says so through its own route, so the kind can never be smuggled past the
// permission that guards it.
type LinkRequest struct {
	Username  string `json:"username" binding:"required"`
	Superuser bool   `json:"superuser"`
}

type LinkResetRequest struct {
	Username string `json:"username" binding:"required"`
}

type PasswordUpdateRequest struct {
	PreviousPassword string `json:"previousPassword" binding:"required"`
	NewPassword      string `json:"newPassword" binding:"required"`
}

type LinkVerifyRequest struct {
	Token string `json:"token" binding:"required"`
}

type LinkPasswordRequest struct {
	Token    string `json:"token" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// SPECIFIC (SERVER)

// User is the stored shape; Password is encrypted with the secret key and never
// leaves the server. CreatedAt tells one account from another under the same
// name, which is what an outstanding reset link is checked against.
type User struct {
	Username  string    `json:"username"`
	Password  string    `json:"password"`
	Superuser bool      `json:"superuser"`
	CreatedAt time.Time `json:"createdAt"`
}

func (u User) toResponse() Response {
	return Response{Username: u.Username, Superuser: u.Superuser}
}

// Link is the stored half of an invitation. The token carries the same facts
// and its own expiration; this record is what makes the token one-time and
// revocable, so a link that is gone from here cannot be used any more.
type Link struct {
	Kind      LinkKind  `json:"kind"`
	Username  string    `json:"username"`
	Superuser bool      `json:"superuser"`
	CreatedAt time.Time `json:"createdAt"`
	ExpiresAt time.Time `json:"expiresAt"`
}

func (l Link) toResponse(id string) LinkResponse {
	return LinkResponse{
		Id:        id,
		Kind:      l.Kind,
		Username:  l.Username,
		Superuser: l.Superuser,
		CreatedAt: l.CreatedAt,
		ExpiresAt: l.ExpiresAt,
	}
}

func (l Link) toPayload() LinkPayload {
	return LinkPayload{Kind: l.Kind, Username: l.Username, Superuser: l.Superuser, ExpiresAt: l.ExpiresAt}
}

func (l Link) toCreatedResponse(id string, token string) LinkCreatedResponse {
	return LinkCreatedResponse{
		Id:        id,
		Kind:      l.Kind,
		Token:     token,
		Username:  l.Username,
		Superuser: l.Superuser,
		CreatedAt: l.CreatedAt,
		ExpiresAt: l.ExpiresAt,
	}
}
