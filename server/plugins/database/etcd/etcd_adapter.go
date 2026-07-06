package etcd

import (
	"errors"
	"ivory/clients/etcd"
	"ivory/plugins/database"
	"strconv"
)

var ErrNotSupported = errors.New("this operation is not supported for etcd")

// NOTE: validate that is matches interface in compile-time
var _ database.Adapter = (*Adapter)(nil)

// Adapter executes etcdctl-style console commands against an etcd cluster.
// It has no databases/schemas/tables hierarchy and no per-query sessions,
// so SchemaInquirer and SessionManager are not supported and the matching
// features are excluded from SupportedFeatures.
type Adapter struct{}

func NewAdapter() *Adapter {
	return &Adapter{}
}

func (a *Adapter) connect(ctx database.Context) (*etcd.Client, string, error) {
	db := ctx.Connection.Config
	if db.Port == 0 || db.Host == "" || db.Host == "-" {
		return nil, "unknown", database.ErrDatabaseHostOrPortNotSpecified
	}

	var username, password string
	if ctx.Connection.Credentials != nil {
		username = ctx.Connection.Credentials.Username
		password = ctx.Connection.Credentials.Password
	}

	url := "etcd://" + db.Host + ":" + strconv.Itoa(db.Port)
	client, err := etcd.Connect(etcd.Config{
		Endpoints: []string{db.Host + ":" + strconv.Itoa(db.Port)},
		Username:  username,
		Password:  password,
		TLS:       ctx.Connection.TlsConfig,
	})
	if err != nil {
		return nil, url, err
	}
	return client, url, nil
}
