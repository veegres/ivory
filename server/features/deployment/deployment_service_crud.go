package deployment

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

// List returns the user's own templates followed by the shipped defaults, so
// the UI renders one list rather than two separate concepts.
func (s *Service) List(criteria ListRequest) ([]Template, error) {
	stored, err := s.repository.List(criteria)
	if err != nil {
		return nil, err
	}
	return append(stored, s.Defaults(criteria)...), nil
}

func (s *Service) Get(key uuid.UUID) (Template, error) {
	return s.repository.Get(key)
}

func (s *Service) Create(r TemplateRequest) (*Template, error) {
	if err := s.validate(r, ""); err != nil {
		return nil, err
	}
	id := uuid.New()
	return s.repository.Create(id, r.toTemplate(id, time.Now().UnixNano()))
}

// Update rejects a change of keeper or platform: the commands are written
// against a specific pair, so switching either would silently invalidate them.
func (s *Service) Update(key uuid.UUID, r TemplateRequest) (*Template, error) {
	current, err := s.repository.Get(key)
	if err != nil {
		return nil, err
	}
	if r.Keeper != current.Keeper || r.Platform != current.Platform {
		return nil, ErrTemplatePluginImmutable
	}
	if err := s.validate(r, current.Name); err != nil {
		return nil, err
	}
	return s.repository.Update(key, r.toTemplate(key, current.CreatedAt))
}

func (s *Service) Delete(key uuid.UUID) error {
	if _, err := s.repository.Get(key); err != nil {
		return err
	}
	return s.repository.Delete(key)
}

func (s *Service) DeleteAll() error {
	return s.repository.DeleteAll()
}

// validate enforces the rules a template must satisfy to be deployable at all:
// a unique name, at least one command, and only variables from the closed
// vocabulary. allowedName is the template's own current name, which a rename
// check must not treat as a collision with itself.
func (s *Service) validate(r TemplateRequest, allowedName string) error {
	if strings.TrimSpace(r.Name) == "" {
		return ErrTemplateNameRequired
	}
	if len(r.Commands) == 0 {
		return ErrTemplateCommandsRequired
	}
	if unknown := r.unknownPlaceholders(); len(unknown) > 0 {
		return fmt.Errorf("unknown variables: %s", strings.Join(unknown, ", "))
	}
	if r.trimmed().Name == allowedName {
		return nil
	}
	return s.validateNameFree(r)
}

// validateNameFree checks the name against stored templates and the shipped
// defaults alike - a copy that reuses a default's name would be ambiguous in
// the one list they share. It looks only within the template's own
// keeper/platform pair, which is the list the user picked the name in: a
// collision they cannot see is not one they can resolve.
func (s *Service) validateNameFree(r TemplateRequest) error {
	name := r.trimmed().Name
	existing, err := s.repository.GetByName(name, r.Keeper, r.Platform)
	if err != nil {
		return err
	}
	if existing != nil {
		return ErrTemplateNameTaken
	}
	for _, d := range s.Defaults(ListRequest{Keeper: &r.Keeper, Platform: &r.Platform}) {
		if d.Name == name {
			return ErrTemplateNameTaken
		}
	}
	return nil
}
