package persistence

import (
	"os"
)

var File file

type file struct {
	CompactTable *CompactTableFile
}

func (f file) Build() {
	compactTableDir := "data/pgcompacttable"
	_ = os.MkdirAll(compactTableDir, os.ModePerm)

	f.CompactTable = &CompactTableFile{
		dir: compactTableDir,
	}
	File = f
}
