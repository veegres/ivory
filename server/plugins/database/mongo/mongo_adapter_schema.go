package mongo

import (
	"context"
	mongoclient "ivory/clients/mongo"
	"ivory/plugins/database"
	"strings"

	"go.mongodb.org/mongo-driver/v2/bson"
)

func (a *Adapter) ListDatabases(ctx database.Context, name string) ([]string, error) {
	client, _, errConnect := a.connect(ctx)
	if errConnect != nil {
		return nil, errConnect
	}
	defer mongoclient.Close(client)

	requestCtx, cancel := context.WithTimeout(context.Background(), client.Timeout)
	defer cancel()

	names, errList := client.ListDatabaseNames(requestCtx, bson.M{})
	if errList != nil {
		return nil, errList
	}
	return filterNames(names, name), nil
}

// ListSchemas has nothing to report: mongo has no schema layer between a
// database and its collections, only ListTables' database -> collection
// hierarchy (schema is left unused there for the same reason).
func (a *Adapter) ListSchemas(ctx database.Context, name string) ([]string, error) {
	return nil, ErrNotSupported
}

func (a *Adapter) ListTables(ctx database.Context, schema string, name string) ([]string, error) {
	client, _, errConnect := a.connect(ctx)
	if errConnect != nil {
		return nil, errConnect
	}
	defer mongoclient.Close(client)

	requestCtx, cancel := context.WithTimeout(context.Background(), client.Timeout)
	defer cancel()

	dbName := databaseName(ctx.Connection.Config.Name)
	names, errList := client.Database(dbName).ListCollectionNames(requestCtx, bson.M{})
	if errList != nil {
		return nil, errList
	}
	return filterNames(names, name), nil
}

func filterNames(names []string, filter string) []string {
	if filter == "" {
		return names
	}
	result := make([]string, 0, len(names))
	for _, n := range names {
		if strings.Contains(n, filter) {
			result = append(result, n)
		}
	}
	return result
}
