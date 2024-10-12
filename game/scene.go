package scene

import (
	"io/fs"

	"github.com/hndada/gosu/draws"
	"github.com/hndada/gosu/game"
)

// const (
// 	SceneSelect = iota
// 	ScenePlay
// )

type Scene interface {
	Update() any
	Draw(screen draws.Image)
	DebugString() string
}

type PlayArgs struct {
	ChartFS        fs.FS // Music file exists in the same directory.
	ChartFilename  string
	Mods           game.Mods
	ReplayFS       fs.FS
	ReplayFilename string
}
