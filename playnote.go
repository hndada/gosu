package gosu

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// PlayNote is for in-game. Handled by pointers to modify its fields easily.
type PlayNote struct {
	Note
	Prev     *PlayNote
	Next     *PlayNote
	Scored   bool
	NextTail *PlayNote // For performance of DrawLongNotes()
}

func NewPlayNotes(c *Chart) ([]*PlayNote, []*PlayNote, []*PlayNote) {
	playNotes := make([]*PlayNote, 0, len(c.Notes))
	firstStagedNotes := make([]*PlayNote, c.KeyCount)
	firstLowestTails := make([]*PlayNote, c.KeyCount)
	prevs := make([]*PlayNote, c.KeyCount)
	prevTails := make([]*PlayNote, c.KeyCount)
	for _, n := range c.Notes {
		prev := prevs[n.Key]
		pn := &PlayNote{
			Note: n,
			Prev: prev,
		}
		if prev != nil { // Next value is set later.
			prev.Next = pn
		}
		prevs[n.Key] = pn
		if firstStagedNotes[n.Key] == nil {
			firstStagedNotes[n.Key] = pn
		}
		if n.Type == Tail {
			if prevTails[n.Key] != nil {
				prevTails[n.Key].NextTail = pn
			}
			prevTails[n.Key] = pn
			if firstLowestTails[n.Key] == nil {
				firstLowestTails[n.Key] = pn
			}
		}
		playNotes = append(playNotes, pn)
	}
	return playNotes, firstStagedNotes, firstLowestTails
}

// top returns top position of long note body.
func (s ScenePlay) Top(tail *PlayNote) int {
	return int(s.NotePosition(tail)+s.TailSprites[tail.Key].H/2) - 1 // Extra 1 pixel for compensating round-down
}

// bottom returns bottom position of long note body.
func (s ScenePlay) Bottom(tail *PlayNote) int {
	return int(s.NotePosition(tail.Prev)-s.HeadSprites[tail.Key].H/2) + 1 // Extra 1 pixel for compensating round-down
}

// DrawLongNotes draws long note before drawing Head or Tail.
// DrawLongNotes just draws sub image of long note body.
func (s *ScenePlay) DrawLongNotes(screen *ebiten.Image) {
	for k, n0 := range s.LowestTails {
		for n := n0; n != nil && s.Top(n) >= screenSizeY; n = n.NextTail {
			s.LowestTails[k] = n
		}
		for n := s.LowestTails[k]; n != nil && s.Bottom(n) >= 0; n = n.NextTail {
			top := s.Top(n)
			bottom := s.Bottom(n)
			if top < 0 {
				top = 0
			}
			if bottom > screenSizeY {
				bottom = screenSizeY
			}
			sprite := s.BodySprites[k]
			sprite.Y = float64(top)
			op := sprite.Op()
			if n.Scored {
				op.ColorM.ChangeHSV(0, 0.3, 0.3)
			}
			rect := sprite.I.Bounds()
			rect.Max.Y = bottom - top
			screen.DrawImage(sprite.I.SubImage(rect).(*ebiten.Image), op)
		}
	}
}

// Time bound for drawing notes in milliseconds.
// The tight values can be calculated by similar way of NotePosition() does.
const (
	up   = 10 * 1000
	down = -2 * 1000
)

func (s *ScenePlay) DrawNotes(screen *ebiten.Image) {
	for _, n := range s.PlayNotes {
		td := n.Time - s.Time()
		if td > up || td < down {
			continue
		}
		var sprite Sprite
		switch n.Type {
		case Head:
			sprite = s.HeadSprites[n.Key]
		case Tail:
			sprite = s.TailSprites[n.Key]
		default:
			sprite = s.NoteSprites[n.Key]
		}
		sprite.Y = s.NotePosition(n) - sprite.H/2
		op := sprite.Op()
		if n.Type == Head {
			if n.Next.Scored {
				op.ColorM.ChangeHSV(0, 0.3, 0.3)
			}
		} else {
			if n.Scored {
				op.ColorM.ChangeHSV(0, 0.3, 0.3)
			}
		}
		screen.DrawImage(sprite.I, op)
	}
}

// NotePosition calculates position, the centered y-axis value.
// y = position - h/2
func (s ScenePlay) NotePosition(n *PlayNote) float64 {
	var distance float64 // Approaching notes have positive distance, vice versa.
	tp := s.TransPoint
	time := s.Time()
	if n.Time-s.Time() > 0 {
		// When there are more than 2 TransPoint in 10 seconds.
		for ; tp.Next != nil && tp.Next.Time < n.Time; tp = tp.Next {
			duration := tp.Next.Time - time
			distance += s.Speed * tp.SpeedFactor * float64(duration)
			time += duration
		}
	} else {
		for ; tp.Prev != nil && tp.Time > n.Time; tp = tp.Prev {
			duration := tp.Time - time // Negative value.
			distance += s.Speed * tp.SpeedFactor * float64(duration)
			time += duration
		}
	}
	// Calculate the remained speed factor (which is farthest from Hint in 10 seconds.)
	distance += s.Speed * tp.SpeedFactor * float64(n.Time-time)
	return HintPosition - distance
}

func (n PlayNote) PlaySE() {}
