package mode

import (
	"fmt"
	"io/fs"
	"time"

	"github.com/hndada/gosu/audios"
	"github.com/hndada/gosu/format/osr"
	"github.com/hndada/gosu/input"
)

type Scene struct {
	ChartHeader ChartHeader
	Dynamics    Dynamics
	SpeedScale  float64
	Timer       Timer
	Keyboard    input.Keyboard
	MusicPlayer audios.MusicPlayer
	SoundPlayer audios.SoundPlayer
}

func NewScene(fsys fs.FS, name string, replay *osr.Format) (s Scene, err error) {
	format, hash, err := LoadChartFile(fsys, name)
	if err != nil {
		return
	}
	s.ChartHeader = NewChartHeader(format, hash)
	s.Dynamics, err = NewDynamics(format)
	if len(s.Dynamics.Dynamics) == 0 {
		err = fmt.Errorf("no Dynamics in the chart")
		return
	}
	return
}

func (s Scene) Speed() float64 {
	return s.Dynamics.Current().Speed * s.SpeedScale
}

func (s *Scene) BeforeUpdate() {
	s.tryPlayMusic()
	s.readInput()
}

func (s *Scene) tryPlayMusic() {
	if s.MusicPlayer.IsPlayed() {
		return
	}
	now := s.Timer.Now()
	if now >= *s.MusicOffset && now < 300 {
		s.MusicPlayer.Play()
		s.Timer.SetMusicPlayed(time.Now())
	}
}

// readInput guarantees that length of return value is at least one.
// The receiver should be pointer for updating replay's index.
func (s *Scene) readInput() []input.KeyboardAction {
	if s.Keyboard != nil {
		return s.Keyboard.Read(s.Timer.Now())
	}
	return s.Keyboard.Reader.Read(s.Timer.Now()) // for replay
}

func (s *Scene) Pause() {
	s.Timer.Pause()
	s.MusicPlayer.Pause()
	s.Keyboard.Pause()
}

func (s *Scene) Resume() {
	s.Timer.Resume()
	s.MusicPlayer.Resume()
	s.Keyboard.Resume()
}

func (s *Scene) Close() {
	// Music keeps playing at result scene.
	// s.MusicPlayer.Close()
	s.Keyboard.Close()
}
