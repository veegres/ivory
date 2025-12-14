package basic

import (
	"errors"
	"ivory/src/clients/auth"
)

var ErrUsernameNotSpecified = errors.New("username is not specified")
var ErrPasswordNotSpecified = errors.New("password is not specified")
var ErrConfigNotConfigured = errors.New("config is not configured")
var ErrCredentialsNotCorrect = errors.New("credentials are not correct")
var ErrConnectionObsolete = errors.New("connection is obsolete")

// NOTE: validate that is matches interface in compile-time
var _ auth.Provider[Config, Login] = (*Provider)(nil)

type Provider struct {
	config *Config
}

func NewProvider() *Provider {
	return &Provider{}
}

func (p *Provider) Configured() bool {
	return p.config != nil
}

func (p *Provider) SetConfig(config Config) error {
	if config.Username == "" {
		return ErrUsernameNotSpecified
	}
	if config.Password == "" {
		return ErrPasswordNotSpecified
	}
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
	if subject.Username != p.config.Username || subject.Password != p.config.Password {
		return "", ErrCredentialsNotCorrect
	}
	return subject.Username, nil
}

func (p *Provider) Connect(_ Config) error {
	return ErrConnectionObsolete
}
