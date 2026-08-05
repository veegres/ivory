package mongo

import (
	"context"
	"errors"
	"fmt"
	mongoclient "ivory/clients/mongo"
	"ivory/plugins/database"
	"strconv"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var ErrEmptyCommand = errors.New("command is empty")
var ErrInvalidSyntax = errors.New(`invalid syntax, expected "<collection>.<verb>(<args>)" (e.g. users.find({}) or db.runCommand({dbStats: 1}))`)
var ErrUnterminatedQuote = errors.New("unterminated quote")
var ErrUnbalancedBrackets = errors.New("unbalanced brackets")
var ErrUnknownVerb = errors.New("unknown verb, supported verbs: find, findOne, insertOne, insertMany, updateOne, updateMany, deleteOne, deleteMany, countDocuments, distinct, aggregate, db.runCommand")
var ErrWrongArgumentCount = errors.New("wrong number of arguments")

func (a *Adapter) GetFields(ctx database.Context, query string, options *database.QueryOptions) (*database.QueryFields, error) {
	startTime := time.Now().UnixMilli()

	cmd, errParse := parseCommand(query)
	if errParse != nil {
		return nil, errParse
	}

	client, url, errConnect := a.connect(ctx)
	if errConnect != nil {
		return nil, errConnect
	}
	defer mongoclient.Close(client)

	requestCtx, cancel := context.WithTimeout(context.Background(), client.Timeout)
	defer cancel()

	dbName := databaseName(ctx.Connection.Config.Name)
	fields, rows, errExecute := a.execute(requestCtx, client, dbName, cmd, options)
	if errExecute != nil {
		return nil, errExecute
	}

	return &database.QueryFields{
		Fields:    fields,
		Rows:      rows,
		Url:       url,
		StartTime: startTime,
		EndTime:   time.Now().UnixMilli(),
		Options:   options,
	}, nil
}

func (a *Adapter) GetMany(ctx database.Context, query string, queryParams []any) ([]string, error) {
	return nil, ErrNotSupported
}

func (a *Adapter) GetOne(ctx database.Context, query string) (any, error) {
	return nil, ErrNotSupported
}

// command is a parsed "<collection>.<verb>(<args>)" console query, e.g.
// `users.find({"age": {"$gt": 21}})` or `db.runCommand({"dbStats": 1})`.
type command struct {
	Collection string
	Verb       string
	Args       []string
}

func parseCommand(query string) (*command, error) {
	trimmed := strings.TrimSpace(query)
	trimmed = strings.TrimSuffix(trimmed, ";")
	trimmed = strings.TrimSpace(trimmed)
	if trimmed == "" {
		return nil, ErrEmptyCommand
	}

	dot := strings.IndexByte(trimmed, '.')
	if dot <= 0 {
		return nil, ErrInvalidSyntax
	}
	collection := trimmed[:dot]
	rest := trimmed[dot+1:]

	open := strings.IndexByte(rest, '(')
	if open <= 0 || !strings.HasSuffix(rest, ")") {
		return nil, ErrInvalidSyntax
	}
	verb := rest[:open]

	args, errSplit := splitArgs(rest[open+1 : len(rest)-1])
	if errSplit != nil {
		return nil, errSplit
	}
	return &command{Collection: collection, Verb: verb, Args: args}, nil
}

// splitArgs splits a command's argument list on top-level commas, tracking
// object/array nesting depth and quotes so that commas inside a JSON
// argument (e.g. `{"a": 1, "b": 2}`) do not split it apart.
func splitArgs(text string) ([]string, error) {
	text = strings.TrimSpace(text)
	if text == "" {
		return nil, nil
	}

	args := make([]string, 0, 2)
	depth := 0
	var quote rune
	start := 0
	for i, r := range text {
		switch {
		case quote != 0:
			if r == quote {
				quote = 0
			}
		case r == '"' || r == '\'':
			quote = r
		case r == '{' || r == '[':
			depth++
		case r == '}' || r == ']':
			depth--
		case r == ',' && depth == 0:
			args = append(args, strings.TrimSpace(text[start:i]))
			start = i + 1
		}
	}
	if quote != 0 {
		return nil, ErrUnterminatedQuote
	}
	if depth != 0 {
		return nil, ErrUnbalancedBrackets
	}
	args = append(args, strings.TrimSpace(text[start:]))
	return args, nil
}

func (cmd *command) arg(i int) string {
	if i >= len(cmd.Args) {
		return ""
	}
	return cmd.Args[i]
}

// decodeArg parses one argument as MongoDB Extended JSON; an empty argument
// (an omitted optional parameter) leaves out untouched, keeping whatever
// zero value the caller pre-filled it with.
func decodeArg(text string, out any) error {
	if strings.TrimSpace(text) == "" {
		return nil
	}
	return bson.UnmarshalExtJSON([]byte(text), false, out)
}

func (a *Adapter) execute(ctx context.Context, client *mongoclient.Client, dbName string, cmd *command, opts *database.QueryOptions) ([]database.QueryField, [][]any, error) {
	if cmd.Collection == "db" {
		return a.executeDbVerb(ctx, client, dbName, cmd)
	}
	coll := client.Database(dbName).Collection(cmd.Collection)
	return a.executeCollectionVerb(ctx, coll, cmd, opts)
}

func (a *Adapter) executeDbVerb(ctx context.Context, client *mongoclient.Client, dbName string, cmd *command) ([]database.QueryField, [][]any, error) {
	if cmd.Verb != "runCommand" {
		return nil, nil, fmt.Errorf("%w: %q", ErrUnknownVerb, cmd.Verb)
	}
	if len(cmd.Args) != 1 {
		return nil, nil, fmt.Errorf("%w: runCommand expects 1 argument (command document)", ErrWrongArgumentCount)
	}
	var doc bson.M
	if errDecode := decodeArg(cmd.arg(0), &doc); errDecode != nil {
		return nil, nil, errDecode
	}

	var result bson.M
	if errCmd := client.Database(dbName).RunCommand(ctx, doc).Decode(&result); errCmd != nil {
		return nil, nil, errCmd
	}
	return documentRows([]bson.M{result})
}

func (a *Adapter) executeCollectionVerb(ctx context.Context, coll *mongo.Collection, cmd *command, opts *database.QueryOptions) ([]database.QueryField, [][]any, error) {
	switch cmd.Verb {
	case "find":
		return executeFind(ctx, coll, cmd, opts)
	case "findOne":
		return executeFindOne(ctx, coll, cmd)
	case "insertOne":
		return executeInsertOne(ctx, coll, cmd)
	case "insertMany":
		return executeInsertMany(ctx, coll, cmd)
	case "updateOne":
		return executeUpdate(ctx, coll, cmd, false)
	case "updateMany":
		return executeUpdate(ctx, coll, cmd, true)
	case "deleteOne":
		return executeDelete(ctx, coll, cmd, false)
	case "deleteMany":
		return executeDelete(ctx, coll, cmd, true)
	case "countDocuments":
		return executeCountDocuments(ctx, coll, cmd)
	case "distinct":
		return executeDistinct(ctx, coll, cmd)
	case "aggregate":
		return executeAggregate(ctx, coll, cmd)
	default:
		return nil, nil, fmt.Errorf("%w: %q", ErrUnknownVerb, cmd.Verb)
	}
}

func executeFind(ctx context.Context, coll *mongo.Collection, cmd *command, opts *database.QueryOptions) ([]database.QueryField, [][]any, error) {
	if len(cmd.Args) > 2 {
		return nil, nil, fmt.Errorf("%w: find expects at most 2 arguments (filter, projection)", ErrWrongArgumentCount)
	}
	filter := bson.M{}
	if errDecode := decodeArg(cmd.arg(0), &filter); errDecode != nil {
		return nil, nil, errDecode
	}
	findOpts := options.Find()
	if text := cmd.arg(1); text != "" {
		var projection bson.M
		if errDecode := decodeArg(text, &projection); errDecode != nil {
			return nil, nil, errDecode
		}
		findOpts.SetProjection(projection)
	}
	if limit := effectiveLimit(opts); limit > 0 {
		findOpts.SetLimit(limit)
	}

	cursor, errFind := coll.Find(ctx, filter, findOpts)
	if errFind != nil {
		return nil, nil, errFind
	}
	defer func() { _ = cursor.Close(ctx) }()

	var docs []bson.M
	if errAll := cursor.All(ctx, &docs); errAll != nil {
		return nil, nil, errAll
	}
	return documentRows(docs)
}

func executeFindOne(ctx context.Context, coll *mongo.Collection, cmd *command) ([]database.QueryField, [][]any, error) {
	if len(cmd.Args) > 2 {
		return nil, nil, fmt.Errorf("%w: findOne expects at most 2 arguments (filter, projection)", ErrWrongArgumentCount)
	}
	filter := bson.M{}
	if errDecode := decodeArg(cmd.arg(0), &filter); errDecode != nil {
		return nil, nil, errDecode
	}
	findOpts := options.FindOne()
	if text := cmd.arg(1); text != "" {
		var projection bson.M
		if errDecode := decodeArg(text, &projection); errDecode != nil {
			return nil, nil, errDecode
		}
		findOpts.SetProjection(projection)
	}

	var doc bson.M
	errFind := coll.FindOne(ctx, filter, findOpts).Decode(&doc)
	if errors.Is(errFind, mongo.ErrNoDocuments) {
		return []database.QueryField{field("index", "int8"), field("document", "text")}, [][]any{}, nil
	}
	if errFind != nil {
		return nil, nil, errFind
	}
	return documentRows([]bson.M{doc})
}

func executeInsertOne(ctx context.Context, coll *mongo.Collection, cmd *command) ([]database.QueryField, [][]any, error) {
	if len(cmd.Args) != 1 {
		return nil, nil, fmt.Errorf("%w: insertOne expects 1 argument (document)", ErrWrongArgumentCount)
	}
	var doc bson.M
	if errDecode := decodeArg(cmd.arg(0), &doc); errDecode != nil {
		return nil, nil, errDecode
	}
	result, errInsert := coll.InsertOne(ctx, doc)
	if errInsert != nil {
		return nil, nil, errInsert
	}
	return []database.QueryField{field("insertedId", "text")}, [][]any{{fmt.Sprintf("%v", result.InsertedID)}}, nil
}

func executeInsertMany(ctx context.Context, coll *mongo.Collection, cmd *command) ([]database.QueryField, [][]any, error) {
	if len(cmd.Args) != 1 {
		return nil, nil, fmt.Errorf("%w: insertMany expects 1 argument (documents array)", ErrWrongArgumentCount)
	}
	var docs []any
	if errDecode := decodeArg(cmd.arg(0), &docs); errDecode != nil {
		return nil, nil, errDecode
	}
	result, errInsert := coll.InsertMany(ctx, docs)
	if errInsert != nil {
		return nil, nil, errInsert
	}

	fields := []database.QueryField{field("index", "int8"), field("insertedId", "text")}
	rows := make([][]any, 0, len(result.InsertedIDs))
	for i, id := range result.InsertedIDs {
		rows = append(rows, []any{int64(i), fmt.Sprintf("%v", id)})
	}
	return fields, rows, nil
}

func executeUpdate(ctx context.Context, coll *mongo.Collection, cmd *command, many bool) ([]database.QueryField, [][]any, error) {
	if len(cmd.Args) != 2 {
		return nil, nil, fmt.Errorf("%w: expected 2 arguments (filter, update)", ErrWrongArgumentCount)
	}
	var filter, update bson.M
	if errDecode := decodeArg(cmd.arg(0), &filter); errDecode != nil {
		return nil, nil, errDecode
	}
	if errDecode := decodeArg(cmd.arg(1), &update); errDecode != nil {
		return nil, nil, errDecode
	}

	var result *mongo.UpdateResult
	var errUpdate error
	if many {
		result, errUpdate = coll.UpdateMany(ctx, filter, update)
	} else {
		result, errUpdate = coll.UpdateOne(ctx, filter, update)
	}
	if errUpdate != nil {
		return nil, nil, errUpdate
	}

	upsertedId := ""
	if result.UpsertedID != nil {
		upsertedId = fmt.Sprintf("%v", result.UpsertedID)
	}
	fields := []database.QueryField{field("matchedCount", "int8"), field("modifiedCount", "int8"), field("upsertedId", "text")}
	rows := [][]any{{result.MatchedCount, result.ModifiedCount, upsertedId}}
	return fields, rows, nil
}

func executeDelete(ctx context.Context, coll *mongo.Collection, cmd *command, many bool) ([]database.QueryField, [][]any, error) {
	if len(cmd.Args) > 1 {
		return nil, nil, fmt.Errorf("%w: expected at most 1 argument (filter)", ErrWrongArgumentCount)
	}
	filter := bson.M{}
	if errDecode := decodeArg(cmd.arg(0), &filter); errDecode != nil {
		return nil, nil, errDecode
	}

	var result *mongo.DeleteResult
	var errDelete error
	if many {
		result, errDelete = coll.DeleteMany(ctx, filter)
	} else {
		result, errDelete = coll.DeleteOne(ctx, filter)
	}
	if errDelete != nil {
		return nil, nil, errDelete
	}
	return []database.QueryField{field("deletedCount", "int8")}, [][]any{{result.DeletedCount}}, nil
}

func executeCountDocuments(ctx context.Context, coll *mongo.Collection, cmd *command) ([]database.QueryField, [][]any, error) {
	if len(cmd.Args) > 1 {
		return nil, nil, fmt.Errorf("%w: countDocuments expects at most 1 argument (filter)", ErrWrongArgumentCount)
	}
	filter := bson.M{}
	if errDecode := decodeArg(cmd.arg(0), &filter); errDecode != nil {
		return nil, nil, errDecode
	}
	count, errCount := coll.CountDocuments(ctx, filter)
	if errCount != nil {
		return nil, nil, errCount
	}
	return []database.QueryField{field("count", "int8")}, [][]any{{count}}, nil
}

func executeDistinct(ctx context.Context, coll *mongo.Collection, cmd *command) ([]database.QueryField, [][]any, error) {
	if len(cmd.Args) < 1 || len(cmd.Args) > 2 {
		return nil, nil, fmt.Errorf("%w: distinct expects 1 or 2 arguments (field, filter)", ErrWrongArgumentCount)
	}
	var fieldName string
	if errDecode := decodeArg(cmd.arg(0), &fieldName); errDecode != nil {
		return nil, nil, errDecode
	}
	filter := bson.M{}
	if errDecode := decodeArg(cmd.arg(1), &filter); errDecode != nil {
		return nil, nil, errDecode
	}

	var values []any
	if errDecode := coll.Distinct(ctx, fieldName, filter).Decode(&values); errDecode != nil {
		return nil, nil, errDecode
	}

	fields := []database.QueryField{field("index", "int8"), field("value", "text")}
	rows := make([][]any, 0, len(values))
	for i, v := range values {
		rows = append(rows, []any{int64(i), fmt.Sprintf("%v", v)})
	}
	return fields, rows, nil
}

func executeAggregate(ctx context.Context, coll *mongo.Collection, cmd *command) ([]database.QueryField, [][]any, error) {
	if len(cmd.Args) != 1 {
		return nil, nil, fmt.Errorf("%w: aggregate expects 1 argument (pipeline array)", ErrWrongArgumentCount)
	}
	var pipeline []bson.M
	if errDecode := decodeArg(cmd.arg(0), &pipeline); errDecode != nil {
		return nil, nil, errDecode
	}

	cursor, errAgg := coll.Aggregate(ctx, pipeline)
	if errAgg != nil {
		return nil, nil, errAgg
	}
	defer func() { _ = cursor.Close(ctx) }()

	var docs []bson.M
	if errAll := cursor.All(ctx, &docs); errAll != nil {
		return nil, nil, errAll
	}
	return documentRows(docs)
}

// documentRows flattens a set of documents into a table: each document is
// rendered as one relaxed Extended JSON string, the same pragmatic
// flattening redis' formatReply uses for its own heterogeneous replies,
// since a mongo result set has no fixed column schema to project onto.
func documentRows(docs []bson.M) ([]database.QueryField, [][]any, error) {
	fields := []database.QueryField{field("index", "int8"), field("document", "text")}
	rows := make([][]any, 0, len(docs))
	for i, doc := range docs {
		text, errMarshal := bson.MarshalExtJSON(doc, false, false)
		if errMarshal != nil {
			return nil, nil, errMarshal
		}
		rows = append(rows, []any{int64(i), string(text)})
	}
	return fields, rows, nil
}

// effectiveLimit reads the console/template limit option the same way
// etcd's normalizeLimit does, applying it only to find, the one verb here
// whose result size is otherwise unbounded by the query itself.
func effectiveLimit(opts *database.QueryOptions) int64 {
	if opts == nil || opts.Limit == nil {
		return 0
	}
	limit, err := strconv.ParseInt(*opts.Limit, 10, 64)
	if err != nil || limit <= 0 {
		return 0
	}
	return limit
}
