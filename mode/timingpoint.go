package mode

import (
	"github.com/hndada/rg-parser/osugame/osu"
)

type TimingPoints struct {
	SpeedFactors []SpeedFactorPoint
	Tempos       []TempoPoint
	Volumes      []VolumePoint
	Effects      []EffectPoint
}

type SpeedFactorPoint struct {
	Time   int64
	Factor float64
}

type TempoPoint struct {
	Time  int64
	BPM   float64
	Meter uint8
}

type VolumePoint struct {
	Time   int64
	Volume uint8
}

type EffectPoint struct {
	Time      int64
	Highlight bool
}

// uninherited: BPM
// should be run at whole slice at once
func newTimingPointsFromOsu(o osu.Format) []TimingPoint {
	var lastBPM float64
	var lastSpeedFactor float64 = 1
	ts := make([]TimingPoint, len(o.TimingPoints))
	for i, tp := range o.TimingPoints {
		var t TimingPoint
		t.Time = int64(tp.Time)
		if tp.Uninherited {
			lastBPM = t.BPM
			t.BPM, _ = tp.BPM()
			t.SpeedFactor = lastSpeedFactor

		} else {
			lastSpeedFactor = t.SpeedFactor
			t.BPM = lastBPM
			t.SpeedFactor, _ = tp.SpeedFactor()
		}
		t.Meter = uint8(tp.Meter)
		t.Volume = uint8(tp.Volume)
		t.Highlight = tp.IsKiai()
		ts[i] = t
	}
	return ts
}
