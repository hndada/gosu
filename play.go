package gosu

import (
	"crypto/md5"
	"fmt"
	"os"
	"path/filepath"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hndada/gosu/audios"
	"github.com/hndada/gosu/ctrl"
	"github.com/hndada/gosu/draws"
	"github.com/hndada/gosu/format/osr"
	"github.com/hndada/gosu/input"
)

type ScenePlay interface {
	SetTick(rf *osr.Format)
}

// Fields can be set at outer function call.
// Update cannot be generalized; each scene use template fields in timely manner.
// Unit of time is a millisecond (1ms = 0.001s).
type BaseScenePlay struct {
	// General
	Tick             int
	Chart            *Chart
	EndTime          int64 // EndTime = Duration + WaitAfter
	Mods             Mods
	BackgroundDrawer BackgroundDrawer

	// Speed, BPM, Volume and Highlight
	MainBPM      float64
	SpeedScale   *float64
	SpeedHandler ctrl.F64Handler
	*TransPoint

	// Audio
	// LastVolume  float64
	Volume      *float64
	MusicPlayer *audio.Player
	MusicCloser func() error
	Sounds      audios.SoundMap // A player for sample sound is generated at a place.

	// Input
	FetchPressed func() []bool
	LastPressed  []bool
	Pressed      []bool

	// Note
	Notes  []*Note
	Staged []*Note // Drum has 3 lanes: normal, Roll, Shake.

	// Score
	Result
	NoteWeights  float64
	Combo        int
	Flow         float64
	DelayedScore ctrl.Delayed
	ScoreDrawer  draws.NumberDrawer
	Meter        Meter
}

func NewBaseScenePlay(
	cpath string,
	mode, subMode int, mods Mods,
	rf *osr.Format,
	speedScale *float64,
	keySettings []input.Key,
	bg draws.Sprite,
) (s *BaseScenePlay, err error) {
	waitBefore := MinWaitBefore
	if rf != nil && rf.BufferTime() < waitBefore {
		waitBefore = rf.BufferTime()
	}
	s.Tick = TimeToTick(waitBefore)
	c, err := NewChart(cpath, mode, subMode)
	if err != nil {
		return nil, err
	}
	s.Chart = c
	s.EndTime = c.Duration + WaitAfter
	s.Mods = mods
	s.BackgroundDrawer = BackgroundDrawer{
		Sprite:  bg,
		Dimness: &BackgroundDimness,
	}

	s.MainBPM, _, _ = BPMs(c.TransPoints, c.Duration)
	s.SpeedScale = speedScale
	s.SpeedHandler = NewSpeedHandler(speedScale)
	s.TransPoint = c.TransPoints[0].FetchLatest()

	s.Volume = &Volume
	if path := c.AudioFilename; path != "virtual" && path != "" {
		s.MusicPlayer, s.MusicCloser, err = audios.NewPlayer(path)
		if err != nil {
			return
		}
		s.MusicPlayer.SetVolume(*s.Volume)
	}
	s.Sounds = audios.NewSoundMap(s.Volume)
	for _, n := range c.Notes {
		if n.SampleName == "" {
			continue
		}
		path := filepath.Join(filepath.Dir(cpath), n.SampleName)
		s.Sounds.Register(path)
	}
	if rf != nil {
		s.FetchPressed = NewReplayListener(rf, subMode, waitBefore)
	} else {
		s.FetchPressed = input.NewListener(keySettings)
	}
	switch mode {
	case ModePiano4, ModePiano7:
		s.LastPressed = make([]bool, subMode)
		s.Pressed = make([]bool, subMode)
		for k := range s.Staged {
			for _, n := range c.Notes {
				if k == n.Key {
					s.Staged[n.Key] = n
					break
				}
			}
		}
	case ModeDrum:
		s.LastPressed = make([]bool, 4)
		s.Pressed = make([]bool, 4)
		for i := range s.Staged {
			for _, n := range c.Notes {
				var j int
				switch n.Type {
				case Head, Tail:
					j = 1
				case Extra:
					j = 2
				default:
					j = 0
				}
				if i == j {
					s.Staged[i] = n
					break
				}
			}
		}
		for _, n := range c.Notes {
			var i int
			switch n.Type {
			case Head, Tail: // Head/Tail of Roll/BigRoll
				i = 1
			case Extra: // Shake
				i = 2
			default: // Don, Kat, BigDon, BigKat
				i = 0
			}
			if s.Staged[i] == nil {
				s.Staged[i] = n
			}
		}
	}

	var b []byte
	b, err = os.ReadFile(cpath)
	if err != nil {
		return
	}
	s.MD5 = md5.Sum(b)
	s.Flow = 1
	s.FlowMarks = make([]float64, 0, c.Duration/1000)
	s.DelayedScore.Mode = ctrl.DelayedModeExp
	s.ScoreDrawer.Sprites = ScoreSprites
	title := fmt.Sprintf("gosu - %s - [%s]", c.MusicName, c.ChartName)
	ebiten.SetWindowTitle(title)
	return
}

// General
const (
	MinWaitBefore int64 = int64(-1.8 * 1000)
	WaitAfter     int64 = 3 * 1000
)

// Todo: put mode then switch Speed depends on the mode?
// func (s BaseScenePlay) Speed() float64 { return s.SpeedScale * s.BeatRatio() * s.BeatLengthScale }
func TimeToTick(time int64) int           { return int(float64(time) / 1000 * float64(TPS)) }
func TickToTime(tick int) int64           { return int64(float64(tick) / float64(TPS) * 1000) }
func (s BaseScenePlay) Time() int64       { return int64(float64(s.Tick) / float64(TPS) * 1000) }
func (s BaseScenePlay) BPMRatio() float64 { return s.TransPoint.BPM / s.MainBPM }
func (s BaseScenePlay) KeyAction(k int) input.KeyAction {
	return input.CurrentKeyAction(s.LastPressed[k], s.Pressed[k])
}

// func (s *BaseScenePlay) SetTick(rf *osr.Format) int64 { // Returns duration of waiting before starts
//
//		waitBefore := MinWaitBefore
//		if rf != nil && rf.BufferTime() < waitBefore {
//			waitBefore = rf.BufferTime()
//		}
//		s.Tick = TimeToTick(waitBefore)
//		return waitBefore
//	}

//	func (s *BaseScenePlay) SetBackground(path string) {
//		if img := draws.NewImage(path); img != nil {
//			sprite := draws.NewSpriteFromImage(img)
//			scale := screenSizeX / sprite.W()
//			sprite.SetScale(scale, scale, ebiten.FilterLinear)
//			sprite.SetPosition(screenSizeX/2, screenSizeY/2, draws.OriginCenter)
//			s.BackgroundDrawer.Sprite = sprite
//		} else {
//			s.BackgroundDrawer.Sprite = DefaultBackground
//		}
//	}

// Speed, BPM, Volume and Highlight
//
//	func (s *BaseScenePlay) SetInitTransPoint(first *TransPoint) {
//		s.TransPoint = first
//		for s.TransPoint.Time == s.TransPoint.Next.Time {
//			s.TransPoint = s.TransPoint.Next
//		}
//		// s.Volume = Volume * s.TransPoint.Volume
//	}
// func (s *BaseScenePlay) UpdateTransPoint() {
// 	for s.TransPoint.Next != nil && s.TransPoint.Next.Time <= s.Time() {
// 		s.TransPoint = s.TransPoint.Next
// 	}
// 	if s.LastVolume != s.TransPoint.Volume {
// 		s.LastVolume = s.Volume
// 		s.Volume = Volume * s.TransPoint.Volume
// 		// s.MusicPlayer.SetVolume(s.Volume)
// 	}
// }

// Audio
// func (s *BaseScenePlay) SetMusicPlayer(apath string) error { // apath stands for audio path.
//
//		if apath == "virtual" || apath == "" {
//			return nil
//		}
//		var err error
//		s.MusicPlayer, s.MusicCloser, err = audios.NewPlayer(apath)
//		if err != nil {
//			return err
//		}
//		s.MusicPlayer.SetVolume(s.Volume)
//		return nil
//	}
// func (s *BaseScenePlay) SetSoundMap(cpath string, names []string) error {
// 	s.Sounds = audios.NewSoundMap(&Volume)
// 	for _, name := range names {
// 		path := filepath.Join(filepath.Dir(cpath), name)
// 		s.Sounds.Register(path)
// 	}
// 	return nil
// }
