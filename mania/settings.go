package mania

import (
	"image/color"

	"github.com/hndada/gosu/engine/kb"
)

// const left = 30

// const (
// 	LNTailModeHead = iota
// 	LNTailModeBody
// 	LNTailModeCustom
// )

var Settings struct {
	KeyLayout    map[int][]kb.Code // todo: 무결성 검사, 겹치는거 있는지 매번 확인
	GeneralSpeed float64

	NoteWidths      map[int][4]float64 // 키마다 width 설정. 단위는 window size 대비 percent
	NoteHeigth      float64            // 두께; 키 관계없이 동일
	StagePosition   float64            // 0 ~ 100; 50 is a center
	HitPosition     float64            // HitPosition*DisplayScale() goes *the center* of Y value
	ComboPosition   float64
	JudgePosition   float64
	JudgeHeight     float64
	SpotlightColor  [4]color.RGBA
	JudgeLineHeight float64
	HPHeight        float64

	LightingScale   float64
	LightingLNScale float64

	// LineInHint        bool
	// LNHeadCustom      bool  // if false, head uses normal note image.
	// LNTailMode        uint8 // 0: Tail=Head 1: Tail=Body 2: Custom

	PlayfieldDimness float64
	ScratchMode      map[int]int
}

func init() {
	Settings.KeyLayout = map[int][]kb.Code{
		4:               {kb.CodeD, kb.CodeF, kb.CodeJ, kb.CodeK},
		5:               {kb.CodeD, kb.CodeF, kb.CodeSpacebar, kb.CodeJ, kb.CodeK},
		6:               {kb.CodeS, kb.CodeD, kb.CodeF, kb.CodeJ, kb.CodeK, kb.CodeL},
		7:               {kb.CodeS, kb.CodeD, kb.CodeF, kb.CodeSpacebar, kb.CodeJ, kb.CodeK, kb.CodeL},
		8 + LeftScratch: {kb.CodeA, kb.CodeS, kb.CodeD, kb.CodeF, kb.CodeSpacebar, kb.CodeJ, kb.CodeK, kb.CodeL},
		8:               {kb.CodeA, kb.CodeS, kb.CodeD, kb.CodeF, kb.CodeJ, kb.CodeK, kb.CodeL, kb.CodeSemicolon},
		9:               {kb.CodeA, kb.CodeS, kb.CodeD, kb.CodeF, kb.CodeSpacebar, kb.CodeJ, kb.CodeK, kb.CodeL, kb.CodeSemicolon},
		10:              {kb.CodeA, kb.CodeS, kb.CodeD, kb.CodeF, kb.CodeV, kb.CodeN, kb.CodeJ, kb.CodeK, kb.CodeL, kb.CodeSemicolon},
	}
	Settings.GeneralSpeed = 0.16

	Settings.NoteWidths = map[int][4]float64{
		4:  {10, 9, 10, 11},
		5:  {10, 9, 10, 11},
		6:  {10, 9, 10, 11},
		7:  {4.5 * 1.8, 4 * 1.8, 5 * 1.8, 5.5 * 1.8}, // {4.67, 3.83, 5.5, 5.5}
		8:  {4.5 * 1.8, 4 * 1.7, 5 * 1.7, 5.5 * 1.7},
		9:  {4.5 * 1.6, 4 * 1.6, 5 * 1.6, 5.5 * 1.6},
		10: {4.5 * 1.5, 4 * 1.5, 5 * 1.5, 5.5 * 1.5},
	}
	Settings.NoteHeigth = 3
	Settings.StagePosition = 50
	Settings.HitPosition = 85
	Settings.ComboPosition = 50
	Settings.JudgePosition = 60
	Settings.SpotlightColor = [4]color.RGBA{
		{64, 0, 0, 0xff},
		{0, 0, 64, 0xff},
		{64, 48, 0, 0xff},
		{40, 0, 40, 0xff},
	}
	Settings.JudgeHeight = 10
	Settings.JudgeLineHeight = 2
	Settings.HPHeight = 65
	Settings.LightingScale = 0.66
	Settings.LightingLNScale = 1
	Settings.PlayfieldDimness = 0.3
	Settings.ScratchMode = map[int]int{
		8: LeftScratch,
	}
}
