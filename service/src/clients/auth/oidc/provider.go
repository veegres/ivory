package oidc

import (
	"context"
	"errors"
	"ivory/src/clients/auth"

	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
)

// NOTE: validate that is matches interface in compile-time
var _ auth.Provider[Config, string] = (*Provider)(nil)

type Provider struct {
	config        *Config
	oauthConfig   *oauth2.Config
	oauthVerifier *oidc.IDTokenVerifier
}

func NewProvider() *Provider {
	return &Provider{}
}

func (p *Provider) SetConfig(config Config) error {
	if config.IssuerURL == "" {
		return errors.New("IssuerURL is not specified")
	}
	if config.ClientID == "" {
		return errors.New("ClientID is not specified")
	}
	if config.ClientSecret == "" {
		return errors.New("ClientSecret is not specified")
	}
	if config.RedirectURL == "" {
		return errors.New("RedirectURL is not specified")
	}
	p.config = &config
	return nil
}

func (p *Provider) DeleteConfig() {
	p.config = nil
	p.oauthConfig = nil
	p.oauthVerifier = nil
}

func (p *Provider) GetCode(state string) (string, error) {
	if p.oauthConfig == nil {
		return "", errors.New("config is not configured")
	}
	return p.oauthConfig.AuthCodeURL(state), nil
}

func (p *Provider) Verify(subject string) (string, error) {
	err := p.initialize()
	if err != nil {
		return "", err
	}
	oauthToken, errExchange := p.oauthConfig.Exchange(context.Background(), subject)
	if errExchange != nil {
		return "", errors.Join(errors.New("failed to exchange token"), errExchange)
	}

	rawIDToken, ok := oauthToken.Extra("id_token").(string)
	if !ok {
		return "", errors.New("no id_token field in oauth2 token")
	}

	idToken, errVerify := p.oauthVerifier.Verify(context.Background(), rawIDToken)
	if errVerify != nil {
		return "", errors.Join(errors.New("failed to verify ID Token"), errVerify)
	}

	var claims struct {
		Email string `json:"email"`
	}
	if err := idToken.Claims(&claims); err != nil {
		return "", err
	}
	return claims.Email, nil
}

func (p *Provider) Connect(config Config) error {
	_, err := p.getConnection(config)
	return err
}

func (p *Provider) getConnection(config Config) (*oidc.Provider, error) {
	provider, errProvider := oidc.NewProvider(context.Background(), config.IssuerURL)
	if errProvider != nil {
		return nil, errProvider
	}
	return provider, nil
}

func (p *Provider) initialize() error {
	if p.config == nil {
		return errors.New("config is not configured")
	}
	provider, errProvider := p.getConnection(*p.config)
	if errProvider != nil {
		return errProvider
	}

	oauth2Config := oauth2.Config{
		ClientID:     p.config.ClientID,
		ClientSecret: p.config.ClientSecret,
		RedirectURL:  p.config.RedirectURL,
		Endpoint:     provider.Endpoint(),
		Scopes:       []string{oidc.ScopeOpenID, "profile", "email"},
	}

	p.oauthVerifier = provider.Verifier(&oidc.Config{ClientID: p.config.ClientID})
	p.oauthConfig = &oauth2Config
	return nil
}
