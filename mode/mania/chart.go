package mania

import (
	"github.com/hndada/gosu/mode"
	"github.com/hndada/rg-parser/osugame/osu"
	"log"
)

type Chart struct {
	mode.BaseChart
	Mods  Mods // mode-specific
	Keys  int
	Notes []Note
}

// 모드를 미리 받으면 안됨
// 모드마다 TimingPoint, Note건듦
func NewChartFromOsu(o *osu.Format, mods Mods) (*Chart, error) {
	var c Chart
	c.BaseChart = mode.NewBaseChartFromOsu(o)
	c.Mods = mods
	c.Keys = int(c.Parameter["Scale"])
	c.Notes = make([]Note, 0, len(o.HitObjects)*2)
	for _, ho := range o.HitObjects {
		ns, err := newNote(ho, c.Keys, c.Mods)
		if err != nil {
			log.Fatal("invalid hit object") // todo: normal log
		}
		c.Notes = append(c.Notes, ns...)
	}
	return &c, nil
}
