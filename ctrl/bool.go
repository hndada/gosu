package ctrl

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
