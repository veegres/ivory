package files

import (
	"errors"
	"os"
	"strings"
)

type Storage struct {
	path string
	ext  string
}

func NewStorage(name string, ext string) *Storage {
	path := "data/" + name

	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		panic(err)
	}

	return &Storage{
		path: path,
		ext:  ext,
	}
}

func (r *Storage) CreateByName(name string) (string, error) {
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

func (r *Storage) ExistByName(name string) bool {
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

func (r *Storage) OpenByName(name string) (*os.File, error) {
	path, errPath := r.getPath(name)
	if errPath != nil {
		return nil, errPath
	}
	return os.OpenFile(path, os.O_RDWR|os.O_APPEND, 0666)
}

func (r *Storage) ReadByName(name string) ([]byte, error) {
	path, errPath := r.getPath(name)
	if errPath != nil {
		return nil, errPath
	}
	return os.ReadFile(path)
}

func (r *Storage) DeleteByName(name string) error {
	path, err := r.getPath(name)
	if err != nil {
		return err
	}
	return os.Remove(path)
}

func (r *Storage) ExistByPath(path string) bool {
	if _, err := os.Stat(path); err == nil {
		return true
	} else {
		return false
	}
}

func (r *Storage) DeleteAll() error {
	errRem := os.RemoveAll(r.path + "/")
	errCreate := os.Mkdir(r.path, os.ModePerm)
	return errors.Join(errRem, errCreate)
}

func (r *Storage) getPath(name string) (string, error) {
	if name == "" {
		return "", errors.New("file name cannot be empty")
	}
	if strings.ContainsAny(name, "./") {
		return "", errors.New("file name contains invalid characters '/', '.'")
	}
	return r.path + "/" + name + r.ext, nil
}
