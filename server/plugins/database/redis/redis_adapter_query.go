package redis

import (
	"context"
	"errors"
	"fmt"
	"ivory/plugins/database"
	"strings"
	"time"

	goredis "github.com/redis/go-redis/v9"
)

var ErrEmptyCommand = errors.New("command is empty")
var ErrUnterminatedQuote = errors.New("unterminated quote")

func (a *Adapter) GetFields(ctx database.Context, query string, options *database.QueryOptions) (*database.QueryFields, error) {
	startTime := time.Now().UnixMilli()

	tokens, errTokens := tokenize(query)
	if errTokens != nil {
		return nil, errTokens
	}
	if len(tokens) == 0 {
		return nil, ErrEmptyCommand
	}

	client, url, errConnect := a.connect(ctx)
	if errConnect != nil {
		return nil, errConnect
	}
	defer client.Close()

	requestCtx, cancel := context.WithTimeout(context.Background(), client.Timeout)
	defer cancel()

	args := make([]any, len(tokens))
	for i, token := range tokens {
		args[i] = token
	}
	result, errDo := client.Do(requestCtx, args...).Result()
	if errDo != nil && !errors.Is(errDo, goredis.Nil) {
		return nil, errDo
	}

	fields, rows := formatReply(result)
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

// formatValue reduces one redis reply value into a display data type and
// value. redis-cli itself always renders replies as text, and command reply
// shapes vary too widely (string/int/bool/array/nested array/nil) to model
// each one as its own typed column the way a fixed-schema SQL result can, so
// scalars keep their natural Go type (for right-aligned numeric display) and
// everything else is flattened to a readable string.
func formatValue(v any) (string, any) {
	switch val := v.(type) {
	case nil:
		return "text", "(nil)"
	case int64:
		return "int8", val
	case float64:
		return "float8", val
	case bool:
		return "bool", val
	case string:
		return "text", val
	case []byte:
		return "text", string(val)
	case []interface{}:
		parts := make([]string, len(val))
		for i, item := range val {
			_, value := formatValue(item)
			parts[i] = fmt.Sprintf("%v", value)
		}
		return "text", "[" + strings.Join(parts, ", ") + "]"
	default:
		return "text", fmt.Sprintf("%v", val)
	}
}

// formatReply turns one redis reply into a table: a top-level array reply
// (e.g. KEYS, LRANGE, CLIENT LIST as parsed by the caller) becomes one row
// per element, anything else becomes a single "result" row.
func formatReply(v any) ([]database.QueryField, [][]any) {
	if items, ok := v.([]interface{}); ok {
		fields := []database.QueryField{field("index", "int8"), field("value", "text")}
		rows := make([][]any, 0, len(items))
		for i, item := range items {
			_, value := formatValue(item)
			rows = append(rows, []any{int64(i), fmt.Sprintf("%v", value)})
		}
		return fields, rows
	}
	dataType, value := formatValue(v)
	return []database.QueryField{field("result", dataType)}, [][]any{{value}}
}

// tokenize splits a redis command line on whitespace with single/double
// quote support, so keys and values may contain spaces (e.g. `set greeting
// "hello world"`).
func tokenize(input string) ([]string, error) {
	tokens := make([]string, 0)
	var current strings.Builder
	started := false
	var quote rune

	for _, r := range input {
		switch {
		case quote != 0:
			if r == quote {
				quote = 0
			} else {
				current.WriteRune(r)
			}
		case r == '\'' || r == '"':
			quote = r
			started = true
		case r == ' ' || r == '\t' || r == '\n' || r == '\r':
			if started {
				tokens = append(tokens, current.String())
				current.Reset()
				started = false
			}
		default:
			current.WriteRune(r)
			started = true
		}
	}
	if quote != 0 {
		return nil, ErrUnterminatedQuote
	}
	if started {
		tokens = append(tokens, current.String())
	}
	return tokens, nil
}
