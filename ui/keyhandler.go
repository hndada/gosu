package ui

import "github.com/hndada/gosu/audios"

// Each Listener returns Control.
// A Controller then updates values by returned Control struct.
type KeyNumberHandler[T Number] struct {
	KeyListener
	*audios.SoundPlayer
	NumberController[T]
	// SoundVolumeScale *float64
}

func (h *KeyNumberHandler[T]) Handle() {
	ctrl, ok := h.KeyListener.Update()
	if !ok {
		return
	}

	h.SoundPlayer.Play(ctrl.SoundFilename)
	switch ctrl.Type {
	case Decrease:
		h.Decrease()
	case Increase:
		h.Increase()
	}
}

// BoolController is preferred to BooleanController.
type KeyBoolHandler struct {
	KeyListener
	*audios.SoundPlayer
	BoolController
}

func (h *KeyBoolHandler) Handle() {
	ctrl, ok := h.KeyListener.Update()
	if !ok {
		return
	}

	h.SoundPlayer.Play(ctrl.SoundFilename)
	h.Toggle()
}

// func NewKeyNumberHandler[T Number](c NumberController[T], s audios.SoundPlayer, l KeyListener) *KeyNumberHandler[T] {
// 	return &KeyNumberHandler[T]{
// 		NumberController: c,
// 		SoundPlayer:      s,
// 		KeyListener:      l,
// 	}
// }
