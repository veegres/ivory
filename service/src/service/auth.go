package service

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"strings"
	"time"
)

type AuthService struct {
	secretService *SecretService

	signingAlgorithm *jwt.SigningMethodHMAC
	issuer           string
}

func NewAuthService(secretService *SecretService) *AuthService {
	return &AuthService{
		secretService:    secretService,
		signingAlgorithm: jwt.SigningMethodHS256,
		issuer:           "ivory",
	}
}

func (s *AuthService) GetIssuer() string {
	return s.issuer
}

func (s *AuthService) GetTokenFromHeader(str string) (string, error) {
	if str == "" {
		return "", errors.New("no authorization token")
	}
	parts := strings.SplitN(str, " ", 2)
	if !(len(parts) == 2 && parts[0] == "Bearer") {
		return "", errors.New("invalid authorisation header")
	}
	return parts[1], nil
}

func (s *AuthService) GenerateToken(username string) (string, time.Time, error) {
	now := time.Now()
	exp := now.Add(time.Hour)
	t := jwt.NewWithClaims(
		s.signingAlgorithm,
		jwt.MapClaims{
			"iss": s.issuer,   // issuer
			"sub": username,   // subject
			"iat": now.Unix(), // issued at
			"exp": exp.Unix(), // expiration time
		})
	token, signErr := t.SignedString(s.secretService.GetByte())
	return token, exp, signErr
}

func (s *AuthService) ValidateToken(token string, username string) error {
	_, err := jwt.Parse(
		token,
		func(t *jwt.Token) (interface{}, error) {
			return s.secretService.GetByte(), nil
		},
		jwt.WithIssuedAt(),
		jwt.WithIssuer(s.issuer),
		jwt.WithSubject(username),
	)
	return err
}
