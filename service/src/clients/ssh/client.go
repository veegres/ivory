package ssh

import (
	"bytes"
	"crypto/ed25519"
	"crypto/rand"
	"errors"
	"net"
	"net/url"
	"strconv"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"
)

var ErrCommandEmpty = errors.New("command cannot be empty")
var ErrHostEmpty = errors.New("vm host cannot be empty")

type Client struct {
	timeout time.Duration
}

func NewClient() *Client {
	return &Client{timeout: 10 * time.Second}
}

func (c *Client) Execute(connection Connection, command string) (*CommandResult, error) {
	trimmed := strings.TrimSpace(command)
	if trimmed == "" {
		return nil, ErrCommandEmpty
	}

	signer, err := ssh.NewSignerFromKey(connection.PrivateKey)
	if err != nil {
		return nil, err
	}

	config := &ssh.ClientConfig{
		User:            connection.Username,
		Auth:            []ssh.AuthMethod{ssh.PublicKeys(signer)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         c.timeout,
	}

	target, err := c.getDialAddress(connection)
	if err != nil {
		return nil, err
	}

	conn, err := ssh.Dial("tcp", target, config)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	session, err := conn.NewSession()
	if err != nil {
		return nil, err
	}
	defer session.Close()

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	session.Stdout = &stdout
	session.Stderr = &stderr

	result := &CommandResult{}
	err = session.Run(trimmed)
	result.Stdout = stdout.String()
	result.Stderr = stderr.String()

	var exitErr *ssh.ExitError
	if errors.As(err, &exitErr) {
		result.ExitCode = exitErr.ExitStatus()
		return result, nil
	}
	if err != nil {
		return nil, err
	}

	result.ExitCode = 0
	return result, nil
}

func (c *Client) ExecuteDocker(connection Connection, command string) (*CommandResult, error) {
	return c.Execute(connection, c.normalizeDockerCommand(command))
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

func (c *Client) normalizeDockerCommand(command string) string {
	trimmed := strings.TrimSpace(command)
	if trimmed == "" {
		return ""
	}
	if strings.HasPrefix(trimmed, "docker ") || trimmed == "docker" || strings.HasPrefix(trimmed, "sudo docker ") {
		return trimmed
	}
	return "docker " + trimmed
}

func (c *Client) getDialAddress(connection Connection) (string, error) {
	host := strings.TrimSpace(connection.Host)
	if host == "" {
		return "", ErrHostEmpty
	}

	if strings.Contains(host, "://") {
		parsed, err := url.Parse(host)
		if err != nil {
			return "", err
		}
		if parsed.Hostname() != "" {
			port := connection.Port
			if parsed.Port() != "" {
				parsedPort, err := strconv.Atoi(parsed.Port())
				if err != nil {
					return "", err
				}
				port = parsedPort
			}
			return net.JoinHostPort(parsed.Hostname(), strconv.Itoa(port)), nil
		}
	}

	if _, _, err := net.SplitHostPort(host); err == nil {
		return host, nil
	}

	return net.JoinHostPort(host, strconv.Itoa(connection.Port)), nil
}
