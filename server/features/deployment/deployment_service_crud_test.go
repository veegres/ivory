package deployment

import (
	"errors"
	"ivory/core/utils"
	"ivory/plugins/keeper"
	"ivory/plugins/keeper/etcd"
	"ivory/plugins/platform"
	"ivory/plugins/platform/docker"
	"strings"
	"testing"

	"github.com/google/uuid"
)

func newTestService(t *testing.T) *Service {
	t.Helper()

	keeperMetadataRegistry := utils.NewRegistry[keeper.Plugin, keeper.Metadata]()
	keeperMetadataRegistry.Register(keeper.NATIVE_ETCD, etcd.NewAdapter())
	platformRegistry := utils.NewRegistry[platform.Plugin, platform.Adapter]()
	platformRegistry.Register(platform.Docker, docker.NewAdapter(nil))

	return NewService(newTestRepository(t), keeperMetadataRegistry, platformRegistry)
}

func testRequest(name string) TemplateRequest {
	return TemplateRequest{
		Name:     name,
		Keeper:   keeper.NATIVE_ETCD,
		Platform: platform.Docker,
		Commands: []TemplateCommand{{
			Command: `docker run -d --name {{name}} -e ETCD_INITIAL_CLUSTER="etcd-1=http://etcd-1:2380" etcd`,
		}},
	}
}

func TestService_Create(t *testing.T) {
	tests := []struct {
		name     string
		mutate   func(r *TemplateRequest)
		expected error
		contains string
	}{
		{
			name:   "a valid template is created",
			mutate: func(r *TemplateRequest) {},
		},
		{
			name:     "a blank name is rejected",
			mutate:   func(r *TemplateRequest) { r.Name = "   " },
			expected: ErrTemplateNameRequired,
		},
		{
			name:     "a template with no commands deploys nothing",
			mutate:   func(r *TemplateRequest) { r.Commands = nil },
			expected: ErrTemplateCommandsRequired,
		},
		{
			// NOTE: this is what keeps the vocabulary closed - a typo must be
			// an error, never a silently-introduced new variable
			name:     "an unknown variable in a command is rejected",
			mutate:   func(r *TemplateRequest) { r.Commands[0].Command = "docker run -d {{clusterHosts}}" },
			contains: "{{clusterHosts}}",
		},
		{
			name:     "an unknown variable in a post script is rejected",
			mutate:   func(r *TemplateRequest) { r.Commands[0].PostScript = "etcdctl user add {{dbUesr}}" },
			contains: "{{dbUesr}}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := newTestService(t)
			r := testRequest("mine")
			tt.mutate(&r)

			created, err := s.Create(r)
			switch {
			case tt.expected != nil:
				if !errors.Is(err, tt.expected) {
					t.Fatalf("expected %v, got %v", tt.expected, err)
				}
			case tt.contains != "":
				if err == nil || !strings.Contains(err.Error(), tt.contains) {
					t.Fatalf("expected an error naming %s, got %v", tt.contains, err)
				}
			default:
				if err != nil {
					t.Fatalf("Create() error = %v", err)
				}
				if created.Creation != Manual {
					t.Errorf("a user-created template must be manual, got %q", created.Creation)
				}
			}
		})
	}
}

func TestService_CreateRejectsTakenNames(t *testing.T) {
	t.Run("a stored name", func(t *testing.T) {
		s := newTestService(t)
		if _, err := s.Create(testRequest("mine")); err != nil {
			t.Fatalf("Create() error = %v", err)
		}
		if _, err := s.Create(testRequest("mine")); !errors.Is(err, ErrTemplateNameTaken) {
			t.Fatalf("expected ErrTemplateNameTaken, got %v", err)
		}
	})

	// NOTE: the name is what the user reads, so surrounding space is not what
	// makes two of them different
	t.Run("a stored name with surrounding space", func(t *testing.T) {
		s := newTestService(t)
		if _, err := s.Create(testRequest("mine")); err != nil {
			t.Fatalf("Create() error = %v", err)
		}
		if _, err := s.Create(testRequest("  mine  ")); !errors.Is(err, ErrTemplateNameTaken) {
			t.Fatalf("expected ErrTemplateNameTaken, got %v", err)
		}
	})

	// NOTE: the list is filtered by keeper and platform, so a name taken under
	// another keeper is one the user can neither see nor resolve
	t.Run("but not a name taken under another keeper", func(t *testing.T) {
		s := newTestService(t)
		if _, err := s.Create(testRequest("mine")); err != nil {
			t.Fatalf("Create() error = %v", err)
		}
		other := testRequest("mine")
		other.Keeper = keeper.NATIVE_REDIS
		if _, err := s.Create(other); err != nil {
			t.Fatalf("Create() error = %v, want the name to be free here", err)
		}
	})

	// NOTE: defaults share one list with the user's own templates, so a copy
	// reusing a shipped name would be ambiguous in the UI
	t.Run("a shipped default's name", func(t *testing.T) {
		s := newTestService(t)
		defaults := s.Defaults(ListRequest{})
		if len(defaults) == 0 {
			t.Fatal("expected the etcd defaults to be available")
		}
		if _, err := s.Create(testRequest(defaults[0].Name)); !errors.Is(err, ErrTemplateNameTaken) {
			t.Fatalf("expected ErrTemplateNameTaken, got %v", err)
		}
	})
}

func TestService_CreateTrimsTheName(t *testing.T) {
	s := newTestService(t)
	created, err := s.Create(testRequest("  mine  "))
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	if created.Name != "mine" {
		t.Errorf("Name = %q, want %q", created.Name, "mine")
	}
}

func TestService_Update(t *testing.T) {
	s := newTestService(t)
	created, err := s.Create(testRequest("mine"))
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	id := uuid.MustParse(created.Id)

	t.Run("keeping its own name is not a collision with itself", func(t *testing.T) {
		r := testRequest("mine")
		r.Commands[0].Command = "docker run -d --name {{name}} redis:7"
		got, err := s.Update(id, r)
		if err != nil {
			t.Fatalf("Update() error = %v", err)
		}
		if got.Commands[0].Command != r.Commands[0].Command {
			t.Errorf("expected the edited command to be saved, got %q", got.Commands[0].Command)
		}
	})

	t.Run("the keeper cannot be changed", func(t *testing.T) {
		r := testRequest("mine")
		r.Keeper = "native_redis"
		if _, err := s.Update(id, r); !errors.Is(err, ErrTemplatePluginImmutable) {
			t.Fatalf("expected ErrTemplatePluginImmutable, got %v", err)
		}
	})

	t.Run("the platform cannot be changed", func(t *testing.T) {
		r := testRequest("mine")
		r.Platform = "k8s"
		if _, err := s.Update(id, r); !errors.Is(err, ErrTemplatePluginImmutable) {
			t.Fatalf("expected ErrTemplatePluginImmutable, got %v", err)
		}
	})

	t.Run("renaming onto a taken name is rejected", func(t *testing.T) {
		if _, err := s.Create(testRequest("other")); err != nil {
			t.Fatalf("Create() error = %v", err)
		}
		if _, err := s.Update(id, testRequest("other")); !errors.Is(err, ErrTemplateNameTaken) {
			t.Fatalf("expected ErrTemplateNameTaken, got %v", err)
		}
	})
}

func TestService_ListReturnsStoredThenDefaults(t *testing.T) {
	s := newTestService(t)
	if _, err := s.Create(testRequest("mine")); err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	list, err := s.List(ListRequest{})
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}
	if len(list) < 2 {
		t.Fatalf("expected the stored template plus the shipped defaults, got %d", len(list))
	}
	if list[0].Name != "mine" || list[0].Creation != Manual {
		t.Errorf("expected the user's own template first, got %+v", list[0])
	}
	for _, template := range list[1:] {
		if template.Creation != System {
			t.Errorf("expected every template after the stored ones to be a system one, got %+v", template)
		}
	}
}

func TestService_Delete(t *testing.T) {
	s := newTestService(t)
	created, err := s.Create(testRequest("mine"))
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	if err := s.Delete(uuid.MustParse(created.Id)); err != nil {
		t.Fatalf("Delete() error = %v", err)
	}
	if _, err := s.Get(uuid.MustParse(created.Id)); err == nil {
		t.Error("expected the template to be gone")
	}
}

// TestIsDefaultId covers the guard that makes a write against a shipped
// template impossible to express: their ids are not uuids at all.
func TestIsDefaultId(t *testing.T) {
	s := newTestService(t)

	for _, template := range s.Defaults(ListRequest{}) {
		if !isDefaultId(template.Id) {
			t.Errorf("shipped template %q has id %q, which a write would accept", template.Name, template.Id)
		}
		if _, err := uuid.Parse(template.Id); err == nil {
			t.Errorf("shipped template %q has a uuid id, which could collide with a stored one", template.Name)
		}
	}

	if isDefaultId(uuid.New().String()) {
		t.Error("a stored template's id must not look like a shipped one")
	}
}

// TestService_CreateRenamesLegacyPlatform covers the one live path that still
// carries the old key: importing a backup written before linux was renamed to
// docker. The template must be stored - and name-checked - as docker.
func TestService_CreateRenamesLegacyPlatform(t *testing.T) {
	s := newTestService(t)

	r := testRequest("imported")
	r.Platform = "linux"

	created, err := s.Create(r)
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	if created.Platform != platform.Docker {
		t.Errorf("Create() platform = %q, want %q", created.Platform, platform.Docker)
	}

	t.Run("the stored name is taken under the current platform too", func(t *testing.T) {
		again := testRequest("imported")
		if _, err := s.Create(again); !errors.Is(err, ErrTemplateNameTaken) {
			t.Fatalf("expected ErrTemplateNameTaken, got %v", err)
		}
	})

	t.Run("it lists under the current platform", func(t *testing.T) {
		list, err := s.List(ListRequest{Platform: platformPtr(platform.Docker)})
		if err != nil {
			t.Fatalf("List() error = %v", err)
		}
		if len(list) == 0 || list[0].Name != "imported" {
			t.Fatalf("List() = %+v, want the imported template first", list)
		}
	})
}
