package user

import (
	"errors"
	"ivory/clients/storage"
	"ivory/features/permission"
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

// Create brings a user into existence by name and by the ways they may sign
// in. Nobody types their password here: asking for basic issues a registration
// along with the user, and the person sets it themselves on the page it opens.
func (s *Service) Create(request CreateRequest, requester string) (*CreateResponse, error) {
	if request.Superuser {
		if err := s.requireSuperuser(requester); err != nil {
			return nil, err
		}
	}
	user, errCreate := s.create(request.Username, "", request.AuthTypes, request.Superuser)
	if errCreate != nil {
		return nil, errCreate
	}
	if !user.allows(AuthBasic) {
		return &CreateResponse{User: user.toResponse()}, nil
	}
	registration, registered, errRegistration := s.registrationIssue(*user)
	if errRegistration != nil {
		return nil, errRegistration
	}
	return &CreateResponse{User: registered.toResponse(), Registration: registration}, nil
}

// CreateOutright creates a user without handing a registration out, which
// two callers need and nobody else may have: setup, where a password is typed
// on somebody else's behalf because at first run there is nobody signed in to
// send a registration to, and a restore, which brings a person back with no
// password at all for the admin to register afresh.
func (s *Service) CreateOutright(username string, password string, authTypes []AuthType, superuser bool) (*Response, error) {
	user, err := s.create(username, password, authTypes, superuser)
	if err != nil {
		return nil, err
	}
	response := user.toResponse()
	return &response, nil
}

// Update changes the ways an existing user may sign in, and nothing else - not
// their name, and never their superuser flag. Asking for basic issues a
// registration; taking it away drops the password along with it, since there is
// nothing left it could be used for.
func (s *Service) Update(username string, authTypes []AuthType, requester string) (*CreateResponse, error) {
	user, errGet := s.get(username)
	if errGet != nil {
		return nil, errGet
	}
	if user.Superuser {
		if errSuper := s.requireSuperuser(requester); errSuper != nil {
			return nil, errSuper
		}
	}
	if errAuth := s.validateAuthTypes(authTypes); errAuth != nil {
		return nil, errAuth
	}

	basicBefore := user.allows(AuthBasic)
	user.AuthTypes = authTypes
	basicAfter := user.allows(AuthBasic)
	if basicBefore && !basicAfter {
		user.Password = ""
		user.RegistrationId = ""
		user.RegistrationExpiresAt = time.Time{}
	}
	if errUpdate := s.userRepository.Update(*user); errUpdate != nil {
		return nil, errUpdate
	}
	if !basicBefore && basicAfter {
		registration, registered, errRegistration := s.registrationIssue(*user)
		if errRegistration != nil {
			return nil, errRegistration
		}
		return &CreateResponse{User: registered.toResponse(), Registration: registration}, nil
	}
	return &CreateResponse{User: user.toResponse()}, nil
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
	user, errGet := s.get(username)
	if errGet != nil {
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
	if errPerm := s.permissionService.DeleteUserPermissions(username); errPerm != nil {
		return errPerm
	}
	// NOTE: the last-superuser rule is decided inside the delete's own
	// transaction, so two deletes at once cannot both find somebody else left
	return s.userRepository.DeleteIf(username, s.checkNotLastSuperuser(username))
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

// UpdatePassword is the only way an existing user changes their own record, and
// only its owner can do it: the password they are replacing is the proof of
// that. A username is never updated - a user is registered and deleted.
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
	user, errGet := s.get(username)
	if errGet != nil {
		return errGet
	}
	encrypted, errEnc := s.encrypt(newPassword)
	if errEnc != nil {
		return errEnc
	}
	user.Password = encrypted
	return s.userRepository.Update(*user)
}

func (s *Service) DeleteAll() error {
	return s.userRepository.DeleteAll()
}

// VerifyPassword implements the basic provider's Verifier: it decrypts the
// stored password and compares it to the one that was typed. A user who was
// never registered for basic, or who has not set a password yet, has nothing to
// compare against.
func (s *Service) VerifyPassword(username string, password string) error {
	user, errGet := s.userRepository.Get(username)
	if errGet != nil {
		if errors.Is(errGet, storage.ErrNotFound) {
			return ErrCredentialsNotCorrect
		}
		return errGet
	}
	if !user.allows(AuthBasic) || user.Password == "" {
		return ErrCredentialsNotCorrect
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

// VerifySignIn is the gate every login goes through once its provider has
// vouched for somebody: a person signs in only under a name Ivory was told
// about, and only the ways that registration said they may.
func (s *Service) VerifySignIn(username string, authType AuthType) error {
	if username == "" {
		return ErrUsernameCannotBeEmpty
	}
	user, errGet := s.userRepository.Get(username)
	if errGet != nil {
		if errors.Is(errGet, storage.ErrNotFound) {
			return ErrUserNotRegistered
		}
		return errGet
	}
	if !user.allows(authType) {
		return ErrAuthTypeNotAllowed
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
		if user.Password == "" {
			continue
		}
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

// create is the one place a user record comes into existence, whether the
// password is set now (setup) or later (a registration).
func (s *Service) create(username string, password string, authTypes []AuthType, superuser bool) (*User, error) {
	name := strings.TrimSpace(username)
	if name == "" {
		return nil, ErrUsernameCannotBeEmpty
	}
	if errAuth := s.validateAuthTypes(authTypes); errAuth != nil {
		return nil, errAuth
	}
	exist, errExist := s.Exists(name)
	if errExist != nil {
		return nil, errExist
	}
	if exist {
		return nil, ErrUserAlreadyExists
	}
	encrypted := ""
	if password != "" {
		value, errEnc := s.encrypt(password)
		if errEnc != nil {
			return nil, errEnc
		}
		encrypted = value
	}
	user, errCreate := s.userRepository.Create(User{
		Username:  name,
		Password:  encrypted,
		AuthTypes: authTypes,
		Superuser: superuser,
		CreatedAt: time.Now(),
	})
	if errCreate != nil {
		return nil, errCreate
	}
	if _, errPerm := s.permissionService.CreateUserPermissions(user.Username, s.statusFor(user)); errPerm != nil {
		return nil, errPerm
	}
	return &user, nil
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

func (s *Service) validateAuthTypes(authTypes []AuthType) error {
	if len(authTypes) == 0 {
		return ErrAuthTypeRequired
	}
	for _, authType := range authTypes {
		if !authType.valid() {
			return ErrAuthTypeInvalid
		}
	}
	return nil
}

func (s *Service) get(username string) (*User, error) {
	user, err := s.userRepository.Get(username)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (s *Service) encrypt(password string) (string, error) {
	return s.encryptionService.Encrypt(password, s.secretService.Get())
}

// syncPermissions gives every user the record they are entitled to: one that is
// missing is created, and a feature added since the last start is filled in
// with the status that user's kind gets. The permission feature is told what to
// default to rather than working it out - it has never heard of a superuser.
func (s *Service) syncPermissions() error {
	users, err := s.userRepository.List()
	if err != nil {
		return err
	}
	for _, user := range users {
		if _, errPerm := s.permissionService.CreateUserPermissions(user.Username, s.statusFor(user)); errPerm != nil {
			return errPerm
		}
	}
	return nil
}

func (s *Service) statusFor(user User) permission.Status {
	if user.Superuser {
		return permission.GRANTED
	}
	return permission.NOT_PERMITTED
}
