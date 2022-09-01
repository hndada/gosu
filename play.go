package gosu

import (
	"crypto/md5"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hndada/gosu/audios"
	"github.com/hndada/gosu/ctrl"
	"github.com/hndada/gosu/draws"
	"github.com/hndada/gosu/input"
)

// Fields can be set at outer function call.
// Update cannot be generalized; each scene use template fields in timely manner.
// Unit of time is a millisecond (1ms = 0.001s).
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

	// Changed speed should be applied after Update().
	SpeedHandler ctrl.F64Handler
	*TransPoint

	Result
	NoteWeights float64
	Combo       int
	Flow        float64

	BackgroundDrawer BackgroundDrawer
	DelayedScore     ctrl.Delayed
	ScoreDrawer      draws.NumberDrawer
	MeterDrawer      MeterDrawer
}

func IsAudioExisted(path string) bool {
	return path != "virtual" && path != ""
}
func MD5(cpath string) (v [16]byte, err error) {
	var b []byte
	b, err = os.ReadFile(cpath)
	if err != nil {
		return
	}
	v = md5.Sum(b)
	return
}

const (
	WaitBefore int64 = -1800
	WaitAfter  int64 = 3000
)

func TimeToTick(time int64) int     { return int(float64(time) / 1000 * float64(TPS)) }
func TickToTime(tick int) int64     { return int64(float64(tick) / float64(TPS) * 1000) }
func (s BaseScenePlay) Time() int64 { return TickToTime(s.Tick) }
func (s BaseScenePlay) KeyAction(k int) input.KeyAction {
	return input.CurrentKeyAction(s.LastPressed[k], s.Pressed[k])
}
func (s *BaseScenePlay) Ticker() { s.Tick++ }
func (s BaseScenePlay) IsDone() bool {
	return (ebiten.IsKeyPressed(ebiten.KeyEscape) ||
		s.Tick >= s.MaxTick)
}
func (s *BaseScenePlay) UpdateTransPoint() {
	for s.TransPoint.Next != nil && s.Time() >= s.TransPoint.Next.Time {
		s.TransPoint = s.TransPoint.Next
	}
}

// Todo: use Effecter in draws.BaseDrawer
type BackgroundDrawer struct {
	Sprite  draws.Sprite
	Dimness *float64
}

func NewBackgroundDrawer(path string, dimness *float64) (d BackgroundDrawer) {
	if sprite := NewBackground(path); sprite.IsValid() {
		d.Sprite = sprite
	} else {
		d.Sprite = DefaultBackground
	}
	d.Dimness = dimness
	return
}
func (d BackgroundDrawer) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.ColorM.ChangeHSV(0, 1, *d.Dimness)
	d.Sprite.Draw(screen, op)
}

func NewScoreDrawer() draws.NumberDrawer {
	return draws.NumberDrawer{
		Sprites:     ScoreSprites,
		SignSprites: SignSprites,
		DigitWidth:  ScoreSprites[0].W(),
		ZeroFill:    1,
		Origin:      ScoreSprites[0].Origin(),
	}
}
