package ui

import (
	"time"

	"github.com/hndada/gosu/input"
	"github.com/hndada/gosu/times"
)

const minUpdateInterval = 500 * time.Millisecond

// In most cases, keyboard shortcuts are designed to be sensitive to
// the order in which the modifier keys are pressed. So, for example,
// Ctrl+Alt+A might perform a different action than Alt+Ctrl+A.
// This implies that game should store a state of the keyboard globally.
type KeyboardStatus struct {
	lastUpdateTime time.Time
	pressedKeys    []input.Key
}

func NewKeyboardStatus() *KeyboardStatus {
	return &KeyboardStatus{}
}

// AreAllKeysPressed is used for checking modifier keys.
func (ks *KeyboardStatus) AreAllKeysPressed(keys []input.Key) bool {
	if times.Now().Sub(ks.lastUpdateTime) > minUpdateInterval {
		ks.update()
	}

	for _, m := range keys {
		if !input.IsKeyPressed(m) {
			return false
		}
	}
	return true
}

// AreAnyKeysPressed is used for checking non-modifier keys.
func (ks *KeyboardStatus) AreAnyKeysPressed(keys []input.Key) (input.Key, bool) {
	if times.Now().Sub(ks.lastUpdateTime) > minUpdateInterval {
		ks.update()
	}

	for _, k := range keys {
		if input.IsKeyPressed(k) {
			return k, true
		}
	}
	return 0, false
}

func (ks *KeyboardStatus) update() {
	newPressedKeys := make([]input.Key, 0, len(ks.pressedKeys))
	for _, pk := range ks.pressedKeys {
		if input.IsKeyPressed(pk) {
			newPressedKeys = append(newPressedKeys, pk)
		}
	}
	ks.pressedKeys = newPressedKeys

	isAlreadyPressed := func(k input.Key) bool {
		for _, pk := range ks.pressedKeys {
			if pk == k {
				return true
			}
		}
		return false
	}

	for k := input.Key(0); k < input.KeyFinal; k++ {
		if input.IsKeyPressed(k) {
			if isAlreadyPressed(k) {
				continue
			}
			ks.pressedKeys = append(ks.pressedKeys, k)
		}
	}

	ks.lastUpdateTime = times.Now()
}
