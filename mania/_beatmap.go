package mania

import (
	"github.com/hndada/gosu/game"
	"github.com/hndada/gosu/game/parser"
)

type Beatmap struct {
	parser.Beatmap
	Mods    game.Mods
	Keymode int
	Notes   []Note
}

func ManiaBeatmap(path string) *Beatmap {
	b, err := parser.ParseBeatmap(path)
	if err != nil {
		panic(err)
	}
	km := int(b.Difficulty["CircleSize"].(float64))
	ns, err := NewNotes(b.HitObjects, km, game.Mods{})
	if err != nil {
		panic(err)
	}
	return &Beatmap{
		Beatmap: b,
		Mods:    game.Mods{}, // todo: implement
		Keymode: km,
		Notes:   ns,
	}
}
