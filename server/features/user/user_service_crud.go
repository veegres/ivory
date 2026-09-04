package user

import (
	"errors"
	"ivory/clients/storage"
	"strings"
	"time"
)

func (s *Service) List() ([]Response, error) {
	users, err := s.userRepository.List()
	if err != nil {
		return nil, err
	}
	result := make([]Response, 0, len(users))
	for _, u := range users {
		result = append(result, u.toResponse())
	}
	return result, nil
}

func (s *Service) Create(username string, password string, superuser bool) (*Response, error) {
	name := strings.TrimSpace(username)
	if name == "" {
		return nil, ErrUsernameCannotBeEmpty
	}
	if password == "" {
		return nil, ErrPasswordCannotBeEmpty
	}
	exist, errExist := s.Exists(name)
	if errExist != nil {
		return nil, errExist
	}
	if exist {
		return nil, ErrUserAlreadyExists
	}
	encrypted, errEnc := s.encryptionService.Encrypt(password, s.secretService.Get())
	if errEnc != nil {
		return nil, errEnc
	}
	user, errCreate := s.userRepository.Create(User{
		Username:  name,
		Password:  encrypted,
		Superuser: superuser,
		CreatedAt: time.Now(),
	})
	if errCreate != nil {
		return nil, errCreate
	}
	if errSync := s.syncSuperusers(); errSync != nil {
		return nil, errSync
	}
	response := user.toResponse()
	return &response, nil
}

// Delete removes a user together with the permissions they were granted.
// Nobody deletes themselves, only a superuser deletes a superuser, and the last
// superuser stays whatever happens - an Ivory nobody can administer is no use
// to its owner.
func (s *Service) Delete(username string, requester string) error {
	if username == "" {
		return ErrUsernameCannotBeEmpty
	}
	if requester != "" && username == requester {
		return ErrCannotDeleteYourself
	}
	user, errGet := s.userRepository.Get(username)
	if errGet != nil {
		if errors.Is(errGet, storage.ErrNotFound) {
			return ErrUserNotFound
		}
		return errGet
	}
	if user.Superuser {
		if errSuper := s.requireSuperuser(requester); errSuper != nil {
			return errSuper
		}
		// NOTE: the floor is checked here as well as inside the delete's own
		// transaction, so a refusal costs nothing - a delete that was never
		// allowed must not have taken the permissions with it on the way
		superusers, errSuper := s.userRepository.Superusers()
		if errSuper != nil {
			return errSuper
		}
		if len(superusers) <= 1 {
			return ErrLastSuperuser
		}
	}

	// NOTE: the permissions go first and the two steps are deliberately not one
	// transaction - they live in different buckets, and the order is what makes
	// that safe. Permissions can be granted again by anybody who holds
	// manage.permission.update, while a half-deleted user who kept their access
	// is the outcome nobody can undo, so the step that would leave access
	// standing is the one that happens last.
	// The permission feature works out which record that is: a superuser's is
	// kept under a prefix of its own, and it is still one while this runs.
	if errPerm := s.permissionService.DeleteBasicUserPermissions(username); errPerm != nil {
		return errPerm
	}
	// NOTE: the last-superuser rule is decided inside the delete's own
	// transaction, so two deletes at once cannot both find somebody else left
	if errDelete := s.userRepository.DeleteIf(username, s.checkNotLastSuperuser(username)); errDelete != nil {
		return errDelete
	}
	return s.syncSuperusers()
}

// checkNotLastSuperuser is the rule the delete transaction decides on: Ivory
// keeps at least one superuser, whoever is asking and whatever else is being
// deleted at the same moment.
func (s *Service) checkNotLastSuperuser(username string) func(users map[string]User) error {
	return func(users map[string]User) error {
		user, ok := users[username]
		if !ok {
			return ErrUserNotFound
		}
		if !user.Superuser {
			return nil
		}
		superusers := 0
		for _, u := range users {
			if u.Superuser {
				superusers++
			}
		}
		if superusers <= 1 {
			return ErrLastSuperuser
		}
		return nil
	}
}

// ResetPassword sets a password without asking for the previous one, which only
// a one-time link may do: the link is the proof, and it is spent on the way. The
// user's permissions and superuser flag stay exactly as they were. issuedAt is
// the link's own timestamp - a name deleted and invited again is a different
// account, and a link written for the old one must not open the new one.
func (s *Service) ResetPassword(username string, newPassword string, issuedAt time.Time) (*Response, error) {
	if newPassword == "" {
		return nil, ErrPasswordCannotBeEmpty
	}
	user, errGet := s.userRepository.Get(username)
	if errGet != nil {
		if errors.Is(errGet, storage.ErrNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, errGet
	}
	if user.CreatedAt.After(issuedAt) {
		return nil, ErrLinkObsolete
	}
	encrypted, errEnc := s.encryptionService.Encrypt(newPassword, s.secretService.Get())
	if errEnc != nil {
		return nil, errEnc
	}
	user.Password = encrypted
	if errUpdate := s.userRepository.Update(user); errUpdate != nil {
		return nil, errUpdate
	}
	response := user.toResponse()
	return &response, nil
}

// UpdatePassword is the only way an existing user changes, and only the owner
// can do it: the password they are replacing is the proof of that. A username
// is never updated - a user is created and deleted, nothing else.
func (s *Service) UpdatePassword(username string, previousPassword string, newPassword string) error {
	if username == "" {
		return ErrUsernameCannotBeEmpty
	}
	if newPassword == "" {
		return ErrPasswordCannotBeEmpty
	}
	if errVerify := s.VerifyPassword(username, previousPassword); errVerify != nil {
		return errVerify
	}
	user, errGet := s.userRepository.Get(username)
	if errGet != nil {
		return errGet
	}
	encrypted, errEnc := s.encryptionService.Encrypt(newPassword, s.secretService.Get())
	if errEnc != nil {
		return errEnc
	}
	user.Password = encrypted
	return s.userRepository.Update(user)
}

func (s *Service) DeleteAll() error {
	if err := s.userRepository.DeleteAll(); err != nil {
		return err
	}
	return s.syncSuperusers()
}

// VerifyPassword implements the basic provider's Verifier: it decrypts the
// stored password and compares it to the one that was typed.
func (s *Service) VerifyPassword(username string, password string) error {
	user, errGet := s.userRepository.Get(username)
	if errGet != nil {
		if errors.Is(errGet, storage.ErrNotFound) {
			return ErrCredentialsNotCorrect
		}
		return errGet
	}
	decrypted, errDec := s.encryptionService.Decrypt(user.Password, s.secretService.Get())
	if errDec != nil {
		return errDec
	}
	if decrypted != password {
		return ErrCredentialsNotCorrect
	}
	return nil
}

func (s *Service) Exists(username string) (bool, error) {
	_, err := s.userRepository.Get(username)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (s *Service) IsSuperuser(username string) (bool, error) {
	user, err := s.userRepository.Get(username)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return false, nil
		}
		return false, err
	}
	return user.Superuser, nil
}

// requireSuperuser answers whether the person asking may manage a superuser. An
// empty requester is Ivory running without authentication, where there is no
// identity to check and everything is permitted anyway.
func (s *Service) requireSuperuser(requester string) error {
	if requester == "" {
		return nil
	}
	isSuperuser, err := s.IsSuperuser(requester)
	if err != nil {
		return err
	}
	if !isSuperuser {
		return ErrSuperuserRequired
	}
	return nil
}

func (s *Service) HasSuperuser() (bool, error) {
	superusers, err := s.userRepository.Superusers()
	if err != nil {
		return false, err
	}
	return len(superusers) > 0, nil
}

func (s *Service) Reencrypt(oldSecret [16]byte, newSecret [16]byte) error {
	users, errList := s.userRepository.List()
	if errList != nil {
		return errList
	}
	for _, user := range users {
		decrypted, errDec := s.encryptionService.Decrypt(user.Password, oldSecret)
		if errDec != nil {
			return errDec
		}
		encrypted, errEnc := s.encryptionService.Encrypt(decrypted, newSecret)
		if errEnc != nil {
			return errEnc
		}
		user.Password = encrypted
		if errUpdate := s.userRepository.Update(user); errUpdate != nil {
			return errUpdate
		}
	}
	return nil
}

func (s *Service) syncSuperusers() error {
	superusers, err := s.userRepository.Superusers()
	if err != nil {
		return err
	}
	return s.permissionService.SetSuperusers(superusers)
}
