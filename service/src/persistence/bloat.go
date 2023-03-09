package persistence

import (
	"github.com/google/uuid"
	"ivory/src/config"
	. "ivory/src/model"
	"strings"
	"time"
)

type CompactTableRepository struct {
	bucket *config.Bucket[CompactTableModel]
	file   *config.FilePersistence
}

func NewCompactTableRepository(bucket *config.Bucket[CompactTableModel], file *config.FilePersistence) *CompactTableRepository {
	return &CompactTableRepository{
		bucket: bucket,
		file:   file,
	}
}

func (r *CompactTableRepository) List() ([]CompactTableModel, error) {
	return r.bucket.GetList(nil, r.sortDescByCreatedAt)
}

func (r *CompactTableRepository) ListByStatus(status JobStatus) ([]CompactTableModel, error) {
	return r.bucket.GetList(func(model CompactTableModel) bool {
		return model.Status == status
	}, r.sortDescByCreatedAt)
}

func (r *CompactTableRepository) ListByCluster(cluster string) ([]CompactTableModel, error) {
	return r.bucket.GetList(func(model CompactTableModel) bool {
		return model.Cluster == cluster
	}, r.sortDescByCreatedAt)
}

func (r *CompactTableRepository) Get(uuid uuid.UUID) (CompactTableModel, error) {
	return r.bucket.Get(uuid.String())
}

func (r *CompactTableRepository) Create(credentialId uuid.UUID, cluster string, args []string) (*CompactTableModel, error) {
	jobUuid := uuid.New()
	compactTableModel := CompactTableModel{
		Uuid:         jobUuid,
		CredentialId: credentialId,
		Cluster:      cluster,
		Status:       PENDING,
		Command:      "pgcompacttable " + strings.Join(args, " "),
		CommandArgs:  args,
		LogsPath:     r.file.Create(jobUuid),
		CreatedAt:    time.Now().UnixNano(),
	}

	err := r.bucket.Update(jobUuid.String(), compactTableModel)
	if err != nil {
		return nil, err
	}
	return &compactTableModel, nil
}

func (r *CompactTableRepository) UpdateStatus(compactTable CompactTableModel, status JobStatus) error {
	tmp := compactTable
	tmp.Status = status
	return r.bucket.Update(tmp.Uuid.String(), tmp)
}

func (r *CompactTableRepository) Delete(uuid uuid.UUID) error {
	var err error
	err = r.file.Delete(uuid)
	err = r.bucket.Delete(uuid.String())
	return err
}

func (r *CompactTableRepository) DeleteAll() error {
	err := r.file.DeleteAll()
	err = r.bucket.DeleteAll()
	return err
}

func (r *CompactTableRepository) sortDescByCreatedAt(list []CompactTableModel, i, j int) bool {
	return list[i].CreatedAt > list[j].CreatedAt
}
