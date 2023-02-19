package config

import (
	"github.com/google/uuid"
	"os"
)

type FilePersistence struct {
	path string
	ext  string
}

func NewFilePersistence(name string, ext string) *FilePersistence {
	path := "data/" + name

	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		panic(err)
	}

	return &FilePersistence{
		path: path,
		ext:  ext,
	}
}

func (r *FilePersistence) Create(uuid uuid.UUID) string {
	path := r.path + "/" + uuid.String() + r.ext
	_, _ = os.Create(path)
	return path
}

func (r *FilePersistence) Open(path string) (*os.File, error) {
	return os.OpenFile(path, os.O_RDWR, 0666)
}

func (r *FilePersistence) Read(path string) ([]byte, error) {
	return os.ReadFile(path)
}

func (r *FilePersistence) Delete(uuid uuid.UUID) error {
	return os.Remove(r.path + "/" + uuid.String() + r.ext)
}

func (r *FilePersistence) DeleteAll() error {
	err := os.RemoveAll(r.path + "/")
	err = os.Mkdir(r.path, os.ModePerm)
	return err
}
