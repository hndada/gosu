package gosu

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// PlayNote is for in-game. Handled by pointers to modify its fields easily.
type PlayNote struct {
	Note
	Prev   *PlayNote
	Next   *PlayNote
	Scored bool
}

// Second output is initial staged notes.
func NewPlayNotes(c *Chart) (playNotes []*PlayNote, stagedNotes []*PlayNote) {
	playNotes = make([]*PlayNote, 0, len(c.Notes))
	stagedNotes = make([]*PlayNote, c.KeyCount)
	prevs := make([]*PlayNote, c.KeyCount)
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
		if stagedNotes[n.Key] == nil {
			stagedNotes[n.Key] = pn
		}
		playNotes = append(playNotes, pn)
	}
	return
}

// DrawLongNote draws long note before drawing Head or Tail.
func (s *ScenePlay) DrawLongNotes(screen *ebiten.Image) {
	for _, n := range s.PlayNotes {
		if n.Type != Tail {
			continue
		}
		// Gives extra 1 pixel for compensating round-down
		top := int(s.NotePosition(n)+s.TailSprites[n.Key].H/2) - 1
		bottom := int(s.NotePosition(n.Prev)-s.HeadSprites[n.Key].H/2) + 1
		if top > screenSizeY || bottom < 0 {
			continue
		}
		if top < 0 {
			top = 0
		}
		if bottom > screenSizeY {
			bottom = screenSizeY
		}

		sprite := s.BodySprites[n.Key]
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(sprite.ScaleW(), sprite.ScaleH())
		op.GeoM.Translate(sprite.X, float64(top))
		if n.Scored {
			op.ColorM.ChangeHSV(0, 0.3, 0.3)
		}
		l, h := bottom-top, int(sprite.H) // Bottom has larger value
		q, r := l/h, l%h
		for i := 0; i < q; i++ {
			screen.DrawImage(sprite.I, op)
			op.GeoM.Translate(0, float64(h))
		}
		// last's bound is not scaled, since it's derived from the source.
		last := sprite.I.Bounds()
		last.Max.Y = last.Max.Y*r/h + 1
		screen.DrawImage(sprite.I.SubImage(last).(*ebiten.Image), op)
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
