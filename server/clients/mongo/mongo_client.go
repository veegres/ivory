package mongo

import (
	"context"
	"crypto/tls"
	"errors"
	"ivory/clients"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var ErrHostOrPortNotSpecified = errors.New("host or port are not specified")

// Config contains the bare minimum details to open a mongo connection.
type Config struct {
	Host     string
	Username string
	Password string
	// AuthDb is the database credentials are verified against, defaulting to
	// "admin" (mongo's own default) when empty.
	AuthDb  string
	Port    int
	TLS     *tls.Config
	Timeout time.Duration
}

type Client struct {
	*mongo.Client
	Timeout time.Duration
}

// Connect opens a connection and returns it together with the connection url
// (without credentials) that can be shown to the user. It always dials the
// given host:port directly (Config.Host/Port name one specific replica set
// member, not the whole set) rather than letting the driver discover and
// route across the replica set topology, and eagerly pings so an unreachable
// host or wrong credentials fail here rather than on the first real command.
func Connect(ctx context.Context, c Config) (*Client, string, error) {
	if c.Port == 0 || c.Host == "" || c.Host == "-" {
		return nil, "unknown", ErrHostOrPortNotSpecified
	}

	timeout := c.Timeout
	if timeout == 0 {
		timeout = clients.IntegrationTimeout
	}

	address := c.Host + ":" + strconv.Itoa(c.Port)
	connUrl := "mongodb://" + address + "/?directConnection=true"

	opts := options.Client().
		SetHosts([]string{address}).
		SetDirect(true).
		SetConnectTimeout(timeout).
		SetServerSelectionTimeout(timeout)
	if c.Username != "" {
		authDb := c.AuthDb
		if authDb == "" {
			authDb = "admin"
		}
		opts.SetAuth(options.Credential{Username: c.Username, Password: c.Password, AuthSource: authDb})
	}
	if c.TLS != nil {
		opts.SetTLSConfig(c.TLS)
	}

	client, err := mongo.Connect(opts)
	if err != nil {
		return nil, connUrl, err
	}

	pingCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	if err := client.Ping(pingCtx, nil); err != nil {
		_ = client.Disconnect(context.Background())
		return nil, connUrl, err
	}
	return &Client{Client: client, Timeout: timeout}, connUrl, nil
}

func Close(client *Client) {
	_ = client.Disconnect(context.Background())
}
