package piano

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hndada/gosu/draws"
	"github.com/hndada/gosu/input"
)

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
	s.Field.Draw(dst, draws.Op{})
}

// Bars are fixed. Lane itself moves, all bars move as same amount.
func (s ScenePlay) drawBars(dst draws.Image) {
	lowerBound := s.cursor - 100
	for b := s.highestBar; b.Position > lowerBound; b = b.Prev {
		pos := b.Position - s.cursor
		sprite := s.Bar
		sprite.Move(0, -pos)
		sprite.Draw(dst, draws.Op{})
		if b.Prev == nil {
			break
		}
	}
}

func (s ScenePlay) drawHint(dst draws.Image) {
	s.Hint.Draw(dst, draws.Op{})
}

// drawLongNoteBody draws stretched long note body sprite.
// Draw long note body before drawing notes.
func (s ScenePlay) drawLongNoteBodies(dst draws.Image) {
	for k, tail := range s.highestNotes {
		if tail.Type != Tail {
			continue
		}
		head := tail.Prev
		body := s.NoteTypes[k][Body][0]

		holding := s.lastKeyActions[k] == input.Hold
		holding = holding && s.Scorer.Staged[k].Type == Tail
		if holding {
			body = s.noteTimers[k].Frame(s.NoteTypes[k][Body])
		}

		length := tail.Position - head.Position
		length += s.NoteHeigth
		if length < 0 {
			length = 0
		}

		body.SetSize(body.W(), length)
		tailY := head.Position - s.cursor
		body.Move(0, -tailY)

		op := draws.Op{}
		if tail.Marked {
			op.ColorM.ChangeHSV(0, 0.3, 0.3)
		}
		body.Draw(dst, op)
	}
}

// Notes are fixed. Lane itself moves, all notes move same amount.
// Draw from farthest to nearest to make nearer notes priorly exposed.
func (s ScenePlay) drawNotes(dst draws.Image) {
	lowerBound := s.cursor - 100
	for k, n := range s.highestNotes {
		for ; n.Position > lowerBound; n = n.Prev {
			sprite := s.noteTimers[k].Frame(s.NoteTypes[k][n.Type])
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
	}
}

func (s ScenePlay) drawKeys(dst draws.Image) {
	for k, sprites := range s.KeysUpDowns {
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
	for k, sprite := range s.KeyLightings {
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
	for k, a := range s.HitLightings {
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
	for k, a := range s.HoldLightings {
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
	sprite := timer.Frame(s.JudgmentKinds[index])
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
	fps := fmt.Sprintf("FPS: %.2f\n", ebiten.ActualFPS())
	tps := fmt.Sprintf("TPS: %.2f\n", ebiten.ActualTPS())
	time := fmt.Sprintf("Time: %.3fs/%.0fs\n", float64(s.Now())/1000, float64(s.Chart.Duration())/1000)

	score := fmt.Sprintf("Score: %.0f \n", s.Scorer.Score)
	combo := fmt.Sprintf("Combo: %d\n", s.Scorer.Combo)
	flow := fmt.Sprintf("Flow: %.2f%%\n", s.Scorer.Flow/MaxFlow*100)
	acc := fmt.Sprintf("Acc: %.2f%%\n", s.Scorer.Acc/MaxAcc*100)
	judgmentCount := fmt.Sprintf("Judgment counts: %v\n", s.Scorer.JudgmentCounts)

	speedScale := fmt.Sprintf("Speed scale (Z/X): %.0f (x%.2f)\n", s.SpeedScale, s.Dynamic.Speed)
	exposureTime := fmt.Sprintf("(Exposure time: %.fms)\n", s.ExposureTime(s.Speed()))

	musicVolume := fmt.Sprintf("Music volume (Ctrl+ Left/Right): %.0f%%\n", s.MusicVolume*100)
	soundVolume := fmt.Sprintf("Sound volume (Alt+ Left/Right): %.0f%%\n", s.SoundVolume*100)
	offset := fmt.Sprintf("Offset (Shift+ Left/Right): %dms\n", s.Offset)

	exit := "Press ESC to back to choose a song.\n"
	pause := "Press TAB to pause.\n"

	ebitenutil.DebugPrint(screen.Image, fps+tps+time+"\n"+
		score+combo+flow+acc+judgmentCount+"\n"+
		speedScale+exposureTime+"\n"+
		musicVolume+soundVolume+offset+"\n"+
		exit+pause,
	)
}
