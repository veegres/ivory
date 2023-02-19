package persistence

import (
	"github.com/google/uuid"
	"ivory/src/config"
	. "ivory/src/model"
	"sort"
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
	return r.list(nil)
}

func (r *CompactTableRepository) ListByStatus(status JobStatus) ([]CompactTableModel, error) {
	return r.list(func(model CompactTableModel) bool {
		return model.Status == status
	})
}

func (r *CompactTableRepository) ListByCluster(cluster string) ([]CompactTableModel, error) {
	return r.list(func(model CompactTableModel) bool {
		return model.Cluster == cluster
	})
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

func (r *CompactTableRepository) list(filter func(model CompactTableModel) bool) ([]CompactTableModel, error) {
	modelList, err := r.bucket.GetList(filter)
	r.sortDescByCreatedAt(modelList)
	return modelList, err
}

func (r *CompactTableRepository) sortDescByCreatedAt(list []CompactTableModel) {
	sort.Slice(list, func(i, j int) bool {
		return list[i].CreatedAt > list[j].CreatedAt
	})
}
