package etcd

import (
	"context"
	"errors"
	"fmt"
	"ivory/clients/etcd"
	"ivory/core/config"
	"ivory/plugins/database"
	"strconv"
	"strings"
	"time"

	"go.etcd.io/etcd/api/v3/etcdserverpb"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
)

var ErrNotSupported = errors.New("this operation is not supported for etcd")

// NOTE: validate that is matches interface in compile-time
var _ database.Adapter = (*Adapter)(nil)

// Adapter executes etcdctl-style console commands against an etcd cluster.
// It has no databases/schemas/tables hierarchy and no per-query sessions,
// so SchemaInquirer and SessionManager are not supported and the matching
// features are excluded from SupportedFeatures.
type Adapter struct{}

func NewAdapter() *Adapter {
	return &Adapter{}
}

func (a *Adapter) SupportedFeatures() map[env.Feature]bool {
	return map[env.Feature]bool{
		env.ViewQueryDbInfo:        false,
		env.ViewQueryDbChart:       false,
		env.ManageQueryDbTemplate:  true,
		env.ManageQueryDbConsole:   true,
		env.ManageQueryDbCancel:    false,
		env.ManageQueryDbTerminate: false,
	}
}

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

func (a *Adapter) ListDatabases(ctx database.Context, name string) ([]string, error) {
	return nil, ErrNotSupported
}

func (a *Adapter) ListSchemas(ctx database.Context, name string) ([]string, error) {
	return nil, ErrNotSupported
}

func (a *Adapter) ListTables(ctx database.Context, schema string, name string) ([]string, error) {
	return nil, ErrNotSupported
}

func (a *Adapter) Cancel(ctx database.Context, pid int) error {
	return ErrNotSupported
}

func (a *Adapter) Terminate(ctx database.Context, pid int) error {
	return ErrNotSupported
}

func (a *Adapter) ActiveQueries(ctx database.Context, options *database.QueryOptions) (*database.QueryFields, error) {
	return nil, ErrNotSupported
}

func (a *Adapter) SystemCharts() map[database.SystemChartType]string {
	return map[database.SystemChartType]string{}
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
		return a.executeEndpointStatus(ctx, client)
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

func (a *Adapter) executeEndpointStatus(ctx context.Context, client *etcd.Client) ([]database.QueryField, [][]any, error) {
	memberList, err := client.MemberList(ctx)
	if err != nil {
		return nil, nil, err
	}

	fields := []database.QueryField{
		field("endpoint", "text"), field("id", "text"), field("version", "text"), field("dbSize", "int8"),
		field("isLeader", "bool"), field("raftTerm", "int8"), field("raftIndex", "int8"), field("error", "text"),
	}
	rows := make([][]any, 0, len(memberList.Members))
	for _, m := range memberList.Members {
		if len(m.ClientURLs) == 0 {
			rows = append(rows, []any{"-", formatID(m.ID), "-", int64(0), false, int64(0), int64(0), "member has not started"})
			continue
		}
		endpoint := m.ClientURLs[0]
		status, errStatus := client.Status(ctx, endpoint)
		if errStatus != nil {
			rows = append(rows, []any{endpoint, formatID(m.ID), "-", int64(0), false, int64(0), int64(0), errStatus.Error()})
			continue
		}
		isLeader := status.Leader == m.ID
		rows = append(rows, []any{endpoint, formatID(m.ID), status.Version, status.DbSize, isLeader, int64(status.RaftTerm), int64(status.RaftIndex), ""})
	}
	return fields, rows, nil
}

func (a *Adapter) connect(ctx database.Context) (*etcd.Client, string, error) {
	db := ctx.Connection.Config
	if db.Port == 0 || db.Host == "" || db.Host == "-" {
		return nil, "unknown", database.ErrDatabaseHostOrPortNotSpecified
	}

	var username, password string
	if ctx.Connection.Credentials != nil {
		username = ctx.Connection.Credentials.Username
		password = ctx.Connection.Credentials.Password
	}

	url := "etcd://" + db.Host + ":" + strconv.Itoa(db.Port)
	client, err := etcd.Connect(etcd.Config{
		Endpoints: []string{db.Host + ":" + strconv.Itoa(db.Port)},
		Username:  username,
		Password:  password,
		TLS:       ctx.Connection.TlsConfig,
	})
	if err != nil {
		return nil, url, err
	}
	return client, url, nil
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
