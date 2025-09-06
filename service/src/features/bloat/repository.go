package bloat

import (
	"errors"
	. "ivory/src/features/bloat/job"
	"ivory/src/storage/db"
	"ivory/src/storage/files"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Repository struct {
	bucket *db.Bucket[Bloat]
	file   *files.Storage
}

func NewRepository(bucket *db.Bucket[Bloat], file *files.Storage) *Repository {
	return &Repository{
		bucket: bucket,
		file:   file,
	}
}

func (r *Repository) List() ([]Bloat, error) {
	return r.bucket.GetList(nil, r.sortDescByCreatedAt)
}

func (r *Repository) ListByStatus(status JobStatus) ([]Bloat, error) {
	return r.bucket.GetList(func(model Bloat) bool {
		return model.Status == status
	}, r.sortDescByCreatedAt)
}

func (r *Repository) ListByCluster(cluster string) ([]Bloat, error) {
	return r.bucket.GetList(func(model Bloat) bool {
		return model.Cluster == cluster
	}, r.sortDescByCreatedAt)
}

func (r *Repository) Get(uuid uuid.UUID) (Bloat, error) {
	return r.bucket.Get(uuid.String())
}

func (r *Repository) GetOpenFile(uuid uuid.UUID) (*os.File, error) {
	return r.file.OpenByName(uuid.String())
}

func (r *Repository) Create(credentialId uuid.UUID, cluster string, args []string) (*Bloat, error) {
	jobUuid := uuid.New()
	logsPath, errCreate := r.file.CreateByName(jobUuid.String())
	if errCreate != nil {
		return nil, errCreate
	}

	compactTableModel := Bloat{
		Uuid:         jobUuid,
		CredentialId: credentialId,
		Cluster:      cluster,
		Status:       PENDING,
		Command:      "pgcompacttable " + strings.Join(args, " "),
		CommandArgs:  args,
		LogsPath:     logsPath,
		CreatedAt:    time.Now().UnixNano(),
	}

	err := r.bucket.Update(jobUuid.String(), compactTableModel)
	return &compactTableModel, err
}

func (r *Repository) UpdateStatus(compactTable Bloat, status JobStatus) error {
	tmp := compactTable
	tmp.Status = status
	return r.bucket.Update(tmp.Uuid.String(), tmp)
}

func (r *Repository) Delete(uuid uuid.UUID) error {
	// NOTE: we shouldn't check error here, if there is no file we should try to remove info
	_ = r.file.DeleteByName(uuid.String())
	return r.bucket.Delete(uuid.String())
}

func (r *Repository) DeleteAll() error {
	errFile := r.file.DeleteAll()
	errBucket := r.bucket.DeleteAll()
	return errors.Join(errFile, errBucket)
}

func (r *Repository) sortDescByCreatedAt(list []Bloat, i, j int) bool {
	return list[i].CreatedAt > list[j].CreatedAt
}
