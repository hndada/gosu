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

// Second output is initial StagedNotes
func NewPlayNotes(c *Chart) ([]*PlayNote, []*PlayNote) {
	pns := make([]*PlayNote, 0, len(c.Notes))
	prevs := make([]*PlayNote, c.KeyCount)
	stagedNotes := make([]*PlayNote, c.KeyCount)
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
		pns = append(pns, pn)
	}
	return pns, stagedNotes
}
func (n PlayNote) PlaySE() {}

// Todo: 노트 이미지
// 노트 WH는 고정, X는 Map 고정, Y만 매번 바뀜
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

		x := float64() // 미리 계산되어 Map에 저장된 것 불러오기
		y := float64(JudgeLine.Y) - d
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(x, y)
		if n.Scored {
			op.ColorM.ChangeHSV(0, 0.3, 0.3)
		}
		screen.DrawImage(s.i, op)
	}
}

// DrawLongNote draws vertically long sprite.
func DrawLongNote(screen *ebiten.Image) {
	x, y := s.X, s.Y
	op.GeoM.Translate(float64(x), float64(y))
	q, r := s.H/h1, s.H%h1+1 // quotient, remainder // TEMP: +1

	first := s.i.Bounds()
	w, h := s.W, r
	first.Min = image.Pt(0, h1-r)
	if !isOut(w, h, x, y, screen.Bounds().Size()) {
		screen.DrawImage(s.i.SubImage(first).(*ebiten.Image), op)
	}
	op.GeoM.Translate(0, float64(h))
	y += h
	h = h1
	for i := 0; i < q; i++ {
		if !s.isOut(w, h, x, y, screen.Bounds().Size()) {
			screen.DrawImage(s.i, op)
		}
		op.GeoM.Translate(0, float64(h))
		y += h
	}
}
func isOut(w, h, x, y int, screenSize image.Point) bool {
	return x+w < 0 || x > screenSize.X || y+h < 0 || y > screenSize.Y
}
