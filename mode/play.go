package mode

import (
	"crypto/md5"
	"time"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hndada/gosu/audioutil"
	"github.com/hndada/gosu/render"
)

// In milliseconds.
// func Time(tick int) int64       { return int64(float64(tick) / float64(MaxTPS) * 1000) }
func TimeToTick(time int64) int { return int(float64(time) / 1000 * float64(MaxTPS)) }
func TickToTime(tick int) int64 { return int64(float64(tick) / float64(MaxTPS) * 1000) }

// This is template struct. Fields can be set at outer function call.
// Update cannot be generalized; each scene use template fields in timely manner.
type ScenePlay struct {
	Play bool // Whether the scene is for play or not
	Tick int
	MD5  [md5.Size]byte // MD5 for raw chart file.

	MusicPlayer *audio.Player
	SoundBytes  map[string][]byte // A player for Sound is generated at a place.

	MainBPM   float64
	SpeedBase float64
	*TransPoint

	FetchPressed func() []bool
	LastPressed  []bool
	Pressed      []bool

	Result
	Combo int
	Flow  float64

	// Graphics
	DelayedScore float64
	Background   render.Sprite
	TimingMarks  []TimingMark
}

const (
	DefaultWaitBefore int64 = int64(-1.8 * 1000)
	DefaultWaitAfter  int64 = 3 * 1000
)

func (s ScenePlay) Time() int64 { return int64(float64(s.Tick) / float64(MaxTPS) * 1000) }

func (s *ScenePlay) SetMusicPlayer(musicPath string) error {
	if musicPath == "virtual" || musicPath == "" {
		return nil
	}
	mbytes, err := audioutil.NewBytes(musicPath)
	if err != nil {
		return err
	}
	s.MusicPlayer = audioutil.Context.NewPlayerFromBytes(mbytes)
	s.MusicPlayer.SetVolume(Volume)
	return nil
}

// Replay is a entire key stroke timed-log.
type PlayToResultArgs struct {
	Time time.Time // Playing finish time
	MD5  [16]byte
	// Replay // Todo: implement
	Result
}
