package mode

import (
	"io/fs"

	"github.com/hndada/gosu/format/osr"
)

type ScenePlayArgs struct {
	FS            fs.FS
	ChartFilename string
	Mods          any
	Replay        *osr.Format
}
