package piano

import (
	"time"

	"github.com/hndada/gosu/format/osr"
	"github.com/hndada/gosu/input"
)

type ReplayListener struct {
	States []input.KeyboardState
	index  int

	StartTime time.Time
	PauseTime time.Time
	paused    bool
}

// ReplayListener supposes the time of first state is 0 ms with no any inputs.
func NewReplayListener(f *osr.Format, keyCount int, bufferTime time.Duration) *ReplayListener {
	actions := f.ReplayData
	// actions := append(f.ReplayData, osr.Action{W: 2e9})

	// clean replay data
	// Osu replay data uses X for storing key count.
	for i := 0; i < 2; i++ {
		if i < len(actions) {
			break
		}
		if a := actions[i]; a.Y == -500 {
			a.X = 0
		}
	}

	getKeyStates := func(a osr.Action) []bool {
		ps := make([]bool, keyCount)
		var k int
		for x := int(a.X); x > 0; x /= 2 {
			if x%2 == 1 {
				ps[k] = true
			}
			k++
		}
		return ps
	}

	states := make([]input.KeyboardState, 0, len(actions)+1)
	var t int32
	for _, a := range actions {
		t += int32(a.W)
		s := input.KeyboardState{Time: t, Pressed: getKeyStates(a)}
		states = append(states, s)
	}

	return &ReplayListener{
		States: states,
		index:  1,

		StartTime: time.Now().Add(bufferTime),
	}
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
	rl.PauseTime = time.Now()
	rl.paused = true
}

func (rl *ReplayListener) Resume() {
	pauseDuration := time.Since(rl.PauseTime)
	rl.StartTime = rl.StartTime.Add(pauseDuration)
	rl.paused = false
}

func (rl *ReplayListener) IsPaused() bool { return rl.paused }

func (rl *ReplayListener) Output() []input.KeyboardState { return rl.States }

// Close does nothing on ReplayListener.
func (rl *ReplayListener) Close() {}
