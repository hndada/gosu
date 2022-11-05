package drum

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hndada/gosu"
	"github.com/hndada/gosu/audios"
	draws "github.com/hndada/gosu/draws2"
	"github.com/hndada/gosu/format/osr"
)

type ScenePlay struct {
	Chart *Chart
	gosu.Timer
	// time int64 // Just a cache.
	gosu.MusicPlayer
	SoundEffectBytes [2][2][]byte // No custom hitsound at Drum mode.
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
	s.Timer = gosu.NewTimer(c.Duration())
	if path, ok := c.MusicPath(cpath); ok {
		s.MusicPlayer, err = gosu.NewMusicPlayer(path, &s.Timer)
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
	s.KeyLogger = gosu.NewKeyLogger(KeySettings[4][:])
	if rf != nil {
		s.KeyLogger.FetchPressed = NewReplayListener(rf, &s.Timer)
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
		Time:   s.Now,
		Bars:   c.Bars,
		Sprite: s.BarSprite,
	}
	s.JudgmentDrawer = JudgmentDrawer{
		Timer:   draws.NewTimer(gosu.TimeToTick(400)),
		Sprites: s.JudgmentSprites,
	}
	s.ShakeDrawer = ShakeDrawer{
		Time:         s.Now,
		Staged:       s.StagedShake,
		BorderSprite: s.ShakeBorderSprite,
		ShakeSprite:  s.ShakeSprite,
	}
	s.RollDrawer = RollDrawer{
		Time:        s.Now,
		Rolls:       c.Rolls,
		Dots:        c.Dots,
		HeadSprites: s.HeadSprites,
		BodySprites: s.BodySprites,
		TailSprites: s.TailSprites,
		DotSprite:   s.DotSprite,
	}
	s.NoteDrawer = NoteDarwer{
		Time:        s.Now,
		Notes:       c.Notes,
		Rolls:       c.Rolls,
		Shakes:      c.Shakes,
		NoteSprites: s.NoteSprites,
	}
	for i, sprites := range s.OverlaySprites {
		s.NoteDrawer.OverlayAnimations[i] = draws.Animation{
			Time:      s.Now,
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
	s.DancerDrawer.AnimationEndTime = s.Now
	if s.Highlight {
		s.DancerDrawer.Mode = DancerHigh
	}
	for i := range s.DancerDrawer.Animations {
		s.DancerDrawer.Animations[i].Sprites = s.DancerSprites[i]
		s.DancerDrawer.Animations[i].Update(s.Now, int64(60000/ScaledBPM(s.BPM)), false)
	}
	s.ScoreDrawer = gosu.NewScoreDrawer()
	s.ComboDrawer = gosu.NumberDrawer{
		Timer:      draws.NewTimer(gosu.TimeToTick(2000)),
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
	s.UpdateKeyActions()

	var (
		judgment gosu.Judgment
		big      bool
	)
	if s.StagedJudgment.Valid() {
		n := s.StagedNote
		j := s.StagedJudgment
		jTime := s.StagedJudgmentTime

		flush := false
		if IsOtherColorHit(s.KeyActions, n.Color) { // s.KeyActions[n.Color] == Regular ||
			flush = true
		}
		if td := jTime - s.Now; td < -MaxBigHitDuration {
			flush = true
		}
		if td := n.Time - s.Now; td < -Miss.Window {
			flush = true
		}
		if flush {
			for _, key := range [][]int{{1, 2}, {0, 3}}[n.Color] {
				s.LastHitTimes[key] = -gosu.Wait
			}
			td := n.Time - jTime
			s.MarkNote(n, j, false)
			s.MeterDrawer.AddMark(int(td), 0)
			judgment = j
			big = false
			s.StagedJudgment = gosu.Judgment{}
		}
	}
	if n := s.StagedNote; n != nil {
		td := n.Time - s.Now // A negative value means late hit.
		if j, b := VerdictNote(n, s.KeyActions, td); j.Window != 0 {
			if n.Size == Big && !b {
				s.StagedJudgment = j
				s.StagedJudgmentTime = s.Now
			} else {
				s.MarkNote(n, j, b)
				s.MeterDrawer.AddMark(int(td), 0)
				judgment = j
				big = b
				if s.StagedJudgment.Valid() {
					s.StagedJudgment = gosu.Judgment{}
				}
			}
		}
	}
	if n := s.StagedDot; n != nil {
		td := n.Time - s.Now
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
		if t := n.Time - s.Now; t > 0 {
			return
		}
		if t := n.Time + n.Duration - s.Now; t < 0 {
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
	s.BarDrawer.Update(s.Now)
	s.JudgmentDrawer.Update(judgment, big)

	s.ShakeDrawer.Update(s.Now, s.StagedShake)
	s.RollDrawer.Update(s.Now)
	s.NoteDrawer.Update(s.Now, s.BPM)

	s.KeyDrawer.Update(s.LastPressed, s.Pressed)
	s.DancerDrawer.Update(s.Now, s.BPM, s.Combo, judgment.Is(Miss),
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
			"Press ESC to select a song.\nPress TAB to pause.\n\n"+
			"FPS: %.2f\nTPS: %.2f\nTime: %.3fs/%.0fs\n\n"+
			"Score: %.0f | %.0f \nFlow: %.0f/100\nCombo: %d\n\n"+
			"Flow rate: %.2f%%\nAccuracy: %.2f%%\nExtra: %.2f%%\n"+
			"Judgment counts: %v\nPartial counts: %v\nTick counts: %v\n\n"+
			"Speed scale (PageUp/Down): %.0f (x%.2f)\n(Exposure time: %.fms)\n\n"+
			"Music volume (Alt+ Left/Right): %.0f%%\nEffect volume (Ctrl+ Left/Right): %.0f%%\n\n"+
			"Offset (Shift+ Left/Right): %dms\n",
		ebiten.ActualFPS(), ebiten.ActualTPS(), float64(s.Now)/1000, float64(s.Chart.Duration())/1000,
		s.Scores[gosu.Total], s.ScoreBounds[gosu.Total], s.Flow*100, s.Combo,
		s.Ratios[0]*100, s.Ratios[1]*100, s.Ratios[2]*100,
		s.JudgmentCounts[:3], s.JudgmentCounts[3:5], s.JudgmentCounts[5:],
		s.SpeedScale*100, s.SpeedScale/s.TransPoint.Speed, ExposureTime(s.CurrentSpeed()),
		gosu.MusicVolume*100, gosu.EffectVolume*100,
		gosu.Offset))
}

// 1 pixel is 1 millisecond.
// Todo: Separate NoteHeight / 2 at piano mode
func ExposureTime(speedScale float64) float64 {
	return (screenSizeX - HitPosition) / speedScale
}
func (s *ScenePlay) UpdateTransPoint() {
	s.TransPoint = s.TransPoint.FetchByTime(s.Now)
}
func (s ScenePlay) Speed()                { s.CurrentSpeed() }
func (s ScenePlay) CurrentSpeed() float64 { return s.TransPoint.Speed * s.SpeedScale }

var DefaultSampleNames = [2][2]string{
	{"red-regular", "red-big"},
	{"blue-regular", "blue-big"},
}
