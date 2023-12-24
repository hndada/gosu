package mode

// Defining Bars is just redundant if it has no additional methods.
type Bar struct {
	Time     int32 // Times are in milliseconds.
	Position float64
}

func (ds Dynamics) NewBars(chartDuration int32) (bs []Bar) {
	// const useDefaultMeter = 0
	times := ds.BeatTimes(chartDuration)
	bs = make([]Bar, 0, len(times))
	for _, t := range times {
		b := Bar{Time: t}
		bs = append(bs, b)
	}
	return
}

// Drum and Piano modes have different bar drawing methods.
// Hence, this method is not defined in mode.go.
