package mania

import (
	"errors"
	"github.com/hndada/gosu/game"
	"github.com/hndada/rg-parser/osugame/osu"
	"sort"
)

const (
	typeNote game.NoteType = 1 << iota
	typeReleaseNote
	typeLongNote
)
const (
	TypeNote   = typeNote
	TypeLNHead = typeLongNote | typeNote
	TypeLNTail = typeLongNote | typeReleaseNote
)

// 난이도 및 점수 관련은 나중에
// 아래와 같이 난이도 계산에만 쓰이는 값들은 unexported로 할듯
type Note struct {
	game.BaseNote
	Key int

	position float64
	score    float64
	karma    float64
	hp       float64
	prev     int // prev note index
	next     int // next note index
	scored   bool

	hand    int
	strain  float64
	stamina float64
	// Read

	chord       []int
	trillJack   []int
	holdImpacts []float64

	baseStrain   float64
	chordPenalty float64
	trillBonus   float64
	jackBonus    float64
	holdBonus    float64 // score 필요함
}

var ErrDuration = errors.New("invalid duration: not a positive value")

// todo: Keys는 아무래도 헷갈리니, KeyCount로?
func (c *Chart) loadNotes(o *osu.Format) error {
	c.Notes = make([]Note, 0, len(o.HitObjects)*2)
	for _, ho := range o.HitObjects {
		ns, err := newNote(ho, c.Keys)
		if err != nil {
			return errors.New("invalid hit object")
		}
		c.Notes = append(c.Notes, ns...)
	}

	sort.Slice(c.Notes, func(i, j int) bool {
		if c.Notes[i].Time == c.Notes[j].Time {
			return c.Notes[i].Key < c.Notes[j].Key
		}
		return c.Notes[i].Time < c.Notes[j].Time
	})

	prevs := make([]int, c.Keys)
	for k := range prevs {
		prevs[k] = -1 // no found
	}
	for next, n := range c.Notes {
		prev := prevs[n.Key]
		c.Notes[next].prev = prev
		if prev != -1 {
			c.Notes[prev].next = next
		}
		prevs[n.Key] = next
	}
	for _, lastIdx := range prevs {
		c.Notes[lastIdx].next = -1
	}
	return nil
}

func newNote(ho osu.HitObject, keys int) ([]Note, error) {
	ns := make([]Note, 0, 2)
	var n Note
	switch ho.NoteType & osu.ComboMask {
	case osu.HitTypeHoldNote:
		n.Type = typeLongNote
	case osu.HitTypeNote:
		n.Type = TypeNote
	default:
		return ns, errors.New("invalid hit object")
	}
	n.Key = ho.Column(keys)
	n.Time = int64(ho.Time)
	n.SampleFilename = ho.HitSample.Filename
	n.SampleVolume = uint8(ho.HitSample.Volume)
	n.initSliceFields(keys)

	if n.Type == typeLongNote {
		n.Type = TypeLNHead
		n.Time2 = int64(ho.EndTime)
		ns = append(ns, n)
		if n.Time2-n.Time <= 0 {
			return ns, ErrDuration
		}

		var n2 Note
		n2.Type = TypeLNTail
		n2.Key = n.Key
		n2.Time = n.Time2
		n2.Time2 = n.Time
		n2.initSliceFields(keys)
		ns = append(ns, n2)
	} else {
		ns = append(ns, n)
	}
	return ns, nil
}

func (n *Note) initSliceFields(keys int) {
	n.chord = make([]int, keys)
	n.trillJack = make([]int, keys)
	n.holdImpacts = make([]float64, keys)
	for k := 0; k < keys; k++ {
		n.chord[k] = -1
		n.trillJack[k] = -1
	}
}
