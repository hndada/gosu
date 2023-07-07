package ctrl

type FloatHandler struct {
	Value *float64
	Min   float64
	Max   float64
	Unit  float64
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
