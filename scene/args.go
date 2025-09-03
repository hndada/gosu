package scene

import (
	"io/fs"
	"os"

	"github.com/hndada/gosu/game"
)

type Args interface{}

type PlayArgsData struct {
	ChartFS       string `json:"chartFS"`
	ChartFilename string `json:"chartFilename"`
}

func (ad PlayArgsData) ToPlayArgs() PlayArgs {
	return PlayArgs{
		ChartFS:       os.DirFS(ad.ChartFS),
		ChartFilename: ad.ChartFilename,
	}
}

type PlayArgs struct {
	// Music file exists in the same directory.
	ChartFS        fs.FS     `json:"chartFS"`
	ChartFilename  string    `json:"chartFilename"`
	Mods           game.Mods `json:"mods"`
	ReplayFS       fs.FS     `json:"replayFS"`
	ReplayFilename string    `json:"replayFilename"`
}
