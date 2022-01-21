package mania

import (
	"errors"
	"math"
	"sort"

	"github.com/hndada/gosu/common"
	"github.com/hndada/gosu/engine/ui"
	"github.com/hndada/rg-parser/osugame/osu"
)

const (
	typeNote common.NoteType = 1 << iota
	typeReleaseNote
	typeLongNote
)
const (
	TypeNote   = typeNote
	TypeLNHead = typeLongNote | typeNote
	TypeLNTail = typeLongNote | typeReleaseNote
)

// A note sprite's WHXY values are dependent of speed, screen size and widths-settings
type Note struct {
	common.Note
	key int

	prev   int // index of previous note
	next   int // index of next note
	score  float64
	karma  float64
	hp     float64
	scored bool

	noteDifficulty

	playSE func()

	ui.Sprite
	position      float64 // sv is applied, unscaled by speed yet
	ui.LongSprite         // temp
}

func (c *Chart) loadNotesFromOsu(o *osu.Format) error {
	c.Notes = make([]Note, 0, len(o.HitObjects)*2)
	for _, ho := range o.HitObjects {
		ns, err := newNoteFromOsu(ho, c.KeyCount, c.Path(ho.HitSample.Filename))
		if err != nil {
			return errors.New("invalid hit object")
		}
		c.Notes = append(c.Notes, ns...)
	}

	sort.Slice(c.Notes, func(i, j int) bool {
		if c.Notes[i].Time == c.Notes[j].Time {
			return c.Notes[i].key < c.Notes[j].key
		}
		return c.Notes[i].Time < c.Notes[j].Time
	})

	prevs := make([]int, c.KeyCount)
	for k := range prevs {
		prevs[k] = -1 // no found
	}
	for next, n := range c.Notes {
		prev := prevs[n.key]
		c.Notes[next].prev = prev
		if prev != -1 {
			c.Notes[prev].next = next
		}
		prevs[n.key] = next
	}
	for _, lastIdx := range prevs {
		c.Notes[lastIdx].next = -1
	}
	return nil
}

// TODO: sound should be lazy loaded
func newNoteFromOsu(ho osu.HitObject, keyCount int, sePath string) ([]Note, error) {
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
	n.key = ho.Column(keyCount)
	n.Time = int64(ho.Time)
	n.SampleFilename = ho.HitSample.Filename
	n.SampleVolume = ho.HitSample.Volume
	// if n.SampleFilename != "" {
	// 	n.playSE = audio.NewSEPlayer(sePath, int(n.SampleVolume))
	// }
	n.init(keyCount)

	if n.Type == typeLongNote { // LNTail has no sample sound
		n.Type = TypeLNHead
		n.Time2 = int64(ho.EndTime)
		ns = append(ns, n)
		if n.Time2-n.Time <= 0 {
			return ns, errors.New("invalid duration: not a positive value")
		}

		var n2 Note
		n2.Type = TypeLNTail
		n2.key = n.key
		n2.Time = n.Time2
		n2.Time2 = n.Time
		n2.init(keyCount)
		ns = append(ns, n2)
	} else {
		ns = append(ns, n)
	}
	return ns, nil
}

func (n *Note) init(keyCount int) {
	n.chord = make([]int, keyCount)
	n.trillJack = make([]int, keyCount)
	n.holdImpacts = make([]float64, keyCount)
	for k := 0; k < keyCount; k++ {
		n.chord[k] = -1
		n.trillJack[k] = -1
	}
}

// should precede to lnotes loading
// There seems little performance difference between inner loop and outer loop about range stamps
func (c *Chart) setNotePosition() {
	var cursor int
	s := c.TimeStamps[0]
	for ni, n := range c.Notes {
		for si := range c.TimeStamps[cursor:] {
			if n.Time < c.TimeStamps[cursor+si].NextTime {
				if si != 0 {
					s = c.TimeStamps[cursor+si]
					cursor += si
				}
				break
			}
		}
		c.Notes[ni].position = float64(n.Time-s.Time)*s.Factor + s.Position
	}
}

// const holdUnitHP = 0.002 // Indicates about how much HP increase per 1ms when holding LN.

func (c *Chart) allotScore() {
	var sumStrain float64
	for _, n := range c.Notes {
		sumStrain += n.strain
	}
	var avgStrain float64
	if len(c.Notes) != 0 {
		avgStrain = sumStrain / float64(len(c.Notes))
	}
	for i := range c.Notes {
		n := c.Notes[i]
		c.Notes[i].score = maxScore * (n.strain / sumStrain)
		c.Notes[i].karma = math.Min(n.strain/avgStrain, 2.5)          // 0 ~ 2.5
		c.Notes[i].hp = math.Min(n.strain/(3*avgStrain)+2.0/3.0, 1.5) // 0 ~ 1.5
	}
}
