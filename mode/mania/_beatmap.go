package mania

import (
	"github.com/hndada/gosu/mode"
	"github.com/hndada/gosu/parser"
)

type Beatmap struct {
	parser.Beatmap
	Mods    mode.Mods
	Keymode int
	Notes   []Note
}

func ManiaBeatmap(path string) *Beatmap {
	b, err := parser.ParseBeatmap(path)
	if err != nil {
		panic(err)
	}
	km := int(b.Difficulty["CircleSize"].(float64))
	ns, err := NewNotes(b.HitObjects, km, mode.Mods{})
	if err != nil {
		panic(err)
	}
	return &Beatmap{
		Beatmap: b,
		Mods:    mode.Mods{}, // todo: implement
		Keymode: km,
		Notes:   ns,
	}
}
