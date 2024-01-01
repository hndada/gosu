package ctrl

import "golang.org/x/exp/constraints"

type Handler interface {
	Decrease()
	Increase()
}

const (
	none = iota - 1
	decrease
	increase
)

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

type Number interface {
	constraints.Integer | constraints.Float
}

type ValueHandler[T Number] struct {
	Value *T
	Min   T
	Max   T
	Unit  T
}

func (h *ValueHandler[T]) Decrease() {
	newValue := *h.Value - h.Unit
	if newValue < h.Min {
		*h.Value = h.Min
	} else {
		*h.Value = newValue
	}
}

func (h *ValueHandler[T]) Increase() {
	newValue := *h.Value + h.Unit
	if newValue > h.Max {
		*h.Value = h.Max
	} else {
		*h.Value = newValue
	}
}
