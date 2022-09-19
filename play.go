package gosu

import (
	"fmt"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hndada/gosu/audios"
	"github.com/hndada/gosu/ctrl"
	"github.com/hndada/gosu/input"
)

func SetTitle(c ChartHeader) {
	title := fmt.Sprintf("gosu - %s - [%s]", c.MusicName, c.ChartName)
	ebiten.SetWindowTitle(title)
}

// Time is a point of time, duration a length of time.
func TimeToTick(time int64) int { return int(float64(time) / 1000 * float64(TPS)) }
func TickToTime(tick int) int64 { return int64(float64(tick) / float64(TPS) * 1000) }

const (
	WaitBefore int64 = -1800
	WaitAfter  int64 = 3000
)

type Timer struct {
	Tick    int
	MaxTick int // A tick corresponding to EndTime = Duration + WaitAfter
	Time    int64
}

// func (t Timer) IsDone() bool { return ebiten.IsKeyPressed(ebiten.KeyEscape) }
func (t Timer) IsDone() bool {
	return ebiten.IsKeyPressed(ebiten.KeyEscape) || t.Tick >= t.MaxTick
}
func (t *Timer) SetTicks(duration int64) {
	t.Tick = TimeToTick(WaitBefore)
	t.MaxTick = TimeToTick(duration + WaitAfter)
	t.Time = TickToTime(t.Tick)
}
func (t *Timer) Ticker() {
	t.Tick++
	t.Time = TickToTime(t.Tick)
}

type MusicPlayer struct {
	VolumeHandler ctrl.F64Handler
	Player        *audio.Player
	Closer        func() error
}

func NewMusicPlayer(mvh ctrl.F64Handler, path string) (MusicPlayer, error) {
	player, closer, err := audios.NewPlayer(path)
	if err != nil {
		return MusicPlayer{}, err
	}
	player.SetVolume(*mvh.Target)
	player.SetBufferSize(100 * time.Millisecond)
	return MusicPlayer{
		VolumeHandler: mvh,
		Player:        player,
		Closer:        closer,
	}, nil
}
func (mp MusicPlayer) Play() {
	if mp.Player == nil {
		return
	}
	mp.Player.Play()
}

// Todo: volume does not increment gradually.
func (mp MusicPlayer) Update() {
	if fired := mp.VolumeHandler.Update(); fired {
		vol := *mp.VolumeHandler.Target
		mp.Player.SetVolume(vol)
	}
}
func (mp MusicPlayer) Close() {
	if mp.Player == nil {
		return
	}
	mp.Player.Close()
	mp.Closer()
}

// Todo: need to refactor
type EffectPlayer struct {
	VolumeHandler ctrl.F64Handler
	Effects       audios.SoundMap // A player for sample sound is generated at a place.
}

func NewEffectPlayer(evh ctrl.F64Handler) EffectPlayer {
	return EffectPlayer{
		VolumeHandler: evh,
		Effects:       audios.NewSoundMap(evh.Target),
	}
}

type KeyLogger struct {
	FetchPressed func() []bool
	LastPressed  []bool
	Pressed      []bool
}

func NewKeyLogger(keySettings []input.Key) (k KeyLogger) {
	keyCount := len(keySettings)
	k.FetchPressed = input.NewListener(keySettings)
	k.LastPressed = make([]bool, keyCount)
	k.Pressed = make([]bool, keyCount)
	return
}
func (l KeyLogger) KeyAction(k int) input.KeyAction {
	return input.CurrentKeyAction(l.LastPressed[k], l.Pressed[k])
}
