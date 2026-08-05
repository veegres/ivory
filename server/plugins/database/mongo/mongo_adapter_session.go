package mongo

import (
	"context"
	mongoclient "ivory/clients/mongo"
	"ivory/plugins/database"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

const adminDb = "admin"

// Cancel has no gentle mongo equivalent: killOp forcibly ends both the
// operation and its connection - the same thing Terminate does - so there is
// nothing left for a separate "cancel but keep the connection" call to do.
func (a *Adapter) Cancel(ctx database.Context, pid int) error {
	return ErrNotSupported
}

func (a *Adapter) Terminate(ctx database.Context, pid int) error {
	client, _, errConnect := a.connect(ctx)
	if errConnect != nil {
		return errConnect
	}
	defer mongoclient.Close(client)

	requestCtx, cancel := context.WithTimeout(context.Background(), client.Timeout)
	defer cancel()

	return client.Database(adminDb).RunCommand(requestCtx, bson.D{{Key: "killOp", Value: 1}, {Key: "op", Value: pid}}).Err()
}

func (a *Adapter) ActiveQueries(ctx database.Context, options *database.QueryOptions) (*database.QueryFields, error) {
	startTime := time.Now().UnixMilli()

	client, url, errConnect := a.connect(ctx)
	if errConnect != nil {
		return nil, errConnect
	}
	defer mongoclient.Close(client)

	requestCtx, cancel := context.WithTimeout(context.Background(), client.Timeout)
	defer cancel()

	var result struct {
		InProg []bson.M `bson:"inprog"`
	}
	errCmd := client.Database(adminDb).RunCommand(requestCtx, bson.D{{Key: "currentOp", Value: 1}}).Decode(&result)
	if errCmd != nil {
		return nil, errCmd
	}

	fields, rows := parseCurrentOp(result.InProg)
	return &database.QueryFields{
		Fields:    fields,
		Rows:      rows,
		Url:       url,
		StartTime: startTime,
		EndTime:   time.Now().UnixMilli(),
		Options:   options,
	}, nil
}

// parseCurrentOp maps currentOp's inprog documents onto a fixed table shape,
// the same pragmatic flattening redis' CLIENT LIST parsing does for its own
// free-form output. The first column is named "pid" (not "opid") so the
// query table's terminate action, which looks up a field literally named
// "pid", can reuse it the same way it does for postgres' backend pid.
func parseCurrentOp(inProg []bson.M) ([]database.QueryField, [][]any) {
	fields := []database.QueryField{
		field("pid", "int8"), field("ns", "text"), field("op", "text"),
		field("secsRunning", "int8"), field("client", "text"), field("desc", "text"),
	}
	rows := make([][]any, 0, len(inProg))
	for _, op := range inProg {
		rows = append(rows, []any{
			int64Field(op["opid"]), stringField(op["ns"]), stringField(op["op"]),
			int64Field(op["secs_running"]), stringField(op["client"]), stringField(op["desc"]),
		})
	}
	return fields, rows
}

func stringField(v any) string {
	s, _ := v.(string)
	return s
}

// int64Field narrows a decoded bson numeric value (int32/int64/float64,
// depending on how the server encoded it) to int64 for display; a value of
// any other shape (e.g. a sharded cluster's compound "shard:id" opid string,
// which never occurs in Ivory's single-replica-set deployments) reports 0.
func int64Field(v any) int64 {
	switch n := v.(type) {
	case int32:
		return int64(n)
	case int64:
		return n
	case float64:
		return int64(n)
	default:
		return 0
	}
}
