package format

import "github.com/hndada/gosu/framework/format"

type File = format.File

const (
	Unknown format.Type = iota - 1
	OsuBeatmap
	OsuReplay
)

func FileType(ext string) format.Type {
	switch ext {
	case ".osu":
		return OsuBeatmap
	case ".osr":
		return OsuReplay
	}
	return Unknown
}
