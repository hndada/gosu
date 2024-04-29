package scene

import (
	"io/fs"

	draws "github.com/hndada/gosu/draws5"
	"github.com/hndada/gosu/format/osr"
)

// const (
// 	SceneSelect = iota
// 	ScenePlay
// )

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
