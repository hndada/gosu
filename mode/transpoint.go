package mode

import (
	"github.com/hndada/rg-parser/osugame/osu"
)

// uninherited: BPM
type TransPoint struct {
	Time      int64
	BPM       float64
	Speed     float64
	Meter     uint8
	Volume    uint8
	Highlight bool
}

// should be run at whole slice at once
func newTransPointsFromOsu(o osu.Format) []TransPoint {
	var lastBPM float64
	var lastSpeed float64 = 1
	ts := make([]TransPoint, len(o.TimingPoints))
	for i, tp := range o.TimingPoints {
		var t TransPoint
		t.Time = int64(tp.Time)
		if tp.Uninherited {
			lastBPM = t.BPM
			t.BPM, _ = tp.BPM()
			t.Speed = lastSpeed

		} else {
			lastSpeed = t.Speed
			t.BPM = lastBPM
			t.Speed, _ = tp.Speed()
		}
		t.Meter = uint8(tp.Meter)
		t.Volume = uint8(tp.Volume)
		t.Highlight = tp.IsKiai()
		ts[i] = t
	}
	return ts
}
