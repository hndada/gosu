package input

import (
	"sync"
	"time"
)

// When to use lock: read or write to shared resources.
// A practical example would be public toilet.
type KeyboardListener struct {
	Keyboard
	keyStatesGetter func() []bool
	pollingRate     time.Duration

	mu        *sync.Mutex
	startTime time.Time
	pauseTime time.Time
	pause     chan struct{}
	resume    chan struct{}
	paused    bool
}

func NewKeyboardListener(keys []Key, startTime time.Time) *KeyboardListener {
	return &KeyboardListener{
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

func (kl *KeyboardListener) Now() int32 {
	return int32(time.Since(kl.startTime).Milliseconds())
}

func (kl *KeyboardListener) Poll() {
	go func() {
		for {
			select {
			case <-kl.pause:
				// Wait until Resume() is called.
				<-kl.resume
			default:
				start := time.Now()

				state := KeyboardState{kl.Now(), kl.keyStatesGetter()}
				kl.mu.Lock()
				kl.states = append(kl.states, state)
				kl.mu.Unlock()

				// It is fine to pass negative value to Sleep().
				elapsed := time.Since(start)
				time.Sleep(kl.pollingRate - elapsed)
			}
		}
	}()
}

func (kl KeyboardListener) IsPaused() bool { return kl.paused }

func (kl *KeyboardListener) Pause() {
	kl.mu.Lock()
	defer kl.mu.Unlock()

	if kl.paused {
		return
	}

	kl.pause <- struct{}{}
	kl.pauseTime = time.Now()
	kl.paused = true
}

func (kl *KeyboardListener) Resume() {
	kl.mu.Lock()
	defer kl.mu.Unlock()

	if !kl.paused {
		return
	}

	kl.resume <- struct{}{}
	elapsedTime := time.Since(kl.pauseTime)
	kl.startTime = kl.startTime.Add(elapsedTime)
	kl.paused = false
}

func (kl *KeyboardListener) Close() {
	kl.mu.Lock()
	defer kl.mu.Unlock()

	if kl.paused {
		kl.resume <- struct{}{}
		kl.paused = false
	}

	close(kl.pause)
	close(kl.resume)
}
