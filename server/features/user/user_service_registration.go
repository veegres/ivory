package user

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// A registration is the one-time way a user sets their password: a token they
// hold plus the id Ivory keeps on their own record. The token carries the name
// and the expiration, the record is what a registration can be revoked and
// spent through - a token whose id no longer matches opens nothing. Issuing one
// is always done by somebody signed in; the initial setup creates its superuser
// with a typed password instead, so there is no public way to ask for one.

// PasswordResetIssue hands an existing user a new registration, which is how a
// forgotten password is answered: nobody ever types somebody else's password,
// and the account keeps its permissions. Taking a superuser's account over is a
// superuser's own right, exactly as deleting one is.
func (s *Service) PasswordResetIssue(username string, requester string) (*RegistrationResponse, error) {
	user, errGet := s.get(username)
	if errGet != nil {
		return nil, errGet
	}
	if user.Superuser {
		if errSuper := s.requireSuperuser(requester); errSuper != nil {
			return nil, errSuper
		}
	}
	if !user.allows(AuthBasic) {
		return nil, ErrRegistrationNotNeeded
	}
	registration, _, err := s.registrationIssue(*user)
	return registration, err
}

// PasswordResetRevoke makes an outstanding registration useless straight away.
// The password the user already has, if any, is left alone.
func (s *Service) PasswordResetRevoke(username string, requester string) error {
	user, errGet := s.get(username)
	if errGet != nil {
		return errGet
	}
	if user.Superuser {
		if errSuper := s.requireSuperuser(requester); errSuper != nil {
			return errSuper
		}
	}
	if user.RegistrationId == "" {
		return ErrRegistrationObsolete
	}
	user.RegistrationId = ""
	user.RegistrationExpiresAt = time.Time{}
	return s.userRepository.Update(*user)
}

func (s *Service) RegistrationVerify(rawToken string) (*RegistrationPayload, error) {
	user, err := s.registrationResolve(rawToken)
	if err != nil {
		return nil, err
	}
	return &RegistrationPayload{Username: user.Username, ExpiresAt: user.RegistrationExpiresAt}, nil
}

// RegistrationApply spends the registration - it sets the password of the user
// it names - and clears the id, so the same token cannot be used twice.
func (s *Service) RegistrationApply(rawToken string, password string) (*Response, error) {
	if password == "" {
		return nil, ErrPasswordCannotBeEmpty
	}
	user, errResolve := s.registrationResolve(rawToken)
	if errResolve != nil {
		return nil, errResolve
	}
	encrypted, errEnc := s.encrypt(password)
	if errEnc != nil {
		return nil, errEnc
	}
	user.Password = encrypted
	user.RegistrationId = ""
	user.RegistrationExpiresAt = time.Time{}
	if errUpdate := s.userRepository.Update(*user); errUpdate != nil {
		return nil, errUpdate
	}
	response := user.toResponse()
	return &response, nil
}

func (s *Service) registrationIssue(user User) (*RegistrationResponse, User, error) {
	id, errId := uuid.NewUUID()
	if errId != nil {
		return nil, user, errId
	}
	rawToken, exp, errToken := s.tokenService.Generate(user.Username, jwt.MapClaims{"jti": id.String()}, s.registrationExpiration)
	if errToken != nil {
		return nil, user, errToken
	}

	user.RegistrationId = id.String()
	user.RegistrationExpiresAt = *exp
	if errUpdate := s.userRepository.Update(user); errUpdate != nil {
		return nil, user, errUpdate
	}
	return &RegistrationResponse{Token: rawToken, Username: user.Username, ExpiresAt: *exp}, user, nil
}

// registrationResolve answers both halves of a registration: the token has to be
// signed by this Ivory and unexpired, and its id has to be the one the user
// still carries. A name deleted and registered again therefore opens nothing -
// the new record holds an id of its own, or none at all.
func (s *Service) registrationResolve(rawToken string) (*User, error) {
	if rawToken == "" {
		return nil, ErrRegistrationInvalid
	}
	claims, errParse := s.tokenService.Parse(rawToken)
	if errParse != nil {
		if errors.Is(errParse, jwt.ErrTokenExpired) {
			return nil, ErrRegistrationExpired
		}
		return nil, ErrRegistrationInvalid
	}
	id, okId := claims["jti"].(string)
	if !okId || id == "" {
		return nil, ErrRegistrationInvalid
	}
	username, errUsername := claims.GetSubject()
	if errUsername != nil || username == "" {
		return nil, ErrRegistrationInvalid
	}

	user, errUser := s.userRepository.Get(username)
	if errUser != nil {
		return nil, ErrRegistrationObsolete
	}
	if user.RegistrationId == "" || user.RegistrationId != id {
		return nil, ErrRegistrationObsolete
	}
	if time.Now().After(user.RegistrationExpiresAt) {
		return nil, ErrRegistrationExpired
	}
	return &user, nil
}
