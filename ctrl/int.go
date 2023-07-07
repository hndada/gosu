package ctrl

type Int interface{ int | int64 }

// Todo: merge with Int64Handler into one
type IntHandler struct {
	Value *int
	Min   int
	Max   int
	Unit  int
	Loop  bool
}

func (h IntHandler) Decrease() {
	*h.Value -= h.Unit
	if *h.Value < h.Min {
		if h.Loop {
			*h.Value = h.Max
		} else {
			*h.Value = h.Min
		}
	}
}
func (h IntHandler) Increase() {
	*h.Value += h.Unit
	if *h.Value > h.Max {
		if h.Loop {
			*h.Value = h.Min
		} else {
			*h.Value = h.Max
		}
	}
}

type Int64Handler struct {
	Value *int64
	Min   int64
	Max   int64
	Unit  int64
	Loop  bool
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
