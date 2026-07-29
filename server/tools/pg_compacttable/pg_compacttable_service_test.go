package pg_compacttable

import (
	"io"
	"ivory/clients/console/shell"
	"ivory/clients/console/ssh"
	"ivory/clients/storage"
	"ivory/core/service/encryption"
	"ivory/core/service/job"
	"ivory/core/service/secret"
	"ivory/core/service/vault"
	"ivory/plugins/database"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/boltdb/bolt"
	"github.com/google/uuid"
)

func TestSupportedFeatures(t *testing.T) {
	// NOTE: method under test does not touch dependencies, so the zero
	// service is enough and avoids the initializer goroutine in NewService
	s := &Service{}

	t.Run("postgres gets all features", func(t *testing.T) {
		features := s.SupportedFeatures(database.POSTGRES)
		for feature, supported := range features {
			if !supported {
				t.Errorf("expected feature %v to be supported for postgres", feature)
			}
		}
		if len(features) != 4 {
			t.Fatalf("expected 4 features for postgres, got %d", len(features))
		}
	})

	t.Run("etcd gets no supported features", func(t *testing.T) {
		features := s.SupportedFeatures(database.ETCD)
		for feature, supported := range features {
			if supported {
				t.Errorf("expected feature %v to not be supported for etcd", feature)
			}
		}
		if len(features) != 4 {
			t.Fatalf("expected 4 features for etcd, got %d", len(features))
		}
	})
}

// fakeJobCommand is a minimal console.Command used to seed a real, controllable
// job into job.Service without needing to spawn a real "pgcompacttable" binary.
type fakeJobCommand struct {
	id      string
	r       *io.PipeReader
	w       *io.PipeWriter
	aborted bool
}

func (c *fakeJobCommand) Id() string      { return c.id }
func (c *fakeJobCommand) KeepAlive() bool { return true }
func (c *fakeJobCommand) Persist() bool   { return false }
func (c *fakeJobCommand) Start() (io.Reader, error) {
	c.r, c.w = io.Pipe()
	return c.r, nil
}
func (c *fakeJobCommand) Wait() error { return nil }
func (c *fakeJobCommand) Abort() error {
	c.aborted = true
	if c.w != nil {
		_ = c.w.Close()
	}
	return nil
}
func (c *fakeJobCommand) Execute() ([]string, error) { return nil, nil }

type testCompactTableEnv struct {
	service      *Service
	repository   *Repository
	jobManager   *job.Service
	vaultService *vault.Service
}

func createTestCompactTableService(t *testing.T) *testCompactTableEnv {
	t.Helper()

	tmpDir, errDir := os.MkdirTemp("", "pg-compacttable-service-test-*")
	if errDir != nil {
		t.Fatalf("failed to create temp dir: %v", errDir)
	}
	t.Cleanup(func() {
		os.RemoveAll(tmpDir)
	})

	db, errOpen := bolt.Open(filepath.Join(tmpDir, "test.db"), 0600, nil)
	if errOpen != nil {
		t.Fatalf("failed to open test database: %v", errOpen)
	}
	t.Cleanup(func() {
		db.Close()
	})

	oldWd, _ := os.Getwd()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to chdir: %v", err)
	}
	t.Cleanup(func() { os.Chdir(oldWd) })

	secretService := secret.NewService(
		secret.NewRepository(storage.NewDbBucket[string](db, "Secret")),
		encryption.NewService(),
	)
	if err := secretService.SetDefault(); err != nil {
		t.Fatalf("failed to set default secret: %v", err)
	}
	vaultService := vault.NewService(
		vault.NewRepository(storage.NewDbBucket[vault.Vault](db, "Vault")),
		ssh.NewClient(),
		secretService,
		encryption.NewService(),
	)

	jobManager := job.NewService(storage.NewFileStorage("jobs", ".log"))
	repository := NewRepository(storage.NewDbBucket[Response](db, "PgCompactTable"))
	// NOTE: built directly (not via NewService) to avoid its background
	// initializer() goroutine racing with this test's db.Close/RemoveAll cleanup.
	service := &Service{
		bloatRepository: repository,
		shellClient:     shell.NewClient(),
		vaultService:    vaultService,
		jobManager:      jobManager,
	}

	return &testCompactTableEnv{service: service, repository: repository, jobManager: jobManager, vaultService: vaultService}
}

func TestServiceStartCreatesAPendingRecord(t *testing.T) {
	env := createTestCompactTableService(t)

	created, err := env.service.Start("cluster1", nil, []string{"--target", "public"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if created.Cluster != "cluster1" {
		t.Fatalf("expected cluster 'cluster1', got %q", created.Cluster)
	}
	if created.Command != "pgcompacttable --target public" {
		t.Fatalf("unexpected command: %q", created.Command)
	}
}

func TestServiceStartFailsImmediatelyWithBadVault(t *testing.T) {
	env := createTestCompactTableService(t)
	badVaultId := uuid.New()

	created, err := env.service.Start("cluster1", &badVaultId, []string{})
	if err != nil {
		t.Fatalf("expected Start itself to succeed (the record is created first), got %v", err)
	}

	got, errGet := env.service.Get(created.Uuid)
	if errGet != nil {
		t.Fatalf("expected no error, got %v", errGet)
	}
	if got.Status != job.FAILED {
		t.Fatalf("expected the record to fail synchronously due to the bad vault id, got %v", got.Status)
	}
}

func TestServiceListMethods(t *testing.T) {
	env := createTestCompactTableService(t)

	if _, err := env.repository.Create("cluster1", nil, []string{}); err != nil {
		t.Fatalf("failed to seed: %v", err)
	}
	c2, err := env.repository.Create("cluster2", nil, []string{})
	if err != nil {
		t.Fatalf("failed to seed: %v", err)
	}
	if err := env.repository.UpdateStatus(*c2, job.FINISHED); err != nil {
		t.Fatalf("failed to update status: %v", err)
	}

	t.Run("List returns everything", func(t *testing.T) {
		list, err := env.service.List()
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(list) != 2 {
			t.Fatalf("expected 2 records, got %d", len(list))
		}
	})

	t.Run("ListByStatus filters", func(t *testing.T) {
		list, err := env.service.ListByStatus(job.FINISHED)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(list) != 1 || list[0].Cluster != "cluster2" {
			t.Fatalf("expected only cluster2, got %v", list)
		}
	})

	t.Run("ListByCluster filters", func(t *testing.T) {
		list, err := env.service.ListByCluster("cluster1")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(list) != 1 || list[0].Cluster != "cluster1" {
			t.Fatalf("expected only cluster1, got %v", list)
		}
	})

	t.Run("Get returns a single record", func(t *testing.T) {
		got, err := env.service.Get(c2.Uuid)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if got.Cluster != "cluster2" {
			t.Fatalf("expected cluster2, got %q", got.Cluster)
		}
	})

	t.Run("Get fails for an unknown id", func(t *testing.T) {
		if _, err := env.service.Get(uuid.New()); err == nil {
			t.Fatalf("expected an error for an unknown id")
		}
	})
}

func TestServiceGetLogsPath(t *testing.T) {
	env := createTestCompactTableService(t)

	created, err := env.repository.Create("cluster1", nil, []string{})
	if err != nil {
		t.Fatalf("failed to seed: %v", err)
	}

	t.Run("known id resolves a path", func(t *testing.T) {
		path, errPath := env.service.GetLogsPath(created.Uuid)
		if errPath != nil {
			t.Fatalf("expected no error, got %v", errPath)
		}
		if path == "" {
			t.Fatalf("expected a non-empty path")
		}
	})

	t.Run("unknown id fails", func(t *testing.T) {
		if _, err := env.service.GetLogsPath(uuid.New()); err == nil {
			t.Fatalf("expected an error for an unknown id")
		}
	})
}

func TestServiceStream(t *testing.T) {
	env := createTestCompactTableService(t)

	t.Run("unknown job reports not found and unknown status", func(t *testing.T) {
		var received []job.Message
		env.service.Stream(uuid.New(), "sub-1", make(<-chan struct{}), func(msg job.Message) {
			received = append(received, msg)
		})
		if len(received) != 2 {
			t.Fatalf("expected 2 messages, got %v", received)
		}
		if received[0].Type != job.SERVER {
			t.Fatalf("expected a SERVER message first, got %v", received[0])
		}
		if received[1].Type != job.STATUS || received[1].Message != job.UNKNOWN.String() {
			t.Fatalf("expected an UNKNOWN status message, got %v", received[1])
		}
	})
}

func TestServiceDelete(t *testing.T) {
	env := createTestCompactTableService(t)

	pending, err := env.repository.Create("cluster1", nil, []string{})
	if err != nil {
		t.Fatalf("failed to seed: %v", err)
	}

	t.Run("active job cannot be deleted", func(t *testing.T) {
		if err := env.service.Delete(pending.Uuid); err != ErrJobIsActive {
			t.Fatalf("expected ErrJobIsActive, got %v", err)
		}
	})

	t.Run("inactive job can be deleted", func(t *testing.T) {
		if err := env.repository.UpdateStatus(*pending, job.FAILED); err != nil {
			t.Fatalf("failed to update status: %v", err)
		}
		if err := env.service.Delete(pending.Uuid); err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if _, err := env.service.Get(pending.Uuid); err == nil {
			t.Fatalf("expected the record to be gone")
		}
	})
}

func TestServiceDeleteAll(t *testing.T) {
	env := createTestCompactTableService(t)

	if _, err := env.repository.Create("cluster1", nil, []string{}); err != nil {
		t.Fatalf("failed to seed: %v", err)
	}
	if _, err := env.repository.Create("cluster2", nil, []string{}); err != nil {
		t.Fatalf("failed to seed: %v", err)
	}

	if err := env.service.DeleteAll(); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	list, err := env.service.List()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(list) != 0 {
		t.Fatalf("expected no records, got %v", list)
	}
}

func TestServiceStop(t *testing.T) {
	env := createTestCompactTableService(t)

	t.Run("unknown job fails", func(t *testing.T) {
		if err := env.service.Stop(uuid.New()); err != ErrNoSuchActiveJob {
			t.Fatalf("expected ErrNoSuchActiveJob, got %v", err)
		}
	})

	t.Run("inactive job fails", func(t *testing.T) {
		created, err := env.repository.Create("cluster1", nil, []string{})
		if err != nil {
			t.Fatalf("failed to seed: %v", err)
		}
		if err := env.repository.UpdateStatus(*created, job.FINISHED); err != nil {
			t.Fatalf("failed to update status: %v", err)
		}
		if err := env.service.Stop(created.Uuid); err != ErrNoSuchActiveJob {
			t.Fatalf("expected ErrNoSuchActiveJob, got %v", err)
		}
	})

	t.Run("active job with a real running job is stopped", func(t *testing.T) {
		cmd := &fakeJobCommand{id: "fake-job-1"}
		jobID, errStart := env.jobManager.Start(cmd)
		if errStart != nil {
			t.Fatalf("failed to start fake job: %v", errStart)
		}
		// Give the job goroutine a moment to reach RUNNING before we stop it.
		for i := 0; i < 100 && env.jobManager.Status(jobID) != job.RUNNING; i++ {
			time.Sleep(time.Millisecond)
		}

		created, errCreate := env.repository.Create("cluster1", nil, []string{})
		if errCreate != nil {
			t.Fatalf("failed to seed: %v", errCreate)
		}
		if err := env.repository.UpdateJobId(*created, jobID); err != nil {
			t.Fatalf("failed to set job id: %v", err)
		}
		withJobId, errGet := env.repository.Get(created.Uuid)
		if errGet != nil {
			t.Fatalf("failed to get: %v", errGet)
		}
		if err := env.repository.UpdateStatus(withJobId, job.RUNNING); err != nil {
			t.Fatalf("failed to set status: %v", err)
		}

		if err := env.service.Stop(created.Uuid); err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if !cmd.aborted {
			t.Fatalf("expected the underlying command to be aborted")
		}
	})
}
