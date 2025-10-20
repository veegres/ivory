package basic

import (
	"errors"
)

type BasicProvider struct {
	config *Config
}

func NewProvider() *BasicProvider {
	return &BasicProvider{}
}

func (p *BasicProvider) SetConfig(config *Config) error {
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

func (p *BasicProvider) DeleteConfig() {
	p.config = nil
}

func (p *BasicProvider) Verify(login Login) (string, error) {
	if login.Username != p.config.Username || login.Password != p.config.Password {
		return "", errors.New("credentials are not correct")
	}
	return login.Username, nil
}
