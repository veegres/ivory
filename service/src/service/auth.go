package service

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	. "ivory/src/model"
	"strings"
	"time"
)

type AuthService struct {
	secretService *SecretService
	encryption    *EncryptionService

	signingAlgorithm *jwt.SigningMethodHMAC
	issuer           string
}

func NewAuthService(
	secretService *SecretService,
	encryption *EncryptionService,
) *AuthService {
	return &AuthService{
		secretService:    secretService,
		encryption:       encryption,
		signingAlgorithm: jwt.SigningMethodHS256,
		issuer:           "ivory",
	}
}

func (s *AuthService) GetIssuer() string {
	return s.issuer
}

func (s *AuthService) Login(login Login, authConfig AuthConfig) (string, *time.Time, error) {
	switch authConfig.Type {
	case NONE:
		return "", nil, errors.New("authentication is not used")
	case BASIC:
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

func (s *AuthService) ValidateHeader(header string, authConfig AuthConfig) (bool, string) {
	switch authConfig.Type {
	case NONE:
		if header != "" {
			return false, "authentication header should be empty"
		}
		return true, ""
	case BASIC:
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

func (s *AuthService) EncryptAuthConfig(authConfig AuthConfig) (AuthConfig, error) {
	switch authConfig.Type {
	case NONE:
		return authConfig, nil
	case BASIC:
		user := authConfig.Body["username"]
		pass := authConfig.Body["password"]
		if pass == "" || user == "" {
			return authConfig, errors.New("username or password are not specified")
		}

		encryptedPassword, errEnc := s.encryption.Encrypt(pass, s.secretService.Get())
		if errEnc != nil {
			return authConfig, errEnc
		}
		authConfig.Body["password"] = encryptedPassword

		return authConfig, nil
	default:
		return authConfig, errors.New("type is unknown")
	}
}

func (s *AuthService) ReEncryptAuthConfig(authConfig AuthConfig, oldSecret [16]byte, newSecret [16]byte) (AuthConfig, error) {
	switch authConfig.Type {
	case NONE:
		return authConfig, nil
	case BASIC:
		user := authConfig.Body["username"]
		pass := authConfig.Body["password"]
		if pass == "" || user == "" {
			return authConfig, errors.New("username or password are not specified")
		}

		decryptedPassword, errDec := s.encryption.Decrypt(pass, oldSecret)
		if errDec != nil {
			return authConfig, errDec
		}
		encryptedPassword, errEnc := s.encryption.Encrypt(decryptedPassword, newSecret)
		if errEnc != nil {
			return authConfig, errEnc
		}
		authConfig.Body["password"] = encryptedPassword

		return authConfig, nil
	default:
		return authConfig, errors.New("type is unknown")
	}
}

func (s *AuthService) generateToken(username string) (string, *time.Time, error) {
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
	return token, &exp, signErr
}

func (s *AuthService) validateToken(token string, username string) error {
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

func (s *AuthService) getTokenFromHeader(str string) (string, error) {
	if str == "" {
		return "", errors.New("no authorization token")
	}
	parts := strings.SplitN(str, " ", 2)
	if !(len(parts) == 2 && parts[0] == "Bearer") {
		return "", errors.New("invalid authorisation header")
	}
	return parts[1], nil
}
