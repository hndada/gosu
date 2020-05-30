package taiko

import (
	"math"
)

var hitWindowKeys = []string{"300", "100", "0"}
var hitWindowsBase = map[string]float64{
	"300": 49, "100": 119, "0": 188, "100_5": 79,
}

func (beatmap *TaikoBeatmap) SetHitWindows() {
	// I'm not very sure about changing HitWindow of Miss regarding of od
	hitWindows := make(map[string]float64)
	od := beatmap.Difficulty["OverallDifficulty"]
	odActual := math.Min(10, od*beatmap.Mods.ODRate)

	for _, k := range hitWindowKeys {
		hitWindows[k] = hitWindowsBase[k] - 3*odActual
	}
	if odActual >= 5 {
		hitWindows["100"] = hitWindowsBase["100_5"] - (odActual-5)*6
	} else {
		hitWindows["100"] = hitWindowsBase["100"] - odActual*8
	}

	for _, k := range hitWindowKeys {
		hitWindows[k] = math.Floor(hitWindows[k]) + 0.5
	}

	switch {
	case beatmap.Mods.HasMods("HT"):
		for _, k := range hitWindowKeys {
			hitWindows[k] *= 4 / 3
			hitWindows[k] += 2 / 3
		}
	case beatmap.Mods.HasMods("DT"):
		for _, k := range hitWindowKeys {
			hitWindows[k] *= 2 / 3
			hitWindows[k] += 1 / 3
		}
	}
	beatmap.HitWindows = hitWindows
}
