package mode

import (
	"io/fs"

	"github.com/hndada/gosu/format/osr"
)

type ScenePlayArgs struct {
	FS            fs.FS
	ChartFilename string
	Mods          interface{}
	Replay        *osr.Format
}
