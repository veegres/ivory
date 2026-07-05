package postgres

import (
	"context"
	"crypto/tls"
	"errors"
	"log/slog"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5"
)

var ErrHostOrPortNotSpecified = errors.New("host or port are not specified")

const defaultConnectTimeout = 5 * time.Second

// Config contains the bare minimum details to open a postgres connection.
type Config struct {
	Host           string
	Port           int
	Database       string // defaults to "postgres" when empty
	Username       string
	Password       string
	AppName        string        // application_name runtime param
	TLS            *tls.Config   // nil = no ssl; non-nil = verify-ca with RootCAs/Certificates override
	ConnectTimeout time.Duration // connection establishment timeout, defaults to 5s
}

// Connect opens a connection and returns it together with the connection
// url (without credentials) that can be shown to the user.
func Connect(ctx context.Context, c Config) (*pgx.Conn, string, error) {
	conConfig, connUrl, err := Parse(c)
	if err != nil {
		return nil, connUrl, err
	}
	conn, err := pgx.ConnectConfig(ctx, conConfig)
	return conn, connUrl, err
}

// Parse builds the pgx connection config and the credential-free url from Config.
func Parse(c Config) (*pgx.ConnConfig, string, error) {
	if c.Port == 0 || c.Host == "" || c.Host == "-" {
		return nil, "unknown", ErrHostOrPortNotSpecified
	}

	dbName := "postgres"
	if c.Database != "" {
		dbName = c.Database
	}

	connProtocol := "postgres://"
	connHost := c.Host + ":" + strconv.Itoa(c.Port) + "/" + dbName
	connUrl := connProtocol + connHost

	if c.TLS != nil {
		// NOTE: verify-ca was chosen because it potentially can protect from machine-in-the-middle attack if
		// it has the right CA policy. More info can be found here https://www.postgresql.org/docs/16/libpq-ssl.html#LIBPQ-SSL-PROTECTION
		connUrl += "?sslmode=verify-ca"
	}

	conConfig, errConfig := pgx.ParseConfig(connUrl)
	if errConfig != nil {
		return nil, connUrl, errConfig
	}
	conConfig.User = c.Username
	conConfig.Password = c.Password
	conConfig.RuntimeParams = map[string]string{
		"application_name": c.AppName,
	}
	// NOTE: without the timeout a host that accepts tcp connections but never
	// answers the postgres handshake (a firewalled or non-postgres port) would
	// hang the caller forever
	conConfig.ConnectTimeout = c.ConnectTimeout
	if conConfig.ConnectTimeout == 0 {
		conConfig.ConnectTimeout = defaultConnectTimeout
	}
	if c.TLS != nil {
		// NOTE: we rewrite only RootCAs and Certificates, because pgx.ParseConfig creates proper
		//  tlsConfig for different `sslmode`. For example `verify-ca` should mark `InsecureSkipVerify=true`
		//  and it always sets `ServerName` it required for `verify-full` mode.
		conConfig.TLSConfig.RootCAs = c.TLS.RootCAs
		conConfig.TLSConfig.Certificates = c.TLS.Certificates
	}

	return conConfig, connUrl, nil
}

func Close(conn *pgx.Conn) {
	err := conn.Close(context.Background())
	if err != nil {
		slog.Warn("postgres close connection", "error", err)
	}
}
