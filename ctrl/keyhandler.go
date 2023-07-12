package ctrl

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hndada/gosu/audios"
	"github.com/hndada/gosu/input"
)

const (
	short = 80
	long  = 200
)

// Todo: Modifiers work strangely when there are plural modifiers.
// Todo: support modifiers for KeyHandler
type KeyHandler struct {
	Handler
	Modifier input.Key // Handler works only when Modifier is pressed.
	Keys     [2]input.Key
	Sounds   [2]audios.SoundPlayer
	Volume   *float64

	holdIndex int
	active    bool
	countdown int // User needs to hold for a while to activate.
}

// Update returns whether the handler has set off (triggered) or not.
// Todo: set off the stricter handler only
func (kh *KeyHandler) Update() (set bool) {
	if kh.countdown > 0 {
		kh.countdown--
		return
	}

	if !ebiten.IsKeyPressed(kh.Modifier) {
		kh.reset()
		return
	}

	if kh.holdIndex > none && !ebiten.IsKeyPressed(kh.Keys[kh.holdIndex]) {
		kh.reset()
	}

	for i, k := range kh.Keys {
		if ebiten.IsKeyPressed(k) {
			kh.holdIndex = i
			break
		}
	}
	if kh.holdIndex == none {
		return
	}

	[]func(){kh.Decrease, kh.Increase}[kh.holdIndex]()

	kh.Sounds[kh.holdIndex].Play(*kh.Volume)

	if kh.active {
		kh.countdown = short
	} else {
		kh.countdown = long
	}
	kh.active = true

	return true
}

func (kh *KeyHandler) reset() {
	kh.active = false
	kh.holdIndex = none
}
