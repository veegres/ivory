package redis

import (
	"context"
	"ivory/plugins/database"
	"strconv"
	"strings"
	"time"
)

// Cancel has no redis equivalent: commands run to completion synchronously
// on redis' single command thread, there is nothing to cancel mid-flight
// short of killing the connection outright, which is what Terminate does.
func (a *Adapter) Cancel(ctx database.Context, pid int) error {
	return ErrNotSupported
}

func (a *Adapter) Terminate(ctx database.Context, pid int) error {
	client, _, errConnect := a.connect(ctx)
	if errConnect != nil {
		return errConnect
	}
	defer client.Close()

	requestCtx, cancel := context.WithTimeout(context.Background(), client.Timeout)
	defer cancel()

	return client.ClientKillByFilter(requestCtx, "ID", strconv.Itoa(pid)).Err()
}

func (a *Adapter) ActiveQueries(ctx database.Context, options *database.QueryOptions) (*database.QueryFields, error) {
	startTime := time.Now().UnixMilli()

	client, url, errConnect := a.connect(ctx)
	if errConnect != nil {
		return nil, errConnect
	}
	defer client.Close()

	requestCtx, cancel := context.WithTimeout(context.Background(), client.Timeout)
	defer cancel()

	list, errList := client.ClientList(requestCtx).Result()
	if errList != nil {
		return nil, errList
	}

	fields, rows := parseClientList(list)
	return &database.QueryFields{
		Fields:    fields,
		Rows:      rows,
		Url:       url,
		StartTime: startTime,
		EndTime:   time.Now().UnixMilli(),
		Options:   options,
	}, nil
}

// parseClientList parses CLIENT LIST's one-line-per-client, space-separated
// key=value output. The first column is named "pid" (not "id") so the query
// table's terminate action, which looks up a field literally named "pid",
// can reuse it the same way it does for postgres' backend pid.
func parseClientList(text string) ([]database.QueryField, [][]any) {
	fields := []database.QueryField{
		field("pid", "int8"), field("addr", "text"), field("name", "text"),
		field("age", "int8"), field("idle", "int8"), field("db", "text"),
		field("cmd", "text"), field("user", "text"),
	}
	lines := strings.Split(strings.TrimRight(text, "\n"), "\n")
	rows := make([][]any, 0, len(lines))
	for _, line := range lines {
		if line == "" {
			continue
		}
		attrs := parseClientAttrs(line)
		pid, _ := strconv.ParseInt(attrs["id"], 10, 64)
		age, _ := strconv.ParseInt(attrs["age"], 10, 64)
		idle, _ := strconv.ParseInt(attrs["idle"], 10, 64)
		rows = append(rows, []any{pid, attrs["addr"], attrs["name"], age, idle, attrs["db"], attrs["cmd"], attrs["user"]})
	}
	return fields, rows
}

func parseClientAttrs(line string) map[string]string {
	attrs := map[string]string{}
	for _, token := range strings.Fields(line) {
		if idx := strings.IndexByte(token, '='); idx >= 0 {
			attrs[token[:idx]] = token[idx+1:]
		}
	}
	return attrs
}
