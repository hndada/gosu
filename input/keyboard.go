package input

import (
	"sync"
	"time"
)

// When to use lock: read or write to shared resources.
// A practical example would be public toilet.

// No additional adjustment for keyboard when offset has changed.
// Both music and keyboard cannot seek at precise position once they start.
type Keyboard struct {
	Reader          KeyboardReader
	keyStatesGetter func() []bool
	pollingRate     time.Duration

	mu        *sync.Mutex
	startTime time.Time
	pauseTime time.Time
	pause     chan struct{}
	resume    chan struct{}
	paused    bool
	done      chan struct{}
}

func NewKeyboard(keys []Key, startTime time.Time) *Keyboard {
	kb := &Keyboard{
		keyStatesGetter: newKeyStatesGetter(keys),
		pollingRate:     PollingRate,

		mu:        &sync.Mutex{},
		startTime: startTime,
		pauseTime: time.Time{},
		pause:     make(chan struct{}),
		resume:    make(chan struct{}),
		paused:    false,
		done:      make(chan struct{}),
	}

	first := KeyboardState{kb.now(), make([]bool, len(keys))}
	kb.Reader.states = append(kb.Reader.states, first)
	return kb
}

func (kb *Keyboard) Listen() {
	go func() {
		for {
			select {
			case <-kb.pause:
				// Wait until Resume() is called.
				<-kb.resume
			case <-kb.done:
				return
			default:
				start := time.Now()

				state := KeyboardState{kb.now(), kb.keyStatesGetter()}
				kb.mu.Lock()
				kb.Reader.states = append(kb.Reader.states, state)
				kb.mu.Unlock()

				// It is fine to pass negative value to Sleep().
				elapsed := time.Since(start)
				time.Sleep(kb.pollingRate - elapsed)
			}
		}
	}()
}

func (kb *Keyboard) now() int32 {
	return int32(time.Since(kb.startTime).Milliseconds())
}

func (kb *Keyboard) Pause() {
	kb.mu.Lock()
	defer kb.mu.Unlock()

	if kb.paused {
		return
	}

	kb.pause <- struct{}{}
	kb.pauseTime = time.Now()
	kb.paused = true
}

func (kb *Keyboard) Resume() {
	kb.mu.Lock()
	defer kb.mu.Unlock()

	if !kb.paused {
		return
	}

	kb.resume <- struct{}{}
	elapsedTime := time.Since(kb.pauseTime)
	kb.startTime = kb.startTime.Add(elapsedTime)
	kb.paused = false
}

func (kb *Keyboard) Close() {
	kb.mu.Lock()
	defer kb.mu.Unlock()

	if kb.paused {
		kb.resume <- struct{}{}
		kb.paused = false
	}

	kb.done <- struct{}{}
	close(kb.pause)
	close(kb.resume)
	close(kb.done)
}
