package ctrl

type Int interface{ int | int64 }

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

// Todo: merge with IntHandler into one
// Todo: add Unit to IntHandler?
type Int64Handler struct {
	Value    *int64
	Min, Max int64
	Loop     bool
	Unit     int64
}

func (h Int64Handler) Decrease() {
	*h.Value -= h.Unit
	if *h.Value < h.Min {
		if h.Loop {
			*h.Value = h.Max
		} else {
			*h.Value = h.Min
		}
	}
}
func (h Int64Handler) Increase() {
	*h.Value += h.Unit
	if *h.Value > h.Max {
		if h.Loop {
			*h.Value = h.Min
		} else {
			*h.Value = h.Max
		}
	}
}
