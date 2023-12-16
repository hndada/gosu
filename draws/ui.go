package draws

// focus는 상위에서 다뤄줘야 한다. 단 하나만 focus될 수 있기 때문.
// Cursor 모양은 Scene에서 바꾼다.
// 필요한 field는 각 element가 추가하기로 한다.
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
