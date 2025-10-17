package auth

import (
	"context"
	"errors"
	"fmt"
	"ivory/src/features/config"
	"ivory/src/features/secret"
	"log/slog"
	"net"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-ldap/ldap/v3"
	"github.com/golang-jwt/jwt/v5"
)

type Service struct {
	secretService *secret.Service
	configService *config.Service

	signingAlgorithm *jwt.SigningMethodHMAC
	issuer           string
	expiration       time.Duration
}

func NewService(
	secretService *secret.Service,
	configService *config.Service,
) *Service {
	return &Service{
		secretService:    secretService,
		configService:    configService,
		signingAlgorithm: jwt.SigningMethodHS256,
		issuer:           "ivory",
		expiration:       time.Hour,
	}
}

func (s *Service) getIssuer() string {
	return s.issuer
}

func (s *Service) getConfig() (*config.AppConfig, error) {
	return s.configService.GetAppConfig()
}

func (s *Service) getOAuthCodeURL(str string) (string, error) {
	oauthConfig, _, err := s.configService.GetOAuthConfig()
	if err != nil {
		return "", err
	}
	return oauthConfig.AuthCodeURL(str), nil
}

func (s *Service) ValidateAuthToken(context *gin.Context, authConfig config.AuthConfig) (bool, error) {
	token, errToken := s.getToken(context)
	if errToken != nil {
		return false, errToken
	}
	switch authConfig.Type {
	case config.NONE:
		if token != "" {
			return false, errors.New("token should be empty")
		}
		return true, nil
	case config.BASIC, config.LDAP, config.OIDC:
		errValidate := s.validateToken(token)
		if errValidate != nil {
			return false, errValidate
		}
		return true, nil
	default:
		return false, errors.New("the authentication type is not specified or incorrect")
	}
}

func (s *Service) GenerateBasicAuthToken(login Login, authConfig config.AuthConfig) (string, *time.Time, error) {
	if login.Username != authConfig.Basic.Username || login.Password != authConfig.Basic.Password {
		return "", nil, errors.New("credentials are not correct")
	}
	return s.generateToken(authConfig.Basic.Username)
}

func (s *Service) GenerateLdapAuthToken(login Login, authConfig config.AuthConfig) (string, *time.Time, error) {
	valid, err := s.authenticateLdap(login, authConfig.Ldap)
	if err != nil {
		return "", nil, err
	}
	if !valid {
		return "", nil, errors.New("credentials are not correct")
	}
	return s.generateToken(login.Username)
}

func (s *Service) GenerateOidcAuthToken(code string) (string, *time.Time, error) {
	oauthConfig, oauthVerifier, err := s.configService.GetOAuthConfig()
	if err != nil {
		return "", nil, err
	}

	oauthToken, errExchange := oauthConfig.Exchange(context.Background(), code)
	if errExchange != nil {
		return "", nil, errors.Join(errors.New("failed to exchange token"), errExchange)
	}

	rawIDToken, ok := oauthToken.Extra("id_token").(string)
	if !ok {
		return "", nil, errors.New("no id_token field in oauth2 token")
	}

	idToken, errVerify := oauthVerifier.Verify(context.Background(), rawIDToken)
	if errVerify != nil {
		return "", nil, errors.Join(errors.New("failed to verify ID Token"), errVerify)
	}

	var claims struct {
		Email string `json:"email"`
	}
	if err := idToken.Claims(&claims); err != nil {
		return "", nil, err
	}
	return s.generateToken(claims.Email)
}

func (s *Service) authenticateLdap(login Login, ldapConfig *config.LdapConfig) (bool, error) {
	// TODO ldap can be configure to to use TLS, consider adding it
	dialer := &net.Dialer{Timeout: 3 * time.Second}
	l, err := ldap.DialURL("ldap://"+ldapConfig.Url, ldap.DialWithDialer(dialer))
	if err != nil {
		return false, err
	}
	defer func(l *ldap.Conn) {
		err := l.Close()
		if err != nil {
			slog.Error(err.Error())
		}
	}(l)

	err = l.Bind(ldapConfig.BindDN, ldapConfig.BindPass)
	if err != nil {
		return false, err
	}

	filter := ldapConfig.Filter
	if filter == "" {
		filter = "(uid=%s)"
	}

	searchRequest := ldap.NewSearchRequest(
		ldapConfig.BaseDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf(filter, login.Username),
		[]string{"dn"},
		nil,
	)

	sr, err := l.Search(searchRequest)
	if err != nil {
		return false, err
	}

	if len(sr.Entries) != 1 {
		return false, errors.New("user not found or too many entries returned")
	}

	userdn := sr.Entries[0].DN

	err = l.Bind(userdn, login.Password)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (s *Service) generateToken(subject string) (string, *time.Time, error) {
	now := time.Now()
	exp := now.Add(s.expiration)
	t := jwt.NewWithClaims(
		s.signingAlgorithm,
		jwt.MapClaims{
			"iss": s.issuer,   // issuer
			"sub": subject,    // subject
			"iat": now.Unix(), // issued at
			"exp": exp.Unix(), // expiration time
		})
	token, signErr := t.SignedString(s.secretService.GetByte())
	return token, &exp, signErr
}

func (s *Service) validateToken(rawIDToken string) error {
	token, err := jwt.Parse(
		rawIDToken,
		func(t *jwt.Token) (interface{}, error) {
			return s.secretService.GetByte(), nil
		},
		jwt.WithIssuedAt(),
		jwt.WithIssuer(s.issuer),
	)
	if err != nil {
		return err
	}
	if token.Valid {
		return nil
	}
	return errors.New("invalid token")
}

func (s *Service) getToken(context *gin.Context) (string, error) {
	authHeader := context.Request.Header.Get("Authorization")
	if authHeader != "" {
		return s.getTokenFromHeader(authHeader)
	}
	cookieToken, errToken := context.Cookie("token")
	if cookieToken == "" || errToken != nil {
		return "", errors.New("no authorization token")
	}
	return cookieToken, nil
}

func (s *Service) getTokenFromHeader(str string) (string, error) {
	if str == "" {
		return "", errors.New("no authorization token")
	}
	parts := strings.SplitN(str, " ", 2)
	if !(len(parts) == 2 && parts[0] == "Bearer") {
		return "", errors.New("invalid authorisation header")
	}
	return parts[1], nil
}
