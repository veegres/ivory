package oidc

import (
	"context"
	"errors"

	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
)

type OidcProvider struct {
	config        *Config
	oauthConfig   *oauth2.Config
	oauthVerifier *oidc.IDTokenVerifier
}

func NewProvider() *OidcProvider {
	return &OidcProvider{}
}

func (p *OidcProvider) SetConfig(config *Config) error {
	if config == nil {
		return errors.New("config is not configured")
	}
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
	p.config = config
	err := p.setOAuthConfig()
	if err != nil {
		return err
	}
	return nil
}

func (p *OidcProvider) DeleteConfig() {
	p.config = nil
	p.oauthConfig = nil
	p.oauthVerifier = nil
}

func (p *OidcProvider) GetCode(state string) (string, error) {
	if p.oauthConfig == nil {
		return "", errors.New("config is not configured")
	}
	return p.oauthConfig.AuthCodeURL(state), nil
}

func (p *OidcProvider) Verify(code string) (string, error) {
	oauthToken, errExchange := p.oauthConfig.Exchange(context.Background(), code)
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

func (p *OidcProvider) setOAuthConfig() error {
	if p.config == nil {
		return errors.New("config is not configured")
	}

	provider, errProvider := oidc.NewProvider(context.Background(), p.config.IssuerURL)
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
