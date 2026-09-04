package auth

import (
	"errors"
	"ivory/clients/auth/basic"
	"ivory/clients/auth/ldap"
	"ivory/clients/auth/oidc"
	"ivory/core/service/token"
	"ivory/features/permission"
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

type Service struct {
	tokenService      *token.Service
	basicProvider     *basic.Provider
	ldapProvider      *ldap.Provider
	oidcProvider      *oidc.Provider
	permissionService *permission.Service

	expiration time.Duration
}

func NewService(
	tokenService *token.Service,
	basicProvider *basic.Provider,
	ldapProvider *ldap.Provider,
	oidcProvider *oidc.Provider,
	permissionService *permission.Service,
) *Service {
	return &Service{
		tokenService:      tokenService,
		basicProvider:     basicProvider,
		ldapProvider:      ldapProvider,
		oidcProvider:      oidcProvider,
		permissionService: permissionService,

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
	_, errPerm := s.permissionService.CreateUserPermissions(BASIC.Prefix(), username)
	if errPerm != nil {
		return "", nil, errPerm
	}
	return s.generateToken(username, BASIC)
}

func (s *Service) GenerateLdapAuthToken(login ldap.Login) (string, *time.Time, error) {
	username, err := s.ldapProvider.Verify(login)
	if err != nil {
		return "", nil, err
	}
	_, errPerm := s.permissionService.CreateUserPermissions(LDAP.Prefix(), username)
	if errPerm != nil {
		return "", nil, errPerm
	}
	return s.generateToken(username, LDAP)
}

func (s *Service) GenerateOidcAuthToken(code string) (string, *time.Time, error) {
	username, err := s.oidcProvider.Verify(code)
	if err != nil {
		return "", nil, err
	}
	_, errPerm := s.permissionService.CreateUserPermissions(OIDC.Prefix(), username)
	if errPerm != nil {
		return "", nil, errPerm
	}
	return s.generateToken(username, OIDC)
}

func (s *Service) generateToken(subject string, authType AuthType) (string, *time.Time, error) {
	return s.tokenService.Generate(subject, jwt.MapClaims{"frm": authType}, s.expiration)
}
