package piano

import (
	"fmt"
	"path/filepath"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hndada/gosu"
	"github.com/hndada/gosu/ctrl"
	"github.com/hndada/gosu/draws"
	"github.com/hndada/gosu/format/osr"
	"github.com/hndada/gosu/input"
)

// ScenePlay: struct, PlayScene: function
// NoteWeights is a sum of weight of marked notes.
// This is also max value of each score sum can get at the time.
type ScenePlay struct {
	Skin // The skin may be applied some custom settings: on/off some sprites
	gosu.BaseScenePlay
	Chart           *Chart
	Staged          []*gosu.Note
	Leadings        []*gosu.Note
	NoteLaneDrawers []gosu.NoteLaneDrawer
	// BarDrawer       gosu.NoteLaneDrawer
	JudgmentDrawer JudgmentDrawer
	ComboDrawer    draws.NumberDrawer
	KeyDrawer      KeyDrawer
}

// Todo: Let users change speed during playing
// Todo: add Mods to input param
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
	s.Skin = Skins[c.KeyCount]
	s.SetBackground(c.BackgroundPath(cpath))

	// Speed, BPM, Volume and Highlight
	s.MainBPM, _, _ = gosu.BPMs(c.TransPoints, c.Duration) // Todo: Need a test
	s.SpeedBase = SpeedBase
	s.SetInitTransPoint(c.TransPoints[0])
	s.SpeedHandler = gosu.NewSpeedHandler(&s.SpeedBase)

	// Audio
	apath := filepath.Join(filepath.Dir(cpath), c.AudioFilename)
	err = s.SetMusicPlayer(apath)
	if err != nil {
		return s, err
	}
	seNames := make([]string, 0, len(c.Notes))
	for _, n := range c.Notes {
		if name := n.SampleName; name != "" {
			seNames = append(seNames, name)
		}
	}
	s.SetSoundMap(cpath, seNames)

	// Input
	if rf != nil {
		s.FetchPressed = gosu.NewReplayListener(rf, c.KeyCount, waitBefore)
	} else {
		s.FetchPressed = input.NewListener(KeySettings[c.KeyCount])
	}
	s.LastPressed = make([]bool, c.KeyCount)
	s.Pressed = make([]bool, c.KeyCount)

	// Note
	prevs := make([]*gosu.Note, c.KeyCount)
	s.Staged = make([]*gosu.Note, s.Chart.KeyCount)

	for i := range c.Notes {
		n := c.Notes[i]
		prev := prevs[n.Key]
		c.Notes[i].Prev = prev
		if prev != nil { // Next value is set later.
			prev.Next = n
		}
		prevs[n.Key] = n
		if s.Staged[n.Key] == nil {
			s.Staged[n.Key] = n
		}
		s.MaxNoteWeights += Weight(*n)
	}
	s.Leadings = make([]*gosu.Note, s.Chart.KeyCount)
	copy(s.Leadings, s.Staged)
	// s.PlayNotes, s.Staged, s.LowestTails, s.MaxNoteWeights = NewPlayNotes(c)
	// Note: Graphics
	s.BarDrawer.Sprite = s.BarLineSprite
	// et, wb, wa := s.EndTime, waitBefore, gosu.DefaultWaitAfter

	// times := gosu.BarTimes(c.TransPoints, s.EndTime, wb, wa)
	// s.BarDrawer.Times =
	// s.BarDrawer.Offset = NoteHeigth / 2
	s.Chart.SetPositions(s.SpeedBase)
	s.BarDrawer.Bars = s.Chart.Bars
	s.KeyDrawer.Sprites[0] = s.KeyUpSprites
	s.KeyDrawer.Sprites[1] = s.KeyDownSprites
	s.KeyDrawer.KeyDownCountdowns = make([]int, c.KeyCount)

	// Score
	s.JudgmentCounts = make([]int, 5)
	s.FlowMarks = make([]float64, 0, c.Duration)
	s.Flow = 1
	// Score: Graphics
	s.DelayedScore.Mode = ctrl.DelayedModeExp
	// s.ScoreDrawer.DelayedScore.Mode = ctrl.DelayedModeExp
	s.ScoreDrawer.Sprites = gosu.ScoreSprites
	s.JudgmentDrawer.Sprites = s.JudgmentSprites
	s.ComboDrawer.Sprites = s.ComboSprites
	s.Meter = gosu.NewMeter(Judgments, JudgmentColors)
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
	if fired := s.SpeedHandler.Update(); fired {
		go s.Chart.SetPositions(s.SpeedBase)
	}

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
	marks := make([]gosu.MeterMark, 0, 7)
	for k, n := range s.Staged {
		if n == nil {
			continue
		}
		if n.Type != gosu.Tail && s.KeyAction(k) == input.Hit {
			if name := n.SampleName; name != "" {
				vol := n.SampleVolume
				if vol == 0 {
					vol = s.TransPoint.Volume
				}
				s.Sounds.PlayWithVolume(name, vol)
			}
		}
		td := n.Time - s.Time() // Time difference. A negative value infers late hit
		if n.Marked {
			if n.Type != gosu.Tail {
				return fmt.Errorf("non-tail note has not flushed")
			}
			if td < Miss.Window { // Keep Tail staged until near ends.
				s.Staged[n.Key] = n.Next
			}
			continue
		}
		if j := Verdict(n.Type, s.KeyAction(n.Key), td); j.Window != 0 {
			s.MarkNote(n, j)
			if worst.Window < j.Window {
				worst = j
			}
			mark := gosu.MeterMark{
				Countdown: s.Meter.MaxCountdown,
				Offset:    int(td),
				ColorType: 0,
			}
			if n.Type == gosu.Tail {
				mark.ColorType = 1
			}
			marks = append(marks, mark)
		}
	}
	s.JudgmentDrawer.Update(worst)
	s.ComboDrawer.Update(s.Combo, 0)
	s.DelayedScore.Set(s.Score())
	s.DelayedScore.Update()
	s.ScoreDrawer.Update(int(s.DelayedScore.Delayed), 0)
	s.Meter.Update(marks)
	return nil
}
func (s ScenePlay) Draw(screen *ebiten.Image) {
	s.BackgroundDrawer.Draw(screen)
	s.FieldSprite.Draw(screen, nil)
	s.HintSprite.Draw(screen, nil)
	s.BarDrawer.Draw(screen)
	// s.DrawLongNoteBodies(screen)
	// s.DrawNotes(screen)
	for _, d := range s.NoteLaneDrawers {
		d.Draw(screen)
	}
	s.KeyDrawer.Draw(screen, s.Pressed)
	s.Meter.Draw(screen)
	s.ComboDrawer.Draw(screen)
	s.JudgmentDrawer.Draw(screen)
	s.ScoreDrawer.Draw(screen)
	s.DebugPrint(screen)
}

func (s ScenePlay) DebugPrint(screen *ebiten.Image) {
	var fr, ar, rr float64 = 1, 1, 1
	if s.NoteWeights > 0 {
		fr = s.Flows / s.NoteWeights
		ar = s.Accs / s.NoteWeights
		rr = s.Extras / s.NoteWeights
	}

	ebitenutil.DebugPrint(screen, fmt.Sprintf(
		"CurrentFPS: %.2f\nCurrentTPS: %.2f\nTime: %.3fs/%.0fs\n\n"+
			"Score: %.0f | %.0f \nFlow: %.0f/100\nCombo: %d\n\n"+
			"Flow rate: %.2f%%\nAccuracy: %.2f%%\n(Kool: %.2f%%)\nJudgment counts: %v\n\n"+
			"Speed: %.0f | %.0f\n(Exposure time: %.fms)\n\n",
		ebiten.CurrentFPS(), ebiten.CurrentTPS(), float64(s.Time())/1000, float64(s.Chart.Duration)/1000,
		s.Score(), s.ScoreBound(), s.Flow*100, s.Combo,
		fr*100, ar*100, rr*100, s.JudgmentCounts,
		s.Speed()*100, s.SpeedBase*100, ExposureTime(s.Speed())))
}
