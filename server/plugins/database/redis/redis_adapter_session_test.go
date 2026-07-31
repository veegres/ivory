package redis

import (
	"errors"
	"ivory/plugins/database"
	"testing"
)

func TestCancelNotSupported(t *testing.T) {
	adapter := NewAdapter()
	if err := adapter.Cancel(database.Context{}, 1); !errors.Is(err, ErrNotSupported) {
		t.Errorf("expected ErrNotSupported, got %v", err)
	}
}

func TestParseClientList(t *testing.T) {
	text := "id=3 addr=127.0.0.1:52368 laddr=127.0.0.1:6379 name= age=10 idle=0 db=0 cmd=client|list user=default\n" +
		"id=7 addr=127.0.0.1:52400 laddr=127.0.0.1:6379 name=worker age=120 idle=5 db=1 cmd=get user=default\n"

	fields, rows := parseClientList(text)

	if len(fields) == 0 || fields[0].Name != "pid" {
		t.Fatalf("expected first field to be pid, got %+v", fields)
	}
	if len(rows) != 2 {
		t.Fatalf("expected 2 rows, got %d", len(rows))
	}
	if rows[0][0] != int64(3) {
		t.Errorf("expected pid 3, got %v", rows[0][0])
	}
	if rows[1][0] != int64(7) || rows[1][2] != "worker" {
		t.Errorf("unexpected second row: %v", rows[1])
	}
}

func TestParseClientListEmpty(t *testing.T) {
	_, rows := parseClientList("")
	if len(rows) != 0 {
		t.Errorf("expected no rows for empty input, got %v", rows)
	}
}

func TestParseClientAttrs(t *testing.T) {
	attrs := parseClientAttrs("id=3 addr=127.0.0.1:1 name=foo age=10")
	if attrs["id"] != "3" || attrs["addr"] != "127.0.0.1:1" || attrs["name"] != "foo" || attrs["age"] != "10" {
		t.Errorf("unexpected attrs: %v", attrs)
	}
}
