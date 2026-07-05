package etcd

import (
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"
)

var ErrEmptyCommand = errors.New("command is empty")
var ErrUnknownCommand = errors.New("unknown command, supported commands: get, put, del, member list, endpoint status, alarm list")
var ErrUnknownFlag = errors.New("unknown flag")
var ErrMissingArgument = errors.New("missing argument")
var ErrUnexpectedArgument = errors.New("unexpected argument")
var ErrInvalidLimit = errors.New("limit must be a positive number")
var ErrUnterminatedQuote = errors.New("unterminated quote")

type verb string

const (
	verbGet            verb = "get"
	verbPut            verb = "put"
	verbDel            verb = "del"
	verbMemberList     verb = "member list"
	verbEndpointStatus verb = "endpoint status"
	verbAlarmList      verb = "alarm list"
)

// command is a parsed etcdctl-style console query, e.g.
// `get /service --prefix --limit 100` or `member list`.
type command struct {
	Verb     verb
	Key      string
	Value    string
	Prefix   bool
	KeysOnly bool
	Limit    int64
}

func parseCommand(query string) (*command, error) {
	tokens, errTokens := tokenize(query)
	if errTokens != nil {
		return nil, errTokens
	}
	if len(tokens) == 0 {
		return nil, ErrEmptyCommand
	}

	switch tokens[0] {
	case "get":
		return parseKeyCommand(verbGet, tokens[1:], 1, []string{"--prefix", "--keys-only", "--limit"})
	case "put":
		return parseKeyCommand(verbPut, tokens[1:], 2, nil)
	case "del":
		return parseKeyCommand(verbDel, tokens[1:], 1, []string{"--prefix"})
	case "member":
		return parseBareCommand(verbMemberList, "list", tokens[1:])
	case "endpoint":
		return parseBareCommand(verbEndpointStatus, "status", tokens[1:])
	case "alarm":
		return parseBareCommand(verbAlarmList, "list", tokens[1:])
	default:
		return nil, fmt.Errorf("%w: %q", ErrUnknownCommand, tokens[0])
	}
}

func parseBareCommand(v verb, subcommand string, args []string) (*command, error) {
	if len(args) == 0 || args[0] != subcommand {
		return nil, fmt.Errorf("%w: expected %q", ErrMissingArgument, subcommand)
	}
	if len(args) > 1 {
		return nil, fmt.Errorf("%w: %q", ErrUnexpectedArgument, args[1])
	}
	return &command{Verb: v}, nil
}

func parseKeyCommand(v verb, args []string, positional int, flags []string) (*command, error) {
	cmd := &command{Verb: v}
	positionals := make([]string, 0, positional)

	for i := 0; i < len(args); i++ {
		arg := args[i]
		if !strings.HasPrefix(arg, "--") {
			positionals = append(positionals, arg)
			continue
		}
		if !slices.Contains(flags, arg) {
			return nil, fmt.Errorf("%w: %q", ErrUnknownFlag, arg)
		}
		switch arg {
		case "--prefix":
			cmd.Prefix = true
		case "--keys-only":
			cmd.KeysOnly = true
		case "--limit":
			if i+1 >= len(args) {
				return nil, fmt.Errorf("%w: --limit requires a value", ErrMissingArgument)
			}
			i++
			limit, errLimit := strconv.ParseInt(args[i], 10, 64)
			if errLimit != nil || limit <= 0 {
				return nil, fmt.Errorf("%w: %q", ErrInvalidLimit, args[i])
			}
			cmd.Limit = limit
		}
	}

	if len(positionals) < positional {
		return nil, fmt.Errorf("%w: %s expects %d argument(s)", ErrMissingArgument, v, positional)
	}
	if len(positionals) > positional {
		return nil, fmt.Errorf("%w: %q", ErrUnexpectedArgument, positionals[positional])
	}

	cmd.Key = positionals[0]
	if positional > 1 {
		cmd.Value = positionals[1]
	}
	return cmd, nil
}

// tokenize splits the query on whitespace with single/double-quote support,
// so keys and values may contain spaces and empty strings can be expressed
// as "" (used by `get "" --prefix` to list all keys).
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
