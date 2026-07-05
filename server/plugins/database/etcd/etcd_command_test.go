package etcd

import (
	"errors"
	"testing"
)

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
