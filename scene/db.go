package scene

import "io/fs"

// A function which load database should not load the entire file system into memory.
// Instead, they should load the minimum information only, such as
// path to the file, and read the file when needed.

type Databases struct {
	Music  []MusicRow
	Chart  map[FileKey][]ChartRow
	Replay []ReplayRow
}

type FileKey struct {
	FS   fs.FS
	Name string
}

type MusicRow struct {
	FileKey
	// FolderName string
}

// Pseudo database
type ChartRow struct {
	FileKey
	MusicName string
	Artist    string
	ChartName string
	Mode      int
	SubMode   int
	ChartHash string
	Level     float64
}

type ReplayRow struct {
	FileKey
	ChartHash string
}

// Memo: archive/zip.OpenReader returns ReadSeeker, which implements Read.
// Both Read and fs.Open are same in type: (name string) (fs.File, error)
// func zipFS(path string) (fs.FS, error) {
// 	r, err := zip.OpenReader(path)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return r, nil
// }
