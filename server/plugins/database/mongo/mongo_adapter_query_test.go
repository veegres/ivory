package mongo

import (
	"errors"
	"ivory/plugins/database"
	"testing"

	"go.mongodb.org/mongo-driver/v2/bson"
)

func TestParseCommand(t *testing.T) {
	tests := []struct {
		name        string
		query       string
		expected    command
		expectedErr error
	}{
		{
			name:     "find with filter",
			query:    `users.find({"age": {"$gt": 21}})`,
			expected: command{Collection: "users", Verb: "find", Args: []string{`{"age": {"$gt": 21}}`}},
		},
		{
			name:     "find with filter and projection",
			query:    `users.find({"age": 21}, {"name": 1})`,
			expected: command{Collection: "users", Verb: "find", Args: []string{`{"age": 21}`, `{"name": 1}`}},
		},
		{
			name:     "find with no arguments",
			query:    `users.find()`,
			expected: command{Collection: "users", Verb: "find", Args: nil},
		},
		{
			name:     "db.runCommand",
			query:    `db.runCommand({"dbStats": 1})`,
			expected: command{Collection: "db", Verb: "runCommand", Args: []string{`{"dbStats": 1}`}},
		},
		{
			name:     "trailing semicolon and whitespace tolerated",
			query:    "  users.find({});  \n",
			expected: command{Collection: "users", Verb: "find", Args: []string{"{}"}},
		},
		{name: "empty query", query: "   ", expectedErr: ErrEmptyCommand},
		{name: "missing dot", query: "find({})", expectedErr: ErrInvalidSyntax},
		{name: "missing parens", query: "users.find", expectedErr: ErrInvalidSyntax},
		{name: "missing closing paren", query: "users.find({}", expectedErr: ErrInvalidSyntax},
		{name: "unterminated quote", query: `users.find({"a": "b})`, expectedErr: ErrUnterminatedQuote},
		{name: "unbalanced brackets", query: `users.find({"a": 1)`, expectedErr: ErrUnbalancedBrackets},
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
			if parsed.Collection != tt.expected.Collection || parsed.Verb != tt.expected.Verb {
				t.Fatalf("expected collection/verb %q/%q, got %q/%q", tt.expected.Collection, tt.expected.Verb, parsed.Collection, parsed.Verb)
			}
			if len(parsed.Args) != len(tt.expected.Args) {
				t.Fatalf("expected args %v, got %v", tt.expected.Args, parsed.Args)
			}
			for i := range parsed.Args {
				if parsed.Args[i] != tt.expected.Args[i] {
					t.Errorf("expected arg[%d] %q, got %q", i, tt.expected.Args[i], parsed.Args[i])
				}
			}
		})
	}
}

func TestSplitArgs(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		expected []string
	}{
		{name: "empty", text: "", expected: nil},
		{name: "single", text: `{"a": 1}`, expected: []string{`{"a": 1}`}},
		{name: "two args with nested commas", text: `{"a": 1, "b": 2}, {"c": 1}`, expected: []string{`{"a": 1, "b": 2}`, `{"c": 1}`}},
		{name: "array argument", text: `[{"$match": {"a": 1, "b": 2}}]`, expected: []string{`[{"$match": {"a": 1, "b": 2}}]`}},
		{name: "quoted comma ignored", text: `"a,b", 1`, expected: []string{`"a,b"`, "1"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args, err := splitArgs(tt.text)
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
			if len(args) != len(tt.expected) {
				t.Fatalf("expected %v, got %v", tt.expected, args)
			}
			for i := range args {
				if args[i] != tt.expected[i] {
					t.Errorf("expected arg[%d] %q, got %q", i, tt.expected[i], args[i])
				}
			}
		})
	}
}

func TestDecodeArg(t *testing.T) {
	var m bson.M
	if err := decodeArg(`{"a": 1}`, &m); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if m["a"] != int32(1) {
		t.Errorf("expected a=1, got %v", m["a"])
	}

	var empty bson.M
	if err := decodeArg("  ", &empty); err != nil {
		t.Fatalf("expected no error for blank argument, got %v", err)
	}
	if empty != nil {
		t.Errorf("expected out to be left untouched for a blank argument, got %v", empty)
	}
}

func TestDocumentRows(t *testing.T) {
	docs := []bson.M{{"name": "a"}, {"name": "b"}}
	fields, rows, err := documentRows(docs)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(fields) != 2 || fields[0].Name != "index" || fields[1].Name != "document" {
		t.Fatalf("unexpected fields: %+v", fields)
	}
	if len(rows) != 2 || rows[0][0] != int64(0) || rows[1][0] != int64(1) {
		t.Fatalf("unexpected rows: %+v", rows)
	}
}

func TestEffectiveLimit(t *testing.T) {
	limit := "100"
	invalid := "abc"
	zero := "0"

	tests := []struct {
		name     string
		options  *database.QueryOptions
		expected int64
	}{
		{"nil options", nil, 0},
		{"valid limit", &database.QueryOptions{Limit: &limit}, 100},
		{"invalid limit", &database.QueryOptions{Limit: &invalid}, 0},
		{"zero limit", &database.QueryOptions{Limit: &zero}, 0},
		{"no limit set", &database.QueryOptions{}, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := effectiveLimit(tt.options); got != tt.expected {
				t.Errorf("expected %d, got %d", tt.expected, got)
			}
		})
	}
}

func TestGetManyGetOneNotSupported(t *testing.T) {
	adapter := NewAdapter()
	ctx := database.Context{}

	if _, err := adapter.GetMany(ctx, "users.find({})", nil); !errors.Is(err, ErrNotSupported) {
		t.Errorf("expected ErrNotSupported, got %v", err)
	}
	if _, err := adapter.GetOne(ctx, "users.find({})"); !errors.Is(err, ErrNotSupported) {
		t.Errorf("expected ErrNotSupported, got %v", err)
	}
}

func TestGetFieldsInvalidQuery(t *testing.T) {
	adapter := NewAdapter()
	ctx := database.Context{Connection: &database.Connection{Config: database.Config{Host: "localhost", Port: 27017}}}

	if _, err := adapter.GetFields(ctx, "not a command", nil); !errors.Is(err, ErrInvalidSyntax) {
		t.Errorf("expected ErrInvalidSyntax, got %v", err)
	}
}
