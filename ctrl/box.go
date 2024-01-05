package ctrl

import (
	"time"

	"github.com/hndada/gosu/draws"
)

type Box2 struct {
	draws.Box
	BoxOptions
}

func (b *Box2) Update() []Event2 {
	if b.Draggable {
		b.updateDrag()
	}
	if b.Resizable {
		b.updateResize()
	}
	if b.Scrollable {
		b.updateScroll()
	}
	return nil
}

func (b *Box2) updateDrag() {
	// If cursor is in the box, and mouse is just pressed, and so on
}

type BoxOptions struct {
	Draggable bool
	ResizableOptions
	Scrollable bool
	// It would be better to handle events in each struct.
	// EventHandler map[EventType]func()
}

// Event types: MouseDown, MouseUp, MouseMove, MouseEnter, MouseLeave, MouseScroll,
// KeyDown, KeyUp, KeyPress,
// Focus, Blur, Resize, Move, Close
// HoverIn, HoverOut, Click, DoubleClick, DragStart, Drag, DragEnd,
// ResizeStart, Resize, ResizeEnd
type ResizableOptions struct {
	Resizable bool
	MinSize   draws.Vector2
	MaxSize   draws.Vector2
}

type Event2 struct {
	From *Box
	Time time.Time
	Type EventType
}
