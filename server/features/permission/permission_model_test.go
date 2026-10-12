package permission

import (
	"ivory/core/config"
	"ivory/core/utils"
	"strings"
	"testing"

	"github.com/gin-gonic/gin/binding"
)

func TestRouterPermissionRequestBinding(t *testing.T) {
	utils.RegisterValidators()
	tests := []struct {
		name        string
		permissions []config.Feature
		want        string
	}{
		{name: "known features pass", permissions: []config.Feature{config.ViewClusterList}},
		{name: "at least one feature", permissions: []config.Feature{}, want: "'RouterPermissionRequest.Permissions' Error:Field validation for 'Permissions' failed on the 'min' tag"},
		{name: "a feature that does not exist", permissions: []config.Feature{config.ViewClusterList, "view.cluster.everything"}, want: "'RouterPermissionRequest.Permissions[1]' Error:Field validation for 'Permissions[1]' failed on the 'valid' tag"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := binding.Validator.ValidateStruct(&RouterPermissionRequest{Permissions: test.permissions})
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
