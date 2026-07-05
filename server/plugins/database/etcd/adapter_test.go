package etcd

import (
	"errors"
	"ivory/core/config"
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

func TestNotSupportedOperations(t *testing.T) {
	adapter := NewAdapter()
	ctx := database.Context{}

	type op struct {
		name string
		call func() error
	}
	ops := []op{
		{"GetMany", func() error { _, e := adapter.GetMany(ctx, "get /a", nil); return e }},
		{"GetOne", func() error { _, e := adapter.GetOne(ctx, "get /a"); return e }},
		{"ListDatabases", func() error { _, e := adapter.ListDatabases(ctx, ""); return e }},
		{"ListSchemas", func() error { _, e := adapter.ListSchemas(ctx, ""); return e }},
		{"ListTables", func() error { _, e := adapter.ListTables(ctx, "", ""); return e }},
		{"Cancel", func() error { return adapter.Cancel(ctx, 1) }},
		{"Terminate", func() error { return adapter.Terminate(ctx, 1) }},
		{"ActiveQueries", func() error { _, e := adapter.ActiveQueries(ctx, nil); return e }},
	}

	for _, o := range ops {
		t.Run(o.name, func(t *testing.T) {
			if err := o.call(); !errors.Is(err, ErrNotSupported) {
				t.Errorf("expected ErrNotSupported, got %v", err)
			}
		})
	}
}

func TestSupportedFeaturesExclusions(t *testing.T) {
	features := NewAdapter().SupportedFeatures()

	expected := []env.Feature{env.ManageQueryDbConsole, env.ManageQueryDbTemplate}
	if !slices.Equal(features, expected) {
		t.Fatalf("expected features %v, got %v", expected, features)
	}

	excluded := []env.Feature{env.ViewQueryDbInfo, env.ViewQueryDbChart, env.ManageQueryDbCancel, env.ManageQueryDbTerminate}
	for _, feature := range excluded {
		if slices.Contains(features, feature) {
			t.Errorf("feature %v must not be supported for etcd", feature)
		}
	}
}

func TestSystemChartsEmpty(t *testing.T) {
	if len(NewAdapter().SystemCharts()) != 0 {
		t.Fatal("expected no system charts for etcd")
	}
}
