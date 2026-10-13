package utils

import (
	"errors"
	"testing"
)

func TestCheck(t *testing.T) {
	failing := func() error { return errors.New("refused") }
	passing := func() error { return nil }

	tests := []struct {
		name     string
		build    func(c *Check)
		expected []CheckStep
	}{
		{
			name: "every step done",
			build: func(c *Check) {
				c.Run(CheckStageVault, passing)
				c.RunAfter(CheckStageCert, passing)
			},
			expected: []CheckStep{{Stage: CheckStageVault, Result: CheckResultDone}, {Stage: CheckStageCert, Result: CheckResultDone}},
		},
		{
			name: "a failure carries its error and skips what depends on it",
			build: func(c *Check) {
				c.Run(CheckStageVault, failing)
				c.RunAfter(CheckStageCert, passing)
			},
			expected: []CheckStep{
				{Stage: CheckStageVault, Result: CheckResultFailed, Message: "refused"},
				{Stage: CheckStageCert, Result: CheckResultSkipped, Message: "vault failed"},
			},
		},
		{
			name: "an independent step still runs after a failure",
			build: func(c *Check) {
				c.Run(CheckStageVault, failing)
				c.Run(CheckStageCert, passing)
				c.RunAfter(CheckStageConnect, passing)
			},
			expected: []CheckStep{
				{Stage: CheckStageVault, Result: CheckResultFailed, Message: "refused"},
				{Stage: CheckStageCert, Result: CheckResultDone},
				{Stage: CheckStageConnect, Result: CheckResultSkipped, Message: "vault failed"},
			},
		},
		{
			name: "a skipped step does not stop what depends on it",
			build: func(c *Check) {
				c.Skip(CheckStageVault, "nothing to check")
				c.RunAfter(CheckStageCert, passing)
			},
			expected: []CheckStep{
				{Stage: CheckStageVault, Result: CheckResultSkipped, Message: "nothing to check"},
				{Stage: CheckStageCert, Result: CheckResultDone},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewCheck()
			tt.build(c)
			steps := c.Steps()
			if len(steps) != len(tt.expected) {
				t.Fatalf("expected %d steps, got %+v", len(tt.expected), steps)
			}
			for i, want := range tt.expected {
				got := steps[i]
				got.DurationMs = 0
				if got != want {
					t.Errorf("step %d = %+v, want %+v", i, got, want)
				}
			}
		})
	}

	t.Run("a step that never ran is not called", func(t *testing.T) {
		c := NewCheck()
		c.Run(CheckStageVault, failing)
		called := false
		c.RunAfter(CheckStageCert, func() error { called = true; return nil })
		if called {
			t.Fatalf("expected the dependent step not to run")
		}
	})
}
