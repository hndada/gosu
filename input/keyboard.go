package input

import (
	"time"

	"github.com/hndada/gosu/times"
)

// type Keystrokes struct?
type KeyboardState struct {
	Time        time.Duration // Stands for elapsed time.
	PressedList []bool
}

// KeyboardStateBuffer is supposed to have at least one state.
type KeyboardStateBuffer struct {
	buf []KeyboardState
	idx int
}

func NewKeyboardStateBuffer(buf []KeyboardState) KeyboardStateBuffer {
	return KeyboardStateBuffer{
		buf: buf,
	}
}

// Read returns the last read state and unread states before given time.
func (kb *KeyboardStateBuffer) Read(now time.Duration) (kss []KeyboardState) {
	for _, state := range kb.buf[kb.idx:] {
		if state.Time > now {
			break
		}
		kss = append(kss, state)
	}
	// To make the index pointing at the last state.
	kb.idx += len(kss) - 1
	return kss
}

// Output trims redundant states then returns the states.
func (kb *KeyboardStateBuffer) Trim() {
	trimmed := make([]KeyboardState, 1, len(kb.buf))
	copy(trimmed, kb.buf)

	old := kb.buf[0]
	for _, now := range kb.buf[1:] {
		if isEqual(old, now) {
			continue
		}
		trimmed = append(trimmed, now)
		old = now
	}
	kb.buf = trimmed
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

func (kb KeyboardStateBuffer) Output() []KeyboardState {
	kb.Trim()
	return kb.buf
}

type KeyboardReader interface {
	Read(now time.Duration) []KeyboardState
}

type KeyboardListener interface {
	Listen()
	Stop()
}

// Keyboard should not require additional adjustment when offset has changed,
// Because Keyboard cannot seek at precise position once it starts. Same goes for music.
type Keyboard struct {
	KeyboardStateBuffer
	// mu                 *sync.Mutex // for lock
	fetchKeyboardState func() []bool
	startTime          time.Time
	period             time.Duration
	stop               chan struct{}
}

func NewKeyboard(keys []Key) Keyboard {
	kb := Keyboard{
		fetchKeyboardState: newFetchKeyboardState(keys),
		// mu:                 &sync.Mutex{},
		stop: make(chan struct{}),
	}
	first := KeyboardState{-10 * time.Second, make([]bool, len(keys))}
	kb.buf = append(kb.buf, first)
	kb.SetPollingRate(defaultPollingRate)
	return kb
}

func (kb *Keyboard) SetPollingRate(rate float64) {
	second := float64(time.Second) * times.PlaybackRate()
	kb.period = time.Duration(second / rate)
}

// Listen starts polling keyboard state.
func (kb *Keyboard) Listen(startTime time.Time) {
	kb.startTime = startTime
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

// When to use lock: read or write to shared resources.
// A practical example would be public toilet.
func (kb *Keyboard) poll() {
	t := times.Since(kb.startTime)
	ps := kb.fetchKeyboardState()
	state := KeyboardState{t, ps}

	// kb.mu.Lock()
	// defer kb.mu.Unlock()
	kb.buf = append(kb.buf, state)
}

func (kb *Keyboard) Stop() {
	kb.stop <- struct{}{}
}

// No worry of accessing with nil pointer.
// https://go.dev/play/p/B4Z1LwQC_jP
