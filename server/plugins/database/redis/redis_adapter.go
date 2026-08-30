package redis

import (
	"context"
	"errors"
	"ivory/clients/redis"
	"ivory/plugins/database"
	"strconv"
)

var ErrNotSupported = errors.New("this operation is not supported for redis")

// NOTE: validate that is matches interface in compile-time
var _ database.Adapter = (*Adapter)(nil)

// Adapter executes redis commands as a generic console, the same way the
// etcd database plugin executes etcdctl-style commands. Unlike etcd, redis
// has hundreds of commands with wildly different reply shapes, so instead of
// hand-parsing a fixed verb set, GetFields tokenizes the input and passes it
// straight through to redis via the generic Do() command, formatting
// whatever comes back (see redis_adapter_query.go's formatReply). It has no
// databases/schemas/tables hierarchy beyond numeric DB indices (see
// ListDatabases) and no per-query sessions, so SchemaInquirer's
// schema/table methods and SessionManager's Cancel are not supported.
type Adapter struct{}

func NewAdapter() *Adapter {
	return &Adapter{}
}

func (a *Adapter) connect(ctx database.Context) (*redis.Client, string, error) {
	db := ctx.Connection.Config
	if db.Port == 0 || db.Host == "" || db.Host == "-" {
		return nil, "unknown", database.ErrDatabaseHostOrPortNotSpecified
	}

	var username, password string
	if ctx.Connection.Credentials != nil {
		username = ctx.Connection.Credentials.Username
		password = ctx.Connection.Credentials.Password
	}

	index := 0
	if db.Name != nil {
		if parsed, errParse := strconv.Atoi(*db.Name); errParse == nil {
			index = parsed
		}
	}

	return redis.Connect(context.Background(), redis.Config{
		Host:     db.Host,
		Port:     db.Port,
		Database: index,
		Username: username,
		Password: password,
		TLS:      ctx.Connection.TlsConfig,
	})
}

func field(name string, dataType string) database.QueryField {
	return database.QueryField{Name: name, DataType: dataType, DataTypeOID: 0}
}
