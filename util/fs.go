package util

import (
	"io"
	"io/fs"
)

func ReadFile(fsys fs.FS, name string) ([]byte, error) {
	f, err := fsys.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return io.ReadAll(f)
}
