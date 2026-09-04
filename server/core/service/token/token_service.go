package token

import (
	"errors"
	"ivory/core/service/secret"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var ErrInvalidToken = errors.New("invalid token")

type Service struct {
	secretService *secret.Service

	// NOTE: For HMAC signing method, the key can be any []byte.
	// You need the same key for signing and validating. Whereas for RSA
	// you need public and private key.
	signingAlgorithm *jwt.SigningMethodHMAC
	issuer           string
}

func NewService(secretService *secret.Service) *Service {
	return &Service{
		secretService: secretService,

		signingAlgorithm: jwt.SigningMethodHS256,
		issuer:           "ivory",
	}
}

func (s *Service) Issuer() string {
	return s.issuer
}

// Generate signs a token with the standard claims plus everything in extra,
// which can never overwrite them.
func (s *Service) Generate(subject string, extra jwt.MapClaims, expiration time.Duration) (string, *time.Time, error) {
	now := time.Now()
	exp := now.Add(expiration)
	claims := jwt.MapClaims{}
	for key, value := range extra {
		claims[key] = value
	}
	claims["iss"] = s.issuer   // issuer
	claims["sub"] = subject    // subject
	claims["iat"] = now.Unix() // issued at
	claims["exp"] = exp.Unix() // expiration time

	t := jwt.NewWithClaims(s.signingAlgorithm, claims)
	token, errSign := t.SignedString(s.secretService.GetByte())
	if errSign != nil {
		return "", nil, errSign
	}
	return token, &exp, nil
}

func (s *Service) Parse(rawToken string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(
		rawToken,
		func(t *jwt.Token) (interface{}, error) {
			return s.secretService.GetByte(), nil
		},
		jwt.WithValidMethods([]string{s.signingAlgorithm.Alg()}),
		jwt.WithIssuedAt(),
		jwt.WithIssuer(s.issuer),
	)
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !token.Valid || !ok {
		return nil, ErrInvalidToken
	}
	return claims, nil
}
