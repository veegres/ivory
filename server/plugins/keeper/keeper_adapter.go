// Package keeper defines the plugin boundary for HA management tools
// (patroni, native postgres, native etcd). A keeper.Adapter lets the rest of
// Ivory inspect and operate a cluster's keeper without knowing which
// concrete plugin is running it.
package keeper

import (
	"errors"
	"ivory/clients/http"
)

var ErrBodyShouldBeEmpty = errors.New("body should be empty")
var ErrNotSupported = errors.New("operation is not supported by this keeper plugin")

// Plugin is the whole contract a keeper plugin registers under: operations
// against a running keeper plus its own self-description. It exists so a
// plugin is registered once, in one registry, while consumers keep depending
// on the half they actually use.
type Plugin interface {
	Adapter
	Metadata
}

// Adapter covers operations against a running keeper; plugin self-description
// (supported features, deployment defaults) lives in Metadata. Methods a
// plugin cannot perform (e.g. patroni-only orchestration on native postgres)
// return ErrNotSupported.
type Adapter interface {
	// List returns the cluster overview: one Response per known member.
	List(request Request) ([]Response, int, error)
	// Config returns the keeper's current configuration.
	Config(request Request) (any, int, error)
	// ConfigUpdate patches the keeper's configuration with request.Body.
	ConfigUpdate(request Request) (any, int, error)
	// Switchover schedules or performs a leader change to the candidate
	// named in request.Body.
	Switchover(request Request) (*string, int, error)
	// DeleteSwitchover cancels a pending scheduled switchover.
	DeleteSwitchover(request Request) (*string, int, error)
	// Reinitialize re-creates a member's data directory from the leader.
	Reinitialize(request Request) (*string, int, error)
	// Restart restarts (or schedules a restart of) the keeper-managed process.
	Restart(request Request) (*string, int, error)
	// DeleteRestart cancels a pending scheduled restart.
	DeleteRestart(request Request) (*string, int, error)
	// Reload asks the keeper-managed process to reload its configuration.
	Reload(request Request) (*string, int, error)
	// Failover forces an immediate leader change, bypassing any schedule.
	Failover(request Request) (*string, int, error)
	// Activate resumes automatic failover/switchover after a Pause.
	Activate(request Request) (*string, int, error)
	// Pause suspends automatic failover/switchover.
	Pause(request Request) (*string, int, error)
}

func Map(request Request, path string) http.Request {
	var creds *http.Credentials
	if request.Credentials != nil {
		creds = &http.Credentials{
			Username: request.Credentials.Username,
			Password: request.Credentials.Password,
		}
	}

	return http.Request{
		Host:        request.Host,
		Port:        request.Port,
		Path:        path,
		Body:        request.Body,
		Credentials: creds,
		TLSConfig:   request.TlsConfig,
	}
}
