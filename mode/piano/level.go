package piano

import "github.com/hndada/gosu/mode"

var fingersList = [][]int{
	{},
	{0},
	{1, 1},
	{1, 0, 1},
	{2, 1, 1, 2},
	{2, 1, 0, 1, 2},
	{3, 2, 1, 1, 2, 3},
	{3, 2, 1, 0, 1, 2, 3},
	{4, 3, 2, 1, 1, 2, 3, 4},
	{4, 3, 2, 1, 0, 1, 2, 3, 4},
	{4, 3, 2, 1, 0, 0, 1, 2, 3, 4},
}

func Fingers(keyCount int, scratchMode ScratchMode) []int {
	maxFinger := fingersList[keyCount-1][0] + 1
	switch scratchMode {
	case NoScratch:
		return fingersList[keyCount]
	case LeftScratch:
		return append([]int{maxFinger}, fingersList[keyCount-1]...)
	case RightScratch:
		return append(fingersList[keyCount-1], maxFinger)
	}
	return nil
}

// Weight is for Tail's variadic weight based on its length.
// For example, short long note does not require much strain to release.
// Todo: fine-tuning with replay data
func (n Note) Weight() float64 {
	switch n.Type {
	case Tail:
		head := n.Prev
		d := float64(head.Time + head.Duration)
		switch {
		case d < 50:
			return 0.5 - 0.35*d/50
		case d >= 50 && d < 200:
			return 0.15
		case d >= 200 && d < 800:
			return 0.15 + 0.85*(d-200)/600
		default:
			return 1
		}
	default:
		return 1
	}
}

func (c Chart) Difficulties() (ds []float64) {
	if len(c.Notes) == 0 {
		return
	}

	const standardDuration = 800 // 800ms. 2 beats with 150 BPM
	times, durations := mode.DifficultyPieceTimes(c.Dynamics, c.Duration())
	ds = make([]float64, 0, len(times))

	var i int
	var d float64
	for _, n := range c.Notes {
		for n.Time > times[i] {
			scale := standardDuration / float64(durations[i])
			ds = append(ds, d*scale)
			d = 0
			i++
		}
		d += n.Strain
	}
	return ds
}
