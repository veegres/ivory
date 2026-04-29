package keeper

import (
	"errors"
	"ivory/src/clients/http"
	"ivory/src/features"
)

var ErrBodyShouldBeEmpty = errors.New("body should be empty")
var ErrClientNotImplemented = errors.New("client is not implemented")

type Adapter interface {
	Overview(request Request) ([]Response, int, error)
	Config(request Request) (any, int, error)
	ConfigUpdate(request Request) (any, int, error)
	Switchover(request Request) (*string, int, error)
	DeleteSwitchover(request Request) (*string, int, error)
	Reinitialize(request Request) (*string, int, error)
	Restart(request Request) (*string, int, error)
	DeleteRestart(request Request) (*string, int, error)
	Reload(request Request) (*string, int, error)
	Failover(request Request) (*string, int, error)
	Activate(request Request) (*string, int, error)
	Pause(request Request) (*string, int, error)
	SupportedFeatures() []features.Feature
}

type PluginRegistry struct {
	clients map[Plugin]Adapter
}

func NewPluginRegistry() *PluginRegistry {
	return &PluginRegistry{
		clients: make(map[Plugin]Adapter),
	}
}

func (r *PluginRegistry) Register(t Plugin, client Adapter) {
	r.clients[t] = client
}

func (r *PluginRegistry) Get(t Plugin) (Adapter, error) {
	if client, ok := r.clients[t]; ok {
		return client, nil
	}
	return nil, ErrClientNotImplemented
}

func (r *PluginRegistry) All() map[Plugin]Adapter {
	return r.clients
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
