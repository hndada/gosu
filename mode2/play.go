package mode

import (
	"fmt"
	"io/fs"
	"time"

	"github.com/hndada/gosu/audios"
	"github.com/hndada/gosu/format/osr"
	"github.com/hndada/gosu/input"
)

type ScenePlay struct {
	ChartHeader
	cursor float64

	Timer       Timer
	Keyboard    input.Keyboard
	MusicPlayer audios.MusicPlayer
	SoundPlayer audios.SoundPlayer
	Dynamics    Dynamics

	Update           func([]input.KeyboardAction)
	SpeedScale       float64
	UpdateSpeedScale func(new, old float64)
	Combo            ComboComp
	Score            ScoreComp
}

func NewScenePlay(fsys fs.FS, name string, replay *osr.Format) (s ScenePlay, err error) {
	format, hash, err := LoadChartFile(fsys, name)
	if err != nil {
		return
	}
	s.Dynamics, err = NewDynamics(format)
	if len(s.Dynamics.Data) == 0 {
		err = fmt.Errorf("no Dynamics in the chart")
		return
	}
	s.ChartHeader = NewChartHeader(format, hash)
	s.SetSpeedScale(1)
}

func (s *ScenePlay) BeforeUpdate() {
	s.tryPlayMusic()
	s.readInput()
}

func (s *ScenePlay) AfterUpdate() {
	s.Update(s.readInput())
}

func (s *ScenePlay) tryPlayMusic() {
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
func (s *ScenePlay) readInput() []input.KeyboardAction {
	if s.Keyboard != nil {
		return s.Keyboard.Read(s.Timer.Now())
	}
	return s.Keyboard.Reader.Read(s.Timer.Now()) // for replay
}

func (s ScenePlay) Speed() float64 {
	return s.Dynamics.Current().Speed * s.SpeedScale
}

func (s *ScenePlay) SetSpeedScale(new float64) {
	old := s.SpeedScale
	s.UpdateSpeedScale(new, old)
	s.SpeedScale = new
}

func (s *ScenePlay) Pause() {
	s.Timer.Pause()
	s.MusicPlayer.Pause()
	s.Keyboard.Pause()
}

func (s *ScenePlay) Resume() {
	s.Timer.Resume()
	s.MusicPlayer.Resume()
	s.Keyboard.Resume()
}

func (s *ScenePlay) Close() {
	// Music keeps playing at result scene.
	// s.MusicPlayer.Close()
	s.Keyboard.Close()
}

func (s ScenePlay) Finish() any {
	return s.Scorer
}
