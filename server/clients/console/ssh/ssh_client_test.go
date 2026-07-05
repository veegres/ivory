package ssh

import (
	"crypto/ed25519"
	"crypto/rand"
	"testing"

	"golang.org/x/crypto/ssh"
)

type mockAddr struct{}

func (m mockAddr) Network() string { return "tcp" }
func (m mockAddr) String() string  { return "127.0.0.1:22" }

func TestHostKeyCallback(t *testing.T) {
	client := NewClient()
	hostname := "test-host"
	remote := mockAddr{}

	// Generate a key
	pub, _, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	key, err := ssh.NewPublicKey(pub)
	if err != nil {
		t.Fatal(err)
	}

	// First connection - should be trusted and cached
	err = client.hostKeyCallback(hostname, remote, key)
	if err != nil {
		t.Errorf("expected no error on first connection, got %v", err)
	}

	// Second connection with same key - should pass
	err = client.hostKeyCallback(hostname, remote, key)
	if err != nil {
		t.Errorf("expected no error on second connection with same key, got %v", err)
	}

	// Connection with different key - should fail
	pub2, _, _ := ed25519.GenerateKey(rand.Reader)
	key2, _ := ssh.NewPublicKey(pub2)
	err = client.hostKeyCallback(hostname, remote, key2)
	if err == nil {
		t.Error("expected error on host key mismatch, got nil")
	}

	// Different host with key2 - should be trusted for that host
	err = client.hostKeyCallback("other-host", remote, key2)
	if err != nil {
		t.Errorf("expected no error on first connection for other-host, got %v", err)
	}
}

func TestGenerateKey(t *testing.T) {
	client := NewClient()
	pubStr, prvStr, err := client.GenerateKey()
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if pubStr == "" {
		t.Error("expected public key, got empty string")
	}
	if prvStr == "" {
		t.Error("expected private key, got empty string")
	}

	// Validate private key
	_, err = ssh.NewSignerFromKey(ed25519.PrivateKey(prvStr))
	if err != nil {
		t.Errorf("failed to create signer from generated private key: %v", err)
	}

	// Validate public key
	_, _, _, _, err = ssh.ParseAuthorizedKey([]byte(pubStr))
	if err != nil {
		t.Errorf("failed to parse generated public key: %v", err)
	}
}
