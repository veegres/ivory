package auth

import (
	"errors"
	"ivory/clients/auth/basic"
	"ivory/clients/auth/ldap"
	"ivory/clients/auth/oidc"
	"ivory/core/service/token"
	"ivory/features/user"
	"log/slog"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var ErrAuthDisabled = errors.New("authorization is disabled")
var ErrInvalidTokenCannotParseAuthType = errors.New("invalid token: cannot parse auth type")
var ErrNoAuthorizationToken = errors.New("no authorization token")
var ErrInvalidAuthorizationHeader = errors.New("invalid authorisation header")
var ErrUsernameEmpty = errors.New("username cannot be empty")
var ErrInvalidAuthType = errors.New("invalid auth type")
var ErrStateCookieNotFound = errors.New("state cookie not found")
var ErrInvalidStateParameter = errors.New("invalid state parameter")

// users is the narrow view auth takes of the user store: a provider says who
// somebody is, and the store says whether Ivory knows that name at all and
// whether it may be signed in with this way.
type users interface {
	VerifySignIn(username string, authType user.AuthType) error
	IsSuperuser(username string) (bool, error)
}

type Service struct {
	tokenService  *token.Service
	basicProvider *basic.Provider
	ldapProvider  *ldap.Provider
	oidcProvider  *oidc.Provider
	userService   users

	expiration time.Duration
}

func NewService(
	tokenService *token.Service,
	basicProvider *basic.Provider,
	ldapProvider *ldap.Provider,
	oidcProvider *oidc.Provider,
	userService users,
) *Service {
	return &Service{
		tokenService:  tokenService,
		basicProvider: basicProvider,
		ldapProvider:  ldapProvider,
		oidcProvider:  oidcProvider,
		userService:   userService,

		expiration: time.Hour,
	}
}

func (s *Service) getIssuer() string {
	return s.tokenService.Issuer()
}

func (s *Service) GetSupportedTypes() []AuthType {
	supported := make([]AuthType, 0)
	if s.basicProvider.Configured() {
		supported = append(supported, BASIC)
	}
	if s.oidcProvider.Configured() {
		supported = append(supported, OIDC)
	}
	if s.ldapProvider.Configured() {
		supported = append(supported, LDAP)
	}
	return supported
}

// IsSuperuser reports the one thing about a signed-in person that no permission
// record can overrule. A failure is answered with false: a person who cannot be
// looked up is treated as an ordinary user, never as an administrator.
func (s *Service) IsSuperuser(username string) bool {
	if username == "" {
		return false
	}
	superuser, err := s.userService.IsSuperuser(username)
	if err != nil {
		slog.Error("failed to check whether the user is a superuser", "error", err)
		return false
	}
	return superuser
}

// ParseAuthTokenWithFallback tries primaryToken first and only falls back to fallbackToken if it doesn't validate.
func (s *Service) ParseAuthTokenWithFallback(primaryToken string, primaryErr error, fallbackToken string, fallbackErr error) (bool, string, *AuthType, error) {
	valid, username, authType, errParse := s.ParseAuthToken(primaryToken, primaryErr)
	if valid {
		return valid, username, authType, errParse
	}
	return s.ParseAuthToken(fallbackToken, fallbackErr)
}

func (s *Service) ParseAuthToken(token string, tokenErr error) (bool, string, *AuthType, error) {
	if len(s.GetSupportedTypes()) == 0 {
		return true, "", nil, ErrAuthDisabled
	}
	if tokenErr != nil {
		return false, "", nil, tokenErr
	}
	claims, errParse := s.tokenService.Parse(token)
	if errParse != nil {
		return false, "", nil, errParse
	}
	username, errUsername := claims.GetSubject()
	if errUsername != nil {
		return true, "", nil, errUsername
	}
	frm, ok := claims["frm"].(float64)
	if !ok {
		return true, "", nil, ErrInvalidTokenCannotParseAuthType
	}
	authType := AuthType(frm)
	return true, username, &authType, nil
}

func (s *Service) GenerateBasicAuthToken(login basic.Login) (string, *time.Time, error) {
	username, err := s.basicProvider.Verify(login)
	if err != nil {
		return "", nil, err
	}
	return s.generateToken(username, BASIC)
}

func (s *Service) GenerateLdapAuthToken(login ldap.Login) (string, *time.Time, error) {
	username, err := s.ldapProvider.Verify(login)
	if err != nil {
		return "", nil, err
	}
	return s.generateToken(username, LDAP)
}

func (s *Service) GenerateOidcAuthToken(code string) (string, *time.Time, error) {
	username, err := s.oidcProvider.Verify(code)
	if err != nil {
		return "", nil, err
	}
	return s.generateToken(username, OIDC)
}

// generateToken is where a verified person becomes a session, and the last gate
// before that: a directory may vouch for somebody Ivory was never told about,
// and being registered - for this way of signing in - is what says they are
// welcome here.
func (s *Service) generateToken(subject string, authType AuthType) (string, *time.Time, error) {
	if err := s.userService.VerifySignIn(subject, authType.User()); err != nil {
		return "", nil, err
	}
	return s.tokenService.Generate(subject, jwt.MapClaims{"frm": authType}, s.expiration)
}
