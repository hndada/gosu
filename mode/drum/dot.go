package drum

const (
	DotReady = iota
	DotHit
	DotMiss
)

// No consider situations that multiple rolls are overlapped.
type Dot struct {
	Floater
	// Showtime int64 // Dot will appear at Showtime.
	Marked int
	Next   *Dot
	Prev   *Dot
}

func NewDots(rolls []*Note) (ds []*Dot) {
	for _, n := range rolls {
		var step float64
		if n.Tick >= 2 {
			step = float64(n.Duration) / float64(n.Tick-1)
		}
		for tick := 0; tick < n.Tick; tick++ {
			time := step * float64(tick)
			d := Dot{
				Floater: Floater{
					Time:  n.Time + int64(time),
					Speed: n.Speed,
				},
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
