package mode

import (
	"fmt"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hndada/gosu/times"
)

const Wait = 1800

type Timer struct {
	StartTime time.Time
	Now       int64
	Offset    int64
	offset    *int64
	Tick      int
	MaxTick   int
	TPS       int
}

var (
	ToTick = times.ToTick
	ToTime = times.ToTime
)

func NewTimer(duration int64, offset *int64, tps int) Timer {
	return Timer{
		StartTime: time.Now().Add(Wait * time.Millisecond),
		Now:       -Wait,
		Offset:    *offset,
		offset:    offset,
		Tick:      ToTick(-Wait, tps),
		MaxTick:   ToTick(duration+Wait, tps),
		TPS:       tps,
	}
}
func (t *Timer) Ticker() {
	t.Tick++
	// Adjusting offset in real-time.
	if td := t.Offset - *t.offset; td != 0 {
		t.Offset += td
		t.Tick += ToTick(td, t.TPS)
	}
	if t.Now > 0 && ebiten.ActualTPS() < 0.8*float64(t.TPS) {
		t.sync()
	}
	t.Now = ToTime(t.Tick, t.TPS)
}
func (t *Timer) sync() {
	since := time.Since(t.StartTime).Milliseconds() // - Wait
	if e := since - t.Now; e >= 1 {
		fmt.Printf("adjusting time error at %dms: %d\n", t.Now, e)
		t.Tick += ToTick(e, t.TPS)
	}
}
func (t Timer) IsDone() bool { return t.Tick >= t.MaxTick }
func (t *Timer) SetDone()    { t.Tick = t.MaxTick }
