package redis

import (
	"context"
	"ivory/plugins/database"
	"strconv"
	"strings"
)

// ListDatabases reports redis' numbered databases (0..databases-1, redis has
// no named databases) filtered by name. ListSchemas/ListTables have nothing
// to report: redis has no schema or table hierarchy beyond that.
func (a *Adapter) ListDatabases(ctx database.Context, name string) ([]string, error) {
	client, _, errConnect := a.connect(ctx)
	if errConnect != nil {
		return nil, errConnect
	}
	defer client.Close()

	requestCtx, cancel := context.WithTimeout(context.Background(), client.Timeout)
	defer cancel()

	settings, errConfig := client.ConfigGet(requestCtx, "databases").Result()
	if errConfig != nil {
		return nil, errConfig
	}
	return filterDatabases(databaseCount(settings), name), nil
}

func (a *Adapter) ListSchemas(ctx database.Context, name string) ([]string, error) {
	return nil, ErrNotSupported
}

func (a *Adapter) ListTables(ctx database.Context, schema string, name string) ([]string, error) {
	return nil, ErrNotSupported
}

func databaseCount(settings map[string]string) int {
	if raw, ok := settings["databases"]; ok {
		if parsed, err := strconv.Atoi(raw); err == nil && parsed > 0 {
			return parsed
		}
	}
	return 16
}

func filterDatabases(count int, name string) []string {
	result := make([]string, 0, count)
	for i := 0; i < count; i++ {
		db := strconv.Itoa(i)
		if name == "" || strings.Contains(db, name) {
			result = append(result, db)
		}
	}
	return result
}
