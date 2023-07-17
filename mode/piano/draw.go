package piano

import "github.com/hndada/gosu/draws"

// The name 'sprite' is used for local variable of Sprite instead of 's'
// to avoid confusion with the local variable of Scene.
func (s ScenePlay) Draw(screen draws.Image) {
	s.drawField(screen)
	s.drawBars(screen)
	s.drawHint(screen)

	s.drawLongNoteBodies(screen)
	s.drawNotes(screen)

	s.drawKeys(screen)
	s.drawKeyLightings(screen)
	s.drawHitLightings(screen)
	s.drawHoldLightings(screen)

	s.drawJudgment(screen)
	s.drawScore(screen)
	s.drawCombo(screen)
	// Todo: s.drawMeter(screen)
}

func (s ScenePlay) drawField(dst draws.Image) {
	s.FieldSprite.Draw(dst, draws.Op{})
}

// Bars are fixed. Lane itself moves, all bars move as same amount.
func (s ScenePlay) drawBars(dst draws.Image) {
	lowerBound := s.cursor - 100
	for b := s.highestBar; b != nil && b.Position > lowerBound; b = b.Prev {
		pos := b.Position - s.cursor
		sprite := s.BarSprite
		sprite.Move(0, -pos)
		sprite.Draw(dst, draws.Op{})
		if b.Prev == nil {
			break
		}
	}
}

func (s ScenePlay) drawHint(dst draws.Image) {
	s.HintSprite.Draw(dst, draws.Op{})
}

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

			bodyFrame.SetSize(bodyFrame.W(), length)
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

func (s ScenePlay) drawKeys(dst draws.Image) {
	for k, sprites := range s.KeySprites {
		timer := s.drawKeyTimers[k]
		index := keyUp
		// drawKeys draws for a while even when pressed off very shortly.
		if s.isKeyPresseds[k] || timer.Tick < timer.MaxTick {
			index = keyDown
		}
		sprites[index].Draw(dst, draws.Op{})
	}
}

// drawKeyLightings draws for a while even when pressed off very shortly.
func (s ScenePlay) drawKeyLightings(dst draws.Image) {
	for k, sprite := range s.KeyLightingSprites {
		timer := s.drawKeyLightingTimers[k]
		if s.isKeyPresseds[k] || timer.Tick < timer.MaxTick {
			op := draws.Op{}
			op.ColorM.ScaleWithColor(s.KeyLightingColors[k])
			sprite.Draw(dst, op)
		}
	}
}

// drawHitLightings draws when Normal is Hit or Tail is Release.
func (s ScenePlay) drawHitLightings(dst draws.Image) {
	for k, a := range s.HitLightingAnimations {
		timer := s.drawHitLightingTimers[k]
		if timer.IsDone() {
			continue
		}
		op := draws.Op{}
		// opaque := UserSettings.HitLightingOpacity * (1 - d.Progress(0.75, 1))
		op.ColorM.Scale(1, 1, 1, s.HitLightingOpacity)
		timer.Frame(a).Draw(dst, op)
	}
}

func (s ScenePlay) drawHoldLightings(dst draws.Image) {
	for k, a := range s.HoldLightingAnimations {
		if !s.isLongNoteHoldings[k] {
			continue
		}
		timer := s.drawHoldLightingTimers[k]
		op := draws.Op{}
		op.ColorM.Scale(1, 1, 1, s.HoldLightingOpacity)
		timer.Frame(a).Draw(dst, op)
	}
}

func (s ScenePlay) drawJudgment(dst draws.Image) {
	timer := s.drawJudgmentTimer
	if timer.IsDone() {
		return
	}

	// worstJudgment is guaranteed not to be blank,
	// hence no panicked by index out of range.
	index := s.judgmentIndex(s.worstJudgment)
	sprite := timer.Frame(s.JudgmentAnimations[index])

	age := timer.Age()
	const (
		bound0 = 0.1
		bound1 = 0.2
		bound2 = 0.9
	)
	scale := 1.0
	if age < bound0 {
		scale = 1 + 0.15*timer.Progress(0, bound0)
	}
	if age >= bound0 && age < bound1 {
		scale = 1.15 - 0.15*timer.Progress(bound0, bound1)
	}
	if age >= bound2 {
		scale = 1 - 0.25*timer.Progress(bound2, 1)
	}
	sprite.MultiplyScale(scale)
	sprite.Draw(dst, draws.Op{})
}

// for TimeErrorMeter
// {244, 177, 0, 255},   // Yellow
// var judgmentColors = []color.NRGBA{
// 	{0, 170, 242, 255},   // Blue
// 	{85, 251, 255, 255},  // Skyblue
// 	{51, 255, 40, 255},   // Lime
// 	{109, 120, 134, 255}, // Gray
// }
