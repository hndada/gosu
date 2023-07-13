package input

import (
	"sync"
	"time"
)

type KeyboardListener struct {
	KeySettings  []Key
	StartTime    time.Time
	getKeyStates func() []bool
	PollingRate  time.Duration

	pauseTime     time.Time
	paused        bool
	pauseChannel  chan struct{}
	resumeChannel chan struct{}
	pauseMutex    *sync.Mutex

	States []KeyboardState
	index  int
}

func NewKeyboardListener(keys []Key, bufferTime time.Duration) *KeyboardListener {
	blank := KeyboardState{
		Time:    -5000,
		Pressed: make([]bool, len(keys)),
	}
	return &KeyboardListener{
		KeySettings:  keys,
		StartTime:    time.Now().Add(bufferTime),
		getKeyStates: getKeyStatesFunc(keys),
		PollingRate:  PollingRate,

		pauseChannel:  make(chan struct{}),
		resumeChannel: make(chan struct{}),
		pauseMutex:    &sync.Mutex{},
		States:        []KeyboardState{blank},
		// index:         1, // Suppose blank has already fetched.
	}
}

func (kl *KeyboardListener) Now() int32 {
	return int32(time.Since(kl.StartTime).Milliseconds())
}

func (kl *KeyboardListener) Fetch() (kas []KeyboardAction) {
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
	return nil
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

				kl.pauseMutex.Lock()
				state := KeyboardState{kl.Now(), kl.getKeyStates()}
				kl.pauseMutex.Unlock()

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

func (kl *KeyboardListener) Pause() {
	kl.pauseMutex.Lock()
	defer kl.pauseMutex.Unlock()
	kl.pauseChannel <- struct{}{}

	kl.pauseTime = time.Now()
	kl.paused = true
}

func (kl *KeyboardListener) Resume() {
	kl.pauseMutex.Lock()
	defer kl.pauseMutex.Unlock()
	kl.resumeChannel <- struct{}{}

	pauseDuration := time.Since(kl.pauseTime)
	kl.StartTime = kl.StartTime.Add(pauseDuration)
	kl.paused = false
}
func (kl *KeyboardListener) IsPaused() bool { return kl.paused }

func (kl *KeyboardListener) Output() []KeyboardState { return kl.States }

func (kl *KeyboardListener) Close() {
	kl.pauseMutex.Lock()
	defer kl.pauseMutex.Unlock()
	kl.pauseChannel <- struct{}{}
	close(kl.pauseChannel)
	close(kl.resumeChannel)
}
