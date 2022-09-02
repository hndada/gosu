package gosu

import (
	"crypto/md5"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hndada/gosu/audios"
	"github.com/hndada/gosu/ctrl"
	"github.com/hndada/gosu/input"
)

// Fields can be set at outer function call.
// Update cannot be generalized; each scene use template fields in timely manner.

type BaseScenePlay struct {
	Tick    int
	MaxTick int // Tick corresponding to EndTime = Duration + WaitAfter

	VolumeHandler ctrl.F64Handler
	MusicPlayer   *audio.Player
	MusicCloser   func() error
	Sounds        audios.SoundMap // A player for sample sound is generated at a place.

	FetchPressed func() []bool
	LastPressed  []bool
	Pressed      []bool

	SpeedHandler ctrl.F64Handler
	*TransPoint

	Result
	// NoteWeights is a sum of weight of marked notes.
	// This is also max value of each score sum can get at the time.
	NoteWeights float64
	Combo       int
	Flow        float64

	BackgroundDrawer BackgroundDrawer
	ScoreDrawer      ScoreDrawer
	MeterDrawer      MeterDrawer
}

const (
	WaitBefore int64 = -1800
	WaitAfter  int64 = 3000
)

func MD5(cpath string) (v [16]byte, err error) {
	var b []byte
	b, err = os.ReadFile(cpath)
	if err != nil {
		return
	}
	v = md5.Sum(b)
	return
}
func TimeToTick(time int64) int     { return int(float64(time) / 1000 * float64(TPS)) }
func TickToTime(tick int) int64     { return int64(float64(tick) / float64(TPS) * 1000) }
func (s BaseScenePlay) Time() int64 { return TickToTime(s.Tick) }
func (s *BaseScenePlay) Ticker()    { s.Tick++ }
func (s BaseScenePlay) IsDone() bool {
	return (ebiten.IsKeyPressed(ebiten.KeyEscape) ||
		s.Tick >= s.MaxTick)
}
func (s *BaseScenePlay) UpdateTransPoint() {
	for s.TransPoint.Next != nil && s.Time() >= s.TransPoint.Next.Time {
		s.TransPoint = s.TransPoint.Next
	}
}
func (s BaseScenePlay) KeyAction(k int) input.KeyAction {
	return input.CurrentKeyAction(s.LastPressed[k], s.Pressed[k])
}
