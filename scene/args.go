package scene

import (
	"io/fs"

	"github.com/hndada/gosu/format/osr"
)

type PlayArgs struct {
	MusicFS   fs.FS
	ChartName string
	Replay    *osr.Format
}
