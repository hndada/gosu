package ctrl

import "github.com/hajimehoshi/ebiten/v2"

// Type, Kind, State
// User may implement their own custom events.
type Event [3]int

const (
	TypeKeyButton = iota
	TypeMouseButton
	TypeMouseCursor
	TypeMouseScroll
	TypeWidget
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
	KeyButtonPressed          = Event{TypeKeyButton, Any, Pressed}
	KeyButtonReleased         = Event{TypeKeyButton, Any, Released}
	MouseClick                = Event{TypeMouseButton, Any, Pressed}
	MouseButtonLeftPressed    = Event{TypeMouseButton, MouseButtonLeft, Pressed}
	MouseButtonLeftReleased   = Event{TypeMouseButton, MouseButtonLeft, Released}
	MouseButtonMiddlePressed  = Event{TypeMouseButton, MouseButtonMiddle, Pressed}
	MouseButtonMiddleReleased = Event{TypeMouseButton, MouseButtonMiddle, Released}
	MouseButtonRightPressed   = Event{TypeMouseButton, MouseButtonRight, Pressed}
	MouseButtonRightReleased  = Event{TypeMouseButton, MouseButtonRight, Released}
	MouseCursorIn             = Event{TypeMouseCursor, Any, In}
	MouseCursorOut            = Event{TypeMouseCursor, Any, Out}
	MouseScrollUp             = Event{TypeMouseScroll, Any, Up}
	MouseScrollDown           = Event{TypeMouseScroll, Any, Down}
	WidgetFocus               = Event{TypeWidget, Focus, On}
	WidgetUnfocus             = Event{TypeWidget, Focus, Off}
)
