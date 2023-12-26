package input

import (
	"sync"
	"time"

	"github.com/hndada/gosu/times"
)

// input package should not be tightly adjusted to gosu.
type KeyboardState struct {
	Time        time.Duration
	PressedList []bool
}

// When to use lock: read or write to shared resources.
// A practical example would be public toilet.

// No additional adjustment for keyboard when offset has changed.
// Both music and keyboard cannot seek at precise position once they start.
type Keyboard struct {
	fetchKeyboardState func() []bool
	states             []KeyboardState
	idx                int
	period             time.Duration
	startTime          time.Time
	pauseTime          time.Time
	paused             bool
	mu                 *sync.Mutex // for lock
	pauseChan          chan struct{}
	resumeChan         chan struct{}
	doneChan           chan struct{}
}

func NewKeyboard(keys []Key, pollingRate float64) *Keyboard {
	period := time.Duration(float64(time.Second) / pollingRate)
	kb := &Keyboard{
		fetchKeyboardState: newFetchKeyboardState(keys),
		period:             period,
		mu:                 &sync.Mutex{},
		pauseChan:          make(chan struct{}),
		resumeChan:         make(chan struct{}),
		doneChan:           make(chan struct{}),
	}

	first := KeyboardState{kb.now(), make([]bool, len(keys))}
	kb.Reader.states = append(kb.Reader.states, first)
	return kb
}

func (kb *Keyboard) Poll(startTime time.Time) {
	kb.startTime = startTime
	go func() {
		for {
			select {
			case <-kb.pauseChan:
				<-kb.resumeChan // Wait until Resume() is called.
			case <-kb.doneChan:
				return
			default:
				start := times.Now()

				state := KeyboardState{times.Since(kb.startTime), kb.fetchKeyboardState()}
				kb.mu.Lock()
				kb.Reader.states = append(kb.Reader.states, state)
				kb.mu.Unlock()

				elapsed := times.Since(start)
				time.Sleep(kb.period - elapsed) // ok to pass negative value
			}
		}
	}()
}

func (kb *Keyboard) Pause() {
	kb.mu.Lock()
	defer kb.mu.Unlock()
	if kb.paused {
		return
	}

	kb.pauseChan <- struct{}{}
	kb.pauseTime = times.Now()
	kb.paused = true
}

func (kb *Keyboard) Resume() {
	kb.mu.Lock()
	defer kb.mu.Unlock()
	if !kb.paused {
		return
	}

	kb.resumeChan <- struct{}{}
	elapsedTime := times.Since(kb.pauseTime)
	kb.startTime = kb.startTime.Add(elapsedTime)
	kb.paused = false
}

func (kb *Keyboard) Close() {
	kb.mu.Lock()
	defer kb.mu.Unlock()
	// if kb.paused {
	// 	kb.resumeChan <- struct{}{}
	// 	kb.paused = false
	// }

	kb.doneChan <- struct{}{}
	close(kb.pauseChan)
	close(kb.resumeChan)
	close(kb.doneChan)
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

func isEqual(a, b KeyboardState) bool {
	for k, ap := range a.PressedList {
		bp := b.PressedList[k]
		if ap != bp {
			return false
		}
	}
	return true
}

// No worry of accessing with nil pointer.
// https://go.dev/play/p/B4Z1LwQC_jP
