package auth

import (
	"errors"
	"ivory/src/clients/auth/basic"
	"ivory/src/clients/auth/ldap"
	"ivory/src/clients/auth/oidc"
	"ivory/src/features/secret"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type Service struct {
	secretService *secret.Service
	basicProvider *basic.Provider
	ldapProvider  *ldap.Provider
	oidcProvider  *oidc.Provider

	// NOTE: For HMAC signing method, the key can be any []byte.
	// You need the same key for signing and validating. Whereas for RSA
	// you need public and private key.
	signingAlgorithm *jwt.SigningMethodHMAC
	issuer           string
	expiration       time.Duration
}

func NewService(
	secretService *secret.Service,
	basicProvider *basic.Provider,
	ldapProvider *ldap.Provider,
	oidcProvider *oidc.Provider,
) *Service {
	return &Service{
		secretService: secretService,
		basicProvider: basicProvider,
		ldapProvider:  ldapProvider,
		oidcProvider:  oidcProvider,

		signingAlgorithm: jwt.SigningMethodHS256,
		issuer:           "ivory",
		expiration:       time.Hour,
	}
}

func (s *Service) getIssuer() string {
	return s.issuer
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

func (s *Service) ValidateAuthToken(context *gin.Context) (bool, error) {
	if len(s.GetSupportedTypes()) == 0 {
		return true, errors.New("authorization is disabled")
	}
	token, errToken := s.getToken(context)
	if errToken != nil {
		return false, errToken
	}
	errValidate := s.validateToken(token)
	if errValidate != nil {
		return false, errValidate
	}
	return true, nil
}

func (s *Service) GenerateBasicAuthToken(login basic.Login) (string, *time.Time, error) {
	sub, err := s.basicProvider.Verify(login)
	if err != nil {
		return "", nil, err
	}
	return s.generateToken(sub, BASIC)
}

func (s *Service) GenerateLdapAuthToken(login ldap.Login) (string, *time.Time, error) {
	sub, err := s.ldapProvider.Verify(login)
	if err != nil {
		return "", nil, err
	}
	return s.generateToken(sub, LDAP)
}

func (s *Service) GenerateOidcAuthToken(code string) (string, *time.Time, error) {
	sub, err := s.oidcProvider.Verify(code)
	if err != nil {
		return "", nil, err
	}
	return s.generateToken(sub, OIDC)
}

func (s *Service) generateToken(subject string, authType AuthType) (string, *time.Time, error) {
	now := time.Now()
	exp := now.Add(s.expiration)
	t := jwt.NewWithClaims(
		s.signingAlgorithm,
		jwt.MapClaims{
			"iss": s.issuer,   // issuer
			"sub": subject,    // subject
			"iat": now.Unix(), // issued at
			"exp": exp.Unix(), // expiration time
			"frm": authType,   // generated from
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
		jwt.WithValidMethods([]string{s.signingAlgorithm.Alg()}),
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
