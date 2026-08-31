package clickhouse

import (
	"errors"
	"fmt"
	"io"
	"testing"
)

// TestIsEmptyResultSet holds the one distinction the console depends on. The
// native protocol answers a statement with no result set - every INSERT, and
// every DDL that is not ON CLUSTER - by ending the stream, and the driver
// reports that as a bare io.EOF even though the statement has already been
// applied. A read that genuinely failed carries the same io.EOF wrapped in
// context, and there the statement's fate is unknown, so only the bare one may
// be read as success.
func TestIsEmptyResultSet(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{name: "end of stream with no data block", err: io.EOF, expected: true},
		{name: "no error at all", err: nil, expected: false},
		{
			name:     "a read that failed, wrapped by the driver",
			err:      fmt.Errorf("query processing: failed to read first block packet from 10.0.0.1:9000 (conn_id=7): %w", io.EOF),
			expected: false,
		},
		{name: "an error the server sent", err: errors.New("code: 47, message: Unknown expression identifier"), expected: false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := isEmptyResultSet(test.err); got != test.expected {
				t.Errorf("expected %v, got %v", test.expected, got)
			}
		})
	}
}
