package ui

import "golang.org/x/exp/constraints"

type BoolControl struct {
	Value *bool
}

func (h *BoolControl) Switch() {
	if !*h.Value {
		*h.Value = true
	} else {
		*h.Value = false
	}
}

type Number interface {
	constraints.Integer | constraints.Float
}

type NumberControl[T Number] struct {
	Value *T
	Min   T
	Max   T
	Unit  T
}

func NewNumberControl[T Number](value *T, min, max, unit T) *NumberControl[T] {
	return &NumberControl[T]{
		Value: value,
		Min:   min,
		Max:   max,
		Unit:  unit,
	}
}

func (h *NumberControl[T]) Decrease() {
	newValue := *h.Value - h.Unit
	if newValue < h.Min {
		*h.Value = h.Min
	} else {
		*h.Value = newValue
	}
}

func (h *NumberControl[T]) Increase() {
	newValue := *h.Value + h.Unit
	if newValue > h.Max {
		*h.Value = h.Max
	} else {
		*h.Value = newValue
	}
}
