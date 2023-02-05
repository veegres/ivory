package persistence

import (
	"bytes"
	"encoding/gob"
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
	value, err := r.common.get(r.bucket, uuid.String())
	var model CompactTableModel
	buff := bytes.NewBuffer(value)
	_ = gob.NewDecoder(buff).Decode(&model)
	return model, err
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

	err := r.common.update(r.bucket, jobUuid.String(), compactTableModel)
	if err != nil {
		return nil, err
	}
	return &compactTableModel, nil
}

func (r CompactTableRepository) UpdateStatus(compactTable CompactTableModel, status JobStatus) error {
	tmp := compactTable
	tmp.Status = status
	return r.common.update(r.bucket, tmp.Uuid.String(), tmp)
}

func (r CompactTableRepository) Delete(uuid uuid.UUID) error {
	var err error
	err = File.CompactTable.Delete(uuid)
	err = r.common.delete(r.bucket, uuid.String())
	return err
}

func (r CompactTableRepository) DeleteAll() error {
	err := File.CompactTable.DeleteAll()
	err = r.common.deleteAll(r.bucket)
	return err
}

func (r CompactTableRepository) list(filter func(model CompactTableModel) bool) ([]CompactTableModel, error) {
	bytesList, err := r.common.getList(r.bucket)
	modelList := make([]CompactTableModel, 0)
	for _, el := range bytesList {
		var model CompactTableModel
		buff := bytes.NewBuffer(el.value)
		_ = gob.NewDecoder(buff).Decode(&model)
		if filter == nil || filter(model) {
			modelList = append(modelList, model)
		}
	}
	r.sortDescByCreatedAt(modelList)
	return modelList, err
}

func (r CompactTableRepository) sortDescByCreatedAt(list []CompactTableModel) {
	sort.Slice(list, func(i, j int) bool {
		return list[i].CreatedAt > list[j].CreatedAt
	})
}
