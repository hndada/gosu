package gosu

import (
	"github.com/hndada/gosu/draws"
)

type Bar struct {
	sprite   draws.Sprite
	time     int64
	position float64
	speed    float64
	next     *Bar
	prev     *Bar
}

func (b *Bar) Sprite() draws.Sprite     { return b.sprite }
func (b *Bar) BodySprite() draws.Sprite { return draws.Sprite{} }
func (b *Bar) Position() float64        { return b.position }
func (b *Bar) SetPosition(pos float64)  { b.position = pos }
func (b *Bar) Speed() float64           { return b.speed }
func (b *Bar) IsHead() bool             { return false }
func (b *Bar) IsTail() bool             { return false }
func (b *Bar) Marked() bool             { return false }
func (b *Bar) Next() LaneSubject        { return b.next }
func (b *Bar) Prev() LaneSubject        { return b.prev }

func NewBars(transPoints []*TransPoint, endTime int64, sprite draws.Sprite) []LaneSubject {
	bars := make([]Bar, 0)
	first := transPoints[0]
	first = first.FetchLatest()
	var margin int64 = 5000
	if margin > first.Time {
		margin = first.Time
	}

	speed := first.BPM * first.BeatLengthScale
	for t := float64(first.Time); t >= float64(-margin); t -= float64(first.Meter) * 60000 / first.BPM {
		bar := Bar{
			sprite:   sprite,
			time:     int64(t),
			position: speed * (t - 0),
			speed:    speed,
		}
		bars = append([]Bar{bar}, bars...)
	}

	bars = bars[:len(bars)-1] // Drop for avoiding duplicated
	// var prevPos float64 = first.Position
	var nextPos float64 = first.Position
	for tp := first; tp != nil; tp = tp.NextBPMPoint.FetchLatest() {
		nextTime := endTime + margin
		if tp.NextBPMPoint != nil {
			nextTime = tp.NextBPMPoint.Time
		}
		speed := tp.BPM * tp.BeatLengthScale
		unit := float64(tp.Meter) * 60000 / tp.BPM
		for t := float64(tp.Time); t < float64(nextTime); t += unit {
			// pos := prevPos + speed*unit
			pos := nextPos
			bar := Bar{
				sprite:   sprite,
				time:     int64(t),
				position: pos,
				speed:    speed,
			}
			bars = append(bars, bar)
			// prevPos = pos
			nextPos = pos + speed*unit
		}
	}
	for i := range bars {
		switch i {
		case 0:
			bars[i].next = &bars[i+1]
		case len(bars) - 1:
			bars[i].prev = &bars[i-1]
		default:
			bars[i].next = &bars[i+1]
			bars[i].prev = &bars[i-1]
		}
	}
	bars2 := make([]LaneSubject, len(bars))
	for i := range bars2 {
		bars2[i] = &bars[i]
	}
	return bars2
}

// func BarPositions(transPoints []*TransPoint, endTime int64) []float64 {
// 	ps := make([]float64, 0)
// 	first := transPoints[0]
// 	first = first.FetchLatest()
// 	var margin int64 = 5000
// 	if margin > first.Time {
// 		margin = first.Time
// 	}
// 	speed := first.BPM * first.BeatLengthScale
// 	for t := float64(first.Time); t >= float64(-margin); t -= float64(first.Meter) * 60000 / first.BPM {
// 		p := speed * (t - 0)
// 		ps = append([]float64{p}, ps...)
// 	}

// 	ps = ps[:len(ps)-1] // Drop for avoiding duplicated
// 	for tp := first; tp != nil; tp = tp.NextBPMPoint.FetchLatest() {
// 		nextTime := endTime + margin
// 		if tp.NextBPMPoint != nil {
// 			nextTime = tp.NextBPMPoint.Time
// 		}
// 		speed := tp.BPM * tp.BeatLengthScale
// 		unit := float64(tp.Meter) * 60000 / tp.BPM
// 		for t := float64(tp.Time); t < float64(nextTime); t += unit {
// 			p := ps[len(ps)-1] + speed*unit
// 			ps = append(ps, p)
// 		}
// 	}
// 	return ps
// }

// func BarTimes(transPoints []*TransPoint, endTime int64) []int64 {
// 	ts := make([]int64, 0)
// 	first := transPoints[0]
// 	first.FetchLatest()
// 	var margin int64 = 5000
// 	if margin > first.Time {
// 		margin = first.Time
// 	}
// 	for t := float64(first.Time); t >= float64(-margin); t -= float64(first.Meter) * 60000 / first.BPM {
// 		ts = append([]int64{int64(t)}, ts...)
// 	}
// 	// Bar for the first TransPoint has already appended.
// 	for tp := first.NextBPMPoint; tp != nil; tp = tp.NextBPMPoint {
// 		next := endTime + margin
// 		if tp.NextBPMPoint != nil {
// 			next = tp.NextBPMPoint.Time
// 		}
// 		unit := float64(tp.Meter) * 60000 / tp.BPM
// 		for t := float64(tp.Time); t < float64(next); t += unit {
// 			ts = append(ts, int64(t))
// 		}
// 	}
// 	return ts
// }
