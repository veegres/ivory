package node

import (
	"fmt"
	"ivory/src/plugins/keeper"
	"sync"
)

func (s *Service) ListParallel(r KeeperParallelRequest) ([]KeeperParallelResponse, error) {
	client, tlsConfig, cred, err := s.getKeeperAdapter(r.KeeperOptions)
	if err != nil {
		return nil, err
	}

	keeperAutoResponse := make([]KeeperParallelResponse, len(r.Connections))
	// NOTE: we do not need mutex, because we save always by index
	var wg sync.WaitGroup
	for i, conn := range r.Connections {
		wg.Add(1)
		go func(i int, conn KeeperConnection) {
			defer wg.Done()
			r := keeper.Request{Host: conn.Host, Port: conn.Port, Body: r.Body, TlsConfig: tlsConfig, Credentials: cred}
			response, statusCode, err := client.List(r)
			if err != nil {
				err = fmt.Errorf("host %q failed with code %d: %w", r.Host, statusCode, err)
			}
			keeperAutoResponse[i] = KeeperParallelResponse{Connection: conn, Response: response, Error: err}
		}(i, conn)
	}
	wg.Wait()
	return keeperAutoResponse, nil
}

func (s *Service) List(r KeeperRequest) ([]KeeperResponse, int, error) {
	client, tlsConfig, cred, err := s.getKeeperAdapter(r.KeeperOptions)
	if err != nil {
		return nil, 0, err
	}
	return client.List(keeper.Request{Host: r.Host, Port: r.Port, Body: r.Body, TlsConfig: tlsConfig, Credentials: cred})
}

func (s *Service) Config(r KeeperRequest) (any, int, error) {
	client, tlsConfig, cred, err := s.getKeeperAdapter(r.KeeperOptions)
	if err != nil {
		return nil, 0, err
	}
	return client.Config(keeper.Request{Host: r.Host, Port: r.Port, Body: r.Body, TlsConfig: tlsConfig, Credentials: cred})
}

func (s *Service) ConfigUpdate(r KeeperRequest) (any, int, error) {
	client, tlsConfig, cred, err := s.getKeeperAdapter(r.KeeperOptions)
	if err != nil {
		return nil, 0, err
	}
	return client.ConfigUpdate(keeper.Request{Host: r.Host, Port: r.Port, Body: r.Body, TlsConfig: tlsConfig, Credentials: cred})
}

func (s *Service) Switchover(r KeeperRequest) (*string, int, error) {
	client, tlsConfig, cred, err := s.getKeeperAdapter(r.KeeperOptions)
	if err != nil {
		return nil, 0, err
	}
	return client.Switchover(keeper.Request{Host: r.Host, Port: r.Port, Body: r.Body, TlsConfig: tlsConfig, Credentials: cred})
}

func (s *Service) DeleteSwitchover(r KeeperRequest) (*string, int, error) {
	client, tlsConfig, cred, err := s.getKeeperAdapter(r.KeeperOptions)
	if err != nil {
		return nil, 0, err
	}
	return client.DeleteSwitchover(keeper.Request{Host: r.Host, Port: r.Port, Body: r.Body, TlsConfig: tlsConfig, Credentials: cred})
}

func (s *Service) Reinitialize(r KeeperRequest) (*string, int, error) {
	client, tlsConfig, cred, err := s.getKeeperAdapter(r.KeeperOptions)
	if err != nil {
		return nil, 0, err
	}
	return client.Reinitialize(keeper.Request{Host: r.Host, Port: r.Port, Body: r.Body, TlsConfig: tlsConfig, Credentials: cred})
}

func (s *Service) Restart(r KeeperRequest) (*string, int, error) {
	client, tlsConfig, cred, err := s.getKeeperAdapter(r.KeeperOptions)
	if err != nil {
		return nil, 0, err
	}
	return client.Restart(keeper.Request{Host: r.Host, Port: r.Port, Body: r.Body, TlsConfig: tlsConfig, Credentials: cred})
}

func (s *Service) DeleteRestart(r KeeperRequest) (*string, int, error) {
	client, tlsConfig, cred, err := s.getKeeperAdapter(r.KeeperOptions)
	if err != nil {
		return nil, 0, err
	}
	return client.DeleteRestart(keeper.Request{Host: r.Host, Port: r.Port, Body: r.Body, TlsConfig: tlsConfig, Credentials: cred})
}

func (s *Service) Reload(r KeeperRequest) (*string, int, error) {
	client, tlsConfig, cred, err := s.getKeeperAdapter(r.KeeperOptions)
	if err != nil {
		return nil, 0, err
	}
	return client.Reload(keeper.Request{Host: r.Host, Port: r.Port, Body: r.Body, TlsConfig: tlsConfig, Credentials: cred})
}

func (s *Service) Failover(r KeeperRequest) (*string, int, error) {
	client, tlsConfig, cred, err := s.getKeeperAdapter(r.KeeperOptions)
	if err != nil {
		return nil, 0, err
	}
	return client.Failover(keeper.Request{Host: r.Host, Port: r.Port, Body: r.Body, TlsConfig: tlsConfig, Credentials: cred})
}

func (s *Service) Activate(r KeeperRequest) (*string, int, error) {
	client, tlsConfig, cred, err := s.getKeeperAdapter(r.KeeperOptions)
	if err != nil {
		return nil, 0, err
	}
	return client.Activate(keeper.Request{Host: r.Host, Port: r.Port, Body: r.Body, TlsConfig: tlsConfig, Credentials: cred})
}

func (s *Service) Pause(r KeeperRequest) (*string, int, error) {
	client, tlsConfig, cred, err := s.getKeeperAdapter(r.KeeperOptions)
	if err != nil {
		return nil, 0, err
	}
	return client.Pause(keeper.Request{Host: r.Host, Port: r.Port, Body: r.Body, TlsConfig: tlsConfig, Credentials: cred})
}
