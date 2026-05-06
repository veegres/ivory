package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"
)

func TestJSONRequest_Get(t *testing.T) {
	type Response struct {
		Message string `json:"message"`
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET request, got %s", r.Method)
		}
		if r.URL.Path != "/test" {
			t.Errorf("expected path /test, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(Response{Message: "success"})
	}))
	defer server.Close()

	serverUrl, _ := url.Parse(server.URL)
	host := serverUrl.Hostname()
	port, _ := strconv.Atoi(serverUrl.Port())

	client := NewClient()
	jsonReq := NewJSONRequest[Response](client)

	resp, status, err := jsonReq.Get(Request{
		Host: host,
		Port: port,
		Path: "/test",
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if status != http.StatusOK {
		t.Errorf("expected status 200, got %d", status)
	}
	if resp.Message != "success" {
		t.Errorf("expected message 'success', got '%s'", resp.Message)
	}
}

func TestJSONRequest_Post(t *testing.T) {
	type Payload struct {
		Data string `json:"data"`
	}
	type Response struct {
		Result string `json:"result"`
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST request, got %s", r.Method)
		}

		var p Payload
		if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
			t.Errorf("failed to decode payload: %v", err)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(Response{Result: "received: " + p.Data})
	}))
	defer server.Close()

	serverUrl, _ := url.Parse(server.URL)
	host := serverUrl.Hostname()
	port, _ := strconv.Atoi(serverUrl.Port())

	client := NewClient()
	jsonReq := NewJSONRequest[Response](client)

	resp, status, err := jsonReq.Post(Request{
		Host: host,
		Port: port,
		Path: "/post",
		Body: Payload{Data: "hello"},
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if status != http.StatusCreated {
		t.Errorf("expected status 201, got %d", status)
	}
	if resp.Result != "received: hello" {
		t.Errorf("expected result 'received: hello', got '%s'", resp.Result)
	}
}

func TestJSONRequest_BasicAuth(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, pass, ok := r.BasicAuth()
		if !ok || user != "admin" || pass != "secret" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`"ok"`))
	}))
	defer server.Close()

	serverUrl, _ := url.Parse(server.URL)
	host := serverUrl.Hostname()
	port, _ := strconv.Atoi(serverUrl.Port())

	client := NewClient()
	jsonReq := NewJSONRequest[string](client)

	_, status, err := jsonReq.Get(Request{
		Host: host,
		Port: port,
		Path: "/auth",
		Credentials: &Credentials{
			Username: "admin",
			Password: "secret",
		},
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if status != http.StatusOK {
		t.Errorf("expected status 200, got %d", status)
	}
}

func TestJSONRequest_ErrorHandling(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error": "something went wrong"}`))
	}))
	defer server.Close()

	serverUrl, _ := url.Parse(server.URL)
	host := serverUrl.Hostname()
	port, _ := strconv.Atoi(serverUrl.Port())

	client := NewClient()
	jsonReq := NewJSONRequest[any](client)

	_, status, err := jsonReq.Get(Request{
		Host: host,
		Port: port,
		Path: "/error",
	})

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if status != http.StatusInternalServerError {
		t.Errorf("expected status 500, got %d", status)
	}
}

func TestJSONRequest_StringResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("raw string response"))
	}))
	defer server.Close()

	serverUrl, _ := url.Parse(server.URL)
	host := serverUrl.Hostname()
	port, _ := strconv.Atoi(serverUrl.Port())

	client := NewClient()
	jsonReq := NewJSONRequest[string](client)

	resp, status, err := jsonReq.Get(Request{
		Host: host,
		Port: port,
		Path: "/string",
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if status != http.StatusOK {
		t.Errorf("expected status 200, got %d", status)
	}
	if *resp != "raw string response" {
		t.Errorf("expected 'raw string response', got '%s'", *resp)
	}
}

func TestNewClient_HostEmpty(t *testing.T) {
	client := NewClient()
	_, err := client.Send(Request{Host: ""})
	if err != ErrHostEmpty {
		t.Errorf("expected ErrHostEmpty, got %v", err)
	}
}

func TestJSONRequest_Methods(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(fmt.Sprintf(`"%s"`, r.Method)))
	}))
	defer server.Close()

	serverUrl, _ := url.Parse(server.URL)
	host := serverUrl.Hostname()
	port, _ := strconv.Atoi(serverUrl.Port())

	client := NewClient()

	t.Run("PUT", func(t *testing.T) {
		resp, _, _ := NewJSONRequest[string](client).Put(Request{Host: host, Port: port})
		if *resp != "PUT" {
			t.Errorf("expected PUT, got %s", *resp)
		}
	})

	t.Run("PATCH", func(t *testing.T) {
		resp, _, _ := NewJSONRequest[string](client).Patch(Request{Host: host, Port: port})
		if *resp != "PATCH" {
			t.Errorf("expected PATCH, got %s", *resp)
		}
	})

	t.Run("DELETE", func(t *testing.T) {
		resp, _, _ := NewJSONRequest[string](client).Delete(Request{Host: host, Port: port})
		if *resp != "DELETE" {
			t.Errorf("expected DELETE, got %s", *resp)
		}
	})
}
