package draws

// Focus should be handled by parent. Only one element can be focused.
// Cursor shape should be changed by Scene.
// Each element should add required fields as needed.

// holdTick          int
type Hoverable struct {
	on      bool
	tick    int
	MaxTick int
	// ResetTickWhenOut bool
}

func (e *Hoverable) Update(isMouseIn bool) {
	if isMouseIn {
		e.on = true
		e.tick++
	} else {
		e.on = false
		e.tick--
	}
}

type Clickable struct {
	count int
	tick  int
}

func (e *Clickable) Update(clicked bool) {
	if clicked {
		e.count++
		e.tick++
	} else {
		e.count = 0
		e.tick--
	}
}

// Usage: selecting multiple charts by dragging mouse.
// type Draggable struct {}

// Usage: score board, chat
// type Collapsible struct {}
