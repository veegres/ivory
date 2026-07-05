package job

import (
	"fmt"
	"io"
	"ivory/clients/storage"
	"os"
	"sync"
	"testing"
	"time"
)

type MockCommand struct {
	id        string
	persist   bool
	keepAlive bool
	output    []string
	delay     time.Duration
	startErr  error
	waitErr   error
	abortErr  error

	// Internals
	r       *io.PipeReader
	w       *io.PipeWriter
	running bool
	mu      sync.Mutex
}

func (c *MockCommand) Id() string      { return c.id }
func (c *MockCommand) KeepAlive() bool { return c.keepAlive }
func (c *MockCommand) Persist() bool   { return c.persist }
func (c *MockCommand) Start() (io.Reader, error) {
	if c.startErr != nil {
		return nil, c.startErr
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	c.r, c.w = io.Pipe()
	c.running = true
	go func() {
		for _, line := range c.output {
			if c.delay > 0 {
				time.Sleep(c.delay)
			}
			c.mu.Lock()
			if !c.running {
				c.mu.Unlock()
				return
			}
			_, _ = c.w.Write([]byte(line + "\n"))
			c.mu.Unlock()
		}
		c.mu.Lock()
		if c.running {
			_ = c.w.Close()
		}
		c.mu.Unlock()
	}()
	return c.r, nil
}

func (c *MockCommand) Wait() error {
	return c.waitErr
}

func (c *MockCommand) Abort() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.running = false
	if c.w != nil {
		_ = c.w.CloseWithError(fmt.Errorf("aborted"))
	}
	return c.abortErr
}

func (c *MockCommand) Execute() ([]string, error) {
	return c.output, nil
}

func TestJob_Run_Success(t *testing.T) {
	oldWd, _ := os.Getwd()
	tmpDir := t.TempDir()
	_ = os.Chdir(tmpDir)
	defer os.Chdir(oldWd)

	storage := storage.NewFileStorage("job-tests", ".log")
	cmd := &MockCommand{
		id:      "test-job-success",
		persist: true,
		output:  []string{"line 1", "line 2"},
	}

	job := NewJob(cmd, storage)
	ch, ok := job.addSubscriber("sub-1")
	if !ok {
		t.Fatal("expected subscriber to be added")
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		job.Run()
	}()

	var received []Message
	for msg := range ch {
		received = append(received, msg)
	}

	wg.Wait()

	// Verify status progression: PENDING/RUNNING -> FINISHED/FAILED
	// Wait, we get STATUS message from setStatus when it transitions.
	// Let's verify we got the log lines and status events.
	hasRunning := false
	hasFinished := false
	var logLines []string

	for _, msg := range received {
		if msg.Type == STATUS {
			if msg.Message == RUNNING.String() {
				hasRunning = true
			} else if msg.Message == FINISHED.String() {
				hasFinished = true
			}
		} else if msg.Type == LOG {
			logLines = append(logLines, msg.Message)
		}
	}

	if !hasRunning {
		t.Error("Expected to receive RUNNING status event")
	}
	if !hasFinished {
		t.Error("Expected to receive FINISHED status event")
	}
	if len(logLines) != 2 || logLines[0] != "line 1" || logLines[1] != "line 2" {
		t.Errorf("Expected logs ['line 1', 'line 2'], got %v", logLines)
	}

	// Verify persistence
	persisted, err := storage.ReadByName("test-job-success")
	if err != nil {
		t.Fatalf("Failed to read persisted logs: %v", err)
	}
	expectedPersisted := "line 1\nline 2\n"
	if string(persisted) != expectedPersisted {
		t.Errorf("Expected persisted logs %q, got %q", expectedPersisted, string(persisted))
	}
}

func TestJob_Run_StartError(t *testing.T) {
	cmd := &MockCommand{
		id:       "test-job-starterr",
		startErr: fmt.Errorf("start failed"),
	}

	job := NewJob(cmd, nil)
	ch, ok := job.addSubscriber("sub-1")
	if !ok {
		t.Fatal("expected subscriber to be added")
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		job.Run()
	}()

	var received []Message
	for msg := range ch {
		received = append(received, msg)
	}
	wg.Wait()

	hasFailed := false
	hasServerError := false
	for _, msg := range received {
		if msg.Type == STATUS && msg.Message == FAILED.String() {
			hasFailed = true
		}
		if msg.Type == SERVER && msg.Message == "start failed" {
			hasServerError = true
		}
	}

	if !hasFailed {
		t.Error("Expected status to be FAILED")
	}
	if !hasServerError {
		t.Error("Expected to receive SERVER error message")
	}
}

func TestJob_Run_Stop(t *testing.T) {
	cmd := &MockCommand{
		id:     "test-job-stop",
		output: []string{"line 1", "line 2", "line 3"},
		delay:  10 * time.Millisecond,
	}

	job := NewJob(cmd, nil)
	ch, ok := job.addSubscriber("sub-1")
	if !ok {
		t.Fatal("expected subscriber to be added")
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		job.Run()
	}()

	// Let it start and then stop it
	time.Sleep(5 * time.Millisecond)
	errStop := job.Stop()
	if errStop != nil {
		t.Fatalf("Stop failed: %v", errStop)
	}

	var received []Message
	for msg := range ch {
		received = append(received, msg)
	}
	wg.Wait()

	hasStopped := false
	for _, msg := range received {
		if msg.Type == STATUS && msg.Message == STOPPED.String() {
			hasStopped = true
		}
	}

	if !hasStopped {
		t.Error("Expected status to become STOPPED")
	}
}

func TestJob_Run_KillerStopsJobWithoutSubscribers(t *testing.T) {
	cmd := &MockCommand{
		id:     "test-job-killer-stop",
		output: []string{"line 1", "line 2", "line 3"},
		delay:  50 * time.Millisecond,
	}

	job := NewJob(cmd, nil)
	job.keepAliveDuration = time.Millisecond
	job.keepAliveCheckInterval = time.Millisecond

	done := make(chan struct{})
	go func() {
		defer close(done)
		job.Run()
	}()

	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("expected killer to stop job without deadlocking")
	}

	if status := job.getStatus(); status != STOPPED {
		t.Fatalf("expected status %s, got %s", STOPPED, status)
	}
}

func TestJob_AddSubscriber_FailsWhenJobFinished(t *testing.T) {
	cmd := &MockCommand{id: "test-job-finished"}
	job := NewJob(cmd, nil)
	job.setStatus(FINISHED)

	ch, ok := job.addSubscriber("sub-1")
	if ok {
		t.Fatal("expected subscriber add to fail for finished job")
	}
	if ch != nil {
		t.Fatal("expected nil channel for failed subscriber add")
	}
}

func TestManager(t *testing.T) {
	oldWd, _ := os.Getwd()
	tmpDir := t.TempDir()
	_ = os.Chdir(tmpDir)
	defer os.Chdir(oldWd)

	storage := storage.NewFileStorage("manager-tests", ".log")
	mgr := NewService(storage)

	cmd := &MockCommand{
		id:      "mgr-job-1",
		persist: true,
		output:  []string{"log A", "log B"},
		delay:   5 * time.Millisecond,
	}

	jobID, err := mgr.Start(cmd)
	if err != nil {
		t.Fatalf("Failed to start job: %v", err)
	}

	if jobID != "mgr-job-1" {
		t.Errorf("Expected jobID 'mgr-job-1', got %s", jobID)
	}

	// Verify status is RUNNING/PENDING
	status := mgr.Status(jobID)
	if status != RUNNING && status != PENDING {
		t.Errorf("Expected status to be RUNNING or PENDING, got %s", status)
	}

	// Wait for the job to finish by polling status
	for i := 0; i < 50; i++ {
		if mgr.Status(jobID) == UNKNOWN {
			break
		}
		time.Sleep(5 * time.Millisecond)
	}

	// Stream logs
	var streamedMsgs []Message
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		mgr.Stream(jobID, "stream-sub", make(<-chan struct{}), func(msg Message) {
			streamedMsgs = append(streamedMsgs, msg)
		})
	}()
	wg.Wait()

	hasStartEvent := false
	hasEndEvent := false
	var logs []string

	for _, msg := range streamedMsgs {
		if msg.Type == STREAM {
			if msg.Message == START.String() {
				hasStartEvent = true
			} else if msg.Message == END.String() {
				hasEndEvent = true
			}
		} else if msg.Type == LOG {
			logs = append(logs, msg.Message)
		}
	}

	if !hasStartEvent {
		t.Error("Expected STREAM START event")
	}
	if !hasEndEvent {
		t.Error("Expected STREAM END event")
	}
	if len(logs) != 2 || logs[0] != "log A" || logs[1] != "log B" {
		t.Errorf("Expected logs ['log A', 'log B'], got %v", logs)
	}
}
