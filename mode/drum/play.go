package drum

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hndada/gosu"
	"github.com/hndada/gosu/audios"
	"github.com/hndada/gosu/draws"
	"github.com/hndada/gosu/format/osr"
	"github.com/hndada/gosu/input"
)

// No custom hitsound at Drum mode.
type ScenePlay struct {
	Chart *Chart
	gosu.Timer
	time int64 // Just a cache.
	gosu.MusicPlayer
	SoundEffectBytes [2][2][]byte
	gosu.KeyLogger
	KeyActions [2]int

	*gosu.TransPoint
	SpeedScale         float64
	StagedNote         *Note
	StagedDot          *Dot
	StagedShake        *Note
	LastHitTimes       [4]int64      // For judging big note.
	StagedJudgment     gosu.Judgment // For judging big note.
	StagedJudgmentTime int64
	ShakeWaitingColor  int
	gosu.Scorer
	// WaitingKeys [2]int // For judging big note.
	// StagedJudge  gosu.Judgment
	// IsBigStaged bool // For judging big note.

	// Skin may be applied some custom settings: on/off some sprites
	Skin
	BackgroundDrawer gosu.BackgroundDrawer
	StageDrawer      StageDrawer
	BarDrawer        BarDrawer
	JudgmentDrawer   JudgmentDrawer

	ShakeDrawer ShakeDrawer
	RollDrawer  RollDrawer
	NoteDrawer  NoteDarwer

	KeyDrawer    KeyDrawer
	DancerDrawer DancerDrawer
	ScoreDrawer  gosu.ScoreDrawer
	ComboDrawer  gosu.NumberDrawer
	MeterDrawer  gosu.MeterDrawer
}

func (s ScenePlay) Time() int64 { return s.Timer.Time }

// Todo: actual auto replay generator for gimmick charts
// Todo: support mods: show Piano's ScenePlay during Drum's ScenePlay
func NewScenePlay(cpath string, rf *osr.Format) (scene gosu.Scene, err error) {
	s := new(ScenePlay)
	s.Chart, err = NewChart(cpath)
	if err != nil {
		return
	}
	c := s.Chart
	gosu.SetTitle(c.ChartHeader)
	s.SetTicks(c.Duration())
	s.time = s.Time()
	if path, ok := c.MusicPath(cpath); ok {
		s.MusicPlayer, err = gosu.NewMusicPlayer(path) //(gosu.MusicVolumeHandler, path)
		if err != nil {
			return
		}
	}
	for i, colorName := range []string{"regular", "big"} {
		for j, sizeName := range []string{"red", "blue"} {
			path := fmt.Sprintf("skin/drum/sound/%s/%s.wav", colorName, sizeName)
			b, err := audios.NewBytes(path)
			if err != nil {
				panic(err)
			}
			s.SoundEffectBytes[i][j] = b
		}
	}
	s.KeyLogger = gosu.NewKeyLogger(KeySettings[:])
	if rf != nil {
		s.KeyLogger.FetchPressed = NewReplayListener(rf, s.time)
	}

	s.TransPoint = c.TransPoints[0]
	s.SpeedScale = 1
	s.SetSpeed()
	s.Scorer = gosu.NewScorer(c.ScoreFactors)
	s.JudgmentCounts = make([]int, len(JudgmentCountKinds))
	// s.FlowMarks = make([]float64, 0, c.Duration()/1000)
	for _, n := range c.Notes {
		s.MaxWeights[gosu.Flow] += n.Weight()
	}
	s.MaxWeights[gosu.Acc] = s.MaxWeights[gosu.Flow]
	for _, n := range c.Dots {
		s.MaxWeights[gosu.Extra] += n.Weight()
	}
	for _, n := range c.Shakes {
		s.MaxWeights[gosu.Extra] += n.Weight()
	}
	s.SetMaxScores()
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
		Brightness: &gosu.BackgroundBrightness,
		Sprite:     gosu.DefaultBackground,
	}
	if bg := gosu.NewBackground(c.BackgroundPath(cpath)); bg.IsValid() {
		s.BackgroundDrawer.Sprite = bg
	}
	s.StageDrawer = StageDrawer{
		Hightlight:  s.Highlight,
		FieldSprite: s.FieldSprite,
		HintSprites: s.HintSprites,
	}
	s.BarDrawer = BarDrawer{
		Time:   s.time,
		Bars:   c.Bars,
		Sprite: s.BarSprite,
	}
	s.JudgmentDrawer = JudgmentDrawer{
		BaseDrawer: draws.BaseDrawer{
			MaxCountdown: gosu.TimeToTick(600),
		},
		Sprites: s.JudgmentSprites,
	}
	s.ShakeDrawer = ShakeDrawer{
		Time:         s.time,
		Staged:       s.StagedShake,
		BorderSprite: s.ShakeBorderSprite,
		ShakeSprite:  s.ShakeSprite,
	}
	s.RollDrawer = RollDrawer{
		Time:        s.time,
		Rolls:       c.Rolls,
		Dots:        c.Dots,
		HeadSprites: s.HeadSprites,
		BodySprites: s.BodySprites,
		TailSprites: s.TailSprites,
		DotSprite:   s.DotSprite,
	}
	s.NoteDrawer = NoteDarwer{
		Time:        s.time,
		Notes:       c.Notes,
		Rolls:       c.Rolls,
		Shakes:      c.Shakes,
		NoteSprites: s.NoteSprites,
	}
	for i, sprites := range s.OverlaySprites {
		s.NoteDrawer.OverlayDrawers[i] = draws.AnimationDrawer{
			Time:      s.time,
			Duration:  int64(60000 / ScaledBPM(s.BPM)),
			StartTime: s.TransPoint.Time,
			Sprites:   sprites,
		}
	}
	s.KeyDrawer = KeyDrawer{
		MaxCountdown: gosu.TimeToTick(75),
		Field:        s.KeyFieldSprite,
		Keys:         s.KeySprites,
	}
	s.DancerDrawer.AnimationEndTime = s.time
	if s.Highlight {
		s.DancerDrawer.Mode = DancerHigh
	}
	for i := range s.DancerDrawer.AnimationDrawers {
		s.DancerDrawer.AnimationDrawers[i].Sprites = s.DancerSprites[i]
	}
	s.ScoreDrawer = gosu.NewScoreDrawer()
	s.ComboDrawer = gosu.NumberDrawer{
		BaseDrawer: draws.BaseDrawer{
			MaxCountdown: gosu.TimeToTick(2000),
		},
		Sprites:    s.ComboSprites,
		DigitWidth: s.ComboSprites[0].W(),
		DigitGap:   ComboDigitGap,
		Bounce:     1.25,
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
	for _, tp := range c.TransPoints {
		tp.Speed *= new / old
	}
	for _, b := range c.Bars {
		b.Speed *= new / old
	}
	for _, ns := range [][]*Note{c.Notes, c.Rolls, c.Shakes} {
		for _, n := range ns {
			n.Speed *= new / old
		}
	}
	for _, n := range c.Dots { // Not a Note type.
		n.Speed *= new / old
	}
	s.SpeedScale = new
}

func (s *ScenePlay) Update() any {
	defer func() { s.Ticker(); s.time = s.Time() }()
	if s.IsDone() {
		s.MusicPlayer.Close()
		return gosu.PlayToResultArgs{Result: s.NewResult(s.Chart.MD5)}
	}
	if s.time == 0 {
		s.MusicPlayer.Play()
	}
	s.MusicPlayer.Update()
	// fmt.Printf("game: %dms music: %s\n", s.Time(), s.MusicPlayer.Player.Current())

	s.LastPressed = s.Pressed
	s.Pressed = s.FetchPressed()
	s.UpdateKeyActions()

	var (
		judgment gosu.Judgment
		big      bool
	)
	if j := s.StagedJudgment; j.Window != 0 {
		n := s.StagedNote
		td1 := s.StagedNote.Time - s.time
		td2 := s.StagedJudgmentTime - s.time
		td3 := s.StagedNote.Time - s.StagedJudgmentTime
		if td1 < -Miss.Window || td2 < -BigHitTimeDifferenceBound ||
			s.KeyActions[n.Color] == Regular ||
			IsOtherColorHit(s.KeyActions, n.Color) {
			s.MarkNote(n, s.StagedJudgment, false)
			s.MeterDrawer.AddMark(int(td3), 0)
			judgment = j
			big = false
			s.StagedJudgment = gosu.Judgment{}
		}
	}
	if n := s.StagedNote; n != nil {
		td := n.Time - s.time // A negative value means late hit.
		if j, b := VerdictNote(n, s.KeyActions, td); j.Window != 0 {
			if n.Size == Big && !b {
				s.StagedJudgment = j
				s.StagedJudgmentTime = s.time
			} else {
				s.MarkNote(n, j, b)
				s.MeterDrawer.AddMark(int(td), 0)
				judgment = j
				big = b
			}
		}
	}
	if n := s.StagedDot; n != nil {
		td := n.Time - s.time
		if marked := VerdictDot(n, s.KeyActions, td); marked != DotReady {
			s.MarkDot(n, marked)
			s.MeterDrawer.AddMark(int(td), 1)
		}
	}
	func() {
		n := s.StagedShake
		if n == nil {
			return
		}
		if t := n.Time - s.time; t > 0 {
			return
		}
		if t := n.Time + n.Duration - s.time; t < 0 {
			s.MarkShake(n, true)
			return
		}
		waiting := s.ShakeWaitingColor
		if next := VerdictShake(n, s.KeyActions, waiting); next != waiting {
			s.MarkShake(n, false)
			s.ShakeWaitingColor = next
		}
	}()

	// Todo: apply effect volume change from changer
	for i, size := range s.KeyActions {
		if size == SizeNone {
			continue
		}
		vol := s.TransPoint.Volume
		p := audios.Context.NewPlayerFromBytes(s.SoundEffectBytes[i][size])
		p.SetVolume(vol * gosu.EffectVolume)
		p.Play()
	}
	s.StageDrawer.Update(s.Highlight)
	s.BarDrawer.Update(s.time)
	s.JudgmentDrawer.Update(judgment, big)

	s.ShakeDrawer.Update(s.time, s.StagedShake)
	s.RollDrawer.Update(s.time)
	s.NoteDrawer.Update(s.time, s.BPM)

	s.KeyDrawer.Update(s.LastPressed, s.Pressed)
	s.DancerDrawer.Update(s.time, s.BPM, s.Combo, judgment.Is(Miss),
		!judgment.Is(Miss) && judgment.Valid(), s.Highlight)
	s.ScoreDrawer.Update(s.Scores[gosu.Total])
	s.ComboDrawer.Update(s.Combo)
	s.MeterDrawer.Update()

	// Changed speed should be applied after positions are calculated.
	s.UpdateTransPoint()
	if SpeedScale != s.SpeedScale {
		s.SetSpeed()
	}
	return nil
}
func (s ScenePlay) Draw(screen *ebiten.Image) {
	// screen.Fill(color.NRGBA{0, 255, 0, 255}) // Chroma-key
	s.BackgroundDrawer.Draw(screen)
	s.StageDrawer.Draw(screen)
	s.BarDrawer.Draw(screen)
	s.JudgmentDrawer.Draw(screen)

	s.ShakeDrawer.Draw(screen)
	s.RollDrawer.Draw(screen)
	s.NoteDrawer.Draw(screen)

	s.KeyDrawer.Draw(screen)
	s.DancerDrawer.Draw(screen)
	s.ScoreDrawer.Draw(screen)
	s.ComboDrawer.Draw(screen)
	s.MeterDrawer.Draw(screen)
	s.DebugPrint(screen)
}

func (s ScenePlay) DebugPrint(screen *ebiten.Image) {
	ebitenutil.DebugPrint(screen, fmt.Sprintf(
		"\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n"+
			"FPS: %.2f\nTPS: %.2f\nTime: %.3fs/%.0fs\n\n"+
			"Score: %.0f | %.0f \nFlow: %.0f/100\nCombo: %d\n\n"+
			"Flow rate: %.2f%%\nAccuracy: %.2f%%\nExtra: %.2f%%\nJudgment counts: %v\n\n"+
			"Speed scale (Z/X): %.0f (x%.2f)\n(Exposure time: %.fms)\n\n"+
			"Music volume (Q/W): %.0f%%\nEffect volume (A/S): %.0f%%\n\n",
		ebiten.ActualFPS(), ebiten.ActualTPS(), float64(s.time)/1000, float64(s.Chart.Duration())/1000,
		s.Scores[gosu.Total], s.ScoreBounds[gosu.Total], s.Flow*100, s.Combo,
		s.Ratios[0]*100, s.Ratios[1]*100, s.Ratios[2]*100, s.JudgmentCounts,
		s.SpeedScale*100, s.SpeedScale/s.TransPoint.Speed, ExposureTime(s.CurrentSpeed()),
		gosu.MusicVolume*100, gosu.EffectVolume*100))
}

// 1 pixel is 1 millisecond.
// Todo: Separate NoteHeight / 2 at piano mode
func ExposureTime(speedScale float64) float64 {
	return (screenSizeX - HitPosition) / speedScale
}
func (s *ScenePlay) UpdateTransPoint() {
	s.TransPoint = s.TransPoint.FetchByTime(s.time)
}
func (s ScenePlay) Speed()                { s.CurrentSpeed() }
func (s ScenePlay) CurrentSpeed() float64 { return s.TransPoint.Speed * s.SpeedScale }
func (s *ScenePlay) UpdateKeyActions() {
	var hits [4]bool
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
	}
	for k, hit := range hits {
		if hit {
			s.LastHitTimes[k] = s.time
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
