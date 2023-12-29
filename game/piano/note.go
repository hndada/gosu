package piano

import (
	"fmt"
	"image/color"
	"io/fs"
	"sort"

	"github.com/hndada/gosu/draws"
	"github.com/hndada/gosu/format/osu"
	"github.com/hndada/gosu/game"
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
	Sample   game.Sample
	Duration int32

	position float64 // Scaled x or y value.
	next     int     // For updating staged notes.
	prev     int     // For accessing to Head from Tail.
	scored   bool
}

// The length of the returned slice is 1 or 2.
func newNoteFromOsu(f osu.HitObject, keyCount int) (ns []Note) {
	n := Note{
		Time:   int32(f.Time),
		Type:   Normal,
		Key:    f.Column(keyCount),
		Sample: game.NewSample(f),
	}
	if f.NoteType&osu.ComboMask == osu.HitTypeHoldNote {
		n.Type = Head
		n.Duration = int32(f.EndTime) - n.Time
		n2 := Note{
			Time: n.Time + n.Duration,
			Type: Tail,
			Key:  n.Key,
			// Sample: game.Sample{}, // Tail has no sample sound.
			// Duration: n.Duration, // Todo: 0 or n.Duration?
		}
		ns = append(ns, n, n2)
	} else {
		ns = append(ns, n)
	}
	return ns
}

func NewNotes(f any, dys game.Dynamics, keyCount int) (ns []Note) {
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

	// Position calculation is based on Dynamics.
	// Farther note has larger position.
	// dys.Index = 0
	for i, n := range ns {
		dys.UpdateIndex(n.Time)
		d := dys.Current()
		// Ratio (speed) first, then time difference.
		ns[i].position = d.Position + d.Speed*float64(n.Time-d.Time)

		// Tail's Position should be always equal or larger than Head's.
		if n.Type == Tail {
			head := ns[n.prev]
			if n.position < head.position {
				ns[i].position = head.position
			}
		}
	}

	// linking
	prevs := make([]int, keyCount)
	exists := make([]bool, keyCount)
	for i, n := range ns {
		prev := prevs[n.Key]
		ns[i].prev = prev
		if exists[n.Key] {
			ns[prev].next = i
		}
		prevs[n.Key] = i
		exists[n.Key] = true
	}

	// Set each last note's next.
	for k, prev := range prevs {
		if !exists[k] {
			continue
		}
		ns[prev].next = len(ns)
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
		res.framesList[nt] = draws.NewFramesFromFile(fsys, name)
	}
}

type NotesOpts struct {
	TailOffset int32

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
		TailOffset: 0,

		ws: keys.ws,
		H:  20,
		xs: keys.xs,
		y:  keys.BaselineY,

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

type NotesComp struct {
	notes      []Note
	cursor     float64
	keysHold   []bool
	keysLowest []int // indexes of lowest notes
	keysAnims  [][4]draws.Animation
	h          float64 // used for drawLongNoteBody
	keysColor  []color.NRGBA
}

func NewNotesComp(res NotesRes, opts NotesOpts, ns []Note, dys game.Dynamics) (comp NotesComp) {
	comp.notes = ns
	for i, n := range comp.notes {
		if n.Type != Tail {
			continue
		}
		dys.UpdateIndex(n.Time)
		d := dys.Current()

		// Apply TailOffset to Tail's Position.
		// Tail's Position should be always equal or larger than Head's.
		comp.notes[i].position += float64(opts.TailOffset) * d.Speed
		if head := comp.notes[n.prev]; n.position < head.position {
			comp.notes[i].position = head.position
		}

		// Apply dynamics' volume to note's sample with blank volume.
		if n.Sample.Volume == 0 {
			comp.notes[i].Sample.Volume = d.Volume
		}
	}

	keyCount := len(opts.ws)
	comp.keysHold = make([]bool, keyCount)
	comp.keysLowest = make([]int, keyCount)

	comp.keysAnims = make([][4]draws.Animation, keyCount)
	for k := range comp.keysAnims {
		for nt, frames := range res.framesList {
			a := draws.NewAnimation(frames, 400)
			a.SetSize(opts.ws[k], opts.H)
			a.Locate(opts.xs[k], opts.y, draws.CenterBottom)
			comp.keysAnims[k][nt] = a
		}
	}

	comp.h = opts.H
	comp.keysColor = make([]color.NRGBA, keyCount)
	for k := range comp.keysColor {
		comp.keysColor[k] = opts.Colors[opts.keyOrder[k]]
	}
	return
}

func (comp NotesComp) Span() int32 {
	if len(comp.notes) == 0 {
		return 0
	}
	last := comp.notes[len(comp.notes)-1]
	// No need to add last.Duration, since last is
	// always either Normal or Tail.
	return last.Time
}

func (comp NotesComp) NoteCounts() []int {
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

func (comp *NotesComp) Update(cursor float64, keysHold []bool) {
	lowerBound := cursor - game.ScreenH
	for k, idx := range comp.keysLowest {
		var n Note
		for i := idx; i < len(comp.notes); i = n.next {
			n = comp.notes[idx]
			if n.position > lowerBound {
				break
			}
			// index should be updated outside of if block.
			comp.keysLowest[k] = i
		}

		// When Head is off the screen but Tail is on,
		// update Tail to Head since drawLongNote uses Head.
		if n.Type == Tail {
			comp.keysLowest[k] = n.prev
		}
	}
	comp.cursor = cursor
	comp.keysHold = keysHold
}

// Notes are fixed. Lane itself moves, all notes move as same amount.
func (comp NotesComp) Draw(dst draws.Image) {
	upperBound := comp.cursor + game.ScreenH
	for k, lowest := range comp.keysLowest {
		idxs := []int{}
		var n Note
		for i := lowest; i < len(comp.notes); i = n.next {
			n = comp.notes[i]
			if n.position > upperBound {
				break
			}
			idxs = append(idxs, i)
		}

		// Make farther notes overlapped by nearer notes.
		sort.Reverse(sort.IntSlice(idxs))

		for _, i := range idxs {
			n := comp.notes[i]
			// Make long note's body overlapped by its Head and Tail.
			if n.Type == Head {
				comp.drawLongNoteBody(dst, n)
			}

			a := comp.keysAnims[k][n.Type]
			pos := n.position - comp.cursor
			a.Move(0, -pos)
			if n.scored {
				// op.ColorM.ChangeHSV(0, 0.3, 0.3)
				a.ColorScale.ScaleWithColor(color.Gray{128})
			}
			a.Draw(dst)
		}
	}
}

// drawLongNoteBody draws stretched long note body sprite.
func (comp NotesComp) drawLongNoteBody(dst draws.Image, head Note) {
	tail := comp.notes[head.next]
	if head.Type != Head || tail.Type != Tail {
		return
	}

	a := comp.keysAnims[head.Key][Body]
	if comp.keysHold[head.Key] {
		a.Reset()
	}

	length := tail.position - head.position
	length += comp.h
	if length < 0 {
		length = 0
	}
	a.SetSize(a.W(), length)

	// Use head position because anchor is center bottom.
	pos := head.position - comp.cursor
	a.Move(0, -pos)
	if head.scored {
		a.ColorScale.ScaleWithColor(color.Gray{128})
	}
	a.Draw(dst)
}
