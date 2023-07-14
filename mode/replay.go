package mode

import (
	"time"

	"github.com/hndada/gosu/input"
)

// Todo: handle error message from osr.NewFormat
type ReplayListener struct {
	States []input.KeyboardState
	index  int

	StartTime time.Time
	pauseTime time.Time
	paused    bool
}

func (rl *ReplayListener) Now() int32 {
	return int32(time.Since(rl.StartTime).Milliseconds())
}

func (rl *ReplayListener) Fetch() (kas []input.KeyboardAction) {
	end := rl.index
	now := rl.Now()
	for ; end < len(rl.States)-1; end++ {
		if rl.States[end+1].Time > now {
			break
		}
	}

	last := rl.States[rl.index-1]
	// From last fetched state to latest state
	for _, current := range rl.States[rl.index : end+1] {
		as := input.KeyActions(last.Pressed, current.Pressed)
		for t := last.Time; t < current.Time; t++ {
			ka := input.KeyboardAction{Time: t, Action: as}
			kas = append(kas, ka)
		}
		last = current
	}

	// From latest state to now
	las := input.KeyActions(last.Pressed, last.Pressed)
	for t := last.Time; t < rl.Now(); t++ {
		ka := input.KeyboardAction{Time: t, Action: las}
		kas = append(kas, ka)
	}

	// Update index
	rl.index = end
	return nil
}

// Poll does nothing on ReplayListener.
func (rl *ReplayListener) Poll() {}

func (rl *ReplayListener) Pause() {
	rl.pauseTime = time.Now()
	rl.paused = true
}

func (rl *ReplayListener) Resume() {
	pauseDuration := time.Since(rl.pauseTime)
	rl.StartTime = rl.StartTime.Add(pauseDuration)
	rl.paused = false
}

func (rl *ReplayListener) IsPaused() bool { return rl.paused }

func (rl *ReplayListener) Output() []input.KeyboardState { return rl.States }

// Close does nothing on ReplayListener.
func (rl *ReplayListener) Close() {}
