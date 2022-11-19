package ctrl

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
