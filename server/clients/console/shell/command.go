package shell

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os/exec"
	"strings"
)

type Command struct {
	Name string
	Args []string

	JobID        string
	JobPersist   bool
	JobKeepAlive bool

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

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("stdout: %w", err)
	}
	cmd.Stderr = cmd.Stdout
	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("start: %w", err)
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

	var output []string
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		output = append(output, scanner.Text())
	}

	errWait := c.Wait()
	if errWait != nil {
		if exitErr, ok := errWait.(*exec.ExitError); ok {
			return output, fmt.Errorf("exit code %d: %s", exitErr.ExitCode(), strings.Join(output, "\n"))
		}
		return output, errWait
	}

	return output, nil
}
