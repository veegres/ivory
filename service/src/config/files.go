package config

import (
	"errors"
	"os"
	"strings"
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

func (r *FileGateway) CreateByName(name string) (string, error) {
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

func (r *FileGateway) ExistByName(name string) bool {
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

func (r *FileGateway) OpenByName(name string) (*os.File, error) {
	path, errPath := r.getPath(name)
	if errPath != nil {
		return nil, errPath
	}
	return os.OpenFile(path, os.O_RDWR|os.O_APPEND, 0666)
}

func (r *FileGateway) ReadByName(name string) ([]byte, error) {
	path, errPath := r.getPath(name)
	if errPath != nil {
		return nil, errPath
	}
	return os.ReadFile(path)
}

func (r *FileGateway) DeleteByName(name string) error {
	path, err := r.getPath(name)
	if err != nil {
		return err
	}
	return os.Remove(path)
}

func (r *FileGateway) ExistByPath(path string) bool {
	if _, err := os.Stat(path); err == nil {
		return true
	} else {
		return false
	}
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
	if strings.ContainsAny(name, "./") {
		return "", errors.New("file name contains invalid characters '/', '.'")
	}
	return r.path + "/" + name + r.ext, nil
}
