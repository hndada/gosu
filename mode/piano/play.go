package piano

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hndada/gosu"
	"github.com/hndada/gosu/ctrl"
	"github.com/hndada/gosu/draws"
	"github.com/hndada/gosu/format/osr"
	"github.com/hndada/gosu/input"
)

// ScenePlay: struct, PlayScene: function
type ScenePlay struct {
	gosu.Timer
	Chart *Chart

	gosu.MusicPlayer
	gosu.EffectPlayer
	gosu.KeyLogger

	*gosu.TransPoint
	SpeedHandler ctrl.F64Handler
	// MainBPM              float64
	// NormalizedSpeed float64
	Speed  float64
	Cursor float64

	gosu.Result
	Staged []*Note
	Flow   float64
	Combo  int
	// NoteWeights is a sum of weight of marked notes.
	// This is also max value of each score sum can get at the time.
	NoteWeights float64

	Skin             // The skin may be applied some custom settings: on/off some sprites
	BackgroundDrawer gosu.BackgroundDrawer
	StageDrawer      StageDrawer
	BarDrawer        BarDrawer
	NoteLaneDrawers  []NoteLaneDrawer
	KeyDrawer        KeyDrawer
	ScoreDrawer      gosu.ScoreDrawer
	ComboDrawer      ComboDrawer
	JudgmentDrawer   JudgmentDrawer
	MeterDrawer      gosu.MeterDrawer
}

// Todo: add Mods
func NewScenePlay(cpath string, rf *osr.Format, mvh, evh, sh ctrl.F64Handler) (scene gosu.Scene, err error) {
	s := new(ScenePlay)
	s.Chart, err = NewChart(cpath)
	if err != nil {
		return
	}
	c := s.Chart
	keyCount := c.KeyCount & ScratchMask

	if rf != nil {
		s.SetTicks(rf.BufferTime(), c.Duration())
	} else {
		s.SetTicks(-1800, c.Duration())
	}
	if path, ok := c.MusicPath(cpath); ok {
		s.MusicPlayer, err = gosu.NewMusicPlayer(mvh, path)
		if err != nil {
			return
		}
	}
	s.EffectPlayer = gosu.NewEffectPlayer(evh)
	for _, n := range c.Notes {
		if path, ok := n.SamplePath(cpath); ok {
			_ = s.Effects.Register(path)
		}
	}
	s.KeyLogger = gosu.NewKeyLogger(KeySettings[keyCount])
	if rf != nil {
		s.KeyLogger.FetchPressed = gosu.NewReplayListener(rf, keyCount, s.Time())
	}

	s.TransPoint = c.TransPoints[0]
	s.SpeedHandler = sh
	s.Speed = 1
	s.Cursor = float64(s.Time()) * s.Speed
	s.SetSpeed()

	s.Result.MD5, err = gosu.MD5(cpath)
	if err != nil {
		return
	}
	s.Result.JudgmentCounts = make([]int, len(Judgments))
	s.Result.FlowMarks = make([]float64, 0, c.Duration()/1000)
	for _, n := range c.Notes {
		s.MaxNoteWeights += n.Weight()
	}
	s.Staged = make([]*Note, keyCount)
	for k := range s.Staged {
		for _, n := range c.Notes {
			if k == n.Key {
				s.Staged[n.Key] = n
				break
			}
		}
	}
	s.Flow = 1

	s.Skin = Skins[keyCount]
	s.BackgroundDrawer = gosu.BackgroundDrawer{
		Sprite:  gosu.DefaultBackground,
		Dimness: &gosu.BackgroundDimness,
	}
	if bg := gosu.NewBackground(c.BackgroundPath(cpath)); bg.IsValid() {
		s.BackgroundDrawer.Sprite = bg
	}
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
			Cursor:   s.Cursor,
			Farthest: s.Staged[k],
			Nearest:  s.Staged[k],
		}
	}
	s.BarDrawer = BarDrawer{
		Sprite:   s.BarSprite,
		Cursor:   s.Cursor,
		Farthest: s.Chart.Bars[0],
		Nearest:  s.Chart.Bars[0],
	}
	s.KeyDrawer = NewKeyDrawer(s.KeyUpSprites, s.KeyDownSprites)
	s.ScoreDrawer = gosu.NewScoreDrawer()
	s.ComboDrawer = NewComboDrawer(s.ComboSprites)
	s.ComboDrawer.Sprites = s.ComboSprites
	s.JudgmentDrawer = NewJudgmentDrawer()
	s.MeterDrawer = gosu.NewMeterDrawer(Judgments, JudgmentColors)

	title := fmt.Sprintf("gosu - %s - [%s]", c.MusicName, c.ChartName)
	ebiten.SetWindowTitle(title)
	return s, nil
}

// Farther note has larger position. Tail's Position is always larger than Head's.
// Need to re-calculate positions when Speed has changed.
func (s *ScenePlay) SetSpeed() {
	old := s.Speed
	new := *s.SpeedHandler.Target
	s.Cursor *= new / old
	for _, tp := range s.Chart.TransPoints {
		tp.Position *= new / old
	}
	for _, n := range s.Chart.Notes {
		n.Position *= new / old
	}
	for _, b := range s.Chart.Bars {
		b.Position *= new / old
	}
	s.Speed = new
}

// Todo: apply other values of TransPoint (Volume has finished so far)
// Todo: keep playing music when making SceneResult
func (s *ScenePlay) Update() any {
	defer s.Ticker()
	if s.IsDone() {
		s.MusicPlayer.Close()
		return gosu.PlayToResultArgs{
			Result: s.Result,
		}
	}
	if s.Tick == 0 {
		s.MusicPlayer.Play()
	}

	var worst gosu.Judgment
	s.LastPressed = s.Pressed
	s.Pressed = s.FetchPressed()
	s.KeyDrawer.Update(s.LastPressed, s.Pressed)
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
				s.Effects.PlayWithVolume(name, vol)
			}
		}
		td := n.Time - s.Time() // Time difference. A negative value infers late hit
		if n.Marked {
			if n.Type != Tail {
				return fmt.Errorf("non-Tail note has not flushed")
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
			var colorType int = 0
			if n.Type == Tail {
				colorType = 1
			}
			s.MeterDrawer.AddMark(int(td), colorType)
		}
	}
	s.SetCursor()
	s.BarDrawer.Update(s.Cursor)
	for k, d := range s.NoteLaneDrawers {
		d.Update(s.Cursor)
		if k == 0 {
			fmt.Printf("Key %d: %.2f %.2f (cursor: %.2f)\n", k, d.Nearest.Position, d.Farthest.Position, d.Cursor)
		}
	}
	s.ScoreDrawer.Update(s.Score())
	s.ComboDrawer.Update(s.Combo)
	s.JudgmentDrawer.Update(worst)
	s.MeterDrawer.Update()

	// Changed speed should be applied after positions are calculated.
	s.TransPoint.FetchByTime(s.Time())
	if fired := s.SpeedHandler.Update(); fired {
		s.SetSpeed()
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
	s.KeyDrawer.Draw(screen)
	s.ScoreDrawer.Draw(screen)
	s.ComboDrawer.Draw(screen)
	s.JudgmentDrawer.Draw(screen)
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
		s.Speed*100, *s.SpeedHandler.Target*100, ExposureTime(s.Speed)))
}
func (s ScenePlay) Time() int64 { return s.Timer.Time() }

// func (s ScenePlay) Speed() float64 { return *s.SpeedHandler.Target * s.TransPoint.Speed() }

func (s ScenePlay) CurrentSpeed() float64 { return s.TransPoint.Speed * s.Speed }

// Supposes one current TransPoint can increment cursor precisely.
func (s *ScenePlay) SetCursor() {
	tp := s.TransPoint
	s.Cursor = tp.Position + float64(s.Time()-tp.Time)*s.CurrentSpeed()
}
