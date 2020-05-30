package lv

import (
	"errors"

	"github.com/hndada/gosu/game/beatmap"
	"github.com/hndada/gosu/game/tools"
)

var ErrMode = errors.New("invalid mode input")

type Beatmap interface {
	SetBase(path string, modsBits int)

	AddNotes()
	SortNotes()
	SetHitWindows()
	SetCurves()
	CalcStrain()
	CalcStamina()
	// CalcLegibility()

	AddOldNotes()
	CalcOldStrain()
	InitOldStrainPeak(i int, sectionEndTime int, timeRate float64) float64
	GetOldStrain(i int) float64
}

func CheckMode(mapMode, runMode int) {
	switch mapMode {
	case runMode, element.ModeOsu:
		return
	}
	panic(&tools.ValError{"Mode", tools.Itoa(mapMode), ErrMode})
}
