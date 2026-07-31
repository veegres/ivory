package zookeeper

import (
	"errors"
	"fmt"
	"ivory/plugins/database"
	"sort"
	"strings"
	"time"

	"github.com/go-zookeeper/zk"
)

var ErrEmptyCommand = errors.New("command is empty")
var ErrUnknownCommand = errors.New("unknown command, supported commands: ls, get, create, set, delete, exists")
var ErrMissingArgument = errors.New("missing argument")
var ErrUnexpectedArgument = errors.New("unexpected argument")
var ErrUnterminatedQuote = errors.New("unterminated quote")

type verb string

const (
	verbLs     verb = "ls"
	verbGet    verb = "get"
	verbCreate verb = "create"
	verbSet    verb = "set"
	verbDelete verb = "delete"
	verbExists verb = "exists"
)

// command is a parsed znode console query, e.g. `get /service/config` or
// `create /service/config "enabled"`.
type command struct {
	Verb verb
	Path string
	Data string
}

func (a *Adapter) GetFields(ctx database.Context, query string, options *database.QueryOptions) (*database.QueryFields, error) {
	startTime := time.Now().UnixMilli()

	cmd, errParse := parseCommand(query)
	if errParse != nil {
		return nil, errParse
	}

	conn, url, errConnect := a.connect(ctx)
	if errConnect != nil {
		return nil, errConnect
	}
	defer conn.Close()

	fields, rows, errExecute := execute(conn, cmd)
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

func parseCommand(query string) (*command, error) {
	tokens, errTokens := tokenize(query)
	if errTokens != nil {
		return nil, errTokens
	}
	if len(tokens) == 0 {
		return nil, ErrEmptyCommand
	}

	switch verb(tokens[0]) {
	case verbLs, verbDelete, verbExists:
		return parseCommandArgs(verb(tokens[0]), tokens[1:], 1)
	case verbGet:
		return parseCommandArgs(verbGet, tokens[1:], 1)
	case verbCreate, verbSet:
		return parseCommandArgs(verb(tokens[0]), tokens[1:], 2)
	default:
		return nil, fmt.Errorf("%w: %q", ErrUnknownCommand, tokens[0])
	}
}

func parseCommandArgs(v verb, args []string, arity int) (*command, error) {
	if len(args) < arity {
		return nil, fmt.Errorf("%w: %s expects %d argument(s)", ErrMissingArgument, v, arity)
	}
	if len(args) > arity {
		return nil, fmt.Errorf("%w: %q", ErrUnexpectedArgument, args[arity])
	}
	cmd := &command{Verb: v, Path: args[0]}
	if arity > 1 {
		cmd.Data = args[1]
	}
	return cmd, nil
}

func execute(conn *zk.Conn, cmd *command) ([]database.QueryField, [][]any, error) {
	switch cmd.Verb {
	case verbLs:
		children, _, err := conn.Children(cmd.Path)
		if err != nil {
			return nil, nil, err
		}
		sort.Strings(children)
		rows := make([][]any, 0, len(children))
		for _, child := range children {
			rows = append(rows, []any{child})
		}
		return []database.QueryField{field("child", "text")}, rows, nil
	case verbGet:
		data, stat, err := conn.Get(cmd.Path)
		if err != nil {
			return nil, nil, err
		}
		fields := []database.QueryField{
			field("data", "text"), field("version", "int8"),
			field("czxid", "int8"), field("mzxid", "int8"), field("numChildren", "int8"),
		}
		rows := [][]any{{string(data), int64(stat.Version), stat.Czxid, stat.Mzxid, int64(stat.NumChildren)}}
		return fields, rows, nil
	case verbCreate:
		path, err := conn.Create(cmd.Path, []byte(cmd.Data), zk.FlagPersistent, zk.WorldACL(zk.PermAll))
		if err != nil {
			return nil, nil, err
		}
		return []database.QueryField{field("created", "text")}, [][]any{{path}}, nil
	case verbSet:
		stat, err := conn.Set(cmd.Path, []byte(cmd.Data), -1)
		if err != nil {
			return nil, nil, err
		}
		return []database.QueryField{field("version", "int8")}, [][]any{{int64(stat.Version)}}, nil
	case verbDelete:
		if err := conn.Delete(cmd.Path, -1); err != nil {
			return nil, nil, err
		}
		return []database.QueryField{field("result", "text")}, [][]any{{"OK"}}, nil
	case verbExists:
		exists, stat, err := conn.Exists(cmd.Path)
		if err != nil {
			return nil, nil, err
		}
		var version int64
		if stat != nil {
			version = int64(stat.Version)
		}
		return []database.QueryField{field("exists", "bool"), field("version", "int8")}, [][]any{{exists, version}}, nil
	default:
		return nil, nil, fmt.Errorf("%w: %q", ErrUnknownCommand, cmd.Verb)
	}
}

// tokenize splits a znode command line on whitespace with single/double
// quote support, so a path or value may contain spaces (e.g. `create /motd
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
