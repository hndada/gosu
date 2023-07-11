package mode

import (
	"fmt"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

// TPS affects only on Update(), not on Draw().
var TPS = float64(ebiten.TPS())

func ToTick(ms int32) int   { return int(TPS * float64(ms) / 1000) }
func ToTime(tick int) int32 { return int32(float64(tick) / TPS * 1000) }

// func ToTick(t time.Duration) int { return int(TPS * t.Seconds()) }
// func ToTime(t int) time.Duration { return time.Duration(float64(t)/TPS) * time.Second }

type Timer struct {
	StartTime time.Time
	Offset    time.Duration
	offset    *time.Duration
	Tick      int
	MaxTick   int
}

func NewTimer(ms int32, offset *time.Duration) Timer {
	const wait = 1800
	return Timer{
		StartTime: time.Now().Add(wait * time.Millisecond),
		Offset:    *offset,
		offset:    offset,
		Tick:      ToTick(-wait),
		MaxTick:   ToTick(ms + wait),
	}
}
func (t Timer) Now() time.Duration { return ToTime(t.Tick) }
func (t *Timer) Ticker() {
	t.Tick++

	// Adjusting offset in real-time.
	if td := *t.offset - t.Offset; td != 0 {
		t.Offset += td
		t.Tick += ToTick(td)
	}

	// Try sync after buffer time ends.
	if t.Now() > 0 {
		t.sync()
	}
}
func (t *Timer) sync() {
	const threshold = 30 * time.Millisecond
	since := time.Since(t.StartTime)
	if e := since - t.Now(); e >= threshold {
		fmt.Printf("%dms: adjusting time error (%dms)\n", since, e/time.Millisecond)
		t.Tick += ToTick(e)
	}
}
func (t Timer) IsDone() bool { return t.Tick >= t.MaxTick }
