package ui

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

type Event struct {
	From *Box
	Time time.Time
	Type EventType
}

// It would be better to handle events in each struct.
// type EventHandler func() any

// Type, Button value, State
// User may implement their own custom events.
type EventType [3]int

const (
	TypeKeyButton = iota
	TypeMouseButton
	TypeMouseCursor
	TypeMouseWheel
	TypeBox
)

// MouseButton: input package
const (
	MouseButtonLeft   = int(ebiten.MouseButtonLeft)
	MouseButtonMiddle = int(ebiten.MouseButtonMiddle)
	MouseButtonRight  = int(ebiten.MouseButtonRight)
	// MouseCursor       = ebiten.MouseButtonMax + 1
)

const (
	Focus = iota
)

const Any = 0

const (
	On = iota
	Off

	// All buttons including keyboard and mouse are either Pressed or Released.
	Pressed  = On
	Released = Off

	In  = On
	Out = Off
)

const (
	Up = iota
	Down
)

var (
	KeyButtonPressed          = EventType{TypeKeyButton, Any, Pressed}
	KeyButtonReleased         = EventType{TypeKeyButton, Any, Released}
	MouseClick                = EventType{TypeMouseButton, Any, Pressed}
	MouseButtonLeftPressed    = EventType{TypeMouseButton, MouseButtonLeft, Pressed}
	MouseButtonLeftReleased   = EventType{TypeMouseButton, MouseButtonLeft, Released}
	MouseButtonMiddlePressed  = EventType{TypeMouseButton, MouseButtonMiddle, Pressed}
	MouseButtonMiddleReleased = EventType{TypeMouseButton, MouseButtonMiddle, Released}
	MouseButtonRightPressed   = EventType{TypeMouseButton, MouseButtonRight, Pressed}
	MouseButtonRightReleased  = EventType{TypeMouseButton, MouseButtonRight, Released}
	MouseCursorIn             = EventType{TypeMouseCursor, Any, In}
	MouseCursorOut            = EventType{TypeMouseCursor, Any, Out}
	MouseScrollUp             = EventType{TypeMouseWheel, Any, Up}
	MouseScrollDown           = EventType{TypeMouseWheel, Any, Down}
	WidgetFocus               = EventType{TypeBox, Focus, On}
	WidgetUnfocus             = EventType{TypeBox, Focus, Off}
)
