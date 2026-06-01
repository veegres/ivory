package keeper

import (
	"errors"
	"ivory/clients/http"
	"ivory/core/config"
)

var ErrBodyShouldBeEmpty = errors.New("body should be empty")

type Adapter interface {
	List(request Request) ([]Response, int, error)
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
	SupportedFeatures() []env.Feature
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
