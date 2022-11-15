package audios

import (
	"fmt"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hndada/gosu/input"
)

var tps = ebiten.TPS()

func SetTPS(v int) { tps = v }

// Todo: remove it
func init() {
	SetTPS(1000)
}

// Time is a point of time.
// Duration a length of time.
func ToTick(time int64) int { return int(float64(time) / 1000 * float64(tps)) }
func ToTime(tick int) int64 { return int64(float64(tick) / float64(tps) * 1000) }

const Wait = 1800

type Timer struct {
	StartTime time.Time
	Tick      int
	MaxTick   int
	Offset    *int64
	offset    int64
	Now       int64
	Pause     bool
}

func NewTimer(duration int64, offset *int64) Timer {
	return Timer{
		StartTime: time.Now().Add(Wait * time.Millisecond),
		Offset:    offset,
		offset:    *offset,
		Tick:      ToTick(-Wait),
		MaxTick:   ToTick(duration + Wait),
		Now:       -Wait,
	}
}

func (t Timer) IsDone() bool { return t.Tick >= t.MaxTick }
func (t *Timer) SetDone()    { t.Tick = t.MaxTick }
func (t *Timer) Ticker() {
	if inpututil.IsKeyJustPressed(input.KeyTab) {
		t.Pause = !t.Pause
	}
	if t.Pause {
		return
	}
	t.Tick++
	// Adjusting offset in real-time.
	if td := *t.Offset - t.offset; td != 0 {
		t.offset += td
		t.Tick += ToTick(td)
	}
	if t.Now > 0 && ebiten.ActualTPS() < 0.8*float64(tps) {
		t.sync()
	}
	t.Now = ToTime(t.Tick)
}
func (t *Timer) sync() {
	since := time.Since(t.StartTime).Milliseconds() // - Wait
	if e := since - t.Now; e >= 1 {
		fmt.Printf("adjusting time error at %dms: %d\n", t.Now, e)
		t.Tick += ToTick(e)
	}
}
