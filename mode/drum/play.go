package drum

import (
	"fmt"
	"runtime/debug"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hndada/gosu"
	"github.com/hndada/gosu/audios"
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
	// CustomEffectPlayer gosu.EffectPlayer
	// gosu.EffectPlayer
	SoundEffectBytes [2][2][]byte
	gosu.KeyLogger
	KeyActions [2]int

	*gosu.TransPoint
	SpeedHandler ctrl.F64Handler
	Speed        float64

	gosu.Scorer
	StagedNote         *Note
	StagedDot          *Dot
	StagedShake        *Note
	LastHitTimes       [4]int64      // For judging big note.
	StagedJudgment     gosu.Judgment // For judging big note.
	StagedJudgmentTime int64
	ShakeWaitingColor  int
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
	RollDrawer   RollDrawer
	NoteDrawer   NoteDarwer
	KeyDrawer    KeyDrawer
	DancerDrawer DancerDrawer

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
	// s.EffectPlayer = gosu.NewEffectPlayer(gosu.EffectVolumeHandler)
	for i, colorName := range []string{"red", "blue"} {
		for j, sizeName := range []string{"regular", "big"} {
			path := fmt.Sprintf("skin/drum/sound/%s-%s.wav", colorName, sizeName)
			b, err := audios.NewBytes(path)
			if err != nil {
				panic(err)
			}
			s.SoundEffectBytes[i][j] = b
		}
	}
	// s.CustomEffectPlayer = gosu.NewEffectPlayer(gosu.EffectVolumeHandler)
	// for _, n := range c.Notes {
	// 	if path, ok := n.Sample.Path(cpath); ok {
	// 		_ = s.Effects.Register(path)
	// 	}
	// }

	s.KeyLogger = gosu.NewKeyLogger(KeySettings[:])
	if rf != nil {
		s.KeyLogger.FetchPressed = NewReplayListener(rf, s.time)
	}

	s.TransPoint = c.TransPoints[0]
	s.SpeedHandler = sh
	s.Speed = 1
	s.SetSpeed()

	s.Scorer = gosu.NewScorer(c.ScoreFactors)
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
	s.SetMaxScores()
	// fmt.Println(s.MaxWeights, 10*s.MaxWeights[gosu.Extra]/s.MaxWeights[gosu.Flow])
	if len(c.Notes) > 0 {
		s.StagedNote = c.Notes[0]
	}
	if len(c.Dots) > 0 {
		s.StagedDot = c.Dots[0]
	}
	if len(c.Shakes) > 0 {
		s.StagedShake = c.Shakes[0]
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
		Hints: s.HintSprites,
	}
	s.BarDrawer = BarDrawer{
		Sprite: s.BarSprite,
		Time:   s.time,
		Bars:   c.Bars,
	}
	s.ShakeDrawer = ShakeDrawer{
		BorderSprite: s.ShakeBorderSprite,
		ShakeSprite:  s.ShakeSprite,
		Time:         s.time,
		Staged:       s.StagedShake,
	}
	s.RollDrawer = RollDrawer{
		HeadSprites: s.HeadSprites,
		BodySprites: s.BodySprites,
		TailSprites: s.TailSprites,
		DotSprite:   s.DotSprite,
		Time:        s.time,
		Rolls:       c.Rolls,
		Dots:        c.Dots,
		// StagedDot:   s.StagedDot,
	}
	// s.DotDrawer = DotDrawer{
	// 	Sprite: s.DotSprite,
	// 	Time:   s.time,
	// 	Dots:   c.Dots,
	// 	Staged: s.StagedDot,
	// }
	s.NoteDrawer = NoteDarwer{
		NoteSprites:    s.NoteSprites,
		OverlaySprites: s.OverlaySprites,
		// ShakeNoteSprite: s.ShakeSprites[ShakeNote],
		Time:   s.time,
		Notes:  c.Notes,
		Rolls:  c.Rolls,
		Shakes: c.Shakes,
	}
	s.KeyDrawer = KeyDrawer{
		MaxCountdown: gosu.TimeToTick(50),
		Field:        s.KeyFieldSprite,
		Keys:         s.KeySprites,
	}
	s.DancerDrawer = DancerDrawer{
		Time:           s.time,
		Duration:       50, // Temporary value.
		LastFrameTimes: [4]int64{s.time, s.time, s.time, s.time},
		Sprites:        s.DancerSprites,
	}
	// s.DancerDrawer.Update(s.time, s.BPM, false, false, 0, s.Highlight)
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
	s.UpdateKeyActions()
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
	var (
		samples  [2]gosu.Sample
		judgment gosu.Judgment
		big      bool
	)
	if j := s.StagedJudgment; j.Window != 0 {
		n := s.StagedNote
		td := s.StagedNote.Time - s.time
		td2 := s.StagedJudgmentTime - s.time
		td3 := s.StagedNote.Time - s.StagedJudgmentTime
		if td < -Miss.Window || td2 < -BigHitTimeDifferenceBound ||
			s.KeyActions[n.Color] == Regular ||
			IsOtherColorHit(s.KeyActions, n.Color) {
			s.MarkNote(n, s.StagedJudgment, false)
			s.MeterDrawer.AddMark(int(td3), 0)
			samples[n.Color] = n.Sample
			judgment = j
			big = false

			s.StagedJudgment = gosu.Judgment{}
		}
	}
	if n := s.StagedNote; n != nil {
		td := n.Time - s.time // Time difference. A negative value means late hit.
		if j, b := VerdictNote(n, s.KeyActions, td); j.Window != 0 {
			if n.Size == Big && !b {
				s.StagedJudgment = j
				s.StagedJudgmentTime = s.time
			} else {
				s.MarkNote(n, j, b)
				s.MeterDrawer.AddMark(int(td), 0)
				samples[n.Color] = n.Sample
				judgment = j
				big = b
			}
		}
	}
	if dot := s.StagedDot; dot != nil {
		td := dot.Time - s.time
		// fmt.Println(s.KeyActions)
		if marked := VerdictDot(dot, s.KeyActions, td); marked != DotReady {
			// fmt.Printf("%+v %d %d\n", dot, td, marked)
			s.MarkDot(dot, marked)
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
		if next := VerdictShake(shake, s.KeyActions, waiting); next != waiting {
			s.MarkShake(shake, false)
			s.ShakeWaitingColor = next
		}
	}()

	// Todo: apply effect volume change from changer
	for i, size := range s.KeyActions {
		if size == SizeNone {
			continue
		}
		sample := samples[i]
		vol := sample.Volume
		if vol == 0 {
			vol = s.TransPoint.Volume
		}
		p := audios.Context.NewPlayerFromBytes(s.SoundEffectBytes[i][size])
		p.SetVolume(vol * 0.25)
		p.Play()
		// if sample.Name != "" {
		// 	s.CustomEffectPlayer.Effects.PlayWithVolume(sample.Name, vol)
		// } else {
		// 	s.Effects.PlayWithVolume(DefaultSampleNames[i][size-1], vol)
		// }
	}
	// This works fine.
	// for i, k := range []ebiten.Key{ebiten.KeyQ, ebiten.KeyW, ebiten.KeyZ, ebiten.KeyX} {
	// 	if ebiten.IsKeyPressed(k) {
	// 		p := audios.Context.NewPlayerFromBytes(s.SoundEffectBytes[i/2][i%2])
	// 		p.SetVolume(1)
	// 		p.Play()
	// 	}
	// }

	s.BarDrawer.Update(s.time)
	s.StageDrawer.Update(s.Highlight)
	s.ShakeDrawer.Update(s.time, s.StagedShake)
	// s.BodyDrawer.Update(s.time)
	// s.DotDrawer.Update(s.time, s.StagedDot)
	// s.RollDrawer.Update(s.time, s.StagedDot)
	s.RollDrawer.Update(s.time)
	s.NoteDrawer.Update(s.time, s.BPM)
	s.KeyDrawer.Update(s.LastPressed, s.Pressed)
	{
		miss := judgment.Window == Miss.Window
		hit := !miss && judgment.Window != 0
		s.DancerDrawer.Update(s.time, s.BPM, miss, hit, s.Combo, s.Highlight)
	}
	s.ScoreDrawer.Update(s.Scores[gosu.Total])
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
	s.JudgmentDrawer.Draw(screen)
	s.NoteDrawer.Draw(screen)
	s.KeyDrawer.Draw(screen)
	s.DancerDrawer.Draw(screen)

	s.ScoreDrawer.Draw(screen)
	s.ComboDrawer.Draw(screen)
	// s.DotCountDrawer.Draw(screen)
	// s.ShakeCountDrawer.Draw(screen)
	// s.JudgmentDrawer.Draw(screen)
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
			"Vsync enabled: %v\n",
		ebiten.ActualFPS(), ebiten.ActualTPS(), float64(s.Time())/1000, float64(s.Chart.Duration())/1000,
		s.Scores[gosu.Total], s.ScoreBounds[gosu.Total], s.Flow*100, s.Combo,
		s.Ratios[0]*100, s.Ratios[1]*100, s.Ratios[2]*100, s.JudgmentCounts,
		s.Speed*100, *s.SpeedHandler.Target*100, ExposureTime(s.CurrentSpeed()),
		// gosu.MusicVolume*100, gosu.EffectVolume*100,
		gosu.VsyncSwitch))
}

// Farther note has larger position. Tail's Position is always larger than Head's.
// Need to re-calculate positions when Speed has changed.
func (s *ScenePlay) SetSpeed() {
	c := s.Chart
	old := s.Speed
	new := *s.SpeedHandler.Target
	for _, tp := range c.TransPoints {
		tp.Speed *= new / old
	}
	for _, b := range c.Bars {
		b.Speed *= new / old
	}
	for _, d := range c.Dots {
		d.Speed *= new / old
	}
	for _, ns := range [][]*Note{c.Notes, c.Rolls, c.Shakes} {
		for _, n := range ns {
			n.Speed *= new / old
		}
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
func (s *ScenePlay) UpdateKeyActions() {
	var hits [4]bool
	// This works fine.
	for k := range hits {
		hits[k] = s.KeyLogger.KeyAction(k) == input.Hit
	}
	for color, keys := range [][]int{{1, 2}, {0, 3}} {
		if hits[keys[0]] || hits[keys[1]] {
			// if hits[keys[0]] && s.time-s.LastHitTimes[keys[1]] < BigHitTimeDifferenceBound ||
			// 	hits[keys[1]] && s.time-s.LastHitTimes[keys[0]] < BigHitTimeDifferenceBound {
			// 	s.KeyActions[color] = Big
			// } else {
			// 	s.KeyActions[color] = Regular
			// }
			s.KeyActions[color] = Regular
		} else {
			s.KeyActions[color] = SizeNone
		}

		// fmt.Println(s.time, color, hits)
		// fmt.Println(hits[keys[0]], hits[keys[1]], s.KeyActions[color])
	}
	// fmt.Println()
	// for color, keys := range [][]int{{1, 2}, {0, 3}} {
	// 	switch {
	// 	case !hits[keys[0]] && !hits[keys[1]]:
	// 		s.KeyActions[color] = None
	// 	case hits[keys[0]] && s.time-s.LastHitTimes[keys[1]] < BigHitTimeDifferenceBound,
	// 		hits[keys[1]] && s.time-s.LastHitTimes[keys[0]] < BigHitTimeDifferenceBound:
	// 		s.KeyActions[color] = Big
	// 	default:
	// 		s.KeyActions[color] = Regular
	// 	}
	// }
	for k, hit := range hits {
		if hit {
			s.LastHitTimes[k] = s.time
			// fmt.Printf("%d: hit at %d\n", s.time, k)
		}
	}
	// for color, a := range s.KeyActions {
	// 	if a != SizeNone {
	// 		fmt.Printf("%d: %s at color %s\n", s.time, []string{"regular", "hit"}[a], []string{"red", "blue"}[color])
	// 	}
	// }
}

var DefaultSampleNames = [2][2]string{
	{"red-regular", "red-big"},
	{"blue-regular", "blue-big"},
}
