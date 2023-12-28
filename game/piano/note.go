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
		ns[i].Position = d.Position + d.Speed*float64(n.Time-d.Time)

		// Tail's Position should be always equal or larger than Head's.
		if n.Type == Tail {
			head := ns[n.Prev]
			if n.Position < head.Position {
				ns[i].Position = head.Position
			}
		}
	}

	// linking
	prevs := make([]int, keyCount)
	exists := make([]bool, keyCount)
	for i, n := range ns {
		prev := prevs[n.Key]
		ns[i].Prev = prev
		if exists[n.Key] {
			ns[prev].Next = i
		}
		prevs[n.Key] = i
		exists[n.Key] = true
	}

	// Set each last note's Next.
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
		res.framesList[nt] = draws.NewFramesFromFile(fsys, name)
	}
}

type NotesOpts struct {
	TailOffset int32

	ws []float64
	H  float64 // Applies to all types of notes.
	xs []float64
	y  float64 // center bottom

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
	notes     []Note
	staged    []int // targets of judging
	cursor    float64
	holds     []bool
	lowests   []int // indexes of lowest notes
	animsList [][4]draws.Animation
	h         float64 // used for drawLongNoteBody
	colors    []color.NRGBA
}

func NewNotesComp(res NotesRes, opts NotesOpts, ns []Note, dys game.Dynamics) (comp NotesComp) {
	comp.notes = ns
	comp.applyTailOffset(opts.TailOffset, dys)

	keyCount := len(opts.ws)
	comp.staged = make([]int, keyCount)
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

	comp.holds = make([]bool, keyCount)
	comp.lowests = make([]int, keyCount)

	comp.animsList = make([][4]draws.Animation, keyCount)
	for k := range comp.animsList {
		for nt, frames := range res.framesList {
			a := draws.NewAnimation(frames, 400)
			a.SetSize(opts.ws[k], opts.H)
			a.Locate(opts.xs[k], opts.y, draws.CenterBottom)
			comp.animsList[k][nt] = a
		}
	}

	comp.h = opts.H
	comp.colors = make([]color.NRGBA, keyCount)
	for k := range comp.colors {
		comp.colors[k] = opts.Colors[opts.keyOrder[k]]
	}
	return
}

func (comp *NotesComp) applyTailOffset(duration int32, dys game.Dynamics) {
	// dys.Index = 0
	for i, n := range comp.notes {
		if n.Type != Tail {
			continue
		}
		dys.UpdateIndex(n.Time)
		d := dys.Current()
		comp.notes[i].Position += float64(duration) * d.Speed

		// Tail's Position should be always equal or larger than Head's.
		if head := comp.notes[n.Prev]; n.Position < head.Position {
			comp.notes[i].Position = head.Position
		}
	}
}

func (comp NotesComp) Duration() int32 {
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

func (comp *NotesComp) Update(cursor float64, holds []bool) {
	comp.updatePosition(cursor)
	comp.holds = holds
}

func (comp *NotesComp) updatePosition(cursor float64) {
	lowerBound := cursor - game.ScreenH
	for k, idx := range comp.lowests {
		var n Note
		for i := idx; i < len(comp.notes); i = n.Next {
			n = comp.notes[idx]
			if n.Position > lowerBound {
				break
			}
			// index should be updated outside of if block.
			comp.lowests[k] = i
		}

		// When Head is off the screen but Tail is on,
		// update Tail to Head since drawLongNote uses Head.
		if n.Type == Tail {
			comp.lowests[k] = n.Prev
		}
	}
	comp.cursor = cursor
}

// Notes are fixed. Lane itself moves, all notes move as same amount.
func (comp NotesComp) Draw(dst draws.Image) {
	upperBound := comp.cursor + game.ScreenH
	for k, lowest := range comp.lowests {
		idxs := []int{}
		var n Note
		for i := lowest; i < len(comp.notes); i = n.Next {
			n = comp.notes[i]
			if n.Position > upperBound {
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

			a := comp.animsList[k][n.Type]
			pos := n.Position - comp.cursor
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
	tail := comp.notes[head.Next]
	if head.Type != Head || tail.Type != Tail {
		return
	}

	a := comp.animsList[head.Key][Body]
	if comp.holds[head.Key] {
		a.Reset()
	}

	length := tail.Position - head.Position
	length += comp.h
	if length < 0 {
		length = 0
	}
	a.SetSize(a.W(), length)

	// Use head position because anchor is center bottom.
	pos := head.Position - comp.cursor
	a.Move(0, -pos)
	if head.scored {
		a.ColorScale.ScaleWithColor(color.Gray{128})
	}
	a.Draw(dst)
}
