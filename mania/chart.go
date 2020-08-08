package mania

import (
	"github.com/hndada/gosu/game"
	"github.com/hndada/gosu/game/parser/osu"
)

type Chart struct {
	game.BaseChart
	Mods  game.Mods
	Keys  int
	Notes []Note
}

func NewChart(path string) *Chart { // todo: 더 neat하게 input
	var o osu.OSU // todo:should be pointer
	var err error
	o, err = osu.NewOSU(path)
	if err != nil {
		panic(err)
	}

	var c Chart
	c.BaseChart=game.NewBaseChart(&o)
	c.Mods = game.Mods{} // todo: implement
	c.Keys = int(c.Parameter["Scale"])
	c.Notes, err = NewNotes(o.HitObjects, c.Keys)
	if err != nil {
		panic(err)
	}
	return &c
}
