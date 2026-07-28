package ssh

import (
	"crypto/ed25519"
	"crypto/rand"
	"fmt"
	"net"
	"strconv"
	"sync/atomic"
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

// startEchoSSHServer starts a fake SSH server that replies to any "exec"
// request with "ok\n" and a zero exit status, enough to exercise Command's
// full Start/Wait path against a real (loopback) SSH connection.
func startEchoSSHServer(t *testing.T) net.Listener {
	t.Helper()
	serverConfig := &ssh.ServerConfig{
		PublicKeyCallback: func(conn ssh.ConnMetadata, key ssh.PublicKey) (*ssh.Permissions, error) {
			return nil, nil
		},
	}
	_, serverPrv, _ := ed25519.GenerateKey(rand.Reader)
	signer, _ := ssh.NewSignerFromKey(serverPrv)
	serverConfig.AddHostKey(signer)

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}

	go func() {
		for {
			nConn, err := listener.Accept()
			if err != nil {
				return
			}
			go func(c net.Conn) {
				defer c.Close()
				sConn, chans, reqs, err := ssh.NewServerConn(c, serverConfig)
				if err != nil {
					return
				}
				defer sConn.Close()
				go ssh.DiscardRequests(reqs)
				for newChannel := range chans {
					if newChannel.ChannelType() != "session" {
						newChannel.Reject(ssh.UnknownChannelType, "unknown channel type")
						continue
					}
					channel, requests, err := newChannel.Accept()
					if err != nil {
						continue
					}
					go func(in <-chan *ssh.Request) {
						defer channel.Close()
						for req := range in {
							switch req.Type {
							case "pty-req":
								req.Reply(true, nil)
							case "exec":
								req.Reply(true, nil)
								channel.Write([]byte("ok\n"))
								channel.SendRequest("exit-status", false, ssh.Marshal(struct{ Status uint32 }{0}))
								return
							default:
								req.Reply(false, nil)
							}
						}
					}(requests)
				}
			}(nConn)
		}
	}()

	return listener
}

// TestClient_DialCachedReusesLiveConnection verifies dialCached's liveness
// probe (a lightweight SendRequest) accepts a healthy cached connection
// instead of redialing, so a second command against the same connection
// details reuses the same *ssh.Client.
func TestClient_DialCachedReusesLiveConnection(t *testing.T) {
	listener := startEchoSSHServer(t)
	defer listener.Close()

	host, portStr, _ := net.SplitHostPort(listener.Addr().String())
	port, _ := strconv.Atoi(portStr)

	_, clientPrv, _ := ed25519.GenerateKey(rand.Reader)
	client := NewClient()

	var dialCount int32
	client.dial = func(network, addr string, config *ssh.ClientConfig) (*ssh.Client, error) {
		atomic.AddInt32(&dialCount, 1)
		return ssh.Dial(network, addr, config)
	}

	conn := Connection{Host: host, Port: port, Username: "test", PrivateKey: &clientPrv}

	if _, err := client.Command(conn, "echo hi").Execute(); err != nil {
		t.Fatalf("first Execute failed: %v", err)
	}
	if _, err := client.Command(conn, "echo hi").Execute(); err != nil {
		t.Fatalf("second Execute failed: %v", err)
	}

	if got := atomic.LoadInt32(&dialCount); got != 1 {
		t.Fatalf("expected exactly one dial (the second command should reuse the cached, live connection), got %d", got)
	}
}

// TestClient_DialCachedRedialsDeadConnection verifies that once the cached
// connection is actually closed, dialCached's liveness probe fails and a
// fresh connection is dialed instead of reusing the dead one.
func TestClient_DialCachedRedialsDeadConnection(t *testing.T) {
	listener := startEchoSSHServer(t)
	defer listener.Close()

	host, portStr, _ := net.SplitHostPort(listener.Addr().String())
	port, _ := strconv.Atoi(portStr)

	_, clientPrv, _ := ed25519.GenerateKey(rand.Reader)
	client := NewClient()

	var dialCount int32
	client.dial = func(network, addr string, config *ssh.ClientConfig) (*ssh.Client, error) {
		atomic.AddInt32(&dialCount, 1)
		return ssh.Dial(network, addr, config)
	}

	conn := Connection{Host: host, Port: port, Username: "test", PrivateKey: &clientPrv}

	if _, err := client.Command(conn, "echo hi").Execute(); err != nil {
		t.Fatalf("first Execute failed: %v", err)
	}

	// Kill the cached connection out from under the client.
	key := fmt.Sprintf("%s|%s|%s|%s", net.JoinHostPort(host, portStr), "test", "tcp", credentialHash("", &clientPrv))
	client.mu.RLock()
	cached := client.cachedClients[key]
	client.mu.RUnlock()
	if cached == nil {
		t.Fatal("expected a cached client after the first command")
	}
	if err := cached.Close(); err != nil {
		t.Fatalf("failed to close cached connection: %v", err)
	}

	if _, err := client.Command(conn, "echo hi").Execute(); err != nil {
		t.Fatalf("second Execute failed: %v", err)
	}

	if got := atomic.LoadInt32(&dialCount); got != 2 {
		t.Fatalf("expected a redial once the cached connection died, got %d dial(s)", got)
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
