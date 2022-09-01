package gosu

import (
	"github.com/hndada/gosu/format/osu"
)

// type NoteType int

// const (
// 	Normal = iota
// 	Head
// 	Tail
// 	Body
// 	BodyTick // e.g., Roll tick in Drum mode.
// 	Extra    // e.g., Shake in Drum mode.
// 	ExtraTick
// )

// Strategy of Piano mode
// Calculate position of each note in advance
// Parameter: SpeedScale, BPM Ratio, BeatLengthScale
// Calculate current HitPosition only.
// For other notes, just calculate the difference between HitPosition.
type BaseNote struct {
	Type         int
	Time         int64
	Time2        int64
	SampleName   string // aka SampleFilename.
	SampleVolume float64

	Position float64 // Scaled x or y value.
	Marked   bool
}

func NewBaseNote(f any) (n BaseNote) {
	switch f := f.(type) {
	case osu.HitObject:
		n = BaseNote{
			Time:         int64(f.Time),
			Time2:        int64(f.Time),
			SampleName:   f.HitSample.Filename,
			SampleVolume: float64(f.HitSample.Volume) / 100,
		}
	}
	return n
}
