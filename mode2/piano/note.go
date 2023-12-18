package piano

import (
	"image/color"
	"sort"

	"github.com/hndada/gosu/format/osu"
	mode "github.com/hndada/gosu/mode2"
)

const (
	Normal = iota
	Head
	Tail
	Body
)

type Note struct {
	Time     int32
	Type     int
	Key      int
	Sample   mode.Sample
	Duration int32

	Position float64 // Scaled x or y value.
	Next     *Note
	Prev     *Note // For accessing to Head from Tail.
	scored   bool
}

type Notes []*Note

func NewNotes(f any, keyCount int) (ns Notes) {
	switch f := f.(type) {
	case *osu.Format:
		ns = make([]*Note, 0, len(f.HitObjects)*2)
		for _, ho := range f.HitObjects {
			ns = append(ns, newNoteFromOsu(ho, keyCount)...)
		}
	}

	sort.Slice(ns, func(i, j int) bool {
		if ns[i].Time == ns[j].Time {
			return ns[i].Key < ns[j].Key
		}
		return ns[i].Time < ns[j].Time
	})

	// linking
	prevs := make([]*Note, keyCount)
	for _, n := range ns {
		prev := prevs[n.Key]
		n.Prev = prev
		if prev != nil {
			prev.Next = n
		}
		prevs[n.Key] = n
	}
	return
}

// The length of the returned slice is 1 or 2.
func newNoteFromOsu(f osu.HitObject, keyCount int) (ns []*Note) {
	n := &Note{
		Time:   int32(f.Time),
		Type:   Normal,
		Key:    f.Column(keyCount),
		Sample: mode.NewSample(f),
	}
	if f.NoteType&osu.ComboMask == osu.HitTypeHoldNote {
		n.Type = Head
		n.Duration = int32(f.EndTime) - n.Time
		n2 := &Note{
			Time: n.Time + n.Duration,
			Type: Tail,
			Key:  n.Key,
			// Sample: mode.Sample{}, // Tail has no sample sound.
			// Duration: n.Duration, // Todo: 0 or n.Duration?
		}
		ns = append(ns, n, n2)
	} else {
		ns = append(ns, n)
	}
	return ns
}

type NotesConfig struct {
	SpeedScale            float64
	Heigth                float64 // Applies to all types of notes.
	Colors                [4]color.NRGBA
	TailNoteExtraDuration int32
	LongNoteBodyStyle     int // Stretch or Attach.
	UpsideDown            bool
}

func NewNotesConfig(screen mode.ScreenConfig) NotesConfig {
	return NotesConfig{
		SpeedScale: 1.0,
		Heigth:     0.03 * screen.Size.Y, // 0.03: 27px
		// The following colors are from
		// each note's second outermost pixel.
		Colors: [4]color.NRGBA{
			{255, 255, 255, 255}, // One: white
			{239, 191, 226, 255}, // Two: pink
			{218, 215, 103, 255}, // Mid: yellow
			{224, 0, 0, 255},     // Tip: red
		},
		TailNoteExtraDuration: 0,
		LongNoteBodyStyle:     0,
		UpsideDown:            false,
	}
}

type NotesComponent struct {
	Notes
}

func NewNotesComponent(ns Notes) NotesComponent {
	return NotesComponent{
		Notes: ns,
	}
}
