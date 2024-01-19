package selects

import (
	"io/fs"

	"github.com/hndada/gosu/game"
)

type chartFolder struct {
	name  string
	items []chartItem
	index int
}

type chartItem struct {
	fsys fs.FS
	name string

	game.ChartHeader // Header is supposed to be loaded from database.
	// level            game.Level
	// bestResult       game.Result
}

func newChartItem(fsys fs.FS, name string) (chartItem, error) {
	f, err := fsys.Open(name)
	if err != nil {
		return chartItem{}, err
	}
	defer f.Close()

	format, hash, err := game.LoadChartFile(fsys, name)
	if err != nil {
		return chartItem{}, err
	}

	return chartItem{
		fsys:        fsys,
		name:        name,
		ChartHeader: game.NewChartHeader(format, hash),
	}, nil
}
