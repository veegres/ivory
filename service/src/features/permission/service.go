package permission

import (
	"errors"
	"ivory/src/storage/db"
	"slices"
)

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

	return result, nil
}

func (s *Service) GetUserPermissions(username string) (map[string]PermissionStatus, error) {
	if username == "" {
		return nil, errors.New("username cannot be empty")
	}
	permissions, err := s.permissionRepository.Get(username)
	if err != nil {
		return nil, err
	}
	return permissions, nil
}

func (s *Service) CreateUserPermissions(username string) (map[string]PermissionStatus, error) {
	if username == "" {
		return nil, errors.New("username cannot be empty")
	}
	existingPermissions, err := s.permissionRepository.Get(username)
	if err != nil && !errors.Is(err, db.ErrNotFound) {
		return nil, err
	}
	if existingPermissions != nil {
		return existingPermissions, nil
	}

	status := NOT_PERMITTED
	if slices.Contains(s.superusers, username) {
		status = GRANTED
	}

	permissions := make(map[string]PermissionStatus)
	for _, permissionName := range PermissionList {
		permissions[permissionName] = status
	}

	errCreate := s.permissionRepository.CreateOrUpdate(username, permissions)
	if errCreate != nil {
		return nil, errCreate
	}

	return permissions, nil
}

func (s *Service) RequestUserPermission(username string, permissionName string) error {
	return s.updateUserPermission(username, permissionName, PENDING)
}

func (s *Service) ApproveUserPermission(username string, permissionName string) error {
	return s.updateUserPermission(username, permissionName, GRANTED)
}

func (s *Service) RejectUserPermission(username string, permissionName string) error {
	return s.updateUserPermission(username, permissionName, NOT_PERMITTED)
}

func (s *Service) updateUserPermission(username string, permissionName string, status PermissionStatus) error {
	if username == "" {
		return errors.New("username cannot be empty")
	}
	if slices.Contains(s.superusers, username) {
		return errors.New("cannot change permissions for superusers")
	}
	if !s.isValidPermission(permissionName) {
		return errors.New("invalid permission name: " + permissionName)
	}

	existingPermissions, err := s.permissionRepository.Get(username)
	if err != nil {
		return err
	}

	existingPermissions[permissionName] = status
	return s.permissionRepository.CreateOrUpdate(username, existingPermissions)
}

func (s *Service) DeleteUserPermissions(username string) error {
	if username == "" {
		return errors.New("username cannot be empty")
	}
	return s.permissionRepository.Delete(username)
}

func (s *Service) DeleteAll() error {
	return s.permissionRepository.DeleteAll()
}

func (s *Service) isValidPermission(permissionName string) bool {
	for _, validPerm := range PermissionList {
		if validPerm == permissionName {
			return true
		}
	}
	return false
}

func (s *Service) normalize() error {
	permissionsMap, errMap := s.permissionRepository.GetAll()
	if errMap != nil {
		return errMap
	}
	for username, permissions := range permissionsMap {
		status := NOT_PERMITTED
		if slices.Contains(s.superusers, username) {
			status = GRANTED
		}

		normalisedPermissions := make(map[string]PermissionStatus)
		for _, permissionName := range PermissionList {
			if perm, ok := permissions[permissionName]; !ok {
				normalisedPermissions[permissionName] = status
			} else {
				normalisedPermissions[permissionName] = perm
			}
		}

		errUpdate := s.permissionRepository.CreateOrUpdate(username, normalisedPermissions)
		if errUpdate != nil {
			return errUpdate
		}
	}
	return nil
}
