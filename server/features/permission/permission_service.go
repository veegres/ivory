package permission

import (
	"errors"
	"fmt"
	"ivory/clients/storage"
	"ivory/core/config"
	"slices"
	"sort"
)

var ErrUsernameCannotBeEmpty = errors.New("username cannot be empty")
var ErrUserPermissionsNotFound = errors.New("this user does not exist any more")
var ErrInvalidFeature = errors.New("invalid feature")

type Service struct {
	permissionRepository *Repository
}

func NewService(permissionRepository *Repository) *Service {
	return &Service{permissionRepository: permissionRepository}
}

func (s *Service) GetAllUserPermissions() ([]UserPermissions, error) {
	permissionsMap, err := s.permissionRepository.GetAll()
	if err != nil {
		return nil, err
	}

	result := make([]UserPermissions, 0, len(permissionsMap))
	for username, permissions := range permissionsMap {
		result = append(result, UserPermissions{
			Username:    username,
			Permissions: permissions,
		})
	}

	sort.Slice(result, func(i, j int) bool { return result[i].Username < result[j].Username })
	return result, nil
}

// GetUserPermissions answers what this session may do. Whatever the answer, the
// features the caller withholds are taken back out of it: a permission nobody
// can hold in this Ivory is not this feature's business to decide, so the
// caller states it.
func (s *Service) GetUserPermissions(username string, allowAll bool, withheld []config.Feature) (PermissionMap, error) {
	if allowAll {
		return s.getAllPermissionsWithStatus(GRANTED).without(withheld), nil
	}
	if username == "" {
		return nil, ErrUsernameCannotBeEmpty
	}
	permissions, err := s.permissionRepository.Get(username)
	if err != nil {
		// NOTE: a signed-in person whose record is gone was deleted while their
		// token was still valid - saying so beats reporting a missing element
		if errors.Is(err, storage.ErrNotFound) {
			return nil, ErrUserPermissionsNotFound
		}
		return nil, err
	}
	return permissions.without(withheld), nil
}

// CreateUserPermissions makes sure a user has the record they are entitled to:
// a missing one is created with defaultStatus, and an existing one keeps every
// answer it already carries while a feature added since it was written is
// filled in with defaultStatus too. The caller states that default because the
// caller is the one who knows what kind of user this is - this feature has
// never heard of a superuser.
func (s *Service) CreateUserPermissions(username string, defaultStatus Status) (PermissionMap, error) {
	if username == "" {
		return nil, ErrUsernameCannotBeEmpty
	}
	existing, err := s.permissionRepository.Get(username)
	if err != nil && !errors.Is(err, storage.ErrNotFound) {
		return nil, err
	}

	permissions := make(PermissionMap, len(config.All))
	renamed := existing.renamed()
	for _, feature := range config.All {
		if status, ok := renamed[feature]; ok {
			permissions[feature] = status
		} else {
			permissions[feature] = defaultStatus
		}
	}
	if errCreate := s.permissionRepository.CreateOrUpdate(username, permissions); errCreate != nil {
		return nil, errCreate
	}
	return permissions, nil
}

func (s *Service) RequestUserPermissions(username string, featuresList []config.Feature) error {
	return s.updateUserPermissions(username, featuresList, PENDING)
}

func (s *Service) ApproveUserPermissions(username string, featuresList []config.Feature) error {
	return s.updateUserPermissions(username, featuresList, GRANTED)
}

func (s *Service) RejectUserPermissions(username string, featuresList []config.Feature) error {
	return s.updateUserPermissions(username, featuresList, NOT_PERMITTED)
}

func (s *Service) DeleteUserPermissions(username string) error {
	if username == "" {
		return ErrUsernameCannotBeEmpty
	}
	return s.permissionRepository.Delete(username)
}

func (s *Service) DeleteAll() error {
	return s.permissionRepository.DeleteAll()
}

func (s *Service) UpdateUserPermissions(username string, permissions PermissionMap) error {
	if username == "" {
		return ErrUsernameCannotBeEmpty
	}
	return s.permissionRepository.CreateOrUpdate(username, permissions)
}

func (s *Service) isValidFeature(feature config.Feature) bool {
	return slices.Contains(config.All, feature)
}

func (s *Service) getAllPermissionsWithStatus(status Status) PermissionMap {
	permissions := make(PermissionMap)
	for _, feature := range config.All {
		permissions[feature] = status
	}
	return permissions
}

func (s *Service) updateUserPermissions(username string, featuresList []config.Feature, status Status) error {
	var err error
	for _, feature := range featuresList {
		errPerm := s.updateUserPermission(username, feature, status)
		if errPerm != nil {
			err = errors.Join(err, fmt.Errorf("%s: %w", feature, errPerm))
		}
	}
	return err
}

func (s *Service) updateUserPermission(username string, feature config.Feature, status Status) error {
	if username == "" {
		return ErrUsernameCannotBeEmpty
	}
	if !s.isValidFeature(feature) {
		return ErrInvalidFeature
	}

	existingPermissions, err := s.permissionRepository.Get(username)
	if err != nil {
		return err
	}

	existingPermissions[feature] = status
	return s.permissionRepository.CreateOrUpdate(username, existingPermissions)
}
