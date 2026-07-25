package etcd

import (
	"crypto/tls"
	"ivory/clients"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

// Config contains the bare minimum details to open an etcd client connection.
type Config struct {
	Endpoints []string
	Username  string
	Password  string
	TLS       *tls.Config
	Timeout   time.Duration
}

type Client struct {
	*clientv3.Client
	Timeout time.Duration
}

func Connect(c Config) (*Client, error) {
	timeout := c.Timeout
	if timeout == 0 {
		timeout = clients.IntegrationTimeout
	}
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   c.Endpoints,
		Username:    c.Username,
		Password:    c.Password,
		TLS:         c.TLS,
		DialTimeout: timeout,
	})
	if err != nil {
		return nil, err
	}
	return &Client{Client: client, Timeout: timeout}, nil
}

func (c *Client) Close() error {
	return c.Client.Close()
}
