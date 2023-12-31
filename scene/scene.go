package scene

import (
	"io/fs"

	"github.com/hndada/gosu/draws"
	"github.com/hndada/gosu/format/osr"
)

type Scene interface {
	Update() any
	Draw(screen draws.Image)
	// DebugString() string
}

type PlayArgs struct {
	MusicFS       fs.FS
	ChartFilename string
	Replay        *osr.Format
}
