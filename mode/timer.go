package mode

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

// TPS affects only on Update(), not on Draw().
var TPS = float64(ebiten.TPS())

func ToTick(ms int32) int       { return int(TPS * float64(ms) / 1000) }
func ToTime(tick int) int32     { return int32(float64(tick) / TPS * 1000) }
func ToSecond(ms int32) float64 { return float64(ms) / 1000 }

type Timer struct {
	startTime time.Time
	pauseTime time.Time
	paused    bool
	offset    int32
}

func NewTimer(wait time.Duration) Timer {
	return Timer{startTime: time.Now().Add(wait)}
}

func (t Timer) Now() int32 {
	var duration time.Duration
	if t.paused {
		duration = t.pauseTime.Sub(t.startTime)
	} else {
		duration = time.Since(t.startTime)
	}
	return int32(duration.Milliseconds())
}

func (t Timer) IsPaused() bool { return t.paused }

// Music is hard to seek precisely.
// Hence, we simply add offset to StartTime.
// Positive offset makes notes delayed.
// It is no use to set offset before music starts.
func (t *Timer) SetOffset(new int32) {
	old := t.offset
	diff := time.Duration(new-old) * time.Millisecond
	t.startTime = t.startTime.Add(diff)
	t.offset = new
}

func (t *Timer) Pause() {
	t.pauseTime = time.Now()
	t.paused = true
}

func (t *Timer) Resume() {
	elapsedTime := time.Since(t.pauseTime)
	t.startTime = t.startTime.Add(elapsedTime)
	t.paused = false
}

// func (t *Timer) sync() {
// 	const threshold = 30 * 1000
// 	since := int32(time.Since(t.startTime).Milliseconds())
// 	if e := since - t.Now(); e >= threshold {
// 		fmt.Printf("%dms: adjusting time error (%dms)\n", since, e)
// 		t.Tick += ToTick(e)
// 	}
// }
