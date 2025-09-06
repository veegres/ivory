package auth

import (
	"errors"
	"ivory/src/features/config"
	"ivory/src/features/encryption"
	"ivory/src/features/secret"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Service struct {
	secretService *secret.Service
	encryption    *encryption.Service

	signingAlgorithm *jwt.SigningMethodHMAC
	issuer           string
	expiration       time.Duration
}

func NewService(
	secretService *secret.Service,
	encryption *encryption.Service,
) *Service {
	return &Service{
		secretService:    secretService,
		encryption:       encryption,
		signingAlgorithm: jwt.SigningMethodHS256,
		issuer:           "ivory",
		expiration:       time.Hour,
	}
}

func (s *Service) GetIssuer() string {
	return s.issuer
}

func (s *Service) GenerateAuthToken(login Login, authConfig config.Auth) (string, *time.Time, error) {
	switch authConfig.Type {
	case config.NONE:
		return "", nil, errors.New("authentication is not used")
	case config.BASIC:
		encryptedPassword, errEnc := s.encryption.Encrypt(login.Password, s.secretService.Get())
		if errEnc != nil {
			return "", nil, errEnc
		}
		basicUser := authConfig.Body["username"]
		basicPass := authConfig.Body["password"]
		if login.Username != basicUser || encryptedPassword != basicPass {
			return "", nil, errors.New("credentials are not correct")
		}
		return s.generateToken(basicUser)
	default:
		return "", nil, errors.New("authentication type is not specified or incorrect")
	}
}

func (s *Service) ValidateAuthHeader(header string, authConfig config.Auth) (bool, string) {
	switch authConfig.Type {
	case config.NONE:
		if header != "" {
			return false, "authentication header should be empty"
		}
		return true, ""
	case config.BASIC:
		token, errHeader := s.getTokenFromHeader(header)
		if errHeader != nil {
			return false, errHeader.Error()
		} else {
			errValidate := s.validateToken(token, authConfig.Body["username"])
			if errValidate != nil {
				return false, errValidate.Error()
			}
		}
		return true, ""
	default:
		return false, "authentication type is not specified or incorrect"
	}
}

func (s *Service) generateToken(username string) (string, *time.Time, error) {
	now := time.Now()
	exp := now.Add(s.expiration)
	t := jwt.NewWithClaims(
		s.signingAlgorithm,
		jwt.MapClaims{
			"iss": s.issuer,   // issuer
			"sub": username,   // subject
			"iat": now.Unix(), // issued at
			"exp": exp.Unix(), // expiration time
		})
	token, signErr := t.SignedString(s.secretService.GetByte())
	return token, &exp, signErr
}

func (s *Service) validateToken(token string, username string) error {
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
