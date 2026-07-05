package ssh

import (
	"bufio"
	"crypto/ed25519"
	"errors"
	"fmt"
	"io"
	"log/slog"
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
var ErrRequestPtyFailed = errors.New("request for PTY failed")
var ErrStdoutPipeFailed = errors.New("stdout pipe failed")
var ErrStartCommandFailed = errors.New("start command failed")
var ErrDialFailed = errors.New("dial failed")
var ErrSessionFailed = errors.New("session failed")

type Command struct {
	Host       string
	Port       int
	Username   string
	Password   string
	PrivateKey *ed25519.PrivateKey

	Command         string
	HostKeyCallback ssh.HostKeyCallback
	Timeout         time.Duration

	JobID        string
	JobPersist   bool
	JobKeepAlive bool

	client  *ssh.Client
	session *ssh.Session
	Dial    func(network, addr string, config *ssh.ClientConfig) (*ssh.Client, error)
}

func (c *Command) Id() string {
	if c.JobID != "" {
		return c.JobID
	}
	return fmt.Sprintf("ssh|%s@%s|%s", c.Username, c.Host, c.Command)
}
func (c *Command) Persist() bool {
	return c.JobPersist
}

func (c *Command) KeepAlive() bool {
	return c.JobKeepAlive
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

	var client *ssh.Client
	if c.Dial != nil {
		client, err = c.Dial("tcp", target, config)
	} else {
		client, err = ssh.Dial("tcp", target, config)
	}
	if err != nil {
		return nil, errors.Join(ErrDialFailed, err)
	}

	session, err := client.NewSession()
	if err != nil {
		if c.Dial == nil {
			if err := client.Close(); err != nil {
				slog.Error("failed to close ssh client", "error", err)
			}
		}
		return nil, errors.Join(ErrSessionFailed, err)
	}

	// NOTE: SSH HACK - Explicitly request a PTY from the remote server
	// This forces Docker on the remote machine to treat Go like an interactive user terminal.
	modes := ssh.TerminalModes{
		ssh.ECHO:          0,     // Disable echo so you don't read your own command back
		ssh.TTY_OP_ISPEED: 14400, // Input speed
		ssh.TTY_OP_OSPEED: 14400, // Output speed
	}
	// NOTE: "xterm" simulates a standard linux terminal window (50000 columns, 40 rows)
	if err := session.RequestPty("xterm", 40, 50000, modes); err != nil {
		if err := session.Close(); err != nil {
			slog.Error("failed to close ssh session", "error", err)
		}
		if c.Dial == nil {
			if err := client.Close(); err != nil {
				slog.Error("failed to close ssh client", "error", err)
			}
		}
		return nil, errors.Join(ErrRequestPtyFailed, err)
	}

	stdout, err := session.StdoutPipe()
	if err != nil {
		if err := session.Close(); err != nil {
			slog.Error("failed to close ssh session", "error", err)
		}
		if c.Dial == nil {
			if err := client.Close(); err != nil {
				slog.Error("failed to close ssh client", "error", err)
			}
		}
		return nil, errors.Join(ErrStdoutPipeFailed, err)
	}
	if err := session.Start(c.Command); err != nil {
		if err := session.Close(); err != nil {
			slog.Error("failed to close ssh session", "error", err)
		}
		if c.Dial == nil {
			if err := client.Close(); err != nil {
				slog.Error("failed to close ssh client", "error", err)
			}
		}
		return nil, errors.Join(ErrStartCommandFailed, err)
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
	if c.client != nil && c.Dial == nil {
		errClose := c.client.Close()
		if errClose != nil {
			slog.Error("failed to close ssh client during wait", "error", errClose)
		}
		if errWait == nil {
			return errClose
		}
	}
	return errWait
}

func (c *Command) Abort() error {
	if c.session != nil {
		if err := c.session.Close(); err != nil {
			slog.Error("failed to close ssh session during abort", "error", err)
		}
	}
	if c.client != nil && c.Dial == nil {
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

	errScanner := scanner.Err()
	if errScanner != nil {
		// NOTE: Reading from a PTY might return an EIO error when the process exits.
		if !strings.Contains(errScanner.Error(), "input/output error") {
			return output, errScanner
		}
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
