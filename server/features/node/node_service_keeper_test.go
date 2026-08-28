package node

import (
	"errors"
	"ivory/core/utils"
	"ivory/plugins/keeper"
	"net/http"
	"testing"
)

// fakeKeeperAdapter lets the test control exactly what List() reports, in
// particular returning a usable Response alongside a non-nil error the way
// a real adapter does when it can still describe a node's state despite a
// connection problem (e.g. postgres starting up).
type fakeKeeperAdapter struct {
	keeper.Adapter
	listResponse []keeper.Response
	listStatus   int
	listErr      error
}

func (f *fakeKeeperAdapter) List(request keeper.Request) ([]keeper.Response, int, error) {
	return f.listResponse, f.listStatus, f.listErr
}

func TestKeeperNodeList_KeepsResponseAlongsideError(t *testing.T) {
	host := "db1"
	port := 5432
	state := keeper.StateStarting

	keeperRegistry := utils.NewRegistry[keeper.Plugin, keeper.Adapter]()
	keeperRegistry.Register("fake", &fakeKeeperAdapter{
		listResponse: []keeper.Response{{
			State:                state,
			Role:                 keeper.Unknown,
			DiscoveredHost:       &host,
			DiscoveredKeeperPort: &port,
		}},
		listStatus: http.StatusServiceUnavailable,
		listErr:    errors.New("the database system is starting up"),
	})
	s := NewService(nil, nil, keeperRegistry, nil, nil, nil, nil)

	responses, status, err := s.KeeperNodeList(KeeperOneRequest{
		KeeperConnection: KeeperConnection{Host: host, Port: port},
		KeeperOptions:    KeeperOptions{Plugin: "fake"},
	})

	if err == nil {
		t.Fatalf("expected the adapter error to still be returned, got nil")
	}
	if status != http.StatusServiceUnavailable {
		t.Errorf("expected status 503, got %d", status)
	}
	if len(responses) != 1 || responses[0].State != state {
		t.Fatalf("expected the degraded response to still be returned alongside the error, got %v", responses)
	}
}
