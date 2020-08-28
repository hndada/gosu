package config

import (
	"github.com/hajimehoshi/ebiten"
	"image"
	"image/color"
)

// todo: reset으로 다루기?
// 로컬db 만들고 나면 loading 구현
func LoadSettings() *Settings {
	// var s Settings
	// s.reset()
	// return &s
	return newSettings()
}

// todo: reset 제대로 만들기, 아니면 그냥 newSettings으로
func newSettings() *Settings {
	s := &Settings{
		maxTPS:     240,
		screenSize: image.Pt(1600, 900),
		// screenWidth:       1600,
		// screenHeight:      900,
		DimValue:          25,
		volumeMaster:      100,
		volumeBGM:         50,
		volumeSFX:         50,
		ScrollSpeed:       1.33,
		ManiaKeyLayout:    make(map[int][]ebiten.Key),
		HitPosition:       70,
		ComboPosition:     50,
		HitResultPosition: 60,
		noteWidths:        make(map[int][4]float64),
		noteHeigth:        3,
		lnHeadCustom:      false,
		lnTailMode:        0,
		SpotlightColor: [4]color.RGBA{
			{64, 0, 0, 64},
			{0, 0, 64, 64},
			{64, 48, 0, 64},
			{40, 0, 40, 64},
		},
		lineInJudgeLine:     true,
		SplitGap:            0,
		UpsideDown:          false,
		StagePosition:       0,
		ColumnDivisionWidth: 0,
	}
	s.ManiaKeyLayout[4] = []ebiten.Key{
		ebiten.KeyD, ebiten.KeyF, ebiten.KeyJ, ebiten.KeyK,
	}
	s.ManiaKeyLayout[7] = []ebiten.Key{
		ebiten.KeyS, ebiten.KeyD, ebiten.KeyF,
		ebiten.KeySpace, ebiten.KeyJ, ebiten.KeyK, ebiten.KeyL,
	}

	s.noteWidths[4] = [4]float64{10, 9, 11, 12}
	s.noteWidths[7] = [4]float64{4.5, 4, 5, 5.5} // [4]float64{4.67, 3.83, 5.5, 5.5}
	return s
}