package ssh

import (
	"bufio"
	"crypto/ed25519"
	"errors"
	"fmt"
	"io"
	"net"
	"net/url"
	"strconv"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"
)

var ErrCommandEmpty = errors.New("command cannot be empty")
var ErrHostEmpty = errors.New("vm host cannot be empty")
var ErrPasswordOrKey = errors.New("cannot use both private key and password")

type Command struct {
	Host       string
	Port       int
	Username   string
	Password   string
	PrivateKey *ed25519.PrivateKey

	Command         string
	HostKeyCallback ssh.HostKeyCallback
	Timeout         time.Duration

	JobID string
	Log   bool

	client  *ssh.Client
	session *ssh.Session
	dial    func(network, addr string, config *ssh.ClientConfig) (*ssh.Client, error)
}

func (c *Command) Id() string {
	if c.JobID != "" {
		return c.JobID
	}
	return fmt.Sprintf("ssh|%s@%s|%s",
		c.Username, c.Host, c.Command,
	)
}

func (c *Command) Persist() bool {
	return c.Log
}

func (c *Command) Start() (io.Reader, error) {
	trimmed := strings.TrimSpace(c.Command)
	if trimmed == "" {
		return nil, ErrCommandEmpty
	}

	if c.PrivateKey != nil && c.Password != "" {
		return nil, ErrPasswordOrKey
	}

	var auth []ssh.AuthMethod
	if c.PrivateKey != nil {
		signer, err := ssh.NewSignerFromKey(*c.PrivateKey)
		if err != nil {
			return nil, err
		}
		auth = append(auth, ssh.PublicKeys(signer))
	}
	if c.Password != "" {
		auth = append(auth, ssh.Password(c.Password))
	}

	config := &ssh.ClientConfig{
		User:            c.Username,
		Auth:            auth,
		HostKeyCallback: c.HostKeyCallback,
		Timeout:         c.Timeout,
	}

	target, err := c.getDialAddress()
	if err != nil {
		return nil, err
	}
	dial := c.dial
	if dial == nil {
		dial = ssh.Dial
	}
	client, err := dial("tcp", target, config)
	if err != nil {
		return nil, fmt.Errorf("dial: %w", err)
	}

	session, err := client.NewSession()
	if err != nil {
		client.Close()
		return nil, fmt.Errorf("session: %w", err)
	}

	stdout, err := session.StdoutPipe()
	if err != nil {
		session.Close()
		client.Close()
		return nil, fmt.Errorf("stdout pipe: %w", err)
	}
	if err := session.Start(c.Command); err != nil {
		session.Close()
		client.Close()
		return nil, fmt.Errorf("start command: %w", err)
	}

	c.client = client
	c.session = session
	return stdout, nil
}

func (c *Command) Wait() error {
	if c.session == nil {
		return nil
	}
	errWait := c.session.Wait()
	if c.client != nil {
		errClose := c.client.Close()
		if errWait == nil {
			return errClose
		}
	}
	return errWait
}

func (c *Command) Abort() error {
	if c.session != nil {
		_ = c.session.Close()
	}
	if c.client != nil {
		return c.client.Close()
	}
	return nil
}

func (c *Command) Execute() ([]string, error) {
	reader, err := c.Start()
	if err != nil {
		return nil, err
	}

	var output []string
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		output = append(output, scanner.Text())
	}

	errWait := c.Wait()
	if errWait != nil {
		var exitErr *ssh.ExitError
		if errors.As(errWait, &exitErr) {
			return output, fmt.Errorf("exit code %d: %s", exitErr.ExitStatus(), strings.Join(output, "\n"))
		}
		return output, errWait
	}

	return output, nil
}

func (c *Command) getDialAddress() (string, error) {
	host := strings.TrimSpace(c.Host)
	if host == "" {
		return "", ErrHostEmpty
	}

	if strings.Contains(host, "://") {
		parsed, err := url.Parse(host)
		if err != nil {
			return "", err
		}
		if parsed.Hostname() != "" {
			port := c.Port
			if parsed.Port() != "" {
				parsedPort, errAtoi := strconv.Atoi(parsed.Port())
				if errAtoi != nil {
					return "", errAtoi
				}
				port = parsedPort
			}
			return net.JoinHostPort(parsed.Hostname(), strconv.Itoa(port)), nil
		}
	}

	if _, _, err := net.SplitHostPort(host); err == nil {
		return host, nil
	}

	return net.JoinHostPort(host, strconv.Itoa(c.Port)), nil
}
