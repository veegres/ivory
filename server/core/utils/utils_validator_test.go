package utils

import (
	"strings"
	"testing"

	"github.com/gin-gonic/gin/binding"
)

type testColor string

func (c testColor) Valid() bool { return c == "red" || c == "blue" }

type testNode struct {
	Name string `json:"name" binding:"notblank"`
}

type testRequest struct {
	Name     string      `json:"name" binding:"required"`
	Color    testColor   `json:"color" binding:"valid"`
	Optional *testColor  `json:"optional" binding:"omitempty,valid"`
	Colors   []testColor `json:"colors" binding:"min=1,dive,valid"`
	Nodes    []testNode  `json:"nodes" binding:"dive"`
}

func validRequest() testRequest {
	return testRequest{
		Name:   "prod",
		Color:  "red",
		Colors: []testColor{"blue"},
		Nodes:  []testNode{{Name: "pg1"}},
	}
}

func TestRegisterValidators(t *testing.T) {
	RegisterValidators()
	RegisterValidators()
	green := testColor("green")

	tests := []struct {
		name   string
		modify func(r *testRequest)
		want   string
	}{
		{name: "a complete request passes", modify: func(r *testRequest) {}},
		{name: "unknown enum value", modify: func(r *testRequest) { r.Color = "green" }, want: "'testRequest.Color' Error:Field validation for 'Color' failed on the 'valid' tag"},
		{name: "nil optional enum is skipped", modify: func(r *testRequest) { r.Optional = nil }},
		{name: "set optional enum is checked", modify: func(r *testRequest) { r.Optional = &green }, want: "'testRequest.Optional' Error:Field validation for 'Optional' failed on the 'valid' tag"},
		{name: "dive reaches each element", modify: func(r *testRequest) { r.Colors = []testColor{"red", "pink"} }, want: "'testRequest.Colors[1]' Error:Field validation for 'Colors[1]' failed on the 'valid' tag"},
		{name: "blank name", modify: func(r *testRequest) { r.Nodes[0].Name = "  " }, want: "'testRequest.Nodes[0].Name' Error:Field validation for 'Name' failed on the 'notblank' tag"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := validRequest()
			test.modify(&request)
			err := binding.Validator.ValidateStruct(&request)
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

	t.Run("a type without Valid under the valid rule is a programming error", func(t *testing.T) {
		defer func() {
			if recover() == nil {
				t.Error("expected a panic")
			}
		}()
		_ = binding.Validator.ValidateStruct(struct {
			Name string `binding:"valid"`
		}{Name: "x"})
	})
}
