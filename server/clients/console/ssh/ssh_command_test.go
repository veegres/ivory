package ssh

import (
	"crypto/ed25519"
	"crypto/rand"
	"errors"
	"io"
	"ivory/clients/console"
	"net"
	"strconv"
	"strings"
	"testing"
	"time"

	"golang.org/x/crypto/ssh"
)

func TestCommand_Id(t *testing.T) {
	cmd1 := &Command{Host: "localhost", Username: "user", Command: "ls -la"}
	cmd2 := &Command{Host: "localhost", Username: "user", Command: "ls -la"}
	cmd3 := &Command{Host: "127.0.0.1", Username: "user", Command: "ls -la"}

	if cmd1.Id() != cmd2.Id() {
		t.Errorf("Expected same ID for identical commands, got %s and %s", cmd1.Id(), cmd2.Id())
	}
	if cmd1.Id() == cmd3.Id() {
		t.Errorf("Expected different ID for different hosts")
	}

	customID := "custom-job-123"
	cmdCustom := &Command{Host: "localhost", Command: "ls", JobID: customID}
	if cmdCustom.Id() != customID {
		t.Errorf("Expected custom JobID %s, got %s", customID, cmdCustom.Id())
	}
}

func TestCommand_Errors(t *testing.T) {
	client := NewClient()
	_, prv, _ := ed25519.GenerateKey(rand.Reader)
	conn := Connection{Host: "localhost", Port: 22, Username: "test", PrivateKey: &prv}

	t.Run("empty command", func(t *testing.T) {
		_, err := client.Command(conn, "").Execute()
		if !errors.Is(err, ErrCommandEmpty) {
			t.Errorf("expected ErrCommandEmpty, got %v", err)
		}
	})

	t.Run("invalid private key", func(t *testing.T) {
		invalidConn := conn
		invalidKey := ed25519.PrivateKey(make([]byte, 64))
		invalidConn.PrivateKey = &invalidKey
		_, err := client.Command(invalidConn, "ls").Execute()
		if err == nil {
			t.Error("expected error for invalid private key, got nil")
		}
	})

	t.Run("dial error", func(t *testing.T) {
		dialErr := errors.New("dial failed")
		client.dial = func(network, addr string, config *ssh.ClientConfig) (*ssh.Client, error) {
			return nil, dialErr
		}
		_, err := client.Command(conn, "ls").Execute()
		if !errors.Is(err, dialErr) {
			t.Errorf("expected %v, got %v", dialErr, err)
		}
	})
}

func TestCommand_Execute(t *testing.T) {
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
							switch req.Type {
							case "pty-req":
								req.Reply(true, nil)
							case "exec":
								payload := req.Payload
								if len(payload) < 4 {
									req.Reply(false, nil)
									continue
								}
								cmdLen := int(payload[3]) | int(payload[2])<<8 | int(payload[1])<<16 | int(payload[0])<<24
								cmd := string(payload[4 : 4+cmdLen])

								if cmd == "echo hello" {
									req.Reply(true, nil)
									io.WriteString(channel, "hello\n")
									channel.SendRequest("exit-status", false, ssh.Marshal(struct{ Status uint32 }{0}))
									return
								} else if cmd == "fail" {
									req.Reply(true, nil)
									io.WriteString(channel, "error\n")
									channel.SendRequest("exit-status", false, ssh.Marshal(struct{ Status uint32 }{1}))
									return
								} else {
									req.Reply(false, nil)
								}
							default:
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
		PrivateKey: &clientPrv,
	}

	t.Run("command success", func(t *testing.T) {
		command := client.Command(conn, "echo hello")
		output, err := command.Execute()
		if err != nil {
			t.Fatalf("Execute failed: %v", err)
		}
		if len(output) != 1 || output[0] != "hello" {
			t.Errorf("expected stdout ['hello'], got %v", output)
		}
	})

	t.Run("command failure", func(t *testing.T) {
		command := client.Command(conn, "fail")
		output, err := command.Execute()
		if err == nil {
			t.Fatalf("expected error, got nil")
		}
		if !strings.Contains(err.Error(), "exit code 1") {
			t.Errorf("expected exit code 1 error, got %v", err)
		}
		if len(output) != 1 || output[0] != "error" {
			t.Errorf("expected output 'error' in array, got %v", output)
		}
	})
}

// TestCommand_ExecuteTimeout reproduces a hung remote command (e.g. a stuck
// "docker pull"): the fake server accepts the exec but never replies or
// exits, so without a bound Execute() would block forever. With
// ExecuteTimeout set, it must return console.ErrExecuteTimeout instead.
func TestCommand_ExecuteTimeout(t *testing.T) {
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
					_, requests, _ := newChannel.Accept()
					go func(in <-chan *ssh.Request) {
						for req := range in {
							switch req.Type {
							case "pty-req":
								req.Reply(true, nil)
							case "exec":
								req.Reply(true, nil)
								// NOTE: deliberately never writes output, exits,
								// or sends exit-status - simulates a hung remote
								// command.
							default:
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
		PrivateKey: &clientPrv,
	}

	command := client.Command(conn, "sleep 10")
	command.ExecuteTimeout = 50 * time.Millisecond

	start := time.Now()
	_, err = command.Execute()
	elapsed := time.Since(start)

	if !errors.Is(err, console.ErrExecuteTimeout) {
		t.Fatalf("expected console.ErrExecuteTimeout, got %v", err)
	}
	if elapsed > 2*time.Second {
		t.Fatalf("Execute() took %v, expected it to be bounded by ExecuteTimeout", elapsed)
	}
}

func TestCommand_Abort(t *testing.T) {
	// Start an SSH server that ignores exec for a bit so we can test abort
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
					channel, requests, _ := newChannel.Accept()
					go func(in <-chan *ssh.Request) {
						for req := range in {
							switch req.Type {
							case "pty-req":
								req.Reply(true, nil)
							case "exec":
								req.Reply(true, nil)
								go func() {
									time.Sleep(5 * time.Second) // Block to simulate a long running command
									if channel != nil {
										channel.Close()
									}
								}()
							default:
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
		PrivateKey: &clientPrv,
	}

	command := client.Command(conn, "sleep 10")
	_, errStart := command.Start()
	if errStart != nil {
		t.Fatalf("Start failed: %v", errStart)
	}

	errAbort := command.Abort()
	if errAbort != nil {
		t.Errorf("Abort() error = %v", errAbort)
	}

	errWait := command.Wait()
	if errWait == nil {
		t.Error("Expected error from Wait() after Abort(), got nil")
	}
}
