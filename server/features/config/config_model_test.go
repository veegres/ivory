package config

import (
	"ivory/clients/auth/ldap"
	"ivory/core/utils"
	"ivory/features/user"
	"strings"
	"testing"

	"github.com/gin-gonic/gin/binding"
)

func TestNewAppConfigBinding(t *testing.T) {
	utils.RegisterValidators()
	tests := []struct {
		name    string
		request NewAppConfig
		want    string
	}{
		{name: "a company and nothing else passes", request: NewAppConfig{AppConfig: AppConfig{Company: "acme"}}},
		{name: "empty company name is rejected", request: NewAppConfig{AppConfig: AppConfig{Company: ""}}, want: "'NewAppConfig.AppConfig.Company' Error:Field validation for 'Company' failed on the 'required' tag"},
		{name: "blank company name is rejected", request: NewAppConfig{AppConfig: AppConfig{Company: "  "}}, want: "'NewAppConfig.AppConfig.Company' Error:Field validation for 'Company' failed on the 'notblank' tag"},
		{
			name:    "an ldap config states where to bind",
			request: NewAppConfig{AppConfig: AppConfig{Company: "acme", Auth: AuthConfig{Ldap: &ldap.Config{Url: "ldap://x"}}}},
			want:    "'NewAppConfig.AppConfig.Auth.Ldap.BindDN' Error:Field validation for 'BindDN' failed on the 'required' tag",
		},
		{
			name:    "the setup user signs in some way",
			request: NewAppConfig{AppConfig: AppConfig{Company: "acme"}, User: &UserSetupRequest{Username: "admin"}},
			want:    "'NewAppConfig.User.AuthTypes' Error:Field validation for 'AuthTypes' failed on the 'required' tag",
		},
		{
			name:    "the setup user is not registered for a minted session",
			request: NewAppConfig{AppConfig: AppConfig{Company: "acme"}, User: &UserSetupRequest{Username: "admin", AuthTypes: []user.AuthType{user.AuthSession}}},
			want:    "'NewAppConfig.User.AuthTypes[0]' Error:Field validation for 'AuthTypes[0]' failed on the 'ne' tag",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := binding.Validator.ValidateStruct(&test.request)
			if test.want == "" {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				return
			}
			if err == nil || !strings.Contains(err.Error(), test.want) {
				t.Fatalf("expected an error containing %q, got %v", test.want, err)
			}
		})
	}
}
