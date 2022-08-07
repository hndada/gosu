package main

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
		if prev != nil { // Set Next value later
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
func (n PlayNote) PlaySE() {}

// DrawImageOptions is not commutative. Rotate -> Scale -> Translate.
func (s *ScenePlay) DrawNotes(screen *ebiten.Image) {
	const (
		up   = 10000
		down = -2000
	)
	for _, n := range s.PlayNotes {
		td := n.Time - s.Time()
		if td > up || td < down {
			continue
		}
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
		d += s.Speed * sf.Factor * float64(n.Time-t) // remained speed factor calc

		ns := s.NoteSprites[n.Key]
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(ns.ScaleW(), ns.ScaleH())
		x := ns.X
		y := float64(HitPosition)*Scale() - d - ns.H/2 // A note locates at the center of judge line at the time.
		if n.Type == Head {
			DrawLongNote(screen, s.BodySprites[n.Key], 0, y)
		} else if n.Type == Tail && n.Time2-s.Time() < down { // Avoid drawing long note twice.
			DrawLongNote(screen, s.BodySprites[n.Key], y, float64(ScreenSizeY))
		}
		op.GeoM.Translate(x, y)
		if n.Scored {
			op.ColorM.ChangeHSV(0, 0.3, 0.3)
		}
		screen.DrawImage(s.NoteSprites[n.Key].I, op)
	}
}

// DrawLongNote draws long note before drawing Head or Tail.
func DrawLongNote(screen *ebiten.Image, ns Sprite, top, bottom float64) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(ns.ScaleW(), ns.ScaleH())
	op.GeoM.Translate(ns.X, top)
	length := bottom - top // Bottom has larger value
	for i := 0; i < int(length/ns.H); i++ {
		screen.DrawImage(ns.I, op)
		op.GeoM.Translate(0, ns.H)
	}
	last := ns.I.Bounds()
	last.Max = image.Pt(int(ns.W), int(length)%int(ns.H))
	screen.DrawImage(ns.I.SubImage(last).(*ebiten.Image), op)
}

// func isOut(w, h, x, y int) bool {
// 	return x+w < 0 || x > ScreenSizeX || y+h < 0 || y > ScreenSizeY
// }
