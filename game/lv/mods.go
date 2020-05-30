package lv

import (
	"strings"

	"github.com/hndada/gosu/game/beatmap"
)

type Mods struct {
	Bits                           int
	TimeRate                       float64
	ScoreRate                      float64
	HPRate, CSRate, ODRate, ARRate float64
}

var ModsBitsMap = make(map[string]int)

func init() {
	// ModsBitsMap["NM"] = 0
	var bit int = 1
	for _, abbr := range [...]string{
		"NF", "EZ", "TD", "HD",
		"HR", "SD", "DT", "RL",
		"HT", "NC", "FL", "AT",
		"SO", "AP", "PF", "4K",
		"5K", "6K", "7K", "8K",
		"FI", "RD", "CM", "TP",
		"9K", "Co-op", "1K", "3K",
		"2K", "ScoreV2", "LastMod"} {
		ModsBitsMap[abbr] = bit
		bit = bit << 1
	}
	ModsBitsMap["NC"] += ModsBitsMap["DT"]
	ModsBitsMap["PF"] += ModsBitsMap["SD"]

}
func GetMods(mode int, bits int) Mods {
	mods := Mods{bits, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0}
	if mods.HasMods("EZ") {
		mods.ScoreRate *= 0.5
		mods.HPRate = 0.5
		mods.CSRate = 0.5
		mods.ODRate = 0.5
		mods.ARRate = 0.5
	}
	if mods.HasMods("NF") {
		mods.ScoreRate *= 0.5
	}
	if mods.HasMods("HT") {
		switch mode {
		case element.ModeMania:
			mods.ScoreRate *= 0.5
		default:
			mods.ScoreRate *= 0.3
		}
		mods.TimeRate = 0.75
	}
	if mods.HasMods("HR") {
		switch mode {
		case element.ModeMania:
			break
		case element.ModeCatch:
			mods.ScoreRate *= 1.12
		default:
			mods.ScoreRate *= 1.06
		}
		mods.HPRate = 1.4
		mods.CSRate = 1.3
		mods.ODRate = 1.4
		mods.ARRate = 1.4
	}
	if mods.HasMods("DT") {
		switch mode {
		case element.ModeMania:
			break
		case element.ModeCatch:
			mods.ScoreRate *= 1.06
		default:
			mods.ScoreRate *= 1.12
		}
		mods.TimeRate = 1.5
	}
	if mods.HasMods("HD") {
		switch mode {
		case element.ModeMania:
			break
		default:
			mods.ScoreRate *= 1.06
		}
	}
	if mods.HasMods("FL") {
		switch mode {
		case element.ModeMania:
			break
		default:
			mods.ScoreRate *= 1.12
		}
	}
	return mods
}

func GetModsBits(abbrs []string) int {
	bits := 0
	for _, abbr := range abbrs {
		bits += ModsBitsMap[abbr]
	}
	return bits
}

func (mods Mods) HasMods(abbr string) bool {
	return ModsBitsMap[abbr]&mods.Bits != 0
}

func GetModsName(bits int) string {
	names := make([]string, 0, 2)
	for name, bit := range ModsBitsMap {
		if bit&bits == bit { // NC has DT bit as well
			names = append(names, name)
		}
	}
	return strings.Join(names, ",")
}
