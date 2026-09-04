package basic

import (
	"errors"
	"ivory/clients/auth"
)

var ErrConfigNotConfigured = errors.New("config is not configured")
var ErrUsernameNotSpecified = errors.New("username is not specified")
var ErrPasswordNotSpecified = errors.New("password is not specified")
var ErrConnectionObsolete = errors.New("connection is obsolete")

// NOTE: validate that is matches interface in compile-time
var _ auth.Provider[Config, Login] = (*Provider)(nil)

// Verifier checks a login against the store that owns the credentials, which is
// the user feature. The provider only knows the narrow view it uses.
type Verifier interface {
	VerifyPassword(username string, password string) error
}

type Provider struct {
	config   *Config
	verifier Verifier
}

func NewProvider(verifier Verifier) *Provider {
	return &Provider{verifier: verifier}
}

func (p *Provider) Configured() bool {
	return p.config != nil
}

func (p *Provider) SetConfig(config Config) error {
	p.config = &config
	return nil
}

func (p *Provider) DeleteConfig() {
	p.config = nil
}

func (p *Provider) Verify(subject Login) (string, error) {
	if p.config == nil {
		return "", ErrConfigNotConfigured
	}
	if subject.Username == "" {
		return "", ErrUsernameNotSpecified
	}
	if subject.Password == "" {
		return "", ErrPasswordNotSpecified
	}
	if err := p.verifier.VerifyPassword(subject.Username, subject.Password); err != nil {
		return "", err
	}
	return subject.Username, nil
}

func (p *Provider) Connect(_ Config) error {
	return ErrConnectionObsolete
}
