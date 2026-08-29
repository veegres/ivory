package redis

import (
	"ivory/core/config"
	"strings"
	"testing"
)

func TestSupportedFeaturesExclusions(t *testing.T) {
	features := NewAdapter().SupportedFeatures()

	supported := []config.Feature{config.ViewNodeKeeperOverview, config.ViewNodeKeeperConfig, config.ManageNodeKeeperFailover}
	for _, feature := range supported {
		if !features[feature] {
			t.Errorf("feature %v must be supported for redis", feature)
		}
	}

	excluded := []config.Feature{
		config.ManageNodeKeeperConfigUpdate, config.ManageNodeKeeperSwitchover, config.ManageNodeKeeperReinitialize,
		config.ManageNodeKeeperRestart, config.ManageNodeKeeperReload, config.ManageNodeKeeperActivation,
	}
	for _, feature := range excluded {
		if features[feature] {
			t.Errorf("feature %v must not be supported for redis", feature)
		}
	}
}

func TestRequirements(t *testing.T) {
	req := NewAdapter().Requirements()

	if !req.KeeperCredentials || req.KeeperUser != "default" {
		t.Errorf("expected keeper credentials locked to redis' default username, got %v/%q", req.KeeperCredentials, req.KeeperUser)
	}
	if !req.DbCredentials || req.DbUser != "default" {
		t.Errorf("expected credentials with redis' default username, got %+v", req)
	}
}

// TestDefaultTemplates covers the leader/replica asymmetry: it is expressed as
// a different command at index 0, not as a flag.
func TestDefaultTemplates(t *testing.T) {
	templates := NewAdapter().DefaultTemplates()

	if len(templates) != 2 {
		t.Fatalf("expected a multi-host and a single-host template, got %d", len(templates))
	}

	for _, template := range templates {
		t.Run(template.Name, func(t *testing.T) {
			if strings.Contains(template.Commands[0].Command, "--replicaof") {
				t.Error("the leader must not come up as a replica")
			}
			// NOTE: neither a plain docker run across VMs nor host networking
			// resolves a container name, so the leader is an address either
			// way - literal example text on separate VMs, {{host}} on one
			leader := "--replicaof 10.0.0.1 6379"
			if strings.Contains(template.Name, "Single Host") {
				leader = "--replicaof {{host}} 6379"
			}
			for i, command := range template.Commands[1:] {
				if !strings.Contains(command.Command, leader) {
					t.Errorf("replica %d does not attach to the leader", i+1)
				}
				if !strings.Contains(command.Command, `--masterauth "{{dbPass}}"`) {
					t.Errorf("replica %d cannot authenticate to a password-protected leader", i+1)
				}
			}
			// NOTE: bitnami/redis was retired from Docker Hub, so every
			// command states its port and password as redis-server flags
			// instead - the official image takes no env vars for them
			for i, command := range template.Commands {
				if !strings.Contains(command.Command, "redis:7.4") || strings.Contains(command.Command, "bitnami") {
					t.Errorf("command %d must use the official image, bitnami/redis no longer exists", i)
				}
				if !strings.Contains(command.Command, "--port {{dbPort}}") {
					t.Errorf("command %d does not state its own port", i)
				}
				if !strings.Contains(command.Command, `--requirepass "{{dbPass}}"`) {
					t.Errorf("command %d comes up without a password", i)
				}
			}
		})
	}
}
