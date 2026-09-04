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
var ErrCredentialsNotCorrect = errors.New("credentials are not correct")
var ErrLastSuperuser = errors.New("the last superuser cannot be deleted")
var ErrCannotDeleteYourself = errors.New("you cannot delete yourself")
var ErrSuperuserRequired = errors.New("only a superuser can manage a superuser")
var ErrLinkInvalid = errors.New("the link is not valid")
var ErrLinkExpired = errors.New("the link has expired")
var ErrLinkObsolete = errors.New("the link is not valid any more, it was already used or has been revoked")
var ErrLinkIdCannotBeEmpty = errors.New("link identifier cannot be empty")

type Service struct {
	userRepository    *Repository
	encryptionService *encryption.Service
	secretService     *secret.Service
	permissionService *permission.Service
	tokenService      *token.Service

	linkExpiration time.Duration
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

		linkExpiration: 12 * time.Hour,
	}
	if err := s.syncSuperusers(); err != nil {
		slog.Error("failed to load superusers", "error", err)
	}
	return s
}
