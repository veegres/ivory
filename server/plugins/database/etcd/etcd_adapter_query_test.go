package etcd

import (
	"errors"
	"ivory/plugins/database"
	"slices"
	"testing"

	"go.etcd.io/etcd/api/v3/etcdserverpb"
	"go.etcd.io/etcd/api/v3/mvccpb"
)

func TestMapKvRows(t *testing.T) {
	kvs := []*mvccpb.KeyValue{
		{Key: []byte("/a"), Value: []byte("1"), CreateRevision: 10, ModRevision: 12, Version: 2},
		{Key: []byte("/b"), Value: []byte("2"), CreateRevision: 11, ModRevision: 11, Version: 1},
	}

	t.Run("full rows", func(t *testing.T) {
		fields, rows := mapKvRows(kvs, false)
		if len(fields) != 5 || fields[0].Name != "key" || fields[1].Name != "value" {
			t.Fatalf("unexpected fields: %+v", fields)
		}
		if len(rows) != 2 {
			t.Fatalf("expected 2 rows, got %d", len(rows))
		}
		expected := []any{"/a", "1", int64(10), int64(12), int64(2)}
		if !slices.Equal(rows[0], expected) {
			t.Errorf("expected row %v, got %v", expected, rows[0])
		}
	})

	t.Run("keys only", func(t *testing.T) {
		fields, rows := mapKvRows(kvs, true)
		if len(fields) != 1 || fields[0].Name != "key" {
			t.Fatalf("unexpected fields: %+v", fields)
		}
		if rows[1][0] != "/b" {
			t.Errorf("expected key /b, got %v", rows[1][0])
		}
	})
}

func TestMapMemberRows(t *testing.T) {
	members := []*etcdserverpb.Member{
		{ID: 0xabc, Name: "etcd1", IsLearner: false, PeerURLs: []string{"http://e1:2380"}, ClientURLs: []string{"http://e1:2379", "http://10.0.0.1:2379"}},
		{ID: 0xdef, Name: "etcd2", IsLearner: true, PeerURLs: []string{"http://e2:2380"}, ClientURLs: nil},
	}

	fields, rows := mapMemberRows(members)
	if len(fields) != 5 {
		t.Fatalf("expected 5 fields, got %d", len(fields))
	}
	if rows[0][0] != "abc" {
		t.Errorf("expected hex id abc, got %v", rows[0][0])
	}
	if rows[0][4] != "http://e1:2379, http://10.0.0.1:2379" {
		t.Errorf("unexpected client urls: %v", rows[0][4])
	}
	if rows[1][2] != true {
		t.Errorf("expected learner true, got %v", rows[1][2])
	}
}

func TestNormalizeLimit(t *testing.T) {
	limit := "200"
	invalid := "abc"

	tests := []struct {
		name            string
		cmd             command
		options         *database.QueryOptions
		expectedCmd     int64
		expectedOptions *string
	}{
		{
			name:            "options limit applied to get without limit",
			cmd:             command{Verb: verbGet},
			options:         &database.QueryOptions{Limit: &limit},
			expectedCmd:     200,
			expectedOptions: &limit,
		},
		{
			name:            "command limit wins and clears option",
			cmd:             command{Verb: verbGet, Limit: 50},
			options:         &database.QueryOptions{Limit: &limit},
			expectedCmd:     50,
			expectedOptions: nil,
		},
		{
			name:            "non-get clears option",
			cmd:             command{Verb: verbMemberList},
			options:         &database.QueryOptions{Limit: &limit},
			expectedCmd:     0,
			expectedOptions: nil,
		},
		{
			name:            "invalid option limit is dropped",
			cmd:             command{Verb: verbGet},
			options:         &database.QueryOptions{Limit: &invalid},
			expectedCmd:     0,
			expectedOptions: nil,
		},
		{
			name:        "nil options tolerated",
			cmd:         command{Verb: verbGet},
			options:     nil,
			expectedCmd: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := tt.cmd
			normalizeLimit(&cmd, tt.options)
			if cmd.Limit != tt.expectedCmd {
				t.Errorf("expected command limit %d, got %d", tt.expectedCmd, cmd.Limit)
			}
			if tt.options == nil {
				return
			}
			if tt.expectedOptions == nil && tt.options.Limit != nil {
				t.Errorf("expected options limit to be cleared, got %v", *tt.options.Limit)
			}
			if tt.expectedOptions != nil && (tt.options.Limit == nil || *tt.options.Limit != *tt.expectedOptions) {
				t.Errorf("expected options limit %v, got %v", *tt.expectedOptions, tt.options.Limit)
			}
		})
	}
}

func TestParseCommand(t *testing.T) {
	tests := []struct {
		name        string
		query       string
		expected    command
		expectedErr error
	}{
		{
			name:     "get key",
			query:    `get /service/config`,
			expected: command{Verb: verbGet, Key: "/service/config"},
		},
		{
			name:     "get with all flags",
			query:    `get /service --prefix --keys-only --limit 100`,
			expected: command{Verb: verbGet, Key: "/service", Prefix: true, KeysOnly: true, Limit: 100},
		},
		{
			name:     "get empty quoted key lists everything",
			query:    `get "" --prefix --keys-only --limit 500`,
			expected: command{Verb: verbGet, Key: "", Prefix: true, KeysOnly: true, Limit: 500},
		},
		{
			name:     "put key value",
			query:    `put /service/config enabled`,
			expected: command{Verb: verbPut, Key: "/service/config", Value: "enabled"},
		},
		{
			name:     "put quoted value with spaces",
			query:    `put /motd "hello world"`,
			expected: command{Verb: verbPut, Key: "/motd", Value: "hello world"},
		},
		{
			name:     "put single quoted json",
			query:    `put /config '{"a": 1}'`,
			expected: command{Verb: verbPut, Key: "/config", Value: `{"a": 1}`},
		},
		{
			name:     "del with prefix",
			query:    `del /tmp --prefix`,
			expected: command{Verb: verbDel, Key: "/tmp", Prefix: true},
		},
		{
			name:     "member list",
			query:    `member list`,
			expected: command{Verb: verbMemberList},
		},
		{
			name:     "endpoint status",
			query:    `endpoint status`,
			expected: command{Verb: verbEndpointStatus},
		},
		{
			name:     "alarm list",
			query:    `alarm list`,
			expected: command{Verb: verbAlarmList},
		},
		{
			name:     "whitespace tolerated",
			query:    "  get   /key  \n",
			expected: command{Verb: verbGet, Key: "/key"},
		},
		{name: "empty query", query: "   ", expectedErr: ErrEmptyCommand},
		{name: "unknown command", query: "watch /key", expectedErr: ErrUnknownCommand},
		{name: "unknown flag", query: "get /key --recursive", expectedErr: ErrUnknownFlag},
		{name: "get missing key", query: "get", expectedErr: ErrMissingArgument},
		{name: "put missing value", query: "put /key", expectedErr: ErrMissingArgument},
		{name: "get extra argument", query: "get /key extra", expectedErr: ErrUnexpectedArgument},
		{name: "member missing subcommand", query: "member", expectedErr: ErrMissingArgument},
		{name: "member extra argument", query: "member list extra", expectedErr: ErrUnexpectedArgument},
		{name: "limit missing value", query: "get /key --limit", expectedErr: ErrMissingArgument},
		{name: "limit not a number", query: "get /key --limit abc", expectedErr: ErrInvalidLimit},
		{name: "limit negative", query: "get /key --limit -5", expectedErr: ErrInvalidLimit},
		{name: "unterminated quote", query: `get "/key`, expectedErr: ErrUnterminatedQuote},
		{name: "del flag not allowed keys-only", query: "del /key --keys-only", expectedErr: ErrUnknownFlag},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parsed, err := parseCommand(tt.query)
			if tt.expectedErr != nil {
				if !errors.Is(err, tt.expectedErr) {
					t.Fatalf("expected error %v, got %v", tt.expectedErr, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
			if *parsed != tt.expected {
				t.Errorf("expected %+v, got %+v", tt.expected, *parsed)
			}
		})
	}
}

func TestGetManyGetOneNotSupported(t *testing.T) {
	adapter := NewAdapter()
	ctx := database.Context{}

	if _, err := adapter.GetMany(ctx, "get /a", nil); !errors.Is(err, ErrNotSupported) {
		t.Errorf("expected ErrNotSupported, got %v", err)
	}
	if _, err := adapter.GetOne(ctx, "get /a"); !errors.Is(err, ErrNotSupported) {
		t.Errorf("expected ErrNotSupported, got %v", err)
	}
}
