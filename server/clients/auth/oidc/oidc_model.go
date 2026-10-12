package oidc

type Config struct {
	IssuerURL    string `json:"issuerUrl" binding:"required,notblank"`
	ClientID     string `json:"clientId" binding:"required,notblank"`
	ClientSecret string `json:"clientSecret" binding:"required"`
	RedirectURL  string `json:"redirectUrl" binding:"required,notblank"`
}
