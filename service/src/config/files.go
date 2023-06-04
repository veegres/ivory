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

// Create TODO we should handle error file creation not an a panic way
func (r *FileGateway) Create(name string) string {
	path := r.getPath(name)
	_, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	return path
}

func (r *FileGateway) Exist(name string) bool {
	if _, err := os.Stat(r.getPath(name)); err == nil {
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
	return os.Remove(r.getPath(name))
}

func (r *FileGateway) DeleteAll() error {
	errRem := os.RemoveAll(r.path + "/")
	errCreate := os.Mkdir(r.path, os.ModePerm)
	return errors.Join(errRem, errCreate)
}

// TODO throw error if name is empty
func (r *FileGateway) getPath(name string) string {
	return r.path + "/" + name + r.ext
}
