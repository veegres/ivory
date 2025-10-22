package basic

import (
	"errors"
	"ivory/src/clients/auth"
)

// NOTE: validate that is matches interface in compile-time
var _ auth.Provider[*Config, Login] = (*Provider)(nil)

type Provider struct {
	config *Config
}

func NewProvider() *Provider {
	return &Provider{}
}

func (p *Provider) SetConfig(config *Config) error {
	if config == nil {
		return errors.New("config is not configured")
	}
	if config.Username == "" {
		return errors.New("username is not specified")
	}
	if config.Password == "" {
		return errors.New("password is not specified")
	}
	p.config = config
	return nil
}

func (p *Provider) DeleteConfig() {
	p.config = nil
}

func (p *Provider) Verify(subject Login) (string, error) {
	if subject.Username != p.config.Username || subject.Password != p.config.Password {
		return "", errors.New("credentials are not correct")
	}
	return subject.Username, nil
}

func (p *Provider) Connect(_ *Config) error {
	return errors.New("connection is obsolete")
}
