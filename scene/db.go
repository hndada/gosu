package scene

import "io/fs"

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
