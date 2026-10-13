package node

import (
	"crypto/tls"
	"fmt"
	"ivory/core/utils"
	"ivory/plugins/keeper"
	"net/http"
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
			r := keeper.Request{Cluster: r.Cluster, Host: conn.Host, Port: conn.Port, Body: r.Body, TlsConfig: tlsConfig, Credentials: cred}
			response, statusCode, err := client.List(r)
			var errorMessage string
			if err != nil {
				err = fmt.Errorf("host %q failed with code %d: %w", r.Host, statusCode, err)
				errorMessage = err.Error()
			}
			mapped := make([]KeeperResponse, len(response))
			for j, resp := range response {
				mapped[j] = mapKeeperResponse(resp)
			}
			keeperAutoResponse[i] = KeeperMultiResponse{Connection: conn, Response: mapped, Error: errorMessage}
		}(i, conn)
	}
	wg.Wait()
	return keeperAutoResponse, nil
}

func (s *Service) KeeperNodeList(r KeeperOneRequest) ([]KeeperOneResponse, int, error) {
	client, tlsConfig, cred, err := s.getKeeperAdapter(r.KeeperOptions)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}
	// NOTE: responses is still mapped even when err != nil, since an adapter
	// can report a degraded-but-useful response (e.g. postgres starting up)
	// alongside the error that caused it; callers rely on getting both so the
	// error still ends up in the node's warnings instead of being discarded.
	responses, status, err := client.List(keeper.Request{Cluster: r.Cluster, Host: r.Host, Port: r.Port, Body: r.Body, TlsConfig: tlsConfig, Credentials: cred})
	result := make([]KeeperOneResponse, len(responses))
	for i, resp := range responses {
		result[i] = mapKeeperResponse(resp)
	}
	return result, status, err
}

// KeeperCheck reports how far Ivory gets reaching the keeper, one step at a
// time. The error is only for a plugin Ivory does not know.
func (s *Service) KeeperCheck(r KeeperOneRequest) ([]utils.CheckStep, error) {
	adapter, err := s.keeperRegistry.Get(r.Plugin)
	if err != nil {
		return nil, err
	}
	check := utils.NewCheck()
	var cred *keeper.Credentials
	if r.VaultId != nil {
		check.Run(utils.CheckStageVault, func() error {
			d, errVault := s.vaultService.GetDecrypted(*r.VaultId)
			if errVault == nil {
				cred = &keeper.Credentials{Username: d.Username, Password: d.Secret}
			}
			return errVault
		})
	} else {
		check.Skip(utils.CheckStageVault, "no vault is configured")
	}
	var tlsConfig *tls.Config
	if r.Certs != nil {
		check.Run(utils.CheckStageCert, func() error {
			return s.certService.EnrichTLSConfig(&tlsConfig, r.Certs)
		})
	} else {
		check.Skip(utils.CheckStageCert, "tls is off")
	}
	check.RunAfter(utils.CheckStageConnect, func() error {
		_, _, errList := adapter.List(keeper.Request{Cluster: r.Cluster, Host: r.Host, Port: r.Port, TlsConfig: tlsConfig, Credentials: cred})
		return errList
	})
	return check.Steps(), nil
}

func (s *Service) KeeperConfigGet(r KeeperOneRequest) (any, int, error) {
	client, tlsConfig, cred, err := s.getKeeperAdapter(r.KeeperOptions)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}
	return client.Config(keeper.Request{Cluster: r.Cluster, Host: r.Host, Port: r.Port, Body: r.Body, TlsConfig: tlsConfig, Credentials: cred})
}

func (s *Service) KeeperConfigUpdate(r KeeperOneRequest) (any, int, error) {
	client, tlsConfig, cred, err := s.getKeeperAdapter(r.KeeperOptions)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}
	return client.ConfigUpdate(keeper.Request{Cluster: r.Cluster, Host: r.Host, Port: r.Port, Body: r.Body, TlsConfig: tlsConfig, Credentials: cred})
}

func (s *Service) KeeperSwitchover(r KeeperOneRequest) (*string, int, error) {
	client, tlsConfig, cred, err := s.getKeeperAdapter(r.KeeperOptions)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}
	return client.Switchover(keeper.Request{Cluster: r.Cluster, Host: r.Host, Port: r.Port, Body: r.Body, TlsConfig: tlsConfig, Credentials: cred})
}

func (s *Service) KeeperSwitchoverDelete(r KeeperOneRequest) (*string, int, error) {
	client, tlsConfig, cred, err := s.getKeeperAdapter(r.KeeperOptions)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}
	return client.DeleteSwitchover(keeper.Request{Cluster: r.Cluster, Host: r.Host, Port: r.Port, Body: r.Body, TlsConfig: tlsConfig, Credentials: cred})
}

func (s *Service) KeeperReinitialize(r KeeperOneRequest) (*string, int, error) {
	client, tlsConfig, cred, err := s.getKeeperAdapter(r.KeeperOptions)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}
	return client.Reinitialize(keeper.Request{Cluster: r.Cluster, Host: r.Host, Port: r.Port, Body: r.Body, TlsConfig: tlsConfig, Credentials: cred})
}

func (s *Service) KeeperRestart(r KeeperOneRequest) (*string, int, error) {
	client, tlsConfig, cred, err := s.getKeeperAdapter(r.KeeperOptions)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}
	return client.Restart(keeper.Request{Cluster: r.Cluster, Host: r.Host, Port: r.Port, Body: r.Body, TlsConfig: tlsConfig, Credentials: cred})
}

func (s *Service) KeeperRestartDelete(r KeeperOneRequest) (*string, int, error) {
	client, tlsConfig, cred, err := s.getKeeperAdapter(r.KeeperOptions)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}
	return client.DeleteRestart(keeper.Request{Cluster: r.Cluster, Host: r.Host, Port: r.Port, Body: r.Body, TlsConfig: tlsConfig, Credentials: cred})
}

func (s *Service) KeeperReload(r KeeperOneRequest) (*string, int, error) {
	client, tlsConfig, cred, err := s.getKeeperAdapter(r.KeeperOptions)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}
	return client.Reload(keeper.Request{Cluster: r.Cluster, Host: r.Host, Port: r.Port, Body: r.Body, TlsConfig: tlsConfig, Credentials: cred})
}

func (s *Service) KeeperFailover(r KeeperOneRequest) (*string, int, error) {
	client, tlsConfig, cred, err := s.getKeeperAdapter(r.KeeperOptions)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}
	return client.Failover(keeper.Request{Cluster: r.Cluster, Host: r.Host, Port: r.Port, Body: r.Body, TlsConfig: tlsConfig, Credentials: cred})
}

func (s *Service) KeeperActivate(r KeeperOneRequest) (*string, int, error) {
	client, tlsConfig, cred, err := s.getKeeperAdapter(r.KeeperOptions)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}
	return client.Activate(keeper.Request{Cluster: r.Cluster, Host: r.Host, Port: r.Port, Body: r.Body, TlsConfig: tlsConfig, Credentials: cred})
}

func (s *Service) KeeperPause(r KeeperOneRequest) (*string, int, error) {
	client, tlsConfig, cred, err := s.getKeeperAdapter(r.KeeperOptions)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}
	return client.Pause(keeper.Request{Cluster: r.Cluster, Host: r.Host, Port: r.Port, Body: r.Body, TlsConfig: tlsConfig, Credentials: cred})
}
