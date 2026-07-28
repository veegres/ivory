package console

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"
)

// ErrExecuteTimeout is returned by Execute when a command doesn't finish
// within its timeout and had to be aborted.
var ErrExecuteTimeout = errors.New("command execution timed out")

// Execute runs an already-started command to completion: it scans reader
// line by line, waits for the command to finish via wait, and formats a
// non-nil wait error into "exit code N: <output>" using exitCode (which
// must report a code only when the error is the transport's own exit-status
// error, e.g. *ssh.ExitError or *exec.ExitError - any other wait error is
// returned as-is). Every console.Command implementation shares this logic
// so a transport only supplies Start/Wait/Abort and its own exit-error type.
//
// If timeout is positive and the command doesn't finish in time, abort is
// called to kill it and ErrExecuteTimeout is returned, bounding a hung
// remote/local command (e.g. a stuck "docker pull") instead of blocking the
// caller - and therefore a whole synchronous cluster deploy - forever.
func Execute(reader io.Reader, wait func() error, abort func() error, exitCode func(error) (int, bool), timeout time.Duration) ([]string, error) {
	type result struct {
		output []string
		err    error
	}
	done := make(chan result, 1)
	go func() {
		var output []string
		scanner := bufio.NewScanner(reader)
		for scanner.Scan() {
			output = append(output, scanner.Text())
		}

		if errWait := wait(); errWait != nil {
			if code, ok := exitCode(errWait); ok {
				done <- result{output, fmt.Errorf("exit code %d: %s", code, strings.Join(output, "\n"))}
				return
			}
			done <- result{output, errWait}
			return
		}

		if errScanner := scanner.Err(); errScanner != nil {
			// NOTE: reading from a PTY master after the slave side closes
			// returns an EIO error when the process exits; expected, not a
			// real failure.
			if !strings.Contains(errScanner.Error(), "input/output error") {
				done <- result{output, errScanner}
				return
			}
		}
		done <- result{output, nil}
	}()

	if timeout <= 0 {
		r := <-done
		return r.output, r.err
	}
	select {
	case r := <-done:
		return r.output, r.err
	case <-time.After(timeout):
		_ = abort()
		return nil, ErrExecuteTimeout
	}
}
