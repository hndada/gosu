package mania

import (
	"image"
	"image/color"

	"github.com/moutend/go-hook/pkg/types"
)

// 매냐: 세로 100 기준으로 가로 얼마 만큼 쓸래
// EachDimness map[[16]byte]uint8 -> 얘는 toml 등으로 저장
// EachSpeed map[[16]byte]float64 -> 얘는 toml 등으로 저장
type SettingsTemplate struct {
	KeyLayout    map[int][]types.VKCode // todo: 무결성 검사, 겹치는거 있는지 매번 확인
	GeneralSpeed float64                // todo: fixed decimal?
	GroupSpeeds  []float64

	NoteWidths          map[int][4]float64 // 키마다 width 설정
	NoteHeigth          float64            // 두께; 키 관계없이 동일
	StagePosition       float64            // 0 ~ 100; 50 is a center
	HitPosition         float64            // object which is now set at 'options'
	ComboPosition       float64
	HitResultPosition   float64
	SpotlightColor      [4]color.RGBA
	LineInHint          bool
	LNHeadCustom        bool  // if false, head uses normal note image.
	LNTailMode          uint8 // 0: Tail=Head 1: Tail=Body 2: Custom
	SplitGap            float64
	UpsideDown          bool
	ColumnDivisionWidth float64
}

var Settings SettingsTemplate

const (
	LNTailModeHead = iota
	LNTailModeBody
	LNTailModeCustom
)

func ResetSettings() {
	Settings.KeyLayout = map[int][]types.VKCode{
		4: {types.VK_D, types.VK_F, types.VK_J, types.VK_K},
		7: {types.VK_S, types.VK_D, types.VK_F,
			types.VK_SPACE, types.VK_J, types.VK_K, types.VK_L},
	}
	Settings.GeneralSpeed = 0.115

	Settings.NoteWidths = map[int][4]float64{
		4: {10, 9, 11, 12},
		7: {4.5 * 1.8, 4 * 1.8, 5 * 1.8, 5.5 * 1.8}, // {4.67, 3.83, 5.5, 5.5}
	}
	Settings.NoteHeigth = 3
	Settings.StagePosition = 50
	Settings.HitPosition = 85
	Settings.ComboPosition = 50
	Settings.HitResultPosition = 60
	Settings.SpotlightColor = [4]color.RGBA{
		{64, 0, 0, 64},
		{0, 0, 64, 64},
		{64, 48, 0, 64},
		{40, 0, 40, 64},
	}
	Settings.LineInHint = true
	Settings.LNHeadCustom = false
	Settings.LNTailMode = LNTailModeHead
	Settings.SplitGap = 0
	Settings.UpsideDown = false
	Settings.ColumnDivisionWidth = 0
}

func StageCenter(screenSize image.Point) int {
	return int(float64(screenSize.X) * Settings.StagePosition / 100)
}
