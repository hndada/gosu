package main

import (
	"image"

	"github.com/hndada/gosu/input/hook"
)

var (
	MaxTPS      int64 = 1000 // 60
	ScreenSizeX int   = 800
	ScreenSizeY int   = 600

	MusicVolume float64 = 0.05
	Speed       float64 = 0.16
	KeyLayout           = map[int][]hook.Code{
		4:               {hook.CodeD, hook.CodeF, hook.CodeJ, hook.CodeK},
		5:               {hook.CodeD, hook.CodeF, hook.CodeSpacebar, hook.CodeJ, hook.CodeK},
		6:               {hook.CodeS, hook.CodeD, hook.CodeF, hook.CodeJ, hook.CodeK, hook.CodeL},
		7:               {hook.CodeS, hook.CodeD, hook.CodeF, hook.CodeSpacebar, hook.CodeJ, hook.CodeK, hook.CodeL},
		8 + LeftScratch: {hook.CodeA, hook.CodeS, hook.CodeD, hook.CodeF, hook.CodeSpacebar, hook.CodeJ, hook.CodeK, hook.CodeL},
		8:               {hook.CodeA, hook.CodeS, hook.CodeD, hook.CodeF, hook.CodeJ, hook.CodeK, hook.CodeL, hook.CodeSemicolon},
		9:               {hook.CodeA, hook.CodeS, hook.CodeD, hook.CodeF, hook.CodeSpacebar, hook.CodeJ, hook.CodeK, hook.CodeL, hook.CodeSemicolon},
		10:              {hook.CodeA, hook.CodeS, hook.CodeD, hook.CodeF, hook.CodeV, hook.CodeN, hook.CodeJ, hook.CodeK, hook.CodeL, hook.CodeSemicolon},
	}
	// NoteWidths      map[int][4]float64 // Unit: percentage comparing to screen size.
	// NoteHeigth float64 // Universal to all notes
	// HitPosition float64 // HitPosition*DisplayScale() goes *the center* of Y value
	// ComboPosition   float64
	// JudgePosition   float64
	// PlayfieldDimness float64
	// BgDimness        float64
	// ScratchMode map[int]int
)

// Scale returns scaled value based on screen size
// func Scale(v float64) int     { return int(v * DisplayScale()) }
func DisplayScale() float64   { return float64(ScreenSizeY) / 100 }
func ScreenSize() image.Point { return image.Pt(ScreenSizeX, ScreenSizeY) }
