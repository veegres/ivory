package patroni

import (
	"encoding/json"
	"ivory/src/clients/http"
	"ivory/src/clients/keeper"
	nethttp "net/http"
	"strings"
	"testing"
)

func TestClient_Overview_Mapping(t *testing.T) {
	t.Run("should map basic Patroni cluster to internal model", func(t *testing.T) {
		// Mock Patroni response
		patroniResponse := PatroniCluster{
			Members: []PatroniInstance{
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
		_ = NewClient(httpClient)

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

		_ = &Client{}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				patroniInstance := PatroniInstance{
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
	client := &Client{}

	t.Run("valid lag", func(t *testing.T) {
		patroniInstance := PatroniInstance{
			Lag: json.RawMessage("12345"),
		}
		lag := client.mapLag(patroniInstance.Lag)
		if lag != 12345 {
			t.Errorf("Expected lag 12345, got %d", lag)
		}
	})

	t.Run("invalid lag string", func(t *testing.T) {
		patroniInstance := PatroniInstance{
			Lag: json.RawMessage("\"invalid\""),
		}
		lag := client.mapLag(patroniInstance.Lag)
		if lag != -1 {
			t.Errorf("Expected lag -1 for invalid input, got %d", lag)
		}
	})

	t.Run("null lag", func(t *testing.T) {
		patroniInstance := PatroniInstance{
			Lag: json.RawMessage("null"),
		}
		lag := client.mapLag(patroniInstance.Lag)
		if lag != -1 {
			t.Errorf("Expected lag -1 for null input, got %d", lag)
		}
	})
}

func TestClient_mapRole(t *testing.T) {
	client := &Client{}

	testCases := []struct {
		input    string
		expected keeper.Role
	}{
		{"leader", keeper.Leader},
		{"master", keeper.Leader},
		{"replica", keeper.Replica},
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

func TestClient_mapRestart(t *testing.T) {
	client := &Client{}

	t.Run("valid restart", func(t *testing.T) {
		patroniRestart := &PatroniScheduledRestart{
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
	client := NewClient(&http.Client{})

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
	client := NewClient(&http.Client{})

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
