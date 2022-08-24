package drum

import (
	"fmt"
	"path/filepath"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hndada/gosu"
	"github.com/hndada/gosu/ctrl"
	"github.com/hndada/gosu/format/osr"
	"github.com/hndada/gosu/input"
)

// Only one note is staged in Drum mode.
type ScenePlay struct {
	gosu.BaseScenePlay

	// General
	Chart *Chart
	Skin  // The skin may be applied some custom settings: on/off some sprites

	// Speed, BPM, Volume and Highlight
	// Audio
	// Input
	// Notes
	PlayNotes   []*PlayNote
	StagedNote  *PlayNote
	LeadingTail *PlayNote // For drawing Roll body efficiently.

	// Scores
	JudgmentDrawer JudgmentDrawer
	ComboDrawer    ComboDrawer
	KeyDrawer      KeyDrawer

	NoteOverlayDrawer   NoteOverlayDrawer
	RollTickDrawer      RollTickDrawer
	RollTickComboDrawer RollTickComboDrawer
	ShakeDrawer         ShakeDrawer
}

func NewScenePlay(cpath string, mods gosu.Mods, rf *osr.Format) (gosu.Scene, error) {
	s := new(ScenePlay)

	// General
	waitBefore := s.SetTick(rf)
	s.MD5 = gosu.MD5(cpath)
	c, err := NewChart(cpath, mods) // NewChart must be at first.
	if err != nil {
		return nil, err
	}
	s.Chart = c
	s.EndTime = c.Duration + gosu.DefaultWaitAfter
	// General: Graphics
	s.SetWindowTitle(c.BaseChart)
	s.Skin = DefaultSkin
	s.SetBackground(c.BackgroundPath(cpath))

	// Speed, BPM, Volume and Highlight
	s.MainBPM, _, _ = gosu.BPMs(c.TransPoints, c.Duration) // Todo: Need a test
	s.SpeedBase = SpeedBase
	s.SetInitTransPoint(c.TransPoints[0])

	// Audio
	apath := filepath.Join(filepath.Dir(cpath), c.AudioFilename)
	err = s.SetMusicPlayer(apath)
	if err != nil {
		return s, err
	}
	seNames := make([]string, 0, len(c.Notes))
	for _, n := range c.Notes {
		if name := n.SampleFilename; name != "" {
			seNames = append(seNames, name)
		}
	}
	s.SetSoundMap(cpath, seNames)

	// Input
	if rf != nil {
		s.FetchPressed = gosu.NewReplayListener(rf, 4, waitBefore)
	} else {
		s.FetchPressed = input.NewListener(KeySettings[:])
	}
	s.LastPressed = make([]bool, 4)
	s.Pressed = make([]bool, 4)

	// Note
	s.PlayNotes, s.StagedNote, s.LeadingTail, s.MaxNoteWeights = NewPlayNotes(c)
	et, wb, wa := s.EndTime, waitBefore, gosu.DefaultWaitAfter
	s.BarLineDrawer.Times = gosu.BarLineTimes(c.TransPoints, et, wb, wa)
	// s.BarLineDrawer.Offset = NoteHeigth / 2
	s.BarLineDrawer.Sprite = s.BarLineSprite
	s.BarLineDrawer.Horizontal = true
	s.KeyDrawer.Sprites = s.KeySprites

	// Score
	s.JudgmentCounts = make([]int, 5)
	s.FlowMarks = make([]float64, 0, c.Duration)
	s.Flow = 1
	s.ScoreDrawer.DelayedScore.Mode = ctrl.DelayedModeExp
	s.ScoreDrawer.Sprites = s.ScoreSprites
	s.JudgmentDrawer.Sprites = s.JudgmentSprites
	s.ComboDrawer.Sprites = s.ComboSprites
	s.TimingMeter = gosu.NewTimingMeter(Judgments[0][:], JudgmentColors)
	return s, nil
}

// TPS affects only on Update(), not on Draw().
// Todo: apply other values of TransPoint (Volume has finished so far)
func (s *ScenePlay) Update() any {
	// General
	s.Tick++
	done := ebiten.IsKeyPressed(ebiten.KeyEscape) || s.Time() >= s.EndTime
	if done { // Todo: keep playing music when making SceneResult
		if s.MusicPlayer != nil {
			s.MusicPlayer.Close()
		}
		return gosu.PlayToResultArgs{
			Result: s.Result,
		}
	}

	// Speed, BPM, Volume and Highlight
	s.UpdateTransPoint()

	// Audio
	if s.Tick == 0 && s.MusicPlayer != nil {
		s.MusicPlayer.SetVolume(s.Volume)
		s.MusicPlayer.Play()
	}

	// Input
	s.LastPressed = s.Pressed
	s.Pressed = s.FetchPressed()
	s.KeyDrawer.Update(s.LastPressed, s.Pressed)

	// Notes and Scores
	var worst gosu.Judgment
	marks := make([]gosu.TimingMeterMark, 0, 7)
	if s.StagedNote != nil {
		// if n.Type != Tail && s.KeyAction(k) == input.Hit {
		// 	if name := n.SampleFilename; name != "" {
		// 		vol := n.SampleVolume
		// 		if vol == 0 {
		// 			vol = s.TransPoint.Volume
		// 		}
		// 		s.Sounds.PlayWithVolume(name, vol)
		// 	}
		// }
		n := s.StagedNote
		td := n.Time - s.Time() // Time difference. A negative value infers late hit
		if n.Marked {
			if n.Type != Tail {
				return fmt.Errorf("non-tail note has not flushed")
			}
			if td < Miss.Window { // Keep Tail staged until near ends.
				s.StagedNote = n.Next
			}
		}
		// } else if j := Verdict(n.Type, s.KeyAction(n.Key), td); j.Window != 0 {
		// 	s.MarkNote(n, j)
		// 	if worst.Window < j.Window {
		// 		worst = j
		// 	}

		// 	clr := white
		// 	if n.Type == BigDon || n.Type == BigKat {
		// 		clr = purple
		// 	}
		// 	marks = append(marks, gosu.TimingMeterMark{
		// 		Countdown: gosu.TimingMeterMarkDuration,
		// 		TimeDiff:  td,
		// 		Color:     clr,
		// 	})
		// }
	}
	s.JudgmentDrawer.Update(worst)
	s.ComboDrawer.Update(s.Combo)
	// s.ScoreDrawer.Update(s.Score())
	s.TimingMeter.Update(marks)
	return nil
}
func (s ScenePlay) Draw(screen *ebiten.Image) {
	s.BackgroundDrawer.Draw(screen)
	s.FieldSprite.Draw(screen)
	s.HintSprite.Draw(screen)
	// s.BarLineDrawer.Draw(screen, n.Position)
	s.DrawRollBodies(screen)
	s.DrawNotes(screen)
	s.KeyDrawer.Draw(screen)
	s.TimingMeter.Draw(screen)
	s.ComboDrawer.Draw(screen)
	s.JudgmentDrawer.Draw(screen)
	s.ScoreDrawer.Draw(screen)
	// s.DebugPrint(screen)
}

// func (s ScenePlay) DebugPrint(screen *ebiten.Image) {
// 	var fr, ar, rr float64 = 1, 1, 1
// 	if s.NoteWeights > 0 {
// 		fr = s.Flows / s.NoteWeights
// 		ar = s.Accs / s.NoteWeights
// 		rr = s.Extras / s.NoteWeights
// 	}

// 	ebitenutil.DebugPrint(screen, fmt.Sprintf(
// 		"CurrentFPS: %.2f\nCurrentTPS: %.2f\nTime: %.3fs/%.0fs\n\n"+
// 			"Score: %.0f | %.0f \nFlow: %.0f/100\nCombo: %d\n\n"+
// 			"Flow rate: %.2f%%\nAccuracy: %.2f%%\n(Kool: %.2f%%)\nJudgment counts: %v\n\n"+
// 			"Speed: %.0f | %.0f\n(Exposure time: %.fms)\n\n",
// 		ebiten.CurrentFPS(), ebiten.CurrentTPS(), float64(s.Time())/1000, float64(s.Chart.Duration)/1000,
// 		s.Score(), s.ScoreBound(), s.Flow*100, s.Combo,
// 		fr*100, ar*100, rr*100, s.JudgmentCounts,
// 		s.Speed()*100, s.SpeedBase*100, ExposureTime(s.Speed())))
// }
