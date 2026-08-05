package mongo

import (
	"context"
	"errors"
	mongoclient "ivory/clients/mongo"
	"ivory/plugins/database"
)

var ErrNotSupported = errors.New("this operation is not supported for mongo")

// NOTE: validate that is matches interface in compile-time
var _ database.Adapter = (*Adapter)(nil)

// Adapter runs a small mongosh-like console language against mongo: queries
// are either "<collection>.<verb>(<args>)" (find, insertOne, ... - see
// mongo_adapter_query.go) or "db.runCommand(<command document>)", the one
// generic escape hatch for database/server-level admin commands
// (serverStatus, dbStats, currentOp, ...) that don't target a collection at
// all - the same passthrough role redis' Do()-based GetFields plays for
// commands with no fixed verb set. Arguments are parsed as MongoDB Extended
// JSON (bson.UnmarshalExtJSON), not full mongosh/JS syntax, so object keys
// must be quoted.
type Adapter struct{}

func NewAdapter() *Adapter {
	return &Adapter{}
}

func (a *Adapter) connect(ctx database.Context) (*mongoclient.Client, string, error) {
	db := ctx.Connection.Config
	if db.Port == 0 || db.Host == "" || db.Host == "-" {
		return nil, "unknown", database.ErrDatabaseHostOrPortNotSpecified
	}

	var username, password string
	if ctx.Connection.Credentials != nil {
		username = ctx.Connection.Credentials.Username
		password = ctx.Connection.Credentials.Password
	}

	return mongoclient.Connect(context.Background(), mongoclient.Config{
		Host:     db.Host,
		Port:     db.Port,
		Username: username,
		Password: password,
		TLS:      ctx.Connection.TlsConfig,
	})
}

// databaseName returns the connection's selected database, falling back to
// "test" - mongo's own default database when none is specified.
func databaseName(name *string) string {
	if name != nil && *name != "" {
		return *name
	}
	return "test"
}

func field(name string, dataType string) database.QueryField {
	return database.QueryField{Name: name, DataType: dataType, DataTypeOID: 0}
}
