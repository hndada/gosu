package choose

import (
	"archive/zip"
	"io/fs"
)

// Todo: should .zip be extracted throughly?
func ZipFS(name string) (fs.FS, error) {
	r, err := zip.OpenReader(name)
	if err != nil {
		return nil, err
	}
	return r, nil
}
