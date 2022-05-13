package persistence

import (
	"os"
)

var File *file

type file struct {
	CompactTable *CompactTableFile
}

func (f *file) Build() {
	compactTableDir := "data/pgcompacttable"
	_ = os.MkdirAll(compactTableDir, os.ModePerm)

	File = &file{
		CompactTable: &CompactTableFile{
			dir: compactTableDir,
		},
	}
}
