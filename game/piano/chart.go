package piano

import (
	"io/fs"
	"sort"

	"github.com/hndada/gosu/format/osu"
	"github.com/hndada/gosu/game"
)

// Chart contains all data needed to play a piano chart.
type Chart struct {
	game.ChartHeader
	keyCount int // Number of keys in the chart.

	mods Mods
	game.Dynamics
	bars         []Bar
	notes        []Note
	focusedNotes []int // indexes of focused notes
}

func NewChart(fsys fs.FS, name string, mods Mods) (*Chart, error) {
	c := &Chart{}
	format, hash, err := game.LoadChartFormat(fsys, name)
	if err != nil {
		return c, err
	}
	c.ChartHeader = game.NewChartHeaderFromFormat(format, hash)
	c.keyCount = c.ChartHeader.SubMode

	c.mods = mods
	c.Dynamics, err = game.NewDynamics(format)
	if err != nil {
		return c, err
	}
	c.bars = newChartBars(c.Dynamics)
	c.notes, c.focusedNotes = newChartNotes(c.keyCount, format, c.Dynamics)
	c.setStepIDs()
	return c, nil
}

// Drum and Piano modes have different bar drawing methods.
// Hence, this method is defined per game mode.
type Bar struct {
	position float64
}

// Given Dynamics' index is 0 and it is fine to modify it.
func newChartBars(dys game.Dynamics) []Bar {
	// const useDefaultMeter = 0
	times := dys.BeatTimes()
	bs := make([]Bar, len(times))
	dys.Reset()
	for i, t := range times {
		dys.UpdateIndex(t)
		bs[i] = Bar{position: dys.Position(t)}
	}
	return bs
}

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

	// Derived
	Next     int     // For updating staged notes.
	Prev     int     // For accessing to Head from Tail.
	position float64 // Scaled x or y value.
	scored   bool

	// Derived: Level calculation
	step int // Step is convenient for handling "Bomb" note
	hand int
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

func newChartNotes(keyCount int, format game.ChartFormat, dys game.Dynamics) ([]Note, []int) {
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
			if head := ns[n.Prev]; ns[i].position < head.position {
				ns[i].position = head.position
			}
		}
	}
	dys.Reset()

	// linking
	// Keys-: A slice with a length of key count
	// TODO: fix that ugly prefix

	// I once thought of a struct NoteStep which
	// accepts int value as a index, and yields Note.
	// I didn't go further as it looked too complex.
	focusedNotes := make([]int, keyCount)
	for i := range focusedNotes {
		focusedNotes[i] = -1
	}
	prevNotes := make([]int, keyCount)
	for i := range prevNotes {
		prevNotes[i] = -1
	}

	for i, n := range ns {
		prev := prevNotes[n.Key]
		ns[i].Prev = prev
		if prev != -1 {
			ns[prev].Next = i
		}
		prevNotes[n.Key] = i

		if focusedNotes[n.Key] == -1 {
			focusedNotes[n.Key] = i
		}
	}

	for _, last := range prevNotes {
		if last != -1 {
			ns[last].Next = len(ns)
		}
	}
	return ns, focusedNotes
}

func (c Chart) NoteCounts() []int {
	counts := make([]int, 2)
	for _, n := range c.notes {
		switch n.Kind {
		case Normal:
			counts[0]++
		case Head:
			counts[1]++
		}
	}
	return counts
}

func (c Chart) TotalDuration() int32 {
	ns := c.notes
	if len(ns) == 0 {
		return 0
	}

	// No need to add last.Duration, since last is
	// always either Normal or Tail.
	last := ns[len(ns)-1]
	return last.Time
}
