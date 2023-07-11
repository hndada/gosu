package piano

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

type SceneModePlay struct {
	Chart *Chart
	mode.Timer
	audios.MusicPlayer
	audios.SoundPlayer
	input.KeyLogger
	paused bool

	Dynamic    *mode.Dynamic
	SpeedScale float64
	Cursor     float64
	Scorer

	Sound        []byte
	Field        FieldDrawer
	Bar          BarDrawer
	Note         []NoteDrawer
	Keys         []KeyDrawer
	KeyLighting  []KeyLightingDrawer
	Hint         HintDrawer
	HitLighting  []HitLightingDrawer
	HoldLighting []HoldLightingDrawer
	Judgment     JudgmentDrawer
	Score        mode.ScoreDrawer
	Combo        mode.ComboDrawer
	// Meter        mode.MeterDrawer
}

func NewSceneModePlay(fsys fs.FS, cname string, mods Mods, rf *osr.Format) (s *SceneModePlay, err error) {
	s = new(SceneModePlay)
	s.Chart, err = NewChart(fsys, cname)
	if err != nil {
		return
	}
	c := s.Chart
	s.Timer = mode.NewTimer(c.Duration(), S.offset, TPS)
	s.MusicPlayer, err = audios.NewMusicPlayer(fsys, c.MusicFilename)
	if err != nil {
		return
	}
	s.SoundPlayer = audios.NewSoundPlayer(fsys, S.volumeSound)
	s.KeyLogger = input.NewKeyLogger(S.KeySettings[c.KeyCount])
	if rf != nil {
		s.KeyLogger.FetchPressed = NewReplayListener(rf, c.KeyCount, &s.Timer)
	}

	s.Dynamic = c.Dynamics[0]
	s.SpeedScale = 1
	s.Cursor = float64(s.Now) * s.SpeedScale
	s.Scorer = NewScorer(c)

	skin, ok := UserSkins.Skins[c.KeyCount]
	if !ok {
		UserSkins.loadSkin(c.KeyCount)
		skin = UserSkins.Skins[c.KeyCount]
	}
	s.Sound = skin.Sound

	s.Field = FieldDrawer{
		Sprite: skin.Field,
	}
	s.Bar = BarDrawer{
		Cursor:   s.Cursor,
		Farthest: c.Bars[0],
		Nearest:  c.Bars[0],
		Sprite:   skin.Bar,
	}
	s.Note = make([]NoteDrawer, c.KeyCount)
	s.Keys = make([]KeyDrawer, c.KeyCount)
	s.KeyLighting = make([]KeyLightingDrawer, c.KeyCount)
	s.HitLighting = make([]HitLightingDrawer, c.KeyCount)
	s.HoldLighting = make([]HoldLightingDrawer, c.KeyCount)
	for k := 0; k < c.KeyCount; k++ {
		s.Keys[k] = KeyDrawer{
			Timer:   draws.NewTimer(draws.ToTick(30, TPS), 0),
			Sprites: skin.Key[k],
		}
		s.Note[k] = NoteDrawer{
			Timer:    draws.NewTimer(0, draws.ToTick(400, TPS)), // Todo: make it BPM-dependent?
			Cursor:   s.Cursor,
			Farthest: s.Staged[k],
			Nearest:  s.Staged[k],
			Sprites:  skin.Note[k],
		}
		s.KeyLighting[k] = KeyLightingDrawer{
			Timer:  draws.NewTimer(draws.ToTick(30, TPS), 0),
			Sprite: skin.KeyLighting[k],
			Color:  S.keyLightingColors[KeyTypes[c.KeyCount][k]],
		}
		s.HitLighting[k] = HitLightingDrawer{
			Timer:   draws.NewTimer(draws.ToTick(150, TPS), draws.ToTick(150, TPS)),
			Sprites: skin.HitLighting[k],
		}
		s.HoldLighting[k] = HoldLightingDrawer{
			Timer:   draws.NewTimer(0, draws.ToTick(300, TPS)),
			Sprites: skin.HoldLighting[k],
		}
	}
	s.Hint = HintDrawer{
		Sprite: skin.Hint,
	}
	s.Judgment = NewJudgmentDrawer(skin.Judgment[:])
	s.Score = mode.NewScoreDrawer(skin.Score[:])
	s.Combo = mode.ComboDrawer{
		Timer:      draws.NewTimer(draws.ToTick(2000, TPS), 0),
		DigitWidth: skin.Combo[0].W(),
		DigitGap:   TheSettings.ComboDigitGap,
		Bounce:     0.85,
		Sprites:    skin.Combo,
	}
	// s.Meter = mode.NewMeterDrawer(Judgments, JudgmentColors)
	return s, nil
}

func (s *SceneModePlay) PlayPause() {
	if s.paused {
		s.MusicPlayer.Play()
	} else {
		s.MusicPlayer.Pause()
	}
	s.paused = !s.paused
}
func (s *SceneModePlay) Update() any {
	if !s.paused {
		defer s.Ticker()
	}

	if s.Now == 0+s.Offset {
		s.MusicPlayer.Play()
	}

	var kas []input.KeyboardAction
	for _, ka := range kas {
		s.Scorer.Check(ka)
		for k, n := range s.Staged {
			a := ka.Action[k]
			if n.Type != Tail && a == input.Hit {
				vol := s.Dynamic.Volume
				scale := TheSettings.SoundVolume
				n.Sample.Play(vol, scale)
			}
		}
	}

	s.Bar.Update(s.Cursor)
	for k := 0; k < s.Chart.KeyCount; k++ {
		holding := false
		if s.Staged[k] != nil {
			holding = s.Staged[k].Type == Tail && s.Pressed[k]
		}
		s.Note[k].Update(s.Cursor, holding)
		s.Keys[k].Update(s.Pressed[k])
		s.KeyLighting[k].Update(s.Pressed[k])
		s.HitLighting[k].Update(hits[k])
		s.HoldLighting[k].Update(holding)
	}
	s.Judgment.Update(s.Scorer.worstJudgment)
	s.Score.Update(s.Scorer.Score)
	s.Combo.Update(s.Scorer.Combo)
	// s.Meter.Update()

	// Changed speed should be applied after positions are calculated.
	s.UpdateDynamic()
	s.UpdateCursor()
	return nil
}
func (s SceneModePlay) Speed() float64 { return s.Dynamic.Speed * s.speedScale }
func (s *SceneModePlay) UpdateCursor() {
	duration := float64(s.Now - s.Dynamic.Time)
	s.Cursor = s.Dynamic.Position + duration*s.Speed()
}
func (s *SceneModePlay) UpdateDynamic() {
	dy := s.Dynamic
	for dy.Next != nil && s.Now().Milliseconds() >= dy.Next.Time {
		dy = dy.Next
	}
	s.Dynamic = dy
}

func (s SceneModePlay) Draw(screen draws.Image) {
	s.Field.Draw(screen)
	s.Bar.Draw(screen)
	s.Hint.Draw(screen)
	for k := 0; k < s.Chart.KeyCount; k++ {
		s.Note[k].Draw(screen)
		s.Keys[k].Sprites[0].Draw(screen, draws.Op{})
		s.Keys[k].Draw(screen)
		s.KeyLighting[k].Draw(screen)
	}
	for k := 0; k < s.Chart.KeyCount; k++ {
		s.HitLighting[k].Draw(screen)
		s.HoldLighting[k].Draw(screen)
	}
	s.Judgment.Draw(screen)
	s.Score.Draw(screen)
	s.Combo.Draw(screen)
	// s.Meter.Draw(screen)
}

func (s SceneModePlay) Finish() any {
	s.MusicPlayer.Close()
	s.Scorer.Score += 0.01 // To make sure max score is reachable.
	return s.Scorer
}

func (s SceneModePlay) SetMusicVolume(v float64) {
	TheSettings.MusicVolume = v
	s.MusicPlayer.SetVolume(vol)
}
func (s SceneModePlay) SetSoundVolume(v float64) {
	TheSettings.SoundVolume = v
}

// SetOffset(int64)

// Need to re-calculate positions when Speed has changed.
func (s *SceneModePlay) SetSpeedScale() {
	c := s.Chart
	old := s.SpeedScale
	new := TheSettings.SpeedScale
	s.Cursor *= new / old
	for _, dy := range c.Dynamics {
		dy.Position *= new / old
	}
	for _, n := range c.Notes {
		n.Position *= new / old
	}
	for _, b := range c.Bars {
		b.Position *= new / old
	}
	s.SpeedScale = new
}

func (s SceneModePlay) DebugPrint(screen draws.Image) {
	var scorer Scorer

	fps := fmt.Sprintf("FPS: %.2f\n", ebiten.ActualFPS())
	tps := fmt.Sprintf("TPS: %.2f\n", ebiten.ActualTPS())
	time := fmt.Sprintf("Time: %.3fs/%.0fs\n", float64(s.Now)/1000, float64(s.Chart.Duration())/1000)

	score := fmt.Sprintf("Score: %.0f \n", scorer.Score)
	combo := fmt.Sprintf("Combo: %d\n", scorer.Combo)
	flow := fmt.Sprintf("Flow: %.2f%%\n", scorer.Flow/MaxFlow*100)
	acc := fmt.Sprintf("Acc: %.2f%%\n", scorer.Acc/MaxAcc*100)
	judgmentCount := fmt.Sprintf("Judgment counts: %v\n", scorer.JudgmentCounts)

	speedScale := fmt.Sprintf("Speed scale (Z/X): %.0f (x%.2f)\n", s.SpeedScale, s.Dynamic.Speed)
	exposureTime := fmt.Sprintf("(Exposure time: %.fms)\n", ExposureTime(s.Speed()))

	musicVolume := fmt.Sprintf("Music volume (Ctrl+ Left/Right): %.0f%%\n", TheSettings.MusicVolume*100)
	soundVolume := fmt.Sprintf("Sound volume (Alt+ Left/Right): %.0f%%\n", TheSettings.SoundVolume*100)
	offset := fmt.Sprintf("Offset (Shift+ Left/Right): %dms\n", s.Offset)

	exit := "Press ESC to back to choose a song.\n"
	pause := "Press TAB to pause.\n"

	ebitenutil.DebugPrint(screen.Image, fps+tps+time+"\n"+
		score+combo+flow+acc+judgmentCount+"\n"+
		speedScale+exposureTime+"\n"+
		musicVolume+soundVolume+offset+"\n"+
		exit+pause,
	)
}

// 1 pixel is 1 millisecond.
func ExposureTime(speed float64) float64 {
	return TheSettings.HitPosition / speed
}
