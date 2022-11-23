package persistence

import (
	"github.com/google/uuid"
	"os"
)

var File *file

type file struct {
	CompactTable *FilePersistence
	Certs        *FilePersistence
}

func (f *file) Build() {
	File = &file{
		CompactTable: &FilePersistence{
			path:        "data/pgcompacttable",
			fileNameExt: ".log",
		},
		Certs: &FilePersistence{
			path:        "data/cert",
			fileNameExt: ".crt",
		},
	}

	_ = os.MkdirAll(File.CompactTable.path, os.ModePerm)
	_ = os.MkdirAll(File.Certs.path, os.ModePerm)
}

type FilePersistence struct {
	path        string
	fileNameExt string
}

func (r FilePersistence) Create(uuid uuid.UUID) string {
	path := r.path + "/" + uuid.String() + r.fileNameExt
	_, _ = os.Create(path)
	return path
}

func (r FilePersistence) Open(path string) (*os.File, error) {
	return os.OpenFile(path, os.O_RDWR, 0666)
}

func (r FilePersistence) Delete(uuid uuid.UUID) error {
	return os.Remove(r.path + "/" + uuid.String() + r.fileNameExt)
}
