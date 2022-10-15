package ctrl

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hndada/gosu/audios"
)

type Handler interface {
	Decrease()
	Increase()
}

type BoolHandler struct {
	Value *bool
}

func (h BoolHandler) Decrease() { h.swap() }
func (h BoolHandler) Increase() { h.swap() }
func (h BoolHandler) swap() {
	if !*h.Value {
		*h.Value = true
	} else {
		*h.Value = false
	}
}

type FloatHandler struct {
	Value    *float64
	Min, Max float64
	Unit     float64
}

func (h FloatHandler) Decrease() {
	*h.Value -= h.Unit
	if *h.Value < h.Min {
		*h.Value = h.Min
	}
}
func (h FloatHandler) Increase() {
	*h.Value += h.Unit
	if *h.Value > h.Max {
		*h.Value = h.Max
	}
}

type IntHandler struct {
	Value    *int
	Min, Max int
	Loop     bool
	// Unit     int
}

func (h IntHandler) Decrease() {
	*h.Value--
	if *h.Value < h.Min {
		if h.Loop {
			*h.Value = h.Max
		} else {
			*h.Value = h.Min
		}
	}
}
func (h IntHandler) Increase() {
	*h.Value++
	if *h.Value > h.Max {
		if h.Loop {
			*h.Value = h.Min
		} else {
			*h.Value = h.Max
		}
	}
}

const (
	None = iota - 1
	Decrease
	Increase
)

type KeyHandler struct {
	Handler
	Modifiers []ebiten.Key // Handler works only when all Modifier are pressed.
	Keys      [2]ebiten.Key
	Sounds    [2][]byte // Todo: implement
	Volume    *float64

	holdIndex int
	// holdKey   ebiten.Key
	countdown int // Require to hold for a while to move a cursor.
	active    bool
}

// Update returns whether the handler has set off (triggered) or not.
func (h *KeyHandler) Update() (set bool) {
	if h.countdown > 0 {
		h.countdown--
		return
	}
	for _, k := range h.Modifiers {
		if !ebiten.IsKeyPressed(k) {
			h.reset()
			return
		}
	}
	if h.holdIndex > None && !ebiten.IsKeyPressed(h.Keys[h.holdIndex]) {
		h.reset()
	}
	for i, k := range h.Keys {
		if ebiten.IsKeyPressed(k) {
			h.holdIndex = i
			break
		}
	}
	switch h.holdIndex {
	case None:
		return
	case Decrease:
		h.Decrease()
		audios.PlayEffect(h.Sounds[0], *h.Volume)
	case Increase:
		h.Increase()
		audios.PlayEffect(h.Sounds[1], *h.Volume)
	}
	if h.active {
		h.countdown = shortCountdown
	} else {
		h.countdown = longCountdown
	}
	h.active = true
	return true
}

func (h *KeyHandler) reset() {
	h.active = false
	h.holdIndex = None
}
