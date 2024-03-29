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

// Todo: support custom hitsound?
type ScenePlay struct {
	Chart *Chart
	mode.Timer
	audios.MusicPlayer
	input.KeyLogger
	KeyActions [2]int
	paused     bool

	*mode.Dynamic
	speedScale         float64
	StagedNote         *Note
	StagedDot          *Dot
	StagedShake        *Note
	LastHitTimes       [4]int64      // For judging big note.
	StagedJudgment     mode.Judgment // For judging big note.
	StagedJudgmentTime int64
	ShakeWaitingColor  int
	mode.Scorer

	DrumSound  [2][2][]byte
	Background mode.BackgroundDrawer
	Stage      StageDrawer
	Bar        BarDrawer
	Judgment   JudgmentDrawer
	Shake      ShakeDrawer
	Roll       RollDrawer
	Note       NoteDrawer
	Key        KeyDrawer
	Dancer     DancerDrawer
	Score      mode.ScoreDrawer
	Combo      mode.ComboDrawer
	Meter      mode.MeterDrawer
}

// Todo: actual auto replay generator for gimmick charts
// Todo: support mods: show Piano's ScenePlay during Drum's ScenePlay
func NewScenePlay(fsys fs.FS, cname string, mods interface{}, rf *osr.Format) (s *ScenePlay, err error) {
	s = new(ScenePlay)
	s.Chart, err = LoadChart(fsys, cname)
	if err != nil {
		return
	}
	c := s.Chart
	ebiten.SetWindowTitle(c.WindowTitle())
	s.Timer = mode.NewTimer(c.Duration(), S.offset, TPS)
	s.MusicPlayer, err = audios.NewMusicPlayer(fsys, c.MusicFilename)
	if err != nil {
		return
	}
	s.KeyLogger = input.NewKeyLogger(S.KeySettings[4][:])
	if rf != nil {
		s.KeyLogger.ReadPressed = NewReplayPlayer(rf, &s.Timer)
	}

	s.Dynamic = c.Dynamics[0]
	s.speedScale = 1
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

	skin := UserSkin
	s.DrumSound = skin.DrumSound
	s.Background = mode.BackgroundDrawer{
		Sprite: mode.NewBackground(fsys, c.BackgroundFilename),
	}
	if !s.Background.Sprite.IsValid() {
		s.Background.Sprite = skin.defaultBackground
	}
	s.Stage = StageDrawer{
		Timer:     draws.NewTimer(draws.ToTick(150, TPS), 0),
		Highlight: false, //s.Highlight,
		Field:     skin.Field,
		Hint:      skin.Hint,
	}
	s.Bar = BarDrawer{
		Time:   s.Now,
		Bars:   c.Bars,
		Sprite: skin.Bar,
	}
	s.Judgment = JudgmentDrawer{
		Timer:      draws.NewTimer(draws.ToTick(250, TPS), draws.ToTick(250, TPS)),
		Animations: skin.Judgment,
	}
	s.Shake = ShakeDrawer{
		Timer:  draws.NewTimer(200, 0),
		Time:   s.Now,
		Staged: s.StagedShake,
		Shake:  skin.Shake,
	}
	s.Roll = RollDrawer{
		Time:      s.Now,
		Rolls:     c.Rolls,
		Dots:      c.Dots,
		Head:      skin.Head,
		Tail:      skin.Tail,
		Body:      skin.Body,
		DotSprite: skin.Dot,
	}
	period := int(60000 / ScaledBPM(s.BPM))
	s.Note = NoteDrawer{
		Timer:   draws.NewTimer(0, period),
		Time:    s.Now,
		Notes:   c.Notes,
		Rolls:   c.Rolls,
		Shakes:  c.Shakes,
		Note:    skin.Note,
		Overlay: skin.Overlay,
	}
	s.Key = KeyDrawer{
		MaxCountdown: draws.ToTick(75, TPS),
		Field:        skin.KeyField,
		Keys:         skin.Key,
	}
	s.Dancer = DancerDrawer{
		Timer:       draws.NewTimer(0, period),
		Time:        s.Now,
		Dancer:      skin.Dancer,
		Mode:        DancerIdle,
		ModeEndTime: s.Now,
	}
	s.Score = mode.NewScoreDrawer(skin.score[:])
	s.Combo = mode.ComboDrawer{
		Timer:      draws.NewTimer(draws.ToTick(2000, TPS), 0),
		DigitWidth: skin.Combo[0].Width(),
		DigitGap:   S.ComboDigitGap,
		Combo:      0,
		Bounce:     1.25,
		Sprites:    skin.Combo,
	}
	s.Meter = mode.NewMeterDrawer(Judgments, JudgmentColors)
	return s, nil
}

// Farther note has larger position. Tail's Position is always larger than Head's.
// Need to re-calculate positions when Speed has changed.
func (s *ScenePlay) SetSpeed() {
	c := s.Chart
	old := s.speedScale
	new := S.SpeedScale
	for _, dy := range c.Dynamics {
		dy.Speed *= new / old
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
	s.speedScale = new
}
func (s *ScenePlay) PlayPause() {
	if s.paused {
		s.MusicPlayer.Play()
	} else {
		s.MusicPlayer.Pause()
	}
	s.paused = !s.paused
}
func (s *ScenePlay) Update() any {
	if !s.paused {
		defer s.Ticker()
	}
	if s.Now == 0+s.MusicOffset {
		s.MusicPlayer.Play()
	}
	if vol := *S.musicVolume; S.MusicVolume != vol {
		S.MusicVolume = vol
		s.MusicPlayer.SetVolume(vol)
	}

	s.LastPressed = s.Pressed
	s.Pressed = s.ReadPressed()
	s.UpdateKeyActions()

	var (
		judgment mode.Judgment
		big      bool
	)
	if s.StagedJudgment.IsValid() {
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
				s.LastHitTimes[key] = -mode.Wait
			}
			td := n.Time - jTime
			s.MarkNote(n, j, false)
			s.Meter.AddMark(int(td), 0)
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
				s.Meter.AddMark(int(td), 0)
				judgment = j
				big = b
				if s.StagedJudgment.IsValid() {
					s.StagedJudgment = mode.Judgment{}
				}
			}
		}
	}
	if n := s.StagedDot; n != nil {
		td := n.Time - s.Now
		if marked := VerdictDot(n, s.KeyActions, td); marked != DotReady {
			s.MarkDot(n, marked)
			s.Meter.AddMark(int(td), 1)
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
	for color, size := range s.KeyActions {
		if size == SizeNone {
			continue
		}
		vol2 := s.Dynamic.Volume
		p := audios.Context.NewPlayerFromBytes(s.DrumSound[color][size])
		p.SetVolume((*S.volumeSound) * vol2)
		p.Play()
	}
	if s.Now >= 0 {
		s.Stage.Update(s.Highlight)
	}
	s.Bar.Update(s.Now)
	s.Judgment.Update(judgment, big)
	s.Shake.Update(s.Now, s.StagedShake)
	s.Roll.Update(s.Now)
	s.Note.Update(s.Now, s.BPM)

	s.Key.Update(s.LastPressed, s.Pressed)
	s.Dancer.Update(s.Now, s.BPM, s.Scorer.Combo, judgment.Is(Miss),
		!judgment.Is(Miss) && judgment.IsValid(), s.Highlight)
	s.Score.Update(s.Scores[mode.Total])
	s.Combo.Update(s.Scorer.Combo)
	s.Meter.Update()

	// Changed speed should be applied after positions are calculated.
	s.UpdateDynamic()
	if s.speedScale != S.SpeedScale {
		s.SetSpeed()
	}
	return nil
}
func (s ScenePlay) Finish() any {
	s.MusicPlayer.Close()
	return s.NewResult(s.Chart.MD5)
}
func (s ScenePlay) Speed() float64 { return s.Dynamic.Speed * s.speedScale }
func (s *ScenePlay) UpdateDynamic() {
	s.Dynamic = s.Dynamic.FetchByTime(s.Now)
}
func (s ScenePlay) Draw(screen draws.Image) {
	// screen.Fill(color.NRGBA{0, 255, 0, 255}) // Chroma-key
	s.Background.Draw(screen)
	s.Stage.Draw(screen)
	s.Bar.Draw(screen)
	s.Judgment.Draw(screen)
	s.Shake.Draw(screen)
	s.Roll.Draw(screen)
	s.Note.Draw(screen)
	s.Key.Draw(screen)
	s.Dancer.Draw(screen)
	s.Score.Draw(screen)
	s.Combo.Draw(screen)
	s.Meter.Draw(screen)
	s.DebugPrint(screen)
}

func (s ScenePlay) DebugPrint(screen draws.Image) {
	y := S.FieldPosition + S.FieldHeight/2
	ebitenutil.DebugPrintAt(screen.Image, fmt.Sprintf(
		// "\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n"+
		"Press ESC to select a song.\nPress TAB to pause.\n\n"+
			"FPS: %.2f\nTPS: %.2f\nTime: %.3fs/%.0fs\n\n"+
			"Score: %.0f | %.0f \nFlow: %.0f/100\nCombo: %d\n\n"+
			"Flow rate: %.2f%%\nAccuracy: %.2f%%\nExtra: %.2f%%\n"+
			"Judgment counts: %v\nPartial counts: %v\nTick counts: %v\n\n"+
			"Speed scale (PageUp/Down): %.0f (x%.2f)\n(Exposure time: %.fms)\n\n"+
			"Music volume (Ctrl+ Left/Right): %.0f%%\nSound volume (Alt+ Left/Right): %.0f%%\n\n"+
			"MusicOffset (Shift+ Left/Right): %dms\n",
		ebiten.ActualFPS(), ebiten.ActualTPS(), float64(s.Now)/1000, float64(s.Chart.Duration())/1000,
		s.Scores[mode.Total], s.ScoreBounds[mode.Total], s.Flow*100, s.Scorer.Combo,
		s.Ratios[0]*100, s.Ratios[1]*100, s.Ratios[2]*100,
		s.JudgmentCounts[:3], s.JudgmentCounts[3:5], s.JudgmentCounts[5:],
		s.speedScale*100, s.speedScale/s.Dynamic.Speed, ExposureTime(s.Speed()),
		*S.musicVolume*100, *S.volumeSound*100,
		*S.offset), 0, int(y))
}

var DefaultSampleNames = [2][2]string{{"red", "red-big"}, {"blue", "blue-big"}}
