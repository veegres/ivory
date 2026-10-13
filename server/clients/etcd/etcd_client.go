package etcd

import (
	"context"
	"crypto/tls"
	"ivory/clients"
	"time"

	"go.etcd.io/etcd/api/v3/etcdserverpb"
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

// StatusSelf asks the member this client is connected to, over the existing
// connection, where Status(endpoint) dials a new one every call.
func (c *Client) StatusSelf(ctx context.Context) (*clientv3.StatusResponse, error) {
	resp, err := etcdserverpb.NewMaintenanceClient(c.ActiveConnection()).Status(ctx, &etcdserverpb.StatusRequest{})
	if err != nil {
		return nil, clientv3.ContextError(ctx, err)
	}
	return (*clientv3.StatusResponse)(resp), nil
}

func (c *Client) Close() error {
	return c.Client.Close()
}
