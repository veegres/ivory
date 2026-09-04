package user

import (
	"errors"
	"ivory/core/service/encryption"
	"ivory/core/service/secret"
	"ivory/core/service/token"
	"ivory/features/permission"
	"log/slog"
	"time"
)

var ErrUsernameCannotBeEmpty = errors.New("username cannot be empty")
var ErrPasswordCannotBeEmpty = errors.New("password cannot be empty")
var ErrUserAlreadyExists = errors.New("such a user already exists")
var ErrUserNotFound = errors.New("user is not found")
var ErrUserNotRegistered = errors.New("this user is not registered in ivory")
var ErrCredentialsNotCorrect = errors.New("credentials are not correct")
var ErrLastSuperuser = errors.New("the last superuser cannot be deleted")
var ErrCannotDeleteYourself = errors.New("you cannot delete yourself")
var ErrSuperuserRequired = errors.New("only a superuser can manage a superuser")
var ErrAuthTypeRequired = errors.New("a user needs at least one way to sign in")
var ErrAuthTypeInvalid = errors.New("there is no such way to sign in")
var ErrAuthTypeNotAllowed = errors.New("this user is not registered for this way of signing in")
var ErrRegistrationNotNeeded = errors.New("this user does not sign in with a password")
var ErrRegistrationInvalid = errors.New("the registration is not valid")
var ErrRegistrationExpired = errors.New("the registration has expired")
var ErrRegistrationObsolete = errors.New("the registration is not valid any more, it was already used or has been revoked")

type Service struct {
	userRepository    *Repository
	encryptionService *encryption.Service
	secretService     *secret.Service
	permissionService *permission.Service
	tokenService      *token.Service

	registrationExpiration time.Duration
}

func NewService(
	userRepository *Repository,
	encryptionService *encryption.Service,
	secretService *secret.Service,
	permissionService *permission.Service,
	tokenService *token.Service,
) *Service {
	s := &Service{
		userRepository:    userRepository,
		encryptionService: encryptionService,
		secretService:     secretService,
		permissionService: permissionService,
		tokenService:      tokenService,

		registrationExpiration: 12 * time.Hour,
	}
	if err := s.syncPermissions(); err != nil {
		slog.Error("failed to synchronise user permissions", "error", err)
	}
	return s
}
