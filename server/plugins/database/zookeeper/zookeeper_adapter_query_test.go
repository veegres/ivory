package zookeeper

import (
	"errors"
	"ivory/plugins/database"
	"reflect"
	"testing"
)

func TestTokenize(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expected    []string
		expectedErr error
	}{
		{name: "simple command", input: "ls /service", expected: []string{"ls", "/service"}},
		{name: "quoted value with spaces", input: `create /motd "hello world"`, expected: []string{"create", "/motd", "hello world"}},
		{name: "extra whitespace tolerated", input: "  get   /a  \n", expected: []string{"get", "/a"}},
		{name: "empty input", input: "   ", expected: []string{}},
		{name: "unterminated quote", input: `get "/a`, expectedErr: ErrUnterminatedQuote},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokens, err := tokenize(tt.input)
			if tt.expectedErr != nil {
				if !errors.Is(err, tt.expectedErr) {
					t.Fatalf("expected error %v, got %v", tt.expectedErr, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
			if len(tokens) != len(tt.expected) {
				t.Fatalf("expected %v, got %v", tt.expected, tokens)
			}
			for i := range tokens {
				if tokens[i] != tt.expected[i] {
					t.Errorf("expected %v, got %v", tt.expected, tokens)
				}
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
		{name: "ls", query: "ls /service", expected: command{Verb: verbLs, Path: "/service"}},
		{name: "get", query: "get /service/config", expected: command{Verb: verbGet, Path: "/service/config"}},
		{name: "exists", query: "exists /service", expected: command{Verb: verbExists, Path: "/service"}},
		{name: "delete", query: "delete /service", expected: command{Verb: verbDelete, Path: "/service"}},
		{name: "create", query: `create /motd "hello world"`, expected: command{Verb: verbCreate, Path: "/motd", Data: "hello world"}},
		{name: "set", query: "set /service/config enabled", expected: command{Verb: verbSet, Path: "/service/config", Data: "enabled"}},
		{name: "empty query", query: "   ", expectedErr: ErrEmptyCommand},
		{name: "unknown command", query: "watch /a", expectedErr: ErrUnknownCommand},
		{name: "ls missing path", query: "ls", expectedErr: ErrMissingArgument},
		{name: "create missing data", query: "create /a", expectedErr: ErrMissingArgument},
		{name: "get extra argument", query: "get /a extra", expectedErr: ErrUnexpectedArgument},
		{name: "mntr", query: "mntr", expected: command{Verb: verbMntr}},
		{name: "srvr", query: "srvr", expected: command{Verb: verbSrvr}},
		{name: "stat", query: "stat", expected: command{Verb: verbStat}},
		{name: "ruok", query: "ruok", expected: command{Verb: verbRuok}},
		{name: "conf", query: "conf", expected: command{Verb: verbConf}},
		{name: "admin command takes no argument", query: "mntr /a", expectedErr: ErrUnexpectedArgument},
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

func TestParseAdminOutput(t *testing.T) {
	tests := []struct {
		name     string
		output   string
		expected [][]any
	}{
		{name: "mntr", output: "zk_version\t3.9.1\nzk_server_state\tleader\n", expected: [][]any{{"zk_version", "3.9.1"}, {"zk_server_state", "leader"}}},
		{name: "conf", output: "clientPort=2181\nserver.1=zk1:2888:3888:participant;0.0.0.0:2181\n", expected: [][]any{{"clientPort", "2181"}, {"server.1", "zk1:2888:3888:participant;0.0.0.0:2181"}}},
		{name: "srvr", output: "Zookeeper version: 3.9.1\nMode: follower\n\n", expected: [][]any{{"Zookeeper version", "3.9.1"}, {"Mode", "follower"}}},
		{name: "ruok", output: "imok", expected: [][]any{{"imok", ""}}},
		{name: "refusal is kept whole", output: "mntr is not executed because it is not in the whitelist.\n", expected: [][]any{{"mntr is not executed because it is not in the whitelist.", ""}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseAdminOutput(tt.output)
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("expected %v, got %v", tt.expected, got)
			}
		})
	}
}

func TestAdminRequiresAnEndpoint(t *testing.T) {
	_, _, _, err := NewAdapter().executeAdmin(database.Context{Connection: &database.Connection{}}, verbRuok)
	if !errors.Is(err, database.ErrDatabaseHostOrPortNotSpecified) {
		t.Errorf("expected host/port error, got %v", err)
	}
}

func TestIsAdmin(t *testing.T) {
	for _, v := range []verb{verbMntr, verbSrvr, verbStat, verbRuok, verbConf} {
		if !v.isAdmin() {
			t.Errorf("expected %s to be an admin command", v)
		}
	}
	for _, v := range []verb{verbLs, verbGet, verbCreate, verbSet, verbDelete, verbExists} {
		if v.isAdmin() {
			t.Errorf("expected %s to be a znode command", v)
		}
	}
}

func TestGetManyGetOneNotSupported(t *testing.T) {
	adapter := NewAdapter()
	ctx := database.Context{}

	if _, err := adapter.GetMany(ctx, "ls /", nil); !errors.Is(err, ErrNotSupported) {
		t.Errorf("expected ErrNotSupported, got %v", err)
	}
	if _, err := adapter.GetOne(ctx, "ls /"); !errors.Is(err, ErrNotSupported) {
		t.Errorf("expected ErrNotSupported, got %v", err)
	}
}
