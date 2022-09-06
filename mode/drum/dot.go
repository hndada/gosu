package drum

// Time is a point of time, duartion a length of time.
// No consider situations that multiple rolls are overlapped.
type Dot struct {
	Floater
	// Showtime int64 // Dot will appear at Showtime.
	Marked bool
	Next   *Dot
	Prev   *Dot
}

// Unit of speed is osupixel / 100ms.
// n.SetDots(tp.Speed*speedFactor, bpm)
func NewDots(notes []*Note) (ds []*Dot) {
	for _, n := range notes {
		if n.Type != Head {
			continue
		}
		step := float64(n.Duration) / float64(n.Tick)
		for t := 0.0; t < float64(n.Duration); t += step {
			d := Dot{
				Floater: Floater{
					Time:  n.Time + int64(t),
					Speed: n.Speed,
				},
				// Showtime: Showtime,
			}
			ds = append(ds, &d)
		}
	}

	var prev *Dot
	for _, d := range ds {
		d.Prev = prev
		if prev != nil {
			prev.Next = d
		}
		prev = d
	}
	return
}
