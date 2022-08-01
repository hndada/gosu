package main

import (
	"github.com/hndada/gosu/parse/osu"
)

var DefaultSpeedFactor = SpeedFactorPoint{Time: 0, Factor: 1}

type TransPoints struct {
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

// Uninherited: BPM, meter
// Should be run at whole slice at once
func NewTransPointsFromOsu(o *osu.Format) TransPoints {
	var (
		lastSpeedFactor float64 = 1
		lastBPM         float64
		lastMeter       uint8
		lastVolume      uint8
		lastHighlight   bool
	)
	var tps TransPoints
	for _, tp := range o.TimingPoints {
		time := int64(tp.Time)
		// 다를 때에만 각 slice에 추가
		if tp.Uninherited {
			bpm, _ := tp.BPM()
			m := uint8(tp.Meter)
			if bpm != lastBPM || m != lastMeter {
				t := TempoPoint{time, bpm, m}
				tps.Tempos = append(tps.Tempos, t)
				lastBPM = bpm
				lastMeter = m
			}
		} else {
			sf, _ := tp.SpeedFactor()
			if sf != lastSpeedFactor {
				sfp := SpeedFactorPoint{time, sf}
				tps.SpeedFactors = append(tps.SpeedFactors, sfp)
				lastSpeedFactor = sf
			}
		}
		v := uint8(tp.Volume)
		if v != lastVolume {
			vp := VolumePoint{time, v}
			tps.Volumes = append(tps.Volumes, vp)
			lastVolume = v
		}
		h := tp.IsKiai()
		if h != lastHighlight {
			ep := EffectPoint{time, h}
			tps.Effects = append(tps.Effects, ep)
			lastHighlight = h
		}
	}
	return tps
}
