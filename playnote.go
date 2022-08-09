package gosu

import (
	"image"

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

// Time bound for drawing notes
const (
	up   = 10000
	down = -2000
)

// Todo: top, bottom +H/2? -H/2?
func (s *ScenePlay) DrawNotes(screen *ebiten.Image) {
	for _, n := range s.PlayNotes {
		td := n.Time - s.Time()
		if td > up || td < down {
			continue
		}
		y := s.CalcNoteY(n)
		var sprite Sprite
		switch n.Type {
		case Head:
			top := s.CalcNoteY(n.Next) + s.TailSprites[n.Key].H/2
			bottom := y - s.HeadSprites[n.Key].H/2
			if top < 0 {
				top = 0
			}
			DrawLongNote(screen, s.BodySprites[n.Key], top, bottom, n.Scored)
			sprite = s.HeadSprites[n.Key]
			sprite.Y = y - s.HeadSprites[n.Key].H/2
		case Tail:
			if n.Time2-s.Time() < down { // Avoid drawing long note twice.
				top := y + s.TailSprites[n.Key].H/2
				bottom := s.CalcNoteY(n.Prev) - s.HeadSprites[n.Key].H/2
				if bottom > screenSizeY {
					bottom = screenSizeY
				}
				DrawLongNote(screen, s.BodySprites[n.Key], top, bottom, n.Scored)
			}
			sprite = s.TailSprites[n.Key]
			sprite.Y = y - s.TailSprites[n.Key].H/2
		default:
			sprite = s.NoteSprites[n.Key]
			sprite.Y = y - s.NoteSprites[n.Key].H/2
		}
		op := sprite.Op()
		if n.Scored {
			op.ColorM.ChangeHSV(0, 0.3, 0.3)
		}
		screen.DrawImage(sprite.I, op)
	}
}
func (s ScenePlay) CalcNoteY(n *PlayNote) float64 {
	td := n.Time - s.Time()
	var d float64
	var t int64 = s.Time()
	sf := s.SpeedFactor
	if td > 0 {
		for ; sf.Next != nil && sf.Next.Time < n.Time+up; sf = sf.Next {
			gap := sf.Time - t
			d += s.Speed * sf.Factor * float64(gap) // Speed may be multiplied at once at the last
			t += gap
		}
	} else {
		for ; sf.Prev != nil && sf.Prev.Time > n.Time+down; sf = sf.Prev {
			gap := sf.Time - t
			d += s.Speed * sf.Factor * float64(gap) // Speed may be multiplied at once at the last
			t -= gap
		}
	}
	d += s.Speed * sf.Factor * float64(n.Time-t) // Remained speed factor calc
	return HintPosition - d                      // A note locates at the center of hint at the time.
}

// DrawLongNote draws long note before drawing Head or Tail.
func DrawLongNote(screen *ebiten.Image, sprite Sprite, top, bottom float64, scored bool) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(sprite.ScaleW(), sprite.ScaleH())
	op.GeoM.Translate(sprite.X, top)
	if scored {
		op.ColorM.ChangeHSV(0, 0.3, 0.3)
	}
	// length := int(bottom - top + sprite.H/2) // Bottom has larger value
	// for i := 0; i < length/int(sprite.H); i++ {
	// 	screen.DrawImage(sprite.I, op)
	// 	op.GeoM.Translate(0, sprite.H)
	// }
	// last := sprite.I.Bounds()
	// last.Max = image.Pt(int(sprite.W), length%int(sprite.H))
	l := bottom - top
	for ; l > sprite.H-25; l -= sprite.H { // Todo: resolve -25
		screen.DrawImage(sprite.I, op)
		op.GeoM.Translate(0, sprite.H)
	}
	last := sprite.I.Bounds()
	last.Max = image.Pt(int(sprite.W), int(l))
	screen.DrawImage(sprite.I.SubImage(last).(*ebiten.Image), op)
}

func (n PlayNote) PlaySE() {}
