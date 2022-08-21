package ctrl

import (
	"github.com/hajimehoshi/ebiten/v2"
)

const (
	HandlerKeyIncrease = iota
	HandlerKeyDecrease
	HandlerKeyNext
	HandlerKeyPrev
)

var HandlerKeyTypes = []int{
	HandlerKeyIncrease,
	HandlerKeyDecrease,
	HandlerKeyNext,
	HandlerKeyPrev,
}

type Handler struct {
	Keys       []ebiten.Key
	PlaySounds []func()
	HoldKey    ebiten.Key
	Countdown  int // Require to hold for a while to move a cursor.
	Active     bool
}

const KeyNone = -1

func (h *Handler) Update() {
	h.Countdown--
	if h.Countdown < 0 {
		h.Countdown = 0
	}
	if ebiten.IsKeyPressed(h.HoldKey) {
		return
	}
	h.Active = false
	h.Countdown = 0
	h.HoldKey = KeyNone
	for _, keyType := range HandlerKeyTypes[:len(h.Keys)] {
		key := h.Keys[keyType]
		if ebiten.IsKeyPressed(key) {
			h.HoldKey = key
		}
	}
}

func (h Handler) KeyType() int {
	for i, t := range HandlerKeyTypes[:len(h.Keys)] {
		if h.HoldKey == h.Keys[i] {
			return t
		}
	}
	return -1
}

type F64Handler struct {
	Handler
	Min    float64
	Max    float64
	Unit   float64
	Target *float64
}

// Update returns whether the handler has fired or not.
func (h *F64Handler) Update() bool {
	h.Handler.Update()
	if h.Countdown > 0 || h.KeyType() == -1 {
		return false
	}
	// Now countdown is 0.
	h.PlaySounds[h.KeyType()]()
	switch h.KeyType() {
	case HandlerKeyIncrease:
		*h.Target += h.Unit
		if *h.Target > h.Max {
			*h.Target = h.Max
		}
	case HandlerKeyDecrease:
		*h.Target -= h.Unit
		if *h.Target < h.Min {
			*h.Target = h.Min
		}
	case HandlerKeyNext:
	case HandlerKeyPrev:
	}
	if h.Active {
		h.Countdown = shortCountdown
	} else {
		h.Countdown = longCountdown
	}
	h.Active = true
	return true
}

type IntHandler struct {
	Handler
	Min    int
	Max    int
	Unit   int
	Target *int
	Loop   bool
}

// Update returns whether the handler has fired or not.
func (h *IntHandler) Update() bool {
	h.Handler.Update()
	if h.Countdown > 0 || h.KeyType() == -1 {
		return false
	}
	h.PlaySounds[h.KeyType()]()
	switch h.KeyType() {
	case HandlerKeyIncrease:
		*h.Target += h.Unit
		if *h.Target >= h.Max {
			if h.Loop {
				*h.Target -= h.Max
			} else {
				*h.Target = h.Max
			}
		}
	case HandlerKeyDecrease:
		*h.Target -= h.Unit
		if *h.Target < h.Min {
			if h.Loop {
				*h.Target += h.Max
			} else {
				*h.Target = h.Min
			}
		}
	case HandlerKeyNext:
	case HandlerKeyPrev:
	}
	if h.Active {
		h.Countdown = shortCountdown
	} else {
		h.Countdown = longCountdown
	}
	h.Active = true
	return true
}

type BoolHandler struct {
	Handler
	Target *bool
}

// Update returns whether the handler has fired or not.
func (h *BoolHandler) Update() bool {
	h.Handler.Update()
	// Bool value is updated only once regardless of hold duration.
	if h.Countdown > 0 || h.KeyType() == -1 || h.Active {
		return false
	}
	switch h.KeyType() {
	case HandlerKeyIncrease:
		*h.Target = !*h.Target
	case HandlerKeyDecrease:
	case HandlerKeyNext:
	case HandlerKeyPrev:
	}
	h.Active = true
	return true
}
