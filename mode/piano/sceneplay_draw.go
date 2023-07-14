package piano

import (
	"fmt"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hndada/gosu/draws"
	"github.com/hndada/gosu/input"
	"github.com/hndada/gosu/mode"
)

// These methods are for drawing.
func (s ScenePlay) isKeyHit(k int) bool { return s.lastKeyActions[k] == input.Hit }
func (s ScenePlay) isKeyPressed(k int) bool {
	ka := s.lastKeyActions[k]
	return ka == input.Hit || ka == input.Hold
}

// sceneplay_draw.go: infers it is a part of sceneplay.go.
func (s *ScenePlay) Ticker() {
	for k := 0; k < s.Chart.KeyCount; k++ {
		s.keyTimers[k].Ticker()
		s.noteTimers[k].Ticker()
		s.keyLightingTimers[k].Ticker()
		s.hitLightingTimers[k].Ticker()
		s.holdLightingTimers[k].Ticker()
	}
	s.judgmentTimer.Ticker()
	s.comboTimer.Ticker()
}

// I used 'sprite' for local variable of Sprite instead of 's'
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
	// s.drawMeter(screen)
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

			holding := s.lastKeyActions[k] == input.Hold
			holding = holding && s.Scorer.Staged[k].Type == Tail
			if holding {
				bodyFrame = s.noteTimers[k].Frame(bodyAnim)
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
			if tail.Marked {
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
			sprite := s.noteTimers[k].Frame(s.KeyKindNoteTypeAnimations[k][n.Type])
			pos := n.Position - s.cursor
			sprite.Move(0, -pos)

			op := draws.Op{}
			if n.Marked {
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
		timer := s.keyTimers[k]
		if s.isKeyHit(k) {
			s.keyTimers[k].Reset()
		}
		index := keyUp
		// drawKeys draws for a while even when pressed off very shortly.
		if s.isKeyPressed(k) || timer.Tick < timer.MaxTick {
			index = keyDown
		}
		sprites[index].Draw(dst, draws.Op{})
	}
}

// drawKeyLightings draws for a while even when pressed off very shortly.
func (s ScenePlay) drawKeyLightings(dst draws.Image) {
	for k, sprite := range s.KeyLightingSprites {
		if s.isKeyHit(k) {
			s.keyLightingTimers[k].Reset()
		}
		timer := s.keyLightingTimers[k]
		if s.isKeyPressed(k) || timer.Tick < timer.MaxTick {
			op := draws.Op{}
			op.ColorM.ScaleWithColor(s.KeyLightingColors[k])
			sprite.Draw(dst, op)
		}
	}
}

// drawHitLightings draws when Normal is Hit or Tail is Release.
func (s ScenePlay) drawHitLightings(dst draws.Image) {
	for k, a := range s.HitLightingAnimations {
		if s.isKeyHit(k) {
			s.hitLightingTimers[k].Reset()
		}
		timer := s.hitLightingTimers[k]
		if timer.IsDone() {
			return
		}
		op := draws.Op{}
		// opaque := UserSettings.HitLightingOpaque * (1 - d.Progress(0.75, 1))
		op.ColorM.Scale(1, 1, 1, s.HitLightingOpaque)
		timer.Frame(a).Draw(dst, op)
	}
}

func (s ScenePlay) drawHoldLightings(dst draws.Image) {
	for k, a := range s.HoldLightingAnimations {
		if !s.isKeyPressed(k) {
			return
		}
		if s.isKeyHit(k) {
			s.holdLightingTimers[k].Reset()
		}
		timer := s.holdLightingTimers[k]
		op := draws.Op{}
		op.ColorM.Scale(1, 1, 1, s.HoldLightingOpaque)
		timer.Frame(a).Draw(dst, op)
	}
}

func (s ScenePlay) drawJudgment(dst draws.Image) {
	if !s.Scorer.worstJudgment.IsBlank() {
		s.judgmentTimer.Reset()
	}
	timer := s.judgmentTimer
	if timer.IsDone() {
		return
	}

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

	index := s.Scorer.judgmentIndex(s.Scorer.worstJudgment)
	sprite := timer.Frame(s.JudgmentAnimations[index])
	sprite.MultiplyScale(scale)
	sprite.Draw(dst, draws.Op{})
}

// TimeErrorMeter
// var (
// 	ColorKool = color.NRGBA{0, 170, 242, 255}   // Blue
// 	ColorCool = color.NRGBA{85, 251, 255, 255}  // Skyblue
// 	ColorGood = color.NRGBA{51, 255, 40, 255}   // Lime
// 	ColorBad  = color.NRGBA{244, 177, 0, 255}   // Yellow
// 	ColorMiss = color.NRGBA{109, 120, 134, 255} // Gray
// )

// var JudgmentColors = []color.NRGBA{
// mode.ColorKool, mode.ColorCool, mode.ColorGood, mode.ColorBad, mode.ColorMiss}

func (s ScenePlay) DebugPrint(screen draws.Image) {
	var b strings.Builder
	f := fmt.Fprintf

	f(&b, "FPS: %.2f\n", ebiten.ActualFPS())
	f(&b, "TPS: %.2f\n", ebiten.ActualTPS())
	f(&b, "Time: %.3fs/%.0fs\n", mode.ToSecond(s.Now()), mode.ToSecond(s.Duration()))
	f(&b, "\n")
	f(&b, "Score: %.0f \n", s.Score)
	f(&b, "Combo: %d\n", s.Combo)
	f(&b, "Flow: %.0f/%2d\n", s.Flow, MaxFlow)
	f(&b, " Acc: %.0f/%2d\n", s.Acc, MaxAcc)
	f(&b, "Judgment counts: %v\n", s.JudgmentCounts)
	f(&b, "\n")
	f(&b, "Speed scale (PageUp/Down): x%.2f (x%.2f)\n", s.SpeedScale, s.Speed())
	f(&b, "(Exposure time: %dms)\n", s.NoteExposureDuration())
	f(&b, "\n")
	f(&b, "Music volume (Ctrl+ Left/Right): %.0f%%\n", *s.MusicVolume*100)
	f(&b, "Sound volume (Alt+ Left/Right): %.0f%%\n", *s.SoundVolume*100)
	f(&b, "Offset (Shift+ Left/Right): %dms\n", *s.Offset)
	f(&b, "\n")
	f(&b, "Press ESC to back to choose a song.\n")
	f(&b, "Press TAB to pause.\n")
	f(&b, "Press Ctrl+ O/P to change background brightness\n")
	f(&b, "Press F12 to print debug.\n")

	ebitenutil.DebugPrint(screen.Image, b.String())
}