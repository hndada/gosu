package ctrl

import (
	"time"

	"github.com/hndada/gosu/input"
)

// KeyListener and MouseListener just tell whether the key or mouse is pressed.
// What to do when they are fired is up to each scene:
// Play sound, change volume, draw something, etc.
type KeyListener struct {
	keys           []input.Key
	fireCount      int // Number of consecutive fires.
	firstFiredTime time.Time

	// Required durations to fire the listener:
	// 0, long, short, short, short, ...
	longDuration  time.Duration
	shortDuration time.Duration
}

func NewKeyListener(ks ...input.Key) *KeyListener {
	return &KeyListener{
		keys:          ks,
		longDuration:  500 * time.Millisecond,
		shortDuration: 100 * time.Millisecond,
	}
}

// Avoid using goroutine, it is very hard to sync other Update functions.
func (kl *KeyListener) Update() (fired bool) {
	if !kl.allKeysPressed() {
		kl.fireCount = 0
		return false
	}
	// Now key listener is active.
	switch kl.fireCount {
	case 0:
		// If it was not active, it is instantly fired.
		kl.firstFiredTime = time.Now()
		return kl.fire()
	case 1:
		if time.Since(kl.firstFiredTime) > kl.longDuration {
			return kl.fire()
		}
	default:
		minDuration := kl.longDuration +
			time.Duration(kl.fireCount-1)*kl.shortDuration
		if time.Since(kl.firstFiredTime) > minDuration {
			return kl.fire()
		}
	}
	return false
}

func (kl *KeyListener) allKeysPressed() bool {
	for _, k := range kl.keys {
		if !input.IsKeyPressed(k) {
			return false
		}
	}
	return true
}

func (kl *KeyListener) fire() bool {
	kl.fireCount++
	return true
}
