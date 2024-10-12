package game

import (
	"io/fs"

	"github.com/hndada/gosu/plays"
)

type Args interface{}

type PlayArgs struct {
	ChartFS        fs.FS // Music file exists in the same directory.
	ChartFilename  string
	Mods           plays.Mods
	ReplayFS       fs.FS
	ReplayFilename string
}
