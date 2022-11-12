package game

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/hndada/gosu/format/osr"
)

// Todo: Make own ScenePlay for calculating score from input replay file
// Todo: implement non-playing score simulator
// NewScenePlayCalc(Chart, Mods, *osr.Format); Update returns PlayToResultArgs {} if finished.
func LoadReplays(replayRoot string) ([]*osr.Format, error) {
	rfs := make([]*osr.Format, 0)
	fs, err := os.ReadDir(replayRoot)
	if err != nil {
		return nil, err
	}
	for _, f := range fs {
		if f.IsDir() || filepath.Ext(f.Name()) != ".osr" {
			continue
		}
		rd, err := os.ReadFile(filepath.Join(replayRoot, f.Name()))
		if err != nil {
			fmt.Println(err)
			continue
		}
		rf, err := osr.Parse(rd)
		if err != nil {
			fmt.Println(err)
			continue
		}
		rfs = append(rfs, rf)
	}
	return rfs, nil
}
