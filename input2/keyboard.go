package input

import (
	"time"

	"github.com/hndada/gosu/times"
)

// input package should not be tightly adjusted to gosu.
type KeyboardState struct {
	Time        time.Duration
	PressedList []bool
}

type KeyboardBuffer struct {
	states []KeyboardState
	idx    int
}

type KeyboardReader interface {
	Read(t time.Duration) []KeyboardState
}

type KeyboardListener interface {
	Listen()
}

// When to use lock: read or write to shared resources.
// A practical example would be public toilet.

// No additional adjustment for keyboard when offset has changed.
// Both music and keyboard cannot seek at precise position once they start.
type Keyboard struct {
	KeyboardBuffer
	startTime          time.Time
	fetchKeyboardState func() []bool
	period             time.Duration
	stop               chan struct{}
}

func NewKeyboard(keys []Key, pollingRate float64) *Keyboard {
	second := float64(time.Second) * times.PlaybackRate()
	return &Keyboard{
		fetchKeyboardState: newFetchKeyboardState(keys),
		period:             time.Duration(second / pollingRate),
		stop:               make(chan struct{}),
	}
}

func (kb *Keyboard) Listen() {
	kb.startTime = times.Now()
	go func() {
		for {
			select {
			case <-kb.stop:
				return
			default:
				start := times.Now()
				kb.poll()
				elapsed := times.Since(start)
				// It is fine to pass negative value to time.Sleep.

				// It is fine not to update period by changing playback rate;
				// It would just cause more or less of polling.
				time.Sleep(kb.period - elapsed)
			}
		}
	}()
}

func (kb *Keyboard) poll() {
	t := times.Since(kb.startTime)
	ps := kb.fetchKeyboardState()
	state := KeyboardState{t, ps}

	// kb.mu.Lock()
	kb.states = append(kb.states, state)
	// kb.mu.Unlock()
}

func (kb *Keyboard) Stop() {
	kb.stop <- struct{}{}
}

// Output trims redundant states then returns the states.
func (kb Keyboard) Output() []KeyboardState {
	trimmed := []KeyboardState{}
	for i, s := range kb.states {
		if i == 0 {
			trimmed = append(trimmed, s)
			continue
		}
		last := kb.states[i-1]
		if isEqual(last, s) {
			continue
		}
		trimmed = append(trimmed, s)
	}
	return trimmed
}

// No worry of accessing with nil pointer.
// https://go.dev/play/p/B4Z1LwQC_jP
