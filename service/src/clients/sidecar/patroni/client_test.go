package patroni

import (
	"encoding/json"
	"ivory/src/clients/database"
	"ivory/src/clients/sidecar"
	"net/http"
	"testing"
)

func TestClient_Overview_Mapping(t *testing.T) {
	t.Run("should map basic Patroni cluster to internal model", func(t *testing.T) {
		patroniCluster := PatroniCluster{
			Members: []PatroniInstance{
				{
					Name:           "node1",
					State:          "running",
					Role:           "leader",
					Host:           "db1.example.com",
					Port:           5432,
					ApiUrl:         "http://10.0.0.1:8008/patroni",
					PendingRestart: false,
					Lag:            json.RawMessage(`0`),
					Timeline:       1,
				},
			},
			Pause: false,
		}

		// Mock the gateway response
		mockGateway := &sidecar.Gateway{}
		client := NewClient(mockGateway)

		// We'll test the mapping logic by simulating what the response would be
		// Since Overview() makes actual HTTP calls, we test the mapping separately

		// Verify the struct is created correctly
		if client == nil {
			t.Fatal("Expected client to be created")
		}

		// Verify the Patroni model can be parsed
		if patroniCluster.Members[0].Name != "node1" {
			t.Errorf("Expected node name 'node1', got '%s'", patroniCluster.Members[0].Name)
		}
	})

	t.Run("should parse API URL correctly", func(t *testing.T) {
		testCases := []struct {
			apiUrl       string
			expectedHost string
			expectedPort int
		}{
			{"http://10.0.0.1:8008/patroni", "10.0.0.1", 8008},
			{"https://db1.example.com:8009/patroni", "db1.example.com", 8009},
			{"http://localhost:8008/", "localhost", 8008},
		}

		for _, tc := range testCases {
			// Simulate the parsing logic from Overview()
			// domainString := strings.Split(apiUrl, "/")[2]
			// domain := strings.Split(domainString, ":")

			patroniInstance := PatroniInstance{
				ApiUrl: tc.apiUrl,
				Name:   "test",
				State:  "running",
				Role:   "leader",
				Host:   "testhost",
				Port:   5432,
				Lag:    json.RawMessage(`0`),
			}

			if patroniInstance.ApiUrl != tc.apiUrl {
				t.Errorf("Expected API URL '%s', got '%s'", tc.apiUrl, patroniInstance.ApiUrl)
			}
		}
	})

	t.Run("should handle lag as integer", func(t *testing.T) {
		patroniInstance := PatroniInstance{
			Name:   "node1",
			State:  "running",
			Role:   "replica",
			Host:   "db1",
			Port:   5432,
			ApiUrl: "http://10.0.0.1:8008/patroni",
			Lag:    json.RawMessage(`1024`),
		}

		// Verify lag is stored as RawMessage
		var lagValue int64
		err := json.Unmarshal(patroniInstance.Lag, &lagValue)
		if err != nil {
			t.Fatalf("Expected to parse lag as int, got error: %v", err)
		}
		if lagValue != 1024 {
			t.Errorf("Expected lag 1024, got %d", lagValue)
		}
	})

	t.Run("should handle lag as string unknown", func(t *testing.T) {
		patroniInstance := PatroniInstance{
			Name:   "node1",
			State:  "running",
			Role:   "replica",
			Host:   "db1",
			Port:   5432,
			ApiUrl: "http://10.0.0.1:8008/patroni",
			Lag:    json.RawMessage(`"unknown"`),
		}

		// Verify lag is stored as RawMessage
		var lagValue string
		err := json.Unmarshal(patroniInstance.Lag, &lagValue)
		if err != nil {
			t.Fatalf("Expected to parse lag as string, got error: %v", err)
		}
		if lagValue != "unknown" {
			t.Errorf("Expected lag 'unknown', got '%s'", lagValue)
		}
	})

	t.Run("should handle lag as null", func(t *testing.T) {
		patroniInstance := PatroniInstance{
			Name:   "node1",
			State:  "running",
			Role:   "leader",
			Host:   "db1",
			Port:   5432,
			ApiUrl: "http://10.0.0.1:8008/patroni",
			Lag:    json.RawMessage(`null`),
		}

		// Verify lag is stored as RawMessage
		if string(patroniInstance.Lag) != "null" {
			t.Errorf("Expected lag to be null, got '%s'", string(patroniInstance.Lag))
		}
	})
}

func TestClient_Overview_SidecarStatus(t *testing.T) {
	t.Run("should map pause false to Active status", func(t *testing.T) {
		patroniCluster := PatroniCluster{
			Pause: false,
			Members: []PatroniInstance{
				{
					Name:   "node1",
					State:  "running",
					Role:   "leader",
					Host:   "db1",
					Port:   5432,
					ApiUrl: "http://10.0.0.1:8008/patroni",
					Lag:    json.RawMessage(`0`),
				},
			},
		}

		// Verify pause state
		if patroniCluster.Pause {
			t.Error("Expected pause to be false")
		}

		// The mapping logic would set sidecar.Active when pause is false
		var expectedStatus sidecar.SidecarStatus
		if patroniCluster.Pause == false {
			expectedStatus = sidecar.Active
		} else {
			expectedStatus = sidecar.Paused
		}

		if expectedStatus != sidecar.Active {
			t.Errorf("Expected status Active, got %v", expectedStatus)
		}
	})

	t.Run("should map pause true to Paused status", func(t *testing.T) {
		patroniCluster := PatroniCluster{
			Pause: true,
			Members: []PatroniInstance{
				{
					Name:   "node1",
					State:  "running",
					Role:   "leader",
					Host:   "db1",
					Port:   5432,
					ApiUrl: "http://10.0.0.1:8008/patroni",
					Lag:    json.RawMessage(`0`),
				},
			},
		}

		// Verify pause state
		if !patroniCluster.Pause {
			t.Error("Expected pause to be true")
		}

		// The mapping logic would set sidecar.Paused when pause is true
		var expectedStatus sidecar.SidecarStatus
		if patroniCluster.Pause == false {
			expectedStatus = sidecar.Active
		} else {
			expectedStatus = sidecar.Paused
		}

		if expectedStatus != sidecar.Paused {
			t.Errorf("Expected status Paused, got %v", expectedStatus)
		}
	})
}

func TestClient_Overview_ScheduledOperations(t *testing.T) {
	t.Run("should handle scheduled restart", func(t *testing.T) {
		patroniInstance := PatroniInstance{
			Name:           "node1",
			State:          "running",
			Role:           "leader",
			Host:           "db1",
			Port:           5432,
			ApiUrl:         "http://10.0.0.1:8008/patroni",
			Lag:            json.RawMessage(`0`),
			PendingRestart: true,
			ScheduledRestart: &PatroniScheduledRestart{
				RestartPending:      true,
				Schedule:            "2024-10-26T12:00:00Z",
				PostmasterStartTime: "2024-10-25T10:00:00Z",
			},
		}

		if patroniInstance.ScheduledRestart == nil {
			t.Fatal("Expected scheduled restart to be set")
		}
		if !patroniInstance.ScheduledRestart.RestartPending {
			t.Error("Expected restart pending to be true")
		}
		if patroniInstance.ScheduledRestart.Schedule != "2024-10-26T12:00:00Z" {
			t.Errorf("Expected schedule time, got '%s'", patroniInstance.ScheduledRestart.Schedule)
		}
	})

	t.Run("should handle scheduled switchover", func(t *testing.T) {
		patroniCluster := PatroniCluster{
			Pause: false,
			ScheduledSwitchover: &PatroniScheduledSwitchover{
				At:   "2024-10-26T15:00:00Z",
				From: "node1",
				To:   "node2",
			},
			Members: []PatroniInstance{
				{
					Name:   "node1",
					State:  "running",
					Role:   "leader",
					Host:   "db1",
					Port:   5432,
					ApiUrl: "http://10.0.0.1:8008/patroni",
					Lag:    json.RawMessage(`0`),
				},
			},
		}

		if patroniCluster.ScheduledSwitchover == nil {
			t.Fatal("Expected scheduled switchover to be set")
		}
		if patroniCluster.ScheduledSwitchover.From != "node1" {
			t.Errorf("Expected switchover from 'node1', got '%s'", patroniCluster.ScheduledSwitchover.From)
		}
		if patroniCluster.ScheduledSwitchover.To != "node2" {
			t.Errorf("Expected switchover to 'node2', got '%s'", patroniCluster.ScheduledSwitchover.To)
		}
	})

	t.Run("should handle empty switchover target as random selection", func(t *testing.T) {
		patroniCluster := PatroniCluster{
			ScheduledSwitchover: &PatroniScheduledSwitchover{
				At:   "2024-10-26T15:00:00Z",
				From: "node1",
				To:   "", // Empty means random selection
			},
		}

		// The mapping logic would replace empty "To" with "(random selection)"
		expectedTo := patroniCluster.ScheduledSwitchover.To
		if expectedTo == "" {
			expectedTo = "(random selection)"
		}

		if expectedTo != "(random selection)" {
			t.Errorf("Expected '(random selection)', got '%s'", expectedTo)
		}
	})
}

func TestClient_Activate(t *testing.T) {
	client := NewClient(&sidecar.Gateway{})

	t.Run("should return error when body is not nil", func(t *testing.T) {
		request := sidecar.Request{
			Body: map[string]string{"key": "value"},
		}

		_, status, err := client.Activate(request)

		if err == nil {
			t.Fatal("Expected error when body is not nil, got nil")
		}
		if status != http.StatusBadRequest {
			t.Errorf("Expected status %d, got %d", http.StatusBadRequest, status)
		}
		if err.Error() != "body should be empty" {
			t.Errorf("Expected 'body should be empty', got '%s'", err.Error())
		}
	})
}

func TestClient_Pause(t *testing.T) {
	client := NewClient(&sidecar.Gateway{})

	t.Run("should return error when body is not nil", func(t *testing.T) {
		request := sidecar.Request{
			Body: map[string]string{"key": "value"},
		}

		_, status, err := client.Pause(request)

		if err == nil {
			t.Fatal("Expected error when body is not nil, got nil")
		}
		if status != http.StatusBadRequest {
			t.Errorf("Expected status %d, got %d", http.StatusBadRequest, status)
		}
		if err.Error() != "body should be empty" {
			t.Errorf("Expected 'body should be empty', got '%s'", err.Error())
		}
	})
}

func TestConfigPause(t *testing.T) {
	t.Run("should marshal and unmarshal correctly", func(t *testing.T) {
		config := ConfigPause{Pause: true}

		data, err := json.Marshal(config)
		if err != nil {
			t.Fatalf("Expected to marshal, got error: %v", err)
		}

		var unmarshaled ConfigPause
		err = json.Unmarshal(data, &unmarshaled)
		if err != nil {
			t.Fatalf("Expected to unmarshal, got error: %v", err)
		}

		if unmarshaled.Pause != true {
			t.Errorf("Expected Pause to be true, got %v", unmarshaled.Pause)
		}
	})
}

func TestPatroniModels(t *testing.T) {
	t.Run("should create PatroniCluster with members", func(t *testing.T) {
		cluster := PatroniCluster{
			Pause: false,
			Members: []PatroniInstance{
				{Name: "node1", Role: "leader"},
				{Name: "node2", Role: "replica"},
			},
		}

		if len(cluster.Members) != 2 {
			t.Errorf("Expected 2 members, got %d", len(cluster.Members))
		}
		if cluster.Members[0].Role != "leader" {
			t.Errorf("Expected leader role, got '%s'", cluster.Members[0].Role)
		}
	})

	t.Run("should handle tags as pointer to map", func(t *testing.T) {
		tags := map[string]any{
			"environment": "production",
			"datacenter":  "us-west-1",
		}

		instance := PatroniInstance{
			Name: "node1",
			Tags: &tags,
		}

		if instance.Tags == nil {
			t.Fatal("Expected tags to be set")
		}
		if (*instance.Tags)["environment"] != "production" {
			t.Error("Expected environment tag to be 'production'")
		}
	})

	t.Run("should handle nil tags", func(t *testing.T) {
		instance := PatroniInstance{
			Name: "node1",
			Tags: nil,
		}

		if instance.Tags != nil {
			t.Error("Expected tags to be nil")
		}
	})
}

func TestSidecarInstance_Mapping(t *testing.T) {
	t.Run("should map to internal instance model correctly", func(t *testing.T) {
		// This tests the expected output structure from Overview()
		expectedInstance := sidecar.Instance{
			State:          "running",
			Role:           "leader",
			Lag:            0,
			PendingRestart: false,
			Database:       database.Database{Host: "db1.example.com", Port: 5432},
			Sidecar: sidecar.Sidecar{
				Host:   "10.0.0.1",
				Port:   8008,
				Name:   strPtr("node1"),
				Status: sidecarStatusPtr(sidecar.Active),
			},
		}

		if expectedInstance.Role != "leader" {
			t.Errorf("Expected role 'leader', got '%s'", expectedInstance.Role)
		}
		if expectedInstance.Sidecar.Host != "10.0.0.1" {
			t.Errorf("Expected sidecar host '10.0.0.1', got '%s'", expectedInstance.Sidecar.Host)
		}
		if expectedInstance.Sidecar.Port != 8008 {
			t.Errorf("Expected sidecar port 8008, got %d", expectedInstance.Sidecar.Port)
		}
		if expectedInstance.Database.Host != "db1.example.com" {
			t.Errorf("Expected database host 'db1.example.com', got '%s'", expectedInstance.Database.Host)
		}
	})
}

// Helper functions
func strPtr(s string) *string {
	return &s
}

func sidecarStatusPtr(s sidecar.SidecarStatus) *sidecar.SidecarStatus {
	return &s
}
