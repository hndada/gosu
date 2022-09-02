package gosu

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/hndada/gosu/format/osr"
)

// Todo: Make sure to ReplayListener time is independent of Game's update tick
// ReplayListener supposes closure function is called every 1 ms.
func NewReplayListener(f *osr.Format, keyCount int, waitBefore int64) func() []bool {
	actions := f.TrimmedActions()
	actions = append(actions, osr.Action{W: 2e9})
	var i int // Index of current replay action
	var t = waitBefore
	var next = 0 + 1 + actions[0].W + actions[1].W
	return func() []bool {
		if t >= next {
			i++
			next += actions[i+1].W
		}
		pressed := make([]bool, keyCount)
		var k int
		for x := int(actions[i].X); x > 0; x /= 2 {
			if x%2 == 1 {
				pressed[k] = true
			}
			k++
		}
		t++
		return pressed
	}
}

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
