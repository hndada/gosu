package gosu

import (
	"sort"

	"github.com/hndada/gosu/parse/osu"
)

var DefaultSpeedFactor = SpeedFactorPoint{Time: 0, Factor: 1}

type TransPoints struct {
	SpeedFactors []*SpeedFactorPoint
	Tempos       []*TempoPoint
	Volumes      []*VolumePoint
	Effects      []*EffectPoint
}
type TransPoint struct {
	SpeedFactor *SpeedFactorPoint
	Tempo       *TempoPoint
	Volume      *VolumePoint
	Effect      *EffectPoint
}
type SpeedFactorPoint struct {
	Time   int64
	Factor float64
	Next   *SpeedFactorPoint
	Prev   *SpeedFactorPoint
}
type TempoPoint struct {
	Time  int64
	BPM   float64
	Meter uint8
	Next  *TempoPoint
	Prev  *TempoPoint
}
type VolumePoint struct {
	Time   int64
	Volume uint8
	Next   *VolumePoint
	Prev   *VolumePoint
}
type EffectPoint struct {
	Time      int64
	Highlight bool
	Next      *EffectPoint
	Prev      *EffectPoint
}

// NewTransPointsFromOsu appends to a slice only when following element is different.
// Should be run at whole slice at once since those are all together at one section in osu format.
func NewTransPointsFromOsu(o *osu.Format) TransPoints {
	sort.Slice(o.TimingPoints, func(i int, j int) bool {
		return o.TimingPoints[i].Time < o.TimingPoints[j].Time
	})
	var (
		lastSpeedFactor float64 = 1
		lastBPM         float64
		lastMeter       uint8
		lastVolume      uint8
		lastHighlight   bool

		pTempo       *TempoPoint
		pSpeedFactor *SpeedFactorPoint
		pVolume      *VolumePoint
		pEffect      *EffectPoint
	)
	var tps TransPoints
	for _, tp := range o.TimingPoints {
		time := int64(tp.Time)
		if tp.Uninherited { // Uninherited: BPM, meter
			bpm, _ := tp.BPM()
			m := uint8(tp.Meter)
			if bpm != lastBPM || m != lastMeter {
				p := &TempoPoint{time, bpm, m, nil, pTempo}
				tps.Tempos = append(tps.Tempos, p)
				p.Prev.Next = p
				pTempo = p
			}
			lastBPM = bpm
			lastMeter = m
		} else {
			sf, _ := tp.SpeedFactor()
			if sf != lastSpeedFactor {
				p := &SpeedFactorPoint{time, sf, nil, pSpeedFactor}
				tps.SpeedFactors = append(tps.SpeedFactors, p)
				p.Prev.Next = p
				pSpeedFactor = p
			}
			lastSpeedFactor = sf
		}

		v := uint8(tp.Volume)
		if v != lastVolume {
			p := &VolumePoint{time, v, nil, pVolume}
			tps.Volumes = append(tps.Volumes, p)
			p.Prev.Next = p
			pVolume = p
		}
		lastVolume = v

		h := tp.IsKiai()
		if h != lastHighlight {
			p := &EffectPoint{time, h, nil, pEffect}
			tps.Effects = append(tps.Effects, p)
			p.Prev.Next = p
			pEffect = p
		}
		lastHighlight = h
	}
	return tps
}
