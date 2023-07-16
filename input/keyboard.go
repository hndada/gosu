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
	KeyboardReader
	keyStatesGetter func() []bool
	pollingRate     time.Duration

	mu        *sync.Mutex
	startTime time.Time
	pauseTime time.Time
	pause     chan struct{}
	resume    chan struct{}
	paused    bool
}

func NewKeyboard(keys []Key, startTime time.Time) *Keyboard {
	return &Keyboard{
		keyStatesGetter: newKeyStatesGetter(keys),
		pollingRate:     PollingRate,

		mu:        &sync.Mutex{},
		startTime: startTime,
		pauseTime: time.Time{},
		pause:     make(chan struct{}),
		resume:    make(chan struct{}),
		paused:    false,
	}
}

func (kb *Keyboard) Listen() {
	go func() {
		for {
			select {
			case <-kb.pause:
				// Wait until Resume() is called.
				<-kb.resume
			default:
				start := time.Now()

				state := KeyboardState{kb.now(), kb.keyStatesGetter()}
				kb.mu.Lock()
				kb.states = append(kb.states, state)
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

	close(kb.pause)
	close(kb.resume)
}
