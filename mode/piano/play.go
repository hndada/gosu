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

// ScenePlay: struct, PlayScene: function
// The skin may be applied some custom settings: on/off some sprites
type ScenePlay struct {
	Chart *Chart
	audios.Timer
	audios.MusicPlayer
	// scene.SoundPlayer
	input.KeyLogger

	*mode.TransPoint
	speedScale float64
	Cursor     float64
	Staged     []*Note
	mode.Scorer

	Background    mode.BackgroundDrawer
	Field         FieldDrawer
	Bar           BarDrawer
	Notes         []NoteDrawer
	Keys          []KeyDrawer
	KeyLightings  []KeyLightingDrawer
	Hint          HintDrawer
	HitLightings  []HitLightingDrawer
	HoldLightings []HoldLightingDrawer
	Judgment      JudgmentDrawer
	Score         mode.ScoreDrawer
	Combo         mode.ComboDrawer
	Meter         mode.MeterDrawer
}

func NewScenePlay(fsys fs.FS, cname string, mods interface{}, rf *osr.Format) (s *ScenePlay, err error) {
	S := UserSettings
	s = new(ScenePlay)
	s.Chart, err = NewChart(fsys, cname)
	if err != nil {
		return
	}
	c := s.Chart
	ebiten.SetWindowTitle(c.WindowTitle())
	s.Timer = audios.NewTimer(c.Duration(), S.offset)
	s.MusicPlayer, err = audios.NewMusicPlayer(fsys, c.MusicFilename, &s.Timer, S.volumeMusic)
	if err != nil {
		return
	}
	// s.SoundPlayer = scene.NewSoundPlayer(scene.VolumeSoundHandler)
	// for _, n := range c.Notes {
	// 	if path, ok := n.Sample.Path(cpath); ok {
	// 		_ = s.Sounds.Register(path)
	// 	}
	// }
	s.KeyLogger = input.NewKeyLogger(S.KeySettings[c.KeyCount])
	if rf != nil {
		s.KeyLogger.FetchPressed = NewReplayListener(rf, c.KeyCount, &s.Timer)
	}

	s.TransPoint = c.TransPoints[0]
	s.speedScale = 1
	s.Cursor = float64(s.Now) * s.speedScale
	s.SetSpeed()
	s.Scorer = mode.NewScorer(c.ScoreFactors)
	s.JudgmentCounts = make([]int, len(Judgments))
	// s.Result.FlowMarks = make([]float64, 0, c.Duration()/1000)
	var maxWeight float64
	for _, n := range c.Notes {
		maxWeight += n.Weight()
	}
	for i := range s.MaxWeights {
		s.MaxWeights[i] = maxWeight
	}
	s.Staged = make([]*Note, c.KeyCount)
	for k := range s.Staged {
		for _, n := range c.Notes {
			if k == n.Key {
				s.Staged[n.Key] = n
				break
			}
		}
	}
	PlaySkin.KeyMode = c.KeyMode
	PlaySkin.Load(fsys)
	skin := PlaySkin
	s.Background = mode.BackgroundDrawer{
		Sprite: mode.NewBackground(fsys, c.ImageFilename),
	}
	if !s.Background.Sprite.IsValid() {
		s.Background.Sprite = skin.DefaultBackground
	}
	s.Field = FieldDrawer{
		Sprite: skin.Field,
	}
	s.Bar = BarDrawer{
		Cursor:   s.Cursor,
		Farthest: c.Bars[0],
		Nearest:  c.Bars[0],
		Sprite:   skin.Bar,
	}
	s.Notes = make([]NoteDrawer, c.KeyCount)
	s.Keys = make([]KeyDrawer, c.KeyCount)
	s.KeyLightings = make([]KeyLightingDrawer, c.KeyCount)
	s.HitLightings = make([]HitLightingDrawer, c.KeyCount)
	s.HoldLightings = make([]HoldLightingDrawer, c.KeyCount)
	for k := 0; k < c.KeyCount; k++ {
		s.Keys[k] = KeyDrawer{
			Timer:   draws.NewTimer(draws.ToTick(30, TPS), 0),
			Sprites: skin.Key[k],
		}
		s.Notes[k] = NoteDrawer{
			Timer:    draws.NewTimer(0, draws.ToTick(400, TPS)), // Todo: make it BPM-dependent?
			Cursor:   s.Cursor,
			Farthest: s.Staged[k],
			Nearest:  s.Staged[k],
			Sprites:  skin.Note[k],
		}
		s.KeyLightings[k] = KeyLightingDrawer{
			Timer:  draws.NewTimer(draws.ToTick(30, TPS), 0),
			Sprite: skin.KeyLighting[k],
		}
		s.HitLightings[k] = HitLightingDrawer{
			Timer:   draws.NewTimer(draws.ToTick(150, TPS), draws.ToTick(150, TPS)),
			Sprites: skin.HitLighting[k],
			Color:   S.hitLightingColors[KeyTypes[c.KeyCount][k]],
		}
		s.HoldLightings[k] = HoldLightingDrawer{
			Timer:   draws.NewTimer(0, draws.ToTick(250, TPS)),
			Sprites: skin.HoldLighting[k],
		}
	}
	s.Hint = HintDrawer{
		Sprite: skin.Hint,
	}
	s.Judgment = NewJudgmentDrawer(skin.Judgment[:])
	s.Score = mode.NewScoreDrawer(skin.Score[:])
	s.Combo = mode.ComboDrawer{
		Timer:    draws.NewTimer(draws.ToTick(2000, TPS), 0),
		DigitGap: S.ComboDigitGap,
		Bounce:   0.85,
		Sprites:  skin.Combo,
	}
	s.Meter = mode.NewMeterDrawer(Judgments, JudgmentColors)
	return s, nil
}

// Farther note has larger position. Tail's Position is always larger than Head's.
// Need to re-calculate positions when Speed has changed.
func (s *ScenePlay) SetSpeed() {
	c := s.Chart
	old := s.speedScale
	new := UserSettings.SpeedScale
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
	s.speedScale = new
}

// Todo: apply other values of TransPoint (Volume has finished so far)
// Todo: keep playing music when making SceneResult
func (s *ScenePlay) Update() any {
	defer s.Ticker()
	if s.IsDone() {
		s.MusicPlayer.Close()
		// return scene.PlayToResultArgs{Result: s.NewResult(s.Chart.MD5)}
	}
	s.MusicPlayer.Update()

	s.LastPressed = s.Pressed
	s.Pressed = s.FetchPressed()
	var worst mode.Judgment
	hits := make([]bool, s.Chart.KeyCount)
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
		// 		s.Sounds.PlayWithVolume(name, vol)
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
			s.Meter.AddMark(int(td), kind)
			if !j.Is(Miss) && n.Type != Head {
				hits[n.Key] = true
			}
		}
	}
	s.Bar.Update(s.Cursor)
	for k := 0; k < s.Chart.KeyCount; k++ {
		holding := false
		if s.Staged[k] != nil {
			holding = s.Staged[k].Type == Tail && s.Pressed[k]
		}
		s.Notes[k].Update(s.Cursor, holding)
		s.Keys[k].Update(s.Pressed[k])
		s.KeyLightings[k].Update(s.Pressed[k])
		s.HitLightings[k].Update(hits[k])
		s.HoldLightings[k].Update(holding)
	}
	s.Judgment.Update(worst)
	s.Score.Update(s.Scores[3])
	s.Combo.Update(s.Scorer.Combo)
	s.Meter.Update()

	// Changed speed should be applied after positions are calculated.
	s.UpdateTransPoint()
	s.UpdateCursor()
	if UserSettings.SpeedScale != s.speedScale {
		s.SetSpeed()
	}
	return nil
}
func (s ScenePlay) Draw(screen draws.Image) {
	s.Background.Draw(screen)
	s.Field.Draw(screen)
	s.Bar.Draw(screen)
	for k := 0; k < s.Chart.KeyCount; k++ {
		s.Notes[k].Draw(screen)
		s.Keys[k].Draw(screen)
		s.KeyLightings[k].Draw(screen)
	}
	s.Hint.Draw(screen)
	for k := 0; k < s.Chart.KeyCount; k++ {
		s.HitLightings[k].Draw(screen)
		s.HoldLightings[k].Draw(screen)
	}
	s.Judgment.Draw(screen)
	s.Score.Draw(screen)
	s.Combo.Draw(screen)
	s.Meter.Draw(screen)
	s.DebugPrint(screen)
}

func (s ScenePlay) DebugPrint(screen draws.Image) {
	S := UserSettings
	ebitenutil.DebugPrint(screen.Image, fmt.Sprintf(
		"FPS: %.2f\nTPS: %.2f\nTime: %.3fs/%.0fs\n\n"+
			"Score: %.0f | %.0f \nFlow: %.0f/100\nCombo: %d\n\n"+
			"Flow rate: %.2f%%\nAccuracy: %.2f%%\nExtra: %.2f%%\nJudgment counts: %v\n\n"+
			"Speed scale (Z/X): %.0f (x%.2f)\n(Exposure time: %.fms)\n\n"+
			"Music volume (Alt+ Left/Right): %.0f%%\nSound volume (Ctrl+ Left/Right): %.0f%%\n\n"+
			"Press ESC to select a song.\nPress TAB to pause.\n\n"+
			"Offset (Shift+ Left/Right): %dms\n",
		ebiten.ActualFPS(), ebiten.ActualTPS(), float64(s.Now)/1000, float64(s.Chart.Duration())/1000,
		s.Scores[mode.Total], s.ScoreBounds[mode.Total], s.Flow*100, s.Scorer.Combo,
		s.Ratios[0]*100, s.Ratios[1]*100, s.Ratios[2]*100, s.JudgmentCounts,
		S.SpeedScale*100, s.TransPoint.Speed, ExposureTime(s.Speed()),
		*S.volumeMusic*100, *S.volumeSound*100,
		*S.offset))
}

// 1 pixel is 1 millisecond.
func ExposureTime(speed float64) float64 { return UserSettings.HitPosition / speed }
func (s ScenePlay) Speed() float64       { return s.TransPoint.Speed * s.speedScale }

// Supposes one current TransPoint can increment cursor precisely.
func (s *ScenePlay) UpdateCursor() {
	duration := float64(s.Now - s.TransPoint.Time)
	s.Cursor = s.TransPoint.Position + duration*s.Speed()
}

// Todo: remove it
func (s *ScenePlay) UpdateTransPoint() {
	s.TransPoint = s.TransPoint.FetchByTime(s.Now)
}
