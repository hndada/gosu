package piano

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hndada/gosu/draws"
)

// PlayNote is for in-game. Handled by pointers to modify its fields easily.
type PlayNote struct {
	Note
	Prev     *PlayNote
	Next     *PlayNote
	Marked   bool
	NextTail *PlayNote // For performance of DrawLongNoteBodies()
}

func NewPlayNotes(c *Chart) ([]*PlayNote, []*PlayNote, []*PlayNote, float64) {
	playNotes := make([]*PlayNote, 0, len(c.Notes))
	firstStagedNotes := make([]*PlayNote, c.KeyCount)
	firstLowestTails := make([]*PlayNote, c.KeyCount)
	prevs := make([]*PlayNote, c.KeyCount)
	prevTails := make([]*PlayNote, c.KeyCount)
	var weights float64
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
		weights += pn.Weight()
		playNotes = append(playNotes, pn)
	}
	return playNotes, firstStagedNotes, firstLowestTails, weights
}

// top returns top position of long note body.
func (s ScenePlay) Top(tail *PlayNote) int {
	return int(s.Position(tail.Time)+s.TailSprites[tail.Key].H/2) - 1 // Extra 1 pixel for compensating round-down
}

// bottom returns bottom position of long note body.
func (s ScenePlay) Bottom(tail *PlayNote) int {
	return int(s.Position(tail.Prev.Time)-s.HeadSprites[tail.Key].H/2) + 1 // Extra 1 pixel for compensating round-down
}

// DrawLongNoteBodies draws long sprite with Binary-building method, instead of SubImage.
// DrawLongNoteBodies draws long note before drawing Head or Tail.
// DrawLongNoteBodies just draws sub image of long note body.
func (s *ScenePlay) DrawLongNoteBodies(screen *ebiten.Image) {
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

			var pow int
			y := float64(top)
			for length := bottom - top; length > 0; length /= 2 {
				if length%2 == 0 {
					pow++
					continue
				}
				sprite := s.BodySprites[k][pow]
				sprite.Y = y
				op := sprite.Op()
				if n.Marked {
					op.ColorM.ChangeHSV(0, 0.3, 0.3)
				}
				screen.DrawImage(sprite.I, op)
				y += sprite.H // 1 << pow
				pow++
			}

			// Test case: draw long note by SubImage.
			{
				// sprite := s.BodySpritesTest[k]
				// sprite.Y = float64(top)
				// op := sprite.Op()
				// if n.Marked {
				// 	op.ColorM.ChangeHSV(0, 0.3, 0.3)
				// }
				// rect := sprite.I.Bounds()
				// rect.Max.Y = bottom - top

				// screen.DrawImage(sprite.I.SubImage(rect).(*ebiten.Image), op)
			}
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
		var sprite draws.Sprite
		switch n.Type {
		case Head:
			sprite = s.HeadSprites[n.Key]
		case Tail:
			sprite = s.TailSprites[n.Key]
		default:
			sprite = s.NoteSprites[n.Key]
		}
		sprite.Y = s.Position(n.Time) - sprite.H/2
		op := sprite.Op()
		if n.Type == Head {
			if n.Next.Marked {
				op.ColorM.ChangeHSV(0, 0.3, 0.3)
			}
		} else {
			if n.Marked {
				op.ColorM.ChangeHSV(0, 0.3, 0.3)
			}
		}
		screen.DrawImage(sprite.I, op)
	}
}

// NotePosition calculates position, the centered y-axis value.
// y = position - h/2
func (s ScenePlay) Position(time int64) float64 {
	var distance float64 // Approaching notes have positive distance, vice versa.
	tp := s.TransPoint
	cursor := s.Time()
	if time-s.Time() > 0 {
		// When there are more than 2 TransPoint in bounded time.
		for ; tp.Next != nil && tp.Next.Time < time; tp = tp.Next {
			duration := tp.Next.Time - cursor
			bpmRatio := tp.BPM / s.MainBPM
			distance += s.SpeedBase * (bpmRatio * tp.BeatScale) * float64(duration)
			cursor += duration
		}
	} else {
		for ; tp.Prev != nil && tp.Time > time; tp = tp.Prev {
			duration := tp.Time - cursor // Negative value.
			bpmRatio := tp.BPM / s.MainBPM
			distance += s.SpeedBase * (bpmRatio * tp.BeatScale) * float64(duration)
			cursor += duration
		}
	}
	bpmRatio := tp.BPM / s.MainBPM
	// Calculate the remained (which is farthest from Hint within bound).
	distance += s.SpeedBase * (bpmRatio * tp.BeatScale) * float64(time-cursor)
	return HitPosition - distance
}

// Weight is for Tail's variadic weight based on its length.
// For example, short long note does not require much strain to release.
// Todo: fine-tuning with replay data
func (n PlayNote) Weight() float64 {
	switch n.Type {
	case Tail:
		d := float64(n.Time - n.Time2)
		switch {
		case d < 50:
			return 0.5 - 0.35*d/50
		case d >= 50 && d < 200:
			return 0.15
		case d >= 200 && d < 800:
			return 0.15 + 0.85*(d-200)/600
		default:
			return 1
		}
	default:
		return 1
	}
}
