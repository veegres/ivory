package shell

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"ivory/clients/console"
	"os/exec"
	"strings"
	"time"

	"github.com/creack/pty"
)

type Command struct {
	Name string
	Args []string

	JobID        string
	JobPersist   bool
	JobKeepAlive bool
	// ExecuteTimeout bounds Execute() (not Start/Wait used directly for
	// long-running streaming): zero means no bound.
	ExecuteTimeout time.Duration

	cmd *exec.Cmd
}

func (c *Command) Id() string {
	if c.JobID != "" {
		return c.JobID
	}
	h := sha256.Sum256(fmt.Appendf(nil, "shell|%s|%s", c.Name, strings.Join(c.Args, "\x00")))
	return hex.EncodeToString(h[:])
}

func (c *Command) Persist() bool {
	return c.JobPersist
}

func (c *Command) KeepAlive() bool {
	return c.JobKeepAlive
}

func (c *Command) Start() (io.Reader, error) {
	cmd := exec.Command(c.Name, c.Args...)

	// Allocate a PTY — this makes Docker think it's writing to a real terminal,
	// which disables the internal stdio buffering that can delaying logs.
	stdout, err := pty.Start(cmd)
	if err != nil {
		return nil, fmt.Errorf("stdout: %w", err)
	}

	c.cmd = cmd
	return stdout, nil
}

func (c *Command) Wait() error {
	if c.cmd == nil {
		return nil
	}
	return c.cmd.Wait()
}

func (c *Command) Abort() error {
	if c.cmd == nil || c.cmd.Process == nil {
		return nil
	}
	return c.cmd.Process.Kill()
}

func (c *Command) Execute() ([]string, error) {
	reader, err := c.Start()
	if err != nil {
		return nil, err
	}
	return console.Execute(reader, c.Wait, c.Abort, func(err error) (int, bool) {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			return exitErr.ExitCode(), true
		}
		return 0, false
	}, c.ExecuteTimeout)
}
