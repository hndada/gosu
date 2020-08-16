package osu

import (
	"strconv"
	"strings"

	"github.com/hndada/gosu/internal/tools"
)

// time,beatLength,meter,sampleSet,sampleIndex,volume,uninherited,effects
type TimingPoint struct {
	Time            int
	Bpm, SpeedScale float64
	Meter           int
	SampleSet       int
	SampleIndex     int
	Volume          int
	Uninherited     bool
	Effects         int
	// Kiai            bool
}
func newTimingPoint(line string) (TimingPoint, error) {
	var tp TimingPoint
	vs := strings.Split(line, `,`)

	time, err := tools.Atoi(vs[0])
	if err != nil {
		return tp, err
	}
	tp.Time = time

	beatLength, err := strconv.ParseFloat(vs[1], 64)
	if err != nil {
		return tp, err
	}

	meter, err := tools.Atoi(vs[2])
	if err != nil {
		return tp, err
	}
	tp.Meter = meter

	sampleSet, err := tools.Atoi(vs[3])
	if err != nil {
		return tp, err
	}
	tp.SampleSet = sampleSet

	sampleIndex, err := tools.Atoi(vs[4])
	if err != nil {
		return tp, err
	}
	tp.SampleIndex = sampleIndex

	volume, err := tools.Atoi(vs[5])
	if err != nil {
		return tp, err
	}
	tp.Volume = volume

	uninherited, err := strconv.ParseBool(vs[6])
	if err != nil {
		return tp, err
	}
	tp.Uninherited = uninherited

	effect, err := tools.Atoi(vs[7])
	if err != nil {
		return tp, err
	}
	tp.Effects = effect
	// tp.Kiai = effect&1 != 0

	switch tp.Uninherited {
	case true:
		tp.Bpm = 1000 * 60 / beatLength
	case false:
		tp.SpeedScale = 100 / (-beatLength)
	}
	return tp, nil
}
