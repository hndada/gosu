package mode

import (
	"time"

	"github.com/hndada/gosu/audios"
	"github.com/hndada/gosu/draws"
	"github.com/hndada/gosu/input"
)

// Interface is also used when it uses the unknown struct.
type ScenePlay interface {
	Now() int32
	StartMusic()
	PlaySound(sample Sample, scale float64)
	SetOffset(currentOffset int32)
	SetMusicVolume(vol float64)
	IsPaused() bool
	Pause()
	Resume()
	// UpdateDynamic() // This is used only in each mode.

	// The following methods would go implemented in each mode.
	Speed() float64
	Update() any
	SetSpeedScale()

	WindowTitle() string
	BackgroundFilename() string
	NoteExposureDuration() int32
	Finish() any

	Draw(screen draws.Image)
	DebugPrint(screen draws.Image)
}

// BaseScenePlay implements ScenePlay.
type BaseScenePlay struct {
	StartTime  time.Time
	lastOffset int32
	PauseTime  time.Time
	Dynamic    *Dynamic

	MusicPlayer audios.MusicPlayer
	musicPlayed bool
	SoundMap    audios.SoundMap
	Keyboard    input.Keyboard
	paused      bool
}

func (s BaseScenePlay) Now() int32 {
	return int32(time.Since(s.StartTime).Milliseconds())
}

func (s BaseScenePlay) StartMusic() {
	if !s.musicPlayed && s.Now() >= 0 {
		s.MusicPlayer.Play()
		s.musicPlayed = true
	}
}

func (s BaseScenePlay) PlaySound(sample Sample, scale float64) {
	name := sample.Name
	vol := sample.Volume
	if vol == 0 {
		vol = s.Dynamic.Volume
	}
	vol *= scale
	s.SoundMap.Play(name, vol)
}

// Music is hard to seek precisely.
// Hence, we simply add offset to StartTime.
// Positive offset makes notes delayed.
// It is no use to set offset before music starts.
func (s *BaseScenePlay) SetOffset(currentOffset int32) {
	diff := time.Duration(currentOffset-s.lastOffset) * time.Millisecond
	s.StartTime = s.StartTime.Add(diff)
	s.lastOffset = currentOffset
}

func (s *BaseScenePlay) SetMusicVolume(vol float64) {
	s.MusicPlayer.SetVolume(vol)
}

func (s BaseScenePlay) IsPaused() bool { return s.paused }

func (s *BaseScenePlay) Pause() {
	s.PauseTime = time.Now()
	s.MusicPlayer.Pause()
	s.Keyboard.Pause()
	s.paused = true
}

func (s *BaseScenePlay) Resume() {
	elapsedTime := time.Since(s.PauseTime)
	s.StartTime = s.StartTime.Add(elapsedTime)
	s.MusicPlayer.Play()
	s.Keyboard.Resume()
	s.paused = false
}

func (s *BaseScenePlay) UpdateDynamic() {
	d := s.Dynamic
	for d.Next != nil && s.Now() >= d.Next.Time {
		d = d.Next
	}
	s.Dynamic = d
}

func (s BaseScenePlay) Finish() {
	s.MusicPlayer.Close()
	s.Keyboard.Close()
}
