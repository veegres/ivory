package console

import (
	"strings"
	"testing"
)

func TestExecuteOutputIsNeverNil(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		length int
	}{
		{"no output", "", 0},
		{"lines", "a\nb\n", 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, err := Execute(strings.NewReader(tt.input), func() error { return nil }, func() error { return nil }, func(error) (int, bool) { return 0, false }, 0)
			if err != nil {
				t.Fatal(err)
			}
			if output == nil || len(output) != tt.length {
				t.Errorf("got %#v, want a non-nil slice of %d", output, tt.length)
			}
		})
	}
}
