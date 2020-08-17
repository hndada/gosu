package osu

import (
	"strconv"
	"strings"
)

// time,beatLength,meter,sampleSet,sampleIndex,volume,uninherited,effects
func newTimingPoint(line string) (TimingPoint, error) {
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