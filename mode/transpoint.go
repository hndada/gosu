package mode

import (
	"sort"

	"github.com/hndada/gosu/format/osu"
)

// First BPM is used as temporary main BPM.
// No two TransPoints have same Time.
type TransPoint struct {
	Time      int64
	BPM       float64
	Speed     float64
	Meter     int
	NewBeat   bool    // NewBeat draws a bar.
	Volume    float64 // Range is [0, 1].
	Highlight bool

	Position float64
	Next     *TransPoint
	Prev     *TransPoint
}

func NewTransPoints(f any) []*TransPoint {
	var transPoints []*TransPoint
	switch f := f.(type) {
	case *osu.Format:
		transPoints = newTransPointsFromOsu(f)
	}

	var prev *TransPoint
	for _, tp := range transPoints {
		tp.Prev = prev
		if prev != nil {
			prev.Next = tp
		}
		prev = tp
	}
	return nil
}

// When gathering TransPoints from osu.Format, it should input the whole slice.
// It is because osu.Format.TimingPoints brings some value from previous TimingPoint.
func newTransPointsFromOsu(f *osu.Format) []*TransPoint {
	var transPoints []*TransPoint
	sort.SliceStable(f.TimingPoints, func(i int, j int) bool {
		if f.TimingPoints[i].Time == f.TimingPoints[j].Time {
			return f.TimingPoints[i].Uninherited
		}
		return f.TimingPoints[i].Time < f.TimingPoints[j].Time
	})
	// Inherited points without Uninherited points will go dropped.
	for len(f.TimingPoints) > 0 && !f.TimingPoints[0].Uninherited {
		f.TimingPoints = f.TimingPoints[1:]
	}
	if len(f.TimingPoints) == 0 {
		return transPoints
	}
	tempMainBPM := f.TimingPoints[0].BPM()
	transPoints = make([]*TransPoint, 0, len(f.TimingPoints))
	prevBPM := tempMainBPM
	for _, timingPoint := range f.TimingPoints {
		tp := &TransPoint{
			Time:      int64(timingPoint.Time),
			BPM:       prevBPM,
			Speed:     prevBPM / tempMainBPM,
			Meter:     timingPoint.Meter,
			NewBeat:   timingPoint.Uninherited,
			Volume:    float64(timingPoint.Volume) / 100,
			Highlight: timingPoint.IsKiai(),
		}
		if timingPoint.Uninherited {
			tp.BPM = timingPoint.BPM()
			tp.Speed = tp.BPM / tempMainBPM
		} else {
			tp.Speed *= timingPoint.BeatLengthScale()
		}
		// Drop a TransPoint with a same time
		if len(transPoints) > 0 && transPoints[len(transPoints)-1].Time == tp.Time {
			// Either one makes TransPoint a NewBeat
			tp.NewBeat = transPoints[len(transPoints)-1].NewBeat || tp.NewBeat
			transPoints = transPoints[:len(transPoints)-1]
		}
		transPoints = append(transPoints, tp)
		prevBPM = tp.BPM
	}
	return transPoints
}

func (tp TransPoint) BeatDuration() float64 { return float64(tp.Meter) * (60000 / tp.BPM) }
