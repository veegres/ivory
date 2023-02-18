package persistence

import (
	"github.com/google/uuid"
	. "ivory/model"
	"sort"
	"strings"
	"time"
)

type CompactTableRepository struct {
	common common
	bucket []byte
}

func (r CompactTableRepository) List() ([]CompactTableModel, error) {
	return r.list(nil)
}

func (r CompactTableRepository) ListByStatus(status JobStatus) ([]CompactTableModel, error) {
	return r.list(func(model CompactTableModel) bool {
		return model.Status == status
	})
}

func (r CompactTableRepository) ListByCluster(cluster string) ([]CompactTableModel, error) {
	return r.list(func(model CompactTableModel) bool {
		return model.Cluster == cluster
	})
}

func (r CompactTableRepository) Get(uuid uuid.UUID) (CompactTableModel, error) {
	return Get[CompactTableModel](r.bucket, uuid.String())
}

func (r CompactTableRepository) Create(credentialId uuid.UUID, cluster string, args []string) (*CompactTableModel, error) {
	jobUuid := uuid.New()
	compactTableModel := CompactTableModel{
		Uuid:         jobUuid,
		CredentialId: credentialId,
		Cluster:      cluster,
		Status:       PENDING,
		Command:      "pgcompacttable " + strings.Join(args, " "),
		CommandArgs:  args,
		LogsPath:     File.CompactTable.Create(jobUuid),
		CreatedAt:    time.Now().UnixNano(),
	}

	err := Update(r.bucket, jobUuid.String(), compactTableModel)
	if err != nil {
		return nil, err
	}
	return &compactTableModel, nil
}

func (r CompactTableRepository) UpdateStatus(compactTable CompactTableModel, status JobStatus) error {
	tmp := compactTable
	tmp.Status = status
	return Update(r.bucket, tmp.Uuid.String(), tmp)
}

func (r CompactTableRepository) Delete(uuid uuid.UUID) error {
	var err error
	err = File.CompactTable.Delete(uuid)
	err = Delete(r.bucket, uuid.String())
	return err
}

func (r CompactTableRepository) DeleteAll() error {
	err := File.CompactTable.DeleteAll()
	err = DeleteAll(r.bucket)
	return err
}

func (r CompactTableRepository) list(filter func(model CompactTableModel) bool) ([]CompactTableModel, error) {
	modelList, err := GetList[CompactTableModel](r.bucket, filter)
	r.sortDescByCreatedAt(modelList)
	return modelList, err
}

func (r CompactTableRepository) sortDescByCreatedAt(list []CompactTableModel) {
	sort.Slice(list, func(i, j int) bool {
		return list[i].CreatedAt > list[j].CreatedAt
	})
}
