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
	"github.com/hndada/gosu/input"
)

type ScenePlay struct {
	gosu.Timer
	time  int64 // Just a cache.
	Chart *Chart

	gosu.MusicPlayer
	CustomEffectPlayer gosu.EffectPlayer
	gosu.EffectPlayer
	gosu.KeyLogger

	*gosu.TransPoint
	SpeedHandler ctrl.F64Handler
	Speed        float64

	gosu.Scorer
	StagedNote        *Note
	StagedDot         *Dot
	StagedShake       *Note
	LastHitTimes      [4]int64      // For judging big note.
	StagedJudgement   gosu.Judgment // For judging big note.
	ShakeWaitingColor int
	// WaitingKeys [2]int // For judging big note.
	// StagedJudge  gosu.Judgment
	// IsBigStaged bool // For judging big note.
	// Flow        float64
	// Combo       int
	// DotCount    int
	// ShakeCount  int
	// NoteWeights float64
	// DotWeights   float64
	// ShakeWeights float64

	Skin             // The skin may be applied some custom settings: on/off some sprites
	BackgroundDrawer gosu.BackgroundDrawer
	StageDrawer      StageDrawer

	BarDrawer   BarDrawer
	ShakeDrawer ShakeDrawer
	// BodyDrawer  BodyDrawer
	// DotDrawer   DotDrawer
	RollDrawer RollDrawer
	NoteDrawer NoteDarwer
	KeyDrawer  KeyDrawer

	ScoreDrawer      gosu.ScoreDrawer
	ComboDrawer      gosu.NumberDrawer
	DotCountDrawer   gosu.NumberDrawer
	ShakeCountDrawer gosu.NumberDrawer
	JudgmentDrawer   JudgmentDrawer
	MeterDrawer      gosu.MeterDrawer
}

// Todo: actual auto replay generator for gimmick charts
// Todo: add Mods
// Todo: support mods: show Piano's ScenePlay during Drum's ScenePlay
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
	for _, colorName := range []string{"red", "blue"} {
		for _, sizeName := range []string{"regular", "big"} {
			path := fmt.Sprintf("skin/drum/sound/%s-%s.wav", colorName, sizeName)
			_ = s.Effects.Register(path)
		}
	}
	s.CustomEffectPlayer = gosu.NewEffectPlayer(gosu.EffectVolumeHandler)
	for _, n := range c.Notes {
		if path, ok := n.Sample.Path(cpath); ok {
			_ = s.Effects.Register(path)
		}
	}

	s.KeyLogger = gosu.NewKeyLogger(KeySettings[:])
	if rf != nil {
		s.KeyLogger.FetchPressed = NewReplayListener(rf, s.time)
	}

	s.TransPoint = c.TransPoints[0]
	s.SpeedHandler = sh
	s.Speed = 1
	s.SetSpeed()

	s.Scorer = gosu.NewScorer()
	s.JudgmentCounts = make([]int, len(JudgmentCountKinds))
	// s.FlowMarks = make([]float64, 0, c.Duration()/1000)
	for _, n := range c.Notes {
		s.MaxWeights[gosu.Flow] += n.Weight()
	}
	s.MaxWeights[gosu.Acc] = s.MaxWeights[gosu.Flow]
	for _, dot := range c.Dots {
		s.MaxWeights[gosu.Extra] += dot.Weight()
	}
	for _, shake := range c.Shakes {
		s.MaxWeights[gosu.Extra] += shake.Weight()
	}

	if len(s.Chart.Notes) > 0 {
		s.StagedNote = s.Chart.Notes[0]
	}
	if len(s.Chart.Dots) > 0 {
		s.StagedDot = s.Chart.Dots[0]
	}

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
	s.ShakeDrawer = ShakeDrawer{
		BorderSprite: s.ShakeBorderSprite,
		Sprite:       s.ShakeSprite,
		Time:         s.time,
		Staged:       s.StagedShake,
	}
	s.RollDrawer = RollDrawer{
		BodySprites: s.BodySprites,
		TailSprites: s.TailSprites,
		DotSprite:   s.DotSprite,
		Time:        s.time,
		Rolls:       s.Chart.Notes,
		Dots:        s.Chart.Dots,
		StagedDot:   s.StagedDot,
	}
	// s.DotDrawer = DotDrawer{
	// 	Sprite: s.DotSprite,
	// 	Time:   s.time,
	// 	Dots:   s.Chart.Dots,
	// 	Staged: s.StagedDot,
	// }
	s.NoteDrawer = NoteDarwer{
		NoteSprites:    s.NoteSprites,
		OverlaySprites: s.OverlaySprites,
		// ShakeNoteSprite: s.ShakeSprites[ShakeNote],
		Time:   s.time,
		Notes:  s.Chart.Notes,
		Rolls:  s.Chart.Rolls,
		Shakes: s.Chart.Shakes,
	}
	s.KeyDrawer = KeyDrawer{
		MaxCountdown: gosu.TimeToTick(50),
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
	// s.DotCountDrawer = gosu.NumberDrawer{
	// 	Sprites:    s.DotCountSprites,
	// 	DigitWidth: s.DotCountSprites[0].W(),
	// 	DigitGap:   DotCountDigitGap,
	// 	Bounce:     false,
	// }
	// s.ShakeCountDrawer = gosu.NumberDrawer{
	// 	Sprites:    s.ShakeCountSprites,
	// 	DigitWidth: s.ShakeCountSprites[0].W(),
	// 	DigitGap:   ShakeCountDigitGap,
	// 	Bounce:     false,
	// }
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
		return gosu.PlayToResultArgs{Result: s.NewResult(s.Chart.MD5)}

	}
	if s.Tick == 0 {
		s.MusicPlayer.Play()
	}
	s.MusicPlayer.Update()

	s.LastPressed = s.Pressed
	s.Pressed = s.FetchPressed()
	keyActions := s.KeyActions()
	// var playSample [4]bool // red-regular, red-big, blue-regular, blue-big in order.
	// // Determine whether plays big or regular sample, or none.
	// for color, keys := range [][]int{{1, 2}, {0, 3}} {
	// 	if !hits[keys[0]] && !hits[keys[1]] {
	// 		continue
	// 	}
	// 	if hits[keys[0]] && s.time-s.LastHitTimes[keys[1]] < Good.Window {
	// 		playSample[2*color+1] = true
	// 		continue
	// 	}
	// 	if hits[keys[1]] && s.time-s.LastHitTimes[keys[0]] < Good.Window {
	// 		playSample[2*color+1] = true
	// 		continue
	// 	}
	// 	playSample[2*color] = true
	// }
	if s.StagedJudgement.Window != None {
		color := s.StagedNote.Color
		if keyActions[color-1] == Regular || IsOtherColorHit(keyActions, color) {
			s.MarkNote(s.StagedNote, s.StagedJudgement, false)
			s.StagedJudgement = gosu.Judgment{}
		}
	}
	var (
		samples  [2]gosu.Sample
		judgment gosu.Judgment
		big      bool
	)
	if n := s.StagedNote; n != nil {
		td := n.Time - s.time // Time difference. A negative value means late hit.
		if j, b := VerdictNote(n, keyActions, td); j.Window != 0 {
			s.MarkNote(n, j, b)
			s.MeterDrawer.AddMark(int(td), 0)
			judgment = j
			big = b
			samples[n.Color-1] = n.Sample
		}
	}
	if dot := s.StagedDot; dot != nil {
		td := dot.Time - s.time
		if marked, hit := VerdictDot(dot, keyActions, td); marked {
			s.MarkDot(dot, hit)
			s.MeterDrawer.AddMark(int(td), 1)
		}
	}
	func() {
		shake := s.StagedShake
		if shake == nil {
			return
		}
		if t := shake.Time - s.time; t > 0 {
			return
		}
		if t := shake.Time + shake.Duration - s.time; t < 0 {
			s.MarkShake(shake, true)
			return
		}
		waiting := s.ShakeWaitingColor
		if next := VerdictShake(shake, keyActions, waiting); next != waiting {
			s.MarkShake(shake, false)
			s.ShakeWaitingColor = next
		}
	}()

	// Todo: apply effect volume change from changer
	for i, size := range keyActions {
		if size == None {
			continue
		}
		sample := samples[i]
		vol := sample.Volume
		if vol == 0 {
			vol = s.TransPoint.Volume
		}
		if sample.Name != "" {
			s.CustomEffectPlayer.Effects.PlayWithVolume(sample.Name, vol)
		} else {
			s.Effects.PlayWithVolume(DefaultSampleNames[i][size-1], vol)
		}
	}

	s.BarDrawer.Update(s.time)
	// s.ShakeDrawer.Update()
	// s.BodyDrawer.Update(s.time)
	// s.DotDrawer.Update(s.time, s.StagedDot)
	s.RollDrawer.Update(s.time, s.StagedDot)
	s.NoteDrawer.Update(s.time, s.BPM)
	s.KeyDrawer.Update(s.LastPressed, s.Pressed)

	s.ScoreDrawer.Update(s.Scores[0])
	s.ComboDrawer.Update(s.Combo)
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
	s.ShakeDrawer.Draw(screen)
	s.RollDrawer.Draw(screen)
	s.NoteDrawer.Draw(screen)
	s.KeyDrawer.Draw(screen)

	s.ScoreDrawer.Draw(screen)
	s.ComboDrawer.Draw(screen)
	// s.DotCountDrawer.Draw(screen)
	// s.ShakeCountDrawer.Draw(screen)
	s.JudgmentDrawer.Draw(screen)
	s.MeterDrawer.Draw(screen)

	s.DebugPrint(screen)
}

func (s ScenePlay) DebugPrint(screen *ebiten.Image) {
	ebitenutil.DebugPrint(screen, fmt.Sprintf(
		"FPS: %.2f\nTPS: %.2f\nTime: %.3fs/%.0fs\n\n"+
			"Score: %.0f | %.0f \nFlow: %.0f/100\nCombo: %d\n\n"+
			"Flow rate: %.2f%%\nAccuracy: %.2f%%\n(Kool: %.2f%%)\nJudgment counts: %v\n\n"+
			"Speed (Press 8/9): %.0f | %.0f\n(Exposure time: %.fms)\n\n"+
			// "Music volume (Press 1/2): %.0f%%\nEffect volume (Press 3/4): %.0f%%\n\n"+
			"Vsync: %v\n",
		ebiten.ActualFPS(), ebiten.ActualTPS(), float64(s.Time())/1000, float64(s.Chart.Duration())/1000,
		s.Scores[0], s.ScoreBounds[0], s.Flow*100, s.Combo,
		s.Ratios[0]*100, s.Ratios[1]*100, s.Ratios[2]*100, s.JudgmentCounts,
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

//	func (s *ScenePlay) KeyActions0() (as [2][2]bool) {
//		const (
//			regular = 0
//			big     = 1
//		)
//		var hits [4]bool
//		for k := range hits {
//			hits[k] = s.KeyLogger.KeyAction(k) == input.Hit
//		}
//		for color, keys := range [][]int{{1, 2}, {0, 3}} {
//			switch {
//			case !hits[keys[0]] && !hits[keys[1]]:
//				// Does nothing.
//			case hits[keys[0]] && s.time-s.LastHitTimes[keys[1]] < Good.Window,
//				hits[keys[1]] && s.time-s.LastHitTimes[keys[0]] < Good.Window:
//				as[color][big] = true
//			default:
//				as[color][regular] = true
//			}
//		}
//		for k, hit := range hits {
//			if hit {
//				s.LastHitTimes[k] = s.time
//			}
//		}
//		return
//	}
func (s *ScenePlay) KeyActions() (as [2]int) {
	var hits [4]bool
	for k := range hits {
		hits[k] = s.KeyLogger.KeyAction(k) == input.Hit
	}
	for color, keys := range [][]int{{1, 2}, {0, 3}} {
		switch {
		case !hits[keys[0]] && !hits[keys[1]]:
			as[color] = None
		case hits[keys[0]] && s.time-s.LastHitTimes[keys[1]] < Good.Window,
			hits[keys[1]] && s.time-s.LastHitTimes[keys[0]] < Good.Window:
			as[color] = Big
		default:
			as[color] = Regular
		}
	}
	for k, hit := range hits {
		if hit {
			s.LastHitTimes[k] = s.time
		}
	}
	return
}

var DefaultSampleNames = [2][2]string{
	{"red-regular", "red-big"},
	{"blue-regular", "blue-big"},
}
