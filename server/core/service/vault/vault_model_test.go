package vault

import (
	"ivory/core/utils"
	"strings"
	"testing"

	"github.com/gin-gonic/gin/binding"
)

func TestVaultBinding(t *testing.T) {
	utils.RegisterValidators()
	tests := []struct {
		name  string
		vault Vault
		want  string
	}{
		{name: "an ssh key without a secret passes", vault: Vault{Type: SSH_KEY, Username: "root"}},
		{name: "a username is required", vault: Vault{Type: SSH_PASSWORD, Secret: "x"}, want: "'Vault.Username' Error:Field validation for 'Username' failed on the 'required' tag"},
		{name: "an unknown type is refused", vault: Vault{Type: 9, Username: "root"}, want: "'Vault.Type' Error:Field validation for 'Type' failed on the 'valid' tag"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := binding.Validator.ValidateStruct(&test.vault)
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
