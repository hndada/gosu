package piano

import (
	"fmt"
	"image/color"
	"io/fs"
	"sort"

	"github.com/hndada/gosu/draws"
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
	Next     int     // For updating staged notes.
	Prev     int     // For accessing to Head from Tail.
	scored   bool
}

// The length of the returned slice is 1 or 2.
func newNoteFromOsu(f osu.HitObject, keyCount int) (ns []Note) {
	n := Note{
		Time:   int32(f.Time),
		Type:   Normal,
		Key:    f.Column(keyCount),
		Sample: mode.NewSample(f),
	}
	if f.NoteType&osu.ComboMask == osu.HitTypeHoldNote {
		n.Type = Head
		n.Duration = int32(f.EndTime) - n.Time
		n2 := Note{
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

func NewNotes(f any, keyCount int) (ns []Note) {
	switch f := f.(type) {
	case *osu.Format:
		ns = make([]Note, 0, len(f.HitObjects)*2)
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
	prevs := make([]int, keyCount)
	exists := make([]bool, keyCount)
	for i, n := range ns {
		prev := prevs[n.Key]
		n.Prev = prev
		if exists[n.Key] {
			ns[prev].Next = i
		}
		prevs[n.Key] = i
		exists[n.Key] = true
	}

	// Set Note.Next of the last note of each key.
	for k, prev := range prevs {
		if !exists[k] {
			continue
		}
		ns[prev].Next = len(ns)
	}
	return
}

type NotesRes struct {
	framesList [4]draws.Frames
}

// When note/normal image is not found, use default's note/normal.
// When note/head image is not found, use user's note/normal.
// When note/tail image is not found, let it be blank.
// When note/body image is not found, use user's note/normal.
// Todo: remove key kind folders
func (res *NotesRes) Load(fsys fs.FS) {
	for nt, ntname := range []string{"normal", "head", "tail", "body"} {
		name := fmt.Sprintf("piano/note/one/%s.png", ntname)
		res.framesList[nt] = draws.NewFramesFromFilename(fsys, name)
	}
}

type NotesOpts struct {
	SpeedScale        float64
	lastSpeedScale    float64
	TailExtraDuration int32

	keyCount int
	ws       []float64
	H        float64 // Applies to all types of notes.
	xs       []float64
	y        float64 // center bottom

	keyOrder []KeyKind
	Colors   [4]color.NRGBA
	// LongBodyStyle     int // Stretch or Attach.
}

func NewNotesOpts(keys KeysOpts) NotesOpts {
	return NotesOpts{
		SpeedScale:        1.0,
		lastSpeedScale:    1.0,
		TailExtraDuration: 0,

		keyCount: keys.Count,
		ws:       keys.ws,
		H:        20,
		xs:       keys.xs,
		y:        keys.BaselineY,

		keyOrder: keys.Order(),
		// Colors are from each note's second outermost pixel.
		Colors: [4]color.NRGBA{
			{255, 255, 255, 255}, // One: white
			{239, 191, 226, 255}, // Two: pink
			{218, 215, 103, 255}, // Mid: yellow
			{224, 0, 0, 255},     // Tip: red
		},
		// LongBodyStyle: 0,
	}
}

type NoteComp struct {
	notes     []Note
	staged    []int // targets of judging
	idxs      []int // indexes of lowest notes
	cursor    float64
	animsList [][4]draws.Animation
	colors    []color.NRGBA
}

func NewNoteComp(res NotesRes, opts NotesOpts, ns []Note) (comp NoteComp) {
	comp.notes = ns

	comp.staged = make([]int, opts.keyCount)
	for k := range comp.staged {
		found := false
		for i, n := range ns {
			if k == n.Key {
				comp.staged[n.Key] = i
				found = true
				break
			}
		}
		if !found {
			comp.staged[k] = len(ns)
		}
	}

	comp.animsList = make([][4]draws.Animation, opts.keyCount)
	for k := range comp.animsList {
		for nt, frames := range res.framesList {
			anim := draws.NewAnimation(frames, mode.ToTick(400))
			anim.SetSize(opts.ws[k], opts.H)
			anim.Locate(opts.xs[k], opts.y, draws.CenterBottom)
			comp.animsList[k][nt] = anim
		}
	}

	comp.colors = make([]color.NRGBA, opts.keyCount)
	for k := range comp.colors {
		comp.colors[k] = opts.Colors[opts.keyOrder[k]]
	}
	return
}

func (comp NoteComp) Duration() int32 {
	if len(comp.notes) == 0 {
		return 0
	}
	last := comp.notes[len(comp.notes)-1]
	// No need to add last.Duration, since last is
	// always either Normal or Tail.
	return last.Time
}

func (comp NoteComp) NoteCounts() []int {
	counts := make([]int, 2)
	for _, n := range comp.notes {
		switch n.Type {
		case Normal:
			counts[0]++
		case Head:
			counts[1]++
		}
	}
	return counts
}

func (comp *NoteComp) Update(cursor float64) {
	lowerBound := cursor - mode.ScreenH
	for k, idx := range comp.idxs {
		var n Note
		for i := idx; i < len(comp.notes); i = n.Next {
			n = comp.notes[idx]
			if n.Position > lowerBound {
				break
			}
			// index should be updated outside of if block.
			comp.idxs[k] = i
		}

		// Update Head to Tail since drawLongNoteBody uses Tail.
		if n.Type == Head {
			comp.idxs[k] = n.Next
		}
	}
	comp.cursor = cursor
}

// Notes are fixed. Lane itself moves, all notes move as same amount.
func (comp NoteComp) Draw(dst draws.Image) {
	upperBound := comp.cursor + mode.ScreenH
	for k, idx := range comp.idxs {
		var n Note
		for i := idx; i < len(comp.notes); i = n.Next {
			n = comp.notes[i]
			if n.Position > upperBound {
				break
			}

			pos := n.Position - comp.cursor
			anim := comp.animsList[k][n.Type]
			anim.Move(0, -pos)
			anim.Draw(dst)
		}
	}
}