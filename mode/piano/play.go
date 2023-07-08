package piano

import (
	"encoding/csv"
	"fmt"
	"io/fs"
	"os"
	"strconv"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hndada/gosu/audios"
	"github.com/hndada/gosu/draws"
	"github.com/hndada/gosu/format/osr"
	"github.com/hndada/gosu/input"
	"github.com/hndada/gosu/mode"
)

var silent = true
var backgroundRedMode = true

func init() {
	if silent {
		mode.S.VolumeSound = 0
	}
}

// ScenePlay: struct, PlayScene: function
// The skin may be applied some custom settings: on/off some sprites
type ScenePlay struct {
	Chart *Chart
	mode.Timer
	audios.MusicPlayer
	audios.SoundPlayer
	input.KeyLogger
	paused bool

	*mode.TransPoint
	speedScale float64
	Cursor     float64
	Staged     []*Note
	mode.Scorer
	// // Todo: merge into mode.Scorer
	// NoteCount    int
	// MaxNoteCount int

	Sound         []byte
	Background    mode.BackgroundDrawer
	BackgroundRed BackgroundRedDrawer
	Field         FieldDrawer
	Bar           BarDrawer
	Note          []NoteDrawer
	Keys          []KeyDrawer
	KeyLighting   []KeyLightingDrawer
	Hint          HintDrawer
	HitLighting   []HitLightingDrawer
	HoldLighting  []HoldLightingDrawer
	Judgment      JudgmentDrawer
	Score         mode.ScoreDrawer
	Combo         mode.ComboDrawer
	Meter         mode.MeterDrawer

	// For HCI experiments
	Logs       []Log
	offsetMode bool
}
type Log struct {
	Time   int64
	Key    int
	Offset int64
}

func NewScenePlay(fsys fs.FS, cname string, mods interface{}, rf *osr.Format) (s *ScenePlay, err error) {
	s = new(ScenePlay)
	s.Chart, err = NewChart(fsys, cname)
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
	s.SoundPlayer = audios.NewSoundPlayer(fsys, S.volumeSound)
	s.KeyLogger = input.NewKeyLogger(S.KeySettings[c.KeyCount])
	if rf != nil {
		s.KeyLogger.FetchPressed = NewReplayListener(rf, c.KeyCount, &s.Timer)
	}

	s.TransPoint = c.TransPoints[0]
	s.speedScale = 1
	s.Cursor = float64(s.Now) * s.speedScale
	s.SetSpeed()
	s.Scorer = mode.NewScorer(c.ScoreFactors)
	// s.MaxNoteCount = len(c.Notes)
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
	skin, ok := UserSkins.Skins[c.KeyMode]
	if !ok {
		UserSkins.loadSkin(c.KeyMode)
		skin = UserSkins.Skins[c.KeyMode]
	}
	s.Sound = skin.Sound
	s.Background = mode.BackgroundDrawer{
		Sprite: mode.NewBackground(fsys, c.ImageFilename),
	}
	if !s.Background.Sprite.IsValid() {
		s.Background.Sprite = skin.DefaultBackground
	}
	s.BackgroundRed = NewBackgroundRedDrawer() // HCI
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
		DigitGap:   S.ComboDigitGap,
		Bounce:     0.85,
		Sprites:    skin.Combo,
	}
	s.Meter = mode.NewMeterDrawer(Judgments, JudgmentColors)

	// HCI
	s.offsetMode = true
	//if len(s.Chart.Notes) < 20 {
	// s.offsetMode = true
	// }
	return s, nil
}

// Farther note has larger position. Tail's Position is always larger than Head's.
// Need to re-calculate positions when Speed has changed.
func (s *ScenePlay) SetSpeed() {
	c := s.Chart
	old := s.speedScale
	new := S.SpeedScale
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
	if s.Now == 0+s.Offset {
		s.MusicPlayer.Play()
	}
	// if p.Now == 150+p.Offset {
	// 	p.Player.Seek(time.Duration(150) * time.Millisecond)
	// }

	// HCI
	// if s.offsetMode && s.Staged[3] != nil && s.Now > s.Staged[3].Time {
	// s.Staged[3].passed = true
	// }
	var passed bool
	for k, staged := range s.Staged {
		if staged == nil {
			continue
		}
		if s.Now > staged.Time {
			s.Staged[k].passed = true
		}
		if s.Now-staged.Time < 3 && s.Now-staged.Time > -3 {
			passed = true
		}
	}
	// HCI
	if passed {
		s.BackgroundRed.Update(true)
	} else {
		s.BackgroundRed.Update(false)
	}
	// HCI
	// It might take several tries since Update tick is too short.
	if ebiten.IsKeyPressed(ebiten.KeyHome) {
		if backgroundRedMode {
			backgroundRedMode = false
		} else {
			backgroundRedMode = true
		}
	}

	if vol := *S.volumeMusic; S.VolumeMusic != vol {
		S.VolumeMusic = vol
		s.MusicPlayer.SetVolume(vol)
	}

	s.LastPressed = s.Pressed
	s.Pressed = s.FetchPressed()
	var worst mode.Judgment
	hits := make([]bool, s.Chart.KeyCount)

	for _, n := range s.Staged {
		if n == nil {
			continue
		}
		if n.Type != Tail && s.KeyAction(n.Key) == input.Hit {
			s.PlaySample(n)
		}
		td := n.Time - s.Now      // Time difference. A negative value infers late hit
		td += mode.S.DelayedJudge // For HCI experiment
		if n.Marked {
			if n.Type != Tail {
				return fmt.Errorf("non-Tail note has not flushed")
			}
			if td < Miss.Window { // Keep Tail staged until near ends.
				s.Staged[n.Key] = n.Next
			}
			continue
		}
		if j := Judge(n.Type, s.KeyAction(n.Key), td); j.Window != 0 {
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

			// For HCI experiments
			s.Logs = append(s.Logs, Log{
				Time:   n.Time,
				Key:    n.Key,
				Offset: td,
			})
		}
	}
	for k := 0; k < s.Chart.KeyCount; k++ {
		if s.KeyAction(k) == input.Hit {
			vol2 := s.TransPoint.Volume
			p := audios.Context.NewPlayerFromBytes(s.Sound)
			p.SetVolume((*S.volumeSound) * vol2)
			p.Play()
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
	s.Judgment.Update(worst)
	// s.Score.Update(s.LinearScore())
	s.Score.Update(s.Scores[mode.Total])
	s.Combo.Update(s.Scorer.Combo)
	s.Meter.Update()

	// Changed speed should be applied after positions are calculated.
	s.UpdateTransPoint()
	s.UpdateCursor()
	if S.SpeedScale != s.speedScale {
		s.SetSpeed()
	}
	return nil
}
func (s ScenePlay) Finish() any {
	s.MusicPlayer.Close()
	s.outputLog()
	return s.NewResult(s.Chart.MD5)
}

// For HCI experiments
func (s ScenePlay) outputLog() {
	// Create a file where the CSV data can be saved

	fname := fmt.Sprintf("log/%s[%s]_sp%3d_hp%3d_of%3d_ks%3d_dj%3d_%s.csv",
		s.Chart.MusicName, s.Chart.ChartName, int(S.SpeedScale*100), int(S.HitPosition), s.Offset, mode.S.DelayedJudge,
		int(mode.S.VolumeSound*100),
		time.Now().Format("2006-01-02_15-04-05"))
	// create log directory if not exists
	if _, err := os.Stat("log"); os.IsNotExist(err) {
		os.Mkdir("log", 0744)
	}
	file, err := os.Create(fname)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	writer.Write([]string{"Time", "Key", "Offset"})
	for _, n := range s.Logs {
		t := strconv.Itoa(int(n.Time))
		k := strconv.Itoa(int(n.Key))
		o := strconv.Itoa(int(n.Offset))
		writer.Write([]string{t, k, o})
	}
	writer.Flush()
}
func (s *ScenePlay) UpdateCursor() {
	duration := float64(s.Now - s.TransPoint.Time)
	s.Cursor = s.TransPoint.Position + duration*s.Speed()
}
func (s ScenePlay) Speed() float64 { return s.TransPoint.Speed * s.speedScale }
func (s *ScenePlay) UpdateTransPoint() {
	tp := s.TransPoint
	for tp.Next != nil && s.Now().Milliseconds() >= tp.Next.Time {
		tp = tp.Next
	}
	s.TransPoint = tp
}
func (s ScenePlay) PlaySample(n *Note) {
	name := n.Sample.Name
	if name == "" {
		return
	}
	vol2 := n.Sample.Volume
	if vol2 == 0 {
		vol2 = s.TransPoint.Volume
	}
	s.SoundPlayer.Play(name, vol2)
}

func (s ScenePlay) Draw(screen draws.Image) {
	s.Background.Draw(screen)
	if backgroundRedMode {
		s.BackgroundRed.Draw(screen)
	}
	s.Field.Draw(screen)
	s.Bar.Draw(screen)
	s.Hint.Draw(screen)
	for k := 0; k < s.Chart.KeyCount; k++ {
		s.Note[k].Draw(screen)
		if silent {
		} else {
			s.Keys[k].Sprites[0].Draw(screen, draws.Op{})
			s.Keys[k].Draw(screen)
			s.KeyLighting[k].Draw(screen)
		}
	}
	if *S.debugPrint {
		s.DebugPrint(screen)
	}
	if silent {
		return
	}
	for k := 0; k < s.Chart.KeyCount; k++ {
		s.HitLighting[k].Draw(screen)
		s.HoldLighting[k].Draw(screen)
	}
	s.Judgment.Draw(screen)
	s.Score.Draw(screen)
	s.Combo.Draw(screen)
	s.Meter.Draw(screen)
}

func (s ScenePlay) DebugPrint(screen draws.Image) {
	ebitenutil.DebugPrint(screen.Image, fmt.Sprintf(
		"FPS: %.2f\nTPS: %.2f\nTime: %.3fs/%.0fs\n\n"+
			"Score: %.0f | %.0f \nFlow: %.0f/100\nCombo: %d\n\n"+
			"Flow rate: %.2f%%\nAccuracy: %.2f%%\nExtra: %.2f%%\nJudgment counts: %v\n\n"+
			"Speed scale (Z/X): %.0f (x%.2f)\n(Exposure time: %.fms)\n\n"+
			"Music volume (Ctrl+ Left/Right): %.0f%%\nSound volume (Alt+ Left/Right): %.0f%%\n\n"+
			"Press ESC to select a song.\nPress TAB to pause.\n\n"+
			"Offset (Shift+ Left/Right): %dms\n"+
			"Delayed judge (F9/F10): %vms\n"+ // for HCI experiment
			"Debug print (Ctrl+D): %v\n",
		ebiten.ActualFPS(), ebiten.ActualTPS(), float64(s.Now)/1000, float64(s.Chart.Duration())/1000,
		s.Scores[mode.Total], s.ScoreBounds[mode.Total], s.Flow*100, s.Scorer.Combo,
		s.Ratios[0]*100, s.Ratios[1]*100, s.Ratios[2]*100, s.JudgmentCounts,
		S.SpeedScale*100, s.TransPoint.Speed, ExposureTime(s.Speed()),
		*S.volumeMusic*100, *S.volumeSound*100,
		*S.offset,
		*S.delayedJudge,
		*S.debugPrint))
}
