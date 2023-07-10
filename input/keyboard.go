package input

import (
	"sync"
	"time"
)

type Keyboard interface {
	// Listen() KeyboardState // Not compatible with replay.
	Now() int32              // int32: Maximum duration is around 24 days.
	Fetch() []KeyboardAction // Fetch latest actions with 1ms unit.
	Pause()
	Resume()
	Output() []KeyboardState // Output all states.
}
type KeyboardState struct {
	Time    int32
	Pressed []bool
}
type KeyboardAction struct {
	Time   int32
	Action []KeyAction
}

type KeyboardListener struct {
	KeySettings []Key
	// vkcodes     []uint32 // used in Windows

	StartTime    time.Time // It is used for replaying.
	getKeyStates func() []bool
	PollingRate  time.Duration

	PauseTime     time.Time
	paused        bool
	pauseChannel  chan struct{}
	resumeChannel chan struct{}
	pauseMutex    *sync.Mutex

	States []KeyboardState
	index  int
}

// Windows
func NewKeyboardListener(ks []Key, bufferTime time.Duration) *KeyboardListener {
	const pollingRate = 1 * time.Millisecond

	vkcodes := make([]uint32, len(ks))
	for k, ek := range ks {
		vkcodes[k] = ToVirtualKey(ek)
	}
	getKeyStates := func() []bool {
		ps := make([]bool, len(vkcodes))
		for k, vk := range vkcodes {
			v, _, _ := procGetKeyState.Call(uintptr(vk))
			ps[k] = v&isPressed != 0
		}
		return ps
	}

	startTime := time.Now().Add(bufferTime)
	blank := KeyboardState{
		Time:    -5000,
		Pressed: make([]bool, len(ks)),
	}
	return &KeyboardListener{
		KeySettings:  ks,
		StartTime:    startTime,
		getKeyStates: getKeyStates,
		PollingRate:  pollingRate,

		pauseChannel:  make(chan struct{}),
		resumeChannel: make(chan struct{}),
		pauseMutex:    &sync.Mutex{},
		States:        []KeyboardState{blank},
	}
}

func keyActions(last, current []bool) []KeyAction {
	a := make([]KeyAction, len(current))
	for k, p := range current {
		lp := last[k]
		a[k] = keyAction(lp, p)
	}
	return a
}
func (kl *KeyboardListener) Now() int32 {
	return int32(time.Since(kl.StartTime).Milliseconds())
}
func (kl *KeyboardListener) Poll() {
	go func() {
		for {
			select {
			case <-kl.pauseChannel:
				// Paused, wait until Resume() is called
				<-kl.resumeChannel
			default:
				kl.pauseMutex.Lock()
				state := KeyboardState{kl.Now(), kl.getKeyStates()}
				kl.pauseMutex.Unlock()

				if kl.isStateChanged(state) {
					kl.States = append(kl.States, state)
				}
			}
		}
	}()
}
func (kl *KeyboardListener) Fetch() (kas []KeyboardAction) {
	last := kl.States[kl.index]

	for _, current := range kl.States[kl.index:] {
		as := keyActions(last.Pressed, current.Pressed)
		for t := last.Time; t < current.Time; t++ {
			ka := KeyboardAction{t, as}
			kas = append(kas, ka)
		}
		last = current
	}

	kl.index = len(kl.States)
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
func (kl *KeyboardListener) Pause() {
	// Add current state no matter it is same as last state.
	// This is for updating start time when resume.
	// kl.States = append(kl.States, kl.Listen())

	kl.pauseMutex.Lock()
	defer kl.pauseMutex.Unlock()
	// Signal the pause by sending a value to the pauseChannel
	kl.pauseChannel <- struct{}{}

	kl.PauseTime = time.Now()
	kl.paused = true
}
func (kl *KeyboardListener) Resume() {
	// Update start time
	// since := int32(time.Since(kl.StartTime).Milliseconds())
	// last := kl.States[kl.index].Time
	// pauseDuration := time.Duration(since-last) * time.Millisecond
	// kl.StartTime = kl.StartTime.Add(pauseDuration)

	kl.pauseMutex.Lock()
	defer kl.pauseMutex.Unlock()

	// Signal the resume by sending a value to the resumeChannel
	kl.resumeChannel <- struct{}{}

	pauseDuration := time.Since(kl.PauseTime)
	kl.StartTime = kl.StartTime.Add(pauseDuration)
	kl.paused = false
}
func (kl *KeyboardListener) Output() []KeyboardState { return kl.States }
