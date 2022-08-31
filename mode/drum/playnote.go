package drum

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

func NewPlayNotes(c *Chart) ([]*PlayNote, *PlayNote, *PlayNote, float64) {
	playNotes := make([]*PlayNote, 0, len(c.Notes))
	var (
		prev            *PlayNote
		prevTail        *PlayNote
		firstStagedNote *PlayNote
		firstTail       *PlayNote
		weights         float64
	)
	for _, n := range c.Notes {
		pn := &PlayNote{
			Note: n,
			Prev: prev,
		}
		if prev != nil { // Next value is set later.
			prev.Next = pn
		}
		prev = pn
		if firstStagedNote == nil {
			firstStagedNote = pn
		}
		if n.Type == Tail {
			if prevTail != nil {
				prevTail.NextTail = pn
			}
			prevTail = pn
			if firstTail == nil {
				firstTail = pn
			}
		}
		weights += pn.Weight()
		playNotes = append(playNotes, pn)
	}
	return playNotes, firstStagedNote, firstTail, weights
}

// Right returns right position of Roll body.
// Extra 1 pixel for compensating round-down
func (s ScenePlay) Right(tail *PlayNote) int {
	return int(tail.Position(tail.Time)) - 1 // +s.TailSprites[tail.NoteKind()].W/2
}

// Left returns left position of Roll body.
// Extra 1 pixel for compensating round-down
func (s ScenePlay) Left(tail *PlayNote) int {
	return int(tail.Position(tail.Prev.Time)) + 1 // -s.HeadSprites[tail.NoteKind()].W/2
}

// DrawRollBodies draws long sprite with Binary-building method, instead of SubImage.
// DrawRollBodies draws long note before drawing Head or Tail.
// DrawRollBodies just draws sub image of long note body.
// Right is Tail, and Left is Head.
func (s *ScenePlay) DrawRollBodies(screen *ebiten.Image) {
	for n := s.LeadingTail; n != nil && s.Right(n) < 0; n = n.NextTail {
		s.LeadingTail = n
	}
	for n := s.LeadingTail; n != nil && s.Left(n) >= screenSizeX; n = n.NextTail {
		right := s.Right(n)
		left := s.Left(n)
		if right >= screenSizeX {
			right = screenSizeX
		}
		if left < 0 {
			left = 0
		}

		var pow int
		x := float64(left)
		// y := float64(top)
		for length := right - left; length > 0; length /= 2 {
			if length%2 == 0 {
				pow++
				continue
			}
			sprite := s.BodySprites[n.NoteKind()][pow]
			sprite.X = x
			op := sprite.Op()
			// Todo: tint to orange
			screen.DrawImage(sprite.I, op)
			x += sprite.W // 1 << pow
			pow++
		}
	}
}

// Time bound for drawing notes in milliseconds.
// The tight values can be calculated by similar way of NotePosition() does.
const (
	up   = 10 * 1000
	down = -2 * 1000
)

// Todo: Seek around staged notes only
func (s *ScenePlay) DrawNotes(screen *ebiten.Image) {
	for _, n := range s.PlayNotes {
		td := n.Time - s.Time()
		if td > up || td < down {
			continue
		}
		var sprite draws.Sprite
		switch n.Type {
		case Don, BigDon:
			sprite = s.DonSprites[n.NoteKind()][NoteLayerGround]
		case Kat, BigKat:
			sprite = s.KatSprites[n.NoteKind()][NoteLayerGround]
		case Head, BigHead:
			sprite = s.HeadSprites[n.NoteKind()]
		case Tail, BigTail:
			sprite = s.TailSprites[n.NoteKind()]
		case Shake:
			sprite = s.ShakeSprites[ShakeNote]
		}
		sprite.SetCenterX(n.Position(s.Time()))
		// sprite.X = s.Position(n.Time) - sprite.W/2
		op := sprite.Op()
		if n.Marked { // When a note is marked, next note is going to be staged.
			switch n.Type {
			case Don, BigDon, Kat, BigKat, Shake: // Todo: Shake fade out effect?
				continue
			case Head, BigHead, Tail, BigTail:
				// Todo: color the Roll based on hit count
			}
		}
		screen.DrawImage(sprite.I, op)
	}
}

// NotePosition calculates position, the centered y-axis value.
// y = position - h/2
// Todo: 2 type on note position. Mania type, Taiko type.
func (n PlayNote) Position(time int64) float64 {
	td := n.Time - time
	return HitPosition + SpeedScale*n.Speed*float64(td)
}

// func (s ScenePlay) Position(time int64) float64 {
// 	var distance float64 // Approaching notes have positive distance, vice versa.
// 	tp := s.TransPoint
// 	cursor := s.Time()
// 	if time-s.Time() > 0 {
// 		// When there are more than 2 TransPoint in bounded time.
// 		for ; tp.Next != nil && tp.Next.Time < time; tp = tp.Next {
// 			duration := tp.Next.Time - cursor
// 			bpmRatio := tp.BPM / s.MainBPM
// 			distance += s.SpeedScale * (bpmRatio * tp.BeatLengthScale) * float64(duration)
// 			cursor += duration
// 		}
// 	} else {
// 		for ; tp.Prev != nil && tp.Time > time; tp = tp.Prev {
// 			duration := tp.Time - cursor // Negative value.
// 			bpmRatio := tp.BPM / s.MainBPM
// 			distance += s.SpeedScale * (bpmRatio * tp.BeatLengthScale) * float64(duration)
// 			cursor += duration
// 		}
// 	}
// 	bpmRatio := tp.BPM / s.MainBPM
// 	// Calculate the remained (which is farthest from Hint within bound).
// 	distance += s.SpeedScale * (bpmRatio * tp.BeatLengthScale) * float64(time-cursor)
// 	return HitPosition - distance
// }

// Weight is for Tail's variadic weight based on its length.
// For example, short long note does not require much strain to release.
// Todo: fine-tuning with replay data
func (n PlayNote) Weight() float64 {
	switch n.Type {
	case Don, Kat:
		return 1
	case BigDon, BigKat:
		return 1.1
	default:
		return 0
	}
}
