package input

import (
	"sync"
	"time"
)

type KeyboardListener struct {
	KeySettings  []Key
	startTime    time.Time
	getKeyStates func() []bool
	PollingRate  time.Duration

	pauseTime     time.Time
	paused        bool
	pauseChannel  chan struct{}
	resumeChannel chan struct{}
	mu            *sync.Mutex

	States []KeyboardState
	index  int
}

func NewKeyboardListener(keys []Key) *KeyboardListener {
	var bufferTime time.Duration
	blank := KeyboardState{
		Time:    -5000,
		Pressed: make([]bool, len(keys)),
	}
	return &KeyboardListener{
		KeySettings:  keys,
		startTime:    time.Now().Add(bufferTime),
		getKeyStates: getKeyStatesFunc(keys),
		PollingRate:  PollingRate,

		pauseChannel:  make(chan struct{}),
		resumeChannel: make(chan struct{}),
		mu:            &sync.Mutex{},
		States:        []KeyboardState{blank},
		// index:         1, // Suppose blank has already fetched.
	}
}

func (kl *KeyboardListener) Now() int32 {
	return int32(time.Since(kl.startTime).Milliseconds())
}

func (kl *KeyboardListener) Fetch(now int32) []KeyboardAction {
	var kas []KeyboardAction
	// last := kl.States[kl.index-1]
	last := kl.States[kl.index]

	// From last fetched state to latest state
	for _, current := range kl.States[kl.index+1:] {
		as := KeyActions(last.Pressed, current.Pressed)
		for t := last.Time; t < current.Time; t++ {
			ka := KeyboardAction{t, as}
			kas = append(kas, ka)
		}
		last = current
	}

	// From latest state to now
	las := KeyActions(last.Pressed, last.Pressed)
	for t := last.Time; t < kl.Now(); t++ {
		ka := KeyboardAction{t, las}
		kas = append(kas, ka)
	}

	// Update index
	kl.index = len(kl.States) - 1
	return kas
}

func (kl KeyboardListener) isStateChanged(state KeyboardState) bool {
	last := kl.States[len(kl.States)-1].Pressed
	current := state.Pressed
	for k, lp := range last {
		p := current[k]
		if lp != p {
			return true
		}
	}
	return false
}

func (kl *KeyboardListener) Poll() {
	go func() {
		for {
			select {
			case <-kl.pauseChannel:
				<-kl.resumeChannel // wait until Resume() is called
			default:
				start := time.Now()

				kl.mu.Lock()
				state := KeyboardState{kl.Now(), kl.getKeyStates()}
				kl.mu.Unlock()

				if kl.isStateChanged(state) {
					kl.States = append(kl.States, state)
				}

				elapsed := time.Since(start)
				// It is fine to pass negative value to Sleep().
				time.Sleep(kl.PollingRate - elapsed)
			}
		}
	}()
}

func (kl *KeyboardListener) IsPaused() bool { return kl.paused }

func (kl *KeyboardListener) Pause() {
	kl.mu.Lock()
	kl.pauseChannel <- struct{}{}
	kl.pauseTime = time.Now()
	kl.paused = true
	kl.mu.Unlock()
}

func (kl *KeyboardListener) Resume() {
	kl.mu.Lock()
	kl.resumeChannel <- struct{}{}
	elapsedTime := time.Since(kl.pauseTime)
	kl.startTime = kl.startTime.Add(elapsedTime)
	kl.paused = false
	kl.mu.Unlock()
}

func (kl *KeyboardListener) Output() []KeyboardState { return kl.States }

func (kl *KeyboardListener) Close() {
	kl.mu.Lock()
	kl.pauseChannel <- struct{}{}
	close(kl.pauseChannel)
	close(kl.resumeChannel)
	kl.mu.Unlock()
}
