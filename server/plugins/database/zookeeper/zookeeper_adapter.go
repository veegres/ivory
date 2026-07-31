package zookeeper

import (
	"errors"
	zkclient "ivory/clients/zookeeper"
	"ivory/plugins/database"

	"github.com/go-zookeeper/zk"
)

var ErrNotSupported = errors.New("this operation is not supported for zookeeper")

// NOTE: validate that is matches interface in compile-time
var _ database.Adapter = (*Adapter)(nil)

// Adapter executes znode commands as a generic console, the same way the
// etcd database plugin executes etcdctl-style commands - a small, fixed
// verb set (ls/get/create/set/delete/exists) hand-parsed against zookeeper's
// hierarchical znode tree, the closest equivalent to etcd's flat key space.
// It has no databases/schemas/tables hierarchy beyond that znode tree and no
// per-query sessions, so SchemaInquirer's methods and SessionManager's
// Cancel/Terminate/ActiveQueries are not supported.
type Adapter struct{}

func NewAdapter() *Adapter {
	return &Adapter{}
}

func (a *Adapter) connect(ctx database.Context) (*zk.Conn, string, error) {
	db := ctx.Connection.Config
	if db.Port == 0 || db.Host == "" || db.Host == "-" {
		return nil, "unknown", database.ErrDatabaseHostOrPortNotSpecified
	}

	var username, password string
	if ctx.Connection.Credentials != nil {
		username = ctx.Connection.Credentials.Username
		password = ctx.Connection.Credentials.Password
	}

	return zkclient.Connect(zkclient.Config{
		Host:     db.Host,
		Port:     db.Port,
		Username: username,
		Password: password,
	})
}

func field(name string, dataType string) database.QueryField {
	return database.QueryField{Name: name, DataType: dataType, DataTypeOID: 0}
}
