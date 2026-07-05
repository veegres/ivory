package console

import "io"

// Command is a generic interface for executing commands on a remote or local console.
// NOTE: Implementations of this interface are not guaranteed to be thread-safe.
type Command interface {
	// Id returns a stable string used to identify the job.
	Id() string
	// KeepAlive keeps job alive even if there is no subscribers.
	KeepAlive() bool
	// Persist returns true if the job output should be saved to a file.
	Persist() bool
	// Start begins the command and returns a reader over its output.
	Start() (io.Reader, error)
	// Wait blocks until the command exits and returns its exit error.
	Wait() error
	// Abort kills the running command immediately.
	Abort() error
	// Execute runs the command synchronously and returns the output as string array.
	Execute() ([]string, error)
}
