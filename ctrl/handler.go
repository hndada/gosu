package ctrl

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hndada/gosu/audios"
	"github.com/hndada/gosu/input"
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

// type SliceHandler struct {
// 	IntHandler
// }

// func NewSliceHandler(slice []any) (h SliceHandler) {
// 	var index int
// 	h.IntHandler = IntHandler{
// 		Value: &index,
// 		Min:   0,
// 		Max:   len(slice),
// 		Loop:  false,
// 	}
// 	return
// }

const (
	none = iota - 1
	decrease
	increase
)

type KeyHandler struct {
	Handler
	Modifiers []input.Key // Handler works only when all Modifier are pressed.
	Keys      [2]input.Key
	Sounds    [2]audios.Sound // Todo: implement
	Volume    *float64

	holdIndex int
	// holdKey   input.Key
	countdown int // Require to hold for a while to move a cursor.
	active    bool
}

// Update returns whether the handler has set off (triggered) or not.
// Todo: set off the stricter handler only
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
	if h.holdIndex > none && !ebiten.IsKeyPressed(h.Keys[h.holdIndex]) {
		h.reset()
	}
	for i, k := range h.Keys {
		if ebiten.IsKeyPressed(k) {
			h.holdIndex = i
			break
		}
	}
	if h.holdIndex == none {
		return
	}
	[]func(){h.Decrease, h.Increase}[h.holdIndex]()
	h.Sounds[h.holdIndex].Play(*h.Volume)
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
	h.holdIndex = none
}
