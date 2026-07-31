package redis

import (
	"context"
	"crypto/tls"
	"errors"
	"ivory/clients"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

var ErrHostOrPortNotSpecified = errors.New("host or port are not specified")

// Config contains the bare minimum details to open a redis connection.
type Config struct {
	Host           string
	Port           int
	Database       int // redis numeric database index, defaults to 0
	Username       string
	Password       string
	TLS            *tls.Config
	ConnectTimeout time.Duration
}

type Client struct {
	*redis.Client
	Timeout time.Duration
}

// Connect opens a connection and returns it together with the connection
// url (without credentials) that can be shown to the user. It eagerly pings
// the server so an unreachable host or wrong credentials fail here rather
// than on the first real command.
func Connect(ctx context.Context, c Config) (*Client, string, error) {
	if c.Port == 0 || c.Host == "" || c.Host == "-" {
		return nil, "unknown", ErrHostOrPortNotSpecified
	}

	timeout := c.ConnectTimeout
	if timeout == 0 {
		timeout = clients.IntegrationTimeout
	}

	connUrl := "redis://" + c.Host + ":" + strconv.Itoa(c.Port) + "/" + strconv.Itoa(c.Database)

	client := redis.NewClient(&redis.Options{
		Addr:        c.Host + ":" + strconv.Itoa(c.Port),
		Username:    c.Username,
		Password:    c.Password,
		DB:          c.Database,
		TLSConfig:   c.TLS,
		DialTimeout: timeout,
	})
	if err := client.Ping(ctx).Err(); err != nil {
		_ = client.Close()
		return nil, connUrl, err
	}
	return &Client{Client: client, Timeout: timeout}, connUrl, nil
}

func Close(client *Client) {
	_ = client.Close()
}
