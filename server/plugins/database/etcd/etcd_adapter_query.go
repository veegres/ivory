package etcd

import (
	"context"
	"errors"
	"fmt"
	"ivory/clients/etcd"
	"ivory/plugins/database"
	"slices"
	"strconv"
	"strings"
	"sync"
	"time"

	"go.etcd.io/etcd/api/v3/etcdserverpb"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
)

var ErrEmptyCommand = errors.New("command is empty")
var ErrUnknownCommand = errors.New("unknown command, supported commands: get, put, del, member list, endpoint status, endpoint health, alarm list")
var ErrUnknownFlag = errors.New("unknown flag")
var ErrMissingArgument = errors.New("missing argument")
var ErrUnexpectedArgument = errors.New("unexpected argument")
var ErrInvalidLimit = errors.New("limit must be a positive number")
var ErrUnterminatedQuote = errors.New("unterminated quote")

func (a *Adapter) GetFields(ctx database.Context, query string, options *database.QueryOptions) (*database.QueryFields, error) {
	startTime := time.Now().UnixMilli()

	cmd, errParse := parseCommand(query)
	if errParse != nil {
		return nil, errParse
	}
	normalizeLimit(cmd, options)

	client, url, errConnect := a.connect(ctx)
	if errConnect != nil {
		return nil, errConnect
	}
	defer client.Close()

	requestCtx, cancel := context.WithTimeout(context.Background(), client.Timeout)
	defer cancel()

	fields, rows, errExecute := a.execute(requestCtx, client, cmd)
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

type verb string

const (
	verbGet            verb = "get"
	verbPut            verb = "put"
	verbDel            verb = "del"
	verbMemberList     verb = "member list"
	verbEndpointStatus verb = "endpoint status"
	verbEndpointHealth verb = "endpoint health"
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
		if len(tokens) > 1 && tokens[1] == "health" {
			return parseBareCommand(verbEndpointHealth, "health", tokens[1:])
		}
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
	if len(positionals) == 0 {
		return nil, fmt.Errorf("%w: %s expects at least 1 argument(s)", ErrMissingArgument, v)
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

// normalizeLimit applies the console/template limit option when the command
// itself has no --limit, and clears the option when the command controls it,
// mirroring how the postgres adapter reports the effective limit back.
func normalizeLimit(cmd *command, options *database.QueryOptions) {
	if options == nil {
		return
	}
	if cmd.Verb != verbGet || cmd.Limit > 0 {
		options.Limit = nil
		return
	}
	if options.Limit == nil {
		return
	}
	limit, err := strconv.ParseInt(*options.Limit, 10, 64)
	if err != nil || limit <= 0 {
		options.Limit = nil
		return
	}
	cmd.Limit = limit
}

func (a *Adapter) execute(ctx context.Context, client *etcd.Client, cmd *command) ([]database.QueryField, [][]any, error) {
	switch cmd.Verb {
	case verbGet:
		opts := make([]clientv3.OpOption, 0)
		if cmd.Prefix {
			opts = append(opts, clientv3.WithPrefix())
		}
		if cmd.KeysOnly {
			opts = append(opts, clientv3.WithKeysOnly())
		}
		if cmd.Limit > 0 {
			opts = append(opts, clientv3.WithLimit(cmd.Limit))
		}
		response, err := client.Get(ctx, cmd.Key, opts...)
		if err != nil {
			return nil, nil, err
		}
		fields, rows := mapKvRows(response.Kvs, cmd.KeysOnly)
		return fields, rows, nil
	case verbPut:
		if _, err := client.Put(ctx, cmd.Key, cmd.Value); err != nil {
			return nil, nil, err
		}
		return []database.QueryField{field("result", "text")}, [][]any{{"OK"}}, nil
	case verbDel:
		opts := make([]clientv3.OpOption, 0)
		if cmd.Prefix {
			opts = append(opts, clientv3.WithPrefix())
		}
		response, err := client.Delete(ctx, cmd.Key, opts...)
		if err != nil {
			return nil, nil, err
		}
		return []database.QueryField{field("deleted", "int8")}, [][]any{{response.Deleted}}, nil
	case verbMemberList:
		response, err := client.MemberList(ctx)
		if err != nil {
			return nil, nil, err
		}
		fields, rows := mapMemberRows(response.Members)
		return fields, rows, nil
	case verbEndpointStatus:
		probes, err := a.probeEndpoints(ctx, client)
		if err != nil {
			return nil, nil, err
		}
		fields, rows := mapStatusRows(probes)
		return fields, rows, nil
	case verbEndpointHealth:
		probes, err := a.probeEndpoints(ctx, client)
		if err != nil {
			return nil, nil, err
		}
		fields, rows := mapHealthRows(probes)
		return fields, rows, nil
	case verbAlarmList:
		response, err := client.AlarmList(ctx)
		if err != nil {
			return nil, nil, err
		}
		fields := []database.QueryField{field("memberId", "text"), field("alarm", "text")}
		rows := make([][]any, 0, len(response.Alarms))
		for _, alarm := range response.Alarms {
			rows = append(rows, []any{formatID(alarm.MemberID), alarm.Alarm.String()})
		}
		return fields, rows, nil
	default:
		return nil, nil, fmt.Errorf("%w: %q", ErrUnknownCommand, cmd.Verb)
	}
}

// endpointProbe is one member's answer to a Status call, or why there is none.
type endpointProbe struct {
	Endpoint string
	ID       uint64
	Status   *clientv3.StatusResponse
	Err      error
	Took     time.Duration
}

// probeEndpoints asks every member for its status at once, each under its own
// deadline: sharing the command's one deadline let a dead member spend the
// budget the members after it needed, and they came back timed out too.
func (a *Adapter) probeEndpoints(ctx context.Context, client *etcd.Client) ([]endpointProbe, error) {
	memberList, err := client.MemberList(ctx)
	if err != nil {
		return nil, err
	}
	probes := make([]endpointProbe, len(memberList.Members))
	var wg sync.WaitGroup
	for i, m := range memberList.Members {
		probes[i] = endpointProbe{Endpoint: "-", ID: m.ID, Err: errMemberNotStarted}
		if len(m.ClientURLs) == 0 {
			continue
		}
		probes[i].Endpoint = m.ClientURLs[0]
		wg.Go(func() {
			probeCtx, cancel := context.WithTimeout(context.Background(), client.Timeout)
			defer cancel()
			start := time.Now()
			status, errStatus := client.Status(probeCtx, probes[i].Endpoint)
			probes[i].Status, probes[i].Err, probes[i].Took = status, errStatus, time.Since(start)
		})
	}
	wg.Wait()
	return probes, nil
}

var errMemberNotStarted = errors.New("member has not started")

// mapStatusRows reports each member's own status. errors is what the member
// says about itself - its alarms, a lost leader - and appliedIndex trailing
// raftIndex is a member that commits faster than it can apply.
func mapStatusRows(probes []endpointProbe) ([]database.QueryField, [][]any) {
	fields := []database.QueryField{
		field("endpoint", "text"), field("id", "text"), field("version", "text"),
		field("dbSize", "int8"), field("dbSizeInUse", "int8"), field("isLeader", "bool"), field("isLearner", "bool"),
		field("raftTerm", "int8"), field("raftIndex", "int8"), field("raftAppliedIndex", "int8"),
		field("errors", "text"), field("error", "text"),
	}
	rows := make([][]any, 0, len(probes))
	for _, p := range probes {
		if p.Err != nil {
			rows = append(rows, []any{p.Endpoint, formatID(p.ID), "-", int64(0), int64(0), false, false, int64(0), int64(0), int64(0), "", p.Err.Error()})
			continue
		}
		st := p.Status
		rows = append(rows, []any{
			p.Endpoint, formatID(p.ID), st.Version, st.DbSize, st.DbSizeInUse, st.Leader == p.ID, st.IsLearner,
			int64(st.RaftTerm), int64(st.RaftIndex), int64(st.RaftAppliedIndex), strings.Join(st.Errors, "; "), "",
		})
	}
	return fields, rows
}

// mapHealthRows answers what etcdctl's endpoint health does: a member is
// healthy when it answers, knows a leader and reports no alarm. took is the
// round trip, which is what a slow disk or network shows up as first.
func mapHealthRows(probes []endpointProbe) ([]database.QueryField, [][]any) {
	fields := []database.QueryField{
		field("endpoint", "text"), field("id", "text"), field("health", "bool"), field("took", "text"), field("error", "text"),
	}
	rows := make([][]any, 0, len(probes))
	for _, p := range probes {
		took := p.Took.Round(time.Millisecond).String()
		switch {
		case p.Err != nil:
			rows = append(rows, []any{p.Endpoint, formatID(p.ID), false, took, p.Err.Error()})
		case p.Status.Leader == 0:
			rows = append(rows, []any{p.Endpoint, formatID(p.ID), false, took, "member has no leader"})
		case len(p.Status.Errors) > 0:
			rows = append(rows, []any{p.Endpoint, formatID(p.ID), false, took, strings.Join(p.Status.Errors, "; ")})
		default:
			rows = append(rows, []any{p.Endpoint, formatID(p.ID), true, took, ""})
		}
	}
	return fields, rows
}

func mapKvRows(kvs []*mvccpb.KeyValue, keysOnly bool) ([]database.QueryField, [][]any) {
	if keysOnly {
		fields := []database.QueryField{field("key", "text")}
		rows := make([][]any, 0, len(kvs))
		for _, kv := range kvs {
			rows = append(rows, []any{string(kv.Key)})
		}
		return fields, rows
	}

	fields := []database.QueryField{
		field("key", "text"), field("value", "text"),
		field("createRevision", "int8"), field("modRevision", "int8"), field("version", "int8"),
	}
	rows := make([][]any, 0, len(kvs))
	for _, kv := range kvs {
		rows = append(rows, []any{string(kv.Key), string(kv.Value), kv.CreateRevision, kv.ModRevision, kv.Version})
	}
	return fields, rows
}

func mapMemberRows(members []*etcdserverpb.Member) ([]database.QueryField, [][]any) {
	fields := []database.QueryField{
		field("id", "text"), field("name", "text"), field("isLearner", "bool"),
		field("peerURLs", "text"), field("clientURLs", "text"),
	}
	rows := make([][]any, 0, len(members))
	for _, m := range members {
		rows = append(rows, []any{
			formatID(m.ID), m.Name, m.IsLearner,
			strings.Join(m.PeerURLs, ", "), strings.Join(m.ClientURLs, ", "),
		})
	}
	return fields, rows
}

func formatID(id uint64) string {
	return fmt.Sprintf("%x", id)
}

func field(name string, dataType string) database.QueryField {
	return database.QueryField{Name: name, DataType: dataType, DataTypeOID: 0}
}
