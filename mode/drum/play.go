package drum

import (
	"fmt"
	"io/fs"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hndada/gosu/audios"
	"github.com/hndada/gosu/draws"
	"github.com/hndada/gosu/format/osr"
	"github.com/hndada/gosu/input"
	"github.com/hndada/gosu/mode"
)

type ScenePlay struct {
	Chart        *Chart
	audios.Timer // Todo: MusicPlayer has it
	audios.MusicPlayer
	SoundSoundBytes [2][2][]byte // No custom hitsound at Drum mode.
	input.KeyLogger
	KeyActions [2]int

	*mode.TransPoint
	SpeedScale         float64
	StagedNote         *Note
	StagedDot          *Dot
	StagedShake        *Note
	LastHitTimes       [4]int64      // For judging big note.
	StagedJudgment     mode.Judgment // For judging big note.
	StagedJudgmentTime int64
	ShakeWaitingColor  int
	mode.Scorer

	// Skin may be applied some custom settings: on/off some sprites
	Skin
	BackgroundDrawer mode.BackgroundDrawer
	StageDrawer      StageDrawer
	BarDrawer        BarDrawer
	JudgmentDrawer   JudgmentDrawer

	ShakeDrawer ShakeDrawer
	RollDrawer  RollDrawer
	NoteDrawer  NoteDrawer

	KeyDrawer    KeyDrawer
	DancerDrawer DancerDrawer
	ScoreDrawer  mode.ScoreDrawer
	ComboDrawer  mode.NumberDrawer
	MeterDrawer  mode.MeterDrawer
}

// Todo: actual auto replay generator for gimmick charts
// Todo: support mods: show Piano's ScenePlay during Drum's ScenePlay
func NewScenePlay(fsys fs.FS, cname string, mods interface{}, rf *osr.Format) (s *ScenePlay, err error) {
	s = new(ScenePlay)
	s.Chart, err = NewChart(fsys, cname)
	if err != nil {
		return
	}
	c := s.Chart
	ebiten.SetWindowTitle(c.WindowTitle())
	s.Timer = audios.NewTimer(c.Duration(), &mode.Offset)
	s.MusicPlayer, err = audios.NewMusicPlayer(fsys, c.MusicFilename, &s.Timer, &mode.VolumeMusic)
	if err != nil {
		return
	}
	s.KeyLogger = input.NewKeyLogger(KeySettings[4][:])
	if rf != nil {
		s.KeyLogger.FetchPressed = NewReplayListener(rf, &s.Timer)
	}

	s.TransPoint = c.TransPoints[0]
	s.SpeedScale = 1
	s.SetSpeed()
	s.Scorer = mode.NewScorer(c.ScoreFactors)
	s.JudgmentCounts = make([]int, len(JudgmentCountKinds))
	// s.FlowMarks = make([]float64, 0, c.Duration()/1000)
	for _, n := range c.Notes {
		s.MaxWeights[mode.Flow] += n.Weight()
	}
	s.MaxWeights[mode.Acc] = s.MaxWeights[mode.Flow]
	for _, n := range c.Dots {
		s.MaxWeights[mode.Extra] += n.Weight()
	}
	for _, n := range c.Shakes {
		s.MaxWeights[mode.Extra] += n.Weight()
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
	s.SoundSoundBytes = s.Skin.SoundSoundBytes
	s.BackgroundDrawer = mode.BackgroundDrawer{
		Brightness: &mode.BackgroundBrightness,
		Sprite:     mode.DefaultBackground,
	}
	if bg := mode.NewBackground(fsys, c.ImageFilename); bg.IsValid() {
		s.BackgroundDrawer.Sprite = bg
	}
	s.StageDrawer = StageDrawer{
		Timer:        draws.NewTimer(draws.ToTick(150), 0),
		Highlight:    false, //s.Highlight,
		FieldSprites: s.FieldSprites,
		HintSprites:  s.HintSprites,
	}
	s.BarDrawer = BarDrawer{
		Time:   s.Now,
		Bars:   c.Bars,
		Sprite: s.BarSprite,
	}
	s.JudgmentDrawer = JudgmentDrawer{
		Timer:   draws.NewTimer(draws.ToTick(250), draws.ToTick(250)),
		Sprites: s.JudgmentSprites,
	}
	s.ShakeDrawer = ShakeDrawer{
		Timer:   draws.NewTimer(200, 0),
		Time:    s.Now,
		Staged:  s.StagedShake,
		Sprites: s.ShakeSprites,
	}
	s.RollDrawer = RollDrawer{
		Time:        s.Now,
		Rolls:       c.Rolls,
		Dots:        c.Dots,
		HeadSprites: s.HeadSprites,
		TailSprites: s.TailSprites,
		BodySprites: s.BodySprites,
		DotSprite:   s.DotSprite,
	}
	period := int(60000 / ScaledBPM(s.BPM))
	s.NoteDrawer = NoteDrawer{
		Timer:          draws.NewTimer(0, period),
		Time:           s.Now,
		Notes:          c.Notes,
		Rolls:          c.Rolls,
		Shakes:         c.Shakes,
		NoteSprites:    s.NoteSprites,
		OverlaySprites: s.OverlaySprites,
	}
	s.KeyDrawer = KeyDrawer{
		MaxCountdown: draws.ToTick(75),
		Field:        s.KeyFieldSprite,
		Keys:         s.KeySprites,
	}
	s.DancerDrawer = DancerDrawer{
		Timer:       draws.NewTimer(0, period),
		Time:        s.Now,
		Sprites:     s.DancerSprites,
		Mode:        DancerIdle,
		ModeEndTime: s.Now,
	}
	s.ScoreDrawer = mode.NewScoreDrawer()
	s.ComboDrawer = mode.NumberDrawer{
		Timer:      draws.NewTimer(draws.ToTick(2000), 0),
		Sprites:    s.ComboSprites,
		DigitWidth: s.ComboSprites[0].W(),
		DigitGap:   ComboDigitGap,
		Bounce:     1.25,
	}
	s.MeterDrawer = mode.NewMeterDrawer(Judgments, JudgmentColors)
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
		// return scene.PlayToResultArgs{Result: s.NewResult(s.Chart.MD5)}
	}
	s.MusicPlayer.Update()

	s.LastPressed = s.Pressed
	s.Pressed = s.FetchPressed()
	s.UpdateKeyActions()

	var (
		judgment mode.Judgment
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
				s.LastHitTimes[key] = -audios.Wait
			}
			td := n.Time - jTime
			s.MarkNote(n, j, false)
			s.MeterDrawer.AddMark(int(td), 0)
			judgment = j
			big = false
			s.StagedJudgment = mode.Judgment{}
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
					s.StagedJudgment = mode.Judgment{}
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
		p := audios.Context.NewPlayerFromBytes(s.SoundSoundBytes[i][size])
		p.SetVolume(vol * mode.VolumeSound)
		p.Play()
	}
	if s.Now >= 0 {
		s.StageDrawer.Update(s.Highlight)
	}
	s.BarDrawer.Update(s.Now)
	s.JudgmentDrawer.Update(judgment, big)
	s.ShakeDrawer.Update(s.Now, s.StagedShake)
	s.RollDrawer.Update(s.Now)
	s.NoteDrawer.Update(s.Now, s.BPM)

	s.KeyDrawer.Update(s.LastPressed, s.Pressed)
	s.DancerDrawer.Update(s.Now, s.BPM, s.Combo, judgment.Is(Miss),
		!judgment.Is(Miss) && judgment.Valid(), s.Highlight)
	s.ScoreDrawer.Update(s.Scores[mode.Total])
	s.ComboDrawer.Update(s.Combo)
	s.MeterDrawer.Update()

	// Changed speed should be applied after positions are calculated.
	s.UpdateTransPoint()
	if SpeedScale != s.SpeedScale {
		s.SetSpeed()
	}
	return nil
}
func (s ScenePlay) Draw(screen draws.Image) {
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

func (s ScenePlay) DebugPrint(screen draws.Image) {
	ebitenutil.DebugPrint(screen.Image, fmt.Sprintf(
		"\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n"+
			"Press ESC to select a song.\nPress TAB to pause.\n\n"+
			"FPS: %.2f\nTPS: %.2f\nTime: %.3fs/%.0fs\n\n"+
			"Score: %.0f | %.0f \nFlow: %.0f/100\nCombo: %d\n\n"+
			"Flow rate: %.2f%%\nAccuracy: %.2f%%\nExtra: %.2f%%\n"+
			"Judgment counts: %v\nPartial counts: %v\nTick counts: %v\n\n"+
			"Speed scale (PageUp/Down): %.0f (x%.2f)\n(Exposure time: %.fms)\n\n"+
			"Music volume (Alt+ Left/Right): %.0f%%\nSound volume (Ctrl+ Left/Right): %.0f%%\n\n"+
			"Offset (Shift+ Left/Right): %dms\n",
		ebiten.ActualFPS(), ebiten.ActualTPS(), float64(s.Now)/1000, float64(s.Chart.Duration())/1000,
		s.Scores[mode.Total], s.ScoreBounds[mode.Total], s.Flow*100, s.Combo,
		s.Ratios[0]*100, s.Ratios[1]*100, s.Ratios[2]*100,
		s.JudgmentCounts[:3], s.JudgmentCounts[3:5], s.JudgmentCounts[5:],
		s.SpeedScale*100, s.SpeedScale/s.TransPoint.Speed, ExposureTime(s.Speed()),
		mode.VolumeMusic*100, mode.VolumeSound*100,
		mode.Offset))
}

// 1 pixel is 1 millisecond.
// Todo: Separate NoteHeight / 2 at piano mode
func ExposureTime(speedScale float64) float64 {
	return (ScreenSizeX - HitPosition) / speedScale
}
func (s *ScenePlay) UpdateTransPoint() {
	s.TransPoint = s.TransPoint.FetchByTime(s.Now)
}
func (s ScenePlay) Speed() float64 { return s.TransPoint.Speed * s.SpeedScale }

var DefaultSampleNames = [2][2]string{
	{"red-regular", "red-big"}, {"blue-regular", "blue-big"},
}
