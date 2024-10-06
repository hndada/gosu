package piano

import (
	"image/color"
	"sort"

	"github.com/hndada/gosu/draws"
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
	keyCount  int
	data      []Note
	keysFocus []int // indexes of focused notes
	// sampleBuffer []game.Sample
	// none         int   // index of none value. It is same as len(notes).
}

func NewNotes(keyCount int, format game.ChartFormat, dys game.Dynamics) Notes {
	var ns []Note
	switch format := format.(type) {
	case *osu.Format:
		ns = make([]Note, 0, len(format.HitObjects)*2)
		for _, ho := range format.HitObjects {
			ns = append(ns, newNoteFromOsu(ho, keyCount)...)
		}
		// keyCount = int(format.CircleSize)
	}

	sort.Slice(ns, func(i, j int) bool {
		if ns[i].Time == ns[j].Time {
			return ns[i].Key < ns[j].Key
		}
		return ns[i].Time < ns[j].Time
	})

	// Position calculation is based on Dynamics.
	// Farther note has larger position.
	// Todo: dys.Reset() looks not pretty.
	dys.Reset()
	for i, n := range ns {
		dys.UpdateIndex(n.Time)
		ns[i].position = dys.Position(n.Time)

		// Tail's Position should be always equal or larger than Head's.
		if ns[i].Kind == Tail {
			if head := ns[n.prev]; ns[i].position < head.position {
				ns[i].position = head.position
			}
		}
	}
	dys.Reset()

	// linking
	// none := len(ns)
	none := -1
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
		keyCount:  keyCount,
		data:      ns,
		keysFocus: keysFocus,
	}
}

type NotesComponent struct {
	notes            Notes
	keysAnims        [][4]draws.Animation
	keysLowest       []int // indexes of lowest notes
	cursor           float64
	scaledScreenSize float64 // added in 241006
	keysColor        []color.NRGBA
	keysHolding      []bool
	// h           float64 // used for drawLongNoteBody
}

func NewNotesComponent(res *Resources, opts *Options, c *Chart) (cmp NotesComponent) {
	cmp.keysAnims = make([][4]draws.Animation, c.keyCount)
	for k := range cmp.keysAnims {
		for nk, frames := range res.NotesFramesList {
			a := draws.NewAnimation(frames, 400)
			w := opts.keyWidthsMap[c.keyCount][k]
			h := opts.NoteHeight
			a.SetSize(w, h)

			x := opts.keyPositionXsMap[c.keyCount][k]
			y := opts.KeyPositionY
			if nk == int(Body) {
				a.Locate(x, y, draws.CenterTop)
			} else {
				a.Locate(x, y, draws.CenterBottom)
			}
			cmp.keysAnims[k][nk] = a
		}
	}

	// Apply default sample values.
	cd := c.FuncCurrentDynamic()
	// dys.Reset()
	for i, n := range c.data {
		d := cd(n.Time)
		// d := c.Dynamics.UpdateIndex(n.Time)
		// if n.Sample.Filename == "" {
		// 	c.notes.data[i].Sample.Filename = res.defaultSampleName
		// }
		if n.Sample.Volume == 0 {
			c.data[i].Sample.Volume = d.Volume
		}
	}

	// Apply TailOffset to Tail's Position.
	cd = c.FuncCurrentDynamic()
	// dys.Reset()
	for i, n := range c.Notes.data {
		if n.Kind != Tail {
			continue
		}
		d := cd(n.Time)
		// d := dys.UpdateIndex(n.Time)
		c.Notes.data[i].position += float64(opts.TailNoteOffset) * d.Speed

		// Tail's Position should be always equal or larger than Head's.
		if n.prev == -1 {
			continue
		}
		if head := c.Notes.data[n.prev]; n.position < head.position {
			c.Notes.data[i].position = head.position
		}
	}

	cmp.notes = c.Notes
	cmp.keysLowest = make([]int, c.keyCount)
	copy(cmp.keysLowest, cmp.notes.keysFocus)
	cmp.scaledScreenSize = opts.screenSizeY * opts.SpeedScale

	cmp.keysColor = make([]color.NRGBA, c.keyCount)
	order := opts.KeyOrders[c.keyCount]
	for k := range cmp.keysColor {
		cmp.keysColor[k] = opts.NoteColors[order[k]]
	}
	cmp.keysHolding = make([]bool, c.keyCount)
	// cmp.h = opts.H

	return
}

func (cmp *NotesComponent) Update(ka game.KeyboardAction, cursor float64) {
	lowermost := cursor - cmp.scaledScreenSize
	for k, lowest := range cmp.keysLowest {
		for ni := lowest; ni < len(cmp.notes.data); ni = cmp.notes.data[ni].next {
			if ni < 0 {
				break
			}
			n := cmp.notes.data[ni]
			if n.position > lowermost {
				break
			}
			// index should be updated outside of if block.
			cmp.keysLowest[k] = ni
		}
		// When Head is off the screen but Tail is on,
		// update Tail to Head since drawLongNote uses Head.
		ni := cmp.keysLowest[k]
		if ni < 0 {
			continue
		}
		if n := cmp.notes.data[ni]; n.Kind == Tail {
			cmp.keysLowest[k] = n.prev
		}
	}
	cmp.cursor = cursor
	cmp.keysHolding = ka.KeysHolding()
}

// Notes are fixed. Lane itself moves, all notes move as same amount.
func (cmp NotesComponent) Draw(dst draws.Image) {
	uppermost := cmp.cursor + cmp.scaledScreenSize
	for k, lowest := range cmp.keysLowest {
		var nis []int
		for ni := lowest; ni < len(cmp.notes.data); ni = cmp.notes.data[ni].next {
			if ni < 0 {
				break
			}
			n := cmp.notes.data[ni]
			if n.position > uppermost {
				break
			}
			nis = append(nis, ni)
		}

		// Make farther notes overlapped by nearer notes.
		sort.Sort(sort.Reverse(sort.IntSlice(nis)))

		for _, ni := range nis {
			if ni < 0 {
				break
			}
			n := cmp.notes.data[ni]
			// Make long note's body overlapped by its Head and Tail.
			if n.Kind == Head {
				cmp.drawLongNoteBody(dst, n)
			}

			a := cmp.keysAnims[k][n.Kind]
			pos := n.position - cmp.cursor
			a.Move(0, -pos)
			a.ColorScale.ScaleWithColor(cmp.keysColor[k])
			if n.scored {
				a.ColorScale.ScaleWithColor(color.Gray{128})
				// op.ColorM.ChangeHSV(0, 0.3, 0.3)
			}
			a.Draw(dst)
		}
	}
}

// drawLongNoteBody draws stretched long note body sprite.
func (cmp NotesComponent) drawLongNoteBody(dst draws.Image, head Note) {
	tail := cmp.notes.data[head.next]
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

	pos := tail.position - cmp.cursor
	a.Move(0, -pos)
	a.ColorScale.ScaleWithColor(cmp.keysColor[tail.Key])
	if head.scored {
		a.ColorScale.ScaleWithColor(color.Gray{128})
	}
	a.Draw(dst)
}
