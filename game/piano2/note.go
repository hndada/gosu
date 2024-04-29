package piano

import (
	"fmt"
	"image/color"
	"io/fs"
	"sort"

	draws "github.com/hndada/gosu/draws5"
	"github.com/hndada/gosu/format/osu"
	"github.com/hndada/gosu/game"
)

type NoteKind int

const (
	Normal NoteKind = iota
	Head
	Tail
	Body
)

type Note struct {
	Time   int32
	Kind   NoteKind
	Key    int
	Sample game.Sample

	position float64 // Scaled x or y value.
	next     int     // For updating staged notes.
	prev     int     // For accessing to Head from Tail.
	scored   bool
}

// The length of the returned slice is 1 or 2.
func newNoteFromOsu(f osu.HitObject, keyCount int) (ns []Note) {
	n := Note{
		Time:   int32(f.Time),
		Kind:   Normal,
		Key:    f.Column(keyCount),
		Sample: game.NewSample(f),
	}
	if f.NoteType&osu.ComboMask == osu.HitTypeHoldNote {
		n.Kind = Head
		d := int32(f.EndTime) - n.Time
		n2 := Note{
			Time: n.Time + d,
			Kind: Tail,
			Key:  n.Key,
			// Tail has no sample sound.
		}
		ns = append(ns, n, n2)
	} else {
		ns = append(ns, n)
	}
	return ns
}

type Notes struct {
	keyCount     int
	notes        []Note
	keysFocus    []int // indexes of focused notes
	sampleBuffer []game.Sample
	// none         int   // index of none value. It is same as len(notes).
}

func NewNotes(keyCount int, dys game.Dynamics, chart game.ChartFormat) Notes {
	var ns []Note
	switch chart := chart.(type) {
	case *osu.Format:
		ns = make([]Note, 0, len(chart.HitObjects)*2)
		for _, ho := range chart.HitObjects {
			ns = append(ns, newNoteFromOsu(ho, keyCount)...)
		}
		// keyCount = int(chart.CircleSize)
	}

	sort.Slice(ns, func(i, j int) bool {
		if ns[i].Time == ns[j].Time {
			return ns[i].Key < ns[j].Key
		}
		return ns[i].Time < ns[j].Time
	})

	// Position calculation is based on Dynamics.
	// Farther note has larger position.
	dys.Reset()
	for i, n := range ns {
		dys.UpdateIndex(n.Time)
		ns[i].position = dys.Position(n.Time)

		// Tail's Position should be always equal or larger than Head's.
		if n.Kind == Tail {
			if head := ns[n.prev]; n.position < head.position {
				ns[i].position = head.position
			}
		}
	}

	// linking
	none := len(ns)
	keysNone := make([]int, keyCount)
	for k := range keysNone {
		keysNone[k] = none
	}
	keysFocus := make([]int, keyCount)
	copy(keysFocus, keysNone)
	keysPrev := make([]int, keyCount)
	copy(keysPrev, keysNone)

	for i, n := range ns {
		prev := keysPrev[n.Key]
		ns[i].prev = prev
		if prev != none {
			ns[prev].next = i
		}
		keysPrev[n.Key] = i

		if keysFocus[n.Key] == none {
			keysFocus[n.Key] = i
		}
	}
	// Set each last note's next with none.
	for _, last := range keysPrev {
		if last != none {
			ns[last].next = none
		}
	}

	return Notes{
		keyCount: keyCount,
		notes:    ns,
		// none:         none,
		keysFocus:    keysFocus,
		sampleBuffer: nil,
	}
}

func (ns Notes) keysFocusNote() []Note {
	kn := make([]Note, ns.keyCount)
	for k, ni := range ns.keysFocus {
		if ni == len(ns.notes) {
			continue
		}
		kn[k] = ns.notes[ni]
	}
	return kn
}

type NotesResources struct {
	framesList [4]draws.Frames
	// defaultSampleData []byte
}

// When note/normal image is not found, use default's note/normal.
// When note/head image is not found, use user's note/normal.
// When note/tail image is not found, let it be blank.
// When note/body image is not found, use user's note/normal.
func (res *NotesResources) Load(fsys fs.FS) {
	for nk, nkn := range []string{"normal", "head", "tail", "body"} {
		name := fmt.Sprintf("piano/note/%s.png", nkn)
		res.framesList[nk] = draws.NewFramesFromFile(fsys, name)
	}
}

type NotesOptions struct {
	keyCount   int
	keyOrder   []KeyKind
	keysW      []float64
	H          float64 // Applies to all types of notes.
	keysX      []float64
	y          float64 // center bottom
	TailOffset int32
	Colors     [4]color.NRGBA
	// LongBodyStyle     int // Stretch or Attach.
}

func NewNotesOptions(stage StageOptions, keys KeysOptions) NotesOptions {
	return NotesOptions{
		keyCount:   stage.keyCount,
		keyOrder:   keys.Order(),
		keysW:      keys.w,
		H:          20,
		keysX:      keys.x,
		y:          keys.y,
		TailOffset: 0,
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

type NotesComponent struct {
	keysAnims [][4]draws.Animation
	Notes
	keysLowest  []int // indexes of lowest notes
	cursor      float64
	keysColor   []color.NRGBA
	keysHolding []bool
	// h           float64 // used for drawLongNoteBody
}

func NewNotesComponent(res *Resources, opts *Options, c *Chart) (cmp *NotesComponent) {
	cmp.keysAnims = make([][4]draws.Animation, c.keyCount)
	for k := range cmp.keysAnims {
		for nk, frames := range res.NotesFramesList {
			a := draws.NewAnimation(frames, 400)
			a.SetSize(opts.keysW[k], opts.H)
			if nk == int(Body) {
				a.Locate(opts.keysX[k], opts.y, draws.CenterTop)
			} else {
				a.Locate(opts.keysX[k], opts.y, draws.CenterBottom)
			}
			cmp.keysAnims[k][nk] = a
		}
	}

	// Apply default sample values.
	dys.Reset()
	for i, n := range ns.notes {
		d := dys.UpdateIndex(n.Time)
		// if n.Sample.Filename == "" {
		// 	ns.notes[i].Sample.Filename = res.defaultSampleName
		// }
		if n.Sample.Volume == 0 {
			ns.notes[i].Sample.Volume = d.Volume
		}
	}

	// Apply TailOffset to Tail's Position.
	dys.Reset()
	for i, n := range ns.notes {
		if n.Kind != Tail {
			continue
		}
		d := dys.UpdateIndex(n.Time)
		ns.notes[i].position += float64(opts.TailOffset) * d.Speed

		// Tail's Position should be always equal or larger than Head's.
		if head := ns.notes[n.prev]; n.position < head.position {
			ns.notes[i].position = head.position
		}
	}
	cmp.Notes = ns

	cmp.keysLowest = make([]int, opts.keyCount)
	cmp.keysColor = make([]color.NRGBA, opts.keyCount)
	for k := range cmp.keysColor {
		cmp.keysColor[k] = opts.Colors[opts.keyOrder[k]]
	}
	cmp.keysHolding = make([]bool, opts.keyCount)
	// cmp.h = opts.H
	return
}

func (cmp *NotesComponent) Update(ka game.KeyboardAction, cursor float64) {
	lowermost := cursor - game.ScreenSizeY
	for k, lowest := range cmp.keysLowest {
		for ni := lowest; ni < len(cmp.notes); ni = cmp.notes[ni].next {
			n := cmp.notes[lowest]
			if n.position > lowermost {
				break
			}
			// index should be updated outside of if block.
			cmp.keysLowest[k] = ni
		}
		// When Head is off the screen but Tail is on,
		// update Tail to Head since drawLongNote uses Head.
		ni := cmp.keysLowest[k]
		if n := cmp.notes[ni]; n.Kind == Tail {
			cmp.keysLowest[k] = n.prev
		}
	}
	cmp.cursor = cursor
	cmp.keysHolding = ka.KeysHolding()
}

// Notes are fixed. Lane itself moves, all notes move as same amount.
func (cmp NotesComponent) Draw(dst draws.Image) {
	uppermost := cmp.cursor + game.ScreenSizeY
	for k, lowest := range cmp.keysLowest {
		var nis []int
		for ni := lowest; ni < len(cmp.notes); ni = cmp.notes[ni].next {
			n := cmp.notes[ni]
			if n.position > uppermost {
				break
			}
			nis = append(nis, ni)
		}

		// Make farther notes overlapped by nearer notes.
		sort.Sort(sort.Reverse(sort.IntSlice(nis)))

		for _, ni := range nis {
			n := cmp.notes[ni]
			// Make long note's body overlapped by its Head and Tail.
			if n.Kind == Head {
				cmp.drawLongNoteBody(dst, n)
			}

			a := cmp.keysAnims[k][n.Kind]
			pos := n.position - cmp.cursor
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
func (cmp NotesComponent) drawLongNoteBody(dst draws.Image, head Note) {
	tail := cmp.notes[head.next]
	if head.Kind != Head || tail.Kind != Tail {
		return
	}

	a := cmp.keysAnims[head.Key][Body]
	if !cmp.keysHolding[head.Key] {
		a.Reset()
	}

	length := tail.position - head.position
	// length += cmp.h
	if length < 0 {
		length = 0
	}
	a.SetSize(a.W(), length)

	// Use head position because anchor is center bottom.
	pos := head.position - cmp.cursor
	a.Move(0, -pos)
	if head.scored {
		a.ColorScale.ScaleWithColor(color.Gray{128})
	}
	a.Draw(dst)
}
