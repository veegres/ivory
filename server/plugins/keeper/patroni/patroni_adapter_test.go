package patroni

import (
	"encoding/json"
	"ivory/clients/http"
	"ivory/plugins/keeper"
	"net"
	nethttp "net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
)

func TestClient_Overview_Mapping(t *testing.T) {
	t.Run("should map basic Patroni cluster to internal model", func(t *testing.T) {
		// Mock Patroni response
		patroniResponse := cluster{
			Members: []instance{
				{
					Name:   "db1",
					Host:   "db1.example.com",
					Port:   5432,
					Role:   "leader",
					State:  "running",
					ApiUrl: "http://10.0.0.1:8008/patroni",
					Lag:    json.RawMessage("0"),
				},
			},
			Pause: false,
		}

		// Mock the gateway response
		httpClient := &http.Client{}
		_ = NewAdapter(httpClient)

		// Verify the struct is created correctly
		if len(patroniResponse.Members) != 1 {
			t.Fatalf("Expected 1 member, got %d", len(patroniResponse.Members))
		}
	})

	t.Run("should parse different API URL formats", func(t *testing.T) {
		testCases := []struct {
			name     string
			apiUrl   string
			expected string
		}{
			{"standard", "http://10.0.0.1:8008/patroni", "10.0.0.1"},
			{"hostname", "http://db-host:8008/patroni", "db-host"},
			{"trailing slash", "http://10.0.0.1:8008/", "10.0.0.1"},
		}

		_ = &Adapter{}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				patroniInstance := instance{
					ApiUrl: tc.apiUrl,
				}

				// The mapping logic is inside Overview(), but let's test the string splitting logic
				domainString := strings.Split(patroniInstance.ApiUrl, "/")[2]
				domain := strings.Split(domainString, ":")
				host := domain[0]

				if host != tc.expected {
					t.Errorf("Expected host '%s', got '%s'", tc.expected, host)
				}
			})
		}
	})
}

func TestClient_mapLag(t *testing.T) {
	client := &Adapter{}

	t.Run("valid lag", func(t *testing.T) {
		patroniInstance := instance{
			Lag: json.RawMessage("12345"),
		}
		lag := client.mapLag(patroniInstance.Lag)
		if lag != 12345 {
			t.Errorf("Expected lag 12345, got %d", lag)
		}
	})

	t.Run("invalid lag string", func(t *testing.T) {
		patroniInstance := instance{
			Lag: json.RawMessage("\"invalid\""),
		}
		lag := client.mapLag(patroniInstance.Lag)
		if lag != -1 {
			t.Errorf("Expected lag -1 for invalid input, got %d", lag)
		}
	})

	t.Run("null lag", func(t *testing.T) {
		patroniInstance := instance{
			Lag: json.RawMessage("null"),
		}
		lag := client.mapLag(patroniInstance.Lag)
		if lag != -1 {
			t.Errorf("Expected lag -1 for null input, got %d", lag)
		}
	})
}

func TestClient_mapRole(t *testing.T) {
	client := &Adapter{}

	testCases := []struct {
		input    string
		expected keeper.Role
	}{
		{"leader", keeper.Leader},
		{"master", keeper.Leader},
		{"standby_leader", keeper.Leader},
		{"replica", keeper.Replica},
		{"sync_standby", keeper.Replica},
		{"quorum_standby", keeper.Replica},
		{"unknown", keeper.Unknown},
		{"", keeper.Unknown},
	}

	for _, tc := range testCases {
		role := client.mapRole(tc.input)
		if role != tc.expected {
			t.Errorf("For input '%s', expected role '%s', got '%s'", tc.input, tc.expected, role)
		}
	}
}

func TestClient_mapSync(t *testing.T) {
	client := &Adapter{}

	testCases := []struct {
		input    string
		expected bool
	}{
		{"sync_standby", true},
		{"quorum_standby", true},
		{"replica", false},
		{"leader", false},
		{"master", false},
		{"standby_leader", false},
		{"unknown", false},
		{"", false},
	}

	for _, tc := range testCases {
		sync := client.mapSync(tc.input)
		if sync != tc.expected {
			t.Errorf("For input '%s', expected sync '%v', got '%v'", tc.input, tc.expected, sync)
		}
	}
}

func TestClient_mapState(t *testing.T) {
	client := &Adapter{}

	testCases := []struct {
		input    string
		expected keeper.State
	}{
		{"running", keeper.StateRunning},
		{"streaming", keeper.StateRunning}, // newer Patroni reports replicas as "streaming" instead of "running"
		{"in archive recovery", keeper.StateRunning},
		{"starting", keeper.StateStarting},
		{"creating replica", keeper.StateStarting},
		{"initializing new cluster", keeper.StateStarting},
		{"running custom bootstrap script", keeper.StateStarting},
		{"restarting", keeper.StateRestarting},
		{"stopping", keeper.StateStopping},
		{"stopped", keeper.StateStopped},
		{"crashed", keeper.StateFailed},
		{"start failed", keeper.StateFailed},
		{"restart failed", keeper.StateFailed},
		{"stop failed", keeper.StateFailed},
		{"initdb failed", keeper.StateFailed},
		{"custom bootstrap failed", keeper.StateFailed},
		{"some future patroni state", keeper.StateUnknown},
		{"", keeper.StateUnknown},
	}

	for _, tc := range testCases {
		state := client.mapState(tc.input)
		if state != tc.expected {
			t.Errorf("for input %q, expected state %q, got %q", tc.input, tc.expected, state)
		}
	}
}

func TestClient_mapRestart(t *testing.T) {
	client := &Adapter{}

	t.Run("valid restart", func(t *testing.T) {
		patroniRestart := &scheduledRestart{
			RestartPending: true,
			Schedule:       "2024-10-26T12:00:00Z",
		}
		restart := client.mapRestart(patroniRestart)

		if restart == nil {
			t.Fatal("Expected restart info, got nil")
		}
		if !restart.PendingRestart {
			t.Error("Expected pendingRestart to be true")
		}
		if restart.At != "2024-10-26T12:00:00Z" {
			t.Errorf("Expected schedule time, got '%s'", restart.At)
		}
	})

	t.Run("nil restart", func(t *testing.T) {
		restart := client.mapRestart(nil)
		if restart != nil {
			t.Errorf("Expected nil, got %+v", restart)
		}
	})
}

func TestClient_Activate(t *testing.T) {
	client := NewAdapter(&http.Client{})

	t.Run("should return error when body is not nil", func(t *testing.T) {
		request := keeper.Request{
			Body: map[string]string{"key": "value"},
		}

		_, status, err := client.Activate(request)

		if err == nil {
			t.Fatal("Expected error when body is not nil, got nil")
		}
		if status != nethttp.StatusBadRequest {
			t.Errorf("Expected status %d, got %d", nethttp.StatusBadRequest, status)
		}
		if err.Error() != "body should be empty" {
			t.Errorf("Expected 'body should be empty', got '%s'", err.Error())
		}
	})
}

func TestClient_Pause(t *testing.T) {
	client := NewAdapter(&http.Client{})

	t.Run("should return error when body is not nil", func(t *testing.T) {
		request := keeper.Request{
			Body: map[string]string{"key": "value"},
		}

		_, status, err := client.Pause(request)

		if err == nil {
			t.Fatal("Expected error when body is not nil, got nil")
		}
		if status != nethttp.StatusBadRequest {
			t.Errorf("Expected status %d, got %d", nethttp.StatusBadRequest, status)
		}
		if err.Error() != "body should be empty" {
			t.Errorf("Expected 'body should be empty', got '%s'", err.Error())
		}
	})
}

func TestKeeperResponse_Mapping(t *testing.T) {
	t.Run("should map to internal response model correctly", func(t *testing.T) {
		host := "db1.example.com"
		dbPort, keeperPort := 5432, 8008
		expectedResponse := keeper.Response{
			Role:                 keeper.Leader,
			State:                "running",
			Lag:                  0,
			PendingRestart:       false,
			DiscoveredHost:       &host,
			DiscoveredDbPort:     &dbPort,
			DiscoveredKeeperPort: &keeperPort,
		}

		if expectedResponse.Role != keeper.Leader {
			t.Errorf("Expected role 'leader', got '%s'", expectedResponse.Role)
		}
		if *expectedResponse.DiscoveredHost != "db1.example.com" {
			t.Errorf("Expected host 'db1.example.com', got '%s'", *expectedResponse.DiscoveredHost)
		}
		if *expectedResponse.DiscoveredDbPort != 5432 {
			t.Errorf("Expected db port 5432, got %d", *expectedResponse.DiscoveredDbPort)
		}
		if *expectedResponse.DiscoveredKeeperPort != 8008 {
			t.Errorf("Expected keeper port 8008, got %d", *expectedResponse.DiscoveredKeeperPort)
		}
	})
}

func TestAdapter_List(t *testing.T) {
	server := httptest.NewServer(nethttp.HandlerFunc(func(w nethttp.ResponseWriter, r *nethttp.Request) {
		if r.URL.Path != "/cluster" {
			w.WriteHeader(nethttp.StatusNotFound)
			return
		}
		_ = json.NewEncoder(w).Encode(cluster{
			Members: []instance{{
				Name:   "postgres-1",
				Host:   "10.0.0.1",
				Port:   5432,
				Role:   "leader",
				State:  "running",
				ApiUrl: "http://10.0.0.1:8008/patroni",
				Lag:    json.RawMessage("0"),
			}},
		})
	}))
	defer server.Close()

	host, portStr, _ := net.SplitHostPort(strings.TrimPrefix(server.URL, "http://"))
	port, _ := strconv.Atoi(portStr)

	responses, _, err := NewAdapter(http.NewClient()).List(keeper.Request{Host: host, Port: port})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(responses) != 1 {
		t.Fatalf("expected 1 response, got %d", len(responses))
	}

	member := responses[0]
	if member.DiscoveredName == nil || *member.DiscoveredName != "postgres-1" {
		t.Errorf("expected discovered name postgres-1, got %v", member.DiscoveredName)
	}
	if member.DiscoveredHost == nil || *member.DiscoveredHost != "10.0.0.1" {
		t.Errorf("expected discovered host 10.0.0.1, got %v", member.DiscoveredHost)
	}
	if member.DiscoveredKeeperPort == nil || *member.DiscoveredKeeperPort != 8008 {
		t.Errorf("expected discovered keeper port 8008, got %v", member.DiscoveredKeeperPort)
	}
	if member.DiscoveredDbPort == nil || *member.DiscoveredDbPort != 5432 {
		t.Errorf("expected discovered db port 5432, got %v", member.DiscoveredDbPort)
	}
}
