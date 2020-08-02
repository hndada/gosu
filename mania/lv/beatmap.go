package lv

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
