package mania

import (
	"github.com/hndada/gosu/game"
	"github.com/hndada/gosu/game/beatmap"
)

const Mode = 3

type Beatmap struct {
	beatmap.Beatmap
	Mods    game.Mods
	Keymode int
	Notes   []Note
}

func ManiaBeatmap(path string) *Beatmap {
	b, err := beatmap.ParseBeatmap(path)
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
