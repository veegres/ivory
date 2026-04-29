package ssh

import (
	"crypto/ed25519"
	"crypto/rand"
	"errors"
	"io"
	"net"
	"strconv"
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

func TestGetDialAddress(t *testing.T) {
	client := NewClient()
	tests := []struct {
		name       string
		connection Connection
		want       string
		wantErr    bool
	}{
		{
			name:       "simple host and port",
			connection: Connection{Host: "localhost", Port: 22},
			want:       "localhost:22",
		},
		{
			name:       "host with protocol",
			connection: Connection{Host: "ssh://localhost", Port: 22},
			want:       "localhost:22",
		},
		{
			name:       "host with protocol and port",
			connection: Connection{Host: "ssh://localhost:2222", Port: 22},
			want:       "localhost:2222",
		},
		{
			name:       "host with port",
			connection: Connection{Host: "localhost:2222", Port: 22},
			want:       "localhost:2222",
		},
		{
			name:       "empty host",
			connection: Connection{Host: "", Port: 22},
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := client.getDialAddress(tt.connection)
			if (err != nil) != tt.wantErr {
				t.Errorf("getDialAddress() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("getDialAddress() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExecute_Errors(t *testing.T) {
	client := NewClient()
	_, prv, _ := ed25519.GenerateKey(rand.Reader)
	conn := Connection{Host: "localhost", Port: 22, Username: "test", PrivateKey: []byte(prv)}

	t.Run("empty command", func(t *testing.T) {
		_, err := client.Execute(conn, "")
		if !errors.Is(err, ErrCommandEmpty) {
			t.Errorf("expected ErrCommandEmpty, got %v", err)
		}
	})

	t.Run("invalid private key", func(t *testing.T) {
		invalidConn := conn
		invalidConn.PrivateKey = make([]byte, 64)
		_, err := client.Execute(invalidConn, "ls")
		if err == nil {
			t.Error("expected error for invalid private key, got nil")
		}
	})

	t.Run("dial error", func(t *testing.T) {
		dialErr := errors.New("dial failed")
		client.dial = func(network, addr string, config *ssh.ClientConfig) (*ssh.Client, error) {
			return nil, dialErr
		}
		_, err := client.Execute(conn, "ls")
		if !errors.Is(err, dialErr) {
			t.Errorf("expected %v, got %v", dialErr, err)
		}
	})
}

func TestExecute_Success(t *testing.T) {
	// Server setup
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
	defer listener.Close()

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
							if req.Type == "exec" {
								payload := req.Payload
								if len(payload) < 4 {
									req.Reply(false, nil)
									continue
								}
								cmdLen := int(payload[3]) | int(payload[2])<<8 | int(payload[1])<<16 | int(payload[0])<<24
								cmd := string(payload[4 : 4+cmdLen])

								if cmd == "echo hello" {
									io.WriteString(channel, "hello")
									req.Reply(true, nil)
									channel.SendRequest("exit-status", false, ssh.Marshal(struct{ Status uint32 }{0}))
									return
								} else if cmd == "fail" {
									io.WriteString(channel.Stderr(), "error")
									req.Reply(true, nil)
									channel.SendRequest("exit-status", false, ssh.Marshal(struct{ Status uint32 }{1}))
									return
								} else {
									req.Reply(false, nil)
								}
							} else {
								req.Reply(false, nil)
							}
						}
					}(requests)
				}
			}(nConn)
		}
	}()

	addr := listener.Addr().String()
	host, portStr, _ := net.SplitHostPort(addr)
	port, _ := strconv.Atoi(portStr)

	_, clientPrv, _ := ed25519.GenerateKey(rand.Reader)
	client := NewClient()
	conn := Connection{
		Host:       host,
		Port:       port,
		Username:   "test",
		PrivateKey: []byte(clientPrv),
	}

	t.Run("command success", func(t *testing.T) {
		res, err := client.Execute(conn, "echo hello")
		if err != nil {
			t.Fatalf("Execute failed: %v", err)
		}
		if res.ExitCode != 0 {
			t.Errorf("expected exit code 0, got %d", res.ExitCode)
		}
		if res.Stdout != "hello" {
			t.Errorf("expected stdout 'hello', got '%s'", res.Stdout)
		}
	})

	t.Run("command failure", func(t *testing.T) {
		res, err := client.Execute(conn, "fail")
		if err != nil {
			t.Fatalf("Execute failed: %v", err)
		}
		if res.ExitCode != 1 {
			t.Errorf("expected exit code 1, got %d", res.ExitCode)
		}
		if res.Stderr != "error" {
			t.Errorf("expected stderr 'error', got '%s'", res.Stderr)
		}
	})
}
