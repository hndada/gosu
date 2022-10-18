package piano

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hndada/gosu"
	"github.com/hndada/gosu/draws"
	"github.com/hndada/gosu/format/osr"
)

// ScenePlay: struct, PlayScene: function
type ScenePlay struct {
	Chart *Chart
	gosu.Timer
	gosu.MusicPlayer
	// gosu.EffectPlayer
	gosu.KeyLogger

	*gosu.TransPoint
	SpeedScale float64
	Cursor     float64
	Staged     []*Note
	gosu.Scorer

	Skin             // The skin may be applied some custom settings: on/off some sprites
	BackgroundDrawer gosu.BackgroundDrawer
	StageDrawer      StageDrawer
	BarDrawer        BarDrawer

	NoteDrawers    []NoteDrawer
	KeyDrawer      KeyDrawer
	JudgmentDrawer JudgmentDrawer

	ScoreDrawer gosu.ScoreDrawer
	ComboDrawer gosu.NumberDrawer
	MeterDrawer gosu.MeterDrawer
}

// Todo: add Mods
func NewScenePlay(cpath string, rf *osr.Format) (scene gosu.Scene, err error) {
	s := new(ScenePlay)
	s.Chart, err = NewChart(cpath)
	if err != nil {
		return
	}
	c := s.Chart
	gosu.SetTitle(c.ChartHeader)
	keyCount := c.KeyCount & ScratchMask
	s.Timer = gosu.NewTimer(c.Duration())
	// s.SetTicks(c.Duration())
	if path, ok := c.MusicPath(cpath); ok {
		s.MusicPlayer, err = gosu.NewMusicPlayer(path, &s.Timer) //(gosu.MusicVolumeHandler, path)
		if err != nil {
			return
		}
	}
	// s.EffectPlayer = gosu.NewEffectPlayer(gosu.EffectVolumeHandler)
	// for _, n := range c.Notes {
	// 	if path, ok := n.Sample.Path(cpath); ok {
	// 		_ = s.Effects.Register(path)
	// 	}
	// }
	s.KeyLogger = gosu.NewKeyLogger(KeySettings[keyCount])
	if rf != nil {
		s.KeyLogger.FetchPressed = NewReplayListener(rf, keyCount, &s.Timer)
	}

	s.TransPoint = c.TransPoints[0]
	s.SpeedScale = 1
	s.Cursor = float64(s.Now) * s.SpeedScale
	s.SetSpeed()
	s.Scorer = gosu.NewScorer(c.ScoreFactors)
	s.JudgmentCounts = make([]int, len(Judgments))
	// s.Result.FlowMarks = make([]float64, 0, c.Duration()/1000)
	var maxWeight float64
	for _, n := range c.Notes {
		maxWeight += n.Weight()
	}
	for i := range s.MaxWeights {
		s.MaxWeights[i] = maxWeight
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

	s.Skin = Skins[keyCount]
	s.BackgroundDrawer = gosu.BackgroundDrawer{
		Brightness: &gosu.BackgroundBrightness,
		Sprite:     gosu.DefaultBackground,
	}
	if bg := gosu.NewBackground(c.BackgroundPath(cpath)); bg.IsValid() {
		s.BackgroundDrawer.Sprite = bg
	}
	s.StageDrawer = StageDrawer{
		FieldSprite: s.FieldSprite,
		HintSprite:  s.HintSprite,
	}
	s.NoteDrawers = make([]NoteDrawer, keyCount)
	for k := range s.NoteDrawers {
		s.NoteDrawers[k] = NoteDrawer{
			Cursor:   s.Cursor,
			Farthest: s.Staged[k],
			Nearest:  s.Staged[k],
			Sprites: [4]draws.Sprite{
				s.NoteSprites[k], s.HeadSprites[k],
				s.TailSprites[k], s.BodySprites[k],
			},
		}
	}
	s.BarDrawer = BarDrawer{
		Cursor:   s.Cursor,
		Farthest: c.Bars[0],
		Nearest:  c.Bars[0],
		Sprite:   s.BarSprite,
	}
	s.KeyDrawer = KeyDrawer{
		MinCountdown:   gosu.TimeToTick(30),
		Countdowns:     make([]int, keyCount),
		KeyUpSprites:   s.KeyUpSprites,
		KeyDownSprites: s.KeyDownSprites,
	}
	s.JudgmentDrawer = NewJudgmentDrawer()
	s.ScoreDrawer = gosu.NewScoreDrawer()
	s.ComboDrawer = gosu.NumberDrawer{
		BaseDrawer: draws.BaseDrawer{
			MaxCountdown: gosu.TimeToTick(2000),
		},
		DigitWidth: s.ComboSprites[0].W(),
		DigitGap:   ComboDigitGap,
		Bounce:     0.85,
		Sprites:    s.ComboSprites,
	}
	s.MeterDrawer = gosu.NewMeterDrawer(Judgments, JudgmentColors)
	return s, nil
}

// Farther note has larger position. Tail's Position is always larger than Head's.
// Need to re-calculate positions when Speed has changed.
func (s *ScenePlay) SetSpeed() {
	c := s.Chart
	old := s.SpeedScale
	new := SpeedScale
	s.Cursor *= new / old
	for _, tp := range c.TransPoints {
		tp.Position *= new / old
	}
	for _, n := range c.Notes {
		n.Position *= new / old
	}
	for _, b := range c.Bars {
		b.Position *= new / old
	}
	s.SpeedScale = new
}

// Todo: apply other values of TransPoint (Volume has finished so far)
// Todo: keep playing music when making SceneResult
func (s *ScenePlay) Update() any {
	defer s.Ticker()
	if s.IsDone() {
		s.MusicPlayer.Close()
		return gosu.PlayToResultArgs{Result: s.NewResult(s.Chart.MD5)}
	}
	// if s.Now == 0 {
	// 	s.MusicPlayer.Play()
	// }
	// if s.Now == 150 {
	// 	s.MusicPlayer.Player.Seek(time.Duration(s.Now) * time.Millisecond)
	// }
	s.MusicPlayer.Update()
	// fmt.Printf("game: %dms music: %s\n", s.Now, s.MusicPlayer.Player.Current())

	s.LastPressed = s.Pressed
	s.Pressed = s.FetchPressed()
	var worst gosu.Judgment
	for _, n := range s.Staged {
		if n == nil {
			continue
		}
		// if n.Type != Tail && s.KeyAction(k) == input.Hit {
		// 	if name := n.Sample.Name; name != "" {
		// 		vol := n.Sample.Volume
		// 		if vol == 0 {
		// 			vol = s.TransPoint.Volume
		// 		}
		// 		// Todo: apply effect volume change
		// 		s.Effects.PlayWithVolume(name, vol)
		// 	}
		// }
		td := n.Time - s.Now // Time difference. A negative value infers late hit
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
			var kind int = 0
			if n.Type == Tail {
				kind = 1
			}
			s.MeterDrawer.AddMark(int(td), kind)
		}
	}

	s.BarDrawer.Update(s.Cursor)
	for i := range s.NoteDrawers {
		s.NoteDrawers[i].Update(s.Cursor)
	}
	s.KeyDrawer.Update(s.LastPressed, s.Pressed)
	s.JudgmentDrawer.Update(worst)
	s.ScoreDrawer.Update(s.Scores[3])
	s.ComboDrawer.Update(s.Combo)
	s.MeterDrawer.Update()

	// Changed speed should be applied after positions are calculated.
	s.UpdateTransPoint()
	s.UpdateCursor()
	if SpeedScale != s.SpeedScale {
		s.SetSpeed()
	}
	return nil
}
func (s ScenePlay) Draw(screen *ebiten.Image) {
	s.BackgroundDrawer.Draw(screen)
	s.StageDrawer.Draw(screen)
	s.BarDrawer.Draw(screen)
	for _, d := range s.NoteDrawers {
		d.Draw(screen)
	}
	s.KeyDrawer.Draw(screen)
	s.JudgmentDrawer.Draw(screen)
	s.ScoreDrawer.Draw(screen)
	s.ComboDrawer.Draw(screen)
	s.MeterDrawer.Draw(screen)
	s.DebugPrint(screen)
}

func (s ScenePlay) DebugPrint(screen *ebiten.Image) {
	ebitenutil.DebugPrint(screen, fmt.Sprintf(
		"FPS: %.2f\nTPS: %.2f\nTime: %.3fs/%.0fs\n\n"+
			"Score: %.0f | %.0f \nFlow: %.0f/100\nCombo: %d\n\n"+
			"Flow rate: %.2f%%\nAccuracy: %.2f%%\nExtra: %.2f%%\nJudgment counts: %v\n\n"+
			"Speed scale (Z/X): %.0f (x%.2f)\n(Exposure time: %.fms)\n\n"+
			"Music volume (Alt+ Left/Right): %.0f%%\nEffect volume (Ctrl+ Left/Right): %.0f%%\n\n"+
			"Press ESC to select a song.\nPress TAB to pause.\n\n"+
			"Offset (Shift+ Left/Right): %dms\n",
		ebiten.ActualFPS(), ebiten.ActualTPS(), float64(s.Now)/1000, float64(s.Chart.Duration())/1000,
		s.Scores[gosu.Total], s.ScoreBounds[gosu.Total], s.Flow*100, s.Combo,
		s.Ratios[0]*100, s.Ratios[1]*100, s.Ratios[2]*100, s.JudgmentCounts,
		s.SpeedScale*100, s.TransPoint.Speed, ExposureTime(s.CurrentSpeed()),
		gosu.MusicVolume*100, gosu.EffectVolume*100,
		gosu.Offset))
}

// 1 pixel is 1 millisecond.
func ExposureTime(speed float64) float64 { return HitPosition / speed }

// func (s ScenePlay) Time() int64           { return s.Timer.Time() }
func (s ScenePlay) Speed()                { s.CurrentSpeed() }
func (s ScenePlay) CurrentSpeed() float64 { return s.TransPoint.Speed * s.SpeedScale }

// Supposes one current TransPoint can increment cursor precisely.
func (s *ScenePlay) UpdateCursor() {
	duration := float64(s.Now - s.TransPoint.Time)
	s.Cursor = s.TransPoint.Position + duration*s.CurrentSpeed()
}
func (s *ScenePlay) UpdateTransPoint() {
	s.TransPoint = s.TransPoint.FetchByTime(s.Now)
}
