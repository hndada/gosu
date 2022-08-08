package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hndada/gosu/input/hook"
)

// Todo: separate hook
var (
	MaxTPS      = 1000
	ScreenSizeX = 800
	ScreenSizeY = 600

	Volume    = 0.05
	Speed     = 0.16
	KeyLayout = map[int][]hook.Code{
		4:               {hook.CodeD, hook.CodeF, hook.CodeJ, hook.CodeK},
		5:               {hook.CodeD, hook.CodeF, hook.CodeSpacebar, hook.CodeJ, hook.CodeK},
		6:               {hook.CodeS, hook.CodeD, hook.CodeF, hook.CodeJ, hook.CodeK, hook.CodeL},
		7:               {hook.CodeS, hook.CodeD, hook.CodeF, hook.CodeSpacebar, hook.CodeJ, hook.CodeK, hook.CodeL},
		8 + LeftScratch: {hook.CodeA, hook.CodeS, hook.CodeD, hook.CodeF, hook.CodeSpacebar, hook.CodeJ, hook.CodeK, hook.CodeL},
		8:               {hook.CodeA, hook.CodeS, hook.CodeD, hook.CodeF, hook.CodeJ, hook.CodeK, hook.CodeL, hook.CodeSemicolon},
		9:               {hook.CodeA, hook.CodeS, hook.CodeD, hook.CodeF, hook.CodeSpacebar, hook.CodeJ, hook.CodeK, hook.CodeL, hook.CodeSemicolon},
		10:              {hook.CodeA, hook.CodeS, hook.CodeD, hook.CodeF, hook.CodeV, hook.CodeN, hook.CodeJ, hook.CodeK, hook.CodeL, hook.CodeSemicolon},
	}
	// Scaled to 800 x 600.
	NoteWidths = map[int][4]float64{
		4:  {50, 50, 50, 50},
		5:  {50, 50, 50, 50},
		6:  {50, 50, 50, 50},
		7:  {50, 50, 50, 50},
		8:  {45, 45, 45, 45},
		9:  {45, 45, 45, 45},
		10: {45, 45, 45, 45},
	}
	NoteHeigth    float64 = 30 // Applies all notes
	ComboPosition float64 = 180
	JudgePosition float64 = 250
	HintPosition  float64 = 550 // The middle position of Judge line, not a topmost.
	FieldDark     float64 = 0.95
	BgDimness     float64 = 0.3
	// ScratchMode map[int]int
	ComboWidth    float64 = 40
	ScoreWidth    float64 = 33
	ComboGap      float64 = -2
	ScoreGap      float64 = -2
	HintHeight    float64 = 5
	JudgmentWidth float64 = 65
)

// Scale returns scaled value based on screen size
// func Scale(v float64) int     { return int(v * DisplayScale()) }
// func DisplayScale() float64   { return float64(ScreenSizeY) / 100 }
// func ScreenSize() image.Point { return image.Pt(ScreenSizeX, ScreenSiz eY) }
func Scale() float64 { return float64(ScreenSizeY) / 800 } // Value of Scale() is 1 in 800 x 600
// Todo: ebiten -> general
var KeySettings = map[int][]ebiten.Key{
	4: {ebiten.KeyD, ebiten.KeyF, ebiten.KeyJ, ebiten.KeyK},
	7: {ebiten.KeyS, ebiten.KeyD, ebiten.KeyF,
		ebiten.KeySpace, ebiten.KeyJ, ebiten.KeyK, ebiten.KeyL},
}
