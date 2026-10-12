package cert

import (
	"ivory/core/utils"
	"strings"
	"testing"

	"github.com/gin-gonic/gin/binding"
)

func TestRouterAddRequestBinding(t *testing.T) {
	utils.RegisterValidators()
	tests := []struct {
		name    string
		request RouterAddRequest
		want    string
	}{
		{name: "a path and a known type pass", request: RouterAddRequest{Path: "/etc/ca.crt", Type: CLIENT_CA}},
		{name: "a path is required", request: RouterAddRequest{Type: CLIENT_KEY}, want: "'RouterAddRequest.Path' Error:Field validation for 'Path' failed on the 'required' tag"},
		{name: "an unknown type is refused", request: RouterAddRequest{Path: "/etc/ca.crt", Type: 5}, want: "'RouterAddRequest.Type' Error:Field validation for 'Type' failed on the 'valid' tag"},
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
