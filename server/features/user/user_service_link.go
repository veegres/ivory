package user

import (
	"errors"
	"ivory/clients/storage"
	"sort"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// A link is the one-time way a user sets their password: a token they hold plus
// a record Ivory keeps. The token carries the payload and the expiration, the
// record is what a link can be revoked and spent through. Every link is handed
// out by somebody signed in - the initial setup creates its superuser outright,
// so there is no public way to issue one. These methods carry a Link prefix
// because the same service owns the users themselves.

// LinkCreateInvite issues a link for a user Ivory does not have yet. Inviting a
// superuser is a superuser's own right, exactly as deleting one is.
func (s *Service) LinkCreateInvite(request LinkRequest, requester string) (*LinkCreatedResponse, error) {
	if request.Superuser {
		if err := s.requireSuperuser(requester); err != nil {
			return nil, err
		}
	}
	return s.linkCreate(strings.TrimSpace(request.Username), LinkInvite, request.Superuser)
}

// LinkCreateReset issues a link that hands an existing user a new password,
// which is how a forgotten one is answered: nobody ever types somebody else's
// password, and the account keeps its permissions. A reset takes the account
// over, so a superuser's is a superuser's own to hand out.
func (s *Service) LinkCreateReset(username string, requester string) (*LinkCreatedResponse, error) {
	name := strings.TrimSpace(username)
	if name == "" {
		return nil, ErrUsernameCannotBeEmpty
	}
	superuser, errSuperuser := s.IsSuperuser(name)
	if errSuperuser != nil {
		return nil, errSuperuser
	}
	if superuser {
		if err := s.requireSuperuser(requester); err != nil {
			return nil, err
		}
	}
	return s.linkCreate(name, LinkReset, superuser)
}

func (s *Service) LinkVerify(rawToken string) (*LinkPayload, error) {
	_, link, err := s.linkResolve(rawToken)
	if err != nil {
		return nil, err
	}
	payload := link.toPayload()
	return &payload, nil
}

// LinkApply spends the link - it creates the user or sets the password of the
// one it names - and drops the record, so the same token cannot be used twice.
func (s *Service) LinkApply(rawToken string, password string) (*Response, error) {
	id, link, errResolve := s.linkResolve(rawToken)
	if errResolve != nil {
		return nil, errResolve
	}
	user, errApply := s.linkApply(*link, password)
	if errApply != nil {
		return nil, errApply
	}
	if errDelete := s.userRepository.LinkDelete(id); errDelete != nil {
		return nil, errDelete
	}
	return user, nil
}

func (s *Service) LinkList() ([]LinkResponse, error) {
	linkMap, err := s.userRepository.LinkMap()
	if err != nil {
		return nil, err
	}
	result := make([]LinkResponse, 0, len(linkMap))
	for id, link := range linkMap {
		if time.Now().After(link.ExpiresAt) {
			if errDelete := s.userRepository.LinkDelete(id); errDelete != nil {
				return nil, errDelete
			}
			continue
		}
		result = append(result, link.toResponse(id))
	}
	sort.Slice(result, func(i, j int) bool { return result[i].CreatedAt.After(result[j].CreatedAt) })
	return result, nil
}

func (s *Service) LinkRevoke(id string, requester string) error {
	if id == "" {
		return ErrLinkIdCannotBeEmpty
	}
	link, err := s.userRepository.LinkGet(id)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return ErrLinkObsolete
		}
		return err
	}
	if link.Superuser {
		if errSuper := s.requireSuperuser(requester); errSuper != nil {
			return errSuper
		}
	}
	return s.userRepository.LinkDelete(id)
}

func (s *Service) LinkDeleteAll() error {
	return s.userRepository.LinkDeleteAll()
}

func (s *Service) linkCreate(username string, kind LinkKind, superuser bool) (*LinkCreatedResponse, error) {
	if username == "" {
		return nil, ErrUsernameCannotBeEmpty
	}
	exist, errExist := s.Exists(username)
	if errExist != nil {
		return nil, errExist
	}
	if kind == LinkReset && !exist {
		return nil, ErrUserNotFound
	}
	if kind == LinkInvite && exist {
		return nil, ErrUserAlreadyExists
	}

	id, errId := uuid.NewUUID()
	if errId != nil {
		return nil, errId
	}
	rawToken, exp, errToken := s.tokenService.Generate(username, jwt.MapClaims{"jti": id.String()}, s.linkExpiration)
	if errToken != nil {
		return nil, errToken
	}

	link := Link{Kind: kind, Username: username, Superuser: superuser, CreatedAt: time.Now(), ExpiresAt: *exp}
	stored, errCreate := s.userRepository.LinkCreate(id.String(), link)
	if errCreate != nil {
		return nil, errCreate
	}
	response := stored.toCreatedResponse(id.String(), rawToken)
	return &response, nil
}

func (s *Service) linkApply(link Link, password string) (*Response, error) {
	if link.Kind == LinkReset {
		return s.ResetPassword(link.Username, password, link.CreatedAt)
	}
	return s.Create(link.Username, password, link.Superuser)
}

// linkResolve answers both halves of a link: the token has to be signed by this
// Ivory and unexpired, and the record has to still be there.
func (s *Service) linkResolve(rawToken string) (string, *Link, error) {
	if rawToken == "" {
		return "", nil, ErrLinkInvalid
	}
	claims, errParse := s.tokenService.Parse(rawToken)
	if errParse != nil {
		if errors.Is(errParse, jwt.ErrTokenExpired) {
			return "", nil, ErrLinkExpired
		}
		return "", nil, ErrLinkInvalid
	}
	id, okId := claims["jti"].(string)
	if !okId || id == "" {
		return "", nil, ErrLinkInvalid
	}
	username, errUsername := claims.GetSubject()
	if errUsername != nil || username == "" {
		return "", nil, ErrLinkInvalid
	}

	link, errLink := s.userRepository.LinkGet(id)
	if errLink != nil {
		if errors.Is(errLink, storage.ErrNotFound) {
			return "", nil, ErrLinkObsolete
		}
		return "", nil, errLink
	}
	if link.Username != username {
		return "", nil, ErrLinkInvalid
	}
	if time.Now().After(link.ExpiresAt) {
		return "", nil, ErrLinkExpired
	}
	return id, &link, nil
}
