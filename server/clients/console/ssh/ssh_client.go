package ssh

import (
	"bytes"
	"crypto/ed25519"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"ivory/clients"
	"log/slog"
	"net"
	"strings"
	"sync"
	"time"

	"golang.org/x/crypto/ssh"
)

type Client struct {
	mu            sync.RWMutex
	timeout       time.Duration
	knownHosts    map[string][]byte
	cachedClients map[string]*ssh.Client
	dial          func(network, addr string, config *ssh.ClientConfig) (*ssh.Client, error)
}

func NewClient() *Client {
	return &Client{
		timeout:       clients.IntegrationTimeout,
		knownHosts:    make(map[string][]byte),
		cachedClients: make(map[string]*ssh.Client),
	}
}

func (c *Client) Command(con Connection, command string) *Command {
	return &Command{
		Host:            con.Host,
		Port:            con.Port,
		Username:        con.Username,
		Password:        con.Password,
		PrivateKey:      con.PrivateKey,
		JobPersist:      false,
		JobKeepAlive:    true,
		Command:         command,
		HostKeyCallback: c.hostKeyCallback,
		Timeout:         c.timeout,
		Dial: func(network, addr string, config *ssh.ClientConfig) (*ssh.Client, error) {
			return c.dialCached(network, addr, con.Password, con.PrivateKey, config)
		},
	}
}

// credentialHash is a cache-key fingerprint for cachedClients, not a password store,
// so a fast hash is fine here; it only needs to avoid keying the map on a raw secret.
func credentialHash(password string, privateKey *ed25519.PrivateKey) string {
	h := sha256.New()
	h.Write([]byte(password))
	if privateKey != nil {
		h.Write(*privateKey)
	}
	return hex.EncodeToString(h.Sum(nil)[:8])
}

func (c *Client) dialCached(network, addr, password string, privateKey *ed25519.PrivateKey, config *ssh.ClientConfig) (*ssh.Client, error) {
	key := fmt.Sprintf("%s|%s|%s|%s", addr, config.User, network, credentialHash(password, privateKey))

	c.mu.RLock()
	client, ok := c.cachedClients[key]
	c.mu.RUnlock()

	if ok {
		// Check if the connection is still alive by opening a temporary session
		// or sending a keepalive. Opening a session is the most reliable check.
		if session, err := client.NewSession(); err == nil {
			if err := session.Close(); err != nil {
				slog.Error("failed to close temporary ssh session", "error", err)
			}
			return client, nil
		}
		// Connection is dead, remove it
		c.mu.Lock()
		delete(c.cachedClients, key)
		c.mu.Unlock()
	}

	// Dial a new one
	dial := c.dial
	if dial == nil {
		dial = ssh.Dial
	}
	client, err := dial(network, addr, config)
	if err != nil {
		return nil, err
	}

	c.mu.Lock()
	c.cachedClients[key] = client
	c.mu.Unlock()

	return client, nil
}

func (c *Client) GenerateKey() (string, string, error) {
	pubKey, prvKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return "", "", err
	}
	sshPubKey, err := ssh.NewPublicKey(pubKey)
	if err != nil {
		return "", "", err
	}
	// NOTE: it always adds `\n` at the end, so we need to trim it
	sshPubKeyAuth := strings.TrimSuffix(string(ssh.MarshalAuthorizedKey(sshPubKey)), "\n")
	sshPubKeyAuthComment := sshPubKeyAuth + " " + "ivory"
	return sshPubKeyAuthComment, string(prvKey), nil
}

func (c *Client) hostKeyCallback(hostname string, _ net.Addr, key ssh.PublicKey) error {
	marshaledKey := key.Marshal()

	c.mu.RLock()
	knownKey, ok := c.knownHosts[hostname]
	c.mu.RUnlock()

	if !ok {
		c.mu.Lock()
		c.knownHosts[hostname] = marshaledKey
		c.mu.Unlock()
		return nil
	}

	if !bytes.Equal(knownKey, marshaledKey) {
		return fmt.Errorf("ssh: host key mismatch for %s", hostname)
	}

	return nil
}
