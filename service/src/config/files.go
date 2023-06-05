package config

import (
	"errors"
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

func (r *FileGateway) Create(name string) (string, error) {
	path, err := r.getPath(name)
	if err != nil {
		return "", err
	}
	_, errCreate := os.Create(path)
	if errCreate != nil {
		return "", errCreate
	}
	return path, nil
}

func (r *FileGateway) Exist(name string) bool {
	path, errPath := r.getPath(name)
	if errPath != nil {
		return false
	}
	if _, err := os.Stat(path); err == nil {
		return true
	} else {
		return false
	}
}

// Open TODO change it to use name of the file
func (r *FileGateway) Open(path string) (*os.File, error) {
	return os.OpenFile(path, os.O_RDWR, 0666)
}

// TODO change it to use name of the file
func (r *FileGateway) Read(path string) ([]byte, error) {
	return os.ReadFile(path)
}

func (r *FileGateway) Delete(name string) error {
	path, err := r.getPath(name)
	if err != nil {
		return err
	}
	return os.Remove(path)
}

func (r *FileGateway) DeleteAll() error {
	errRem := os.RemoveAll(r.path + "/")
	errCreate := os.Mkdir(r.path, os.ModePerm)
	return errors.Join(errRem, errCreate)
}

func (r *FileGateway) getPath(name string) (string, error) {
	if name == "" {
		return "", errors.New("file name cannot be empty")
	}
	return r.path + "/" + name + r.ext, nil
}
