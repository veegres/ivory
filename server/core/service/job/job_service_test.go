package job

import (
	"ivory/clients/storage"
	"os"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestService_Start_ReturnsExistingJobIDForSameCommand(t *testing.T) {
	s := NewService(nil)
	cmd := &MockCommand{id: "dup-job", output: []string{"line 1"}, delay: 20 * time.Millisecond}

	first, err := s.Start(cmd)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	second, err := s.Start(cmd)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if first != second {
		t.Fatalf("expected the same job id, got %s and %s", first, second)
	}

	s.mu.Lock()
	count := len(s.jobs)
	s.mu.Unlock()
	if count != 1 {
		t.Fatalf("expected a single tracked job, got %d", count)
	}

	if err := s.Stop(first); err != nil {
		t.Fatalf("failed to stop job: %v", err)
	}
}

func TestService_Stop_NotFound(t *testing.T) {
	s := NewService(nil)
	if err := s.Stop("missing-job"); err == nil {
		t.Fatal("expected an error stopping an unknown job")
	}
}

func TestService_Stop_StopsRunningJob(t *testing.T) {
	s := NewService(nil)
	cmd := &MockCommand{id: "stoppable-job", output: []string{"line 1", "line 2"}, delay: 10 * time.Millisecond}

	jobID, err := s.Start(cmd)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	time.Sleep(5 * time.Millisecond)
	if err := s.Stop(jobID); err != nil {
		t.Fatalf("expected no error stopping the job, got %v", err)
	}

	for i := 0; i < 100; i++ {
		if s.Status(jobID) == UNKNOWN {
			return
		}
		time.Sleep(2 * time.Millisecond)
	}
	t.Fatal("expected job to eventually be removed after stopping")
}

func TestService_Subscribe_NotFound(t *testing.T) {
	s := NewService(nil)
	if _, err := s.Subscribe("missing-job", "sub-1"); err == nil {
		t.Fatal("expected an error subscribing to an unknown job")
	}
}

func TestService_Subscribe_JobNotRunning(t *testing.T) {
	s := NewService(nil)
	cmd := &MockCommand{id: "finished-job"}
	job := NewJob(cmd, nil)
	job.setStatus(FINISHED)
	s.mu.Lock()
	s.unsafeAddJob(JobID(cmd.id), job)
	s.mu.Unlock()

	if _, err := s.Subscribe(JobID(cmd.id), "sub-1"); err == nil {
		t.Fatal("expected an error subscribing to a finished job")
	}
}

func TestService_Status_UnknownForMissingJob(t *testing.T) {
	s := NewService(nil)
	if status := s.Status("missing-job"); status != UNKNOWN {
		t.Fatalf("expected UNKNOWN status, got %s", status)
	}
}

func TestService_GetLogsPath(t *testing.T) {
	t.Run("returns an error when storage is not initialized", func(t *testing.T) {
		s := NewService(nil)
		if _, err := s.GetLogsPath("some-job"); err != ErrStorageNotInitialized {
			t.Fatalf("expected ErrStorageNotInitialized, got %v", err)
		}
	})

	t.Run("returns the path when storage is available", func(t *testing.T) {
		oldWd, _ := os.Getwd()
		tmpDir := t.TempDir()
		_ = os.Chdir(tmpDir)
		defer os.Chdir(oldWd)

		fs := storage.NewFileStorage("logs-path-test", ".log")
		s := NewService(fs)

		path, err := s.GetLogsPath("some-job")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if path == "" {
			t.Fatal("expected a non-empty path")
		}
	})
}

func TestService_Stream_JobNotFoundStopsAfterServerError(t *testing.T) {
	s := NewService(nil)

	var received []Message
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		s.Stream("missing-job", "sub-1", make(<-chan struct{}), func(msg Message) {
			received = append(received, msg)
		})
	}()
	wg.Wait()

	if len(received) == 0 || received[0].Type != STREAM || received[0].Message != START.String() {
		t.Fatalf("expected a STREAM START event first, got %v", received)
	}
	last := received[len(received)-1]
	if last.Type != STREAM || last.Message != END.String() {
		t.Fatalf("expected a STREAM END event last, got %v", received)
	}

	foundConsoleError := false
	for _, msg := range received {
		if msg.Type == SERVER && msg.Message == "streaming from the console error: job missing-job not found" {
			foundConsoleError = true
		}
		if msg.Type == LOG {
			t.Fatalf("did not expect log messages for a missing job, got %v", msg)
		}
	}
	if !foundConsoleError {
		t.Fatalf("expected a console error message, got %v", received)
	}
}

func TestService_Stream_SkipsFileWhenMissing(t *testing.T) {
	oldWd, _ := os.Getwd()
	tmpDir := t.TempDir()
	_ = os.Chdir(tmpDir)
	defer os.Chdir(oldWd)

	fs := storage.NewFileStorage("stream-skip-test", ".log")
	s := NewService(fs)

	var received []Message
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		s.Stream("no-such-job", "sub-1", make(<-chan struct{}), func(msg Message) {
			received = append(received, msg)
		})
	}()
	wg.Wait()

	foundSkip := false
	for _, msg := range received {
		if msg.Type == SERVER && strings.HasPrefix(msg.Message, "streaming from the file skipped, error:") {
			foundSkip = true
		}
	}
	if !foundSkip {
		t.Fatalf("expected a file-skipped message when no log file exists, got %v", received)
	}
}
