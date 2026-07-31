package clickhouse

import (
	"context"
	"crypto/tls"
	"errors"
	"ivory/clients"
	"strconv"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
)

var ErrHostOrPortNotSpecified = errors.New("host or port are not specified")

// Config contains the bare minimum details to open a clickhouse connection.
type Config struct {
	Host           string
	Port           int
	Database       string // defaults to "default" when empty
	Username       string
	Password       string
	TLS            *tls.Config
	ConnectTimeout time.Duration
}

// Connect opens a connection over the native protocol and returns it
// together with the connection url (without credentials) that can be shown
// to the user. It eagerly pings the server so an unreachable host or wrong
// credentials fail here rather than on the first real query.
func Connect(ctx context.Context, c Config) (driver.Conn, string, error) {
	if c.Port == 0 || c.Host == "" || c.Host == "-" {
		return nil, "unknown", ErrHostOrPortNotSpecified
	}

	dbName := "default"
	if c.Database != "" {
		dbName = c.Database
	}

	timeout := c.ConnectTimeout
	if timeout == 0 {
		timeout = clients.IntegrationTimeout
	}

	connUrl := "clickhouse://" + c.Host + ":" + strconv.Itoa(c.Port) + "/" + dbName

	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{c.Host + ":" + strconv.Itoa(c.Port)},
		Auth: clickhouse.Auth{
			Database: dbName,
			Username: c.Username,
			Password: c.Password,
		},
		TLS:         c.TLS,
		DialTimeout: timeout,
	})
	if err != nil {
		return nil, connUrl, err
	}
	if err := conn.Ping(ctx); err != nil {
		return nil, connUrl, err
	}
	return conn, connUrl, nil
}

func Close(conn driver.Conn) {
	_ = conn.Close()
}
