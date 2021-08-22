package mania

import (
	"image/color"

	"github.com/moutend/go-hook/pkg/types"
)

// const left = 30

var Settings struct {
	KeyLayout    map[int][]types.VKCode // todo: 무결성 검사, 겹치는거 있는지 매번 확인
	GeneralSpeed float64

	NoteWidths     map[int][4]float64 // 키마다 width 설정. 단위는 window size 대비 percent
	NoteHeigth     float64            // 두께; 키 관계없이 동일
	StagePosition  float64            // 0 ~ 100; 50 is a center
	HitPosition    float64            // object which is now set at 'options'
	ComboPosition  float64
	JudgePosition  float64
	JudgeHeight    float64
	SpotlightColor [4]color.RGBA
	// LineInHint        bool
	// LNHeadCustom      bool  // if false, head uses normal note image.
	// LNTailMode        uint8 // 0: Tail=Head 1: Tail=Body 2: Custom
}

const (
	LNTailModeHead = iota
	LNTailModeBody
	LNTailModeCustom
)

func init() {
	Settings.KeyLayout = map[int][]types.VKCode{
		4: {types.VK_D, types.VK_F, types.VK_J, types.VK_K},
		7: {types.VK_S, types.VK_D, types.VK_F,
			types.VK_SPACE, types.VK_J, types.VK_K, types.VK_L},
	}
	Settings.GeneralSpeed = 0.115

	Settings.NoteWidths = map[int][4]float64{
		4: {10, 9, 11, 12},
		7: {4.5 * 1.8, 4 * 1.8, 5 * 1.8, 5.5 * 1.8}, // {4.67, 3.83, 5.5, 5.5}
		8: {4.5 * 1.8, 4 * 1.8, 5 * 1.8, 5.5 * 1.8},
		9: {4.5 * 1.6, 4 * 1.6, 5 * 1.6, 5.5 * 1.6},
		10: {4.5 * 1.5, 4 * 1.5, 5 * 1.5, 5.5 * 1.5},
	}
	Settings.NoteHeigth = 3
	Settings.StagePosition = 50
	Settings.HitPosition = 85
	Settings.ComboPosition = 50
	Settings.JudgePosition = 60
	Settings.SpotlightColor = [4]color.RGBA{
		{64, 0, 0, 64},
		{0, 0, 64, 64},
		{64, 48, 0, 64},
		{40, 0, 40, 64},
	}
	Settings.JudgeHeight = 10
}
