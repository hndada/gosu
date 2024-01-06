package input

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

func MouseCursorPosition() (float64, float64) {
	x, y := ebiten.CursorPosition()
	return float64(x), float64(y)
}

// functions
var IsMouseButtonPressed = ebiten.IsMouseButtonPressed
var IsMouseButtonJustPressed = inpututil.IsMouseButtonJustPressed
var MouseWheelPosition = ebiten.Wheel

type MouseButton = ebiten.MouseButton

const (
	MouseButtonLeft   MouseButton = ebiten.MouseButtonLeft
	MouseButtonMiddle MouseButton = ebiten.MouseButtonMiddle
	MouseButtonRight  MouseButton = ebiten.MouseButtonRight
)
