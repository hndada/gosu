package ui

import (
	"time"

	"github.com/hndada/gosu/input"
)

func IsEnterJustPressed() bool {
	return input.IsKeyJustPressed(input.KeyEnter) || input.IsKeyJustPressed(input.KeyNumpadEnter)
}
func IsEscapeJustPressed() bool {
	return input.IsKeyJustPressed(input.KeyEscape)
}

// KeyListener and MouseListener just tell whether the key or mouse is pressed.
// What to do when they are fired is up to each scene:
// Play sound, change volume, draw something, etc.

// It is common that key listeners check the very last key press.
// Tested with Ctrl+B and Ctrl+N in VS Code:
// Ctrl+B+N: opened a new tab repeatedly.
// Ctrl+N+B: collapsed the sidebar and uncollapsed it repeatedly.
type KeyListener struct {
	*KeyboardStatus
	Modifiers []input.Key
	Controls  []Control
	keys      []input.Key
	// Keys           []input.Key
	fireCount      int // Number of consecutive fires.
	firstFiredTime time.Time

	// Required durations to fire the listener:
	// 0, long, short, short, short, ...
	longDuration  time.Duration
	shortDuration time.Duration
}

func NewKeyListener(ks *KeyboardStatus, mods []input.Key, ctrls []Control) *KeyListener {
	keys := make([]input.Key, len(ctrls))
	for i, c := range ctrls {
		keys[i] = c.Key
	}

	return &KeyListener{
		KeyboardStatus: ks,
		Modifiers:      mods,
		Controls:       ctrls,
		keys:           keys,
		longDuration:   500 * time.Millisecond,
		shortDuration:  100 * time.Millisecond,
	}
}

// Avoid using goroutine, it is very hard to sync other Update functions.
// Declaring local functions in a method instead of separating them as methods seems fine.
func (kh *KeyListener) Update() (Control, bool) {
	reset := func() (Control, bool) {
		kh.fireCount = 0
		return Control{}, false
	}
	if !kh.AreAllKeysPressed(kh.Modifiers) {
		return reset()
	}

	k, ok := kh.AreAnyKeysPressed(kh.keys)
	if !ok {
		return reset()
	}

	var c Control
	for i := 0; i < len(kh.keys); i++ {
		if k == kh.keys[i] {
			c = kh.Controls[i]
			break
		}
	}

	fire := func() (Control, bool) {
		kh.fireCount++
		return c, true
	}

	// Now key listener is active.
	switch kh.fireCount {
	case 0:
		// If it was not active, it is instantly fired.
		kh.firstFiredTime = time.Now()
		return fire()
	case 1:
		if time.Since(kh.firstFiredTime) > kh.longDuration {
			return fire()
		}
	default:
		minDuration := kh.longDuration +
			time.Duration(kh.fireCount-1)*kh.shortDuration
		if time.Since(kh.firstFiredTime) > minDuration {
			return fire()
		}
	}
	return Control{}, false
}
