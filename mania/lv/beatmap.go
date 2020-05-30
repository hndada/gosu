package mania

import (
	"github.com/hndada/gosu/game/beatmap"
	"github.com/hndada/gosu/game/lv"
	"github.com/hndada/gosu/game/tools"
)

const modeMania = element.ModeMania

type ManiaBeatmap struct {
	diff.Beatmap
	Keymode       int
	Notes         []ManiaNote
	OldNotes      []ManiaOldNote
	maxChunkDelta int // max distance that yields chord penalty
	maxTrillDelta int // max distance that yields trill bonus
	maxJackDelta  int // max distance that yields jack bonus
}

// raw is required not only for mode, getting metadata but for converting as well
func (beatmap *ManiaBeatmap) SetBase(path string, modsBits int) {
	beatmap.Beatmap = element.ParseBeatmap(path)
	diff.CheckMode(beatmap.Mode, modeMania)
	beatmap.Mods = diff.GetMods(modeMania, modsBits)
	beatmap.Keymode = beatmap.getKeymode()
}

func (beatmap ManiaBeatmap) getKeymode() int {
	cs := beatmap.Difficulty["CircleSize"]
	switch beatmap.Mode {
	case element.ModeOsu:
		// originally, default keymode of converted map is affected by
		// ratio of slider and spinner and od along with cs; mostly it is 7 after all.
		return 7
	case modeMania:
		// dual mode does not fully supported; only former half goes processed.
		if cs > 10 {
			// cs /= 2
			panic(&tools.ValError{"Keymode", tools.Ftoa(cs), diff.ErrMode})
		}
		return int(cs)
	default:
		panic(&tools.ValError{"Mode", tools.Itoa(beatmap.Mode), tools.ErrFlow})
	}
}
