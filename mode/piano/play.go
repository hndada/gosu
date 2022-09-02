package piano

import (
	"fmt"

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
type ScenePlay struct {
	*gosu.BaseScenePlay
	Chart  *Chart
	Staged []*Note

	Skin                 // The skin may be applied some custom settings: on/off some sprites
	MainBPM              float64
	NormalizedSpeedScale float64
	Cursor               float64
	NoteLaneDrawers      []NoteLaneDrawer
	BarDrawer            BarDrawer

	StageDrawer
	JudgmentDrawer JudgmentDrawer
	ComboDrawer    draws.NumberDrawer
	KeyDrawer      KeyDrawer
}

// Todo: add Mods
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

	s.VolumeHandler = vh
	if path, ok := c.MusicPath(cpath); ok {
		s.MusicPlayer, s.MusicCloser, err = audios.NewPlayer(path)
		if err != nil {
			return
		}
		s.MusicPlayer.SetVolume(*vh.Target)
	}
	s.Sounds = audios.NewSoundMap(vh.Target)
	for _, n := range c.Notes {
		if path, ok := n.SamplePath(cpath); ok {
			_ = s.Sounds.Register(path)
		}
	}

	if rf != nil {
		bufferTime := rf.BufferTime()
		if bufferTime > gosu.WaitBefore {
			s.Tick = gosu.TimeToTick(bufferTime)
		}
		s.FetchPressed = gosu.NewReplayListener(rf, keyCount, bufferTime)
	} else {
		s.FetchPressed = input.NewListener(KeySettings[keyCount])
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

	s.SpeedHandler = sh
	s.TransPoint = c.TransPoints[0].FetchLatest()
	s.MainBPM, _, _ = c.BPMs()
	s.Cursor = float64(gosu.TickToTime(s.Tick)) * s.TransPoint.Speed()
	s.NormalizedSpeedScale = 1
	s.SetSpeedScale()

	s.Result.MD5, err = gosu.MD5(cpath)
	if err != nil {
		return
	}
	s.Result.JudgmentCounts = make([]int, 5)
	s.Result.FlowMarks = make([]float64, 0, c.Duration()/1000)
	for _, n := range c.Notes {
		s.MaxNoteWeights += n.Weight()
	}
	s.Flow = 1

	s.Skin = Skins[keyCount]
	s.StageDrawer = StageDrawer{
		Field: s.FieldSprite,
		Hint:  s.HintSprite,
	}
	s.NoteLaneDrawers = make([]NoteLaneDrawer, keyCount)
	for k := range s.NoteLaneDrawers {
		s.NoteLaneDrawers[k] = NoteLaneDrawer{
			Sprites: [4]draws.Sprite{
				s.NoteSprites[k], s.HeadSprites[k],
				s.TailSprites[k], s.BodySprites[k],
			},
			Cursor:   &s.Cursor,
			Farthest: s.Staged[k],
			Nearest:  s.Staged[k],
		}
	}
	s.BarDrawer = BarDrawer{
		Sprite:   s.BarSprite,
		Cursor:   &s.Cursor,
		Farthest: s.Chart.Bars[0],
		Nearest:  s.Chart.Bars[0],
	}
	s.BackgroundDrawer = gosu.BackgroundDrawer{
		Sprite:  gosu.DefaultBackground,
		Dimness: &gosu.BackgroundDimness,
	}
	if bg := gosu.NewBackground(c.BackgroundPath(cpath)); bg.IsValid() {
		s.BackgroundDrawer.Sprite = bg
	}

	// s.KeyDrawer.Sprites[0] = s.KeyUpSprites
	// s.KeyDrawer.Sprites[1] = s.KeyDownSprites
	// s.KeyDrawer.KeyDownCountdowns = make([]int, c.KeyCount)

	s.ScoreDrawer = gosu.NewScoreDrawer()
	s.ComboDrawer.Sprites = s.ComboSprites
	s.JudgmentDrawer.Sprites = s.JudgmentSprites
	s.MeterDrawer = gosu.NewMeterDrawer(Judgments, JudgmentColors)

	title := fmt.Sprintf("gosu - %s - [%s]", c.MusicName, c.ChartName)
	ebiten.SetWindowTitle(title)
	return s, nil
}

// Farther note has larger position. Tail's Position is always larger than Head's.
// Need to re-calculate positions when SpeedScale has changed.
func (s *ScenePlay) SetSpeedScale() {
	normalized := *s.SpeedHandler.Target / s.MainBPM

	s.Cursor *= normalized / s.NormalizedSpeedScale
	for _, tp := range s.Chart.TransPoints {
		tp.Position *= normalized / s.NormalizedSpeedScale
	}
	for _, n := range s.Chart.Notes {
		n.Position *= normalized / s.NormalizedSpeedScale
	}
	for _, b := range s.Chart.Bars {
		b.Position *= normalized / s.NormalizedSpeedScale
	}

	s.NormalizedSpeedScale = normalized
}

// Todo: apply other values of TransPoint (Volume has finished so far)
// Todo: keep playing music when making SceneResult
func (s *ScenePlay) Update() any {
	defer s.Ticker()
	if s.IsDone() {
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
	s.ScoreDrawer.Update(s.Score())
	s.ComboDrawer.Update(s.Combo, 0)
	// s.JudgmentDrawer.Update(worst)
	s.MeterDrawer.Update(marks)

	// Changed speed should be applied after positions are calculated.
	// Supposes one current TransPoint can increment cursor precisely.
	s.Cursor += s.TransPoint.Speed() * gosu.TimeStep
	s.UpdateTransPoint()
	if fired := s.SpeedHandler.Update(); fired {
		s.SetSpeedScale()
	}
	return nil
}
func (s ScenePlay) Draw(screen *ebiten.Image) {
	s.BackgroundDrawer.Draw(screen)
	s.StageDrawer.Draw(screen)
	s.BarDrawer.Draw(screen)
	for _, d := range s.NoteLaneDrawers {
		d.Draw(screen)
	}
	// s.KeyDrawer.Draw(screen, s.Pressed)
	s.ScoreDrawer.Draw(screen)
	s.ComboDrawer.Draw(screen)
	// s.JudgmentDrawer.Draw(screen)
	s.MeterDrawer.Draw(screen)
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
