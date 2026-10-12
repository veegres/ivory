package query

import (
	"ivory/core/utils"
	"strings"
	"testing"

	"github.com/gin-gonic/gin/binding"
	"github.com/google/uuid"
)

func TestRequestBinding(t *testing.T) {
	utils.RegisterValidators()
	bloat, unknown := BLOAT, Type(42)
	connection := Connection{Db: DbConfig{Plugin: "postgres", Host: "db1", Port: 5432}}

	tests := []struct {
		name    string
		request any
		want    string
	}{
		{name: "query passes", request: &Request{Name: "bloat", Type: &bloat, Query: "select 1", Varieties: []VarietyType{MasterOnly}}},
		{name: "query plugin may be left out", request: &Request{Name: "bloat", Type: &bloat, Query: "select 1"}},
		{name: "query needs a type", request: &Request{Name: "bloat", Query: "select 1"}, want: "'Request.Type' Error:Field validation for 'Type' failed on the 'required' tag"},
		{name: "query type must be known", request: &Request{Name: "bloat", Type: &unknown, Query: "select 1"}, want: "'Request.Type' Error:Field validation for 'Type' failed on the 'valid' tag"},
		{name: "query plugin must be known", request: &Request{Name: "bloat", Type: &bloat, Plugin: "oracle", Query: "select 1"}, want: "'Request.Plugin' Error:Field validation for 'Plugin' failed on the 'valid' tag"},
		{name: "query varieties must be known", request: &Request{Name: "bloat", Type: &bloat, Query: "select 1", Varieties: []VarietyType{7}}, want: "'Request.Varieties[0]' Error:Field validation for 'Varieties[0]' failed on the 'valid' tag"},
		{name: "query text is not blank", request: &Request{Name: "bloat", Type: &bloat, Query: " "}, want: "'Request.Query' Error:Field validation for 'Query' failed on the 'notblank' tag"},
		{name: "connection needs a plugin", request: &Connection{Db: DbConfig{Host: "db1", Port: 5432}}, want: "'Connection.Db.Plugin' Error:Field validation for 'Plugin' failed on the 'required' tag"},
		{name: "connection port is a port", request: &Connection{Db: DbConfig{Plugin: "postgres", Host: "db1", Port: 0}}, want: "'Connection.Db.Port' Error:Field validation for 'Port' failed on the 'required' tag"},
		{name: "template names its query", request: &RouterTemplateRequest{Connection: connection}, want: "'RouterTemplateRequest.QueryUuid' Error:Field validation for 'QueryUuid' failed on the 'required' tag"},
		{name: "template with a query passes", request: &RouterTemplateRequest{Connection: connection, QueryUuid: uuid.New()}},
		{name: "kill names a pid", request: &RouterKillRequest{Connection: connection}, want: "'RouterKillRequest.Pid' Error:Field validation for 'Pid' failed on the 'required' tag"},
		{name: "chart type must be known", request: &RouterChartRequest{Connection: connection, Type: "Uptime"}, want: "'RouterChartRequest.Type' Error:Field validation for 'Type' failed on the 'valid' tag"},
		{name: "chart passes", request: &RouterChartRequest{Connection: connection, Type: Connections}},
		{name: "id request needs an id", request: &RouterIdRequest{}, want: "'RouterIdRequest.Id' Error:Field validation for 'Id' failed on the 'required' tag"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := binding.Validator.ValidateStruct(test.request)
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
