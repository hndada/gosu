package gosu

import (
	"archive/zip"
	"io/fs"
)

// Todo: should .zip be extracted throughly?
func ZipFS(name string) fs.FS {
	r, err := zip.OpenReader(name)
	if err != nil {
		panic(err)
	}
	return r
}
