package ui

import "golang.org/x/exp/constraints"

type BoolController struct {
	Value *bool
}

func (h *BoolController) Toggle() {
	if !*h.Value {
		*h.Value = true
	} else {
		*h.Value = false
	}
}

type Number interface {
	constraints.Integer | constraints.Float
}

type NumberController[T Number] struct {
	Value *T
	Min   T
	Max   T
	Unit  T
}

func NewNumberController[T Number](value *T, min, max, unit T) *NumberController[T] {
	return &NumberController[T]{
		Value: value,
		Min:   min,
		Max:   max,
		Unit:  unit,
	}
}

func (h *NumberController[T]) Decrease() {
	newValue := *h.Value - h.Unit
	if newValue < h.Min {
		*h.Value = h.Min
	} else {
		*h.Value = newValue
	}
}

func (h *NumberController[T]) Increase() {
	newValue := *h.Value + h.Unit
	if newValue > h.Max {
		*h.Value = h.Max
	} else {
		*h.Value = newValue
	}
}
