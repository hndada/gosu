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
type ScenePlay struct {
	Chart *Chart
	audios.Timer
	audios.MusicPlayer
	// scene.SoundPlayer
	input.KeyLogger

	*mode.TransPoint
	SpeedScale float64
	Cursor     float64
	Staged     []*Note
	mode.Scorer

	Skin                // The skin may be applied some custom settings: on/off some sprites
	BackgroundDrawer    mode.BackgroundDrawer
	FieldDrawer         FieldDrawer
	BarDrawer           BarDrawer
	NoteDrawers         []NoteDrawer
	KeyDrawers          []KeyDrawer
	KeyLightingDrawers  []KeyLightingDrawer
	HintDrawer          HintDrawer
	HitLightingDrawers  []HitLightingDrawer
	HoldLightingDrawers []HoldLightingDrawer
	JudgmentDrawer      JudgmentDrawer
	ScoreDrawer         mode.ScoreDrawer
	ComboDrawer         mode.NumberDrawer
	MeterDrawer         mode.MeterDrawer
}

func NewScenePlay(fsys fs.FS, cname string, mods interface{}, rf *osr.Format) (s *ScenePlay, err error) {
	s = new(ScenePlay)
	s.Chart, err = NewChart(fsys, cname)
	if err != nil {
		return
	}
	c := s.Chart
	ebiten.SetWindowTitle(c.WindowTitle())
	keyCount := c.KeyCount & ScratchMask
	s.Timer = audios.NewTimer(c.Duration(), &mode.Offset)
	s.MusicPlayer, err = audios.NewMusicPlayer(fsys, c.MusicFilename, &s.Timer, &mode.VolumeMusic)
	if err != nil {
		return
	}
	// s.SoundPlayer = scene.NewSoundPlayer(scene.VolumeSoundHandler)
	// for _, n := range c.Notes {
	// 	if path, ok := n.Sample.Path(cpath); ok {
	// 		_ = s.Sounds.Register(path)
	// 	}
	// }
	s.KeyLogger = input.NewKeyLogger(KeySettings[keyCount])
	if rf != nil {
		s.KeyLogger.FetchPressed = NewReplayListener(rf, keyCount, &s.Timer)
	}

	s.TransPoint = c.TransPoints[0]
	s.SpeedScale = 1
	s.Cursor = float64(s.Now) * s.SpeedScale
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
	s.BackgroundDrawer = mode.BackgroundDrawer{
		Brightness: &mode.BackgroundBrightness,
		Sprite:     mode.DefaultBackground,
	}
	if bg := mode.NewBackground(fsys, c.ImageFilename); bg.IsValid() {
		s.BackgroundDrawer.Sprite = bg
	}
	s.FieldDrawer = FieldDrawer{
		Sprite: s.FieldSprite,
	}
	// s.StageDrawer = StageDrawer{
	// 	FieldSprite: s.FieldSprite,
	// 	HintSprite:  s.HintSprite,
	// }
	s.BarDrawer = BarDrawer{
		Cursor:   s.Cursor,
		Farthest: c.Bars[0],
		Nearest:  c.Bars[0],
		Sprite:   s.BarSprite,
	}
	s.NoteDrawers = make([]NoteDrawer, keyCount)
	s.KeyDrawers = make([]KeyDrawer, keyCount)
	s.KeyLightingDrawers = make([]KeyLightingDrawer, keyCount)
	s.HitLightingDrawers = make([]HitLightingDrawer, keyCount)
	s.HoldLightingDrawers = make([]HoldLightingDrawer, keyCount)
	for k := 0; k < keyCount; k++ {
		s.NoteDrawers[k] = NoteDrawer{
			Timer:    draws.NewTimer(0, draws.ToTick(400)), // Todo: make it BPM-dependent?
			Cursor:   s.Cursor,
			Farthest: s.Staged[k],
			Nearest:  s.Staged[k],
			Sprites:  s.NoteSprites[k],
		}
		s.KeyDrawers[k] = KeyDrawer{
			Timer:   draws.NewTimer(draws.ToTick(30), 0),
			Sprites: s.KeySprites[k],
		}
		s.KeyLightingDrawers[k] = KeyLightingDrawer{
			Timer:  draws.NewTimer(draws.ToTick(30), 0),
			Sprite: s.KeyLightingSprites[k],
		}
		s.HitLightingDrawers[k] = HitLightingDrawer{
			Timer:   draws.NewTimer(draws.ToTick(150), draws.ToTick(150)),
			Sprites: s.HitLightingSprites[k],
		}
		s.HoldLightingDrawers[k] = HoldLightingDrawer{
			Timer:   draws.NewTimer(0, draws.ToTick(250)),
			Sprites: s.HoldLightingSprites[k],
		}
	}
	s.HintDrawer = HintDrawer{
		Sprite: s.HintSprite,
	}
	s.JudgmentDrawer = NewJudgmentDrawer()
	s.ScoreDrawer = mode.NewScoreDrawer()
	s.ComboDrawer = mode.NumberDrawer{
		Timer:      draws.NewTimer(draws.ToTick(2000), 0),
		DigitWidth: s.ComboSprites[0].Size().X,
		DigitGap:   ComboDigitGap,
		Bounce:     0.85,
		Sprites:    s.ComboSprites,
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
			s.MeterDrawer.AddMark(int(td), kind)
			if !j.Is(Miss) && n.Type != Head {
				hits[n.Key] = true
			}
		}
	}
	s.BarDrawer.Update(s.Cursor)
	for k := 0; k < s.Chart.KeyCount; k++ {
		holding := false
		if s.Staged[k] != nil {
			holding = s.Staged[k].Type == Tail && s.Pressed[k]
		}
		s.NoteDrawers[k].Update(s.Cursor, holding)
		s.KeyDrawers[k].Update(s.Pressed[k])
		s.KeyLightingDrawers[k].Update(s.Pressed[k])
		s.HitLightingDrawers[k].Update(hits[k])
		s.HoldLightingDrawers[k].Update(holding)
	}
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
func (s ScenePlay) Draw(screen draws.Image) {
	s.BackgroundDrawer.Draw(screen)
	s.FieldDrawer.Draw(screen)
	s.BarDrawer.Draw(screen)
	for k := 0; k < s.Chart.KeyCount; k++ {
		s.NoteDrawers[k].Draw(screen)
		s.KeyDrawers[k].Draw(screen)
		s.KeyLightingDrawers[k].Draw(screen)
	}
	s.HintDrawer.Draw(screen)
	for k := 0; k < s.Chart.KeyCount; k++ {
		s.HitLightingDrawers[k].Draw(screen)
		s.HoldLightingDrawers[k].Draw(screen)
	}
	s.JudgmentDrawer.Draw(screen)
	s.ScoreDrawer.Draw(screen)
	s.ComboDrawer.Draw(screen)
	s.MeterDrawer.Draw(screen)
	s.DebugPrint(screen)
}

func (s ScenePlay) DebugPrint(screen draws.Image) {
	ebitenutil.DebugPrint(screen.Image, fmt.Sprintf(
		"FPS: %.2f\nTPS: %.2f\nTime: %.3fs/%.0fs\n\n"+
			"Score: %.0f | %.0f \nFlow: %.0f/100\nCombo: %d\n\n"+
			"Flow rate: %.2f%%\nAccuracy: %.2f%%\nExtra: %.2f%%\nJudgment counts: %v\n\n"+
			"Speed scale (Z/X): %.0f (x%.2f)\n(Exposure time: %.fms)\n\n"+
			"Music volume (Alt+ Left/Right): %.0f%%\nSound volume (Ctrl+ Left/Right): %.0f%%\n\n"+
			"Press ESC to select a song.\nPress TAB to pause.\n\n"+
			"Offset (Shift+ Left/Right): %dms\n",
		ebiten.ActualFPS(), ebiten.ActualTPS(), float64(s.Now)/1000, float64(s.Chart.Duration())/1000,
		s.Scores[mode.Total], s.ScoreBounds[mode.Total], s.Flow*100, s.Combo,
		s.Ratios[0]*100, s.Ratios[1]*100, s.Ratios[2]*100, s.JudgmentCounts,
		s.SpeedScale*100, s.TransPoint.Speed, ExposureTime(s.Speed()),
		mode.VolumeMusic*100, mode.VolumeSound*100,
		mode.Offset))
}

// 1 pixel is 1 millisecond.
func ExposureTime(speed float64) float64 { return HitPosition / speed }
func (s ScenePlay) Speed() float64       { return s.TransPoint.Speed * s.SpeedScale }

// Supposes one current TransPoint can increment cursor precisely.
func (s *ScenePlay) UpdateCursor() {
	duration := float64(s.Now - s.TransPoint.Time)
	s.Cursor = s.TransPoint.Position + duration*s.Speed()
}
func (s *ScenePlay) UpdateTransPoint() {
	s.TransPoint = s.TransPoint.FetchByTime(s.Now)
}
