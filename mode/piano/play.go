package piano

import (
	"fmt"
	"path/filepath"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hndada/gosu"
	"github.com/hndada/gosu/audios"
	"github.com/hndada/gosu/ctrl"
	"github.com/hndada/gosu/draws"
	"github.com/hndada/gosu/format/osr"
	"github.com/hndada/gosu/input"
)

// ScenePlay: struct, PlayScene: function
// NoteWeights is a sum of weight of marked notes.
// This is also max value of each score sum can get at the time.
type ScenePlay struct {
	*gosu.BaseScenePlay
	Chart  *Chart
	Staged []*Note
	// Leadings       []*gosu.Note
	Skin           // The skin may be applied some custom settings: on/off some sprites
	NoteDrawers    []gosu.FixedLaneDrawer
	BarDrawer      gosu.FixedLaneDrawer
	JudgmentDrawer JudgmentDrawer
	ComboDrawer    draws.NumberDrawer
	KeyDrawer      KeyDrawer
}

// Todo: Let users change speed during playing
// Todo: add Mods to input param
func NewScenePlay(cpath string, rf *osr.Format, vh, sh ctrl.F64Handler) (scene gosu.Scene, err error) {
	s := new(ScenePlay)
	s.Chart, err = NewChart(cpath)
	if err != nil {
		return
	}
	c := s.Chart
	keyCount := c.KeyCount & ScratchMask
	s.Tick = gosu.TimeToTick(gosu.WaitBefore)
	s.MaxTick = gosu.TimeToTick(c.Duration() + gosu.WaitAfter)

	s.Skin = Skins[keyCount]
	s.VolumeHandler = vh
	if gosu.IsAudioExisted(c.AudioFilename) {
		s.MusicPlayer, s.MusicCloser, err = audios.NewPlayer(c.AudioFilename)
		if err != nil {
			return
		}
		s.MusicPlayer.SetVolume(*vh.Target)
	}
	s.Sounds = audios.NewSoundMap(vh.Target)
	for _, n := range c.Notes {
		if n.SampleName == "" {
			continue
		}
		path := filepath.Join(filepath.Dir(cpath), n.SampleName)
		s.Sounds.Register(path)
	}

	s.SpeedHandler = sh
	s.TransPoint = c.TransPoints[0].FetchLatest()

	if rf != nil {
		bufferTime := rf.BufferTime()
		if bufferTime > gosu.WaitBefore {
			s.Tick = gosu.TimeToTick(bufferTime)
		}
		s.FetchPressed = gosu.NewReplayListener(rf, keyCount, bufferTime)
	} else {
		s.FetchPressed = input.NewListener(KeySettings[c.KeyCount])
	}
	s.LastPressed = make([]bool, keyCount)
	s.Pressed = make([]bool, keyCount)
	for k := range s.Staged {
		for _, n := range c.Notes {
			if k == n.Key {
				s.Staged[n.Key] = n
				break
			}
		}
	}

	s.MD5, err = gosu.MD5(cpath)
	if err != nil {
		return
	}
	s.Flow = 1
	s.FlowMarks = make([]float64, 0, c.Duration()/1000)
	for _, n := range c.Notes {
		s.MaxNoteWeights += n.Weight()
	}
	s.NoteDrawers = make([]gosu.FixedLaneDrawer, keyCount)
	for k := range s.NoteDrawers {
		sprites := []draws.Sprite{
			s.NoteSprites[k], s.HeadSprites[k],
			s.TailSprites[k], s.BodySprites[k],
		}
		s.NoteDrawers[k] = gosu.NewFixedLaneDrawer(
			gosu.Downward,
			sprites,
			HitPosition,
			// beatSpeed,  // Speed calculation is each mode's task.
			*s.SpeedScale,
			draws.Grayer,
			gosu.TickToTime(s.Tick),
			s.TransPoint,
			s.Staged[k].LaneObject,
		)
	}
	s.BarDrawer = gosu.NewFixedLaneDrawer(
		gosu.Downward,
		[]draws.Sprite{s.BarLineSprite},
		HitPosition,
		// beatSpeed,  // Speed calculation is each mode's task.
		*s.SpeedScale,
		nil,
		gosu.TickToTime(s.Tick),
		s.TransPoint,
		s.Chart.Bars[0],
	)
	s.BackgroundDrawer = gosu.NewBackgroundDrawer(
		c.BackgroundPath(cpath), &gosu.BackgroundDimness,
	)
	s.DelayedScore.Mode = ctrl.DelayedModeExp
	s.ScoreDrawer = gosu.NewScoreDrawer()
	// s.KeyDrawer.Sprites[0] = s.KeyUpSprites
	// s.KeyDrawer.Sprites[1] = s.KeyDownSprites
	// s.KeyDrawer.KeyDownCountdowns = make([]int, c.KeyCount)
	s.JudgmentCounts = make([]int, 5)
	s.JudgmentDrawer.Sprites = s.JudgmentSprites
	s.ComboDrawer.Sprites = s.ComboSprites
	s.MeterDrawer = gosu.NewMeterDrawer(Judgments, JudgmentColors)
	title := fmt.Sprintf("gosu - %s - [%s]", c.MusicName, c.ChartName)
	ebiten.SetWindowTitle(title)
	return s, nil
}

// TPS affects only on Update(), not on Draw().
// Todo: apply other values of TransPoint (Volume has finished so far)
func (s *ScenePlay) Update() any {
	defer s.Ticker()
	if s.IsDone() {
		// Todo: keep playing music when making SceneResult
		if s.MusicPlayer != nil {
			s.MusicPlayer.Close()
		}
		return gosu.PlayToResultArgs{
			Result: s.Result,
		}
	}
	if s.Tick == 0 && s.MusicPlayer != nil {
		s.MusicPlayer.Play()
	}

	s.LastPressed = s.Pressed
	s.Pressed = s.FetchPressed()
	s.KeyDrawer.Update(s.LastPressed, s.Pressed)

	var worst gosu.Judgment
	marks := make([]gosu.MeterMark, 0, 7)
	for k, n := range s.Staged {
		if n == nil {
			continue
		}
		if n.Type != Tail && s.KeyAction(k) == input.Hit {
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
			if n.Type != Tail {
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
				Countdown: s.MeterDrawer.MaxCountdown,
				Offset:    int(td),
				ColorType: 0,
			}
			if n.Type == Tail {
				mark.ColorType = 1
			}
			marks = append(marks, mark)
		}
	}
	// Speed, BPM, Volume and Highlight
	s.UpdateTransPoint()
	if fired := s.SpeedHandler.Update(); fired {
		speedScale := *s.SpeedHandler.Target
		for k := range s.NoteDrawers {
			s.NoteDrawers[k].SetSpeedScale(speedScale)
		}
		go s.Chart.SetSpeedScale(speedScale)
		s.BarDrawer.SetSpeedScale(speedScale)
	}
	s.BarDrawer.Update(s.TransPoint.Speed(), speedScale)

	// s.JudgmentDrawer.Update(worst)
	s.ComboDrawer.Update(s.Combo, 0)
	s.DelayedScore.Set(s.Score())
	s.DelayedScore.Update()
	s.ScoreDrawer.Update(int(s.DelayedScore.Delayed), 0)
	s.MeterDrawer.Update(marks)
	return nil
}
func (s ScenePlay) Draw(screen *ebiten.Image) {
	s.BackgroundDrawer.Draw(screen)
	s.FieldSprite.Draw(screen, nil)
	s.HintSprite.Draw(screen, nil)
	s.BarDrawer.Draw(screen)
	for _, d := range s.NoteDrawers {
		d.Draw(screen)
	}
	// s.KeyDrawer.Draw(screen, s.Pressed)
	s.MeterDrawer.Draw(screen)
	s.ComboDrawer.Draw(screen)
	// s.JudgmentDrawer.Draw(screen)
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
		ebiten.CurrentFPS(), ebiten.CurrentTPS(), float64(s.Time())/1000, float64(s.Chart.Duration())/1000,
		s.Score(), s.ScoreBound(), s.Flow*100, s.Combo,
		fr*100, ar*100, rr*100, s.JudgmentCounts,
		s.Speed()*100, *s.SpeedHandler.Target*100, ExposureTime(s.Speed())))
}
