package ssh

import (
	"bytes"
	"crypto/ed25519"
	"crypto/rand"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"

	"golang.org/x/crypto/ssh"
)

type Client struct {
	mu         sync.RWMutex
	timeout    time.Duration
	knownHosts map[string][]byte
	dial       func(network, addr string, config *ssh.ClientConfig) (*ssh.Client, error)
}

func NewClient() *Client {
	return &Client{
		timeout:    10 * time.Second,
		knownHosts: make(map[string][]byte),
		dial:       ssh.Dial,
	}
}

func (c *Client) GenerateKey() (string, string, error) {
	pubKey, prvKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return "", "", err
	}
	sshPubKey, err := ssh.NewPublicKey(pubKey)
	if err != nil {
		return "", "'", err
	}
	// NOTE: it always adds `\n` at the end, so we need to trim it
	sshPubKeyAuth := strings.TrimSuffix(string(ssh.MarshalAuthorizedKey(sshPubKey)), "\n")
	sshPubKeyAuthComment := sshPubKeyAuth + " " + "ivory"
	return sshPubKeyAuthComment, string(prvKey), nil
}

func (c *Client) Command(connection Command, command string) *Command {
	return &Command{
		Host:            connection.Host,
		Port:            connection.Port,
		Username:        connection.Username,
		Password:        connection.Password,
		PrivateKey:      connection.PrivateKey,
		Command:         command,
		HostKeyCallback: c.hostKeyCallback,
		Timeout:         c.timeout,
		dial:            c.dial,
	}
}

func (c *Client) getDialAddress(connection Command) (string, error) {
	return connection.getDialAddress()
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
