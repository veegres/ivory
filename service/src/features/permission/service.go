package permission

import (
	"errors"
	"ivory/src/storage/db"
	"slices"
	"sort"
	"strings"
)

var ErrAtLeastOneSuperuser = errors.New("there should be at least 1 superuser")
var ErrSuperusersCannotHaveEmptyName = errors.New("superusers cannot have empty name")
var ErrUsernameCannotBeEmpty = errors.New("username cannot be empty")
var ErrPrefixCannotBeEmpty = errors.New("prefix cannot be empty")
var ErrCannotChangePermissionsForSuperusers = errors.New("cannot change permissions for superusers")

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

func (s *Service) SetSuperusers(superusers []string) error {
	if len(superusers) == 0 {
		return ErrAtLeastOneSuperuser
	}
	if slices.Contains(superusers, "") {
		return ErrSuperusersCannotHaveEmptyName
	}
	s.superusers = superusers
	return s.normalize()
}

func (s *Service) DeleteAdmins() {
	s.superusers = []string{}
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

func (s *Service) GetUserPermissions(prefix string, username string, allowAll bool) (map[string]PermissionStatus, error) {
	if allowAll {
		return s.getAllPermissionsWithStatus(GRANTED), nil
	}
	permUsername, errName := s.getFullUsername(prefix, username)
	if errName != nil {
		return nil, errName
	}
	permissions, err := s.permissionRepository.Get(permUsername)
	if err != nil {
		return nil, err
	}
	return permissions, nil
}

func (s *Service) CreateUserPermissions(prefix string, username string) (map[string]PermissionStatus, error) {
	permUsername, errName := s.getFullUsername(prefix, username)
	if errName != nil {
		return nil, errName
	}
	existingPermissions, err := s.permissionRepository.Get(permUsername)
	if err != nil && !errors.Is(err, db.ErrNotFound) {
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

func (s *Service) RequestUserPermission(prefix string, username string, permissionName string) error {
	permUsername, errName := s.getFullUsername(prefix, username)
	if errName != nil {
		return errName
	}
	return s.updateUserPermission(permUsername, permissionName, PENDING)
}

func (s *Service) ApproveUserPermission(permUsername string, permissionName string) error {
	return s.updateUserPermission(permUsername, permissionName, GRANTED)
}

func (s *Service) RejectUserPermission(permUsername string, permissionName string) error {
	return s.updateUserPermission(permUsername, permissionName, NOT_PERMITTED)
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

func (s *Service) getFullUsername(prefix string, username string) (string, error) {
	if username == "" {
		return "", ErrUsernameCannotBeEmpty
	}
	if prefix == "" {
		return "", ErrPrefixCannotBeEmpty
	}
	if slices.Contains(s.superusers, username) {
		prefix = "superuser"
	}
	return prefix + ":" + username, nil
}

func (s *Service) isValidPermission(permissionName string) bool {
	for _, validPerm := range PermissionList {
		if validPerm == permissionName {
			return true
		}
	}
	return false
}

func (s *Service) getAllPermissionsWithStatus(status PermissionStatus) map[string]PermissionStatus {
	permissions := make(map[string]PermissionStatus)
	for _, permissionName := range PermissionList {
		permissions[permissionName] = status
	}
	return permissions
}

func (s *Service) getStatus(permUsername string) PermissionStatus {
	split := strings.Split(permUsername, ":")
	username := split[1]
	if slices.Contains(s.superusers, username) {
		return GRANTED
	}
	return NOT_PERMITTED
}

func (s *Service) updateUserPermission(permUsername string, permissionName string, status PermissionStatus) error {
	if permUsername == "" {
		return ErrUsernameCannotBeEmpty
	}
	split := strings.Split(permUsername, ":")
	username := split[1]
	if slices.Contains(s.superusers, username) {
		return ErrCannotChangePermissionsForSuperusers
	}
	if !s.isValidPermission(permissionName) {
		return errors.New("invalid permission name: " + permissionName)
	}

	existingPermissions, err := s.permissionRepository.Get(permUsername)
	if err != nil {
		return err
	}

	existingPermissions[permissionName] = status
	return s.permissionRepository.CreateOrUpdate(permUsername, existingPermissions)
}

func (s *Service) normalize() error {
	permissionsMap, errMap := s.permissionRepository.GetAll()
	if errMap != nil {
		return errMap
	}
	for permUsername, permissions := range permissionsMap {
		status := s.getStatus(permUsername)
		normalisedPermissions := make(map[string]PermissionStatus)
		for _, permissionName := range PermissionList {
			if perm, ok := permissions[permissionName]; !ok {
				normalisedPermissions[permissionName] = status
			} else {
				normalisedPermissions[permissionName] = perm
			}
		}

		errUpdate := s.permissionRepository.CreateOrUpdate(permUsername, normalisedPermissions)
		if errUpdate != nil {
			return errUpdate
		}
	}
	return nil
}
