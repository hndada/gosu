package ui

import (
	draws "github.com/hndada/gosu/draws5"
	"github.com/hndada/gosu/input"
)

// draws package is for drawing.
// ui: draws + input
// ui package is defining some well-known widgets.

// Cursor, Button, Wheel
// Focus should be handled by scene. Only one box can be focused.
type MouseListener struct {
	box                *draws.Box
	lastCursorIn       bool
	cursorIn           bool
	lastButtonsPressed [3]bool
	buttonsPressed     [3]bool
	lastWheelPosition  draws.XY
	wheelPosition      draws.XY
	// scrollable         bool
	// scrollScale        float64
}

// Returning pointer would be better for manipulating
// the same box from other places.
func NewMouseListener(box *draws.Box) *MouseListener {
	return &MouseListener{box: box}
}

// func (ml *MouseListener) SetScrollOptions(scrollable bool, scale float64) {
// 	ml.scrollable = scrollable
// 	ml.scrollScale = scale
// }

func (ml *MouseListener) Update() {
	ml.updateCursor()
	ml.updateButtons()
	ml.updateWheel()
}

func (ml *MouseListener) updateCursor() {
	cp := draws.NewXY(input.MouseCursorPosition())
	ml.lastCursorIn = ml.cursorIn
	ml.cursorIn = ml.box.In(cp)
	if ml.cursorIn {
		wp := draws.NewXY(input.MouseWheelPosition())
		ml.lastWheelPosition = wp
	}
}

func (ml *MouseListener) updateButtons() {
	ml.lastButtonsPressed = ml.buttonsPressed
	if ml.IsCursorEntered() {
		ml.buttonsPressed = [3]bool{
			input.IsMouseButtonPressed(input.MouseButtonLeft),
			input.IsMouseButtonPressed(input.MouseButtonMiddle),
			input.IsMouseButtonPressed(input.MouseButtonRight),
		}
	} else {
		ml.buttonsPressed = [3]bool{false, false, false}
	}
}

// Todo: tweening
func (ml *MouseListener) updateWheel() {
	if !ml.IsCursorEntered() {
		return
	}
	ml.lastWheelPosition = ml.wheelPosition
	wp := draws.NewXY(input.MouseWheelPosition())
	ml.wheelPosition = wp

	// dx, dy := input.MouseWheelPosition()
	// ml.box.AddPixelToX(dx * ml.scrollScale)
	// ml.box.AddPixelToY(dy * ml.scrollScale)
}

// There are 4 * 2 = 8 basic functions.
// State: In, Out, Pressed, Released
// Tense: Just or Keep
func (ml MouseListener) IsCursorJustEntered() bool {
	return !ml.lastCursorIn && ml.cursorIn
}

func (ml MouseListener) IsCursorEntered() bool {
	return ml.cursorIn
}

func (ml MouseListener) IsCursorJustExited() bool {
	return ml.lastCursorIn && !ml.cursorIn
}

func (ml MouseListener) IsCursorExited() bool {
	return !ml.cursorIn
}

func (ml MouseListener) IsButtonJustPressed(kind input.MouseButton) bool {
	return !ml.lastButtonsPressed[kind] && ml.buttonsPressed[kind]
}

func (ml MouseListener) IsButtonPressed(kind input.MouseButton) bool {
	return ml.buttonsPressed[kind]
}

// Released requires cursor to be entered.
func (ml MouseListener) IsButtonJustReleased(kind input.MouseButton) bool {
	if !ml.IsCursorEntered() {
		return false
	}
	return ml.lastButtonsPressed[kind] && !ml.buttonsPressed[kind]
}

// Released requires cursor to be entered.
func (ml MouseListener) IsButtonReleased(kind input.MouseButton) bool {
	if !ml.IsCursorEntered() {
		return false
	}
	return !ml.buttonsPressed[kind]
}

func (ml MouseListener) IsClicked(kind input.MouseButton) bool {
	return ml.IsButtonJustReleased(kind)
}

func (ml MouseListener) MouseWheelMovement() draws.XY {
	return ml.wheelPosition.Sub(ml.lastWheelPosition)
}

// Usage: selecting multiple charts by dragging mouse.
// type Draggable struct {}

// Usage: score board, chat
// type Collapsible struct {}
