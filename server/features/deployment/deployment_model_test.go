package deployment

import (
	"ivory/core/utils"
	"strings"
	"testing"

	"github.com/gin-gonic/gin/binding"
)

func TestTemplateRequestBinding(t *testing.T) {
	utils.RegisterValidators()
	valid := func() TemplateRequest {
		return TemplateRequest{
			Name:     "etcd",
			Keeper:   "native_etcd",
			Platform: "docker",
			Commands: []TemplateCommand{{Command: "-d etcd", Defaults: CommandDefaults{Name: "etcd1", KeeperPort: 2379}}},
		}
	}

	tests := []struct {
		name   string
		modify func(r *TemplateRequest)
		want   string
	}{
		{name: "a complete template passes", modify: func(r *TemplateRequest) {}},
		{name: "unstated ports are allowed", modify: func(r *TemplateRequest) { r.Commands[0].Defaults.KeeperPort = 0 }},
		{name: "name is not blank", modify: func(r *TemplateRequest) { r.Name = " " }, want: "'TemplateRequest.Name' Error:Field validation for 'Name' failed on the 'notblank' tag"},
		{name: "keeper must be known", modify: func(r *TemplateRequest) { r.Keeper = "pgpool" }, want: "'TemplateRequest.Keeper' Error:Field validation for 'Keeper' failed on the 'valid' tag"},
		{name: "platform is required", modify: func(r *TemplateRequest) { r.Platform = "" }, want: "'TemplateRequest.Platform' Error:Field validation for 'Platform' failed on the 'required' tag"},
		{name: "at least one command", modify: func(r *TemplateRequest) { r.Commands = []TemplateCommand{} }, want: "'TemplateRequest.Commands' Error:Field validation for 'Commands' failed on the 'min' tag"},
		{name: "a command is not blank", modify: func(r *TemplateRequest) { r.Commands[0].Command = "\n" }, want: "'TemplateRequest.Commands[0].Command' Error:Field validation for 'Command' failed on the 'notblank' tag"},
		{name: "a stated port is a port", modify: func(r *TemplateRequest) { r.Commands[0].Defaults.DbPort = 70000 }, want: "'TemplateRequest.Commands[0].Defaults.DbPort' Error:Field validation for 'DbPort' failed on the 'max' tag"},
		{name: "a negative port is refused", modify: func(r *TemplateRequest) { r.Commands[0].Defaults.SshPort = -1 }, want: "'TemplateRequest.Commands[0].Defaults.SshPort' Error:Field validation for 'SshPort' failed on the 'min' tag"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := valid()
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
}

func TestListRequestBinding(t *testing.T) {
	utils.RegisterValidators()
	unknown := KeeperPlugin("pgpool")
	if err := binding.Validator.ValidateStruct(&ListRequest{}); err != nil {
		t.Fatalf("an empty filter must pass, got %v", err)
	}
	err := binding.Validator.ValidateStruct(&ListRequest{Keeper: &unknown})
	if err == nil || !strings.Contains(err.Error(), "'ListRequest.Keeper' Error:Field validation for 'Keeper' failed on the 'valid' tag") {
		t.Fatalf("got %v", err)
	}
}
