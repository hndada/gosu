package piano

import "github.com/hndada/gosu/draws"

// drawLongNoteBodies draws stretched long note body sprite.
// Draw long note body before drawing notes.
func (s ScenePlay) drawLongNoteBodies(dst draws.Image) {
	lowerBound := s.cursor - 100
	for k, tail := range s.highestNotes {
		for ; tail != nil && tail.Position > lowerBound; tail = tail.Prev {
			if tail.Type != Tail {
				continue
			}
			head := tail.Prev

			bodyAnim := s.KeyKindNoteTypeAnimations[k][Body]
			bodyFrame := bodyAnim[0]

			if s.isKeyHolds[k] { // || s.stagedNotes[k].Type == Tail
				bodyFrame = s.drawNoteTimers[k].Frame(bodyAnim)
			}

			length := tail.Position - head.Position
			length += s.NoteHeigth
			if length < 0 {
				length = 0
			}

			bodyFrame.SetSize(bodyFrame.Width(), length)
			tailY := head.Position - s.cursor
			bodyFrame.Move(0, -tailY)

			op := draws.Op{}
			if tail.scored {
				op.ColorM.ChangeHSV(0, 0.3, 0.3)
			}
			bodyFrame.Draw(dst, op)
		}
	}
}

// Notes are fixed. Lane itself moves, all notes move same amount.
// Draw from farthest to nearest to make nearer notes priorly exposed.
func (s ScenePlay) drawNotes(dst draws.Image) {
	lowerBound := s.cursor - 100
	for k, n := range s.highestNotes {
		for ; n != nil && n.Position > lowerBound; n = n.Prev {
			// if n.Type == Tail {
			// 	s.drawLongNoteBody(dst, n)
			// }
			sprite := s.drawNoteTimers[k].Frame(s.KeyKindNoteTypeAnimations[k][n.Type])
			pos := n.Position - s.cursor
			sprite.Move(0, -pos)

			op := draws.Op{}
			if n.scored {
				op.ColorM.ChangeHSV(0, 0.3, 0.3)
			}
			sprite.Draw(dst, op)

			if n.Prev == nil {
				break
			}
		}
		// There is a case that head is off the screen
		// but tail is on the screen.
		// if n.Type == Head {
		// 	s.drawLongNoteBody(dst, n.Next)
		// }
	}
}
