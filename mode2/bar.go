package mode

type Bar struct {
	Time     int32 // Times are in milliseconds.
	Position float64
	Next     *Bar
	Prev     *Bar
}
type Bars []*Bar

func (ds Dynamics) NewBars(chartDuration int32) (bs Bars) {
	// const useDefaultMeter = 0
	times := ds.BeatTimes(chartDuration)
	bs = make(Bars, 0, len(times))
	for _, t := range times {
		b := Bar{Time: t}
		bs = append(bs, &b)
	}

	// linking
	var prev *Bar
	for _, b := range bs {
		b.Prev = prev
		if prev != nil {
			prev.Next = b
		}
		prev = b
	}
	return
}
