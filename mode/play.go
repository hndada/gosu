package mode

import (
	"io"
	"time"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hndada/gosu/render"
)

// In milliseconds.
// func Time(tick int) int64       { return int64(float64(tick) / float64(MaxTPS) * 1000) }
func TimeToTick(msec int64) int { return int(float64(msec) * float64(MaxTPS) / 1000) }
func TickToTime(tick int) int64 { return int64(1000 * float64(tick) / float64(MaxTPS)) }

// This is template struct. Fields can be set at outer function call.
type ScenePlay struct {
	Play bool // Whether the scene is for play or not
	Tick int

	// MusicFile   io.ReadSeekCloser
	// MusicStreamer  io.ReadSeeker
	MusicPlayer    *audio.Player
	MusicCloser    io.Closer
	SoundStreamers map[string]io.ReadSeeker
	SoundClosers   []io.Closer
	// A player for Sound is generated at a place.

	MainBPM   float64
	BaseSpeed float64 // Todo: BaseSpeed -> SpeedBase
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

// func NewScenePlay(c *Chart, cpath string, play bool) ScenePlay {
// 	var s ScenePlay
// 	s.MainBPM, _, _ = BPMs(c.TransPoints, c.Duration) // Todo: Need a test
// 	s.TransPoint = c.TransPoints[0]
// 	for s.TransPoint.Time == s.TransPoint.Next.Time {
// 		s.TransPoint = s.TransPoint.Next
// 	}
// 	s.Flow = 1

// 	s.Play = play
// 	if !s.Play {
// 		return s
// 	}
// 	s.MusicFile, s.MusicPlayer = audio.NewPlayer(c.MusicPath(cpath))
// 	s.MusicPlayer.SetVolume(Volume)
// 	return s
// }

// Update cannot be generalized; each scene use template fields in timely manner.

// Replay is a entire key stroke timed-log.
type PlayToResultArgs struct {
	Time time.Time // Playing finish time
	MD5  [16]byte
	// Replay // Todo: implement
	Result
}
