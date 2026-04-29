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

type Registry struct {
	clients map[Type]Adapter
}

func NewRegistry() *Registry {
	return &Registry{
		clients: make(map[Type]Adapter),
	}
}

func (r *Registry) Register(t Type, client Adapter) {
	r.clients[t] = client
}

func (r *Registry) Get(t Type) (Adapter, error) {
	if client, ok := r.clients[t]; ok {
		return client, nil
	}
	return nil, ErrClientNotImplemented
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
