package piano

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hndada/gosu"
	"github.com/hndada/gosu/draws"
	"github.com/hndada/gosu/format/osr"
	"github.com/hndada/gosu/input"
)

// ScenePlay: struct, PlayScene: function
// NoteWeights is a sum of weight of marked notes.
// This is also max value of each score sum can get at the time.
type ScenePlay struct {
	*gosu.BaseScenePlay
	// Chart  *Chart
	// Staged []*gosu.Note
	// Leadings        []*gosu.Note
	Skin           // The skin may be applied some custom settings: on/off some sprites
	NoteDrawers    []gosu.FixedLaneDrawer
	BarDrawer      gosu.FixedLaneDrawer
	JudgmentDrawer JudgmentDrawer
	ComboDrawer    draws.NumberDrawer
	KeyDrawer      KeyDrawer
}

// Todo: Let users change speed during playing
// Todo: add Mods to input param
func NewScenePlay(cpath string,
	mode, subMode int, mods gosu.Mods,
	rf *osr.Format,
	speedScale *float64,
	keySettings []input.Key,
	bg draws.Sprite) (gosu.Scene, error) {
	s := new(ScenePlay)
	var err error
	s.BaseScenePlay, err = gosu.NewBaseScenePlay(
		cpath, mode, subMode, mods,
		rf, speedScale, keySettings, bg,
	)
	if err != nil {
		return s, err
	}
	s.Skin = Skins[subMode]
	for _, n := range s.Chart.Notes {
		s.MaxNoteWeights += Weight(*n)
	}
	s.NoteDrawers = make([]gosu.FixedLaneDrawer, subMode)
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
	// s.KeyDrawer.Sprites[0] = s.KeyUpSprites
	// s.KeyDrawer.Sprites[1] = s.KeyDownSprites
	// s.KeyDrawer.KeyDownCountdowns = make([]int, c.KeyCount)
	s.JudgmentCounts = make([]int, 5)
	s.JudgmentDrawer.Sprites = s.JudgmentSprites
	s.ComboDrawer.Sprites = s.ComboSprites
	s.Meter = gosu.NewMeter(Judgments, JudgmentColors)
	return s, nil
}

// TPS affects only on Update(), not on Draw().
// Todo: apply other values of TransPoint (Volume has finished so far)
func (s *ScenePlay) Update() any {
	defer s.Ticker()
	if args := s.CheckFinished(); args != nil {
		return args
	}
	s.PlayMusicIfStarted()

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
	// Speed, BPM, Volume and Highlight
	s.UpdateTransPoint()
	if fired := s.SpeedHandler.Update(); fired {
		for k := range s.NoteDrawers {
			s.NoteDrawers[k].SetSpeedScale(*s.SpeedScale)
		}
		go s.Chart.SetSpeedScale(*s.SpeedScale)
	}
	s.BarDrawer.Update(s.TransPoint.Speed()/s.MainBPM, *s.SpeedScale)

	// s.JudgmentDrawer.Update(worst)
	s.ComboDrawer.Update(s.Combo, 0)
	s.DelayedScore.Set(s.Score())
	s.DelayedScore.Update()
	s.ScoreDrawer.Update(int(s.DelayedScore.Delayed), 0)
	// fmt.Println(s.Combo, s.DelayedScore.Delayed)
	s.Meter.Update(marks)
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
	s.Meter.Draw(screen)
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
		ebiten.CurrentFPS(), ebiten.CurrentTPS(), float64(s.Time())/1000, float64(s.Chart.Duration)/1000,
		s.Score(), s.ScoreBound(), s.Flow*100, s.Combo,
		fr*100, ar*100, rr*100, s.JudgmentCounts,
		s.Speed()*100, *s.SpeedScale*100, ExposureTime(s.Speed())))
}
