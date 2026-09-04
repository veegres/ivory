package permission

import (
	"errors"
	"fmt"
	"ivory/clients/storage"
	"ivory/core/config"
	"slices"
	"sort"
	"strings"
)

var ErrSuperusersCannotHaveEmptyName = errors.New("superusers cannot have empty name")
var ErrUsernameCannotBeEmpty = errors.New("username cannot be empty")
var ErrUserPermissionsNotFound = errors.New("this user does not exist any more")
var ErrPrefixCannotBeEmpty = errors.New("prefix cannot be empty")
var ErrCannotChangePermissionsForSuperusers = errors.New("cannot change permissions for superusers")
var ErrInvalidFeature = errors.New("invalid feature")

type Service struct {
	permissionRepository *Repository

	superusers []string
}

func NewService(permissionRepository *Repository) *Service {
	return &Service{
		permissionRepository: permissionRepository,
		superusers:           []string{},
	}
}

// SetSuperusers replaces the list of usernames that hold every permission. An
// empty list is allowed: the user feature is what guarantees the last superuser
// cannot be deleted, and Ivory starts with none at all.
func (s *Service) SetSuperusers(superusers []string) error {
	if slices.Contains(superusers, "") {
		return ErrSuperusersCannotHaveEmptyName
	}
	s.superusers = superusers
	return s.normalizeDatabase()
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

func (s *Service) GetUserPermissions(prefix Prefix, username string, allowAll bool) (PermissionMap, error) {
	if allowAll {
		return s.getAllPermissionsWithStatus(GRANTED), nil
	}
	permUsername, errName := s.getFullUsername(prefix, username)
	if errName != nil {
		return nil, errName
	}
	permissions, err := s.permissionRepository.Get(permUsername)
	if err != nil {
		// NOTE: a signed-in person whose record is gone was deleted while their
		// token was still valid - saying so beats reporting a missing element
		if errors.Is(err, storage.ErrNotFound) {
			return nil, ErrUserPermissionsNotFound
		}
		return nil, err
	}
	return permissions, nil
}

func (s *Service) CreateUserPermissions(prefix Prefix, username string) (PermissionMap, error) {
	permUsername, errName := s.getFullUsername(prefix, username)
	if errName != nil {
		return nil, errName
	}
	existingPermissions, err := s.permissionRepository.Get(permUsername)
	if err != nil && !errors.Is(err, storage.ErrNotFound) {
		return nil, err
	}
	if existingPermissions != nil {
		return existingPermissions, nil
	}

	status := s.getStatus(permUsername)

	permissions := s.getAllPermissionsWithStatus(status)
	errCreate := s.permissionRepository.CreateOrUpdate(permUsername, permissions)
	if errCreate != nil {
		return nil, errCreate
	}

	return permissions, nil
}

func (s *Service) RequestUserPermissions(prefix Prefix, username string, featuresList []config.Feature) error {
	permUsername, errName := s.getFullUsername(prefix, username)
	if errName != nil {
		return errName
	}
	return s.updateUserPermissions(permUsername, featuresList, PENDING)
}

func (s *Service) ApproveUserPermissions(permUsername string, featuresList []config.Feature) error {
	return s.updateUserPermissions(permUsername, featuresList, GRANTED)
}

func (s *Service) RejectUserPermissions(permUsername string, featuresList []config.Feature) error {
	return s.updateUserPermissions(permUsername, featuresList, NOT_PERMITTED)
}

// DeleteBasicUserPermissions drops the record of an Ivory user, under whichever
// prefix it is stored: they hold the basic one, or the superuser one that takes
// its place. It exists so the user feature can say a person is gone without
// knowing anything about how permissions are keyed.
func (s *Service) DeleteBasicUserPermissions(username string) error {
	permUsername, errName := s.getFullUsername(PrefixBasic, username)
	if errName != nil {
		return errName
	}
	return s.permissionRepository.Delete(permUsername)
}

func (s *Service) DeleteUserPermissions(permUsername string) error {
	if permUsername == "" {
		return ErrUsernameCannotBeEmpty
	}
	return s.permissionRepository.Delete(permUsername)
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

// getFullUsername composes the key one user's permissions are stored under. A
// superuser gets its own prefix, so the same person holds one record however
// they sign in. How a key is built is nobody else's business: callers ask for
// what they want done, they do not name records.
func (s *Service) getFullUsername(prefix Prefix, username string) (string, error) {
	if username == "" {
		return "", ErrUsernameCannotBeEmpty
	}
	if prefix == "" {
		return "", ErrPrefixCannotBeEmpty
	}
	if slices.Contains(s.superusers, username) {
		prefix = PrefixSuperuser
	}
	return string(prefix) + ":" + username, nil
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

func (s *Service) getStatus(permUsername string) Status {
	split := strings.Split(permUsername, ":")
	username := split[1]
	if slices.Contains(s.superusers, username) {
		return GRANTED
	}
	return NOT_PERMITTED
}

func (s *Service) updateUserPermissions(permUsername string, featuresList []config.Feature, status Status) error {
	var err error
	for _, feature := range featuresList {
		errPerm := s.updateUserPermission(permUsername, feature, status)
		if errPerm != nil {
			err = errors.Join(err, fmt.Errorf("%s: %w", feature, errPerm))
		}
	}
	return err
}

func (s *Service) updateUserPermission(permUsername string, feature config.Feature, status Status) error {
	if permUsername == "" {
		return ErrUsernameCannotBeEmpty
	}
	prefix, username, found := strings.Cut(permUsername, ":")
	if !found || prefix == "" {
		return ErrPrefixCannotBeEmpty
	}
	if username == "" {
		return ErrUsernameCannotBeEmpty
	}
	if slices.Contains(s.superusers, username) {
		return ErrCannotChangePermissionsForSuperusers
	}
	if !s.isValidFeature(feature) {
		return ErrInvalidFeature
	}

	existingPermissions, err := s.permissionRepository.Get(permUsername)
	if err != nil {
		return err
	}

	existingPermissions[feature] = status
	return s.permissionRepository.CreateOrUpdate(permUsername, existingPermissions)
}

func (s *Service) normalizeDatabase() error {
	permissionsMap, errMap := s.permissionRepository.GetAll()
	if errMap != nil {
		return errMap
	}
	for permUsername, permissions := range permissionsMap {
		status := s.getStatus(permUsername)
		permissions = permissions.renamed()
		normalisedPermissions := make(PermissionMap)
		for _, feature := range config.All {
			if perm, ok := permissions[feature]; !ok {
				normalisedPermissions[feature] = status
			} else {
				normalisedPermissions[feature] = perm
			}
		}
		errUpdate := s.permissionRepository.CreateOrUpdate(permUsername, normalisedPermissions)
		if errUpdate != nil {
			return errUpdate
		}
	}
	return nil
}
