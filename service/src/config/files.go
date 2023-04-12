package config

import (
	"errors"
	"github.com/google/uuid"
	"os"
)

type FileGateway struct {
	path string
	ext  string
}

func NewFileGateway(name string, ext string) *FileGateway {
	path := "data/" + name

	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		panic(err)
	}

	return &FileGateway{
		path: path,
		ext:  ext,
	}
}

func (r *FileGateway) Create(uuid uuid.UUID) string {
	path := r.path + "/" + uuid.String() + r.ext
	_, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	return path
}

func (r *FileGateway) Open(path string) (*os.File, error) {
	return os.OpenFile(path, os.O_RDWR, 0666)
}

func (r *FileGateway) Read(path string) ([]byte, error) {
	return os.ReadFile(path)
}

func (r *FileGateway) Delete(uuid uuid.UUID) error {
	return os.Remove(r.path + "/" + uuid.String() + r.ext)
}

func (r *FileGateway) DeleteAll() error {
	errRem := os.RemoveAll(r.path + "/")
	errCreate := os.Mkdir(r.path, os.ModePerm)
	return errors.Join(errRem, errCreate)
}
