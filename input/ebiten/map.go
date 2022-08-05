package ebiten

import (
	"github.com/hajimehoshi/ebiten/v2"
)

var KeyMap = map[string]ebiten.Key{
	"A": ebiten.KeyA,
	"B": ebiten.KeyB,
	"C": ebiten.KeyC,
	"D": ebiten.KeyD,
	"E": ebiten.KeyE,
	"F": ebiten.KeyF,
	"G": ebiten.KeyG,
	"H": ebiten.KeyH,
	"I": ebiten.KeyI,
	"J": ebiten.KeyJ,
	"K": ebiten.KeyK,
	"L": ebiten.KeyL,
	"M": ebiten.KeyM,
	"N": ebiten.KeyN,
	"O": ebiten.KeyO,
	"P": ebiten.KeyP,
	"Q": ebiten.KeyQ,
	"R": ebiten.KeyR,
	"S": ebiten.KeyS,
	"T": ebiten.KeyT,
	"U": ebiten.KeyU,
	"V": ebiten.KeyV,
	"W": ebiten.KeyW,
	"X": ebiten.KeyX,
	"Y": ebiten.KeyY,
	"Z": ebiten.KeyZ,

	"0": ebiten.Key0,
	"1": ebiten.Key1,
	"2": ebiten.Key2,
	"3": ebiten.Key3,
	"4": ebiten.Key4,
	"5": ebiten.Key5,
	"6": ebiten.Key6,
	"7": ebiten.Key7,
	"8": ebiten.Key8,
	"9": ebiten.Key9,

	"Space": ebiten.KeySpace,
	".":     ebiten.KeyPeriod,
	",":     ebiten.KeyComma,
	"Ctrl":  ebiten.KeyControl,
	"Esc":   ebiten.KeyEscape,
	"Alt":   ebiten.KeyAlt,
	"Shift": ebiten.KeyShift,
	"Left":  ebiten.KeyLeft,
	"Right": ebiten.KeyRight,
	"Up":    ebiten.KeyUp,
	"Down":  ebiten.KeyDown,
	"Enter": ebiten.KeyEnter,
}

// for c := '0'; c <= '9'; c++ {
// fmt.Printf(`"%c": ebiten.Key%c,
// `, c, c)
// }
