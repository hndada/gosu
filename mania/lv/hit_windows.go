package mania

import (
	"math"

	"github.com/hndada/gosu/game/beatmap"
)

var hitWindowKeys = []string{"320", "300", "200", "100", "50", "0"}
var hitWindowsBase = map[string]float64{
	"320": 16, "300": 64, "200": 97, "100": 127, "50": 151, "0": 188,
}

func (beatmap *ManiaBeatmap) SetHitWindows() {
	// I'm not very sure about changing HitWindow of Miss regarding of od
	hitWindows := make(map[string]float64)
	od := beatmap.Difficulty["OverallDifficulty"]
	odActual := math.Min(10, od*beatmap.Mods.ODRate)
	converted := beatmap.Mode == element.ModeOsu

	hitWindows["320"] = hitWindowsBase["320"]
	switch {
	case converted: // a converted map doesn't follow od parameter strictly.
		for _, k := range hitWindowKeys[1:] {
			hitWindows[k] = hitWindowsBase[k] - 3*10
		}
		if od <= 4.5 {
			hitWindows["300"] = hitWindowsBase["300"] - 17
			hitWindows["200"] = hitWindowsBase["200"] - 20
		}
	default:
		for _, k := range hitWindowKeys[1:] {
			hitWindows[k] = hitWindowsBase[k] - 3*odActual
		}
	}
	switch {
	case beatmap.Mods.HasMods("EZ"):
		for _, k := range hitWindowKeys {
			hitWindows[k] *= 7 / 5
		}
	case beatmap.Mods.HasMods("HR"):
		for _, k := range hitWindowKeys {
			hitWindows[k] *= 5 / 7
		}
	}
	// add 0.5 to hitWindow after rounding down
	for _, k := range hitWindowKeys {
		hitWindows[k] = math.Floor(hitWindows[k]) + 0.5
	}
	beatmap.HitWindows = hitWindows
}
