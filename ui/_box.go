package ui

import (
	"github.com/hndada/gosu/draws"
)

type Box struct {
	draws.Box
	BoxOptions
}

type BoxOptions struct {
	Scrollable bool
	// Draggable  bool
	// ResizableOptions
}

func NewBox(dbox draws.Box, opts BoxOptions) *Box {
	return &Box{Box: dbox, BoxOptions: opts}
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

func (b *Box) Update() []Event {
	if b.Scrollable {
		b.updateScroll()
	}
	return nil
}

func (b *Box) updateScroll() {
	// If cursor is in the box, and mouse is just pressed, and so on
}
