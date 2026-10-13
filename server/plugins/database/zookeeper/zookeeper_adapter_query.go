package zookeeper

import (
	"context"
	"errors"
	"fmt"
	"ivory/clients/zookeeper"
	"ivory/plugins/database"
	"sort"
	"strings"
	"time"

	"github.com/go-zookeeper/zk"
)

var ErrEmptyCommand = errors.New("command is empty")
var ErrUnknownCommand = errors.New("unknown command, supported commands: ls, get, create, set, delete, exists, mntr, srvr, stat, ruok, conf")
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
	// NOTE: the admin ("four-letter word") commands are answered by the server
	// itself over a plain connection, with no znode session, and only when the
	// server's 4lw whitelist allows them.
	verbMntr verb = "mntr"
	verbSrvr verb = "srvr"
	verbStat verb = "stat"
	verbRuok verb = "ruok"
	verbConf verb = "conf"
)

func (v verb) isAdmin() bool {
	switch v {
	case verbMntr, verbSrvr, verbStat, verbRuok, verbConf:
		return true
	}
	return false
}

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

	var fields []database.QueryField
	var rows [][]any
	var url string
	if cmd.Verb.isAdmin() {
		var errAdmin error
		fields, rows, url, errAdmin = a.executeAdmin(ctx, cmd.Verb)
		if errAdmin != nil {
			return nil, errAdmin
		}
	} else {
		conn, connUrl, errConnect := a.connect(ctx)
		if errConnect != nil {
			return nil, errConnect
		}
		defer conn.Close()
		url = connUrl

		var errExecute error
		fields, rows, errExecute = execute(conn, cmd)
		if errExecute != nil {
			return nil, errExecute
		}
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
	case verbMntr, verbSrvr, verbStat, verbRuok, verbConf:
		return parseCommandArgs(verb(tokens[0]), tokens[1:], 0)
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
	cmd := &command{Verb: v}
	if arity > 0 {
		cmd.Path = args[0]
	}
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

func (a *Adapter) executeAdmin(ctx database.Context, v verb) ([]database.QueryField, [][]any, string, error) {
	db := ctx.Connection.Config
	if db.Port == 0 || db.Host == "" || db.Host == "-" {
		return nil, nil, "unknown", database.ErrDatabaseHostOrPortNotSpecified
	}
	url := fmt.Sprintf("%s:%d", db.Host, db.Port)
	output, err := zookeeper.FourLetterCommand(context.Background(), zookeeper.Config{Host: db.Host, Port: db.Port}, string(v))
	if err != nil {
		return nil, nil, url, err
	}
	return []database.QueryField{field("key", "text"), field("value", "text")}, parseAdminOutput(output), url, nil
}

// parseAdminOutput turns a four-letter reply into key/value rows. Each command
// separates its own way - mntr with a tab, conf with "=", srvr and stat with
// ": " - and a line with none of them (ruok's "imok", stat's client list, a
// whitelist refusal) is kept whole as a key, so nothing the server said is lost.
func parseAdminOutput(output string) [][]any {
	rows := make([][]any, 0)
	for _, line := range strings.Split(output, "\n") {
		line = strings.TrimRight(line, "\r")
		if strings.TrimSpace(line) == "" {
			continue
		}
		key, value := line, ""
		for _, sep := range []string{"\t", "=", ": "} {
			if k, v, found := strings.Cut(line, sep); found {
				key, value = k, v
				break
			}
		}
		rows = append(rows, []any{strings.TrimSpace(key), strings.TrimSpace(value)})
	}
	return rows
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
