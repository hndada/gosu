package gosu

import (
	"crypto/md5"
	"fmt"
	"os"
	"path/filepath"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hndada/gosu/audios"
	"github.com/hndada/gosu/draws"
	"github.com/hndada/gosu/format/osr"
	"github.com/hndada/gosu/input"
)

// Fields can be set at outer function call.
// Update cannot be generalized; each scene use template fields in timely manner.
// Unit of time is a millisecond (1ms = 0.001s).
type BaseScenePlay struct {
	// General
	Tick             int
	EndTime          int64 // EndTime = Duration + WaitAfter
	Mods             Mods
	BackgroundDrawer BackgroundDrawer

	// Speed, BPM, Volume and Highlight
	MainBPM   float64
	SpeedBase float64
	*TransPoint

	// Audio
	LastVolume  float64
	Volume      float64
	MusicPlayer *audio.Player
	MusicCloser func() error
	Sounds      audios.SoundMap // A player for sample sound is generated at a place.

	// Input
	FetchPressed func() []bool
	LastPressed  []bool
	Pressed      []bool

	// Note
	BarLineDrawer BarLineDrawer

	// Score
	Result
	NoteWeights float64
	Combo       int
	Flow        float64
	ScoreDrawer ScoreDrawer
	TimingMeter Meter
}

// General
const (
	DefaultWaitBefore int64 = int64(-1.8 * 1000)
	DefaultWaitAfter  int64 = 3 * 1000
)

func (s BaseScenePlay) BeatRatio() float64 { return s.TransPoint.BPM / s.MainBPM }
func (s BaseScenePlay) Speed() float64     { return s.SpeedBase * s.BeatRatio() * s.BeatScale }
func MD5(path string) [md5.Size]byte {
	b, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return md5.Sum(b)
}

func TimeToTick(time int64) int     { return int(float64(time) / 1000 * float64(MaxTPS)) }
func TickToTime(tick int) int64     { return int64(float64(tick) / float64(MaxTPS) * 1000) }
func (s BaseScenePlay) Time() int64 { return int64(float64(s.Tick) / float64(MaxTPS) * 1000) }
func (s *BaseScenePlay) SetTick(rf *osr.Format) int64 { // Returns duration of waiting before starts
	waitBefore := DefaultWaitBefore
	if rf != nil && rf.BufferTime() < waitBefore {
		waitBefore = rf.BufferTime()
	}
	s.Tick = TimeToTick(waitBefore) - 1 // s.Update() starts with s.Tick++
	return waitBefore
}
func (s BaseScenePlay) SetWindowTitle(c BaseChart) {
	title := fmt.Sprintf("gosu - %s - [%s]", c.MusicName, c.ChartName)
	ebiten.SetWindowTitle(title)
}
func (s *BaseScenePlay) SetBackground(path string) {
	if img := draws.NewImage(path); img != nil {
		sprite := draws.NewSpriteFromImage(img)
		scale := screenSizeX / sprite.W()
		sprite.SetScale(scale, scale, ebiten.FilterLinear)
		sprite.SetPosition(screenSizeX/2, screenSizeY/2, draws.OriginModeCenter)
		s.BackgroundDrawer.Sprite = sprite
	} else {
		s.BackgroundDrawer.Sprite = DefaultBackground
	}
}

// Speed, BPM, Volume and Highlight
func (s *BaseScenePlay) SetInitTransPoint(first *TransPoint) {
	s.TransPoint = first
	for s.TransPoint.Time == s.TransPoint.Next.Time {
		s.TransPoint = s.TransPoint.Next
	}
	s.Volume = Volume * s.TransPoint.Volume
}
func (s *BaseScenePlay) UpdateTransPoint() {
	for s.TransPoint.Next != nil && s.TransPoint.Next.Time <= s.Time() {
		s.TransPoint = s.TransPoint.Next
	}
	if s.LastVolume != s.TransPoint.Volume {
		s.LastVolume = s.Volume
		s.Volume = Volume * s.TransPoint.Volume
		// s.MusicPlayer.SetVolume(s.Volume)
	}
}

// Audio
func (s *BaseScenePlay) SetMusicPlayer(apath string) error { // apath stands for audio path.
	if apath == "virtual" || apath == "" {
		return nil
	}
	var err error
	s.MusicPlayer, s.MusicCloser, err = audios.NewPlayer(apath)
	if err != nil {
		return err
	}
	s.MusicPlayer.SetVolume(s.Volume)
	return nil
}
func (s *BaseScenePlay) SetSoundMap(cpath string, names []string) error {
	s.Sounds = audios.NewSoundMap(&Volume)
	for _, name := range names {
		path := filepath.Join(filepath.Dir(cpath), name)
		s.Sounds.Register(path)
	}
	return nil
}

// Input
func (s BaseScenePlay) KeyAction(k int) input.KeyAction {
	return input.CurrentKeyAction(s.LastPressed[k], s.Pressed[k])
}
