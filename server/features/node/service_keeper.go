package node

import (
	"fmt"
	"ivory/plugins/keeper"
	"sync"
)

func (s *Service) KeeperNodeListMulti(r KeeperMultiRequest) ([]KeeperMultiResponse, error) {
	client, tlsConfig, cred, err := s.getKeeperAdapter(r.KeeperOptions)
	if err != nil {
		return nil, err
	}

	keeperAutoResponse := make([]KeeperMultiResponse, len(r.Connections))
	// NOTE: we do not need mutex, because we save always by index
	var wg sync.WaitGroup
	for i, conn := range r.Connections {
		wg.Add(1)
		go func(i int, conn KeeperConnection) {
			defer wg.Done()
			r := keeper.Request{Host: conn.Host, Port: conn.Port, Body: r.Body, TlsConfig: tlsConfig, Credentials: cred}
			response, statusCode, err := client.List(r)
			var errorMessage string
			if err != nil {
				err = fmt.Errorf("host %q failed with code %d: %w", r.Host, statusCode, err)
				errorMessage = err.Error()
			}
			keeperAutoResponse[i] = KeeperMultiResponse{Connection: conn, Response: response, Error: errorMessage}
		}(i, conn)
	}
	wg.Wait()
	return keeperAutoResponse, nil
}

func (s *Service) KeeperNodeList(r KeeperOneRequest) ([]KeeperOneResponse, int, error) {
	client, tlsConfig, cred, err := s.getKeeperAdapter(r.KeeperOptions)
	if err != nil {
		return nil, 0, err
	}
	return client.List(keeper.Request{Host: r.Host, Port: r.Port, Body: r.Body, TlsConfig: tlsConfig, Credentials: cred})
}

func (s *Service) KeeperConfigGet(r KeeperOneRequest) (any, int, error) {
	client, tlsConfig, cred, err := s.getKeeperAdapter(r.KeeperOptions)
	if err != nil {
		return nil, 0, err
	}
	return client.Config(keeper.Request{Host: r.Host, Port: r.Port, Body: r.Body, TlsConfig: tlsConfig, Credentials: cred})
}

func (s *Service) KeeperConfigUpdate(r KeeperOneRequest) (any, int, error) {
	client, tlsConfig, cred, err := s.getKeeperAdapter(r.KeeperOptions)
	if err != nil {
		return nil, 0, err
	}
	return client.ConfigUpdate(keeper.Request{Host: r.Host, Port: r.Port, Body: r.Body, TlsConfig: tlsConfig, Credentials: cred})
}

func (s *Service) KeeperSwitchover(r KeeperOneRequest) (*string, int, error) {
	client, tlsConfig, cred, err := s.getKeeperAdapter(r.KeeperOptions)
	if err != nil {
		return nil, 0, err
	}
	return client.Switchover(keeper.Request{Host: r.Host, Port: r.Port, Body: r.Body, TlsConfig: tlsConfig, Credentials: cred})
}

func (s *Service) KeeperSwitchoverDelete(r KeeperOneRequest) (*string, int, error) {
	client, tlsConfig, cred, err := s.getKeeperAdapter(r.KeeperOptions)
	if err != nil {
		return nil, 0, err
	}
	return client.DeleteSwitchover(keeper.Request{Host: r.Host, Port: r.Port, Body: r.Body, TlsConfig: tlsConfig, Credentials: cred})
}

func (s *Service) KeeperReinitialize(r KeeperOneRequest) (*string, int, error) {
	client, tlsConfig, cred, err := s.getKeeperAdapter(r.KeeperOptions)
	if err != nil {
		return nil, 0, err
	}
	return client.Reinitialize(keeper.Request{Host: r.Host, Port: r.Port, Body: r.Body, TlsConfig: tlsConfig, Credentials: cred})
}

func (s *Service) KeeperRestart(r KeeperOneRequest) (*string, int, error) {
	client, tlsConfig, cred, err := s.getKeeperAdapter(r.KeeperOptions)
	if err != nil {
		return nil, 0, err
	}
	return client.Restart(keeper.Request{Host: r.Host, Port: r.Port, Body: r.Body, TlsConfig: tlsConfig, Credentials: cred})
}

func (s *Service) KeeperRestartDelete(r KeeperOneRequest) (*string, int, error) {
	client, tlsConfig, cred, err := s.getKeeperAdapter(r.KeeperOptions)
	if err != nil {
		return nil, 0, err
	}
	return client.DeleteRestart(keeper.Request{Host: r.Host, Port: r.Port, Body: r.Body, TlsConfig: tlsConfig, Credentials: cred})
}

func (s *Service) KeeperReload(r KeeperOneRequest) (*string, int, error) {
	client, tlsConfig, cred, err := s.getKeeperAdapter(r.KeeperOptions)
	if err != nil {
		return nil, 0, err
	}
	return client.Reload(keeper.Request{Host: r.Host, Port: r.Port, Body: r.Body, TlsConfig: tlsConfig, Credentials: cred})
}

func (s *Service) KeeperFailover(r KeeperOneRequest) (*string, int, error) {
	client, tlsConfig, cred, err := s.getKeeperAdapter(r.KeeperOptions)
	if err != nil {
		return nil, 0, err
	}
	return client.Failover(keeper.Request{Host: r.Host, Port: r.Port, Body: r.Body, TlsConfig: tlsConfig, Credentials: cred})
}

func (s *Service) KeeperActivate(r KeeperOneRequest) (*string, int, error) {
	client, tlsConfig, cred, err := s.getKeeperAdapter(r.KeeperOptions)
	if err != nil {
		return nil, 0, err
	}
	return client.Activate(keeper.Request{Host: r.Host, Port: r.Port, Body: r.Body, TlsConfig: tlsConfig, Credentials: cred})
}

func (s *Service) KeeperPause(r KeeperOneRequest) (*string, int, error) {
	client, tlsConfig, cred, err := s.getKeeperAdapter(r.KeeperOptions)
	if err != nil {
		return nil, 0, err
	}
	return client.Pause(keeper.Request{Host: r.Host, Port: r.Port, Body: r.Body, TlsConfig: tlsConfig, Credentials: cred})
}
