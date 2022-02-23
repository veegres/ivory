package persistence

import (
	"bytes"
	"encoding/gob"
	"github.com/google/uuid"
	. "ivory/model"
	"os"
	"sort"
)

type CompactTableRepository struct {
	common common
	bucket []byte
}

func (r CompactTableRepository) List() ([]CompactTableModel, error) {
	bytesList, err := r.common.getList(r.bucket)
	modelList := make([]CompactTableModel, len(bytesList))
	for i, el := range bytesList {
		var model CompactTableModel
		buff := bytes.NewBuffer(el.value)
		_ = gob.NewDecoder(buff).Decode(&model)
		modelList[i] = model
	}
	sort.Slice(modelList, func(i, j int) bool {
		return modelList[i].CreatedAt > modelList[j].CreatedAt
	})
	return modelList, err
}

func (r CompactTableRepository) ListByStatus(status JobStatus) ([]CompactTableModel, error) {
	bytesList, err := r.common.getList(r.bucket)
	modelList := make([]CompactTableModel, 0)
	for _, el := range bytesList {
		var model CompactTableModel
		buff := bytes.NewBuffer(el.value)
		_ = gob.NewDecoder(buff).Decode(&model)
		if model.Status == status {
			modelList = append(modelList, model)
		}
	}
	sort.Slice(modelList, func(i, j int) bool {
		return modelList[i].CreatedAt > modelList[j].CreatedAt
	})
	return modelList, err
}

func (r CompactTableRepository) Get(uuid uuid.UUID) (CompactTableModel, error) {
	value, err := r.common.get(r.bucket, uuid.String())
	var model CompactTableModel
	buff := bytes.NewBuffer(value)
	_ = gob.NewDecoder(buff).Decode(&model)
	return model, err
}

func (r CompactTableRepository) Update(compactTable CompactTableModel) error {
	return r.common.update(r.bucket, compactTable.Uuid.String(), compactTable)
}

func (r CompactTableRepository) UpdateStatus(compactTable CompactTableModel, status JobStatus) error {
	tmp := compactTable
	tmp.Status = status
	return r.common.update(r.bucket, tmp.Uuid.String(), tmp)
}

func (r CompactTableRepository) Delete(uuid uuid.UUID) error {
	return r.common.delete(r.bucket, uuid.String())
}

type CompactTableFile struct {
	dir string
}

func (r CompactTableFile) Create(uuid uuid.UUID) string {
	path := r.dir + "/" + uuid.String() + ".log"
	_, _ = os.Create(path)
	return path
}

func (r CompactTableFile) Open(path string) (*os.File, error) {
	return os.OpenFile(path, os.O_RDWR, 0666)
}

func (r CompactTableFile) Delete(uuid uuid.UUID) error {
	return os.Remove(r.dir + "/" + uuid.String() + ".log")
}
