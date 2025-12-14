package oidc

import (
	"context"
	"errors"
	"ivory/src/clients/auth"

	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
)

var ErrIssuerURLNotSpecified = errors.New("IssuerURL is not specified")
var ErrClientIDNotSpecified = errors.New("ClientID is not specified")
var ErrClientSecretNotSpecified = errors.New("ClientSecret is not specified")
var ErrRedirectURLNotSpecified = errors.New("RedirectURL is not specified")
var ErrFailedToExchangeToken = errors.New("failed to exchange token")
var ErrNoIDTokenField = errors.New("no id_token field in oauth2 token")
var ErrFailedToVerifyIDToken = errors.New("failed to verify ID Token")
var ErrNoRightClaim = errors.New("there is no right claim in ID Token (it should have any of these preferred_username, email, nickname)")
var ErrConfigNotConfigured = errors.New("config is not configured")

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

func (p *Provider) Configured() bool {
	return p.config != nil
}

func (p *Provider) SetConfig(config Config) error {
	if config.IssuerURL == "" {
		return ErrIssuerURLNotSpecified
	}
	if config.ClientID == "" {
		return ErrClientIDNotSpecified
	}
	if config.ClientSecret == "" {
		return ErrClientSecretNotSpecified
	}
	if config.RedirectURL == "" {
		return ErrRedirectURLNotSpecified
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
	err := p.initialize()
	if err != nil {
		return "", err
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
		return "", errors.Join(ErrFailedToExchangeToken, errExchange)
	}

	rawIDToken, ok := oauthToken.Extra("id_token").(string)
	if !ok {
		return "", ErrNoIDTokenField
	}

	idToken, errVerify := p.oauthVerifier.Verify(context.Background(), rawIDToken)
	if errVerify != nil {
		return "", errors.Join(ErrFailedToVerifyIDToken, errVerify)
	}

	// NOTE: standard for scope and claims can be found here https://openid.net/specs/openid-connect-core-1_0.html#ScopeClaims
	var claims struct {
		Username string `json:"preferred_username"`
		Nickname string `json:"nickname"`
		Email    string `json:"email"`
	}
	if err := idToken.Claims(&claims); err != nil {
		return "", err
	}
	if claims.Username != "" {
		return claims.Username, nil
	}
	if claims.Email != "" {
		return claims.Email, nil
	}
	if claims.Nickname != "" {
		return claims.Nickname, nil
	}
	return "", ErrNoRightClaim
}

func (p *Provider) Connect(config Config) error {
	// NOTE: just check if connection is possible, it doesn't check the validity of ClientID and ClientSecret
	//  we cannot verify it for web, but can machine-to-machine, but it's not a case for us
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
		return ErrConfigNotConfigured
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
