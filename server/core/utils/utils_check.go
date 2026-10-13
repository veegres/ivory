package utils

import "time"

// CheckResult is how one step of a Check ended. Skipped means the step was not
// tried: a step it depends on failed, or there was nothing to check.
type CheckResult string

const (
	CheckResultDone    CheckResult = "done"
	CheckResultFailed  CheckResult = "failed"
	CheckResultSkipped CheckResult = "skipped"
)

// CheckStage is how far reaching something got: its credentials out of the
// vault, its certificates, then the connection itself. Which thing was
// reached is said by whoever holds the steps, never by the stage.
type CheckStage string

const (
	CheckStageVault   CheckStage = "vault"
	CheckStageCert    CheckStage = "cert"
	CheckStageConnect CheckStage = "connect"
)

// CheckStep is one step of a Check, in the order it was taken.
type CheckStep struct {
	Stage      CheckStage  `json:"stage"`
	Result     CheckResult `json:"result"`
	DurationMs int64       `json:"durationMs"`
	Message    string      `json:"message,omitempty"`
}

// Check records the ordered, individually timed steps of reaching something,
// so a failure says at which step it happened rather than only that it did.
type Check struct {
	steps []CheckStep
}

func NewCheck() *Check {
	return &Check{steps: make([]CheckStep, 0)}
}

// Run takes a step that depends on nothing recorded before it.
func (c *Check) Run(stage CheckStage, step func() error) {
	start := time.Now()
	err := step()
	result := CheckStep{Stage: stage, Result: CheckResultDone, DurationMs: time.Since(start).Milliseconds()}
	if err != nil {
		result.Result = CheckResultFailed
		result.Message = err.Error()
	}
	c.steps = append(c.steps, result)
}

// RunAfter takes a step that depends on every step recorded before it, and
// skips it naming the first of them that failed.
func (c *Check) RunAfter(stage CheckStage, step func() error) {
	for _, p := range c.steps {
		if p.Result == CheckResultFailed {
			c.Skip(stage, string(p.Stage)+" failed")
			return
		}
	}
	c.Run(stage, step)
}

// Skip records a step that had nothing to check.
func (c *Check) Skip(stage CheckStage, reason string) {
	c.steps = append(c.steps, CheckStep{Stage: stage, Result: CheckResultSkipped, Message: reason})
}

func (c *Check) Steps() []CheckStep {
	return c.steps
}
