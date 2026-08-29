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
			if strings.Contains(template.Commands[0].Command, "REDIS_REPLICATION_MODE") {
				t.Error("the leader must not come up as a replica")
			}
			// NOTE: on one VM the leader is reached at {{host}} - host
			// networking joins no docker network, so its container name
			// resolves to nothing
			leader := `REDIS_MASTER_HOST="redis1"`
			if strings.Contains(template.Name, "Single Host") {
				leader = `REDIS_MASTER_HOST="{{host}}"`
			}
			for i, command := range template.Commands[1:] {
				if !strings.Contains(command.Command, "REDIS_REPLICATION_MODE") {
					t.Errorf("replica %d does not attach to the leader", i+1)
				}
				if !strings.Contains(command.Command, leader) {
					t.Errorf("replica %d has no leader host to attach to", i+1)
				}
			}
			// NOTE: the official image takes port/password as CLI flags only,
			// so the leader - which runs no bootstrap - could not be
			// configured through env vars at all
			if !strings.Contains(template.Commands[0].Command, "bitnami/redis") {
				t.Error("expected bitnami/redis, the only image configurable through env vars alone")
			}
		})
	}
}
