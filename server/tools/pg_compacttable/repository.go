package pg_compacttable

import (
	"ivory/clients/storage"
	"ivory/core/service/job"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Repository struct {
	bucket *storage.DbBucket[Response]
}

func NewRepository(bucket *storage.DbBucket[Response]) *Repository {
	return &Repository{
		bucket: bucket,
	}
}

func (r *Repository) List() ([]Response, error) {
	return r.bucket.GetList(nil, r.sortDescByCreatedAt)
}

func (r *Repository) ListByStatus(status job.Status) ([]Response, error) {
	return r.bucket.GetList(func(model Response) bool {
		return model.Status == status
	}, r.sortDescByCreatedAt)
}

func (r *Repository) ListByCluster(cluster string) ([]Response, error) {
	return r.bucket.GetList(func(model Response) bool {
		return model.Cluster == cluster
	}, r.sortDescByCreatedAt)
}

func (r *Repository) Get(uuid uuid.UUID) (Response, error) {
	return r.bucket.Get(uuid.String())
}

func (r *Repository) Create(cluster string, vaultId *uuid.UUID, args []string) (*Response, error) {
	jobUuid := uuid.New()

	compactTableModel := Response{
		Uuid:        jobUuid,
		VaultId:     vaultId,
		Cluster:     cluster,
		Status:      job.PENDING,
		Command:     "pgcompacttable " + strings.Join(args, " "),
		CommandArgs: args,
		JobId:       job.JobID(jobUuid.String()),
		CreatedAt:   time.Now().UnixNano(),
	}

	err := r.bucket.Update(jobUuid.String(), compactTableModel)
	return &compactTableModel, err
}

func (r *Repository) UpdateStatus(compactTable Response, status job.Status) error {
	tmp := compactTable
	tmp.Status = status
	return r.bucket.Update(tmp.Uuid.String(), tmp)
}

func (r *Repository) UpdateJobId(compactTable Response, id job.JobID) error {
	tmp := compactTable
	tmp.JobId = id
	return r.bucket.Update(tmp.Uuid.String(), tmp)
}

func (r *Repository) Delete(uuid uuid.UUID) error {
	return r.bucket.Delete(uuid.String())
}

func (r *Repository) DeleteAll() error {
	return r.bucket.DeleteAll()
}

func (r *Repository) sortDescByCreatedAt(list []Response, i, j int) bool {
	return list[i].CreatedAt > list[j].CreatedAt
}
