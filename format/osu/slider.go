package osu

import (
	"fmt"
	"strings"
)

type Point = [2]int // delimiter:

type SliderParams struct { // delimiter,
	CurveType   string  // one letter
	CurvePoints []Point // delimiter|
	Slides      int
	Length      float64
	EdgeSounds  []int   // delimiter|
	EdgeSets    []Point // delimiter|
}

func newSliderParams(data string) (sp SliderParams, err error) {
	// curveType|curvePoints,slides,length,edgeSounds,edgeSets
	vs := strings.Split(data, `,`)
	if len(vs) < 3 {
		return sp, fmt.Errorf("slider params has no enough length: %s", data)
	}

	// curveType|curvePoints
	// B|200:200|250:200
	curveData := strings.Split(vs[0], `|`)
	sp.CurveType = curveData[0]
	sp.CurvePoints = make([][2]int, len(curveData)-1)
	for i, p := range curveData[1:] {
		if sp.CurvePoints[i], err = parsePoint(p); err != nil {
			return
		}
	}
	// slides
	if sp.Slides, err = parseInt(vs[1]); err != nil {
		return
	}
	// length
	if sp.Length, err = parseFloat(vs[2]); err != nil {
		return
	}

	if len(vs) < 5 {
		return
	}

	// edgeSounds
	// 2|1|2
	edgeSounds := strings.Split(vs[3], `|`)
	sp.EdgeSounds = make([]int, len(edgeSounds))
	for i := 0; i < len(edgeSounds); i++ {
		if sp.EdgeSounds[i], err = parseInt(edgeSounds[i]); err != nil {
			return
		}
	}
	// edgeSets
	// 0:0|0:0|0:2
	edgeSets := strings.Split(vs[4], `|`)
	sp.EdgeSets = make([][2]int, len(edgeSets))
	for i, p := range edgeSets {
		if sp.EdgeSets[i], err = parsePoint(p); err != nil {
			return
		}
	}

	return
}

// SliderDuration returns duration of slider in milliseconds.
func (h HitObject) SliderDuration(speed float64) int {
	// speed := (bpm / 60000) * beatScale * (multiplier * 100)
	return int(h.SliderLength() / speed)
}

// If hit object is not slider, both count and unit will be zero.
func (h HitObject) SliderLength() float64 {
	count := float64(h.SliderParams.Slides)
	unit := h.SliderParams.Length
	return count * unit
}

func parsePoint(s string) (p [2]int, err error) {
	xy := strings.Split(s, `:`)
	if len(xy) < 2 {
		return p, fmt.Errorf("point has not enough length: %s", s)
	}
	for i := 0; i < 2; i++ {
		if p[i], err = parseInt(xy[i]); err != nil {
			return
		}
	}
	return
}
