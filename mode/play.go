package mode

import (
	"time"

	"github.com/hndada/gosu/audios"
	"github.com/hndada/gosu/input"
)

// type ScenePlayArgs struct {
// 	FS            fs.FS
// 	ChartFilename string
// 	Mods          any
// 	Replay        *osr.Format
// }

// interface is also used when it uses the unknown struct.
type BaseScenePlay struct {
	StartTime  time.Time
	lastOffset int32
	PauseTime  time.Time
	Dynamic    *Dynamic

	MusicPlayer     audios.MusicPlayer
	hasMusicStarted bool
	SoundPlayer     audios.SoundPlayer
	Keyboard        input.Keyboard
	paused          bool
}

func NewBaseScenePlay() BaseScenePlay {
	return BaseScenePlay{}
}

func (s *BaseScenePlay) Pause() {
	s.PauseTime = time.Now()
	s.MusicPlayer.Pause()
	s.Keyboard.Pause()
	s.paused = true
}

func (s *BaseScenePlay) Resume() {
	elapsedTime := time.Now().Sub(s.PauseTime)
	s.StartTime = s.StartTime.Add(elapsedTime)
	s.MusicPlayer.Play()
	s.Keyboard.Resume()
	s.paused = false
}

func (s BaseScenePlay) IsPaused() bool { return s.paused }

func (s BaseScenePlay) Now() int32 {
	return int32(time.Since(s.StartTime).Milliseconds())
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

func (s *BaseScenePlay) UpdateDynamic() {
	dy := s.Dynamic
	for dy.Next != nil && s.Now() >= dy.Next.Time {
		dy = dy.Next
	}
	s.Dynamic = dy
}

func (s BaseScenePlay) StartMusic() {
	if !s.hasMusicStarted && s.Now() >= 0 {
		s.MusicPlayer.Play()
		s.hasMusicStarted = true
	}
}

func (s BaseScenePlay) Finish() {
	s.MusicPlayer.Close()
	s.Keyboard.Close()
}
