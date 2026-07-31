package zookeeper

import (
	"bufio"
	"context"
	"errors"
	"ivory/clients"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/go-zookeeper/zk"
)

var ErrHostOrPortNotSpecified = errors.New("host or port are not specified")
var ErrConnectTimeout = errors.New("timed out waiting for the zookeeper session to establish")

// sessionTimeout is the zookeeper session's own validity window (how long a
// session survives a lost connection before its ephemeral nodes/watches are
// dropped) - unrelated to how long Connect itself waits to establish that
// session (see ConnectTimeout). Servers commonly enforce a several-second
// minimum (the official image defaults minSessionTimeout to 4x tickTime,
// 4s), so this must stay well above clients.IntegrationTimeout.
const sessionTimeout = 10 * time.Second

// Config contains the bare minimum details to reach a zookeeper server.
type Config struct {
	Host           string
	Port           int
	Username       string
	Password       string
	Timeout        time.Duration
	ConnectTimeout time.Duration
}

// FourLetterCommand sends a single zookeeper "four-letter word" admin
// command (e.g. "mntr", "srvr", "conf") and returns its raw text response.
// Unlike a normal client connection, zookeeper closes the socket right after
// answering one of these commands, so a fresh connection is opened per call
// - there is no persistent client to keep alive or explicitly close beyond
// this one request. The server must have the command whitelisted via
// ZOO_4LW_COMMANDS_WHITELIST (recent zookeeper versions refuse unlisted
// commands by default).
func FourLetterCommand(ctx context.Context, c Config, command string) (string, error) {
	if c.Port == 0 || c.Host == "" || c.Host == "-" {
		return "", ErrHostOrPortNotSpecified
	}
	timeout := c.Timeout
	if timeout == 0 {
		timeout = clients.IntegrationTimeout
	}

	dialer := net.Dialer{Timeout: timeout}
	conn, err := dialer.DialContext(ctx, "tcp", c.Host+":"+strconv.Itoa(c.Port))
	if err != nil {
		return "", err
	}
	defer conn.Close()
	_ = conn.SetDeadline(time.Now().Add(timeout))

	if _, err := conn.Write([]byte(command + "\n")); err != nil {
		return "", err
	}

	var out strings.Builder
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		out.WriteString(scanner.Text())
		out.WriteByte('\n')
	}
	if err := scanner.Err(); err != nil {
		return "", err
	}
	return out.String(), nil
}

// Connect opens a real znode-protocol session (used for actual data
// operations - Get/Set/Create/Delete/Children - as opposed to
// FourLetterCommand's plain-text admin interface) and returns it together
// with the connection url (without credentials) that can be shown to the
// user. zk.Connect itself only starts a background connection attempt and
// returns immediately regardless of whether the server is reachable, so an
// Exists("/") round-trip - the root znode always exists - is used to
// eagerly prove connectivity within ConnectTimeout, the same way other
// clients eagerly ping, instead of only failing on the first real command.
func Connect(c Config) (*zk.Conn, string, error) {
	if c.Port == 0 || c.Host == "" || c.Host == "-" {
		return nil, "unknown", ErrHostOrPortNotSpecified
	}
	connectTimeout := c.ConnectTimeout
	if connectTimeout == 0 {
		connectTimeout = clients.IntegrationTimeout
	}
	addr := c.Host + ":" + strconv.Itoa(c.Port)
	connUrl := "zookeeper://" + addr

	conn, _, err := zk.Connect([]string{addr}, sessionTimeout)
	if err != nil {
		return nil, connUrl, err
	}

	result := make(chan error, 1)
	go func() {
		_, _, existsErr := conn.Exists("/")
		result <- existsErr
	}()
	select {
	case err := <-result:
		if err != nil {
			conn.Close()
			return nil, connUrl, err
		}
	case <-time.After(connectTimeout):
		conn.Close()
		return nil, connUrl, ErrConnectTimeout
	}

	if c.Username != "" {
		if err := conn.AddAuth("digest", []byte(c.Username+":"+c.Password)); err != nil {
			conn.Close()
			return nil, connUrl, err
		}
	}
	return conn, connUrl, nil
}

func Close(conn *zk.Conn) {
	conn.Close()
}
