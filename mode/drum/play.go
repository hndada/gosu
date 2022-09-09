package drum

import (
	"fmt"
	"runtime/debug"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hndada/gosu"
	"github.com/hndada/gosu/ctrl"
	"github.com/hndada/gosu/draws"
	"github.com/hndada/gosu/format/osr"
)

type ScenePlay struct {
	gosu.Timer
	time  int64 // Just a cache.
	Chart *Chart

	gosu.MusicPlayer
	gosu.EffectPlayer
	gosu.KeyLogger

	*gosu.TransPoint
	SpeedHandler ctrl.F64Handler
	Speed        float64

	gosu.Result
	StagedNote  *Note
	StagedDot   *Dot
	StagedShake *Note
	Flow        float64
	Combo       int
	DotCount    int
	ShakeCount  int
	// NoteWeights is a sum of weight of marked notes.
	// This is also max value of each score sum can get at the time.
	NoteWeights float64

	Skin             // The skin may be applied some custom settings: on/off some sprites
	BackgroundDrawer gosu.BackgroundDrawer
	StageDrawer      StageDrawer

	BarDrawer   BarDrawer
	ShakeDrawer ShakeDrawer
	BodyDrawer  BodyDrawer
	DotDrawer   DotDrawer
	NoteDrawer  NoteDarwer
	KeyDrawer   KeyDrawer

	ScoreDrawer      gosu.ScoreDrawer
	ComboDrawer      gosu.NumberDrawer
	DotCountDrawer   gosu.NumberDrawer
	ShakeCountDrawer gosu.NumberDrawer
	JudgmentDrawer   JudgmentDrawer
	MeterDrawer      gosu.MeterDrawer
}

// Todo: reverse notes and bars.
// Todo: make 3 kinds of note can be stages at once
// Todo: actual auto replay generator for gimmick charts
// Todo: add Mods
func NewScenePlay(cpath string, rf *osr.Format, sh ctrl.F64Handler) (scene gosu.Scene, err error) {
	s := new(ScenePlay)
	s.Chart, err = NewChart(cpath)
	if err != nil {
		return
	}
	c := s.Chart

	if rf != nil {
		s.SetTicks(rf.BufferTime(), c.Duration())
	} else {
		s.SetTicks(-1800, c.Duration())
	}
	s.time = s.Time()
	if path, ok := c.MusicPath(cpath); ok {
		s.MusicPlayer, err = gosu.NewMusicPlayer(gosu.MusicVolumeHandler, path)
		if err != nil {
			return
		}
	}
	s.EffectPlayer = gosu.NewEffectPlayer(gosu.EffectVolumeHandler)
	for _, n := range c.Notes {
		if path, ok := n.Sample.Path(cpath); ok {
			_ = s.Effects.Register(path)
		}
	}
	s.KeyLogger = gosu.NewKeyLogger(KeySettings[:])
	if rf != nil {
		s.KeyLogger.FetchPressed = gosu.NewReplayListener(rf, 4, s.time)
	}

	s.TransPoint = c.TransPoints[0]
	s.SpeedHandler = sh
	s.Speed = 1
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
	if len(s.Chart.Notes) > 0 {
		s.StagedNote = s.Chart.Notes[0]
	}
	if len(s.Chart.Dots) > 0 {
		s.StagedDot = s.Chart.Dots[0]
	}
	s.Flow = 1

	s.Skin = DefaultSkin
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
	s.BarDrawer = BarDrawer{
		Sprite: s.BarSprite,
		Time:   s.time,
		Bars:   s.Chart.Bars,
	}
	// ShakeDrawer
	s.BodyDrawer = BodyDrawer{
		BodySprites: s.BodySprites,
		TailSprite:  s.TailSprites,
		Time:        s.time,
		Notes:       s.Chart.Notes,
	}
	s.DotDrawer = DotDrawer{
		Sprite: s.DotSprite,
		Time:   s.time,
		Dots:   s.Chart.Dots,
		Staged: s.StagedDot,
	}
	s.NoteDrawer = NoteDarwer{
		NoteSprites:     s.NoteSprites,
		OverlaySprites:  s.OverlaySprites,
		ShakeNoteSprite: s.ShakeSprites[ShakeNote],
		Time:            s.time,
		Notes:           s.Chart.Notes,
		Shakes:          s.Chart.Shakes,
	}
	s.KeyDrawer = KeyDrawer{
		MaxCountdown: gosu.TimeToTick(30),
		Field:        s.KeyFieldSprite,
		Keys:         s.KeySprites,
	}

	s.ScoreDrawer = gosu.NewScoreDrawer()
	s.ComboDrawer = gosu.NumberDrawer{
		BaseDrawer: draws.BaseDrawer{
			MaxCountdown: gosu.TimeToTick(2000),
		},
		Sprites:    s.ComboSprites,
		DigitWidth: s.ComboSprites[0].W(),
		DigitGap:   ComboDigitGap,
		Bounce:     true,
	}
	s.DotCountDrawer = gosu.NumberDrawer{
		Sprites:    s.DotCountSprites,
		DigitWidth: s.DotCountSprites[0].W(),
		DigitGap:   DotCountDigitGap,
		Bounce:     false,
	}
	s.ShakeCountDrawer = gosu.NumberDrawer{
		Sprites:    s.ShakeCountSprites,
		DigitWidth: s.ShakeCountSprites[0].W(),
		DigitGap:   ShakeCountDigitGap,
		Bounce:     false,
	}
	s.JudgmentDrawer = JudgmentDrawer{
		BaseDrawer: draws.BaseDrawer{
			MaxCountdown: gosu.TimeToTick(600),
		},
		Sprites: s.JudgmentSprites,
	}
	s.MeterDrawer = gosu.NewMeterDrawer(Judgments, JudgmentColors)

	title := fmt.Sprintf("gosu - %s - [%s]", c.MusicName, c.ChartName)
	ebiten.SetWindowTitle(title)
	debug.SetGCPercent(0)
	return s, nil
}

// Todo: flush Roll immediately
// Todo: apply Highlight
// Todo: keep playing music when making SceneResult
func (s *ScenePlay) Update() any {
	defer s.Ticker()
	s.time = s.Time()
	if s.IsDone() {
		debug.SetGCPercent(100)
		s.MusicPlayer.Close()
		return gosu.PlayToResultArgs{
			Result: s.Result,
		}
	}
	if s.Tick == 0 {
		s.MusicPlayer.Play()
	}
	s.MusicPlayer.Update()

	s.LastPressed = s.Pressed
	s.Pressed = s.FetchPressed()

	var judgment gosu.Judgment
	var big bool
	// for k, n := range s.Staged {
	// 	if n == nil {
	// 		continue
	// 	}
	// 	if n.Type != Tail && s.KeyAction(k) == input.Hit {
	// 		if name := n.Sample.Name; name != "" {
	// 			vol := n.Sample.Volume
	// 			if vol == 0 {
	// 				vol = s.TransPoint.Volume
	// 			}
	// 			// Todo: apply effect volume change
	// 			s.Effects.PlayWithVolume(name, vol)
	// 		}
	// 	}
	// 	td := n.Time - s.time // Time difference. A negative value infers late hit
	// 	if n.Marked {
	// 		if n.Type != Tail {
	// 			return fmt.Errorf("non-Tail note has not flushed")
	// 		}
	// 		if td < Miss.Window { // Keep Tail staged until near ends.
	// 			s.Staged[n.Key] = n.Next
	// 		}
	// 		continue
	// 	}
	// 	if j := Verdict(n.Type, s.KeyAction(n.Key), td); j.Window != 0 {
	// 		s.MarkNote(n, j)
	// 		judgment = j
	// 		var colorType int = 0
	// 		if n.Type == Tail {
	// 			colorType = 1
	// 		}
	// 		s.MeterDrawer.AddMark(int(td), colorType)
	// 	}
	// }

	s.BarDrawer.Update(s.time)
	// s.ShakeDrawer.Update()
	s.BodyDrawer.Update(s.time)
	s.DotDrawer.Update(s.time, s.StagedDot)
	s.NoteDrawer.Update(s.time, 0) // Temporary value at overlay
	s.KeyDrawer.Update(s.LastPressed, s.Pressed)

	// s.ScoreDrawer.Update(s.Score())
	s.ComboDrawer.Update(s.Combo)
	s.DotCountDrawer.Update(s.DotCount)
	s.ShakeCountDrawer.Update(s.ShakeCount)
	s.JudgmentDrawer.Update(judgment, big)
	s.MeterDrawer.Update()

	// Changed speed should be applied after positions are calculated.
	s.UpdateTransPoint()
	if fired := s.SpeedHandler.Update(); fired {
		s.SetSpeed()
	}
	return nil
}
func (s ScenePlay) Draw(screen *ebiten.Image) {
	s.BackgroundDrawer.Draw(screen)
	s.StageDrawer.Draw(screen)

	s.BarDrawer.Draw(screen)
	// s.ShakeDrawer.Draw(screen)
	s.BodyDrawer.Draw(screen)
	s.DotDrawer.Draw(screen)
	s.NoteDrawer.Draw(screen)
	s.KeyDrawer.Draw(screen)

	s.ScoreDrawer.Draw(screen)
	s.ComboDrawer.Draw(screen)
	s.DotCountDrawer.Draw(screen)
	s.ShakeCountDrawer.Draw(screen)
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
		"FPS: %.2f\nTPS: %.2f\nTime: %.3fs/%.0fs\n\n"+
			"Score: %.0f | %.0f \nFlow: %.0f/100\nCombo: %d\n\n"+
			// "Flow rate: %.2f%%\nAccuracy: %.2f%%\n(Kool: %.2f%%)\nJudgment counts: %v\n\n"+
			"Speed (Press 8/9): %.0f | %.0f\n(Exposure time: %.fms)\n\n"+
			// "Music volume (Press 1/2): %.0f%%\nEffect volume (Press 3/4): %.0f%%\n\n"+
			"Vsync: %v\n",
		ebiten.ActualFPS(), ebiten.ActualTPS(), float64(s.Time())/1000, float64(s.Chart.Duration())/1000,
		// s.Score(), s.ScoreBound(), s.Flow*100, s.Combo,
		fr*100, ar*100, rr*100, s.JudgmentCounts,
		s.Speed*100, *s.SpeedHandler.Target*100, ExposureTime(s.CurrentSpeed()),
		// gosu.MusicVolume*100, gosu.EffectVolume*100,
		gosu.VsyncSwitch))
}

// Farther note has larger position. Tail's Position is always larger than Head's.
// Need to re-calculate positions when Speed has changed.
func (s *ScenePlay) SetSpeed() {
	old := s.Speed
	new := *s.SpeedHandler.Target
	for _, tp := range s.Chart.TransPoints {
		tp.Speed *= new / old
	}
	for _, n := range s.Chart.Notes {
		n.Speed *= new / old
	}
	for _, b := range s.Chart.Bars {
		b.Speed *= new / old
	}
	for _, d := range s.Chart.Dots {
		d.Speed *= new / old
	}
	s.Speed = new
}

// 1 pixel is 1 millisecond.
// Todo: Separate NoteHeight / 2 at piano mode
func ExposureTime(speedScale float64) float64 {
	return (screenSizeX - HitPosition) / speedScale
}
func (s *ScenePlay) UpdateTransPoint() {
	s.TransPoint = s.TransPoint.FetchByTime(s.Time())
}

func (s ScenePlay) Time() int64           { return s.Timer.Time() }
func (s ScenePlay) CurrentSpeed() float64 { return s.TransPoint.Speed * s.Speed }
