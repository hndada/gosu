package osu

import (
	"strconv"
	"strings"
)

func newTimingPoint(line string) (TimingPoint, error) {
	// time,beatLength,meter,sampleSet,sampleIndex,volume,uninherited,effects
	var tp TimingPoint
	vs := strings.Split(line, `,`)
	{
		f, err := strconv.ParseFloat(vs[0], 64)
		if err != nil {
			return tp, err
		}
		tp.Time = int(f)
	}
	{
		f, err := strconv.ParseFloat(vs[1], 64)
		if err != nil {
			return tp, err
		}
		tp.BeatLength = f
	}
	{
		f, err := strconv.ParseFloat(vs[2], 64)
		if err != nil {
			return tp, err
		}
		tp.Meter = int(f)
	}
	{
		i, err := strconv.Atoi(vs[3])
		if err != nil {
			return tp, err
		}
		tp.SampleSet = i
	}
	{
		i, err := strconv.Atoi(vs[4])
		if err != nil {
			return tp, err
		}
		tp.SampleIndex = i
	}
	{
		f, err := strconv.ParseFloat(vs[5], 64)
		if err != nil {
			return tp, err
		}
		tp.Volume = int(f)
	}
	{
		b, err := strconv.ParseBool(vs[6])
		if err != nil {
			return tp, err
		}
		tp.Uninherited = b
	}
	{
		i, err := strconv.Atoi(vs[7])
		if err != nil {
			return tp, err
		}
		tp.Effects = i
	}
	return tp, nil
}

func (tp TimingPoint) IsInherited() bool { return !tp.Uninherited }

func (tp TimingPoint) BPM() (bpm float64, ok bool) {
	if !tp.Uninherited {
		return 0, false
	}
	return 1000 * 60 / tp.BeatLength, true
}
// Speed returns speed scale. The standard speed value is 1.
func (tp TimingPoint) Speed() (speed float64, ok bool) {
	if tp.Uninherited {
		return 0, false
	}
	return 100 / (-tp.BeatLength), true
}

func (tp TimingPoint) isKiai() bool { return tp.Effects&1 != 0 }
func (tp TimingPoint) isFirstBarOmitted() bool {
	return tp.Effects&(1<<3) != 0
}
